package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func serveHttp(addr string, exitChan <-chan struct{}, cancelFunc context.CancelFunc) error {
	defer func() {
		cancelFunc()
	}()

	srv := &http.Server{
		Addr: addr,
	}

	go func() {
		select {
		case <-exitChan:
			ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx1)
			log.Printf("%s server shutdown... \n", addr)
		}
	}()
	log.Printf("%s server start... \n", addr)
	err := srv.ListenAndServe()
	return err
}

func handleSignal(exitChan <-chan struct{}, cancelFunc context.CancelFunc) error {
	defer func() {
		cancelFunc()
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signalCode := <-sig:
		return fmt.Errorf("get exit signal %d", signalCode)
	case <-exitChan:
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		return serveHttp(":8080", ctx.Done(), cancel)
	})
	group.Go(func() error {
		return serveHttp(":8081", ctx.Done(), cancel)
	})
	group.Go(func() error {
		return handleSignal(ctx.Done(), cancel)
	})

	err := group.Wait()
	log.Println("all server shut down...", err)
}
