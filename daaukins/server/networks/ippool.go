package networks

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/andreaswachs/bachelors-project/daaukins/server/config"
	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
	"github.com/rs/zerolog/log"
)

const (
	octetMinimum = 4 // first four IPs are reserved for common network services
	octetMaximum = 255
)

var (
	allowedLeftmostOctets = []int{172, 10}
	ipBank                *IPPool
	ipPoolArbiter         = &sync.Mutex{}

	ErrEmptySubnetPool = fmt.Errorf("no more subnets available")
)

type IPPool struct {
	subnetsInUse map[string]bool
	ipsInUse     map[string]bool
	freeIps      map[string][]int
}

func ipPool() *IPPool {
	if ipBank == nil {
		ipPoolArbiter.Lock()
		defer ipPoolArbiter.Unlock()
		ipBank = initIPPool()

		if config.IsUsingDockerCompose() {
			// Remove the subnet that is used by the docker-compose network
			// Get the subnet from the docker-compose network
			network, err := virtual.DockerClient().NetworkInfo("server_default")
			if err != nil {
				log.Error().Err(err).Msg("failed to get docker-compose network \"server_default\"")
				return ipBank
			}

			subnet := network.IPAM.Config[0].Subnet
			octets := strings.Split(subnet, ".")
			for i := 0; i <= 255; i++ {
				subnetToIgnore := fmt.Sprintf("%s.%s.%d.0/24", octets[0], octets[1], i)
				ipBank.subnetsInUse[subnetToIgnore] = true
			}

		}
	}

	return ipBank
}

func initIPPool() *IPPool {
	return &IPPool{
		subnetsInUse: make(map[string]bool),
		freeIps:      make(map[string][]int),
	}
}

func (ipbank *IPPool) FreeIP(ip string) error {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return fmt.Errorf("invalid ip format")
	}

	// Remove the ip from the in use list
	delete(ipbank.ipsInUse, ip)

	// Reconstruct the subnet
	subnet := fmt.Sprintf("%s.%s.%s.0/24", octets[0], octets[1], octets[2])

	// add the ip back to the free list
	lastOctet, err := strconv.Atoi(octets[3])
	if err != nil {
		return err
	}

	ipbank.freeIps[subnet] = append(ipbank.freeIps[subnet], lastOctet)
	return nil
}

func (ipbank *IPPool) GetFreeIP(subnet string) (string, error) {
	if _, ok := ipbank.freeIps[subnet]; !ok {
		ipbank.freeIps[subnet] = generateNewSubnetList()
	}

	if len(ipbank.freeIps[subnet]) == 0 {
		return "", fmt.Errorf("no more ips available")
	}

	octets := strings.Split(subnet, ".")
	if len(octets) != 4 {
		return "", fmt.Errorf("invalid subnet format")
	}

	ipPrefix := fmt.Sprintf("%s.%s.%s", octets[0], octets[1], octets[2])

	usedOctet := ipbank.freeIps[subnet][0]
	ipbank.freeIps[subnet] = ipbank.freeIps[subnet][1:]

	return fmt.Sprintf("%s.%d", ipPrefix, usedOctet), nil
}

// GetSubnet returns a random subnet that is not in use
// this has the format "172.x.y.0/24" with the last octet missing
func (ipbank *IPPool) GetUnusedSubnet() (string, error) {
	for safety := 0; safety < 10000; safety++ {
		leftmostOctet := getRandomLeftmostOctet()

		octet1, err := getRandomOctet(leftmostOctet)
		if err != nil {
			continue
		}

		octet2, err := getRandomOctet(-1)
		if err != nil {
			continue
		}

		ip := fmt.Sprintf("%d.%d.%d.0/24", leftmostOctet, octet1, octet2)

		ipPoolArbiter.Lock()
		defer ipPoolArbiter.Unlock()

		if _, ok := ipbank.subnetsInUse[ip]; !ok {
			ipbank.subnetsInUse[ip] = true
			return ip, nil
		}
	}

	return "", ErrEmptySubnetPool
}

func getRandomLeftmostOctet() int {
	index := int(rand.Int31n(int32(len(allowedLeftmostOctets))))
	return allowedLeftmostOctets[index]
}

func getRandomOctet(leftmostOctet int) (int, error) {
	switch leftmostOctet {
	case 172:
		// In the case of leftmostOctet being 172, ensure that we don't generate
		// a random octet that will collide with the bridge network's second octet
		bridgeNetwork, err := virtual.DockerClient().NetworkInfo("bridge")
		if err != nil {
			return 0, err
		}

		bridgeSubnet := bridgeNetwork.IPAM.Config[0].Subnet
		bridgeOctets := strings.Split(bridgeSubnet, ".")
		secondOctet, err := strconv.Atoi(bridgeOctets[1])
		if err != nil {
			return 0, err
		}

		for safety := 0; safety < 10000; safety++ {
			octet := rand.Intn(16) + 16
			if octet != secondOctet {
				return octet, nil
			}
		}

		return 0, fmt.Errorf("failed to generate random octet that does not collide with bridge network")
	case 10:
		return rand.Intn(255), nil
	}

	return rand.Intn(octetMaximum), nil
}

func generateNewSubnetList() []int {
	subnets := make([]int, 0, 255)

	for i := octetMinimum; i < octetMaximum; i++ {
		subnets = append(subnets, i)
	}

	rand.Shuffle(len(subnets), func(i, j int) { subnets[i], subnets[j] = subnets[j], subnets[i] })

	return subnets
}
