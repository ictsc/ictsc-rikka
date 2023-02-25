package rc

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type GrowiSessionCookieRepository struct {
	rc *redis.Client
}

func NewGrowiSessionCookieRepository(rc *redis.Client) *GrowiSessionCookieRepository {
	return &GrowiSessionCookieRepository{rc: rc}
}

var sessionCookieKey = "connect.sid"

func (r *GrowiSessionCookieRepository) Get(ctx context.Context) (string, error) {
	return r.rc.Get(ctx, sessionCookieKey).Result()
}

func (r *GrowiSessionCookieRepository) Set(ctx context.Context, cookie string) error {
	err := r.rc.Set(ctx, sessionCookieKey, cookie, 0).Err()
	return err
}
