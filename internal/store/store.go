package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/TheEditor/keyp/internal/model"
)

// SearchOptions holds filtering options for search queries
type SearchOptions struct {
	Tags  []string
	Limit int
}

// Store handles SQLite database operations
type Store struct {
	db *sql.DB
}

// Open opens or creates a SQLite database
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// SetMeta stores a metadata key-value pair
func (s *Store) SetMeta(key, value string) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO vault_meta (key, value) VALUES (?, ?)",
		key, value,
	)
	return err
}

// GetMeta retrieves a metadata value by key
func (s *Store) GetMeta(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM vault_meta WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return value, err
}

func (s *Store) initSchema() error {
	schema := `
    CREATE TABLE IF NOT EXISTS vault_meta (
        key TEXT PRIMARY KEY,
        value TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS secrets (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        tags TEXT DEFAULT '[]',
        notes TEXT DEFAULT '',
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS fields (
        id TEXT PRIMARY KEY,
        secret_id TEXT NOT NULL,
        label TEXT NOT NULL,
        value TEXT NOT NULL,
        sensitive INTEGER DEFAULT 1,
        type TEXT DEFAULT 'text',
        sort_order INTEGER DEFAULT 0,
        FOREIGN KEY (secret_id) REFERENCES secrets(id) ON DELETE CASCADE,
        UNIQUE(secret_id, label)
    );

    CREATE VIRTUAL TABLE IF NOT EXISTS secrets_fts USING fts5(
        name, tags, notes, field_labels
    );

    CREATE TRIGGER IF NOT EXISTS secrets_ai AFTER INSERT ON secrets BEGIN
        INSERT INTO secrets_fts(rowid, name, tags, notes, field_labels)
        VALUES (new.rowid, new.name, new.tags, new.notes, '');
    END;

    CREATE TRIGGER IF NOT EXISTS secrets_ad AFTER DELETE ON secrets BEGIN
        INSERT INTO secrets_fts(secrets_fts, rowid, name, tags, notes, field_labels)
        VALUES ('delete', old.rowid, old.name, old.tags, old.notes, '');
    END;

    CREATE TRIGGER IF NOT EXISTS secrets_au AFTER UPDATE ON secrets BEGIN
        INSERT INTO secrets_fts(secrets_fts, rowid, name, tags, notes, field_labels)
        VALUES ('delete', old.rowid, old.name, old.tags, old.notes, '');
        INSERT INTO secrets_fts(rowid, name, tags, notes, field_labels)
        VALUES (new.rowid, new.name, new.tags, new.notes, '');
    END;
    `
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	// Migrate existing data: rebuild FTS index with field labels
	return s.rebuildFTSIndex()
}

