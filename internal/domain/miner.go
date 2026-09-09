package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MinerInfo struct {
	HireCost    int  `json:"hire_cost"`
	Energy      int  `json:"energy"`
	Yield       int  `json:"yield"`
	IntervalSec int  `json:"interval_sec"`
	Progressive bool `json:"progressive"`
	Progression int  `json:"progression,omitempty"` // только для STRONG
}

type MinerInterface interface {
	Run(ctx context.Context) <-chan int
	Info() MinerInfo
}

type Miner struct {
	ID        int
	Class     string
	MinerInfo MinerInfo
	HiredAt   time.Time
}

func NewMiner(minerInfo MinerInfo, class string) *Miner {
	return &Miner{
		ID:        int(uuid.New().ID()),
		Class:     class,
		MinerInfo: minerInfo,
		HiredAt:   time.Now(),
	}
}

func (m *Miner) Run(ctx context.Context) <-chan int {
	coalChan := make(chan int)
	go func() {
		for m.MinerInfo.Energy > 0 {
			select {
			case <-ctx.Done():
				return
			case coalChan <- m.GetYield():
				m.MinerInfo.Energy--
				if m.MinerInfo.Progressive {
					m.MinerInfo.Yield += m.MinerInfo.Progression
				}
				time.Sleep(time.Duration(m.MinerInfo.IntervalSec) * time.Second)
			}
		}
	}()

	return coalChan
}

func (m *Miner) Info() MinerInfo {
	return m.MinerInfo
}

func (m *Miner) GetYield() int {
	return m.MinerInfo.Yield
}

const (
	ClassSmall  = "SMALL"
	ClassNormal = "NORMAL"
	ClassStrong = "STRONG"
)

func GetWagesInfo() map[string]MinerInfo {
	return map[string]MinerInfo{
		ClassSmall: {
			HireCost:    5,
			Energy:      30,
			Yield:       1,
			IntervalSec: 3,
			Progressive: false,
		},
		ClassNormal: {
			HireCost:    50,
			Energy:      45,
			Yield:       3,
			IntervalSec: 2,
			Progressive: false,
		},
		ClassStrong: {
			HireCost:    450,
			Energy:      60,
			Yield:       10,
			IntervalSec: 1,
			Progressive: true,
			Progression: 3,
		},
	}
}
