package domain

import (
	"context"
	"sync"
	"time"

	"github.com/emp2ty0/coal-mine/internal/features/device"
	"github.com/emp2ty0/coal-mine/internal/features/miners"
)

type Enterprise struct {
	balance int
	miners  []*miners.Miner
	devices []*device.Device
	ctx     context.Context
	mtx     *sync.Mutex
}

func NewEnterprise(ctx context.Context, mtx *sync.Mutex) *Enterprise {
	return &Enterprise{
		balance: 0,
		miners:  make([]*miners.Miner, 0),
		devices: make([]*device.Device, 0),
		ctx:     ctx,
		mtx:     mtx,
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

func (e *Enterprise) GetMiners() []*miners.Miner {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	result := make([]*miners.Miner, len(e.miners))
	copy(result, e.miners)

	return result
}

func (e *Enterprise) GetDevices() []*device.Device {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	result := make([]*device.Device, len(e.devices))
	copy(result, e.devices)

	return result
}

func (e *Enterprise) GetBalance() int {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	return e.balance
}

func (e *Enterprise) HireMiner(minerInfo miners.MinerInfo, class string) {
	miner := miners.NewMiner(minerInfo, class)
	e.mtx.Lock()
	e.balance -= minerInfo.HireCost
	e.miners = append(e.miners, miner)
	e.mtx.Unlock()
	coalChan := miner.Run(e.ctx)
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
	device := device.NewDevice(class)
	e.devices = append(e.devices, device)
}
