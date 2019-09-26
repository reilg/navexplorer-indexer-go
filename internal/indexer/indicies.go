package indexer

import "fmt"

type Index string

var (
	BlockIndex            Index = "block"
	BlockTransactionIndex Index = "blocktransaction"
)

// Sets the network and returns the full string
func (i *Index) Get(network string) string {
	return fmt.Sprintf("%s.%s", network, string(*i))
}
