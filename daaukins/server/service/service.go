package service

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedServiceServer
}

func (s *Server) HaveCapacity(context.Context, *HaveCapacityRequest) (*HaveCapacityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HaveCapacity not implemented")
}
func (s *Server) ScheduleLab(context.Context, *ScheduleLabRequest) (*ScheduleLabResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleLab not implemented")
}
func (s *Server) GetLab(context.Context, *GetLabRequest) (*GetLabResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLab not implemented")
}
func (s *Server) GetLabs(context.Context, *GetLabsRequest) (*GetLabsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLabs not implemented")
}
