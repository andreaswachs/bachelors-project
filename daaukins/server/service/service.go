package service

import (
	context "context"
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	service "github.com/andreaswachs/daaukins-service"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

var (
	connectedFollowers    []*follower
	disconnectedFollowers []*follower
	server                *grpc.Server
	port                  int
)

type follower struct {
	client service.ServiceClient
	config config.FollowerConfig
}

type askHasCapacityResponse struct {
	response *service.HaveCapacityResponse
	follower *follower
	isSelf   bool
}

type askGetLabResponse struct {
	response *service.GetLabResponse
	follower *follower
	isSelf   bool
}

type Server struct {
	service.UnimplementedServiceServer
}

func Initialize() {
	port = config.GetServicePort()

	server = grpc.NewServer(
	// grpc.Creds(LoadKeyPair()),
	// grpc.UnaryInterceptor(middlefunc),
	)

	service.RegisterServiceServer(server, new(Server))
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Panic().Err(err).Msgf("failed to listen on port %d", port)
		}
		log.Info().Msgf("listening on port %d", port)
		if err := server.Serve(l); err != nil {
			log.Panic().Err(err).Msg("failed to serve")
		}
	}()

	connectedFollowers, disconnectedFollowers = ConnectFollowers()
}

func Stop() {
	server.GracefulStop()
}

