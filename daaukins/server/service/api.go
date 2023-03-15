package service

import (
	context "context"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func HaveCapacity(context context.Context, request *HaveCapacityRequest) (*HaveCapacityResponse, error) {
	if request.Lab == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lab is empty")
	}

	tempFile, err := os.CreateTemp("", "daaukins-lab.yaml")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(request.Lab)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write lab to temp file: %v", err)
	}
	defer tempFile.Close()

	hasCapacity, err := labs.HasCapacity(tempFile.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check if lab has capacity: %v", err)
	}

	response := &HaveCapacityResponse{
		HasCapacity: hasCapacity,
	}

	return response, status.Errorf(codes.Unimplemented, "method HaveCapacity not implemented")
}

func ScheduleLab(context.Context, *ScheduleLabRequest) (*ScheduleLabResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method ScheduleLab not implemented")
}

func GetLab(context.Context, *GetLabRequest) (*GetLabResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method GetLab not implemented")
}

func GetLabs(context.Context, *GetLabsRequest) (*GetLabsResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method GetLabs not implemented")
}