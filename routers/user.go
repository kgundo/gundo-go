package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kgundo/gundo-go/controllers"
)

func UserRouter(v1 *gin.RouterGroup) {
	userController := controllers.InitUserController()
	v1.GET("/users", userController.GetAllUsers)
	v1.GET("/user/:id", userController.GetUserByID)
	v1.POST("/user", userController.CreateUser)
	v1.DELETE("/user/:id", userController.DeleteUser)
	v1.PUT("/user/:id", userController.UpdateUser)
}
