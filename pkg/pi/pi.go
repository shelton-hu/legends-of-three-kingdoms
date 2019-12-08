package pi

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/config"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/constants"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/etcd"
	"github.com/shelton-hu/legends-of-three-kingdoms/pkg/logger"
	prisma "github.com/shelton-hu/legends-of-three-kingdoms/pkg/prisma/mysql-prisma-client"
)

var (
	ErrPrismaEndpoint  = errors.New("Prisma's endpoint is unknown")
	ErrUsedWrongPrisma = errors.New("Used wrong prisma")
)

type globalCfgWatcher func(*config.GlobalConfig)

type Pi struct {
	cfg              *config.Config
	globalCfg        *config.GlobalConfig
	globalCfgWatcher []globalCfgWatcher
	etcd             *etcd.Etcd
	mysqlPrisma      *prisma.Client
}

var global *Pi
var mutex sync.RWMutex
var globalMutex sync.RWMutex

func SetGlobal(cfg *config.Config) {
	globalMutex.Lock()
	global = NewPi(cfg)
	globalMutex.Unlock()
}

func Global() *Pi {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return global
}

func NewPi(cfg *config.Config) *Pi {
	p := &Pi{cfg: cfg}
	p.openPrisma()
	p.openEtcd()
	p.watchGlobalCfg()
	return p
}

func (p *Pi) GlobalConfig() *config.GlobalConfig {
	mutex.RLock()
	globalCfg := p.globalCfg
	mutex.RUnlock()
	return globalCfg
}

func (p *Pi) Cfg(ctx context.Context) *config.Config {
	return p.cfg
}

func (p *Pi) Etcd(ctx context.Context) *etcd.Etcd {
	return p.etcd
}

func (p *Pi) MysqlPrisma(ctx context.Context) *prisma.Client {
	if !strings.HasSuffix(p.cfg.Prisma.MysqlEndpoint, constants.MysqlPrismaEndpointSuffix) {
		logger.Error(ctx, ErrUsedWrongPrisma.Error())
		return nil
	}
	return p.mysqlPrisma
}

func (p *Pi) TheardWatchGlobalConfig(cb globalCfgWatcher) {
	p.globalCfgWatcher = append(p.globalCfgWatcher, cb)
}

func (p *Pi) watchGlobalCfg() *Pi {
	ctx := context.Background()
	watcher := make(Watcher)

	go func() {
		err := WatchGlobalConfig(p.Etcd(ctx), watcher)
		if err != nil {
			logger.Critical(ctx, "failed to watch global config")
			panic(err)
		}
	}()

	globalCfg := <-watcher
	p.setGlobalCfg(globalCfg)
	logger.Debug(ctx, "Pi not global config: [%+v]", p.globalCfg)

	go func() {
		for globalCfg := range watcher {
			p.setGlobalCfg(globalCfg)
			logger.Debug(ctx, "Pi not global config: [%+v]", p.globalCfg)
		}
	}()

	return p
}

func (p *Pi) setGlobalCfg(globalCfg *config.GlobalConfig) {
	mutex.Lock()
	p.globalCfg = globalCfg
	for _, cb := range p.globalCfgWatcher {
		go cb(globalCfg)
	}
	mutex.Unlock()
}

func (p *Pi) openPrisma() *Pi {
	if p.cfg.Prisma.Disable {
		return p
	}
	mysqlEndpoint := p.cfg.Prisma.MysqlEndpoint
	p.mysqlPrisma = prisma.New(&prisma.Options{
		Endpoint: mysqlEndpoint,
	})

	return p
}

func (p *Pi) openEtcd() *Pi {
	endpoints := strings.Split(p.cfg.Etcd.Endpoints, ",")
	e, err := etcd.Connect(endpoints, constants.ConfigEtcdPrefix)
	if err != nil {
		logger.Critical(context.Background(), "failed to connect etcd")
		panic(err)
	}
	p.etcd = e
	return p
}
