package controllers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cw2/backend/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)
var trades Trades
var BidQueue []models.Order
var AskQueue []models.Order
var tradejsloc = "trade/trades.json"

func (c *Controller) SortOrders(symbol string,otype string, ptype string) ([]models.Order,[]models.Order) {
	var tempBQ,tempSQ []models.Order
	switch strings.ToLower(otype) {
	case "limit" :
		switch strings.ToLower(ptype) {
		case "standard" :
			tempBQ = c.TradeListOrder(symbol,true,true,true)
			tempSQ = c.TradeListOrder(symbol,false,true,true)
		case "bid" :
			tempBQ = c.TradeListOrder(symbol,true,true,false)
			tempSQ = c.TradeListOrder(symbol,false,true,false)
		default:
			log.Panic("error: wrong PlaceType")
		}
	case "price":
		switch strings.ToLower(ptype) {
		case "standard" :
			tempBQ = c.TradeListOrder(symbol,true,false,true)
			tempSQ = c.TradeListOrder(symbol,false,false,true)
		case "bid" :
			tempBQ = c.TradeListOrder(symbol,true,false,false)
			tempSQ = c.TradeListOrder(symbol,false,false,false)
		default:
			log.Panic("error: wrong PlaceType")
		}	
	}
	return tempBQ,tempSQ 
}

func (c *Controller)PairOrder() {
	var trade models.Trade
	stocks = ReadStockJson(mkdloc)
	for _, s := range stocks {
		for i:=0;i<4;i++{
			switch i {
			case 0 :
				BidQueue,AskQueue = c.SortOrders(*s.Symbol,"limit","standard")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity == *SO.Quantity {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *BO.Quantity)
							trades = append(trades, trade)
							if bi == len(BidQueue) {
								BidQueue = BidQueue[:bi]	
							} else {
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
							}
							if si == len(AskQueue){
								AskQueue = AskQueue[:si]
							} else {
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
							}	
						}
					}
				}
			case 1 :
				BidQueue,AskQueue = c.SortOrders(*s.Symbol,"limit","bid")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity == *SO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*BO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}	
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*BO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}		
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*BO.Quantity)
								trades = append(trades, trade)
								if bi == len(BidQueue) {
									BidQueue = BidQueue[:bi]	
								} else {
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
								}
								if si == len(AskQueue){
									AskQueue = AskQueue[:si]
								} else {
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}	
							}
						}
					}
				}
			case 2 :
				BidQueue,AskQueue = c.SortOrders(*s.Symbol,"price","standard")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity > *SO.Quantity {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *SO.Quantity)
							trades = append(trades, trade)
							var NewQ int = *BO.Quantity - *SO.Quantity
							newBO := models.Order{
								OID: BO.OID,
								UID: BO.UID,
								Status: BO.Status,
								Method: BO.Method,
								OrderType: BO.OrderType,
								PlaceType: BO.PlaceType,
								Symbol: BO.Symbol,
								Price: BO.Price,
								Quantity: &NewQ,
								CreatedAt: BO.CreatedAt,
								UpdatedAt: BO.UpdatedAt,
							}
							
							BidQueue[bi] = newBO
							if si+1 == len(AskQueue){
								AskQueue = AskQueue[:si]
							} else {
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
							}	
						} else if *SO.Quantity > *BO.Quantity {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *BO.Quantity)
							trades = append(trades, trade)
							var NewQ int = *SO.Quantity - *BO.Quantity
							newSO := models.Order{
								OID: SO.OID,
								UID: SO.UID,
								Status: SO.Status,
								Method: SO.Method,
								OrderType: SO.OrderType,
								PlaceType: SO.PlaceType,
								Symbol: SO.Symbol,
								Price: SO.Price,
								Quantity: &NewQ,
								CreatedAt: SO.CreatedAt,
								UpdatedAt: SO.UpdatedAt,
							}
							AskQueue[si] = newSO
							if bi == len(BidQueue) {
								BidQueue = BidQueue[:bi]	
							} else {
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
							}
						} else {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *BO.Quantity)
							trades = append(trades, trade)
							if bi == len(BidQueue) {
								BidQueue = BidQueue[:bi]	
							} else {
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
							}
							if si == len(AskQueue){
								AskQueue = AskQueue[:si]
							} else {
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
							}	
						}
					}
				}
			case 3 :
				BidQueue,AskQueue = c.SortOrders(*s.Symbol,"price","bid")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity > *SO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*SO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}	
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*SO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*SO.Quantity)
								trades = append(trades, trade)
								if bi == len(BidQueue) {
									BidQueue = BidQueue[:bi]	
								} else {
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
								}
								if si == len(AskQueue){
									AskQueue = AskQueue[:si]
								} else {
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}
							}
						} else if *SO.Quantity > *BO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*BO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*BO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*BO.Quantity)
								trades = append(trades, trade)
								if bi == len(BidQueue) {
									BidQueue = BidQueue[:bi]	
								} else {
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
								}
								if si == len(AskQueue){
									AskQueue = AskQueue[:si]
								} else {
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}
							}	
						} else {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*SO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*SO.Quantity)
									trades = append(trades, trade)
									if bi == len(BidQueue) {
										BidQueue = BidQueue[:bi]	
									} else {
										BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
									}
									if si == len(AskQueue){
										AskQueue = AskQueue[:si]
									} else {
										AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
									}	
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*SO.Quantity)
								trades = append(trades, trade)
								if bi == len(BidQueue) {
									BidQueue = BidQueue[:bi]	
								} else {
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)	
								}
								if si == len(AskQueue){
									AskQueue = AskQueue[:si]
								} else {
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}
							}	
						}
					}
				}
			}
		}
	}
	err := trades.WriteTradeJson(tradejsloc)
	if err != nil {
		log.Panic(err)
		return
	}
}

