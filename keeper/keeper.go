package keeper

import (
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Codec
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

func (k Keeper) CreateDenom(ctx sdk.Context,
	id, name, schema string,
	creator sdk.AccAddress) error {
	return k.SetDenom(ctx, types.NewDenom(id, name, schema, creator))
}

func (k Keeper) MintONFT(ctx sdk.Context, denomID, onftID string, metadata *types.Metadata, assetType types.AssetType,
	transferable bool, owner sdk.AccAddress) error {
	if !k.HasPermissionToMint(ctx, denomID, owner) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "only creator of denom has permission to mint")
	}
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	if k.HasONFT(ctx, denomID, onftID) {
		return sdkerrors.Wrapf(types.ErrONFTAlreadyExists, "ONFT %s already exists in collection %s", onftID, denomID)
	}

	k.setONFT(ctx, denomID, types.NewONFT(
		onftID,
		metadata,
		assetType,
		transferable,
		owner,
	))
	k.increaseSupply(ctx, denomID)
	return nil
}

func (k Keeper) EditONFT(ctx sdk.Context, denomID, onftID string, metadata *types.Metadata, assetType string,
	transferable string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	onft, err := k.Authorize(ctx, denomID, onftID, owner)
	if err != nil {
		return err
	}

	if metadata.Name != types.DoNotModify {
		onft.Metadata.Name = metadata.Name
	}
	if metadata.Description != types.DoNotModify {
		onft.Metadata.Description = metadata.Description
	}
	if metadata.Preview != types.DoNotModify {
		onft.Metadata.Preview = metadata.Preview
	}
	if metadata.Media != types.DoNotModify {
		onft.Metadata.Media = metadata.Media
	}
	if assetType != types.DoNotModify {
		switch _type := strings.ToLower(assetType); _type {
		case "artwork":
			onft.Type = types.ARTWORK
		case "audio":
			onft.Type = types.AUDIO
		case "video":
			onft.Type = types.VIDEO
		default:
			onft.Type = types.ARTWORK
		}
	}
	if transferable != types.DoNotModify {
		switch transferable := strings.ToLower(transferable); transferable {
		case "yes":
			onft.TransferEnabled = true
		case "no":
			onft.TransferEnabled = false
		default:
			onft.TransferEnabled = true
		}
	}

	k.setONFT(ctx, denomID, onft)
	return nil
}

func (k Keeper) TransferOwnership(ctx sdk.Context, denomID, onftID string, srcOwner, dstOwner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	onft, err := k.Authorize(ctx, denomID, onftID, srcOwner)
	if err != nil {
		return err
	}
	if !onft.IsTransferable() {
		return sdkerrors.Wrap(types.ErrNotTransferable, onft.GetID())
	}

	onft.Owner = dstOwner

	k.setONFT(ctx, denomID, onft)
	return nil
}

func (k Keeper) BurnONFT(ctx sdk.Context,
	denomID, onftID string,
	owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	onft, err := k.Authorize(ctx, denomID, onftID, owner)
	if err != nil {
		return err
	}

	k.deleteONFT(ctx, denomID, onft)
	k.decreaseSupply(ctx, denomID)
	return nil
}
