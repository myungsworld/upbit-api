# upbit-api

## 환경 설정
```
.env 파일 생성후 Upbit Key 기입
AccessKey=""
SecretKey=""

# .env 로드 , Market 상태 수집
func init() {
    config.Init()
}
```

## 웹소켓을 이용한 시세 수신()

```go

# 현재가 
conn := connect.Socket(config.Ticker)

# 체결 
conn := connect.Socket(config.Trade)

# 호가
conn := connect.Socket(config.OrderBook)

go func() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	
		trade := models.Trade{}
		ticker := models.Ticker{}[README.md](README.md)
		orderbook := models.OrderBook{}
		
		var anyInterface interface
		
		if err = json.Unmarshal(message, &anyInterface); err != nil {
			panic(err)
            }
        
		// 매수 , 매도 , 모니터링 등등
    }
}()

```