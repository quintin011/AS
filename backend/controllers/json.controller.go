package controllers

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cw2/backend/models"
	cp "github.com/otiai10/copy"
)
type Stocks []models.Stock
type Trades []models.Trade

func ReadStockJson(fname string) Stocks {
	var stocks Stocks
	JSON, err:= os.Open(fname)
	if err != nil {
		log.Panic(err)
	}
	defer JSON.Close()
	Item, _ := os.ReadFile(fname)
	json.Unmarshal(Item, &stocks)
	return stocks
}

func ReadTradeJson(fname string) Trades {
	var trades Trades
	JSON, err:= os.Open(fname)
	if err != nil {
		log.Panic(err)
	}
	defer JSON.Close()
	Item, _ := os.ReadFile(fname)
	json.Unmarshal(Item, &trades)
	return trades
}

func (T *Trades)WriteTradeJson(fname string) error{
	JSON, err := json.Marshal(T)
	if err != nil {
		return err
	}
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		err = os.WriteFile(fname,JSON,0644)
		if err != nil {
			return err
		}
	} else {
		err := cp.Copy(fname,fname+"_"+time.Now().Format("20060102_150405s999999"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (S *Stocks)WriteStockJson(fname string) error{
	JSON, err := json.Marshal(S)
	if err != nil {
		return err
	}
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		err = os.WriteFile(fname,JSON,0644)
		if err != nil {
			return err
		}
	} else {
		err := cp.Copy(fname,fname+"_"+time.Now().Format("20060102_150405s999999"))
		if err != nil {
			return err
		}
	}
	return nil
}