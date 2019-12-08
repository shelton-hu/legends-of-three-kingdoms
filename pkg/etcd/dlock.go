package etcd

import (
	"context"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
)

type callback func() error

func (etcd *Etcd) Dlock(ctx context.Context, key string, cb callback) error {
	logger.Debug(ctx, "Create dlock with key [%s]", key)
	mutex, err := etcd.NewMutex(key)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to create mutex: %+v", err)
		return err
	}
	err = mutex.Lock(ctx)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to lock mutex: %+v", err)
		return err
	}
	defer func() {
		if err := mutex.Unlock(ctx); err != nil {
			logger.Error(ctx, err.Error())
		}
	}()
	err = cb()
	return err
}
