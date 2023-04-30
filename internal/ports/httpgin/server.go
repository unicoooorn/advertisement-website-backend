package httpgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"homework9/internal/app"
)

func customLogger(c *gin.Context) {
	t := time.Now()
	c.Next()

	log.Printf("%s method, %s path, %d status, %s latency",
		c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(t).String())
}

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}

	s.app.Use(gin.Recovery())
	s.app.Use(customLogger)

	api := s.app.Group("/api/v1")
	s.app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	AppRouter(api, a)

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
