package service

import (
	context "context"
	"fmt"
	"os"
	sync "sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/labs"
	"github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

var (
	connectedMinions    []*minion
	disconnectedMinions []*minion
)

type minion struct {
	client ServiceClient
	config config.MinionConfig
}

type askHasCapacityResponse struct {
	response *HaveCapacityResponse
	minion   *minion
	isSelf   bool
}

type askGetLabResponse struct {
	response *GetLabResponse
	minion   *minion
	isSelf   bool
}

type Server struct {
	UnimplementedServiceServer
}

func (s *Server) HaveCapacity(context context.Context, request *HaveCapacityRequest) (*HaveCapacityResponse, error) {
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
			response: &HaveCapacityResponse{
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
			return &HaveCapacityResponse{
				HasCapacity: false,
				Capacity:    int32(maxCapacity),
			}, nil
		}

		log.Debug().Msg("At least one minion has capacity")
		return &HaveCapacityResponse{
			HasCapacity: true,
			Capacity:    int32(maxCapacity),
		}, nil
	}

	// If we're a minion, then we will check if we have capacity and return that
	log.Debug().Msg("Checking if we have capacity")
	return HaveCapacity(context, request)
}
func (s *Server) ScheduleLab(context context.Context, request *ScheduleLabRequest) (*ScheduleLabResponse, error) {
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
				response: &HaveCapacityResponse{
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
				response, err := m.client.HaveCapacity(context, &HaveCapacityRequest{
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
func (s *Server) GetLab(context context.Context, request *GetLabRequest) (*GetLabResponse, error) {
	if config.GetServerMode() == config.ModeLeader {
		// Ask all minions for the given lab
		// If we find it, return it
		// If we don't find it, return an error

		wg := sync.WaitGroup{}
		responses := make([]*askGetLabResponse, 0)
		responseLock := sync.Mutex{}

		for _, connMinion := range connectedMinions {
			wg.Add(1)
			go func(m *minion) {
				defer wg.Done()
				response, err := m.client.GetLab(context, &GetLabRequest{
					Id: request.Id,
				})
				if err != nil {
					log.Error().Err(err).Msgf("Failed to get lab from minion %s:%d", m.config.Address, m.config.Port)
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
func (s *Server) GetLabs(context context.Context, request *GetLabsRequest) (*GetLabsResponse, error) {
	return GetLabs(context, request)
}
func (s *Server) RemoveLab(context context.Context, request *RemoveLabRequest) (*RemoveLabResponse, error) {
	return RemoveLab(context, request)
}

// ConnectMinions attepmts to connect to all minions in the config
func ConnectMinions() {
	for _, minionConfig := range config.GetMinions() {
		minionBuffer := &minion{
			config: minionConfig,
		}

		// TODO: mTLS
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", minionConfig.Address, minionConfig.Port), grpc.WithInsecure())
		if err != nil {
			disconnectedMinions = append(disconnectedMinions, minionBuffer)
			log.Error().Err(err).Msgf("Failed to connect to minion %s:%d", minionConfig.Address, minionConfig.Port)
		}

		serviceClient := NewServiceClient(conn)
		minionBuffer.client = serviceClient

		connectedMinions = append(connectedMinions, minionBuffer)
		log.Info().Msgf("Connected to minion %s:%d", minionConfig.Address, minionConfig.Port)
	}
}
