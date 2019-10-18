package explorer

type Signal struct {
	Address   string   `json:"address"`
	Height    uint64   `json:"height"`
	SoftForks []string `json:"softforks"`
}

func (s *Signal) IsSignalling() bool {
	return len(s.SoftForks) > 0
}
