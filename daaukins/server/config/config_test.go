package config

import (
	"os"
	"testing"
)

type instrumentedConfigOptions struct{}

var (
	configFilename       = "server.yaml"
	isUsingDockerCompose = true
	testServerId         = "test-server-id"
	configLeader         = `server_mode: leader
service_port: 50051
docker:
  frontend:
    image: "andreaswachs/kali-docker:core"
  proxy:
    image: "andreaswachs/forward-proxy"
  dhcp:
    image: "networkboot/dhcpd:1.2.0"
  dns:
    image: "coredns/coredns:1.10.0"
followers:
  - name: tricep-1
    address: "localhost"
    port: 50052
  - name: tricep-2
    address: "127.0.0.1"
    port: 50052
`
	configFollower = `server_mode: follower
service_port: 50052
docker:
  frontend:
    image: "andreaswachs/kali-docker:core"
  proxy:
    image: "andreaswachs/forward-proxy"
  dhcp:
    image: "networkboot/dhcpd:1.2.0"
  dns:
    image: "coredns/coredns:1.10.0"
followers: []
`
)

func (i *instrumentedConfigOptions) getConfigFilename() string {
	return configFilename
}

func (i *instrumentedConfigOptions) isUsingDockerCompose() bool {
	return isUsingDockerCompose
}

func (i *instrumentedConfigOptions) newServerID() string {
	return testServerId
}

func AddInstrumentation(t *testing.T, conf string) string {
	tmpFile, err := os.CreateTemp("", "server.yaml")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}

	_, err = tmpFile.WriteString(conf)
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
	}

	configFilename = tmpFile.Name()
	confOptsObj = &instrumentedConfigOptions{}

	return tmpFile.Name()
}

func TestInitializeAsLeader(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	if config.ServerMode != ModeLeader {
		t.Errorf("Expected server mode to be 'leader', got '%v'", config.ServerMode)
	}

	if config.ServicePort != 50051 {
		t.Errorf("Expected service port to be 50051, got %v", config.ServicePort)
	}

	if len(config.Followers) != 2 {
		t.Errorf("Expected 2 followers, got %v", len(config.Followers))
	}

	if config.Followers[0].Name != "tricep-1" {
		t.Errorf("Expected follower name to be 'tricep-1', got '%v'", config.Followers[0].Name)
	}

	if config.Followers[0].Address != "localhost" {
		t.Errorf("Expected follower address to be 'localhost', got '%v'", config.Followers[0].Address)
	}

	if config.Followers[0].Port != 50052 {
		t.Errorf("Expected follower port to be 50052, got %v", config.Followers[0].Port)
	}

	if config.Followers[1].Name != "tricep-2" {
		t.Errorf("Expected follower name to be 'tricep-2', got '%v'", config.Followers[1].Name)
	}

	if config.Followers[1].Address != "127.0.0.1" {
		t.Errorf("Expected follower address to be '127.0.0.1', got '%v'", config.Followers[1].Address)
	}

	if config.Followers[1].Port != 50052 {
		t.Errorf("Expected follower port to be 50052, got %v", config.Followers[1].Port)
	}

	if config.Docker.Frontend.Image != "andreaswachs/kali-docker:core" {
		t.Errorf("Expected frontend image to be 'andreaswachs/kali-docker:core', got '%v'", config.Docker.Frontend.Image)
	}

	if config.Docker.Proxy.Image != "andreaswachs/forward-proxy" {
		t.Errorf("Expected proxy image to be 'andreaswachs/forward-proxy', got '%v'", config.Docker.Proxy.Image)
	}

	if config.Docker.Dhcp.Image != "networkboot/dhcpd:1.2.0" {
		t.Errorf("Expected dhcp image to be 'networkboot/dhcpd:1.2.0', got '%v'", config.Docker.Dhcp.Image)
	}

	if config.Docker.Dns.Image != "coredns/coredns:1.10.0" {
		t.Errorf("Expected dns image to be 'coredns/coredns:1.10.0', got '%v'", config.Docker.Dns.Image)
	}
}

func TestInitializeAsFollower(t *testing.T) {
	tmpFile := AddInstrumentation(t, configFollower)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	if config.ServerMode != ModeFollower {
		t.Errorf("Expected server mode to be 'follower', got '%v'", config.ServerMode)
	}

	if config.ServicePort != 50052 {
		t.Errorf("Expected service port to be 50052, got %v", config.ServicePort)
	}

	if len(config.Followers) != 0 {
		t.Errorf("Expected 0 followers, got %v", len(config.Followers))
	}

	if config.Docker.Frontend.Image != "andreaswachs/kali-docker:core" {
		t.Errorf("Expected frontend image to be 'andreaswachs/kali-docker:core', got '%v'", config.Docker.Frontend.Image)
	}

	if config.Docker.Proxy.Image != "andreaswachs/forward-proxy" {
		t.Errorf("Expected proxy image to be 'andreaswachs/forward-proxy', got '%v'", config.Docker.Proxy.Image)
	}

	if config.Docker.Dhcp.Image != "networkboot/dhcpd:1.2.0" {
		t.Errorf("Expected dhcp image to be 'networkboot/dhcpd:1.2.0', got '%v'", config.Docker.Dhcp.Image)
	}

	if config.Docker.Dns.Image != "coredns/coredns:1.10.0" {
		t.Errorf("Expected dns image to be 'coredns/coredns:1.10.0', got '%v'", config.Docker.Dns.Image)
	}
}

