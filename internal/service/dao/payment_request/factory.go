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
		ProposalHash:        paymentRequest.ProposalHash,
		Description:         paymentRequest.Description,
		RequestedAmount:     convertStringToFloat(paymentRequest.RequestedAmount),
		Status:              "pending",
		State:               paymentRequest.State,
		StateChangedOnBlock: paymentRequest.StateChangedOnBlock,
		PaidOnBlock:         paymentRequest.PaidOnBlock,
		Height:              height,
		UpdatedOnBlock:      height,
	}
}

func UpdatePaymentRequest(paymentRequest navcoind.PaymentRequest, height uint64, p *explorer.PaymentRequest) {
	if p.Status != explorer.PaymentRequestStatus(paymentRequest.Status) {
		p.Status = explorer.PaymentRequestStatus(paymentRequest.Status)
		p.UpdatedOnBlock = height
	}
	if p.State != paymentRequest.State {
		p.State = paymentRequest.State
		p.UpdatedOnBlock = height
	}
	if p.StateChangedOnBlock != paymentRequest.StateChangedOnBlock {
		p.StateChangedOnBlock = paymentRequest.StateChangedOnBlock
		p.UpdatedOnBlock = height
	}
	if p.PaidOnBlock != paymentRequest.PaidOnBlock {
		p.PaidOnBlock = paymentRequest.PaidOnBlock
		p.UpdatedOnBlock = height
	}
}

func convertStringToFloat(input string) float64 {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
