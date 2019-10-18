package explorer

type SoftForks []SoftFork

type SoftFork struct {
	Name             string          `json:"name"`
	SignalBit        uint            `json:"signalBit"`
	State            SoftForkState   `json:"state"`
	LockedInHeight   uint64          `json:"lockedinheight,omitempty"`
	ActivationHeight uint64          `json:"activationheight,omitempty"`
	SignalHeight     uint64          `json:"signalheight,omitempty"`
	Cycles           []SoftForkCycle `json:"cycles,omitempty"`
}

type SoftForkCycle struct {
	Cycle            uint `json:"cycle"`
	BlocksSignalling uint `json:"blocks"`
}

func (s *SoftFork) IsOpen() bool {
	return s.State == SoftForkDefined || s.State == SoftForkStarted || s.State == SoftForkFailed
}

func (s *SoftFork) GetCycle(cycle uint) *SoftForkCycle {
	for i, _ := range s.Cycles {
		if (s.Cycles)[i].Cycle == cycle {
			return &s.Cycles[i]
		}
	}
	return nil
}

func (s SoftForks) GetSoftFork(name string) *SoftFork {
	for i, _ := range s {
		if s[i].Name == name {
			return &s[i]
		}
	}
	return nil
}
