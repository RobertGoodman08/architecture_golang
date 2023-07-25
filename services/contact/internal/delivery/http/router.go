package http

import (
	"strings"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"architecture_go/pkg/type/logger"
	docs "architecture_go/services/contact/internal/delivery/http/swagger/docs"
)

func (d *Delivery) initRouter() *gin.Engine {

	if viper.GetBool("IS_PRODUCTION") {
		switch strings.ToUpper(strings.TrimSpace(viper.GetString("LOG_LEVEL"))) {
		case "DEBUG":
			gin.SetMode(gin.DebugMode)
		default:
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		gin.SetMode(gin.DebugMode)
	}

	var router = gin.New()

	router.Use(Tracer())

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(logger.GetLogger(), true))

	d.routerDocs(router.Group("/docs"))

	router.Use(checkAuth)

	d.routerContacts(router.Group("/contacts"))

	d.routerGroups(router.Group("/groups"))

	return router
}

func (d *Delivery) routerContacts(router *gin.RouterGroup) {
	router.POST("/", d.CreateContact)
	router.PUT("/:id", d.UpdateContact)
	router.DELETE("/:id", d.DeleteContact)
	router.GET("/", d.ListContact)
	router.GET("/:id", d.ReadContactByID)
}

func (d *Delivery) routerGroups(router *gin.RouterGroup) {
	router.POST("/", d.CreateGroup)
	router.PUT("/:id", d.UpdateGroup)
	router.DELETE("/:id", d.DeleteGroup)
	router.GET("/", d.ListGroup)
	router.GET("/:id", d.ReadGroupByID)

	router.POST("/:id/contacts/", d.CreateContactIntoGroup)
	router.POST("/:id/contacts/:contactId", d.AddContactToGroup)
	router.DELETE("/:id/contacts/:contactId", d.DeleteContactFromGroup)
}
func (d *Delivery) routerDocs(router *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = "/"

	router.Any("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
