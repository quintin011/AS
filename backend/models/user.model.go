package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	UID        uuid.UUID  `gorm:"type:uuid;primary_key"`
	Name       *string    `gorm:"type:varchar(26);not null"`
	Email      *string    `gorm:"uniqueIndex;not null"`
	Password   *string    `gorm:"not null"`
	Mobile     *string    `gorm:"not null"`
	HKID       *string    `gorm:"not null"`
	Balance    float32    `gorm:"type:numeric(10,2);default:0;check: balance >= 0"`
	BankCode   string     `gorm:"type:varchar(3)"`
	BranchCode string     `gorm:"type:varchar(3)"`
	AccountNo  string     `gorm:"type:varchar(9)"`
	Position   []Position `gorm:"foreignkey:UID"`
	Order      []Order    `gorm:"foreignkey:UID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Position struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UID    uuid.UUID `gorm:"type:uuid;"`
	SID    string    `gorm:"type:varchar(6);not null"`
	Volume int       `gorm:"type:numeric;default:1;check: volume > 0 "`
}

func (pos *Position) BeforeCreate(*gorm.DB) error {
	pos.ID = uuid.NewV4()
	return nil
}
func (user *User) BeforeCreate(*gorm.DB) error {
	user.UID = uuid.NewV4()
	return nil
}

type RegUsrIn struct {
	Name     string `json:"name" binding:required`
	Email    string `json:"email" binding:required`
	Password string `json:"password" binding:required`
	Mobile   string `json:"mobile" binding:required`
	HKID     string `json:"hkid" binding:required`
}

type LoginIn struct {
	Email    string `json:"email" binding:required`
	Password string `json:"password" binding:required`
}

type BankIn struct {
	BankCode   string `json:"bank" binding:required`
	BranchCode string `json:"branch" binding:required`
	AccountNo  string `json:"account" binding:required`
}

type Changepwd struct {
	Currpwd string `json:"currpwd" binding:required`
	Newpwd  string `json:"newpwd" binding:required`
}

type Changeinfo struct {
	Name   string `json:"name,omitempty" binding:required`
	Mobile string `json:"mobile,omitempty" binding:required`
	Email  string `json:"email,omitempty" binding:required`
}
type UserinfoOut struct {
	Name       string    `json:"name,omitempty" binding:required`
	Email      string    `json:"email,omitempty" binding:required`
	Mobile     string    `json:"mobile,omitempty" binding:required`
	HKID       string    `json:"hk_id,omitempty" binding:required`
	Balance    float32    `json:"balance,omitempty" binding:required`
	BankAccount Bank	`json:"bank,omitempty" binding:required`	
}
type Bank struct {
	BankCode   string     `json:"code,omitempty" binding:required`
	BranchCode string     `json:"branch,omitempty" binding:required`
	AccountNo  string     `json:"ac,omitempty" binding:required`
}

type POSout struct {
	Symbol string `json:"symbol,omitempty" binding:required`
	Quantity int `json:"quantity,omitempty" binding:required`
}

type AddBal struct {
	Balance float32 `json:"balance,omitempty" binding:required`
}