package chat

import (
	"chat-server/internal/cache"
	"chat-server/internal/service"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type PersistWorker struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	persistChan chan []redis.XMessage
	svc         *service.Services
}

func NewPersistWorker(parentCtx context.Context, persistChan chan []redis.XMessage, svc *service.Services) *PersistWorker {
	ctx, cancel := context.WithCancel(parentCtx)
	p := &PersistWorker{
		ctx:         ctx,
		ctxCancel:   cancel,
		persistChan: persistChan,
		svc:         svc,
	}
	p.run()
	return p
}

func (p *PersistWorker) handlePersist() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case msgs, ok := <-p.persistChan:
			if !ok {
				return
			}
			cachedMsgs, streamIDs, err := cache.BulkStreamToCacheMessages(msgs)
			if err != nil {
				log.Println("error:", err)
				continue
			}
			if len(cachedMsgs) == 0 || len(cachedMsgs) == len(msgs) {
				if err := p.svc.Message.FlushCacheBatchMessagesToDB(p.ctx, cachedMsgs); err != nil {
					log.Println("error:", err)
					continue
				}
				if err := p.svc.Message.MarkCacheBatchAsPersisted(p.ctx, streamIDs); err != nil {
					log.Println("error:", err)
					continue
				}
			}
		}
	}
}

func (p *PersistWorker) run() {
	go p.handlePersist()
}
