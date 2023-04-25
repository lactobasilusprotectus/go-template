package delivery

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lactobasilusprotectus/go-template/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// RootHandler represents root handler
type RootHandler struct {
	Title          string `json:"title"`
	Env            string `json:"env"`
	Version        string `json:"version"`
	BuildTimestamp string `json:"build_timestamp"`
}

func NewRootHandler(env string) *RootHandler {
	return &RootHandler{
		Title:          "Projek Akhir DTS",
		Env:            env,
		Version:        "n/a",
		BuildTimestamp: time.Now().Format(time.RFC3339),
	}
}

// Register registers root handler
func (h *RootHandler) Register(g *gin.Engine) {
	g.GET("/", h.Get)
	g.GET("/swagger/*any", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}), ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Get				godoc
//
//	@Summary		Check if server is alive
//	@Description	Check Server is alive.
//	@Produce		json
//	@Tags			ping
//	@Success		200	{object}	RootHandler
//	@Router			/ [get]
func (h *RootHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, h)
}
