package auth

import "github.com/gin-gonic/gin"

func GetCurrentUser(c *gin.Context) (*Claims, bool) {
	userAny, exists := c.Get(ContextUserKey)
	if !exists {
		return nil, false
	}

	claims, ok := userAny.(*Claims)
	if !ok {
		return nil, false
	}

	return claims, true
}
