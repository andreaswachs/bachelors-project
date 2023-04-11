package api

import service "github.com/andreaswachs/daaukins-service"

func ScheduleLab(lab, serverId string) (*service.ScheduleLabResponse, error) {
	ctx, cancelFunc := newCtx()
	defer cancelFunc()

	return getClient().ScheduleLab(ctx, &service.ScheduleLabRequest{Lab: lab, ServerId: serverId})
}
