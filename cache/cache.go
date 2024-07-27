package cache

import "context"

type Cache interface {
	Get(ctx context.Context, key string) error
	Set(ctx context.Context, key, value string) error
}
