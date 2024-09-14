package routes

import "github.com/gin-gonic/gin"

func (r *Routes) Order(rg *gin.RouterGroup) {
	router := rg.Group("order", r.Ctrl.HandlerCheck())
	router.POST("/create", r.Ctrl.CreateOrder)
	router.GET("", r.Ctrl.ListOrder)
	router.POST("/:oid/cancel", r.Ctrl.CancelOrder)
	router.GET("/:oid", r.Ctrl.GetOrder)
}
