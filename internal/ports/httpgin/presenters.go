package httpgin

import (
	"github.com/gin-gonic/gin"
	validation "github.com/unicoooorn/tag_validation"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title" validate:"between:1,100"`
	Text   string `json:"text" validate:"between:1,500"`
	UserID int64  `json:"user_id"`
}

func (c createAdRequest) Validate() error {
	return validation.Validate(c)
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type adResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

type findAdsRequest struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title" validate:"between:1,100"`
	Text   string `json:"text" validate:"between:1,500"`
	UserID int64  `json:"user_id"`
}

type updateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserID   int64  `json:"user_id"`
}

func (u updateAdRequest) Validate() error {
	return validation.Validate(u)
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreatedTime:  ad.Created,
			ModifiedTime: ad.Modified,
		},
		"error": nil,
	}
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func MultipleAdsSuccessResponse(ads []ads.Ad) *gin.H {
	multipleAdsResponse := make([]adResponse, 0)
	for _, ad := range ads {
		multipleAdsResponse = append(multipleAdsResponse, adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreatedTime:  ad.Created,
			ModifiedTime: ad.Modified,
		})
	}
	return &gin.H{
		"data":  multipleAdsResponse,
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
