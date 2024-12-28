package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InjectDatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
