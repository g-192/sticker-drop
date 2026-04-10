package graph

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Resolver struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}
