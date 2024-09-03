package backup

import (
	"log"
	"os"
	"time"

	"github.com/cw2/backend/backup/Func"
)

var (
	currtime time.Time
	ltime    []time.Time
	currdate time.Time
	ldate    time.Time
)

func Backup() {
	currtime = time.Now()
	logfile, err := os.OpenFile("backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
	for range time.Tick(time.Minute * 1) {
		if !currtime.IsZero() {
			ltime = append(ltime, currtime)
		}
		if !currdate.IsZero() {
			ldate = currdate
		}
		currtime = time.Now()
		currdate, _ = time.Parse("2006-01-02", currtime.Format("2006-01-02"))
		if !currdate.Equal(ldate) && !ldate.IsZero() {
			log.Println("start backup apps")
			err := Func.BackupTrade()
			if err != nil {
				log.Panic(err)
			}
			log.Println("end backup apps")
		}
		if len(ltime) == 5 {
			if currtime.Sub(ltime[len(ltime)-5]).Round(time.Minute) == (time.Minute * 5) {
				log.Println(ltime[len(ltime)-5].Format("15:04:05"))
				log.Println("after 5 min")
				log.Println(currtime.Format("15:04:05"))
			}
			ltime = ltime[4:]
		}
	}
}
