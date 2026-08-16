package eorder

import (
	"context"
	"fmt"

	"github.com/ARTI7876/worker-service/internal/app/entity"
	eventhandler "github.com/ARTI7876/worker-service/internal/app/handler/event"
	"github.com/ARTI7876/worker-service/internal/pkg/broker"
	"github.com/rs/zerolog/log"
)

type handler struct{}

func NewHandler() eventhandler.OrderCreated {
	return &handler{}
}

func (h *handler) CallbackOrderCreated(
	ctx context.Context,
	ev *entity.EventOrderCreated,
	_ []broker.Header,
) error {
	log.Info().
		Ctx(ctx).
		EmbedObject(ev).
		Msg("Получено событие order.created")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(
			fmt.Errorf("order.created: empty order_guid"),
		)
	}

	return nil
}
