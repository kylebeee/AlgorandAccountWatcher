package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type SlimmedAccountState struct {
	Amount uint64            `json:"amount"`
	Assets map[uint64]uint64 `json:"assets"`
}

func (s *Server) ProcessAddress(address string) error {
	const operation = "ProcessAddress"

	response, err := s.Algod.AccountInformation(address).Do(context.Background())
	if err != nil {
		return errors.Join(fmt.Errorf("operation: %v", operation), err)
	}

	previousState, ok := s.WatchList.Subs.Load(address)
	if !ok {
		// doesn't yet exist, just save the information

		sas := SlimmedAccountState{
			Amount: response.Amount,
			Assets: map[uint64]uint64{},
		}

		for _, asset := range response.Assets {
			sas.Assets[asset.AssetId] = asset.Amount
		}

		s.WatchList.Subs.Store(address, sas)
		fmt.Printf("[ADDED] - %v\n", address)
		return nil
	}

	changes := []string{}
	ps := previousState.(SlimmedAccountState)

	// check for changes in balances
	if ps.Amount != response.Amount {
		changes = append(changes, fmt.Sprintf("         [0]: %v -> %v\n", ps.Amount, response.Amount))
	}

	newAssetAmounts := map[uint64]uint64{}
	for _, asset := range response.Assets {
		offset := 10 - len(fmt.Sprintf("%v", asset.AssetId))
		newAssetAmounts[asset.AssetId] = asset.Amount
		previousAssetAmount, ok := ps.Assets[asset.AssetId]
		if !ok {
			changes = append(changes, fmt.Sprintf("%v[%v]: 0 -> %v\n", strings.Repeat(" ", offset), asset.AssetId, asset.Amount))
			continue
		}

		if previousAssetAmount != asset.Amount {
			changes = append(changes, fmt.Sprintf("%v[%v]: %v -> %v\n", strings.Repeat(" ", offset), asset.AssetId, previousAssetAmount, asset.Amount))
		}
	}

	if len(changes) > 0 {
		fmt.Printf("\n[BALANCE CHANGE] - %v\n", address)
		for _, change := range changes {
			fmt.Printf("%v", change)
		}
		fmt.Println()
		s.WatchList.Subs.Store(address, SlimmedAccountState{
			Amount: response.Amount,
			Assets: newAssetAmounts,
		})
	}

	return nil
}
