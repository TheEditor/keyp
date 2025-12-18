package server

import (
	"encoding/json"
	"time"

	"github.com/TheEditor/keyp/internal/model"
)

// Response envelope for all API responses
type Response struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error *ErrorDetail    `json:"error,omitempty"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error codes
const (
	ErrCodeBadRequest      = "bad_request"
	ErrCodeUnauthorized    = "unauthorized"
	ErrCodeNotFound        = "not_found"
	ErrCodeConflict        = "conflict"
	ErrCodeInternalError   = "internal_error"
)

// SuccessResponse creates a success response envelope
func SuccessResponse(data interface{}) *Response {
	jsonData, _ := json.Marshal(data)
	return &Response{
		OK:   true,
		Data: jsonData,
	}
}

// ErrorResponse creates an error response envelope
func ErrorResponse(code, message string) *Response {
	return &Response{
		OK: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// Session represents an authenticated session
type Session struct {
	Token     string
	Handle    interface{} // *vault.VaultHandle
	CreatedAt time.Time
	ExpiresAt time.Time
}

// UnlockRequest for POST /v1/unlock
type UnlockRequest struct {
	Password string `json:"password"`
}

// UnlockResponse for successful unlock
type UnlockResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshResponse for POST /v1/refresh
type RefreshResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
}

// SecretListItem for list responses (minimal info)
type SecretListItem struct {
	Name      string    `json:"name"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SecretDetail for get responses (full info, redacted by default)
type SecretDetail struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Tags      []string   `json:"tags"`
	Fields    []Field    `json:"fields"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Field in secret response
type Field struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
}

// FieldInput for request body
type FieldInput struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type,omitempty"`
}

// CreateSecretRequest for POST /v1/secrets
type CreateSecretRequest struct {
	Name   string       `json:"name"`
	Tags   []string     `json:"tags,omitempty"`
	Fields []FieldInput `json:"fields"`
	Notes  string       `json:"notes,omitempty"`
}

// UpdateSecretRequest for PUT /v1/secrets/:name (partial updates)
type UpdateSecretRequest struct {
	Tags   *[]string    `json:"tags,omitempty"`
	Fields *[]FieldInput `json:"fields,omitempty"`
	Notes  *string      `json:"notes,omitempty"`
}

// ToSecretListItem converts model to API type
func ToSecretListItem(s *model.SecretObject) SecretListItem {
	return SecretListItem{
		Name:      s.Name,
		Tags:      s.Tags,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSecretDetail converts model to API type with optional redaction
func ToSecretDetail(s *model.SecretObject, redact bool) SecretDetail {
	if redact {
		s = s.Redacted()
	}

	fields := make([]Field, len(s.Fields))
	for i, f := range s.Fields {
		fields[i] = Field{
			Label:     f.Label,
			Value:     f.Value,
			Sensitive: f.Sensitive,
			Type:      f.Type,
		}
	}

	return SecretDetail{
		ID:        s.ID,
		Name:      s.Name,
		Tags:      s.Tags,
		Fields:    fields,
		Notes:     s.Notes,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSecretObject converts request to model type
func (r *CreateSecretRequest) ToSecretObject() *model.SecretObject {
	secret := model.NewSecretObject(r.Name)
	secret.Tags = r.Tags
	secret.Notes = r.Notes

	for _, f := range r.Fields {
		field := model.Field{
			ID:        model.NewField(f.Label, f.Value).ID,
			Label:     f.Label,
			Value:     f.Value,
			Sensitive: f.Sensitive,
			Type:      f.Type,
			SortOrder: len(secret.Fields),
		}
		if field.Type == "" {
			field.Type = model.FieldTypeText
		}
		secret.Fields = append(secret.Fields, field)
	}

	return secret
}

// HealthResponse for GET /health
type HealthResponse struct {
	Status string `json:"status"`
}

// VersionResponse for GET /version
type VersionResponse struct {
	Version string `json:"version"`
	Git     string `json:"git,omitempty"`
}
