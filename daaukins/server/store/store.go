package store

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	store *Store

	ErrorChallengeNotFound       = fmt.Errorf("challenge not found")
	ErrorChallengeNameEmpty      = fmt.Errorf("challenge name is empty")
	ErrorChallengeIdEmpty        = fmt.Errorf("challenge id is empty")
	ErrorChallengeImageEmpty     = fmt.Errorf("challenge image is empty")
	ErrorChallengeMemoryEmpty    = fmt.Errorf("challenge memory is empty")
	ErrorChallengeNegativeMemory = fmt.Errorf("challenge memory is negative")
)

type ChallengeTemplate struct {
	Name   string `yaml:"name"`
	Id     string `yaml:"id"`
	Image  string `yaml:"image"`
	Memory int    `yaml:"memory"`
}

type Store struct {
	challenges map[string]ChallengeTemplate
}

type storeDTO struct {
	Challenges []ChallengeTemplate `yaml:"challenges"`
}

func init() {
	store = &Store{
		challenges: make(map[string]ChallengeTemplate),
	}

	if err := Load("store.yaml"); err != nil {
		log.Err(err).Msg("failed to load store from default path \"store.yaml\"")
	}
}

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dto, err := loadStore(data)
	if err != nil {
		return err
	}

	if err = validateStoreDTO(dto); err != nil {
		return err
	}

	transferChallenges(dto)

	return nil
}

func GetChallenge(name string) (ChallengeTemplate, error) {
	if _, ok := store.challenges[name]; !ok {
		return ChallengeTemplate{}, fmt.Errorf("%w: %s", ErrorChallengeNotFound, name)
	}

	return store.challenges[name], nil
}

func ChallengeExists(name string) bool {
	if _, ok := store.challenges[name]; !ok {
		return false
	}

	return true
}

func loadStore(data []byte) (*storeDTO, error) {
	var dto storeDTO

	err := yaml.Unmarshal(data, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func transferChallenges(dto *storeDTO) {
	var wg sync.WaitGroup

	for _, c := range dto.Challenges {
		store.challenges[c.Id] = c

		wg.Add(1)
		go func(c ChallengeTemplate) {
			defer wg.Done()

			images, err := virtual.DockerClient().ListImages(docker.ListImagesOptions{All: true})
			if err != nil {
				log.Err(err).Msg("failed to list images on the server")
			}

			// Ensure that we don't pull images that are already present
			for _, image := range images {
				for _, tag := range image.RepoTags {
					if tag == c.Image {
						log.Info().Msgf("image already present: %s", c.Image)
						return
					}
				}
			}

			if flag.Lookup("test.v") != nil {
				return
			}

			err = virtual.DockerClient().PullImage(docker.PullImageOptions{
				Repository: c.Image,
			}, docker.AuthConfiguration{})
			if err != nil {
				log.Err(err).Msgf("failed to pull image: %s", c.Image)
			}

			log.Info().Msgf("pulled image: %s", c.Image)
		}(c)

		wg.Wait()
	}
}

func validateStoreDTO(dto *storeDTO) error {
	for _, c := range dto.Challenges {
		if c.Name == "" {
			return ErrorChallengeNameEmpty
		}

		if c.Id == "" {
			return ErrorChallengeIdEmpty
		}

		if c.Image == "" {
			return ErrorChallengeImageEmpty
		}

		if c.Memory == 0 {
			return ErrorChallengeMemoryEmpty
		}

		if c.Memory < 0 {
			return ErrorChallengeNegativeMemory
		}
	}

	return nil
}
