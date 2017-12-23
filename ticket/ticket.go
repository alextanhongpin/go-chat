package ticket

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	secret = []byte("secret")
	issuer = "go-chat"
)

// Ticket represents the token required for the chat
type Ticket struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// Sign returns an encrypted ticket string
func Sign(tic Ticket) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tic)
	return t.SignedString(secret)
}

// Verify the ticket and check for expiry
func Verify(ticket string) (*Ticket, error) {
	token, err := jwt.ParseWithClaims(ticket, &Ticket{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Ticket); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// New returns a new pointer to the Ticket
func New(id string, duration time.Duration) Ticket {
	return Ticket{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    issuer,
		},
	}
}
