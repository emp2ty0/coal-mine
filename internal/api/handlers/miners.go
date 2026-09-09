package api_handlers

import (
	"net/http"

	"github.com/emp2ty0/coal-mine/internal/domain"
	"github.com/gin-gonic/gin"
)

type RequestDTO struct {
	Class string `json:"class"`
}

func (h *HTTPHandlers) MinersWages(c *gin.Context) {
	response := domain.GetWagesInfo()

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

	minerInfo := domain.GetWagesInfo()[minerClass.Class]

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

	var response []domain.Miner

	for _, miner := range h.enterprise.GetMiners() {
		if miner.Class == q {
			response = append(response, *miner)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"miners": response,
	})
}
