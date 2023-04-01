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
	connectedMinions    []*minion
	disconnectedMinions []*minion
	server              *grpc.Server
	port                int
)

type minion struct {
	client service.ServiceClient
	config config.MinionConfig
}

type askHasCapacityResponse struct {
	response *service.HaveCapacityResponse
	minion   *minion
	isSelf   bool
}

type askGetLabResponse struct {
	response *service.GetLabResponse
	minion   *minion
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

	connectedMinions, disconnectedMinions = ConnectMinions()
}

func Stop() {
	server.GracefulStop()
}

func (s *Server) HaveCapacity(context context.Context, request *service.HaveCapacityRequest) (*service.HaveCapacityResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// If we're the leader, then we will query all minions and return true if a single minion has capacity

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
			minion: &minion{},
			isSelf: true,
		})

		for _, connectedMinion := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()

				response, err := m.client.HaveCapacity(context, request)
				if err != nil {
					log.Error().Err(err).Msg("failed to ask minion for capacity")
					return
				}

				responseLock.Lock()
				defer responseLock.Unlock()

				responses = append(responses, &askHasCapacityResponse{
					response: response,
					minion:   m,
				})
			}(connectedMinion)
		}

		wg.Wait()

		// Find out if any of the minions have capacity (or self)
		hasAnyCapacity := false
		maxCapacity := 0
		for _, response := range responses {
			// The logic for the following two if statments is as follows:
			// If there is capacity, then the highest capacity will be enough for the lab
			// If we don't have capacity, then we will return the highest capacity of all minions

			if response.response.HasCapacity {
				hasAnyCapacity = true
			}

			if response.response.Capacity > int32(maxCapacity) {
				maxCapacity = int(response.response.Capacity)
			}
		}

		if !hasAnyCapacity {
			log.Debug().Msg("No minion has capacity")
			return &service.HaveCapacityResponse{
				HasCapacity: false,
				Capacity:    int32(maxCapacity),
			}, nil
		}

		log.Debug().Msg("At least one minion has capacity")
		return &service.HaveCapacityResponse{
			HasCapacity: true,
			Capacity:    int32(maxCapacity),
		}, nil
	}

	// If we're a minion, then we will check if we have capacity and return that
	log.Debug().Msg("Checking if we have capacity")
	return HaveCapacity(context, request)
}
func (s *Server) ScheduleLab(context context.Context, request *service.ScheduleLabRequest) (*service.ScheduleLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// Ask all minions what their capacity is and compare them including our own capacity.
		// Schedule the lab on the minion with the most capacity (this instance included)

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
				minion: &minion{},
				isSelf: true,
			})
		}

		for _, m := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()
				response, err := m.client.HaveCapacity(context, &service.HaveCapacityRequest{
					Lab: request.Lab,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get capacity from minion %s:%d", m.config.Address, m.config.Port)
				}

				if response.HasCapacity {
					responseLock.Lock()
					defer responseLock.Unlock()

					responses = append(responses, &askHasCapacityResponse{
						response: response,
						minion:   m,
					})
				}

				log.Debug().
					Int("capacity", int(response.Capacity)).
					Bool("hasCapacity", response.HasCapacity).
					Msgf("Got capacity from minion %s:%d", m.config.Address, m.config.Port)
			}(m)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Error().Msg("No minions have capacity")
			return nil, fmt.Errorf("no minions have capacity")
		}

		// Find the minion with the most capacity
		bestMinion := responses[0]
		for _, r := range responses {
			if r.response.Capacity > bestMinion.response.Capacity {
				bestMinion = r
			}
		}

		if bestMinion.isSelf {
			log.Info().Msg("Scheduling lab on self")
			return ScheduleLab(context, request)
		}

		log.Info().Msgf("Scheduling lab on minion %s:%d", bestMinion.minion.config.Address, bestMinion.minion.config.Port)
		return bestMinion.minion.client.ScheduleLab(context, request)
	}

	// In this case, this server instance is a minion and we'll just schedule the lab on ourself
	return ScheduleLab(context, request)
}
func (s *Server) GetLab(context context.Context, request *service.GetLabRequest) (*service.GetLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		wg := sync.WaitGroup{}
		responses := make([]*askGetLabResponse, 0)
		responseLock := sync.Mutex{}

		// Check if we have the lab
		lab, _ := labs.GetById(request.Id)
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
		for _, connMinion := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()
				response, err := m.client.GetLab(context, &service.GetLabRequest{
					Id: request.Id,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get lab from minion %s:%d", m.config.Address, m.config.Port)
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
						minion:   m,
					})
				}

				log.Debug().
					Msgf("Got lab from minion %s:%d", m.config.Address, m.config.Port)
			}(connMinion)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Error().Msg("No minions have the lab")
			return nil, fmt.Errorf("no minions have the lab")
		}

		// Find the minion who has the lab
		theMinion := responses[0]
		for _, r := range responses {
			if r.response.Lab != nil {
				theMinion = r
				break
			}
		}

		// TODO: expand response type with location of lab
		return theMinion.response, nil
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

		for _, connMinion := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()
				response, err := m.client.GetLabs(context, &service.GetLabsRequest{})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get labs from minion %s:%d", m.config.Address, m.config.Port)
				}

				if len(response.Labs) > 0 {
					responseLock.Lock()
					defer responseLock.Unlock()

					responses = append(responses, response)
				}

				log.Debug().
					Msgf("Got labs from minion %s:%d", m.config.Address, m.config.Port)
			}(connMinion)
		}
		wg.Wait()

		if len(responses) == 0 {
			log.Debug().Msg("No minions have labs")
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
		// Ask all minions for the given lab
		// If we find it, remove it
		// If we don't find it, return an error

		wg := sync.WaitGroup{}
		labFound := false
		theMinion := &minion{}
		responseLock := sync.Mutex{}

		// Check if the lab is hosted on this server
		lab, err := labs.GetById(request.Id)
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

		for _, connMinion := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()
				response, err := m.client.GetLab(context, &service.GetLabRequest{
					Id: request.Id,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get lab from minion %s:%d", m.config.Address, m.config.Port)
				}

				if response == nil {
					log.Error().Msg("Get lab response from follower was nil!")
					return
				}

				if response.Lab != nil {
					responseLock.Lock()
					defer responseLock.Unlock()

					labFound = true
					theMinion = m
				}

				log.Debug().
					Msgf("Got lab from minion %s:%d", m.config.Address, m.config.Port)
			}(connMinion)
		}

		wg.Wait()

		if !labFound {
			log.Error().Msg("No minions have the lab")
			return nil, fmt.Errorf("no minions have the lab")

		}

		// Remove the lab from the minion
		_, err = theMinion.client.RemoveLab(context, &service.RemoveLabRequest{
			Id: request.Id,
		})
		if err != nil {
			log.Error().Err(err).Msgf("Failed to remove lab from minion %s:%d", theMinion.config.Address, theMinion.config.Port)
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
