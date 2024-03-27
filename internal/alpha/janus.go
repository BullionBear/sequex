package alpha

import "github.com/BullionBear/crypto-trade/internal/model"

type Janus struct {
	sourceChan chan *model.Kline
	resultChan chan JanusAlpha
}

func NewJanus() *Janus {
	return &Janus{
		sourceChan: make(chan *model.Kline),
		resultChan: make(chan JanusAlpha),
	}
}

func (j *Janus) Name() string {
	return "Janus"
}

func (j *Janus) SourceChannel() chan<- *model.Kline {
	return j.sourceChan
}

func (j *Janus) Start() {
	for model := range j.sourceChan {
		processedData := j.processModel(*model)
		j.resultChan <- processedData
	}
}

func (j *Janus) OutputChannel() <-chan JanusAlpha {
	return j.resultChan
}

func (j *Janus) processModel(kline model.Kline) JanusAlpha {
	return JanusAlpha{Alpha: []float64{float64(kline.StartTime), float64(kline.EndTime)}}
}
