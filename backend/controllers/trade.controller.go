package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/cw2/backend/models"
)
var trades Trades
var BidQueue []models.Order
var AskQueue []models.Order
var tradejsloc = "trade/trades.json"

func (c *Controller) SortOrders(symbol string,otype string, ptype string) {
	switch strings.ToLower(otype) {
	case "limit" :
		switch strings.ToLower(ptype) {
		case "standard" :
			BidQueue = c.TradeListOrder(symbol,true,true,true)
			AskQueue = c.TradeListOrder(symbol,false,true,true)
		case "bid" :
			BidQueue = c.TradeListOrder(symbol,true,true,false)
			AskQueue = c.TradeListOrder(symbol,false,true,false)
		default:
			log.Panic("error: wrong PlaceType")
		}
	case "price":
		switch strings.ToLower(ptype) {
		case "standard" :
			BidQueue = c.TradeListOrder(symbol,true,false,true)
			AskQueue = c.TradeListOrder(symbol,false,false,true)
		case "bid" :
			BidQueue = c.TradeListOrder(symbol,true,false,false)
			AskQueue = c.TradeListOrder(symbol,false,false,false)
		default:
			log.Panic("error: wrong PlaceType")
		}	
	}
}

func (c *Controller)PairOrder() {
	var trade models.Trade
	stocks = ReadStockJson(mkdloc)
	for _, s := range stocks {
		for i:=0;i<4;i++{
			switch i {
			case 0 :
				c.SortOrders(*s.Symbol,"limit","standard")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity == *SO.Quantity {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *BO.Quantity)
							trades = append(trades, trade)
							BidQueue = append(BidQueue[:bi],BidQueue[bi+1:]... )
							AskQueue = append(AskQueue[:si],AskQueue[si+1:]... )
						}
					}
				}
			case 1 :
				c.SortOrders(*s.Symbol,"limit","bid")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity == *SO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*BO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]... )
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*BO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*BO.Quantity)
								trades = append(trades, trade)
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
							}
						}
					}
				}
			case 2 :
				c.SortOrders(*s.Symbol,"price","standard")
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
							AskQueue = append(AskQueue[:si],AskQueue[si+1:]... )	
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
							BidQueue = append(BidQueue[:bi],AskQueue[si+1:]... )
						} else {
							trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade, *BO.Quantity)
							trades = append(trades, trade)
							BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
							AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
						}
					}
				}
			case 3 :
				c.SortOrders(*s.Symbol,"price","bid")
				for bi, BO := range BidQueue {
					for si, SO := range AskQueue {
						if *BO.Quantity > *SO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*SO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*SO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*SO.Quantity)
								trades = append(trades, trade)
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
							}
						} else if *SO.Quantity > *BO.Quantity {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*BO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*BO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*BO.Quantity)
								trades = append(trades, trade)
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
							}	
						} else {
							if *BidQueue[len(BidQueue)-1].Price > s.High {
								if *SO.Price >= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*BO.Price,*SO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
								} 
							} else if *AskQueue[len(AskQueue)-1].Price < s.Low {
								if *BO.Price <= s.LastTrade {
									trade.Create(BO.OID.String(),SO.OID.String(),*SO.Price,*SO.Quantity)
									trades = append(trades, trade)
									BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
									AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)	
								}
							} else {
								trade.Create(BO.OID.String(),SO.OID.String(),s.LastTrade,*SO.Quantity)
								trades = append(trades, trade)
								BidQueue = append(BidQueue[:bi], BidQueue[bi+1:]...)
								AskQueue = append(AskQueue[:si], AskQueue[si+1:]...)
							}	
						}
					}
				}
			}
		}
	}
	trades.WriteTradeJson(tradejsloc)
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
			updateSOStatus := c.DB.First(&SO,"o_id = ?",trade.SellOID).Update("status","Conflicted")
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
			updateSOStatus := c.DB.First(&SO,"o_id = ?",trade.SellOID).Update("status","Conflicted")
			if updateSOStatus != nil {
				log.Panic(updateSOStatus.Error)
			}
		}
		if usr.Balance < (*trade.Price * float32(*trade.TVol)) {
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Rejected")
			trades= append(trades[:ti],trades[ti+1:]... )
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}
		} else {
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Progressing")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}
		}
		getSO := c.DB.Find(&SO,"o_id = ?",trade.SellOID)
		if getSO.Error != nil {
			log.Panic(getSO.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Pending")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}	
		}
		getSeller := c.DB.Find(&usr,"uid = ?",SO.UID)
		if getSeller.Error != nil {
			log.Panic(getSO.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Pending")	
			if updateBOStatus.Error != nil {
				log.Panic(updateBOStatus.Error)
			}	
		}
		chkSellerPosition := c.DB.Find(&POS,models.Position{UID: usr.UID,SID: *SO.Symbol})
		if chkSellerPosition.Error != nil {
			log.Panic(chkSellerPosition.Error)
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Pending")	
			if updateBOStatus.Error != nil {
			 	log.Panic(updateBOStatus.Error)
			}
			updateSOStatus := c.DB.First(&SO,"o_id = ?",SO.OID).Update("status","Rejected")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
			}
		}
		if POS.Volume < *trade.TVol {
			log.Panic(fmt.Errorf("Seller(%s) had not enough volume of Symbol(%s)",SO.UID,*SO.Symbol))
			trades = append(trades[:ti],trades[ti+1:]... )
			updateBOStatus := c.DB.First(&BO,"o_id = ?",BO.OID).Update("status","Pending")	
			if updateBOStatus.Error != nil {
			 	log.Panic(updateBOStatus.Error)
			}
			updateSOStatus := c.DB.First(&SO,"o_id = ?",SO.OID).Update("status","Rejected")
			if updateSOStatus.Error != nil {
				log.Panic(updateSOStatus.Error)
			}
		} else {
			
		}
	}
}