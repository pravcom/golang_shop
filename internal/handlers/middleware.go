package handlers

import "github.com/gin-gonic/gin"

func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang") // Из URL ?lang=ru

		// Сохраняем язык в контексте
		c.Set("lang", lang)
		c.Next()
	}
}
