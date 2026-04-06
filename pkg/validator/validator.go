package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Bind(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	return validate.Struct(req)
}

func BindQuery(c *gin.Context, req any) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return err
	}
	return validate.Struct(req)
}
