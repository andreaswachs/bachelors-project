package service

import (
	context "context"
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/frontend"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	service "github.com/andreaswachs/daaukins-service"
	"github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func HaveCapacity(context context.Context, request *service.HaveCapacityRequest) (*service.HaveCapacityResponse, error) {
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

	capacity, err := labs.GetCapacity()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get capacity: %v", err)
	}

	response := &service.HaveCapacityResponse{
		HasCapacity: hasCapacity,
		Capacity:    int32(capacity),
	}

	return response, nil
}

func ScheduleLab(context context.Context, request *service.ScheduleLabRequest) (*service.ScheduleLabResponse, error) {
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

	return &service.ScheduleLabResponse{Id: lab.GetId()}, nil

}

func GetLab(context context.Context, request *service.GetLabRequest) (*service.GetLabResponse, error) {
	if request.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	lab, err := labs.WithId(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return nil, status.Errorf(codes.NotFound, "lab not found")
	}

	response := &service.GetLabResponse{
		Lab: &service.LabDescription{
			Name:          lab.GetName(),
			Id:            lab.GetId(),
			NumChallenges: int32(len(lab.GetChallenges())),
			NumUsers:      1, // PoC limitation: only deploy one frontend to each lab
			ServerId:      config.GetServerID(),
		},
	}

	return response, nil
}

func GetLabs(context context.Context, _request *service.GetLabsRequest) (*service.GetLabsResponse, error) {
	labs := labs.All()

	log.Debug().Int("numLabs", len(labs)).Msg("GetLabs")

	labsResponse := make([]*service.LabDescription, len(labs))
	for i, lab := range labs {
		labsResponse[i] = &service.LabDescription{
			Name:          lab.GetName(),
			Id:            lab.GetId(),
			NumChallenges: int32(len(lab.GetChallenges())),
			NumUsers:      1, // PoC limitation: only deploy one frontend to each lab
			ServerId:      config.GetServerID(),
		}
	}

	response := &service.GetLabsResponse{
		Labs: labsResponse,
	}

	return response, nil
}

func RemoveLab(context context.Context, request *service.RemoveLabRequest) (*service.RemoveLabResponse, error) {
	if request.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	lab, err := labs.WithId(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return nil, status.Errorf(codes.NotFound, "lab not found")
	}

	if err := lab.Remove(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to stop lab: %v", err)
	}

	return &service.RemoveLabResponse{
		Ok: true,
	}, nil
}

func RemoveAllLabs(ctx context.Context, request *service.RemoveLabsRequest) (*service.RemoveLabsResponse, error) {
	// If the server ID is not empty (empty meaning that we want to remove labs from all servers) and it is not equal to this server's ID, return an error
	if request.ServerId != "" && request.ServerId != config.GetServerID() {
		log.Error().
			Str("requested server id", request.ServerId).
			Str("serverID", config.GetServerID()).
			Msg("RemoveAllLabs: ServerID was neither empty nor equal to this server's ID")

		return &service.RemoveLabsResponse{}, status.Errorf(codes.InvalidArgument, "ServerID was neither empty nor equal to this server's ID")
	}

	err := labs.RemoveAll()
	if err != nil {
		return &service.RemoveLabsResponse{}, status.Errorf(codes.Internal, "failed to remove all labs: %v", err)
	}

	return &service.RemoveLabsResponse{
		Ok: true,
	}, nil
}

func GetFrontends(ctx context.Context, request *service.GetFrontendsRequest) (*service.GetFrontendsResponse, error) {
	frontends := make([]frontend.T, 0)
	for _, lab := range labs.All() {
		frontends = append(frontends, lab.GetFrontend())
	}

	log.Debug().Int("numFrontends", len(frontends)).Msg("GetFrontends")

	frontendsResponse := make([]*service.Frontend, len(frontends))
	for i, frontend := range frontends {
		frontendsResponse[i] = &service.Frontend{
			Port:     fmt.Sprint(frontend.GetProxyPort()),
			Host:     "replace-me",
			ServerId: config.GetServerID(),
		}
	}

	return &service.GetFrontendsResponse{
		Frontends: frontendsResponse,
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
