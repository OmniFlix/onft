package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QuerySupply     = "supply"
	QueryCollection = "collection"
	QueryDenoms     = "denoms"
	QueryDenom      = "denom"
	QueryONFT       = "onft"
)

type QuerySupplyParams struct {
	Denom string
	Owner sdk.AccAddress
}

func (q QuerySupplyParams) Bytes() []byte {
	return []byte(q.Denom)
}

type QueryOwnerParams struct {
	Denom string
	Owner sdk.AccAddress
}

type QueryCollectionParams struct {
	Denom string
}

type QueryDenomParams struct {
	ID string
}

type QueryONFTParams struct {
	Denom  string
	ONFTID string
}
