package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/OmniFlix/onft/exported"
	"github.com/OmniFlix/onft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

func (k Keeper) SaveNFT(
	ctx sdk.Context,
	denomID,
	nftID,
	name,
	description,
	mediaURI,
	uriHash,
	previewURI,
	nftData string,
	createdAt time.Time,
	transferable,
	extensible,
	nsfw bool,
	royaltyShare sdk.Dec,
	receiver sdk.AccAddress,
) error {
	nftMetadata := &types.ONFTMetadata{
		Name:         name,
		Description:  description,
		PreviewURI:   previewURI,
		Data:         nftData,
		Transferable: transferable,
		Extensible:   extensible,
		Nsfw:         nsfw,
		CreatedAt:    createdAt,
		RoyaltyShare: royaltyShare,
	}
	data, err := codectypes.NewAnyWithValue(nftMetadata)
	if err != nil {
		return err
	}
	return k.nk.Mint(ctx, nft.NFT{
		ClassId: denomID,
		Id:      nftID,
		Uri:     mediaURI,
		UriHash: uriHash,
		Data:    data,
	}, receiver)
}

func (k Keeper) TransferOwnership(ctx sdk.Context, denomID, onftID string, srcOwner, dstOwner sdk.AccAddress) error {
	if !k.nk.HasClass(ctx, denomID) {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	onft, exist := k.nk.GetNFT(ctx, denomID, onftID)
	if !exist {
		return errorsmod.Wrapf(types.ErrInvalidONFT, "nft ID %s not exists", onftID)
	}

	err := k.Authorize(ctx, denomID, onftID, srcOwner)
	if err != nil {
		return err
	}
	onftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, onft.Data.GetValue())
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidONFTMetadata, "unable to parse nft metadata")
	}

	if !onftMetadata.Transferable {
		return errorsmod.Wrap(types.ErrNotTransferable, onft.GetId())
	}
	return k.nk.Transfer(ctx, denomID, onftID, dstOwner)
}

func (k Keeper) BurnONFT(
	ctx sdk.Context,
	denomID,
	onftID string,
	owner sdk.AccAddress,
) error {
	if !k.nk.HasClass(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	_, exist := k.nk.GetNFT(ctx, denomID, onftID)
	if !exist {
		return sdkerrors.Wrapf(types.ErrInvalidONFT, "nft ID %s not exists", onftID)
	}

	err := k.Authorize(ctx, denomID, onftID, owner)
	if err != nil {
		return err
	}

	return k.nk.Burn(ctx, denomID, onftID)
}

func (k Keeper) GetONFT(ctx sdk.Context, denomID, onftID string) (nft exported.ONFTI, err error) {
	if !k.nk.HasClass(ctx, denomID) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	onft, exist := k.nk.GetNFT(ctx, denomID, onftID)
	if !exist {
		return nil, sdkerrors.Wrapf(types.ErrInvalidONFT, "not found NFT: %s", onftID)
	}

	nftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, onft.Data.GetValue())
	if err != nil {
		return nil, err
	}

	owner := k.nk.GetOwner(ctx, denomID, onftID)
	metadata := types.Metadata{
		Name:        nftMetadata.Name,
		Description: nftMetadata.Description,
		MediaURI:    onft.Uri,
		PreviewURI:  nftMetadata.PreviewURI,
	}
	return types.ONFT{
		Id:           onft.Id,
		Metadata:     metadata,
		Data:         nftMetadata.Data,
		Owner:        owner.String(),
		Transferable: nftMetadata.Transferable,
		Extensible:   nftMetadata.Extensible,
		Nsfw:         nftMetadata.Extensible,
		CreatedAt:    nftMetadata.CreatedAt,
		RoyaltyShare: nftMetadata.RoyaltyShare,
	}, nil
}

func (k Keeper) GetONFTs(ctx sdk.Context, denomID string) (onfts []exported.ONFTI, err error) {
	nfts := k.nk.GetNFTsOfClass(ctx, denomID)
	for _, nft := range nfts {

		nftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, nft.Data.GetValue())
		if err != nil {
			return nil, err
		}

		owner := k.nk.GetOwner(ctx, denomID, nft.GetId())
		metadata := types.Metadata{
			Name:        nftMetadata.Name,
			Description: nftMetadata.Description,
			MediaURI:    nft.Uri,
			PreviewURI:  nftMetadata.PreviewURI,
		}
		onfts = append(onfts, types.ONFT{
			Id:           nft.GetId(),
			Metadata:     metadata,
			Data:         nftMetadata.Data,
			Owner:        owner.String(),
			Transferable: nftMetadata.Transferable,
			Extensible:   nftMetadata.Extensible,
			Nsfw:         nftMetadata.Extensible,
			CreatedAt:    nftMetadata.CreatedAt,
			RoyaltyShare: nftMetadata.RoyaltyShare,
		})
	}
	return onfts, nil
}

func (k Keeper) GetOwnerONFTs(ctx sdk.Context, denomID string, owner sdk.AccAddress) (onfts []exported.ONFTI, err error) {
	nfts := k.nk.GetNFTsOfClassByOwner(ctx, denomID, owner)
	for _, nft := range nfts {

		nftMetadata, err := types.UnmarshalNFTMetadata(k.cdc, nft.Data.GetValue())
		if err != nil {
			return nil, err
		}

		owner := k.nk.GetOwner(ctx, denomID, nft.GetId())
		metadata := types.Metadata{
			Name:        nftMetadata.Name,
			Description: nftMetadata.Description,
			MediaURI:    nft.Uri,
			PreviewURI:  nftMetadata.PreviewURI,
		}
		onfts = append(onfts, types.ONFT{
			Id:           nft.GetId(),
			Metadata:     metadata,
			Data:         nftMetadata.Data,
			Owner:        owner.String(),
			Transferable: nftMetadata.Transferable,
			Extensible:   nftMetadata.Extensible,
			Nsfw:         nftMetadata.Extensible,
			CreatedAt:    nftMetadata.CreatedAt,
			RoyaltyShare: nftMetadata.RoyaltyShare,
		})
	}
	return onfts, nil
}

func (k Keeper) Authorize(
	ctx sdk.Context,
	denomID,
	onftID string,
	owner sdk.AccAddress,
) error {
	if !owner.Equals(k.nk.GetOwner(ctx, denomID, onftID)) {
		return sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}
	return nil
}

func (k Keeper) HasONFT(ctx sdk.Context, denomID, onftID string) bool {
	return k.nk.HasNFT(ctx, denomID, onftID)
}
