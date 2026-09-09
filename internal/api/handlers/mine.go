package api_handlers

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandlers) GetCoal(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"coal": h.enterprise.GetBalance(),
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
