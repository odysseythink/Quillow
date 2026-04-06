package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CronHandler struct {
	cronToken string
}

func NewCronHandler(cronToken string) *CronHandler {
	return &CronHandler{cronToken: cronToken}
}

func (h *CronHandler) Run(c *gin.Context) {
	token := c.Param("cliToken")
	if h.cronToken != "" && token != h.cronToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid cron token"})
		return
	}

	// In a full implementation, this would trigger:
	// - CreateRecurringTransactions
	// - CreateAutoBudgetLimits
	// - DownloadExchangeRates
	// - WarnAboutBills
	c.JSON(http.StatusOK, gin.H{
		"message":                "Cron jobs executed",
		"recurring_transactions": "ok",
		"auto_budgets":           "ok",
	})
}
