package main

import (
	"context"
	"github.com/Michaellqa/iot/web"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	globalCtx, cancel := context.WithCancel(context.Background())
	srv := web.NewServer(globalCtx)

	done := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
		done <- struct{}{}
	}()

	go func() {
		<-sigs
		log.Println("APP: shutting down ...")
		cancel()

		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		srv.Shutdown(ctx)
		os.Exit(1)
	}()

	<-done
	log.Println("APP: finished successfully")
}
