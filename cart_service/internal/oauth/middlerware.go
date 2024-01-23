package oauth

import (
	"store/cart_service/internal/jwt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	BearerPrefix = "Bearer "
)

func Middleware(parser *jwt.JwtParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.Next()
			return
		}

		token, ok := strings.CutPrefix(auth, BearerPrefix)
		if !ok {
			c.Next()
			return
		}

		claims, err := parser.ParseToken(token)
		if err != nil {
			c.Next()
			return
		}

		if expiresAt := time.Unix(claims.ExpiresAt, 0); expiresAt.Before(time.Now()) {
			c.Next()
			return
		}

		c.Set("user", OauthUser{Id: claims.UserId, Roles: claims.UserRoleNames})
		c.Next()
	}
}

func ExtractUser(c *gin.Context) (user OauthUser, ok bool) {
	userAny, ok1 := c.Get("user")

	user, ok2 := userAny.(OauthUser)

	return user, (ok1 && ok2)
}
