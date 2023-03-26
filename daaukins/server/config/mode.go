package config

type Mode string

const (
	ModeLeader   Mode = "leader"
	ModeFollower Mode = "follower"
)

func (m Mode) String() string {
	return string(m)
}
