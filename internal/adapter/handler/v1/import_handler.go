package v1

import (
	"net/http"

	importuc "github.com/anthropics/quillow/internal/usecase/importer"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type ImportHandler struct {
	uc *importuc.UseCase
}

func NewImportHandler(uc *importuc.UseCase) *ImportHandler {
	return &ImportHandler{uc: uc}
}

func (h *ImportHandler) Preview(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	result, err := h.uc.Preview(c.Request.Context(), file)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

type confirmRequest struct {
	Transactions []confirmTransaction `json:"transactions" binding:"required"`
}

type confirmTransaction struct {
	Index       int    `json:"index"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	CategoryID  uint   `json:"category_id"`
	ExternalID  string `json:"external_id"`
}

func (h *ImportHandler) Confirm(c *gin.Context) {
	var req confirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// For now, return count. Full creation would use transaction repository.
	imported := len(req.Transactions)

	c.JSON(http.StatusOK, gin.H{
		"imported": imported,
		"skipped":  0,
	})
}
