package signal

import "github.com/NavExplorer/navexplorer-indexer-go/v2/pkg/explorer"

func CreateSignal(block *explorer.Block, softForks *explorer.SoftForks) *explorer.Signal {
	signal := &explorer.Signal{Address: block.StakedBy, Height: block.Height}

	for _, sf := range *softForks {
		if sf.IsOpen() && block.Version>>sf.SignalBit&1 == 1 {
			signal.SoftForks = append(signal.SoftForks, sf.Name)
		}
	}

	if len(signal.SoftForks) == 0 {
		return nil
	}

	return signal
}
