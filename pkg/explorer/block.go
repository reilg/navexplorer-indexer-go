package explorer

import "math"

type RawBlock struct {
	Hash              string   `json:"hash"`
	Confirmations     uint64   `json:"confirmations"`
	StrippedSize      uint64   `json:"strippedsize"`
	Size              uint64   `json:"size"`
	Weight            uint64   `json:"weight"`
	Height            uint64   `json:"height"`
	Version           uint32   `json:"version"`
	VersionHex        string   `json:"versionHex"`
	Merkleroot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              int64    `json:"time"`
	MedianTime        int64    `json:"mediantime"`
	Nonce             uint64   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork,omitempty"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
}

type Block struct {
	RawBlock
	Stake       uint64 `json:"stake"`
	StakedBy    string `json:"stakedBy"`
	Spend       uint64 `json:"spend"`
	Fees        uint64 `json:"fees"`
	CFundPayout uint64 `json:"cfundPayout"`

	// Transient
	Best bool `json:"best,omitempty"`
}

type BlockCycle struct {
	Cycle  uint
	Index  uint
	Quorum uint
}

func (b *Block) BlockCycle(size uint) BlockCycle {
	cycle := GetCycleForHeight(b.Height, size)

	return BlockCycle{
		Cycle:  cycle,
		Index:  GetCycleIndex(b.Height, cycle, size),
		Quorum: GetQuorum(size),
	}
}

func GetCycleForHeight(height uint64, size uint) uint {
	return uint(math.Ceil(float64(uint(height)/(size+1)))) + 1
}

func GetCycleIndex(height uint64, cycle uint, size uint) uint {
	return uint(height) - (cycle * size) - size
}

func GetQuorum(size uint) uint {
	return uint(float64(size) * 0.75)
}
