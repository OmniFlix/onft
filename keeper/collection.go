package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/OmniFlix/onft/exported"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/OmniFlix/onft/types"
)

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	denom := collection.Denom
	creator, err := sdk.AccAddressFromBech32(denom.Creator)
	if err != nil {
		return err
	}
	if err := k.SaveDenom(ctx,
		denom.Id,
		denom.Symbol,
		denom.Name,
		denom.Schema,
		creator,
		denom.Description,
		denom.PreviewURI,
		denom.Uri,
		denom.UriHash,
		denom.Data,
	); err != nil {
		return err
	}

	for _, onft := range collection.ONFTs {
		if err := k.SaveNFT(ctx,
			denom.Id,
			onft.GetID(),
			onft.GetName(),
			onft.GetDescription(),
			onft.GetMediaURI(),
			onft.GetURIHash(),
			onft.GetPreviewURI(),
			onft.GetData(),
			onft.GetCreatedTime(),
			onft.IsTransferable(),
			onft.IsExtensible(),
			onft.IsNSFW(),
			onft.GetRoyaltyShare(),
			onft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetCollection(ctx sdk.Context, denomID string) (types.Collection, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Collection{}, errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	onfts, err := k.GetONFTs(ctx, denomID)
	if err != nil {
		return types.Collection{}, err
	}
	return types.NewCollection(denom, onfts), nil
}

func (k Keeper) GetCollections(ctx sdk.Context) (collections []types.Collection, err error) {
	for _, class := range k.nk.GetClasses(ctx) {
		onfts, err := k.GetONFTs(ctx, class.Id)
		if err != nil {
			return nil, err
		}

		denom, err := k.GetDenomInfo(ctx, class.Id)
		if err != nil {
			return nil, err
		}

		collections = append(collections, types.NewCollection(*denom, onfts))
	}
	return collections, nil
}

func (k Keeper) GetPaginateCollection(ctx sdk.Context,
	request *types.QueryCollectionRequest, denomId string,
) (types.Collection, *query.PageResponse, error) {
	denom, err := k.GetDenomInfo(ctx, denomId)
	if err != nil {
		return types.Collection{}, nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denomId %s not existed ", denomId)
	}
	var oNFTs []exported.ONFTI
	store := ctx.KVStore(k.storeKey)
	onftStore := prefix.NewStore(store, types.KeyONFT(denomId, ""))
	pagination, err := query.Paginate(onftStore, request.Pagination, func(key []byte, value []byte) error {
		var oNFT types.ONFT
		k.cdc.MustUnmarshal(value, &oNFT)
		oNFTs = append(oNFTs, oNFT)
		return nil
	})
	if err != nil {
		return types.Collection{}, nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return types.NewCollection(*denom, oNFTs), pagination, nil
}

func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	return k.nk.GetTotalSupply(ctx, denomID)
}

// GetBalance returns the amount of NFTs owned in a class by an account
func (k Keeper) GetBalance(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	return k.nk.GetBalance(ctx, id, owner)
}
