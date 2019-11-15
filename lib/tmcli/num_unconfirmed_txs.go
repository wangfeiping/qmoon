// Copyright 2018 The QOS Authors

package tmcli

import (
	"context"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {

}

const numUnconfirmedTxsURI = "num_unconfirmed_txs"

type numUnconfirmedTxsService service

func (s *numUnconfirmedTxsService) Retrieve(ctx context.Context) (*tmctypes.ResultUnconfirmedTxs, error) {
	u := numUnconfirmedTxsURI

	u, err := addOptions(u, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var res tmctypes.ResultUnconfirmedTxs
	_, err = s.client.Do(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
