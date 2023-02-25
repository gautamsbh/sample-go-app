package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gautamsbh/sample-go-app/config"
	"github.com/gautamsbh/sample-go-app/shared"
	"github.com/gautamsbh/sample-go-app/user"
)

func main() {
	var (
		ctx  = context.Background()
		addr = fmt.Sprintf("%s:%d", config.AppConfig.Host, config.AppConfig.Port)
	)

	var (
		genericRouter = shared.NewGenericRouter()
		userService   = user.NewService()
		userHandler   = user.NewHandler(userService)
		server        = &http.Server{
			Addr:    addr,
			Handler: nil,
		}
	)

	// register routes in generic router
	userHandler.RegisterRoutes(genericRouter)

	// direct all request to generic router
	http.Handle("/", genericRouter)

	// start http server
	go func() {
		log.Print("Server running on port: ", config.AppConfig.Port)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Print("Server closed")
			return
		}
		if err != nil {
			log.Print("Server failed to start")
			os.Exit(1)
		}
	}()

	// listen for interrupt or terminate signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	// graceful shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Print("Some error occurred while server shutdown op")
	}
}
