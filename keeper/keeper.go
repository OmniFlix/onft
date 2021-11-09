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
	cdc      codec.BinaryCodec
}

func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

func (k Keeper) CreateDenom(ctx sdk.Context,
	id, symbol, name, schema string,
	creator sdk.AccAddress, description, previewUri string) error {
	return k.SetDenom(ctx, types.NewDenom(id, symbol, name, schema, creator, description, previewUri))
}
func (k Keeper) UpdateDenom(ctx sdk.Context) {
	// TODO: Implement Update Denom keeper functionality
}

func (k Keeper) TransferDenomOwner(ctx sdk.Context) {
	// TODO: impletment transfer denom
}

func (k Keeper) MintONFT(
	ctx sdk.Context, denomID, onftID string,
	metadata types.Metadata, data string, transferable, extensible bool,
	sender, recipient sdk.AccAddress) error {
	if !k.HasPermissionToMint(ctx, denomID, sender) {
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
		data,
		transferable,
		extensible,
		recipient,
		ctx.BlockHeader().Time,
	))
	k.increaseSupply(ctx, denomID)
	return nil
}

func (k Keeper) EditONFT(ctx sdk.Context, denomID, onftID string, metadata types.Metadata,
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
	if metadata.PreviewURI != types.DoNotModify {
		onft.Metadata.PreviewURI = metadata.PreviewURI
	}
	if metadata.MediaURI != types.DoNotModify {
		onft.Metadata.MediaURI = metadata.MediaURI
	}
	if transferable != types.DoNotModify {
		denom, err := k.GetDenom(ctx, denomID)
		if err != nil {
			return err
		}
		if denom.Creator != onft.Owner {
			return sdkerrors.Wrapf(
				types.ErrNotEditable,
				"onft %s: transferability can be modified only when creator is the owner of onft.",
				onftID,
			)
		}
		switch transferable := strings.ToLower(transferable); transferable {
		case "yes":
			onft.Transferable = true
		case "no":
			onft.Transferable = false
		default:
			onft.Transferable = true
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

	onft.Owner = dstOwner.String()

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
