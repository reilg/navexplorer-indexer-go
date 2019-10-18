package explorer

import "math"

type Block struct {
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

	MetaData struct {
		Stake       uint64 `json:"stake"`
		StakedBy    string `json:"stakedBy"`
		Spend       uint64 `json:"spend"`
		Fees        uint64 `json:"fees"`
		CFundPayout uint64 `json:"cfundPayout"`
	}
}

func (b *Block) BlockCycle(size uint) uint {
	return GetCycleForHeight(b.Height, size)
}

func (b *Block) CycleIndex(size uint) uint {
	end := size * b.BlockCycle(size)
	start := end - size

	return uint(b.Height) - start
}

func GetCycleForHeight(height uint64, size uint) uint {
	return uint(math.Ceil(float64(uint(height)/(size+1)))) + 1
}
