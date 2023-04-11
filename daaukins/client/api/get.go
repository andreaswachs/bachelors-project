package api

import (
	"github.com/andreaswachs/bachelors-project/daaukins/client/config"
	service "github.com/andreaswachs/daaukins-service"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetLabs(serverId string) (*service.GetLabsResponse, error) {
	// if id is "", then all getting all labs is implied
	ctx, cancelFunc := shortLivedCtx()
	defer cancelFunc()

	return getClient().GetLabs(ctx, &service.GetLabsRequest{})
}

func GetLab(id string) (*service.GetLabResponse, error) {
	ctx, cancelFunc := shortLivedCtx()
	defer cancelFunc()

	return getClient().GetLab(ctx, &service.GetLabRequest{Id: id})
}

func GetServers() (*service.GetServersResponse, error) {
	ctx, cancelFunc := shortLivedCtx()
	defer cancelFunc()

	return getClient().GetServers(ctx, &emptypb.Empty{})
}

func GetFrontends() (*service.GetFrontendsResponse, error) {
	ctx, cancelFunc := shortLivedCtx()
	defer cancelFunc()

	response, err := getClient().GetFrontends(ctx, &service.GetFrontendsRequest{})
	if err != nil {
		return &service.GetFrontendsResponse{}, err
	}

	// There might be responses with frontend hosts being 'replace-me' which
	// will be of the Leader server. We will replace this value with the leader's IP
	// that is stored in the dkn cli tool configuration
	for _, frontend := range response.Frontends {
		if frontend.Host == "replace-me" {
			frontend.Host = config.ServerAddress()
		}
	}

	return response, nil
}
