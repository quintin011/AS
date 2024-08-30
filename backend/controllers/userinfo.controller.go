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
func (c *Controller)UpdateBankInfo(ctx *gin.Context) {
	var payload *models.BankIn
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status":"fail","message": err.Error()})
		return
	}
	uid := ctx.Param("uid")
	rslt := c.DB.Model(&models.User{}).Where("uid = ?", uid).Updates(models.User{
			BankCode: payload.BankCode, 
			BranchCode: payload.BranchCode, 
			AccountNo: payload.AccountNo,
		},
	)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error","message": rslt.Error.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted,gin.H{"status": "success"})
}
func (c *Controller)ChangePassword(ctx *gin.Context) {
	var payload *models.Changepwd
	var usr models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status":"fail","message": err.Error()})
		return
	}
	uid := ctx.Param("uid")
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid uid"})
		return
	}
	pwd, err := encryption.Decrypt(usr.Password)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError,gin.H{"status":"fail", "message": "Failed to Decrypt."})
		return
	} 
	err = bcrypt.CompareHashAndPassword([]byte(pwd),[]byte(payload.Currpwd))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"status":"fail","message":"password"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*&payload.Newpwd),10)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest,gin.H{"status":"fail", "message": "Failed to hash password."})
		return
	}
	newpwd,err := encryption.Encrypt(hash)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError,gin.H{"status":"fail", "message": "Failed to Encrypt."})
		return
	}
	rslt = c.DB.Model(&usr).Update("password", newpwd)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error","message": rslt.Error.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted,gin.H{"status": "success"})	
}
func (c *Controller)ChangeUserInfo(ctx *gin.Context) {
	var payload *models.Changeinfo
	var usr models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status":"fail","message": err.Error()})
		return
	}
	uid := ctx.Param("uid")
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid uid"})
		return
	}
	var ninfo models.User
	if payload.Email != "" {
		NewEmail := strings.ToLower(payload.Email)
		ninfo.Email = &NewEmail
	}
	if payload.Mobile != "" {
		last4no, err := encryption.Encrypt([]byte(payload.Mobile)[4:])
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"status":"fail", "message": "Failed to Encrypt."})
			return
		}
		NewMobile := string([]byte(payload.Mobile)[:4])+last4no
		ninfo.Mobile = &NewMobile
	}
	if payload.Name != "" {
		ninfo.Name = &payload.Name
	}
	rslt = c.DB.Model(&usr).Updates(ninfo)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error","message": rslt.Error.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted,gin.H{"status": "success"})
}