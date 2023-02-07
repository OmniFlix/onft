package keeper

import (
	"context"
	"github.com/OmniFlix/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if err := m.Keeper.SaveDenom(ctx,
		msg.Id,
		msg.Symbol,
		msg.Name,
		msg.Schema,
		sender,
		msg.Description,
		msg.PreviewURI,
		msg.Uri,
		msg.UriHash,
		msg.Data,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventCreateDenom{
			Id:      msg.Id,
			Symbol:  msg.Symbol,
			Name:    msg.Name,
			Creator: msg.Sender,
		},
	)

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	// TODO: update denom functionality
	/*
		err = m.Keeper.UpdateDenom(ctx, msg.Id, msg.Name, msg.Description, msg.PreviewURI, sender)
		if err != nil {
			return nil, err
		}

		ctx.EventManager().EmitTypedEvent(
			&types.EventUpdateDenom{
				Id:      msg.Id,
				Name:    msg.Name,
				Creator: msg.Sender,
			},
		)
	*/
	return &types.MsgUpdateDenomResponse{}, nil
}

func (m msgServer) TransferDenom(goCtx context.Context, msg *types.MsgTransferDenom) (*types.MsgTransferDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.TransferDenomOwner(ctx, msg.Id, sender, recipient)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferDenom{
			Id:        msg.Id,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

	return &types.MsgTransferDenomResponse{}, nil
}

func (m msgServer) MintONFT(goCtx context.Context,
	msg *types.MsgMintONFT) (*types.MsgMintONFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if !m.Keeper.HasPermissionToMint(ctx, msg.DenomId, sender) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to mint nft under denom %s",
			sender.String(),
			msg.DenomId,
		)
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	if m.Keeper.HasONFT(ctx, msg.DenomId, msg.Id) {
		return nil, sdkerrors.Wrapf(
			types.ErrONFTAlreadyExists,
			"ONFT %s already exists in collection %s", msg.Id, msg.DenomId)
	}

	if err := m.Keeper.SaveNFT(ctx,
		msg.DenomId,
		msg.Id,
		msg.Metadata.Name,
		msg.Metadata.Description,
		msg.Metadata.MediaURI,
		msg.Metadata.UriHash,
		msg.Metadata.PreviewURI,
		msg.Data,
		ctx.BlockTime(),
		msg.Transferable,
		msg.Extensible,
		msg.Nsfw,
		msg.RoyaltyShare,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventMintONFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			URI:     msg.Metadata.MediaURI,
			Owner:   msg.Recipient,
		},
	)

	return &types.MsgMintONFTResponse{}, nil
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
	if err := m.Keeper.TransferOwnership(ctx, msg.DenomId, msg.Id,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferONFT{
			Id:        msg.Id,
			DenomId:   msg.DenomId,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

	return &types.MsgTransferONFTResponse{}, nil
}

func (m msgServer) BurnONFT(goCtx context.Context,
	msg *types.MsgBurnONFT) (*types.MsgBurnONFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnONFT(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventBurnONFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			Owner:   msg.Sender,
		},
	)

	return &types.MsgBurnONFTResponse{}, nil
}
