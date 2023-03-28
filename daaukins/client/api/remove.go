package api

import "github.com/andreaswachs/bachelors-project/daaukins/service"

func RemoveLab(id string) (*service.RemoveLabResponse, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return getClient().RemoveLab(ctx, &service.RemoveLabRequest{Id: id})
}
