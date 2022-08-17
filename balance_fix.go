package main

import (
	"github.com/navcoin/navexplorer-indexer-go/v2/generated/dic"
	"github.com/navcoin/navexplorer-indexer-go/v2/internal/config"
)

var container *dic.Container

func main() {
	config.Init()

	container, _ = dic.NewContainer()
	container.GetElastic().InstallMappings()
	container.GetSoftforkService().InitSoftForks()
	container.GetDaoConsensusService().InitConsensusParameters()

	addressRepo := container.GetAddressRepo()
	rewinder := container.GetAddressRewinder()

	{
		addresses, err := addressRepo.GetAllAddressesWithBalance("spendable")
		if err == nil {
			for i := range addresses {
				rewinder.ResetAddress(addresses[i])
			}
		}
	}
	{
		addresses, err := addressRepo.GetAllAddressesWithBalance("stakable")
		if err == nil {
			for i := range addresses {
				rewinder.ResetAddress(addresses[i])
			}
		}
	}
	{
		addresses, err := addressRepo.GetAllAddressesWithBalance("voting_weight")
		if err == nil {
			for i := range addresses {
				rewinder.ResetAddress(addresses[i])
			}
		}
	}
}
