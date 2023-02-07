package middleware

import (
	"net/http"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/repository"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/adapter/service"

	"github.com/gin-gonic/gin"
)

// AccessTokenAuthentication middleware is a security middleware that check is the requested user if logged in.
// To check is the user is logged, this middleware use a jwt header token X-Access-Token. This token must be valid (not
// expired) and parsed by secret key.
// If the token is not valid, the user is rejected and http.StatusUnauthorized is sent back.
// Otherwise, the userID from the jwt claims is injected into the request context.
func AccessTokenAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqAccessToken := c.Request.Header.Get("X-Access-Token")
		token, err := service.ValidateJWTToken(reqAccessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		claims := service.ExtractJWTClaims(token)
		if claims["type"].(string) != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "cannot use refresh token as access token"})
			return
		}

		c.Set("userID", claims["userID"].(string))
		c.Next()
	}
}

// RefreshTokenAuthentication middleware is a security middleware that check is the requested user has a valid long
// session token.
// This middleware use a jwt header token X-Refresh-Token. This token must be valid (not
// expired) and parsed by secret key.
// If the token is not valid, the user is rejected and http.StatusUnauthorized is sent back.
// Otherwise, the userID from the jwt claims is injected into the request context.
func RefreshTokenAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqRefreshToken := c.Request.Header.Get("X-Refresh-Token")
		token, err := service.ValidateJWTToken(reqRefreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		claims := service.ExtractJWTClaims(token)
		if claims["type"].(string) != "refresh" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "cannot use refresh token as access token"})
			return
		}

		c.Set("userID", claims["userID"].(string))
		c.Next()
	}
}

// AdminOnly middleware checks the user role that match admin value
func AdminOnly(repositories repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		userFound, err := repositories.UserRepository.FindByIDOrFail(claims.UserID)
		if err != nil {
			return
		}

		if userFound.Role != model.RoleAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "member role cannot do this operation"})
			c.Abort()
		}
	}
}
