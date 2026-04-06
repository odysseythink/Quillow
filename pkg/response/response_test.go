package response

import (
	"testing"

	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

func TestSingleResource(t *testing.T) {
	attrs := map[string]any{"email": "test@example.com", "blocked": false}
	result := Single("users", "1", attrs)

	data, ok := result["data"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "users", data["type"])
	assert.Equal(t, "1", data["id"])
	assert.Equal(t, attrs, data["attributes"])
}

func TestCollectionResource(t *testing.T) {
	items := []Resource{
		{Type: "users", ID: "1", Attributes: map[string]any{"email": "a@b.com"}},
		{Type: "users", ID: "2", Attributes: map[string]any{"email": "c@d.com"}},
	}
	pg := pagination.Meta{Total: 10, Count: 2, PerPage: 2, CurrentPage: 1, TotalPages: 5}

	result := Collection(items, pg)
	data, ok := result["data"].([]Resource)
	assert.True(t, ok)
	assert.Len(t, data, 2)

	meta, ok := result["meta"].(map[string]any)
	assert.True(t, ok)
	pag := meta["pagination"].(pagination.Meta)
	assert.Equal(t, 10, pag.Total)
}

func TestErrorResponse(t *testing.T) {
	result := Error("Not Found", 404)
	assert.Equal(t, "Not Found", result["message"])
	assert.Equal(t, 404, result["exception"])
}
