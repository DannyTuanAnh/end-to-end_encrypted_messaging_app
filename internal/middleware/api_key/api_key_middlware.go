package api_key_middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"sync"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
)

var (
	api_keys  = make(map[string]bool)
	apiKeysMu sync.RWMutex
)

func ApiKeyMiddleware(db sqlc.Querier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-API-KEY"})
			return
		}

		hash := utils.HashAPIKey(apiKey)

		// Read cache
		apiKeysMu.RLock()
		cached := api_keys[hash]
		apiKeysMu.RUnlock()

		if !cached {
			// Check if the API key exists in the database
			isActive, err := db.ValidateAPIKey(ctx, hash)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
					return
				}

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check API key"})
				return
			}

			if !isActive {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Inactive API key"})
				return
			}

			// Cache the valid API key hash in memory for future requests
			apiKeysMu.Lock()
			api_keys[hash] = true
			apiKeysMu.Unlock()
		}

		ctx.Next()
	}
}

func InvalidateAPIKey(hash string) {
	apiKeysMu.Lock()
	delete(api_keys, hash)
	apiKeysMu.Unlock()
}

func InvalidateAllAPIKeys() {
	apiKeysMu.Lock()
	api_keys = make(map[string]bool)
	apiKeysMu.Unlock()
}
