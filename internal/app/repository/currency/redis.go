package currency

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ARTI7876/worker-service/internal/app/entity"
	"github.com/ARTI7876/worker-service/internal/app/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const rateKeyPrefix = "rates:"

// repoRedis — кэш курсов валют поверх Redis.
type repoRedis struct {
	client   *redis.Client
	cacheTTL time.Duration
}

func NewRepoFromRedis(client *redis.Client, cacheTTL time.Duration) repository.CurrencyRate {
	return &repoRedis{
		client:   client,
		cacheTTL: cacheTTL,
	}
}

func (r *repoRedis) buildKey(from, to string) string {
	return fmt.Sprintf("%s%s:%s", rateKeyPrefix, from, to)
}

func (r *repoRedis) GetRate(ctx context.Context, from, to string) (float64, error) {
	val, err := r.client.Get(ctx, r.buildKey(from, to)).Result()
	if errors.Is(err, redis.Nil) {
		return 0, entity.ErrRateNotFound
	}
	if err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("parse cached rate %q: %w", val, err)
	}

	return rate, nil
}

func (r *repoRedis) SetRate(ctx context.Context, from, to string, rate float64) error {
	val := strconv.FormatFloat(rate, 'f', -1, 64)

	return r.client.Set(ctx, r.buildKey(from, to), val, r.cacheTTL).Err()
}

func (r *repoRedis) SetRates(ctx context.Context, from string, rates map[string]float64) error {
	pipe := r.client.Pipeline()

	for to, rate := range rates {
		val := strconv.FormatFloat(rate, 'f', -1, 64)
		pipe.Set(ctx, r.buildKey(from, to), val, r.cacheTTL)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("set rates: %w", err)
	}

	log.Info().
		Ctx(ctx).
		Str("from", from).
		Int("count", len(rates)).
		Msg("Курсы валют сохранены в кэш")

	return nil
}
