package pi

import (
	"context"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/config"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/constants"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/etcd"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
)

type Watcher chan *config.GlobalConfig

func WatchGlobalConfig(etcd *etcd.Etcd, watcher Watcher) error {
	ctx := context.Background()
	var globalConfig config.GlobalConfig
	err := etcd.Dlock(ctx, constants.DlockKey, func() error {
		// get value
		get, err := etcd.Get(ctx, constants.GlobalConfigKey)
		if err != nil {
			return err
		}
		// parse value
		if get.Count == 0 {
			logger.Debug(ctx, "Cannot get global config, put the initial string. [%s]", config.InitialGlobalConfig)
			globalConfig = config.DecodeInitConfig()
			_, err = etcd.Put(ctx, constants.GlobalConfigKey, config.InitialGlobalConfig)
			if err != nil {
				return err
			}
		} else {
			globalConfig, err = config.ParseGlobalConfig(get.Kvs[0].Value)
			if err != nil {
				return err
			}
		}
		logger.Debug(ctx, "Global config update to [%+v]", globalConfig)
		// send it back
		watcher <- &globalConfig
		return nil
	})

	// watch
	go func() {
		logger.Debug(ctx, "Start watch global config")
		watchRes := etcd.Watch(ctx, constants.GlobalConfigKey)
		for res := range watchRes {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.PUT {
					//logger.Debug(nil, "Got updated global config from etcd, try to decode with yaml")
					globalConfig, err := config.ParseGlobalConfig(ev.Kv.Value)
					if err != nil {
						logger.Error(ctx, "Watch global config from etcd found error: %+v", err)
					} else {
						watcher <- &globalConfig
					}
				}
			}
		}
	}()
	return err
}
