package controllers

import (
	"net/http"
	"time"

	"github.com/cw2/backend/models"
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

func (S *Stocks)UpdateStock(newStock *models.Stock) {
	for i, s := range *S {
		if s.Symbol == newStock.Symbol {
			(*S)[i] = *newStock
		}
	}
}

func (c *Controller) GetStocks(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
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
	ctx.JSON(http.StatusOK, stocks)
}
