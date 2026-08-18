package eprocessor

import (
	"context"
	"sync"

	"github.com/ARTI7876/worker-service/internal/app/entity"
	handler "github.com/ARTI7876/worker-service/internal/app/handler/event"
	"github.com/ARTI7876/worker-service/internal/app/processor"
	"github.com/ARTI7876/worker-service/internal/pkg/broker"
	"github.com/rs/zerolog/log"
)

type orderCreatedProc struct {
	h   handler.OrderCreated
	bus broker.Bus[entity.EventOrderCreated]
}

func NewOrderCreatedEventsCatcher(
	h handler.OrderCreated,
	bus broker.Bus[entity.EventOrderCreated],
) processor.Processor {
	return &orderCreatedProc{
		h:   h,
		bus: bus,
	}
}

func (p *orderCreatedProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	if err := p.bus.Subscribe(ctx, wg, p.h.CallbackOrderCreated); err != nil {
		log.Fatal().
			Err(err).
			Str("topic_name", p.bus.QueueName()).
			Msg("Не удалось запустить подписку")
	}

	log.Info().
		Str("topic_name", p.bus.QueueName()).
		Msg("Подписка запущена")
}
