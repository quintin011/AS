package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	UID uuid.UUID `gorm:"type:uuid;primary_key"`
	Name *string `gorm:"type:varchar(26);not null"`
	Email *string `gorm:"uniqueIndex;not null"`
	Password *string `gorm:"not null"`
	Mobile *string `gorm:"not null"`
	HKID *string `gorm:"not null"`
	Balance float32 `gorm:"type:numeric(10,2);default:0;check: balance >= 0"`
	BankCode string
	BranchCode string
	AccountNo string
	Position []Position `gorm:"foreignkey:UID"`
	Order []Order `gorm:"foreignkey:UID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Position struct{
	gorm.Model
	UID *uuid.UUID `gorm:"type:uuid;"`
	SID *string `gorm:"type:varchar(6);not null"`
	Volume *int `gorm:"type:numeric;default:1;check: volume > 0 "` 
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.UID = uuid.NewV4()
	return nil
}

type RegUsrIn struct {
	Name *string `json:"name" binding:required`
	Email *string `json:"email" binding:required`
	Password *string `json:"password" binding:required`
	Mobile *string `json:"mobile" binding:required`
	HKID *string `json:"hkid" binding:required`
}

type LoginIn struct {
	Email string `json:"email" binding:required`
	Password string `json:"password" binding:required`	
}

type BankIn struct {
	BankCode string `json:"bank" binding:required`
	BranchCode string `json:"branch" binding:required`
	AccountNo string `json:"account" binding:required` 
}

type Changepwd struct {
	Currpwd string `json:"currpwd" binding:required`
	Newpwd string `json:"newpwd" binding:required`
}

type Changeinfo struct {
	Name string `json:"name,omitempty" binding:required`
	Mobile string `json:"mobile,omitempty" binding:required`
	Email string `json:"email,omitempty" binding:required`
}