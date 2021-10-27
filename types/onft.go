package types

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/exported"
)

var _ exported.ONFT = ONFT{}

func NewONFT(
	id string, metadata Metadata, assetType AssetType,
	transferable, extensible bool, owner sdk.AccAddress,
	createdTime time.Time) ONFT {
	return ONFT{
		Id:           id,
		Metadata:     metadata,
		Type:         assetType,
		Owner:        owner.String(),
		Transferable: transferable,
		Extensible:   extensible,
		CreatedAt:    createdTime,
	}
}

func (onft ONFT) GetID() string {
	return onft.Id
}

func (onft ONFT) GetName() string {
	return onft.Metadata.Name
}

func (onft ONFT) GetDescription() string {
	return onft.Metadata.Description
}

func (onft ONFT) GetMediaURI() string {
	return onft.Metadata.Media
}

func (onft ONFT) GetPreviewURI() string {
	return onft.Metadata.Preview
}

func (onft ONFT) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(onft.Owner)
	return owner
}

func (onft ONFT) GetMetadata() string {
	return onft.Metadata.String()
}

func (onft ONFT) GetType() string {
	return onft.Type.String()
}
func (onft ONFT) IsTransferable() bool {
	return onft.Transferable
}
func (onft ONFT) IsExtensible() bool {
	return onft.Extensible
}
func (onft ONFT) GetCreatedTime() time.Time {
	return onft.CreatedAt
}

// ONFT

type ONFTs []exported.ONFT

func NewONFTs(onfts ...exported.ONFT) ONFTs {
	if len(onfts) == 0 {
		return ONFTs{}
	}
	return onfts
}

func ValidateONFTID(onftID string) error {
	onftID = strings.TrimSpace(onftID)
	if len(onftID) < MinDenomLen || len(onftID) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidONFTID, "invalid onftID %s, only accepts value [%d, %d]", onftID, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(onftID) || !IsAlphaNumeric(onftID) {
		return sdkerrors.Wrapf(ErrInvalidONFTID, "invalid onftID %s, only accepts alphanumeric characters,and begin with an english letter", onftID)
	}
	return nil
}

func ValidateMediaURI(mediaURI string) error {
	if len(mediaURI) > MaxMediaURILen {
		return sdkerrors.Wrapf(ErrInvalidMediaURI, "invalid mediaURI %s, only accepts value [0, %d]", mediaURI, MaxMediaURILen)
	}
	return nil
}

func ValidatePreviewURI(previewURI string) error {
	if len(previewURI) > MaxPreviewURILen {
		return sdkerrors.Wrapf(ErrInvalidPreviewURI, "invalid previewURI %s, only accepts value [0, %d]", previewURI, MaxPreviewURILen)
	}
	return nil
}