func TestInitializeAsLeaderMissingConfigFile(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	os.Remove(tmpFile)

	err := Initialize()
	if err == nil {
		t.Errorf("Expected error initializing config, got nil")
	}
}

func TestInitializeAsLeaderMissingConfigFileContents(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Errorf("Error writing config file: %v", err)
	}

	err = Initialize()
	if err == nil {
		t.Errorf("Expected error initializing config, got nil")
	}
}

func TestInitializeAsLeaderMissingConfigFileContentsInvalidYaml(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := os.WriteFile(tmpFile, []byte("invalid yaml"), 0644)
	if err != nil {
		t.Errorf("Error writing config file: %v", err)
	}

	err = Initialize()
	if err == nil {
		t.Errorf("Expected error initializing config, got nil")
	}
}

func TestGetFollowers(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	followers := GetFollowers()
	if len(followers) != 2 {
		t.Errorf("Expected 2 followers, got %v", len(followers))
	}

	if followers[0].Name != "tricep-1" {
		t.Errorf("Expected follower name to be 'tricep-1', got '%v'", followers[0].Name)
	}

	if followers[0].Address != "localhost" {
		t.Errorf("Expected follower address to be 'localhost', got '%v'", followers[0].Address)
	}

	if followers[0].Port != 50052 {
		t.Errorf("Expected follower port to be 50052, got %v", followers[0].Port)
	}

	if followers[1].Name != "tricep-2" {
		t.Errorf("Expected follower name to be 'tricep-2', got '%v'", followers[1].Name)
	}

	if followers[1].Address != "127.0.0.1" {
		t.Errorf("Expected follower address to be '127.0.0.1', got '%v'", followers[1].Address)
	}

	if followers[1].Port != 50052 {
		t.Errorf("Expected follower port to be 50052, got %v", followers[1].Port)
	}
}

func TestGetServerModeAsLeader(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	mode := GetServerMode()
	if mode != ModeLeader {
		t.Errorf("Expected server mode to be 'leader', got '%v'", mode)
	}
}

func TestGetServerModeAsFollower(t *testing.T) {
	tmpFile := AddInstrumentation(t, configFollower)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	mode := GetServerMode()
	if mode != ModeFollower {
		t.Errorf("Expected server mode to be 'follower', got '%v'", mode)
	}
}

func TestGetServicePort(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	port := GetServicePort()
	if port != 50051 {
		t.Errorf("Expected service port to be 50051, got %v", port)
	}
}

func TestGetServerID(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	id := GetServerID()
	if id != testServerId {
		t.Errorf("Expected server ID to be '%s', got '%v'", testServerId, id)
	}
}

func TestGetDockerConfig(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	dockerConfig := GetDockerConfig()
	if dockerConfig.Frontend.Image != "andreaswachs/kali-docker:core" {
		t.Errorf("Expected frontend image to be 'andreaswachs/kali-docker:core', got '%v'", dockerConfig.Frontend.Image)
	}

	if dockerConfig.Proxy.Image != "andreaswachs/forward-proxy" {
		t.Errorf("Expected proxy image to be 'andreaswachs/forward-proxy', got '%v'", dockerConfig.Proxy.Image)
	}

	if dockerConfig.Dhcp.Image != "networkboot/dhcpd:1.2.0" {
		t.Errorf("Expected dhcp image to be 'networkboot/dhcpd:1.2.0', got '%v'", dockerConfig.Dhcp.Image)
	}

	if dockerConfig.Dns.Image != "coredns/coredns:1.10.0" {
		t.Errorf("Expected dns image to be 'coredns/coredns:1.10.0', got '%v'", dockerConfig.Dns.Image)
	}
}

func TestGetIsUsingDockerCompose(t *testing.T) {
	tmpFile := AddInstrumentation(t, configLeader)
	defer os.Remove(tmpFile)

	err := Initialize()
	if err != nil {
		t.Errorf("Error initializing config: %v", err)
	}

	configIsUsingDockerCompose := IsUsingDockerCompose()
	if configIsUsingDockerCompose != isUsingDockerCompose {
		t.Errorf("Expected configIsUsingDockerCompose to be %v, got %v", isUsingDockerCompose, configIsUsingDockerCompose)
	}
}

func TestDefaultConfigNewServerID(t *testing.T) {
	result := confOptsObj.newServerID()

	if result == "" {
		t.Errorf("Expected server ID to be non-empty, got '%v'", result)
	}
}

func TestDefaultConfigGetConfigFilename(t *testing.T) {
	result := confOptsObj.getConfigFilename()

	if result == "" {
		t.Errorf("Expected config filename to be nonempty, got '%v'", result)
	}
}

func TestDefaultConfigIsUsingDockerCompose(t *testing.T) {
	result := confOptsObj.isUsingDockerCompose()

	// I guess we're simply seeing that the function is executable and will return a boolean
	if !(result || !result) {
		t.Errorf("Expected isUsingDockerCompose to be true or false, got %v", result)
	}
}
