package main

import (
	"net"

	"github.com/rs/zerolog"
  "github.com/rs/zerolog/log"

	"tsai.eu/solar/controller/gRPC"
)

//------------------------------------------------------------------------------

func main() {
	// be verbose
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// open TCP port 10000
	lis, err := net.Listen("tcp", ":10000")
	if err != nil {
    log.Fatal().Str("Context", "main").Str("Module", "controller").Msg("failed to listen:" + err.Error())
	}

	// create a gRPC server
	grpcServer, _ := gRPC.NewController()

	// register controller and start listening
	grpcServer.Serve(lis)
}

//------------------------------------------------------------------------------
