package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
	"unicode/utf8"
)

const (
	TypeMsgCreateDenom  = "create_denom"
	TypeMsgMintONFT     = "mint_onft"
	TypeMsgEditONFT     = "edit_onft"
	TypeMsgTransferONFT = "transfer_onft"
	TypeMsgBurnONFT     = "burn_onft"
)

var (
	_ sdk.Msg = &MsgCreateDenom{}
	_ sdk.Msg = &MsgMintONFT{}
	_ sdk.Msg = &MsgEditONFT{}
	_ sdk.Msg = &MsgTransferONFT{}
	_ sdk.Msg = &MsgBurnONFT{}
)

func NewMsgCreateDenom(symbol, name, schema, sender string) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender: sender,
		Id:     GenUniqueID("onftdenom"),
		Symbol: symbol,
		Name:   name,
		Schema: schema,
	}
}

func (msg MsgCreateDenom) Route() string { return RouterKey }

func (msg MsgCreateDenom) Type() string { return TypeMsgCreateDenom }

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

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgMintONFT(denom, sender, recipient string, metadata Metadata, assetType AssetType, transferable bool) *MsgMintONFT {

	return &MsgMintONFT{
		Id:           GenUniqueID("onft"),
		Denom:        denom,
		Metadata:     metadata,
		AssetType:    assetType,
		Transferable: transferable,
		Sender:       sender,
		Recipient:    recipient,
	}
}

func (msg MsgMintONFT) Route() string { return RouterKey }

func (msg MsgMintONFT) Type() string { return TypeMsgMintONFT }

func (msg MsgMintONFT) ValidateBasic() error {

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address; %s", err)
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
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgTransferONFT(id, denom, sender, recipient string) *MsgTransferONFT {

	return &MsgTransferONFT{
		Id:        strings.ToLower(strings.TrimSpace(id)),
		Denom:     strings.TrimSpace(denom),
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg MsgTransferONFT) Route() string { return RouterKey }

func (msg MsgTransferONFT) Type() string { return TypeMsgTransferONFT }

func (msg MsgTransferONFT) ValidateBasic() error {
	if err := ValidateDenomID(msg.Denom); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address; %s", err)
	}
	return ValidateONFTID(msg.Id)
}

func (msg MsgTransferONFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgTransferONFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgEditONFT(id, denom string, metadata Metadata, assetType, transferable, sender string) *MsgEditONFT {
	return &MsgEditONFT{
		Id:           id,
		Denom:        denom,
		Metadata:     metadata,
		AssetType:    assetType,
		Transferable: transferable,
		Sender:       sender,
	}
}

func (msg MsgEditONFT) Route() string { return RouterKey }

func (msg MsgEditONFT) Type() string { return TypeMsgEditONFT }

func (msg MsgEditONFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
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
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func NewMsgBurnONFT(denom, id, sender string) *MsgBurnONFT {
	return &MsgBurnONFT{
		Denom:  denom,
		Id:     id,
		Sender: sender,
	}
}

func (msg MsgBurnONFT) Route() string { return RouterKey }

func (msg MsgBurnONFT) Type() string { return TypeMsgBurnONFT }

func (msg MsgBurnONFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address; %s", err)
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
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
