package routes

import "github.com/gin-gonic/gin"

func (r *Routes) JWT(rg *gin.RouterGroup) {
<<<<<<< HEAD
		authRouter := rg.Group("jwt",r.Ctrl.HandlerCheck())
		authRouter.GET("/refresh",r.Ctrl.ReToken)
}
=======
	authRouter := rg.Group("jwt", r.Ctrl.HandlerCheck())
	authRouter.GET("/refresh", r.Ctrl.ReToken)
}
>>>>>>> v0.0.2
