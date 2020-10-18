package payment_request

import "github.com/NavExplorer/navexplorer-indexer-go/v2/pkg/explorer"

var PaymentRequests paymentRequests

type paymentRequests []*explorer.PaymentRequest

func (p *paymentRequests) Delete(hash string) {
	for i := range PaymentRequests {
		if PaymentRequests[i].Hash == hash {
			PaymentRequests[i] = PaymentRequests[len(PaymentRequests)-1]                                           // Copy last element to index i.
			PaymentRequests[len(PaymentRequests)-1] = nil                                                          // Erase last element (write zero value).
			PaymentRequests = append([]*explorer.PaymentRequest(nil), PaymentRequests[:len(PaymentRequests)-1]...) // Truncate slice.
			break
		}
	}
}
