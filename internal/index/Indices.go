package index

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
)

type Indices string

var (
	AddressTransactionIndex Indices = "addresstransaction"
	BlockIndex              Indices = "block"
	BlockTransactionIndex   Indices = "blocktransaction"
	ProposalIndex           Indices = "proposal"
	ProposalVoteIndex       Indices = "proposalvotes"
	PaymentRequestIndex     Indices = "paymentrequest"
	SignalIndex             Indices = "signal"
	SoftForkIndex           Indices = "softfork"
)

// Sets the network and returns the full string
func (i *Indices) Get() string {
	return fmt.Sprintf("%s.%s", config.Get().Network, string(*i))
}

func All() []Indices {
	return []Indices{
		AddressTransactionIndex,
		BlockIndex,
		BlockTransactionIndex,
		ProposalIndex,
		PaymentRequestIndex,
		SignalIndex,
		SoftForkIndex,
	}
}
