package synmedreader

import (
	"fmt"
	"time"
)

func makeRxNum(din string, storeNum int) string {
	//GetUnixTime
	firstOfYear, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", time.Now().Format("2006")))
	dayNum := int(time.Now().Sub(firstOfYear).Hours() / 24)
	his := time.Now().Format("15") + time.Now().Format("4") + time.Now().Format("5")
	return fmt.Sprintf("T%d%d%s%s", storeNum, dayNum, his, din)

}
