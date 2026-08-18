package handler

import (
	"context"

	"github.com/ARTI7876/worker-service/internal/app/entity"
	"github.com/ARTI7876/worker-service/internal/pkg/broker"
)

type (
	OrderCreated interface {
		CallbackOrderCreated(
			ctx context.Context,
			ev *entity.EventOrderCreated,
			headers []broker.Header,
		) error
	}
)
