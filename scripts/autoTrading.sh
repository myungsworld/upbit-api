# 기존 프로그램 종료
pids=$(pgrep autoTrading)

if [ -z "$pids" ]; then
    echo "실행중인 자동매매가 없습니다."
else
    echo "기존 자동매매 프로세스를 종료합니다."
    kill -9 $pids
fi

# output 제거
rm -rf autoTrading.out

# 빌드
go build -o bin/autoTrading cmd/autoTrading/main.go

# 실행
./bin/autoTrading > autoTrading.out 2>&1 &

newPids=$(pgrep autoTrading)

echo "자동매매 프로그램 실행 프로세스 ID $newPids"