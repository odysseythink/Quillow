package v1

import (
	"net/http"

	aiuc "github.com/anthropics/quillow/internal/usecase/ai"
	"github.com/anthropics/quillow/pkg/response"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	uc *aiuc.UseCase
}

func NewAIHandler(uc *aiuc.UseCase) *AIHandler {
	return &AIHandler{uc: uc}
}

type suggestRequest struct {
	Description string `json:"description" binding:"required"`
}

func (h *AIHandler) Suggest(c *gin.Context) {
	var req suggestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.uc.Suggest(c.Request.Context(), userID, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

type classifyBatchRequest struct {
	TransactionIDs []uint `json:"transaction_ids" binding:"required,min=1"`
}

func (h *AIHandler) ClassifyBatch(c *gin.Context) {
	var req classifyBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")

	type batchResult struct {
		TransactionID uint   `json:"transaction_id"`
		CategoryID    *uint  `json:"category_id"`
		Source        string `json:"source"`
	}

	var results []batchResult
	classified := 0
	skipped := 0

	for _, txID := range req.TransactionIDs {
		// For batch, we'd need to load transaction description — simplified here
		_ = txID
		_ = userID
		skipped++
		results = append(results, batchResult{TransactionID: txID, CategoryID: nil, Source: "none"})
	}

	c.JSON(http.StatusOK, gin.H{
		"classified": classified,
		"skipped":    skipped,
		"results":    results,
	})
}

type learnRequest struct {
	Pattern    string `json:"pattern" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
}

func (h *AIHandler) Learn(c *gin.Context) {
	var req learnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := h.uc.LearnPattern(c.Request.Context(), userID, req.Pattern, req.CategoryID, req.TagIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pattern learned"})
}

type chatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.uc.Chat(c.Request.Context(), userID, req.Message)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AIHandler) Insight(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.uc.Chat(c.Request.Context(), userID, req.Message)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
