package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()
	u1, err := client.createUser("Vasya", "vasya123@mail.ru")
	assert.NoError(t, err)
	u2, err := client.createUser("Vova", "vova555@mail.ru")
	assert.NoError(t, err)

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(u2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	u1, _ := client.createUser("Petya", "pe@mail.com")
	u2, _ := client.createUser("Vasya", "va@mail.com")

	resp, err := client.createAd(u1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(u2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()
	u, err := client.createUser("Fedya", "fedya@mail.ru")
	assert.NoError(t, err)

	resp, err := client.createAd(u.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(u.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(u.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
