package chat

import (
	"chat-server/internal/constants"
	"chat-server/internal/service"
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type PersistWorkerPool struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	persistChan chan []redis.XMessage
	svc         *service.Services
	workers     []*PersistWorker
	mu          sync.Mutex
}

func NewPersistWorkerPool(svc *service.Services) *PersistWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &PersistWorkerPool{
		ctx:         ctx,
		ctxCancel:   cancel,
		persistChan: make(chan []redis.XMessage),
		svc:         svc,
		workers:     []*PersistWorker{},
	}
	p.run()
	return p
}

func (wp *PersistWorkerPool) addWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	w := NewPersistWorker(wp.ctx, wp.persistChan, wp.svc)
	wp.workers = append(wp.workers, w)
}

func (wp *PersistWorkerPool) removeWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.workers[len(wp.workers)-1].ctxCancel()
	wp.workers = wp.workers[:len(wp.workers)-1]
}

func (wp *PersistWorkerPool) workerCount() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	return len(wp.workers)
}

func (wp *PersistWorkerPool) handleCacheStream() {
	for {
		stream, err := wp.svc.Message.GetCacheMessageStream(wp.ctx, 10, 0)
		if err != nil {
			log.Println("error:", err)
			continue
		}
		for _, msgs := range stream {
			wp.persistChan <- msgs.Messages
		}
	}
}

func (wp *PersistWorkerPool) handleWorkerLoad() {
	wp.addWorker()

	ticker := time.NewTicker(constants.PERSIST_PENDING_CHECK_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-ticker.C:
			pending, err := wp.svc.Message.GetPendingCacheMesages(wp.ctx)
			if err != nil {
				log.Println("error:", err)
			}
			if pending.Count > constants.HIGH_PENDING_COUNT {
				wp.addWorker()
			} else if pending.Count < constants.LOW_PENDING_COUNT && wp.workerCount() > constants.MIN_NUM_OF_WORKERS {
				wp.removeWorker()
			}
		}
	}
}

func (wp *PersistWorkerPool) run() {
	go wp.handleCacheStream()
	go wp.handleWorkerLoad()
}
