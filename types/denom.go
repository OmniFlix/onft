package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewDenom(id, symbol, name, schema string, creator sdk.AccAddress, description, previewURI string) Denom {
	return Denom{
		Id:      id,
		Symbol:  symbol,
		Name:    name,
		Schema:  schema,
		Creator: creator.String(),
		Description: description,
		PreviewURI: previewURI,
	}
}

func ValidateDenomID(denomID string) error {
	denomID = strings.TrimSpace(denomID)
	if len(denomID) < MinDenomLen || len(denomID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom %s, only accepts value [%d, %d]", denomID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom %s, only accepts alphanumeric characters,and begin with an english letter", denomID)
	}
	return nil
}
func ValidateDenomSymbol(denomSymbol string) error {
	denomSymbol = strings.TrimSpace(denomSymbol)
	if len(denomSymbol) < MinDenomLen || len(denomSymbol) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom symbol %s, only accepts value [%d, %d]", denomSymbol, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomSymbol) || !IsAlpha(denomSymbol) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom symbol %s, only accepts alphabetic characters", denomSymbol)
	}
	return nil
}
