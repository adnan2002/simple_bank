package token

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofrs/uuid"
)


// ClaimStrings is just a slice of strings (used for audiences in JWTs).
type ClaimStrings []string

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

// =======================
// JWT Claims-like methods
// =======================

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	if p.ExpiredAt.IsZero() {
		return nil, errors.New("expiration time not set")
	}
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	if p.IssuedAt.IsZero() {
		return nil, errors.New("issued at not set")
	}
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (p *Payload) GetSubject() (string, error) {
	// Typically "subject" is the user ID, here we map it to Username
	if p.Username == "" {
		return "", errors.New("subject not set")
	}
	return p.Username, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
// Valid checks if the token payload is valid (not expired, required fields set, etc.).
func (p *Payload) Valid() error {
    if time.Now().After(p.ExpiredAt) {
        return errors.New("token has expired")
    }

    if p.Username == "" {
        return errors.New("username (subject) not set")
    }

    return nil
}
