package deps

import (
	"chat-server/db"
	"chat-server/db/gen"
	"chat-server/internal/cache"
	"chat-server/internal/repo"
	"chat-server/internal/service"
	"context"
)

type Container struct {
	Services *service.Services
	Cache    *cache.Cache
}

func NewContainer(ctx context.Context) *Container {
	db := db.NewDB(ctx)
	queries := gen.New(db.DBPool)
	repos := repo.NewRepos(queries)
	services := service.NewServices(repos)
	cache := cache.NewCache(ctx)
	return &Container{
		Services: services,
		Cache:    cache,
	}
}
