package api

import "github.com/andreaswachs/bachelors-project/daaukins/service"

func ScheduleLab(lab string) (*service.ScheduleLabResponse, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return getClient().ScheduleLab(ctx, &service.ScheduleLabRequest{Lab: lab})
}
