package cmd

import (
	"github.com/ARTI7876/worker-service/internal/app/builder"
	"github.com/urfave/cli/v2"
)

func SubscribeOrderCreated() *cli.Command {
	return &cli.Command{
		Name:            "subscribe-order-created",
		Aliases:         []string{"consume-order-created"},
		Usage:           "Запускает consumer топика order.created",
		Action:          cmdSubscribeOrderCreated,
		HideHelpCommand: true,
	}
}

func cmdSubscribeOrderCreated(cCtx *cli.Context) error {
	app := builder.NewBuilder(cCtx)

	app.BuildConfig()
	app.BuildMonitorOpenTelemetry()
	app.BuildConnRedis()
	app.BuildClientFixer()
	app.BuildRepoCurrencyRate()
	app.BuildServiceCurrency()
	app.BuildBrokerKafka()
	app.BuildConsumerOrderCreated()
	app.BuildProcHttp()

	app.Run()

	return nil
}
