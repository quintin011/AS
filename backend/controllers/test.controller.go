package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/cw2/backend/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Controller)AddPos(ctx *gin.Context){
	var usr models.User
	var POS models.Position
	user, _ := uuid.FromString(ctx.GetHeader("X-Uid"))
	sym := ctx.Query("symbol")
	quan,_:= strconv.Atoi(ctx.Query("quan"))
	if err := c.DB.First(&usr, "uid = ?", user).Error; err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err := c.DB.First(&POS,models.Position{UID: user,SID: sym}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		newPOS := models.Position{
			UID: user,
			SID: sym,
			Volume: quan,
		}
		rslt := c.DB.Model(&POS).Create(&newPOS)
		if rslt.Error != nil {
			log.Panic(rslt.Error)
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.Status(http.StatusCreated)
	} else {
		newPOS := models.Position{
			UID: user,
			SID: sym,
			Volume: POS.Volume + quan,
		}
		rslt := c.DB.Model(&POS).Where(models.Position{UID:user,SID: sym}).Updates(newPOS)
		if rslt.Error != nil {
			log.Panic(rslt.Error)
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}