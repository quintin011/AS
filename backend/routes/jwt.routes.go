package routes

import "github.com/gin-gonic/gin"

func (r *Routes) JWT(rg *gin.RouterGroup) {
	authRouter := rg.Group("jwt", r.Ctrl.HandlerCheck())
	authRouter.GET("/refresh", r.Ctrl.ReToken)
}
