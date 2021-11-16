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

func (k Keeper) CreateDenom(
	ctx sdk.Context, id, symbol, name, schema string,
	creator sdk.AccAddress, description, previewUri string) error {

	if k.HasDenomID(ctx, id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", id)
	}

	if k.HasDenomSymbol(ctx, symbol) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomSymbol %s has already exists", symbol)
	}
	return k.SetDenom(ctx, types.NewDenom(id, symbol, name, schema, creator, description, previewUri))
}

func (k Keeper) UpdateDenom(ctx sdk.Context, id, name, description, previewURI string, sender sdk.AccAddress) error {
	if !k.HasDenomID(ctx, id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom id %s not exists", id)
	}
	denom, err := k.AuthorizeDenomCreator(ctx, id, sender)
	if err != nil {
		return err
	}
	if len(name) > 0 && name != types.DoNotModify {
		denom.Name = name
	}
	if len(description) > 0 && description != types.DoNotModify {
		denom.Description = description
	}
	if len(previewURI) > 0 && previewURI != types.DoNotModify {
		denom.PreviewURI = previewURI
	}
	return k.SetDenom(ctx, denom)
}

func (k Keeper) TransferDenomOwner(ctx sdk.Context, id string, curOwner, newOwner sdk.AccAddress) error {
	denom, err := k.GetDenom(ctx, id)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", id)
	}

	if curOwner.String() != denom.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", curOwner.String())
	}

	denom.Creator = newOwner.String()

	err = k.SetDenom(ctx, denom)
	if err != nil {
		return err
	}
	return nil
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
	k.setOwner(ctx, denomID, onftID, recipient)
	k.increaseSupply(ctx, denomID)
	return nil
}

func (k Keeper) EditONFT(ctx sdk.Context, denomID, onftID string, metadata types.Metadata,
	data, transferable, extensible string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return err
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
	if data != types.DoNotModify {
		onft.Data = data
	}
	if transferable != types.DoNotModify {
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
	if len(extensible) > 0 && extensible != types.DoNotModify {
		if denom.Creator != onft.Owner {
			return sdkerrors.Wrapf(
				types.ErrNotEditable,
				"onft %s: extensibility can be modified only when creator is the owner of the onft.",
				onftID,
			)
		}
		switch extensible := strings.ToLower(extensible); extensible {
		case "yes":
			onft.Extensible = true
		case "no":
			onft.Extensible = false
		default:
			return sdkerrors.Wrapf(
				types.ErrInvalidOption,
				"%s is invalid option for extensible.",
				extensible,
			)
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
	k.swapOwner(ctx, denomID, onftID, srcOwner, dstOwner)
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
	k.deleteOwner(ctx, denomID, onftID, owner)
	k.decreaseSupply(ctx, denomID)
	return nil
}
