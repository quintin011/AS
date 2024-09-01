package routes

import "github.com/gin-gonic/gin"

func (r *Routes) Stock(rg *gin.RouterGroup) {
	router := rg.Group("stock")
	router.GET("/", r.Ctrl.ListStocks)
	router.GET("/:symbol", r.Ctrl.GetStocks)
}
