package service

import (
	context "context"
)

type Server struct {
	UnimplementedServiceServer
}

func (s *Server) HaveCapacity(context context.Context, request *HaveCapacityRequest) (*HaveCapacityResponse, error) {
	return HaveCapacity(context, request)
}
func (s *Server) ScheduleLab(context context.Context, request *ScheduleLabRequest) (*ScheduleLabResponse, error) {
	return ScheduleLab(context, request)
}
func (s *Server) GetLab(context context.Context, request *GetLabRequest) (*GetLabResponse, error) {
	return GetLab(context, request)
}
func (s *Server) GetLabs(context context.Context, request *GetLabsRequest) (*GetLabsResponse, error) {
	return GetLabs(context, request)
}
