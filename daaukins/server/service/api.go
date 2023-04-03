// api.go implements the Daaukins service gRPC API for the local server instance.
package service

import (
	context "context"
	"fmt"
	"os"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	service "github.com/andreaswachs/daaukins-service"
	"github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var (
	emptyGetFrontendsResponse = &service.GetFrontendsResponse{}
	emptyHaveCapacityResponse = &service.HaveCapacityResponse{}
	emptyScheduleLabResponse  = &service.ScheduleLabResponse{}
	emptyGetLabResponse       = &service.GetLabResponse{}
	emptyGetLabsResponse      = &service.GetLabsResponse{}
	emptyRemoveLabResponse    = &service.RemoveLabResponse{}
)

func HaveCapacity(context context.Context, request *service.HaveCapacityRequest) (*service.HaveCapacityResponse, error) {
	if request.Lab == "" {
		return emptyHaveCapacityResponse, status.Errorf(codes.InvalidArgument, "lab is empty")
	}

	tempFile, err := saveLabFile(request.Lab)
	if err != nil {
		return emptyHaveCapacityResponse, err
	}

	hasCapacity, err := labs.HasCapacity(tempFile.Name())
	if err != nil {
		return emptyHaveCapacityResponse, status.Errorf(codes.Internal, "failed to check if lab has capacity: %v", err)
	}

	capacity, err := labs.GetCapacity()
	if err != nil {
		return emptyHaveCapacityResponse, status.Errorf(codes.Internal, "failed to get capacity: %v", err)
	}

	response := &service.HaveCapacityResponse{
		HasCapacity: hasCapacity,
		Capacity:    int32(capacity),
	}

	return response, nil
}

func ScheduleLab(context context.Context, request *service.ScheduleLabRequest) (*service.ScheduleLabResponse, error) {
	if request.Lab == "" {
		return emptyScheduleLabResponse, status.Errorf(codes.InvalidArgument, "lab is empty")
	}

	tempFile, err := saveLabFile(request.Lab)
	if err != nil {
		return emptyScheduleLabResponse, err
	}

	lab, err := labs.Provision(tempFile.Name())
	if err != nil {
		return emptyScheduleLabResponse, status.Errorf(codes.Internal, "failed to provision lab: %v", err)
	}

	if err = lab.Start(); err != nil {
		return emptyScheduleLabResponse, status.Errorf(codes.Internal, "failed to start lab: %v", err)
	}

	return &service.ScheduleLabResponse{Id: lab.GetId()}, nil

}

func GetLab(context context.Context, request *service.GetLabRequest) (*service.GetLabResponse, error) {
	if request.Id == "" {
		return emptyGetLabResponse, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	lab, err := labs.GetById(request.Id)
	if err != nil {
		return emptyGetLabResponse, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return emptyGetLabResponse, status.Errorf(codes.NotFound, "lab not found")
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
	labs := labs.GetAll()

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
		return emptyRemoveLabResponse, status.Errorf(codes.InvalidArgument, "id is empty")
	}

	lab, err := labs.GetById(request.Id)
	if err != nil {
		return emptyRemoveLabResponse, status.Errorf(codes.Internal, "failed to get lab: %v", err)
	}

	if lab == nil {
		return emptyRemoveLabResponse, status.Errorf(codes.NotFound, "lab not found")
	}

	if err := lab.Remove(); err != nil {
		return emptyRemoveLabResponse, status.Errorf(codes.Internal, "failed to stop lab: %v", err)
	}

	return &service.RemoveLabResponse{
		Ok: true,
	}, nil
}

func GetFrontends(ctx context.Context, request *service.GetFrontendsRequest) (*service.GetFrontendsResponse, error) {
	// If the IDs dont match or if the request's server ID is not empty, then return an error
	// The empty server ID means that we want frontends from all servers
	if config.GetServerID() != request.GetServerId() || request.GetServerId() != "" {
		return emptyGetFrontendsResponse, status.Errorf(codes.InvalidArgument, "server id does not match")
	}

	frontends := make([]*service.Frontend, 0)
	allLabs := labs.GetAll()
	for _, lab := range allLabs {
		fe := lab.GetFrontend()

		if fe == nil {
			continue
		}

		frontends = append(frontends, &service.Frontend{
			Host:     "todo",
			Port:     fmt.Sprint(fe.GetProxyPort()),
			ServerId: config.GetServerID(),
		})

	}

	return &service.GetFrontendsResponse{
		Frontends: frontends,
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
