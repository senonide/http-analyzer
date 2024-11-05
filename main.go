package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/senonide/http-analyzer/pkg/rawhttp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	cert := os.Getenv("CERT")
	key := os.Getenv("KEY")
	ja3 := ""
	tlsConfig, err := rawhttp.LoadTLSConfig(cert, key, &ja3)
	if err != nil {
		log.Fatalf("error loading TLS config: %v", err)
	}

	port := os.Getenv("PORT")
	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		log.Fatalf("error starting TCP server: %v", err)
		return
	}
	defer listener.Close()
	log.Printf("listening on \033[34mhttps://localhost:%s\033[0m", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}
		reqData := new(rawhttp.RequestData)
		reqData.JA3 = ja3
		go reqData.HandleTcpConnection(conn)
	}
}
