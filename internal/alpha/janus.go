package alpha

type Janus struct {
	sourceChan chan Kline
	resultChan chan JanusAlpha
}

func NewJanus() *Janus {
	return &Janus{
		sourceChan: make(chan Kline),
		resultChan: make(chan JanusAlpha),
	}
}

func (j *Janus) Name() string {
	return "Janus"
}

func (j *Janus) SourceChannel() chan<- Kline {
	return j.sourceChan
}

func (j *Janus) Start() {
	for model := range j.sourceChan {
		processedData := j.processModel(model)
		j.resultChan <- processedData
	}
}

func (j *Janus) OutputChannel() <-chan JanusAlpha {
	return j.resultChan
}

func (j *Janus) processModel(kline Kline) JanusAlpha {
	return JanusAlpha{Alpha: []float64{float64(kline.StartTime), float64(kline.EndTime)}}
}
