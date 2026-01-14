# Upbit Auto Trading

업비트(Upbit) 암호화폐 거래소 API를 활용한 자동매매 시스템

## 목차

- [주요 기능](#주요-기능)
- [시스템 요구사항](#시스템-요구사항)
- [설치](#설치)
- [환경 설정](#환경-설정)
- [사용법](#사용법)
- [프로젝트 구조](#프로젝트-구조)
- [매매 전략](#매매-전략)
- [API 참고](#api-참고)

## 주요 기능

- **스윙 자동매매** - 3일 저점/종가/고가 평균 기반 지정가 매매
- **시간별 자동매매** - 전일 대비 하락 코인 자동 매수
- **실시간 모니터링** - WebSocket 기반 실시간 가격 추적
- **이메일 알림** - Gmail 연동 거래 알림

## 시스템 요구사항

- Go 1.19 이상
- MySQL 8.0 이상
- 업비트 API 키 (Open API 발급 필요)

## 설치

```bash
# 저장소 클론
git clone https://github.com/myungsworld/upbit-api.git
cd upbit-api

# 의존성 설치
go mod download

# 빌드
go build -o bin/autoTrading2 ./cmd/autoTrading2
```

## 환경 설정

프로젝트 루트에 `.env` 파일을 생성하고 다음 내용을 입력:

```env
# 업비트 API 키 (필수)
AccessKey=YOUR_ACCESS_KEY
SecretKey=YOUR_SECRET_KEY

# Gmail 알림 (선택)
GmailAppPassword=YOUR_APP_PASSWORD
```

### 데이터베이스 설정

MySQL에 `upbit` 데이터베이스 생성:

```sql
CREATE DATABASE upbit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 사용법

### 스윙 자동매매 (권장)

```bash
go run cmd/autoTrading2/main.go
```

### 시간별 자동매매

```bash
go run cmd/autoTrading/main.go
```

### 전체 매도

```bash
go run cmd/askAll/main.go
```

## 프로젝트 구조

```
upbit-api/
├── cmd/                    # 실행 프로그램
│   ├── autoTrading/        # 시간별 자동매매
│   ├── autoTrading2/       # 스윙 자동매매 (주력)
│   ├── autoTrading3/       # 스윙 자동매매 개선 버전
│   └── askAll/             # 전체 매도
├── config/                 # 환경 설정
├── internal/
│   ├── api/                # 업비트 API 통신
│   │   ├── accounts/       # 계좌 조회
│   │   ├── candle/         # 캔들 데이터
│   │   └── orders/         # 주문 처리
│   ├── connect/            # WebSocket 연결
│   ├── datastore/          # 데이터베이스 (MySQL/GORM)
│   ├── handlers/           # 매매 로직
│   ├── middlewares/        # JWT 인증
│   └── models/             # 데이터 모델
└── scripts/                # 유틸리티 스크립트
```

## 매매 전략

### 스윙 자동매매 (autoTrading2)

1. **초기화** (매일 09:00:01 KST)
   - 전체 코인의 3일 저점/종가/고가 평균 계산
   - 상태값 메모리에 저장

2. **매수 조건**
   - 종가 평균이 저점 평균보다 높고 고가 평균보다 낮음
   - 시작가와 종가 평균 편차 0.5% 이내
   - 고가-저가 평균 차이 5% 이상

3. **매수 실행**
   - 저점 평균 가격에 지정가 매수 주문

4. **매도 실행**
   - 매수 체결 시 고가 평균 -1% 가격에 지정가 매도 주문

5. **일일 정리** (매일 08:55 KST)
   - 미체결 매수 주문 취소
   - 체결된 매수 건은 시장가 매도

### 시간별 자동매매 (autoTrading)

- 매시간 전일 대비 하락폭이 큰 코인 6,000원 매수
- 손실률 -12% 초과 시 6시간마다 추가 매수
- 수익률 +20% 도달 시 전량 매도

## API 참고

### 사용 중인 업비트 API

| 엔드포인트 | 용도 |
|-----------|------|
| `GET /v1/accounts` | 계좌 잔고 조회 |
| `POST /v1/orders` | 주문 생성 |
| `DELETE /v1/order` | 주문 취소 |
| `GET /v1/candles/days` | 일봉 캔들 조회 |
| `GET /v1/market/all` | 마켓 목록 조회 |

### WebSocket

- `wss://api.upbit.com/websocket/v1` - 실시간 시세 수신

## 제외 코인

다음 코인은 가격 변동성이 커서 매매 대상에서 제외:

- KRW-BTC (비트코인)
- KRW-BTT (비트토렌트)
- KRW-SHIB (시바이누)
- KRW-XEC

## 주의사항

- 실제 자금이 거래되므로 충분한 테스트 후 사용
- API 요청 제한 (초당 30회) 준수
- 투자 손실에 대한 책임은 사용자에게 있음

## 라이선스

MIT License
