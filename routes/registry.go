package routes

import (
	"user-service/controllers"
	userRoutes "user-service/routes/user"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegistry {
	return &Registry{
		controller: controller,
		group:      group,
	}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}

func (r *Registry) userRoute() userRoutes.IUserRoute {
	return userRoutes.NewUserRoute(r.controller, r.group)
}
