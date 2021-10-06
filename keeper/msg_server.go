package keeper

import (
	"context"
	"github.com/OmniFlix/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the NFT MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateDenom(goCtx context.Context,
	msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CreateDenom(ctx, msg.Id, msg.Symbol, msg.Name, msg.Schema, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDenom,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomName, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyDenomName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) MintONFT(goCtx context.Context,
	msg *types.MsgMintONFT) (*types.MsgMintONFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.MintONFT(ctx,
		msg.Denom, msg.Id,
		msg.Metadata,
		msg.AssetType,
		msg.Transferable,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintONFT,
			sdk.NewAttribute(types.AttributeKeyONFTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyMediaURI, msg.Metadata.Media),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgMintONFTResponse{}, nil
}

func (m msgServer) EditONFT(goCtx context.Context,
	msg *types.MsgEditONFT) (*types.MsgEditONFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.EditONFT(ctx, msg.Denom, msg.Id,
		msg.Metadata,
		msg.AssetType,
		msg.Transferable,
		sender,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditONFT,
			sdk.NewAttribute(types.AttributeKeyONFTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyMediaURI, msg.Metadata.Media),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})
	return &types.MsgEditONFTResponse{}, nil
}

func (m msgServer) TransferONFT(goCtx context.Context,
	msg *types.MsgTransferONFT) (*types.MsgTransferONFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferOwnership(ctx, msg.Denom, msg.Id,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferONFT,
			sdk.NewAttribute(types.AttributeKeyONFTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Denom),
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgTransferONFTResponse{}, nil
}

func (m msgServer) BurnONFT(goCtx context.Context,
	msg *types.MsgBurnONFT) (*types.MsgBurnONFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnONFT(ctx, msg.Denom, msg.Id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnONFT,
			sdk.NewAttribute(types.AttributeKeyDenomID, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyONFTID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnONFTResponse{}, nil
}