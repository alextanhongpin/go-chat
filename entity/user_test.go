package entity_test

import (
	"testing"

	"github.com/alextanhongpin/go-chat/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserPassword(t *testing.T) {
	assert := assert.New(t)
	user := new(entity.User)
	password := "hello world"
	err := user.SetPassword(password)
	assert.Nil(err)

	err = user.ComparePassword(password)
	assert.Nil(err)

	err = user.ComparePassword("random password")
	assert.NotNil(err)
	assert.Equal("password do not match", err.Error())
}
