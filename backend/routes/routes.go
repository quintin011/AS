package routes

import (
	"github.com/cw2/backend/controllers"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	Ctrl controllers.Controller
}

func NewRoutes(Ctrl controllers.Controller) Routes {
	return Routes{Ctrl}
}

func (r *Routes) MainR(rg *gin.RouterGroup) {
	router := rg.Group("v1")
	router.POST("/register", r.Ctrl.SignUp)
	router.POST("/login", r.Ctrl.LoginUser)
	r.JWT(router)
	r.Order(router)
	r.User(router)
	r.Stock(router)
}
