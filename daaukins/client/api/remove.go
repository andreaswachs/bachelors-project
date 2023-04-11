package api

import service "github.com/andreaswachs/daaukins-service"

func RemoveLab(id string) (*service.RemoveLabResponse, error) {
	ctx, cancelFunc := longLivedCtx()
	defer cancelFunc()

	return getClient().RemoveLab(ctx, &service.RemoveLabRequest{Id: id})
}

func RemoveLabs(serverId string) (*service.RemoveLabsResponse, error) {
	ctx, cancelFunc := longLivedCtx()
	defer cancelFunc()

	return getClient().RemoveLabs(ctx, &service.RemoveLabsRequest{ServerId: serverId})
}
