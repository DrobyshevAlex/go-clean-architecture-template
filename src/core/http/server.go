package http

import (
	"context"
	"log"
	v1_0 "main/src/core/http/controllers/v1.0"
	"main/src/core/http/ws"
	user "main/src/features/user/presentation/controllers"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server

	// Controllers
	// v1.0.0
	apiController  *v1_0.ApiController
	userController *user.UserController

	wg             *sync.WaitGroup
	wsShutdownChan chan struct{}
	isReady        bool

	ws *ws.WSHub
}

func (s *Server) Init(addr string, wg *sync.WaitGroup, isDebug bool, wsShutdownChan chan struct{}) {
	s.wsShutdownChan = wsShutdownChan

	log.Println("webServer.init", addr)

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// get new router in init function (not in constructor) because git.SetMode global function
	// and gin.SetMode don't affect to router structure
	s.router = gin.New()
	s.wg = wg
	s.registerValidators()
	s.registerGlobalMiddlewares()
	s.initRoutes()

	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	httpShutdownCh := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Println("webServer.shutdown: init")

		graceCtx, graceCancel := context.WithTimeout(ctx, 1*time.Second)
		defer graceCancel()

		if err := s.srv.Shutdown(graceCtx); err != nil {
			log.Println(err)
		}

		httpShutdownCh <- struct{}{}
	}()

	go func() {
		defer s.wg.Done()

		s.isReady = true
		err := s.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}

		s.isReady = false
		<-httpShutdownCh

		log.Println("webServer.shutdown: complete")
		close(s.wsShutdownChan)
	}()

	s.wg.Add(1)

	return nil
}

func (s *Server) IsReady() bool {
	return s.isReady
}

func (s *Server) registerValidators() {
}

func (s *Server) registerGlobalMiddlewares() {
}

func (s *Server) initRoutes() {
	v1_0 := s.router.Group("/v1.0/")
	{
		v1_0.GET("/", s.apiController.Version)

		users := v1_0.Group("/users")
		{
			users.GET("/:id", s.userController.Profile)
		}
	}

	s.router.GET("/ws", func(c *gin.Context) {
		if s.isReady {
			s.ws.ServeWs(c.Writer, c.Request)
		} else {
			c.Status(503)
		}
	})
}

func NewServer(
	apiController *v1_0.ApiController,
	userController *user.UserController,
	ws *ws.WSHub,
) *Server {
	return &Server{
		apiController:  apiController,
		userController: userController,
		ws:             ws,
	}
}
