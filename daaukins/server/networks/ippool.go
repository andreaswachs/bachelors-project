package networks

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

		ip := fmt.Sprintf("%d.%d.%d.0/24",
			getRandomLeftmostOctet(),
			getRandomOctet(),
			getRandomOctet())

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

func getRandomOctet() int {
	return int(rand.Int31n(octetMaximum))
}

func generateNewSubnetList() []int {
	subnets := make([]int, 0, 255)

	for i := octetMinimum; i < octetMaximum; i++ {
		subnets = append(subnets, i)
	}

	rand.Shuffle(len(subnets), func(i, j int) { subnets[i], subnets[j] = subnets[j], subnets[i] })

	return subnets
}
