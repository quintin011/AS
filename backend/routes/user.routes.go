package routes

import "github.com/gin-gonic/gin"

func (r *Routes) User(rg *gin.RouterGroup) {
	router := rg.Group("user", r.Ctrl.HandlerCheck())
	router.POST("/update/bankinfo", r.Ctrl.UpdateBankInfo)
	router.POST("/update/password", r.Ctrl.ChangePassword)
	router.POST("/update/userinfo", r.Ctrl.ChangeUserInfo)
	router.POST("/addbalance", r.Ctrl.AddBalance)
	router.GET("/",r.Ctrl.GetUser)
	router.GET("/pos",r.Ctrl.ListPos)
	r.TestUser(router)
}

func (r *Routes)TestUser(rg *gin.RouterGroup){
	router := rg.Group("test")
	router.POST("/add/positions",r.Ctrl.AddPos)
}
