package networks

import (
	"os"
	"strings"
	"testing"

	"github.com/andreaswachs/bachelors-project/daaukins/server/virtual"
)

func init() {
	if err := virtual.Initialize(); err != nil {
		panic(err)
	}
}

func TestGetUnusedSubnetNoSubnetsAreEqual(t *testing.T) {
	subnets := make(map[string]bool)

	// The function `ipPool().GetUnusedSubnet()` is quite slow, so we
	// test for a lower number of subnets to be returned
	for i := 0; i < 500; i++ {
		subnet, err := ipPool().GetUnusedSubnet()
		if err != nil {
			// we've depleted the subnet pool
			if err == ErrEmptySubnetPool {
				break
			}

			t.Fatal(err)
		}

		if subnets[subnet] {
			t.Fatal("Subnet is not unique")
		}

		subnets[subnet] = true
	}
}
func TestGetUnusedSubnetNoSubnetsAreEqualUsingDockerCompose(t *testing.T) {
	os.Setenv("DAAUKINS_USING_DOCKER_COMPOSE", "true")
	subnets := make(map[string]bool)

	// The function `ipPool().GetUnusedSubnet()` is quite slow, so we
	// test for a lower number of subnets to be returned
	for i := 0; i < 500; i++ {
		subnet, err := ipPool().GetUnusedSubnet()
		if err != nil {
			// we've depleted the subnet pool
			if err == ErrEmptySubnetPool {
				break
			}

			t.Fatal(err)
		}

		if subnets[subnet] {
			t.Fatal("Subnet is not unique")
		}

		subnets[subnet] = true
	}
}
func TestGetFreeIpWithinSubnet(t *testing.T) {
	subnet, err := ipPool().GetUnusedSubnet()
	if err != nil {
		t.Fatal(err)
	}

	ip, err := ipPool().GetFreeIP(subnet)
	if err != nil {
		t.Fatal(err)
	}

	subnetOctets := strings.Split(subnet, ".")
	ipOctets := strings.Split(ip, ".")

	for i := 0; i < 3; i++ {
		if subnetOctets[i] != ipOctets[i] {
			t.Fatalf("IP %s is not within subnet %s", ip, subnet)
		}
	}
}

func TestGetFreeIPErrorWhenInvalidSubnet(t *testing.T) {
	_, err := ipPool().GetFreeIP("123123124124:123123")
	if err == nil {
		t.Fatal("Expected error when getting ip from invalid subnet")
	}
}

func TestGetFreeIPErrorWhenSubnetIsEmpty(t *testing.T) {
	_, err := ipPool().GetFreeIP("")
	if err == nil {
		t.Fatal("Expected error when getting ip from empty subnet")
	}
}

func TestFreeIp(t *testing.T) {
	subnet, err := ipPool().GetUnusedSubnet()
	if err != nil {
		t.Fatal(err)
	}

	ip, err := ipPool().GetFreeIP(subnet)
	if err != nil {
		t.Fatal(err)
	}

	err = ipPool().FreeIP(ip)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFreeIpErrorWhenInvalidSubnet(t *testing.T) {
	err := ipPool().FreeIP("123123124124:123123")
	if err == nil {
		t.Fatal("Expected error when freeing invalid ip")
	}
}

func TestFreeIpErrorWhenSubnetIsEmpty(t *testing.T) {
	err := ipPool().FreeIP("")
	if err == nil {
		t.Fatal("Expected error when freeing empty ip")
	}
}

func TestFreeIpWhenIpContainsLetters(t *testing.T) {
	err := ipPool().FreeIP("123.123.123.123a")
	if err == nil {
		t.Fatal("Expected error when freeing ip with letters")
	}
}
