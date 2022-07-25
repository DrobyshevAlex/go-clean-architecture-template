package app

import (
	"context"
	"log"
	"main/src/core/config"
	"main/src/core/db"
	"main/src/core/http"
	"main/src/core/http/ws"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	webServer      *http.Server
	ws             *ws.WSHub
	db             *db.Db
	config         *config.Config
	wg             sync.WaitGroup
	sigs           chan os.Signal
	wsShutdownChan chan struct{}
}

func (a *Application) Init(ctx context.Context) error {
	log.Printf("application: init")

	a.sigs = make(chan os.Signal, 1)
	a.wsShutdownChan = make(chan struct{}, 1)

	err := a.config.Init()
	if err != nil {
		return err
	}

	signal.Notify(a.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.db.Init(a.config.GetDbConfig())

	a.webServer.Init(
		a.config.GetWebServerAddress(),
		&a.wg,
		a.config.IsDebug(),
		a.wsShutdownChan,
	)

	return nil
}

func (a *Application) Run(ctx context.Context) error {
	log.Printf("application.run: start")

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go a.processSignals(cancelFunc)

	a.db.Run()

	err := a.webServer.Run(cancelCtx)
	if err != nil {
		return err
	}

	a.wg.Add(1)

	a.ws.Run(cancelCtx, &a.wg)

	go func() {
		<-a.wsShutdownChan
		a.wg.Done()
	}()

	log.Println("application.run: running")

	a.wg.Wait()

	log.Println("application: graceful shutdown.")

	return nil
}

func (a *Application) processSignals(cancelFun context.CancelFunc) {
	select {
	case <-a.sigs:
		log.Println("application: received shutdown signal from OS")
		cancelFun()
		break
	}
}

func NewApplication(
	config *config.Config,
	db *db.Db,
	webServer *http.Server,
	ws *ws.WSHub,
) *Application {
	return &Application{
		config:    config,
		webServer: webServer,
		ws:        ws,
		db:        db,
	}
}
