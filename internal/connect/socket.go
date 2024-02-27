package connect

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"upbit-api/config"
)

// Socket 업비트 웹소켓을 이용한 시세 수신
// https://docs.upbit.com/docs/upbit-quotation-websocket
func Socket(socketType string, code ...string) *websocket.Conn {

	// socketType 은 config 패키지의
	// Ticker    = "ticker"
	// Trade     = "trade"
	// OrderBook = "orderbook"
	// 중 하나

	conn, _, err := websocket.DefaultDialer.Dial(config.UpbitWebSocketURL, nil)
	if err != nil {
		log.Println("Error connecting to WebSocket:", err)
		return nil
		//log.Fatal("Error connecting to WebSocket:", err)
	}

	var subscription string
	var markets string

	if len(code) == 0 {
		// code 없을시 기본 config 값에서 가져옴
		for _, market := range config.Markets {
			markets = markets + "\"" + market + "\"" + ","
		}

	} else {
		// code 있을시 해당 마켓만 조회
		for _, market := range code {
			markets = markets + "\"" + market + "\"" + ","
		}
	}

	markets = markets[:len(markets)-1]
	subscription = fmt.Sprintf(`[{"ticket":"myungsworld"},{"type":"%s","codes":[%s]}]`, socketType, markets)

	// Send the subscription message to the WebSocket server
	err = conn.WriteMessage(websocket.TextMessage, []byte(subscription))
	if err != nil {
		log.Fatal("Error subscribing to market data:", err)
	}

	return conn
}
