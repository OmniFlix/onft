package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Supply(c context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.Denom))
	ctx := sdk.UnwrapSDKContext(c)

	var supply uint64
	switch {
	case len(request.Owner) == 0 && len(denom) > 0:
		supply = k.GetTotalSupply(ctx, denom)
	default:
		onfts := k.GetOwnerONFTs(ctx, denom, request.Owner)
		supply = uint64(len(onfts))
	}
	return &types.QuerySupplyResponse{
		Amount: supply,
	}, nil
}

func (k Keeper) Collection(c context.Context, request *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.Denom))
	ctx := sdk.UnwrapSDKContext(c)

	collection, err := k.GetCollection(ctx, denom)
	if err != nil {
		return nil, err
	}
	return &types.QueryCollectionResponse{
		Collection: &collection,
	}, nil
}

func (k Keeper) Denom(c context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.Denom))
	ctx := sdk.UnwrapSDKContext(c)

	denomObject, err := k.GetDenom(ctx, denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomResponse{
		Denom: &denomObject,
	}, nil
}

func (k Keeper) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	denoms := k.GetDenoms(ctx)
	return &types.QueryDenomsResponse{
		Denoms: denoms,
	}, nil
}

func (k Keeper) ONFT(c context.Context, request *types.QueryONFTRequest) (*types.QueryONFTResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.Denom))
	onftID := strings.ToLower(strings.TrimSpace(request.Id))
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetONFT(ctx, denom, onftID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownONFT, "invalid ONFT %s from collection %s", request.Id, request.Denom)
	}

	oNFT, ok := nft.(types.ONFT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrUnknownONFT, "invalid type NFT %s from collection %s", request.Id, request.Denom)
	}

	return &types.QueryONFTResponse{
		ONFT: &oNFT,
	}, nil
}
func (k Keeper) OwnerONFTs(c context.Context, request *types.QueryOwnerONFTsRequest) (*types.QueryOwnerONFTsResponse, error) {
	denom := strings.ToLower(strings.TrimSpace(request.Denom))
	ctx := sdk.UnwrapSDKContext(c)

	onfts := k.GetOwnerONFTs(ctx, denom, request.Owner)
	return &types.QueryOwnerONFTsResponse{
		Onfts: onfts,
	}, nil
}
