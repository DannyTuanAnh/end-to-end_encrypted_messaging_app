package rate_limiter_middleware

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rateLimitMemory struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

var (
	clients     = make(map[string]*rateLimitMemory)
	clientsMu   sync.RWMutex
	cleanupOnce sync.Once
)

func NewRateLimitMemory() *rateLimitMemory {
	return &rateLimitMemory{}
}

func (r *rateLimitMemory) Allow(ctx context.Context, key string) bool {
	startCleanupWorker(ctx)
	return getRateLimiter(key).Allow()
}

func getRateLimiter(ip string) *rate.Limiter {
	// khóa toàn bộ map để tránh race condition
	clientsMu.RLock()
	client, exists := clients[ip]
	clientsMu.RUnlock()

	if !exists {
		clientsMu.Lock()
		defer clientsMu.Unlock()

		// double-check để tránh ghi đè limiter nếu đã có goroutine khác tạo rồi
		client, exists = clients[ip]
		if !exists {
			limiter := rate.NewLimiter(5, 10) // refill 5 request per second, burst size 10
			newClient := &rateLimitMemory{limiter: limiter, lastSeen: time.Now()}
			clients[ip] = newClient

			return limiter
		}
	}

	// chỉ lock client cụ thể để tránh khóa toàn bộ map
	client.mu.Lock()
	client.lastSeen = time.Now()
	client.mu.Unlock()

	return client.limiter
}

func cleanupClients(ctx context.Context) {
	for {

		ticker := time.NewTicker(time.Minute)

		select {
		case <-ticker.C:
			clientsMu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			clientsMu.Unlock()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func startCleanupWorker(ctx context.Context) {
	cleanupOnce.Do(func() {
		go cleanupClients(ctx)
	})
}

// sử dụng ApacheBench để test rate limiting
// ab -n 11 -c 1 -H "X-API-KEY: 4f4c48fb-665a-4e6b-a498-01e72e89db7c" localhost:8080/api/v1/users
