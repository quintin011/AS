package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/cw2/backend/encryption"
	"github.com/cw2/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) SignUp(ctx *gin.Context) {
	var payload models.RegUsrIn
	fencmsg := "Failed to Encrypt."
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Failed to hash password."})
		return
	}
	pwd, err := encryption.Encrypt(hash)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fencmsg})
		return
	}
	
	id, err := encryption.Encrypt([]byte(payload.HKID))
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fencmsg})
		return
	}
	last4no, err := encryption.Encrypt([]byte(payload.Mobile)[4:])
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": fencmsg})
		return
	}
	NewEmail := strings.ToLower(payload.Email)
	NewMobile := string([]byte(payload.Mobile)[:4]) + last4no
	newUser := models.User{
		Name:     &payload.Name,
		Email:    &NewEmail,
		Password: &pwd,
		Mobile:   &NewMobile,
		HKID:     &id,
	}

	rslt := c.DB.Create(&newUser)

	if rslt.Error != nil && strings.Contains(rslt.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Email already exist, please use another email address"})
		return
	} else if rslt.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Registered successfully, please login"})
}

func (c *Controller) LoginUser(ctx *gin.Context) {
	var payload models.LoginIn
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var usr models.User
	rslt := c.DB.First(&usr, "email = ?", strings.ToLower(payload.Email))
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}
	pwd, err := encryption.Decrypt(usr.Password)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to Decrypt."})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(payload.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}
	jwt := encryption.GenToken(&usr)
	ctx.Header("Authorization", "Bearer "+encryption.Signstring(jwt))
	ctx.Header("x-uid", usr.UID.String())
	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}
