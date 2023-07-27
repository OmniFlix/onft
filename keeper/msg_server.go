package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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

func (m msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	denomCreationFee := m.Keeper.GetDenomCreationFee(ctx)
	if !msg.CreationFee.Equal(denomCreationFee) {
		if msg.CreationFee.Denom != denomCreationFee.Denom {
			return nil, errorsmod.Wrapf(types.ErrInvalidFeeDenom, "invalid creation fee denom %s",
				msg.CreationFee.Denom)
		}
		if msg.CreationFee.Amount.LT(denomCreationFee.Amount) {
			return nil, errorsmod.Wrapf(
				types.ErrNotEnoughFeeAmount,
				"%s fee is not enough, to create %s fee is required",
				msg.CreationFee.String(),
				denomCreationFee.String(),
			)
		}
		return nil, errorsmod.Wrapf(
			types.ErrInvalidDenomCreationFee,
			"given fee (%s) not matched with  denom creation fee. %s required to create onft denom",
			msg.CreationFee.String(),
			denomCreationFee.String(),
		)
	}
	if err := m.Keeper.CreateDenom(ctx,
		msg.Id,
		msg.Symbol,
		msg.Name,
		msg.Schema,
		sender,
		msg.Description,
		msg.PreviewURI,
		msg.CreationFee,
	); err != nil {
		return nil, err
	}

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.UpdateDenom(ctx, msg.Id, msg.Name, msg.Description, msg.PreviewURI, sender)
	if err != nil {
		return nil, err
	}

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

	return &types.MsgTransferDenomResponse{}, nil
}

func (m msgServer) MintONFT(goCtx context.Context, msg *types.MsgMintONFT) (*types.MsgMintONFTResponse, error) {
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
		msg.DenomId,
		msg.Id,
		msg.Metadata,
		msg.Data,
		msg.Transferable,
		msg.Extensible,
		msg.Nsfw,
		msg.RoyaltyShare,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	return &types.MsgMintONFTResponse{}, nil
}

func (m msgServer) TransferONFT(goCtx context.Context,
	msg *types.MsgTransferONFT,
) (*types.MsgTransferONFTResponse, error) {
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

	return &types.MsgTransferONFTResponse{}, nil
}

func (m msgServer) BurnONFT(goCtx context.Context,
	msg *types.MsgBurnONFT,
) (*types.MsgBurnONFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnONFT(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	return &types.MsgBurnONFTResponse{}, nil
}
