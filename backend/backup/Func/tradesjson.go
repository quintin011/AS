package Func

import (
	"os"
	"time"

	c "github.com/otiai10/copy"
)

func BackupTrade() error {
	if _, err := os.Stat("bak"); !os.IsNotExist(err) {
		os.Mkdir("bak",0755)
	}
	if _, err := os.Stat("trade"); !os.IsNotExist(err) {
		copyerr := c.Copy("trade", "bak/trades_"+time.Now().Format("20060102"))
		if err != nil {
			return copyerr
		}
	} else {
		return err
	}
	return nil
}
