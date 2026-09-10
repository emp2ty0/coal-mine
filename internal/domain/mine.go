package domain

import (
	"context"
	"sync"
	"time"
)

type Enterprise struct {
	balance int
	miners  []*Miner
	devices []*Device
	mtx     sync.Mutex
}

func NewEnterprise() *Enterprise {
	return &Enterprise{
		balance: 0,
		miners:  make([]*Miner, 0),
		devices: make([]*Device, 0),
	}
}

func (e *Enterprise) Run() {
	for {
		e.mtx.Lock()
		e.balance++
		e.mtx.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (e *Enterprise) GetMiners() []*Miner {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	result := make([]*Miner, len(e.miners))
	copy(result, e.miners)

	return result
}

func (e *Enterprise) GetDevices() []*Device {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	result := make([]*Device, len(e.devices))
	copy(result, e.devices)

	return result
}

func (e *Enterprise) GetBalance() int {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	return e.balance
}

func (e *Enterprise) HireMiner(minerInfo MinerInfo, class string, ctx context.Context) {
	miner := NewMiner(minerInfo, class)
	e.mtx.Lock()
	e.balance -= minerInfo.HireCost
	e.miners = append(e.miners, miner)
	e.mtx.Unlock()
	coalChan := miner.Run(ctx)
	go func() {
		for coal := range coalChan {
			e.mtx.Lock()
			e.balance += coal
			e.mtx.Unlock()
		}
	}()
}

func (e *Enterprise) BuyDevice(class string, cost int) {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.balance -= cost
	device := NewDevice(class)
	e.devices = append(e.devices, device)
}
