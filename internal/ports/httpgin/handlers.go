package httpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"homework9/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		_, err = a.GetUserByID(c, reqBody.UserID)
		if err == app.ErrNotFound {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		if err := reqBody.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(u))
	}
}

// метод для получения объявления
func getPublishedAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, err := a.ListPublishedAds(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, MultipleAdsSuccessResponse(l))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		_, err := a.GetUserByID(c, reqBody.UserID)
		if err == app.ErrNotFound {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.UpdateStatusById(c, int64(adID), reqBody.Published, reqBody.UserID)
		if err == app.ErrAccessDenied {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(u))
	}
}

func getAdById(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		l, err := a.GetAdsByFilter(c, app.FilterOpts{ID: int64(adID), Hidden: false})
		if err != nil {
			c.JSON(http.StatusNotFound, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&l[0]))
	}
}

func getAdByTitle(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		l, err := a.GetAdsByFilter(c, app.FilterOpts{Title: title})
		if err == nil && len(l) == 0 {
			c.JSON(http.StatusNotFound, AdErrorResponse(nil))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, MultipleAdsSuccessResponse(l))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		_, err := a.GetUserByID(c, reqBody.UserID)
		if err == app.ErrNotFound {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		if err := reqBody.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.UpdateAdById(c, int64(adID), reqBody.Title, reqBody.Text, reqBody.UserID)
		if err == app.ErrAccessDenied {
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(u))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		u, err := a.CreateUser(c, reqBody.Nickname, reqBody.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		_, err := a.GetUserByID(c, reqBody.UserID)
		if err == app.ErrNotFound {
			c.JSON(http.StatusNotFound, AdErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		u, err := a.UpdateUserByID(c, int64(userID), reqBody.Nickname, reqBody.Email, reqBody.UserID)
		if err == app.ErrAccessDenied {
			c.JSON(http.StatusForbidden, UserErrorResponse(err))
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

func getAdsByFilter(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody findAdsRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		opts := app.FilterOpts{
			ID:           reqBody.ID,
			Title:        reqBody.Title,
			CreatedTime:  reqBody.CreatedTime,
			AuthorID:     reqBody.AuthorID,
			ModifiedTime: reqBody.ModifiedTime,
			Hidden:       !reqBody.Published,
		}

		l, err := a.GetAdsByFilter(c, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, MultipleAdsSuccessResponse(l))
	}
}
