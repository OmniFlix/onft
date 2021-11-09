package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/types"
)

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	for _, onft := range collection.ONFTs {
		metadata := types.Metadata{
			Name:        onft.GetName(),
			Description: onft.GetDescription(),
			MediaURI:       onft.GetMediaURI(),
			PreviewURI:     onft.GetPreviewURI(),
		}

		if err := k.MintONFT(ctx,
			collection.Denom.Id,
			onft.GetID(),
			metadata,
			onft.GetData(),
			onft.IsTransferable(),
			onft.IsExtensible(),
			onft.GetOwner(),
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
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	onfts := k.GetONFTs(ctx, denomID)
	return types.NewCollection(denom, onfts), nil
}

func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection) {
	for _, denom := range k.GetDenoms(ctx) {
		onfts := k.GetONFTs(ctx, denom.Id)
		cs = append(cs, types.NewCollection(denom, onfts))
	}
	return cs
}

func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(owner, id, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		supply++
	}
	return supply
}

func (k Keeper) increaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}

func (k Keeper) decreaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply--

	store := ctx.KVStore(k.storeKey)
	if supply == 0 {
		store.Delete(types.KeyCollection(denomID))
		return
	}

	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}