func (s *Server) HaveCapacity(context context.Context, request *service.HaveCapacityRequest) (*service.HaveCapacityResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// If we're the leader, then we will query all follower and return true if a single follower has capacity

		wg := sync.WaitGroup{}
		responses := make([]*askHasCapacityResponse, 0)
		responseLock := sync.Mutex{}

		file, err := saveLabFile(request.Lab)
		if err != nil {
			return nil, fmt.Errorf("failed to save lab file: %v", err)
		}

		hasCapacity, err := labs.HasCapacity(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to check if lab has capacity: %v", err)
		}

		capacity, err := labs.GetCapacity()
		if err != nil {
			return nil, fmt.Errorf("failed to get capacity: %v", err)
		}

		responses = append(responses, &askHasCapacityResponse{
			response: &service.HaveCapacityResponse{
				HasCapacity: hasCapacity,
				Capacity:    int32(capacity),
			},
			follower: &follower{},
			isSelf:   true,
		})

		for _, connectedFollower := range connectedFollowers {
			wg.Add(1)
			go func(m *follower) {
				defer wg.Done()

				response, err := m.client.HaveCapacity(context, request)
				if err != nil {
					log.Error().Err(err).Msg("failed to ask follower for capacity")
					return
				}

				responseLock.Lock()
				defer responseLock.Unlock()

				responses = append(responses, &askHasCapacityResponse{
					response: response,
					follower: m,
				})
			}(connectedFollower)
		}

		wg.Wait()

		// Find out if any of the followers have capacity (or self)
		hasAnyCapacity := false
		maxCapacity := 0
		for _, response := range responses {
			// The logic for the following two if statments is as follows:
			// If there is capacity, then the highest capacity will be enough for the lab
			// If we don't have capacity, then we will return the highest capacity of all followers

			if response.response.HasCapacity {
				hasAnyCapacity = true
			}

			if response.response.Capacity > int32(maxCapacity) {
				maxCapacity = int(response.response.Capacity)
			}
		}

		if !hasAnyCapacity {
			log.Debug().Msg("No followers has capacity")
			return &service.HaveCapacityResponse{
				HasCapacity: false,
				Capacity:    int32(maxCapacity),
			}, nil
		}

		log.Debug().Msg("At least one followers has capacity")
		return &service.HaveCapacityResponse{
			HasCapacity: true,
			Capacity:    int32(maxCapacity),
		}, nil
	}

	// If we're a follower, then we will check if we have capacity and return that
	log.Debug().Msg("Checking if we have capacity")
	return HaveCapacity(context, request)
}
func (s *Server) ScheduleLab(context context.Context, request *service.ScheduleLabRequest) (*service.ScheduleLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// Ask all follower what their capacity is and compare them including our own capacity.
		// Schedule the lab on the follower with the most capacity (this instance included)

		wg := sync.WaitGroup{}
		responses := make([]*askHasCapacityResponse, 0)
		responseLock := sync.Mutex{}

		file, err := saveLabFile(request.Lab)
		if err != nil {
			return nil, fmt.Errorf("failed to save lab file: %v", err)
		}
		defer os.Remove(file.Name())

		hasCapacity, err := labs.HasCapacity(file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to check if lab has capacity: %v", err)
		}

		if hasCapacity {
			capacity, err := labs.GetCapacity()
			if err != nil {
				return nil, fmt.Errorf("failed to get capacity: %v", err)
			}

			responses = append(responses, &askHasCapacityResponse{
				response: &service.HaveCapacityResponse{
					HasCapacity: true,
					Capacity:    int32(capacity),
				},
				follower: &follower{},
				isSelf:   true,
			})
		}

		for _, m := range connectedFollowers {
			wg.Add(1)
			go func(m *follower) {
				defer wg.Done()
				response, err := m.client.HaveCapacity(context, &service.HaveCapacityRequest{
					Lab: request.Lab,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get capacity from follower %s:%d", m.config.Address, m.config.Port)
				}

				if response.HasCapacity {
					responseLock.Lock()
					defer responseLock.Unlock()

					responses = append(responses, &askHasCapacityResponse{
						response: response,
						follower: m,
					})
				}

				log.Debug().
					Int("capacity", int(response.Capacity)).
					Bool("hasCapacity", response.HasCapacity).
					Msgf("Got capacity from follower %s:%d", m.config.Address, m.config.Port)
			}(m)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Error().Msg("No follower have capacity")
			return nil, fmt.Errorf("no follower have capacity")
		}

		// Find the follower with the most capacity
		bestFollower := responses[0]
		for _, r := range responses {
			if r.response.Capacity > bestFollower.response.Capacity {
				bestFollower = r
			}
		}

		if bestFollower.isSelf {
			log.Info().Msg("Scheduling lab on self")
			return ScheduleLab(context, request)
		}

		log.Info().Msgf("Scheduling lab on follower %s:%d", bestFollower.follower.config.Address, bestFollower.follower.config.Port)
		return bestFollower.follower.client.ScheduleLab(context, request)
	}

	// In this case, this server instance is a folloer and we'll just schedule the lab on ourself
	return ScheduleLab(context, request)
}
func (s *Server) GetLab(context context.Context, request *service.GetLabRequest) (*service.GetLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		wg := sync.WaitGroup{}
		responses := make([]*askGetLabResponse, 0)
		responseLock := sync.Mutex{}

		// Check if we have the lab
		lab, _ := labs.WithId(request.Id)
		if lab != nil {
			return &service.GetLabResponse{
				Lab: &service.LabDescription{
					Id:            lab.GetId(),
					Name:          lab.GetName(),
					NumChallenges: int32(len(lab.GetChallenges())),
					NumUsers:      1,
					ServerId:      config.GetServerID(),
				},
			}, nil
		}

		// Ask all the followers if the lab is located on one of them
		for _, connFollower := range connectedFollowers {
			wg.Add(1)
			go func(m *follower) {
				defer wg.Done()
				response, err := m.client.GetLab(context, &service.GetLabRequest{
					Id: request.Id,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get lab from follower %s:%d", m.config.Address, m.config.Port)
				}

				if response == nil {
					log.Error().Msg("Get lab response from follower was nil!")
					return
				}

				if response.Lab != nil {
					responseLock.Lock()
					defer responseLock.Unlock()

					responses = append(responses, &askGetLabResponse{
						response: response,
						follower: m,
					})
				}

				log.Debug().
					Msgf("Got lab from folloer %s:%d", m.config.Address, m.config.Port)
			}(connFollower)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Error().Msg("No follower have the lab")
			return nil, fmt.Errorf("no folloer have the lab")
		}

		// Find the follower who has the lab
		theFollower := responses[0]
		for _, r := range responses {
			if r.response.Lab != nil {
				theFollower = r
				break
			}
		}

		// TODO: expand response type with location of lab
		return theFollower.response, nil
	}

	return GetLab(context, request)
}
func (s *Server) GetLabs(context context.Context, request *service.GetLabsRequest) (*service.GetLabsResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		wg := sync.WaitGroup{}
		responses := make([]*service.GetLabsResponse, 0)
		responseLock := sync.Mutex{}

		// Get labs from ourselves
		localLabs, err := GetLabs(context, request)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get labs from self")
		}

		log.Debug().Int("labs", len(localLabs.Labs)).Msg("Got labs from self")

		if localLabs.Labs != nil && len(localLabs.Labs) > 0 {
			responses = append(responses, localLabs)
		}

		for _, connFollower := range connectedFollowers {
			wg.Add(1)
			go func(m *follower) {
				defer wg.Done()
				response, err := m.client.GetLabs(context, &service.GetLabsRequest{})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get labs from follower %s:%d", m.config.Address, m.config.Port)
				}

				if len(response.Labs) > 0 {
					responseLock.Lock()
					defer responseLock.Unlock()

					responses = append(responses, response)
				}

				log.Debug().
					Msgf("Got labs from follower %s:%d", m.config.Address, m.config.Port)
			}(connFollower)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Debug().Msg("No follower have labs")
			return &service.GetLabsResponse{
				Labs: make([]*service.LabDescription, 0),
			}, nil
		}

		response := &service.GetLabsResponse{
			Labs: make([]*service.LabDescription, 0),
		}

		for _, r := range responses {
			response.Labs = append(response.Labs, r.Labs...)
		}

		return response, nil
	}

	return GetLabs(context, request)
}
func (s *Server) RemoveLab(context context.Context, request *service.RemoveLabRequest) (*service.RemoveLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// Ask all follower for the given lab
		// If we find it, remove it
		// If we don't find it, return an error

		wg := sync.WaitGroup{}
		labFound := false
		theFollower := &follower{}
		responseLock := sync.Mutex{}

		// Check if the lab is hosted on this server
		lab, err := labs.WithId(request.Id)
		if err == nil {
			if lab != nil {
				err = lab.Remove()
				if err != nil {
					log.Error().Err(err).Msg("Failed to remove lab")
					return nil, err
				}

				return &service.RemoveLabResponse{}, nil
			}
		}

		for _, connFollower := range connectedFollowers {
			wg.Add(1)
			go func(m *follower) {
				defer wg.Done()
				response, err := m.client.GetLab(context, &service.GetLabRequest{
					Id: request.Id,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get lab from follower %s:%d", m.config.Address, m.config.Port)
				}

				if response == nil {
					log.Error().Msg("Get lab response from follower was nil!")
					return
				}

				if response.Lab != nil {
					responseLock.Lock()
					defer responseLock.Unlock()

					labFound = true
					theFollower = m
				}

				log.Debug().
					Msgf("Got lab from follower %s:%d", m.config.Address, m.config.Port)
			}(connFollower)
		}

		wg.Wait()

		if !labFound {
			log.Error().Msg("No follower have the lab")
			return nil, fmt.Errorf("no followers have the lab")

		}

		// Remove the lab from the follower
		_, err = theFollower.client.RemoveLab(context, &service.RemoveLabRequest{
			Id: request.Id,
		})
		if err != nil {
			log.Error().Err(err).Msgf("Failed to remove lab from follower %s:%d", theFollower.config.Address, theFollower.config.Port)
			return nil, err
		}

		return &service.RemoveLabResponse{}, nil
	}

	return RemoveLab(context, request)
}

func (s *Server) GetServerMode(context context.Context, request *service.GetServerModeRequest) (*service.GetServerModeResponse, error) {
	return &service.GetServerModeResponse{
		Mode:     config.GetServerMode().String(),
		ServerId: config.GetServerID(),
	}, nil
}
