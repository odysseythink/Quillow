package response

import (
	"net/http"

	"github.com/anthropics/quillow/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type Resource struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes map[string]any `json:"attributes"`
}

func Single(resourceType, id string, attributes map[string]any) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"type":       resourceType,
			"id":         id,
			"attributes": attributes,
		},
	}
}

func Collection(items []Resource, pg pagination.Meta) map[string]any {
	return map[string]any{
		"data": items,
		"meta": map[string]any{
			"pagination": pg,
		},
	}
}

func Error(message string, code int) map[string]any {
	return map[string]any{
		"message":   message,
		"exception": code,
	}
}

func JSON(c *gin.Context, status int, data map[string]any) {
	c.JSON(status, data)
}

func JSONApi(c *gin.Context, status int, data map[string]any) {
	c.Header("Content-Type", "application/vnd.api+json")
	c.JSON(status, data)
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Error(message, http.StatusNotFound))
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Error(message, http.StatusUnauthorized))
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(message, http.StatusForbidden))
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Error(message, http.StatusBadRequest))
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
