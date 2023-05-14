package api

import service "github.com/andreaswachs/daaukins-service"

func CreateLab(lab, serverId string) (*service.CreateLabResponse, error) {
	ctx, cancelFunc := longLivedCtx()
	defer cancelFunc()

	return getClient().CreateLab(ctx, &service.CreateLabRequest{Lab: lab, ServerId: serverId})
}
