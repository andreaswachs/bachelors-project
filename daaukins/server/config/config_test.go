package config

import "testing"

func TestParseLeaderNoMinions(t *testing.T) {
	conf := `server_mode: leader
minions: []`

	parsed, err := parse([]byte(conf))
	if err != nil {
		t.Error(err)
	}

	if parsed.ServerMode != ModeLeader {
		t.Error("Expected server mode to be leader")
	}
}

func TestParseMinionNoMinions(t *testing.T) {
	conf := `server_mode: minion
minions: []`

	parsed, err := parse([]byte(conf))
	if err != nil {
		t.Error(err)
	}

	if parsed.ServerMode != ModeMinion {
		t.Error("Expected server mode to be minion")
	}
}

func TestParseLeaderOneMinion(t *testing.T) {
	conf := `server_mode: leader
minions:
  - address: localhost
    port: 8080`

	parsed, err := parse([]byte(conf))
	if err != nil {
		t.Error(err)
	}

	if parsed.ServerMode != ModeLeader {
		t.Error("Expected server mode to be leader")
	}

	if len(parsed.Minions) != 1 {
		t.Error("Expected one minion")
	}

	if parsed.Minions[0].Address != "localhost" {
		t.Error("Expected address to be localhost")
	}

	if parsed.Minions[0].Port != 8080 {
		t.Error("Expected port to be 8080")
	}
}
