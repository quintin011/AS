package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Order struct {
	OID uuid.UUID `gorm:"type:uuid;primary_key;"`
	UID *uuid.UUID `gorm:"type:uuid;"`
	Status *string `gorm:"type:varchar(26);not null"`
	Method *bool `gorm:"type:bool;not null"` // 1 = buy, 0 = sell 
	OrderType *bool `gorm:"type:bool;not null"` // 1 = limit, 0 = market price
	PlaceType *bool `gorm:"type:bool;not null"` // 1 = standard, 0 = bid
	Symbol *string `gorm:"type:varchar(6);not null"` 
	Price *float32 `gorm:"type:numeric(,2);default:0.01;check: price > 0"`
	Volume *int `gorm:"type:numeric;default:1;check: volume > 0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
func (o *Order) BeforeCreate(*gorm.DB) error {
	o.OID = uuid.NewV4()
	return nil
}
type Stock struct {
	Symbol *string `json:"symbol" binding:"required"`
	Timestamp *time.Time `json:"updated_at" binding:"required"`
	CurrBid *float32 `json:"currbid" binding:"required"`
	CurrAsk *float32 `json:"currask" binding:"required"`
	LastTrade float32 `json:"lasttrade" binding:"required"`
	High float32 `json:"high_price" binding:"required"`
	Low float32 `json:"low_price" binding:"required"`
	Volume int `json:"vol" binding:"required"`
}

type Trade struct {
	TID uuid.UUID
	BuyOID string `json:"buyer" binding:"required"`
	SellOID string `json:"seller" binding:"required"`
	Price *float32 `json:"price" binding:"required"`
	TVol *int `json:"vol" binding:"required"`
	Timestamp time.Time `json:"created_at" binding:"required"`
}
func (T *Trade) Create(BOID string,SOID string,price float32,vol int) error {
	T.TID = uuid.NewV4()
	T.BuyOID = BOID
	T.SellOID = SOID
	T.Price = &price
	T.TVol = &vol
	T.Timestamp = time.Now()
	return nil
}