// rebuildFTSIndex rebuilds the FTS5 index with field labels for all secrets
func (s *Store) rebuildFTSIndex() error {
	// Delete all FTS entries
	_, err := s.db.Exec("DELETE FROM secrets_fts")
	if err != nil {
		return err
	}

	// Get all secrets
	rows, err := s.db.Query("SELECT rowid, id FROM secrets")
	if err != nil {
		return err
	}
	defer rows.Close()

	var secrets []struct {
		rowid int64
		id    string
	}
	for rows.Next() {
		var rowid int64
		var id string
		if err := rows.Scan(&rowid, &id); err != nil {
			return err
		}
		secrets = append(secrets, struct {
			rowid int64
			id    string
		}{rowid, id})
	}

	// For each secret, fetch data and rebuild FTS entry
	for _, secret := range secrets {
		var name, tags, notes string
		err := s.db.QueryRow(
			"SELECT name, tags, notes FROM secrets WHERE id = ?",
			secret.id,
		).Scan(&name, &tags, &notes)
		if err != nil {
			return err
		}

		// Get field labels for this secret
		fieldLabels, err := s.getFieldLabelsForSecret(secret.id)
		if err != nil {
			return err
		}

		// Insert FTS entry with field labels
		_, err = s.db.Exec(
			"INSERT INTO secrets_fts(rowid, name, tags, notes, field_labels) VALUES (?, ?, ?, ?, ?)",
			secret.rowid, name, tags, notes, fieldLabels,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// getFieldLabelsForSecret returns space-separated field labels for a secret
func (s *Store) getFieldLabelsForSecret(secretID string) (string, error) {
	rows, err := s.db.Query(
		"SELECT label FROM fields WHERE secret_id = ? ORDER BY sort_order",
		secretID,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var labels []string
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return "", err
		}
		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	// Join labels with spaces for FTS
	var result string
	for _, label := range labels {
		if result != "" {
			result += " "
		}
		result += label
	}
	return result, nil
}

// Create inserts a new secret with its fields
func (s *Store) Create(ctx context.Context, secret *model.SecretObject) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO secrets (id, name, tags, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		secret.ID, secret.Name, secret.TagsJSON(), secret.Notes,
		secret.CreatedAt.Format(time.RFC3339),
		secret.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	for _, f := range secret.Fields {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO fields (id, secret_id, label, value, sensitive, type, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)",
			f.ID, secret.ID, f.Label, f.Value, boolToInt(f.Sensitive), f.Type, f.SortOrder,
		)
		if err != nil {
			return err
		}
	}

	// Rebuild FTS entry with field labels after fields are inserted
	// Get the rowid of the inserted secret
	var rowid int64
	err = tx.QueryRowContext(ctx, "SELECT rowid FROM secrets WHERE id = ?", secret.ID).Scan(&rowid)
	if err != nil {
		return err
	}

	// Delete old FTS entry if it exists
	_, err = tx.ExecContext(ctx,
		"INSERT INTO secrets_fts(secrets_fts, rowid, name, tags, notes, field_labels) VALUES ('delete', ?, '', '', '', '')",
		rowid,
	)
	if err != nil {
		return err
	}

	// Insert new FTS entry with field labels
	fieldLabels := buildFieldLabelsForSecret(secret)
	_, err = tx.ExecContext(ctx,
		"INSERT INTO secrets_fts(rowid, name, tags, notes, field_labels) VALUES (?, ?, ?, ?, ?)",
		rowid, secret.Name, secret.TagsJSON(), secret.Notes, fieldLabels,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetByName retrieves a secret by name
func (s *Store) GetByName(ctx context.Context, name string) (*model.SecretObject, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, name, tags, notes, created_at, updated_at FROM secrets WHERE name = ?",
		name,
	)

	secret, err := s.scanSecret(row)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	fields, err := s.getFields(ctx, secret.ID)
	if err != nil {
		return nil, err
	}
	secret.Fields = fields

	return secret, nil
}

// List returns all secrets with optional filtering
func (s *Store) List(ctx context.Context, opts *SearchOptions) ([]*model.SecretObject, error) {
	query := "SELECT id, name, tags, notes, created_at, updated_at FROM secrets"
	args := []interface{}{}

	// Apply tag filtering if specified
	if opts != nil && len(opts.Tags) > 0 {
		query += " WHERE " + buildTagFilter(opts.Tags, &args)
	}

	query += " ORDER BY name"

	// Apply limit if specified
	if opts != nil && opts.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, opts.Limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*model.SecretObject
	for rows.Next() {
		secret, err := s.scanSecretRows(rows)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, rows.Err()
}

// Search performs full-text search across name, tags, notes
func (s *Store) Search(ctx context.Context, query string, opts *SearchOptions) ([]*model.SecretObject, error) {
	sqlQuery := `
        SELECT s.id, s.name, s.tags, s.notes, s.created_at, s.updated_at
        FROM secrets s
        JOIN secrets_fts fts ON s.rowid = fts.rowid
        WHERE secrets_fts MATCH ?`
	args := []interface{}{query}

	// Apply tag filtering if specified
	if opts != nil && len(opts.Tags) > 0 {
		sqlQuery += " AND " + buildTagFilter(opts.Tags, &args)
	}

	sqlQuery += " ORDER BY rank"

	// Apply limit if specified
	if opts != nil && opts.Limit > 0 {
		sqlQuery += " LIMIT ?"
		args = append(args, opts.Limit)
	}

	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*model.SecretObject
	for rows.Next() {
		secret, err := s.scanSecretRows(rows)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, rows.Err()
}

// Update modifies an existing secret
func (s *Store) Update(ctx context.Context, secret *model.SecretObject) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the main secret record
	secret.UpdatedAt = time.Now()
	result, err := tx.ExecContext(ctx,
		"UPDATE secrets SET name = ?, tags = ?, notes = ?, updated_at = ? WHERE id = ?",
		secret.Name, secret.TagsJSON(), secret.Notes,
		secret.UpdatedAt.Format(time.RFC3339),
		secret.ID,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return ErrNotFound
	}

	// Delete old fields for this secret
	_, err = tx.ExecContext(ctx, "DELETE FROM fields WHERE secret_id = ?", secret.ID)
	if err != nil {
		return err
	}

	// Insert updated fields
	for _, f := range secret.Fields {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO fields (id, secret_id, label, value, sensitive, type, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)",
			f.ID, secret.ID, f.Label, f.Value, boolToInt(f.Sensitive), f.Type, f.SortOrder,
		)
		if err != nil {
			return err
		}
	}

	// Rebuild FTS entry with field labels after fields are inserted
	// Get the rowid of the secret
	var rowid int64
	err = tx.QueryRowContext(ctx, "SELECT rowid FROM secrets WHERE id = ?", secret.ID).Scan(&rowid)
	if err != nil {
		return err
	}

	// Delete old FTS entry
	_, err = tx.ExecContext(ctx,
		"INSERT INTO secrets_fts(secrets_fts, rowid, name, tags, notes, field_labels) VALUES ('delete', ?, '', '', '', '')",
		rowid,
	)
	if err != nil {
		return err
	}

	// Insert new FTS entry with field labels
	fieldLabels := buildFieldLabelsForSecret(secret)
	_, err = tx.ExecContext(ctx,
		"INSERT INTO secrets_fts(rowid, name, tags, notes, field_labels) VALUES (?, ?, ?, ?, ?)",
		rowid, secret.Name, secret.TagsJSON(), secret.Notes, fieldLabels,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// buildFieldLabelsForSecret creates space-separated field labels for FTS indexing
func buildFieldLabelsForSecret(secret *model.SecretObject) string {
	var result string
	for i, field := range secret.Fields {
		if i > 0 {
			result += " "
		}
		result += field.Label
	}
	return result
}

// Delete removes a secret and its fields
func (s *Store) Delete(ctx context.Context, name string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM secrets WHERE name = ?", name)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) getFields(ctx context.Context, secretID string) ([]model.Field, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, label, value, sensitive, type, sort_order FROM fields WHERE secret_id = ? ORDER BY sort_order",
		secretID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []model.Field
	for rows.Next() {
		var f model.Field
		var sensitive int
		err := rows.Scan(&f.ID, &f.Label, &f.Value, &sensitive, &f.Type, &f.SortOrder)
		if err != nil {
			return nil, err
		}
		f.Sensitive = sensitive == 1
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (s *Store) scanSecret(row *sql.Row) (*model.SecretObject, error) {
	var secret model.SecretObject
	var tagsJSON, createdAt, updatedAt string

	err := row.Scan(&secret.ID, &secret.Name, &tagsJSON, &secret.Notes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	secret.Tags = model.ParseTags(tagsJSON)
	secret.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	secret.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &secret, nil
}

func (s *Store) scanSecretRows(rows *sql.Rows) (*model.SecretObject, error) {
	var secret model.SecretObject
	var tagsJSON, createdAt, updatedAt string

	err := rows.Scan(&secret.ID, &secret.Name, &tagsJSON, &secret.Notes, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	secret.Tags = model.ParseTags(tagsJSON)
	secret.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	secret.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &secret, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// buildTagFilter constructs a WHERE clause for tag filtering
// Tags are stored as JSON arrays, so we use JSON functions to check membership
func buildTagFilter(tags []string, args *[]interface{}) string {
	if len(tags) == 0 {
		return "1=1"
	}

	// Build a condition that checks if any of the specified tags exist in the tags JSON array
	conditions := make([]string, len(tags))
	for i, tag := range tags {
		conditions[i] = "json_extract(tags, '$') LIKE ?"
		*args = append(*args, "%\""+tag+"\"%")
	}

	// Use OR to match any of the tags
	if len(conditions) == 1 {
		return conditions[0]
	}

	filter := "("
	for i, cond := range conditions {
		if i > 0 {
			filter += " OR "
		}
		filter += cond
	}
	filter += ")"
	return filter
}