func (c *Controller) ProcessTrading() {
	var BO models.Order
	var SO models.Order
	var usr models.User
	var POS models.Position
	trades := ReadTradeJson(tradejsloc)
	for ti, trade := range trades {
		getBO := c.DB.Find(&BO,"o_id = ?",trade.BuyOID)
		if getBO.Error != nil {
			log.Panic(getBO.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			getSO := c.DB.Find(&SO,"o_id = ?",trade.SellOID)
			if getSO.Error != nil {
				log.Panic(getSO.Error)
				log.Panic(fmt.Errorf("Buy Order(%s) & Sell Order(%s) not Found",trade.BuyOID,trade.SellOID))
			}
			updateSOStatus := c.DB.Model(&SO).Where("o_id = ?",trade.SellOID).Update("status","Conflicted")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
			}
		}
		getBuyer := c.DB.Find(&usr,"uid = ?",BO.UID)
		if getBuyer.Error != nil {
			log.Panic(getBuyer.Error)
			log.Panic(fmt.Errorf("Buyer(%s) of Order(%s) not Found",BO.UID,BO.OID))
			trades = append(trades[:ti],trades[ti+1:]... )
			getSO := c.DB.Find(&SO,"o_id = ?",trade.SellOID)
			if getSO.Error != nil {
				log.Panic(getSO.Error)
				log.Panic(fmt.Errorf("Buy Order(%s) & Sell Order(%s) not Found",trade.BuyOID,trade.SellOID))
			}
			updateSOStatus := c.DB.Model(&SO).Where("o_id = ?",trade.SellOID).Update("status","Conflicted")
			if updateSOStatus != nil {
				log.Panic(updateSOStatus.Error)
			}
		}
		tprice:=decimal.NewFromFloat32(*trade.Price).Mul(decimal.NewFromFloat32(float32(*trade.TVol))) 
		price,_ := tprice.Float64()
		if usr.Balance < float32(price) {
			updateBOStatus := c.DB.Model(&BO).Update("status","Rejected")
			trades= append(trades[:ti],trades[ti+1:]... )
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}
		} else {
			updateBOStatus := c.DB.Model(&BO).Update("status","Progressing")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}
		}
		getSO := c.DB.Find(&SO,"o_id = ?",trade.SellOID)
		if getSO.Error != nil {
			log.Panic(getSO.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.Model(&BO).Update("status","Pending")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}	
		}
		getSeller := c.DB.Find(&usr,"uid = ?",SO.UID)
		if getSeller.Error != nil {
			log.Panic(getSO.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.Model(&BO).Update("status","Pending")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}	
		}
		chkSellerPosition := c.DB.Find(&POS,models.Position{UID: *SO.UID,SID: *SO.Symbol})
		if chkSellerPosition.Error != nil {
			log.Panic(chkSellerPosition.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.Model(&BO).Update("status","Pending")	
			if updateBOStatus.Error != nil {
			 	log.Panic(updateBOStatus.Error)
			}
			updateSOStatus := c.DB.Model(&SO).Update("status","Rejected")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
			}
		}
		if POS.Volume < *trade.TVol {
			log.Panic(fmt.Errorf("Seller(%s) had not enough volume of Symbol(%s)",SO.UID,*SO.Symbol))
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.Model(&BO).Update("status","Pending")	
			if updateBOStatus.Error != nil {
			 	log.Panic(updateBOStatus.Error)
			}
			updateSOStatus := c.DB.Model(&SO).Update("status","Rejected")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
			}
			return
		} else {
			newPOS := models.Position{
				UID: *SO.UID,
				SID: *SO.Symbol,
				Volume: POS.Volume - *trade.TVol,
			}
			updateSellerPosition := c.DB.Model(&POS).Where(models.Position{UID: *SO.UID,SID: *SO.Symbol}).Updates(newPOS)
			if updateSellerPosition.Error != nil {
				log.Panic(updateSellerPosition.Error)
				return
			}
			tprice:=decimal.NewFromFloat32(*trade.Price).Mul(decimal.NewFromFloat32(float32(*trade.TVol))) 
			price,_ := tprice.Float64()
			newBalance := usr.Balance + float32(price)	
			updateSellerBalance := c.DB.Model(&usr).Where(models.User{UID: *SO.UID}).Updates(models.User{Balance: newBalance})
			if updateSellerBalance.Error != nil {
				log.Panic(updateSellerBalance.Error)
				return
			}
			updateSOStatus := c.DB.Model(&SO).Update("status","Accepted")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
				return
			}
			newPOS = models.Position{
				UID: *BO.UID,
				SID: *BO.Symbol,
				Volume: POS.Volume + *trade.TVol,
			}
			chkBuyerPosition := c.DB.Find(&POS,models.Position{UID: *BO.UID,SID: *BO.Symbol})
			if errors.Is(chkBuyerPosition.Error, gorm.ErrRecordNotFound) {
				rslt := c.DB.Create(&newPOS)
				if rslt.Error != nil {
					log.Panic(rslt.Error)
					return
				}
			} else {
				rslt := c.DB.Model(&POS).Where(models.Position{UID: *BO.UID,SID: *BO.Symbol}).Updates(newPOS)
				if rslt.Error != nil {
					log.Panic(rslt.Error)
					return
				}
			}
			tprice = decimal.NewFromFloat32(*trade.Price).Mul(decimal.NewFromFloat32(float32(*trade.TVol))) 
			price,_ = tprice.Float64()
			newBalance = usr.Balance - float32(price)
			updateBuyerBalance := c.DB.Model(&usr).Where(models.User{UID: *BO.UID}).Updates(models.User{Balance: newBalance})
			if updateBuyerBalance.Error != nil {
				log.Panic(updateBuyerBalance.Error)
				return
			}
			updateBOStatus := c.DB.Model(&BO).Update("status","Accepted")
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
				return
			}
			var newHigh float32
			var newLow float32
			var newStock models.Stock
			for _, s := range stocks {
				if s.Symbol == BO.Symbol {
					if *trade.Price > s.High {
						newHigh = *trade.Price
						newLow = s.Low
					} else if *trade.Price < s.Low {
						newHigh = s.High
						newLow = *trade.Price
					} else {
						newHigh = s.High
						newLow = s.Low
					}
					newStock = models.Stock{
						Symbol: BO.Symbol,
						Timestamp: UpdateStockTimestamp(),
						CurrBid: BO.Price,
						CurrAsk: SO.Price,
						LastTrade: *trade.Price,
						High: newHigh,
						Low: newLow,
						Volume: s.Volume,
					}
				}
			}
			stocks.UpdateStock(&newStock)
			if ti+1 > len(trades) {
				trades = trades[:ti]
			} else {
				trades = append(trades[:ti],trades[ti+1:]... )
			}
		}
	}
	stocks.WriteStockJson(mkdloc)
	err := os.Rename(tradejsloc,"trade/trades_DONE"+time.Now().Format("20060102_150405s999999")+".json")
	if err != nil {
		log.Panic(err)
		return
	}
}

func (c *Controller)TradeRun() {
	for {
			c.PairOrder()
			time.Sleep(time.Second * 5)
			c.ProcessTrading()
			time.Sleep(time.Second * 5)
	}
}