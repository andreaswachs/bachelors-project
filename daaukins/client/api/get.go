package api

import (
	"github.com/andreaswachs/bachelors-project/daaukins/service"
)

func GetLabs() (*service.GetLabsResponse, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return getClient().GetLabs(ctx, &service.GetLabsRequest{})
}

func GetLab(id string) (*service.GetLabResponse, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return getClient().GetLab(ctx, &service.GetLabRequest{Id: id})
}
