package explorer

type SoftFork struct {
	Name             string        `json:"name"`
	SignalBit        uint          `json:"signalBit"`
	State            SoftForkState `json:"state"`
	LockedInHeight   uint64        `json:"lockedinheight,omitempty"`
	ActivationHeight uint64        `json:"activationheight,omitempty"`
	SignalHeight     uint64        `json:"signalheight,omitempty"`
	Cycle            []struct {
		Cycle            int    `json:"cycle"`
		BlocksSignalling uint64 `json:"blocks"`
	}
}

func (s *SoftFork) IsOpen() bool {
	return s.State == SoftForkDefined || s.State == SoftForkStarted || s.State == SoftForkFailed
}
