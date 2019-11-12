package payment_request

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CreatePaymentRequest(paymentRequest navcoind.PaymentRequest, height uint64) *explorer.PaymentRequest {
	return &explorer.PaymentRequest{
		Version:             paymentRequest.Version,
		Hash:                paymentRequest.Hash,
		BlockHash:           paymentRequest.BlockHash,
		Description:         paymentRequest.Description,
		RequestedAmount:     convertStringToFloat(paymentRequest.RequestedAmount),
		Status:              "pending",
		State:               0,
		StateChangedOnBlock: paymentRequest.StateChangedOnBlock,
		Height:              height,
	}
}

func CreatePaymentRequestVotes(block *explorer.Block, tx *explorer.BlockTransaction) *explorer.PaymentRequestVote {
	if !tx.IsCoinbase() {
		return nil
	}

	paymentRequestVote := &explorer.PaymentRequestVote{Height: tx.Height, Address: block.StakedBy}
	for _, vout := range tx.Vout {
		if !vout.IsPaymentRequestVote() {
			continue
		}

		vote := explorer.Vote{Hash: vout.ScriptPubKey.Hash, Vote: vout.ScriptPubKey.Type == explorer.VoutPaymentRequestYesVote}
		paymentRequestVote.Votes = append(paymentRequestVote.Votes, vote)
	}

	return paymentRequestVote
}

func convertStringToFloat(input string) float64 {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
