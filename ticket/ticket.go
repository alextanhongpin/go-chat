package ticket

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Ticket represents the token required for the chat
type Ticket struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type Machine struct {
	secret   []byte
	issuer   string
	duration time.Duration
}

type Dispenser interface {
	Sign(t Ticket) (string, error)
	Verify(string) (string, error)
	New(id string) Ticket
}

func NewMachine(secret []byte, issuer string, duration time.Duration) Machine {
	return Machine{
		secret:   secret,
		issuer:   issuer,
		duration: duration,
	}
}

func (t Machine) Sign(ticket Ticket) (string, error) {
	return sign(t.secret, ticket)
}

func (t Machine) Verify(rawTicket string) (string, error) {
	claims, err := verify(t.secret, rawTicket)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

func (t Machine) New(id string) Ticket {
	return new(t.issuer, id, t.duration)
}

// Sign returns an encrypted ticket string
func sign(secret []byte, tic Ticket) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tic)
	return t.SignedString(secret)
}

// Verify the ticket and check for expiry
func verify(secret []byte, ticket string) (*Ticket, error) {
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
func new(issuer, id string, duration time.Duration) Ticket {
	return Ticket{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    issuer,
		},
	}
}
