package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFilteredSearch(t *testing.T) {
	client := getTestClient()
	petya, err := client.createUser("Petya", "petya@yahoo.com")
	assert.NoError(t, err)
	vasya, err := client.createUser("vasya", "vasya@yahoo.com")
	assert.NoError(t, err)

	response, err := client.createAd(petya.Data.ID, "Cats", "are cute")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(petya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(vasya.Data.ID, "Dogs", "are funny")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(vasya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(vasya.Data.ID, "Rats", "are funny too")
	assert.NoError(t, err)
	publishedAd, err := client.changeAdStatus(vasya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	// time filtering
	f := findAdsRequest{
		CreatedTime: publishedAd.Data.Created,
	}
	ads, err := client.listAdsByFilter(f)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)

	// author ID filtering
	f = findAdsRequest{
		AuthorID: vasya.Data.ID,
	}
	ads, err = client.listAdsByFilter(f)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	assert.Equal(t, vasya.Data.ID, ads.Data[0].AuthorID)
	assert.Equal(t, vasya.Data.ID, ads.Data[1].AuthorID)
}

func TestGetByID(t *testing.T) {
	client := getTestClient()
	u, err := client.createUser("Petya", "petya@yahoo.com")
	assert.NoError(t, err)

	response, err := client.createAd(u.Data.ID, "Cats", "are cute")
	assert.NoError(t, err)
	publishedAd, err := client.changeAdStatus(u.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.getAdById(publishedAd.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, ads.Data.ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data.Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data.Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data.AuthorID, publishedAd.Data.AuthorID)
}

func TestGetByTitle(t *testing.T) {
	client := getTestClient()
	petya, err := client.createUser("Petya", "petya@yahoo.com")
	assert.NoError(t, err)
	vasya, err := client.createUser("Vasya", "vasya@yahoo.com")
	assert.NoError(t, err)

	response, err := client.createAd(vasya.Data.ID, "Rats", "are funny")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(vasya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)
	response, err = client.createAd(petya.Data.ID, "Cats", "are cute")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(petya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)
	response, err = client.createAd(vasya.Data.ID, "Cats", "are funny")
	assert.NoError(t, err)
	_, err = client.changeAdStatus(vasya.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.getAdsByTitle(response.Data.Title)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[1].Title, response.Data.Title)
}

func (tc *testClient) listAdsByFilter(f findAdsRequest) (adsResponse, error) {
	body := map[string]any{
		"id":            f.ID,
		"title":         f.Title,
		"text":          f.Text,
		"author_id":     f.AuthorID,
		"published":     f.Published,
		"created_time":  f.CreatedTime,
		"modified_time": f.ModifiedTime,
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/search", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}
