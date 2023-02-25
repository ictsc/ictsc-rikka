package repository

import "context"

type GrowiSessionCookieRepository interface {
	Set(context.Context, string) error
	Get(ctx context.Context) (string, error)
}
