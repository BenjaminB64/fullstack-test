package domain

import "time"

type BridgeService interface {
	GetBridgeSchedule() (*BridgeSchedule, error)
}

type BridgeSchedule struct {
	Closures []BridgeClosure
}

type BridgeClosure struct {
	BoatName   string
	CloseTime  time.Time
	ReopenTime time.Time
}
