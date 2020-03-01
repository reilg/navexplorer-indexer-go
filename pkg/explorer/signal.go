package explorer

type Signal struct {
	MetaData MetaData `json:"-"`

	Address   string   `json:"address"`
	Height    uint64   `json:"height"`
	SoftForks []string `json:"softforks"`
}

func (s *Signal) IsSignalling() bool {
	return len(s.SoftForks) > 0
}

func (s *Signal) DeleteSoftFork(name string) {
	softForks := make([]string, 0)

	for i := range s.SoftForks {
		if s.SoftForks[i] != name {
			softForks = append(softForks, s.SoftForks[i])
		}
	}

	s.SoftForks = softForks
}
