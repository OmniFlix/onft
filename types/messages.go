package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"strings"
	"unicode/utf8"
)

func NewMsgCreateDenom(symbol, name, schema string, sender sdk.AccAddress) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender: sender,
		Id:     fmt.Sprintf("onftdenom%s", strings.ReplaceAll(uuid.New().String(), "-", "")),
		Symbol: strings.ToLower(strings.TrimSpace(symbol)),
		Name:   strings.TrimSpace(name),
		Schema: strings.TrimSpace(schema),
	}
}

func (msg MsgCreateDenom) Route() string { return RouterKey }

func (msg MsgCreateDenom) Type() string { return "create_denom" }

func (msg MsgCreateDenom) ValidateBasic() error {
	if err := ValidateDenomID(msg.Id); err != nil {
		return err
	}
	if err := ValidateDenomSymbol(msg.Symbol); err != nil {
		return err
	}
	name := strings.TrimSpace(msg.Name)
	if len(name) > 0 && !utf8.ValidString(name) {
		return sdkerrors.Wrap(ErrInvalidDenom, "denom name is invalid")
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func NewMsgMintONFT(denom string, metadata *Metadata, assetType AssetType, transferable bool, sender, recipient sdk.AccAddress) *MsgMintONFT {

	return &MsgMintONFT{
		Id:           fmt.Sprintf("onft%s", strings.ReplaceAll(uuid.New().String(), "-", "")),
		Denom:        strings.TrimSpace(denom),
		Metadata:     metadata,
		AssetType:    assetType,
		Transferable: transferable,
		Sender:       sender,
		Recipient:    recipient,
	}
}

func (msg MsgMintONFT) Route() string { return RouterKey }

func (msg MsgMintONFT) Type() string { return "mint_onft" }

func (msg MsgMintONFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing receipt address")
	}
	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if err := ValidateMediaURI(msg.Metadata.Media); err != nil {
		return err
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgMintONFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMintONFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func NewMsgTransferONFT(id, denom string, sender, recipient sdk.AccAddress) *MsgTransferONFT {

	return &MsgTransferONFT{
		Id:        strings.ToLower(strings.TrimSpace(id)),
		Denom:     strings.TrimSpace(denom),
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg MsgTransferONFT) Route() string { return RouterKey }

func (msg MsgTransferONFT) Type() string { return "transfer_onft" }

func (msg MsgTransferONFT) ValidateBasic() error {
	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if msg.Recipient.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgTransferONFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTransferONFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func NewMsgEditONFT(id, denom string, metadata *Metadata, assetType string, transferable string, sender sdk.AccAddress) *MsgEditONFT {
	return &MsgEditONFT{
		Id:           strings.ToLower(strings.TrimSpace(id)),
		Denom:        strings.TrimSpace(denom),
		Metadata:     metadata,
		AssetType:    assetType,
		Transferable: transferable,
		Sender:       sender,
	}
}

func (msg MsgEditONFT) Route() string { return RouterKey }

func (msg MsgEditONFT) Type() string { return "edit_onft" }

func (msg MsgEditONFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if err := ValidateMediaURI(msg.Metadata.Media); err != nil {
		return err
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgEditONFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEditONFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func NewMsgBurnONFT(sender sdk.AccAddress, id string, denom string) *MsgBurnONFT {
	return &MsgBurnONFT{
		Sender: sender,
		Id:     strings.ToLower(strings.TrimSpace(id)),
		Denom:  strings.TrimSpace(denom),
	}
}

func (msg MsgBurnONFT) Route() string { return RouterKey }

func (msg MsgBurnONFT) Type() string { return "burn_onft" }

func (msg MsgBurnONFT) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgBurnONFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgBurnONFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
