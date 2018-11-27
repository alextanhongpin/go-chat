package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Signer signs a user id into token and vice-versa.
type Signer interface {
	Sign(id string) (string, error)
	Verify(token string) (string, error)
	ExpiresIn() int64
}

// SignerOptions represents the options for the SignerImpl.
type SignerOptions struct {
	Now    func() time.Time
	TTL    time.Duration
	Issuer string
	Secret []byte
}

// SignerImpl implements the Signer interface.
type SignerImpl struct {
	opts SignerOptions
}

// New returns a new Signer with the given options.
func New(opts SignerOptions) *SignerImpl {
	return &SignerImpl{opts}
}

// ExpiresIn returns the lifespan of the token in seconds.
func (s *SignerImpl) ExpiresIn() int64 {
	return int64(s.opts.TTL.Seconds())
}

// Sign signs the user as a subject and returns a token.
func (s *SignerImpl) Sign(id string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: s.opts.Now().Add(s.opts.TTL).Unix(),
		Issuer:    s.opts.Issuer,
		Subject:   id,
	}
	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return signer.SignedString(s.opts.Secret)
}

// Verify validates the token.
func (s *SignerImpl) Verify(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.opts.Secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["iss"] != s.opts.Issuer {
			return "", errors.New("invalid token")
		}
		sub, _ := claims["sub"].(string)
		return sub, nil
	}
	return "", err
}
