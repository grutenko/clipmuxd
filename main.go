package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
)

func mustOpenStorage(file string) *Storage {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=wal", file))
	if err != nil {
		panic(err)
	}
	return &Storage{db}
}

func main() {
	configFilename := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	config := MustLoadConfig(*configFilename)
	storage := mustOpenStorage(config.DatabaseFile)
	defer storage.Conn.Close()

	ctx := AppCtx{
		Config:  config,
		Storage: storage,
	}

	handshakeListen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.HandshakePort))
	if err != nil {
		panic(err)
	}
	handshakeGrpcServer := grpc.NewServer()
	handshakeHandler := NewHandshakeGrpcServer(&ctx)
	RegisterHandshakeServer(handshakeGrpcServer, &handshakeHandler)
	defer handshakeGrpcServer.GracefulStop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		<-sigChan
		log.Println("Signal handled")
		done <- true
	}()

	go func() {
		log.Printf("Handshake server started :%d", config.HandshakePort)
		if err := handshakeGrpcServer.Serve(handshakeListen); err != nil {
			panic(err)
		}
	}()

	<-done
}
