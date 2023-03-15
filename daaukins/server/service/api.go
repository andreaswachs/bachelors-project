// api.go contains an implementation of the gRPC service defined in service/service.proto.

package service

import (
	context "context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	"github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func HaveCapacity(context context.Context, request *HaveCapacityRequest) (*HaveCapacityResponse, error) {
	if request.Lab == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lab is empty")
	}

	tempFile, err := saveLabFile(request.Lab)
	if err != nil {
		return nil, err
	}

	hasCapacity, err := labs.HasCapacity(tempFile.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check if lab has capacity: %v", err)
	}

	response := &HaveCapacityResponse{
		HasCapacity: hasCapacity,
	}

	return response, nil
}

func ScheduleLab(context context.Context, request *ScheduleLabRequest) (*ScheduleLabResponse, error) {
	if request.Lab == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lab is empty")
	}

	tempFile, err := saveLabFile(request.Lab)
	if err != nil {
		return nil, err
	}

	lab, err := labs.Provision(tempFile.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to provision lab: %v", err)
	}

	if err = lab.Start(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start lab: %v", err)
	}

	// TODO: is returning the name as ID right?
	return &ScheduleLabResponse{Id: lab.GetName()}, nil

}

func GetLab(context context.Context, request *GetLabRequest) (*GetLabResponse, error) {
	if request.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	// TODO: is using name as ID right??
	lab, err := labs.GetByName(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return nil, status.Errorf(codes.NotFound, "lab not found")
	}

	response := &GetLabResponse{
		Lab: &LabDescription{
			Name:          lab.GetName(),
			Id:            lab.GetName(),
			NumChallenges: int32(len(lab.GetChallenges())),
			NumUsers:      0, // TODO
		},
	}

	return response, nil
}

func GetLabs(context context.Context, _request *GetLabsRequest) (*GetLabsResponse, error) {
	labs := labs.GetAll()
	m, err := json.MarshalIndent(labs, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debug().Int("numLabs", len(labs)).Str("labs", string(m)).Msg("GetLabs")
	labsResponse := make([]*LabDescription, len(labs))
	for i, lab := range labs {
		labsResponse[i] = &LabDescription{
			Name:          lab.GetName(),
			Id:            lab.GetName(),
			NumChallenges: int32(len(lab.GetChallenges())),
			NumUsers:      0, // TODO
		}
	}

	response := &GetLabsResponse{
		Labs: labsResponse,
	}

	return response, nil
}

func RemoveLab(context context.Context, request *RemoveLabRequest) (*RemoveLabResponse, error) {
	if request.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	lab, err := labs.GetByName(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return nil, status.Errorf(codes.NotFound, "lab not found")
	}

	if err := lab.Remove(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stop lab: %v", err)
	}

	return &RemoveLabResponse{
		Ok: true,
	}, nil
}

func saveLabFile(lab string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "daaukins-lab.yaml")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create temp file: %v", err)
	}

	if _, err := tempFile.Write([]byte(lab)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write lab to temp file: %v", err)
	}

	return tempFile, nil
}
