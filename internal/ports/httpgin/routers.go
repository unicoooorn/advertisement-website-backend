package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.GET("/ads", getPublishedAds(a))
	r.POST("/ads", createAd(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAdById(a))
	r.GET("/search/:title", getAdByTitle(a))
	r.POST("/search", getAdsByFilter(a))

	r.POST("/users", createUser(a))
	r.PUT("/users/:user_id", updateUser(a))
}
