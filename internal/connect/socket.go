package connect

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"upbit-api/config"
)

// Socket 업비트 웹소켓을 이용한 시세 수신
// https://docs.upbit.com/docs/upbit-quotation-websocket
func Socket(socketType string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(config.UpbitWebSocketURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}

	if err := subscribeToMarketData(conn, socketType); err != nil {
		log.Fatal("Error subscribing to market data:", err)
	}

	return conn
	// 리턴값에다가 아래 샘플 가져다가 쓰면 됨
	//for {
	//	_, message, err := conn.ReadMessage()
	//	if err != nil {
	//		log.Println("Error reading message:", err)
	//		break
	//	}
	//	fmt.Printf("Received message: %s\n", message)
	//}
}

func subscribeToMarketData(conn *websocket.Conn, socketType string) error {

	// 마켓들 파싱
	var markets string

	for _, market := range config.Markets {
		markets = markets + "\"" + market + "\"" + ","
	}

	markets = markets[:len(markets)-1]

	subscription := fmt.Sprintf(`[{"ticket":"myungsworld"},{"type":"%s","codes":[%s]}]`, socketType, markets)

	//test := "PING"

	// Send the subscription message to the WebSocket server
	err := conn.WriteMessage(websocket.TextMessage, []byte(subscription))
	if err != nil {
		return err
	}

	fmt.Printf("코인 %d개 모니터링 소켓 시작\n", len(config.Markets))
	return nil
}
