package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SecretObject represents a structured secret with multiple fields
type SecretObject struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tags      []string  `json:"tags"`
	Fields    []Field   `json:"fields,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Field represents a single named value within a secret
type Field struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
	SortOrder int    `json:"sort_order"`
}

// FieldType constants for UI hints
const (
	FieldTypeText     = "text"
	FieldTypePassword = "password"
	FieldTypePIN      = "pin"
	FieldTypeURL      = "url"
	FieldTypeEmail    = "email"
)

// NewSecretObject creates a new secret with defaults
func NewSecretObject(name string) *SecretObject {
	now := time.Now()
	return &SecretObject{
		ID:        uuid.New().String(),
		Name:      name,
		Tags:      []string{},
		Fields:    []Field{},
		Notes:     "",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewField creates a new field with defaults
func NewField(label, value string) Field {
	return Field{
		ID:        uuid.New().String(),
		Label:     label,
		Value:     value,
		Sensitive: true,
		Type:      FieldTypeText,
		SortOrder: 0,
	}
}

// AddField appends a field to the secret
func (s *SecretObject) AddField(f Field) {
	f.SortOrder = len(s.Fields)
	s.Fields = append(s.Fields, f)
	s.UpdatedAt = time.Now()
}

// TagsJSON returns tags as JSON string for storage
func (s *SecretObject) TagsJSON() string {
	data, _ := json.Marshal(s.Tags)
	return string(data)
}

// ParseTags parses JSON string into tags slice
func ParseTags(jsonStr string) []string {
	var tags []string
	if jsonStr == "" {
		return []string{}
	}
	json.Unmarshal([]byte(jsonStr), &tags)
	return tags
}

const RedactedValue = "********"

// Redacted returns a copy of the secret with sensitive field values masked
func (s *SecretObject) Redacted() *SecretObject {
	copy := *s
	copy.Fields = make([]Field, len(s.Fields))
	for i, f := range s.Fields {
		copy.Fields[i] = f
		if f.Sensitive {
			copy.Fields[i].Value = RedactedValue
		}
	}
	return &copy
}
