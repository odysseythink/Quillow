package middleware

import (
	"strings"

	i18npkg "github.com/anthropics/firefly-iii-go/pkg/i18n"
	"github.com/gin-gonic/gin"
)

func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := i18npkg.DefaultLocale
		header := c.GetHeader("Accept-Language")
		if header != "" {
			parts := strings.SplitN(header, ",", 2)
			lang := strings.TrimSpace(parts[0])
			lang = strings.SplitN(lang, ";", 2)[0]
			lang = strings.ReplaceAll(lang, "-", "_")
			if lang != "" {
				locale = lang
			}
		}
		c.Set(string(i18npkg.LocaleKey), locale)
		c.Next()
	}
}
