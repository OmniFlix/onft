package onft

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/keeper"
	"github.com/OmniFlix/onft/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateDenom:
			return HandleMsgCreateDenom(ctx, msg, k)
		case *types.MsgMintONFT:
			return HandleMsgMintONFT(ctx, msg, k)
		case *types.MsgTransferONFT:
			return HandleMsgTransferONFT(ctx, msg, k)
		case *types.MsgEditONFT:
			return HandleMsgEditONFT(ctx, msg, k)
		case *types.MsgBurnONFT:
			return HandleMsgBurnONFT(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized onft message type: %T", msg)
		}
	}
}

func HandleMsgCreateDenom(ctx sdk.Context, msg *types.MsgCreateDenom, k keeper.Keeper,
) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	name := strings.ToLower(strings.TrimSpace(msg.Name))
	symbol := strings.ToLower(strings.TrimSpace(msg.Symbol))
	if err := k.CreateDenom(ctx,
		id,
		symbol,
		name,
		msg.Schema,
		msg.Sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func HandleMsgTransferONFT(ctx sdk.Context, msg *types.MsgTransferONFT, k keeper.Keeper,
) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	if err := k.TransferOwnership(ctx,
		denom,
		id,
		msg.Sender,
		msg.Recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyDenomID, denom),
			sdk.NewAttribute(types.AttributeKeyONFTID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func HandleMsgEditONFT(ctx sdk.Context, msg *types.MsgEditONFT, k keeper.Keeper) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	if err := k.EditONFT(ctx,
		denom,
		id,
		msg.Metadata,
		msg.AssetType,
		msg.Transferable,
		msg.Sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditNFT,
			sdk.NewAttribute(types.AttributeKeyDenomID, denom),
			sdk.NewAttribute(types.AttributeKeyONFTID, id),
			sdk.NewAttribute(types.AttributeKeyMediaURI, msg.Metadata.Media),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func HandleMsgMintONFT(ctx sdk.Context, msg *types.MsgMintONFT, k keeper.Keeper,
) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	if err := k.MintONFT(ctx,
		denom,
		id,
		msg.Metadata,
		msg.AssetType,
		msg.Transferable,
		msg.Sender,
		msg.Recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyDenomID, denom),
			sdk.NewAttribute(types.AttributeKeyONFTID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func HandleMsgBurnONFT(ctx sdk.Context, msg *types.MsgBurnONFT, k keeper.Keeper,
) (*sdk.Result, error) {
	id := strings.ToLower(strings.TrimSpace(msg.Id))
	denom := strings.ToLower(strings.TrimSpace(msg.Denom))

	if err := k.BurnONFT(ctx,
		denom,
		id,
		msg.Sender,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyDenomID, denom),
			sdk.NewAttribute(types.AttributeKeyONFTID, id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
