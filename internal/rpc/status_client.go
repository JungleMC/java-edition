package rpc

import (
	"fmt"
	"github.com/junglemc/Service-StatusProvider/pkg/service"
	"google.golang.org/grpc"
)

var (
	StatusConnection *grpc.ClientConn
	Status           service.StatusProviderClient
)

func StatusInit(address string, port int) {
	var err error
	StatusConnection, err = grpc.Dial(fmt.Sprintf("%v:%v", address, port), grpc.WithInsecure())
	if err != nil {
		panic(err) // TODO: tidy error handling
	}

	Status = service.NewStatusProviderClient(StatusConnection)
}

func StatusClose() {
	if StatusConnection != nil {
		Status = nil
		_ = StatusConnection.Close()
		StatusConnection = nil
	}
}
