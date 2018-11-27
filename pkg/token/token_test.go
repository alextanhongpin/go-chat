package token_test

import (
	"testing"
	"time"

	"github.com/alextanhongpin/go-chat/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestSigning(t *testing.T) {
	assert := assert.New(t)
	id := "abc123"
	signer := token.New(token.SignerOptions{
		Now: func() time.Time {
			return time.Now().UTC()
		},
		TTL:    2 * time.Hour,
		Issuer: "go-openid",
		Secret: []byte("secret"),
	})
	token, err := signer.Sign(id)
	assert.Nil(err)

	usrID, err := signer.Verify(token)
	assert.Nil(err)
	assert.Equal(id, usrID)
}
