package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/cw2/backend/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func (c *Controller) ListOrder(ctx *gin.Context) {
	uid := ctx.GetHeader("X-Uid")
	var orders []models.Order
	var usr models.User
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": rslt.Error.Error()})
		return
	}
	rslt = c.DB.Find(&orders, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": rslt.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
func (c *Controller) TradeListOrder(symbol string, method bool,otype bool, ptype bool) []models.Order {
	var orders []models.Order
	if method {
		rslt := c.DB.Where(models.Order{Symbol: &symbol, Method: &method, OrderType: &otype, PlaceType: &ptype, Status: "Pending"}).Order("price desc").Find(&orders)
		if rslt.Error != nil {
			log.Panic(rslt.Error)
			return nil
		}
	} else {
		rslt := c.DB.Where(models.Order{Symbol: &symbol, Method: &method, OrderType: &otype, PlaceType: &ptype, Status: "Pending"}).Order("price asc").Find(&orders)
		if rslt.Error != nil {
			log.Panic(rslt.Error)
			return nil
		}
	}
	return orders
}

func (c *Controller) GetOrder(ctx *gin.Context) {
	uid, _ := uuid.FromString(ctx.GetHeader("X-Uid"))
	var usr models.User
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": rslt.Error.Error()})
		return
	}
	oid, _ := uuid.FromString(ctx.Param("oid"))
	var order models.Order
	rslt = c.DB.First(&order, &models.Order{OID: oid, UID: &uid})
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid oid"})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (c *Controller) CreateOrder(ctx *gin.Context) {
	var payload *models.OrderIn
	uid, err := uuid.FromString(ctx.GetHeader("X-Uid"))
	if err != nil {
		log.Panic(err)
		return
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var method, otype, ptype bool
	switch strings.ToLower(*payload.Method) {
	case "buy":
		method = true
	case "sell":
		method = false
	default:
		method = true
	}
	switch strings.ToLower(*payload.OrderType) {
	case "limit":
		method = true
	case "price":
		method = false
	default:
		method = true
	}
	switch strings.ToLower(*payload.PlaceType) {
	case "standard":
		method = true
	case "bid":
		method = false
	default:
		method = true
	}

	newOrder := models.Order{
		UID:       &uid,
		Status:    "Pending",
		Method:    &method,
		OrderType: &otype,
		PlaceType: &ptype,
		Symbol:    payload.Symbol,
		Price:     payload.Price,
		Quantity:  payload.Quantity,
	}

	rslt := c.DB.Create(&newOrder)
	if rslt.Error != nil {
		log.Panic(rslt.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (c *Controller) CancelOrder(ctx *gin.Context) {
	uid, _ := uuid.FromString(ctx.GetHeader("X-Uid"))
	var Order models.Order
	var usr models.User
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": rslt.Error.Error()})
		return
	}
	oid, _ := uuid.FromString(ctx.Param("oid"))
	rslt = c.DB.First(&Order, &models.Order{OID: oid, UID: &uid})
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid oid"})
		return
	}
	rslt = c.DB.Model(&Order).Where(&models.Order{OID: oid, UID: &uid}).Update("status", "Cancelled")
	if rslt.Error != nil {
		log.Panic(rslt.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
}
