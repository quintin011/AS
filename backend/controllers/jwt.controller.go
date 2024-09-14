package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/cw2/backend/encryption"
	"github.com/cw2/backend/models"
	"github.com/gin-gonic/gin"
)

func (c *Controller)HandlerCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.GetHeader("x-uid")
		jwt := ctx.GetHeader("Authorization")
		if jwt == "" || uid == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Panic("Authorization or X-Uid missing")
			return
		}
		jwtToken := encryption.SplitJWT(jwt)
		if jwtToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Panic("Wrong Header")
			return
		}
		t := encryption.ParseToken(jwtToken)
		ttime,err := t.Claims.GetExpirationTime()
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Panic(err)
			return
		}
		if ttime.Equal(time.Now()) || ttime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return	
		}
		tuid, err := t.Claims.GetSubject()
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Panic(err)
			return
		}
		if uid != tuid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			log.Panic("error: User was wrong!")
			return
		}
		var usr models.User
		rslt := c.DB.First(&usr, "uid = ?", tuid)
		if rslt.Error != nil {
			log.Panic(rslt.Error)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if tiss, err := t.Claims.GetIssuer(); err != nil {
			log.Panic(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if tiss != "Prod" {
			log.Panic("error: Token Issuer was wrong!")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func (c *Controller) ReToken(ctx *gin.Context) {
		uid, jwtToken := ctx.GetHeader("x-uid"), encryption.SplitJWT(ctx.GetHeader("Authorization"))
		t := encryption.ParseToken(jwtToken)
		new_jwt := encryption.RefreshToken(t)
		
		ctx.Header("Authorization","Bearer "+new_jwt)
		ctx.Header("x-uid", uid)
		ctx.JSON(http.StatusAccepted,gin.H{"status":"success"})
}
