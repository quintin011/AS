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
func (c *Controller) GetUser(ctx *gin.Context) {
	uid := ctx.GetHeader("X-Uid")
	var usr models.User
	var err error
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid uid"})
		return
	}
	var mobile,hkid string
	first4no := string([]byte(*usr.Mobile)[:4])
	last4no := string([]byte(*usr.Mobile)[4:])
	last4no, err = encryption.Decrypt(&last4no)
	if err != nil {
		log.Panic(err)
	}
	mobile = first4no+last4no
	hkid, err = encryption.Decrypt(usr.HKID)
	if err != nil {
		log.Panic(err)
	}	
	bank := models.Bank{
		BankCode: usr.BankCode,
		BranchCode: usr.BranchCode,
		AccountNo: usr.AccountNo,
	}
	ousr := models.UserinfoOut{
		Name: *usr.Name,
		Email: *usr.Email,
		Mobile: mobile,
		HKID: hkid,
		Balance: usr.Balance,
		BankAccount: bank,
		Position: usr.Position,
		Order: usr.Order,
	}
	ctx.JSON(http.StatusOK, &ousr)
}

func (c *Controller) UpdateBankInfo(ctx *gin.Context) {
	var payload *models.BankIn
	var usr models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	uid := ctx.GetHeader("X-Uid")
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid uid"})
		return
	}
	rslt = c.DB.Model(&usr).Updates(models.User{
		BankCode:   payload.BankCode,
		BranchCode: payload.BranchCode,
		AccountNo:  payload.AccountNo,
	},
	)
	if rslt.Error != nil {
		log.Panic(rslt.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}
func (c *Controller) ChangePassword(ctx *gin.Context) {
	var payload *models.Changepwd
	var usr models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	uid := ctx.GetHeader("X-Uid")
	rslt := c.DB.First(&usr, "uid = ?", uid)
	if rslt.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid uid"})
		return
	}
	pwd, err := encryption.Decrypt(usr.Password)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(payload.Currpwd))
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*&payload.Newpwd), 10)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	newpwd, err := encryption.Encrypt(hash)
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
		return
	}
	rslt = c.DB.Model(&usr).Update("password", newpwd)
	if rslt.Error != nil {
		log.Panic(rslt.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}
func (c *Controller) ChangeUserInfo(ctx *gin.Context) {
	var payload *models.Changeinfo
	var usr models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	uid := ctx.GetHeader("X-Uid")
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail"})
			return
		}
		NewMobile := string([]byte(payload.Mobile)[:4]) + last4no
		ninfo.Mobile = &NewMobile
	}
	if payload.Name != "" {
		ninfo.Name = &payload.Name
	}
	rslt = c.DB.Model(&usr).Updates(ninfo)
	if rslt.Error != nil {
		log.Panic(rslt.Error)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
}
