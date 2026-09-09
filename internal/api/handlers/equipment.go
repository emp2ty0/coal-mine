package api_handlers

import (
	"net/http"

	"github.com/emp2ty0/coal-mine/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *HTTPHandlers) DeviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"device": domain.GetDeviceInfo(),
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

	deviceCost := domain.GetDeviceInfo()[deviceClass.Class]
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
