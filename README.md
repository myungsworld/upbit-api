# upbit-api , 자동매매 (Backend)

## 환경 설정
```
.env 파일 생성후 Upbit Key 기입 , GCP gmail 알람 사용시 해당 App 패스워드 입력
AccessKey=""
SecretKey=""
GmailAppPassword=""
```

## 기존 자동매매

1. 매시간 마다 전날 대비 많이 내린 코인 6천원 매수
2. 구매한 코인중 -12퍼가 넘어갈시 6천원 다시 매수 , 20퍼 오를시 해당 마켓 코인 매도
   - -12퍼가 넘어가는 코인은 00시,06시,12시,18시 기준으로 6시간 마다 6천원 매수로 수정 (장이 안좋은 경우 대비)
3. 실시간으로 핫한 코인 알람

- ~~서버 없어서 aws 람다로 돌리려다가 집나간 맥북이 돌아와서 람다 코드는 무시~~

```go
cmd/autoTrading/main.go
```

## 스윙 자동매매

1. 지나간 n일의 저점,종가,고가의 평균을 가져와 종가의 평균과 현재가의 편차가 적을때
   - 전체코인 가져오는 부분(시세 캔들 조회 API) API 초당 30회 제공이라고 적혀있지만 아닌듯? 요청횟수 초과 페일오버 구축
2. 저점의 평균에서 매수대기 이후 매수가 되면 고점의 평균 -1퍼에서 매도 
    - 업비트에서 거래내역 API 를 제공하지 않을시에 상태값이나 데이터베이스 추가개발 필요 (2024-03-25) 
    - 거래내역 제공에 따라 부분 수정 매수대기 이후 코드 수정(2024-03-27)
3. 하루동안 모니터링후 매수대기만 걸린 데이터들 주문리스트에서 삭제후 초기화 및 회귀

```go
cmd/autoTrading2/main.go
```