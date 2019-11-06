package events

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/address_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/dao_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/signal_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/softfork_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
)

type Subscriber struct {
	elastic               *index.Index
	addressIndexer        *address_indexer.Indexer
	proposalIndexer       *dao_indexer.ProposalIndexer
	paymentRequestIndexer *dao_indexer.PaymentRequestIndexer
	signalIndexer         *signal_indexer.Indexer
	softForkIndexer       *softfork_indexer.Indexer
}

func NewSubscriber(
	elastic *index.Index,
	addressIndexer *address_indexer.Indexer,
	proposalIndexer *dao_indexer.ProposalIndexer,
	paymentRequestIndexer *dao_indexer.PaymentRequestIndexer,
	signalIndexer *signal_indexer.Indexer,
	softForkIndexer *softfork_indexer.Indexer,
) *Subscriber {
	return &Subscriber{
		elastic,
		addressIndexer,
		proposalIndexer,
		paymentRequestIndexer,
		signalIndexer,
		softForkIndexer,
	}
}

func (s *Subscriber) subscribe() {
	event.On(string(config.EventBlockIndexed), event.ListenerFunc(func(e event.Event) error {
		block := e.Get("block").(*explorer.Block)
		txs := e.Get("txs").(*[]explorer.BlockTransaction)

		s.addressIndexer.IndexAddressesForTransactions(txs)
		s.proposalIndexer.IndexProposals(txs)
		s.paymentRequestIndexer.IndexPaymentRequests(txs)
		s.signalIndexer.IndexSignal(block)

		return nil
	}), event.Normal)

	event.On(string(config.EventSignalIndexed), event.ListenerFunc(func(e event.Event) error {
		signal := e.Get("signal").(*explorer.Signal)
		block := e.Get("block").(*explorer.Block)

		s.softForkIndexer.Update(signal, block)
		s.elastic.PersistRequests(signal.Height)
		return nil
	}), event.Normal)
}
