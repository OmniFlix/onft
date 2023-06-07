package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateONFTID(onftId string) error {
	if len(onftId) < MinIDLen || len(onftId) > MaxIDLen {
		return sdkerrors.Wrapf(
			ErrInvalidONFTID,
			"invalid onftId %s, length must be between [%d, %d]", onftId, MinIDLen, MaxIDLen)
	}
	if !IsBeginWithAlpha(onftId) || !IsAlphaNumeric(onftId) {
		return sdkerrors.Wrapf(
			ErrInvalidONFTID,
			"invalid onftId %s, only accepts alphanumeric characters and begin with an english letter", onftId)
	}
	return nil
}

func ValidateDenomID(denomID string) error {
	if len(denomID) < MinIDLen || len(denomID) > MaxIDLen {
		return sdkerrors.Wrapf(
			ErrInvalidDenom,
			"invalid denom ID %s, length  must be between [%d, %d]",
			denomID,
			MinIDLen,
			MaxIDLen,
		)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(
			ErrInvalidDenom,
			"invalid denom ID %s, only accepts alphanumeric characters,and begin with an english letter",
			denomID,
		)
	}
	return nil
}
func ValidateDenomSymbol(denomSymbol string) error {
	if len(denomSymbol) < MinDenomLen || len(denomSymbol) > MaxDenomLen {
		return sdkerrors.Wrapf(
			ErrInvalidDenom,
			"invalid denom symbol %s, only accepts value [%d, %d]",
			denomSymbol,
			MinDenomLen,
			MaxDenomLen,
		)
	}
	if !IsBeginWithAlpha(denomSymbol) || !IsAlpha(denomSymbol) {
		return sdkerrors.Wrapf(
			ErrInvalidDenom,
			"invalid denom symbol %s, only accepts alphabetic characters",
			denomSymbol,
		)
	}
	return nil
}

func ValidateName(name string) error {
	if len(name) > MaxNameLen {
		return sdkerrors.Wrapf(
			ErrInvalidName,
			"invalid name %s, length must be less than %d",
			name,
			MaxNameLen,
		)
	}
	return nil
}

func ValidateDescription(description string) error {
	if len(description) > MaxDescriptionLen {
		return sdkerrors.Wrapf(
			ErrInvalidDescription,
			"invalid description %s, length must be less than %d",
			description,
			MaxDescriptionLen,
		)
	}
	return nil
}

func ValidateURI(uri string) error {
	if len(uri) > MaxURILen {
		return sdkerrors.Wrapf(
			ErrInvalidURI,
			"invalid uri %s, length must be less than %d",
			uri,
			MaxURILen,
		)
	}
	return nil
}

func ValidateMediaURI(uri string) error {
	if len(uri) == 0 {
		return sdkerrors.Wrapf(
			ErrInvalidURI,
			"invalid uri %s, media uri should not be empty",
			uri,
			0,
		)
	}
	if len(uri) > MaxURILen {
		return sdkerrors.Wrapf(
			ErrInvalidURI,
			"invalid uri %s, length must be less than %d",
			uri,
			MaxURILen,
		)
	}
	return nil
}

func ValidateCreationFee(fee sdk.Coin) error {
	if !fee.IsValid() || fee.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidURI, "invalid creation fee %s, fee must be positive", fee.String())
	}
	return nil
}
