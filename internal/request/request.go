package request

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type HTTPClient struct {
	Client  *http.Client
	Limiter Limiters
}

type Limiters struct {
	Limiter      *rate.Limiter
	BurstLimiter *rate.Limiter
}

func (l *Limiters) IsAllowed() bool {
	limited := l.Limiter.Allow()
	burstLimited := l.BurstLimiter.Allow()
	return !limited && !burstLimited
}

var httpClientInstance *HTTPClient
var once sync.Once

func GetHTTPClient() *HTTPClient {
	once.Do(func() {
		httpClientInstance = &HTTPClient{
			Client: &http.Client{},
			Limiter: Limiters{
				Limiter:      rate.NewLimiter(rate.Limit(20), 20),
				BurstLimiter: rate.NewLimiter(rate.Limit(100), 100),
			},
		}
	})
	return httpClientInstance
}
