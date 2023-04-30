package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("Roman", "abcabc@gmail.com")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Nickname, "Roman")
	assert.Equal(t, response.Data.Email, "abcabc@gmail.com")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("Roman", "abcabc@gmail.com")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, response.Data.ID, "Pasha", "kekkekkek@yandex.ru")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "Pasha")
	assert.Equal(t, response.Data.Email, "kekkekkek@yandex.ru")
}

func TestUpdateAnotherUser(t *testing.T) {
	client := getTestClient()

	roman, err := client.createUser("Roman the CEO", "CE0@tinkoff.ru")
	assert.NoError(t, err)
	oliver, err := client.createUser("Oliver the CEO", "CEO@tinkoff.ru")
	assert.NoError(t, err)

	_, err = client.updateUser(oliver.Data.ID, roman.Data.ID, "Roman the false CEO", "ceo@tinkoff.ru")
	assert.ErrorIs(t, err, ErrForbidden)
}
