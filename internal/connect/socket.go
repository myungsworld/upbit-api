package connect

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"upbit-api/config"
)

// Socket 업비트 웹소켓을 이용한 시세 수신
// https://docs.upbit.com/docs/upbit-quotation-websocket
func Socket() {
	conn, _, err := websocket.DefaultDialer.Dial(config.UpbitWebSocketURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	if err := subscribeToMarketData(conn); err != nil {
		log.Fatal("Error subscribing to market data:", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)
	}
}

func subscribeToMarketData(conn *websocket.Conn) error {

	// 마켓들 파싱
	var markets string
	for _, market := range config.Markets {
		markets = markets + "\"" + market + "\"" + ","
	}
	markets = markets[:len(markets)-1]

	subscription := fmt.Sprintf(`[{"ticket":"test"},{"type":"ticker","codes":[%s]}]`, markets)

	//test := "PING"

	// Send the subscription message to the WebSocket server
	err := conn.WriteMessage(websocket.TextMessage, []byte(subscription))
	if err != nil {
		return err
	}

	fmt.Printf("Subscribed to market data for %s\n", markets)
	return nil
}
