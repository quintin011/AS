package Func

import (
	"os"
	"time"

	c "github.com/otiai10/copy"
)

func BackupTrade() error {
	if _, err := os.Stat("trade/trades.json"); !os.IsNotExist(err) {
		copyerr := c.Copy("trade/trades.json", "trade/trades_"+time.Now().Format("20060102")+".json")
		if err != nil {
			return copyerr
		}
	} else {
		return err
	}
	return nil
}
