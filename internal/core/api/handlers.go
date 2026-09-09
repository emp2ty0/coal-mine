package core_api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"syscall"

	"github.com/emp2ty0/coal-mine/internal/core/domain"
	"github.com/emp2ty0/coal-mine/internal/features/device"
	"github.com/emp2ty0/coal-mine/internal/features/miners"
	"github.com/gin-gonic/gin"
)

type HTTPHandlers struct {
	enterprise *domain.Enterprise
	logger     *slog.Logger
	ctx        context.Context
}

func NewHTTPHandlers(enterprise *domain.Enterprise, logger *slog.Logger, ctx context.Context) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
		logger:     logger,
		ctx:        ctx,
	}
}

func (h *HTTPHandlers) GetCoal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"coal": h.enterprise.GetBalance(),
	})
}

func (h *HTTPHandlers) MinersWages(c *gin.Context) {
	response := miners.GetWagesInfo()

	c.JSON(http.StatusOK, response)
}

func (h *HTTPHandlers) MinerHire(c *gin.Context) {
	var minerClass RequestDTO

	if err := c.ShouldBindJSON(&minerClass); err != nil {
		h.logger.Error("error parsing json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	minerInfo := miners.GetWagesInfo()[minerClass.Class]

	if h.enterprise.GetBalance() < minerInfo.HireCost {
		h.logger.Warn("balance < cost")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "balance < cost",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coal": h.enterprise.GetBalance(),
	})

	h.enterprise.HireMiner(minerInfo, minerClass.Class)
}

func (h *HTTPHandlers) MinersShow(c *gin.Context) {
	q := c.Query("class")
	if q == "" {
		c.JSON(http.StatusOK, gin.H{
			"miners": h.enterprise.GetMiners(),
		})

		return
	}

	var response []miners.Miner

	for _, miner := range h.enterprise.GetMiners() {
		if miner.Class == q {
			response = append(response, *miner)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"miners": response,
	})
}

func (h *HTTPHandlers) DeviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"device": device.GetDeviceInfo(),
	})
}

func (h *HTTPHandlers) BuyDevice(c *gin.Context) {
	var deviceClass RequestDTO

	if err := c.ShouldBindJSON(&deviceClass); err != nil {
		h.logger.Error("error parsing json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	deviceCost := device.GetDeviceInfo()[deviceClass.Class]
	if h.enterprise.GetBalance() < deviceCost {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "balance < device cost",
		})

		return
	}

	h.enterprise.BuyDevice(deviceClass.Class, deviceCost)

	c.JSON(http.StatusOK, gin.H{
		"balance": h.enterprise.GetBalance(),
	})
}

func (h *HTTPHandlers) ShowDevices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"devices": h.enterprise.GetDevices(),
	})
}

func (h *HTTPHandlers) CancelContext(c *gin.Context) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to find process",
		})

		return
	}

	if err := proc.Signal(syscall.SIGINT); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send signal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "shutdown signal sent",
	})
}
