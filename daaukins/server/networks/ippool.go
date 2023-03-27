package networks

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
)

const (
	octetMinimum = 4 // first four IPs are reserved for common network services
	octetMaximum = 255
)

var (
	allowedLeftmostOctets = []int{172, 10}
	ipBank                *IPPool

	ErrEmptySubnetPool = fmt.Errorf("no more subnets available")
)

type IPPool struct {
	subnetsInUse map[string]bool
	ipsInUse     map[string]bool
	freeIps      map[string][]int
}

func ipPool() *IPPool {
	if ipBank == nil {
		ipBank = initIPPool()
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

	return fmt.Sprintf("%s.%d", ipPrefix, ipbank.freeIps[subnet][0]), nil
}

// GetSubnet returns a random subnet that is not in use
// this has the format "172.x.y.0/24" with the last octet missing
func (ipbank *IPPool) GetUnusedSubnet() (string, error) {
	for safety := 0; safety < 10000; safety++ {
		leftmostOctet := getRandomLeftmostOctet()
		ip := fmt.Sprintf("%d.%d.%d.0/24",
			leftmostOctet,
			getRandomOctet(leftmostOctet),
			getRandomOctet(-1)) // Yeah this is a hack, but it works..??

		if _, ok := ipbank.subnetsInUse[ip]; !ok {
			ipbank.subnetsInUse[ip] = true
			return ip, nil
		}
	}

	return "", ErrEmptySubnetPool
}

func getRandomLeftmostOctet() int {
	index := int(rand.Int31n(int32(len(allowedLeftmostOctets) - 1)))
	return allowedLeftmostOctets[index]
}

func getRandomOctet(leftmostOctet int) int {
	switch leftmostOctet {
	case 172:
		// In the case of leftmostOctet being 172, ensure that we don't generate
		// a random octet that will collide with the bridge network's second octet
		bridgeNetwork, err := virtual.DockerClient().NetworkInfo("bridge")
		if err != nil {
			log.Panicf("failed to get bridge network info: %+v", err)
		}

		bridgeSubnet := bridgeNetwork.IPAM.Config[0].Subnet
		bridgeOctets := strings.Split(bridgeSubnet, ".")
		secondOctet, err := strconv.Atoi(bridgeOctets[1])
		if err != nil {
			log.Panicf("failed to parse bridge network subnet: %+v", err)
		}

		for safety := 0; safety < 10000; safety++ {
			octet := rand.Intn(16) + 16
			if octet != secondOctet {
				return octet
			}
		}

		log.Panicf("failed to generate random octet that does not collide with bridge network")
	case 10:
		return rand.Intn(255)
	}

	return rand.Intn(octetMaximum)
}

func generateNewSubnetList() []int {
	subnets := make([]int, 0, 255)

	for i := octetMinimum; i < octetMaximum; i++ {
		subnets = append(subnets, i)
	}

	rand.Shuffle(len(subnets), func(i, j int) { subnets[i], subnets[j] = subnets[j], subnets[i] })

	return subnets
}
