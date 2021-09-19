package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	//"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type ServerHandler struct{}

func signalf() error {
	signalCh := make(chan os.Signal)

	signal.Notify(signalCh)
	s := <-signalCh
	fmt.Println("catch system signal", s)
	switch s {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		return errors.New("Found return signal, exit. ")
	default:
		fmt.Println("Unknow signal received.")
	}

	return nil
}

func (handler ServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {

	// use errgroup with context
	group, errctx := errgroup.WithContext(context.Background())

	// start some goroutine.
	s := new(ServerHandler)
	myServer := http.Server{
		Handler: s,
		Addr:    ":9999",
	}

	group.Go(func() error {
		defer fmt.Println("GoRoutine 1: Now stop listening for requests.")
		return myServer.ListenAndServe()
	})

	// stop server
	group.Go(func() error {
		select {
		case <-errctx.Done():
			fmt.Println("GoRoutine 2: Now shutdown the server.")
			return myServer.Shutdown(errctx)
		}
		return nil
	})

	// add signal interrupt operation
	group.Go(func() error {
		err := signalf()
		if err != nil {
			fmt.Println("GoRoutine 3: Received exit signal.")
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("All goroutines are dead, get errors: ", err)
	}
}
