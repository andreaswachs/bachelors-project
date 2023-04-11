package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	service "github.com/andreaswachs/daaukins-service"
	"google.golang.org/grpc"
)

var (
	conn     *grpc.ClientConn
	clientWg sync.WaitGroup
)

func init() {
	// This is an attempt to stop a race condition where the CLI is too fast to progress
	clientWg.Add(1)
}

func Initialize(address, port string) error {
	defer clientWg.Done()

	clientConnection, err := connect(address, port)
	if err != nil {
		fmt.Println("Could not connect to server: ", err)
		return err
	}

	conn = clientConnection

	return nil
}

func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func getClient() service.ServiceClient {
	clientWg.Wait()

	if conn == nil {
		panic("connection to server not initialized (should not happen)")
	}

	return service.NewServiceClient(conn)
}

func connect(address, ip string) (*grpc.ClientConn, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return grpc.DialContext(ctx, fmt.Sprintf("%s:%s", address, ip),
		grpc.WithInsecure(),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
}
