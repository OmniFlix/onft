package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/types"
)

func (k Keeper) HasDenomID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomID(id))
}

func (k Keeper) HasDenomName(ctx sdk.Context, name string) bool {
	if len(name) == 0 {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomName(name))
}

func (k Keeper) SetDenom(ctx sdk.Context, denom *types.Denom) error {
	if k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", denom.Id)
	}

	if k.HasDenomName(ctx, denom.Name) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomName %s has already exists", denom.Name)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(denom)
	store.Set(types.KeyDenomID(denom.Id), bz)
	if len(denom.Name) > 0 {
		store.Set(types.KeyDenomName(denom.Name), []byte(denom.Id))
	}
	return nil
}

func (k Keeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenomID(id))
	if bz == nil || len(bz) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denomID: %s", id)
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomID(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		k.cdc.MustUnmarshal(iterator.Value(), &denom)
		denoms = append(denoms, denom)
	}
	return denoms
}