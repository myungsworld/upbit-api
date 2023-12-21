package autoTrading

var Queue *queue

type queue struct {
	HealthCheckChannel   chan bool
	BidEveryHourChannel  chan bool
	AveragingDownChannel chan bool
	DetectNewCoinChannel chan bool
	ProfitMarginChannel  chan bool

	workingChannel chan bool
}

func NewQueue() *queue {
	healthCheckChannel := make(chan bool)
	bidEveryHourChannel := make(chan bool)
	averagingDownChannel := make(chan bool)
	detectNewCoinChannel := make(chan bool)
	profitMarginChannel := make(chan bool)

	workingChannel := make(chan bool, 6)

	return &queue{
		HealthCheckChannel:   healthCheckChannel,
		BidEveryHourChannel:  bidEveryHourChannel,
		AveragingDownChannel: averagingDownChannel,
		DetectNewCoinChannel: detectNewCoinChannel,
		ProfitMarginChannel:  profitMarginChannel,
		workingChannel:       workingChannel,
	}
}

func (q *queue) Work() {
	//defer func() {
	//	if r := recover(); r != nil {
	//
	//		go q.Work()
	//	}
	//}()

}
