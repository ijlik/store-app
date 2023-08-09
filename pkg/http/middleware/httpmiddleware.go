package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
)

type DatabaseData struct {
	UserName string
	Password string
	Host     string
	Port     int
	Database string
}

type SecretData func(ctx context.Context, Id string) (DatabaseData, error)

const (
	authorization = "Authorization"
	signatureX    = "SignatureX"
	tokenData     = "tokenData"
)

func WithAllowedCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Menu-Slug, X-Origin-Path, X-Request-Id")
			c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	}
}
