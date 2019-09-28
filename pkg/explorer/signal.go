package explorer

type Signal struct {
	Address string   `json:"address"`
	Height  uint64   `json:"height"`
	Signals []string `json:"signals"`
}
