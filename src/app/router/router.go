package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/kallepan/pcr-backend/app/middleware"
	"gitlab.com/kallepan/pcr-backend/config"
)

func Init(init *config.Initialization) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())

	auth := router.Group("/api")
	{
		auth.POST("/token", init.UserCtrl.LoginUser)
		auth.POST("/register", init.UserCtrl.RegisterUser)
	}

	api := router.Group("/api/v1")
	{
		api.GET("/ping", init.SysCtrl.GetPing)

		// Import
		api.POST("/import", init.ImportCtrl.ImportSample)
	}

	return router
}
