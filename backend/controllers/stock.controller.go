package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var stocks Stocks
var mkdloc = "stocks/marketdata.json"

func UpdateStockTimestamp() *time.Time {
	now := time.Now().Format("2006-01-02T15:04:05.999999")
	nowT, _ := time.Parse("2006-01-02T15:04:05.999999", now)
	nowT = nowT.Add(-(8 * time.Hour))
	return &nowT
}

func (c *Controller) GetStocks(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	fmt.Println(&symbol)
	stocks = ReadStockJson(mkdloc)
	for _, stock := range stocks {
		if *stock.Symbol == symbol {
			ctx.JSON(http.StatusOK, stock)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, symbol+" not found")
}

func (c *Controller) ListStocks(ctx *gin.Context) {
	stocks = ReadStockJson(mkdloc)
	ctx.JSON(http.StatusOK, gin.H{"stocks": stocks})
}
