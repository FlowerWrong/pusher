package middlewares

import (
	"net/http"

	"github.com/FlowerWrong/pusher"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Signature middleware
func Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.Param("app_id")
		if appID != viper.GetString("APP_ID") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid app id"})
			c.Abort()
			return
		}

		ok, err := pusher.Verify(c.Request)
		if ok {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
}
