package routes

import "github.com/gin-gonic/gin"

func (r *Routes) Trade(rg *gin.RouterGroup) {
	router := rg.Group("trade",r.Ctrl.HandlerCheck())
	router.POST("/register",r.Ctrl.SignUp)
	router.POST("/login", r.Ctrl.LoginUser)
}