package autoTrading2

import (
	"fmt"
	"log"
	"time"
)

func BidBigShort() {

	setTicker := time.NewTicker(time.Minute * 3)

	for {
		select {
		case <-setTicker.C:
			log.Println("--현재 돌아가고있는 상태값 확인 시작--")
			PreviousMarketMutex.Lock()
			for market, info := range PreviousMarketInfo {
				fmt.Println(market)
				fmt.Println(info)
			}
			PreviousMarketMutex.Unlock()
			log.Println("--현재 돌아가고 있는 상태값 확인 끝--")

		}
	}

}
