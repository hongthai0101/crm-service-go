package middlewares

import (
	"crm-service-go/pkg"
	json "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account map[string][]string `json:"account"`
}

type TokenPayload struct {
	Exp            int            `json:"exp"`
	Iat            int            `json:"iat"`
	Sub            string         `json:"sub"`
	Jti            string         `json:"jti"`
	Iss            string         `json:"iss"`
	Aud            string         `json:"aud"`
	Typ            string         `json:"typ"`
	Azp            string         `json:"azp"`
	RealmAccess    RealmAccess    `json:"realm_access"`
	ResourceAccess ResourceAccess `json:"resource_access"`
	Scope          string         `json:"scope"`
	EmailVerified  bool           `json:"email_verified"`

	jwt.StandardClaims
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication header is missing",
			})
			return
		}

		temp := strings.Split(authHeader, "Bearer")
		if len(temp) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token not Bearer"})
			return
		}
		tokenString := strings.TrimSpace(temp[1])
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is wrong or Expire"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			payload, _ := json.Marshal(claims)
			tokenPayload := TokenPayload{}
			err = json.Unmarshal(payload, &tokenPayload)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
				return
			}
			c.Set(pkg.ProjectKeyUserInfo, tokenPayload)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token is not valid",
		})
		return
	}
}

func LoggedUser(ctx *gin.Context) *TokenPayload {
	if user, ok := ctx.Get(pkg.ProjectKeyUserInfo); ok {
		userInfo := user.(TokenPayload)
		return &userInfo
	}
	return nil
}
