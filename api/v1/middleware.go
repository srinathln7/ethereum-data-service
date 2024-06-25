package v1

import "github.com/gin-gonic/gin"

// handleFavicon handles the favicon.ico request without logging.
// Abort further processing to prevent logging
func handleFavicon(c *gin.Context) {
	c.Abort()
}
