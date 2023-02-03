package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/OmniFlix/onft/types"
)

// SaveDenom saves a denom
func (k Keeper) SaveDenom(
	ctx sdk.Context,
	id,
	symbol,
	name,
	schema string,
	creator sdk.AccAddress,
	description,
	previewUri string,
	uri,
	uriHash,
	data string,
) error {
	denomMetadata := &types.DenomMetadata{
		Creator:    creator.String(),
		Schema:     schema,
		PreviewUri: previewUri,
		Data:       data,
	}
	metadata, err := codectypes.NewAnyWithValue(denomMetadata)
	if err != nil {
		return err
	}
	return k.nk.SaveClass(ctx, nft.Class{
		Id:          id,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        metadata,
	})
}

// TransferDenomOwner transfers the ownership to new address
func (k Keeper) TransferDenomOwner(
	ctx sdk.Context,
	denomID string,
	srcOwner,
	dstOwner sdk.AccAddress,
) error {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return err
	}

	// authorize
	if srcOwner.String() != denom.Creator {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%s is not allowed to transfer denom %s", srcOwner.String(),
			denomID,
		)
	}

	denomMetadata := &types.DenomMetadata{
		Creator:    dstOwner.String(),
		Schema:     denom.Schema,
		PreviewUri: denom.PreviewURI,
		Data:       denom.Data,
	}
	data, err := codectypes.NewAnyWithValue(denomMetadata)
	if err != nil {
		return err
	}
	class := nft.Class{
		Id:          denom.Id,
		Name:        denom.Name,
		Symbol:      denom.Symbol,
		Description: denom.Description,
		Uri:         denom.Uri,
		UriHash:     denom.UriHash,
		Data:        data,
	}

	return k.nk.UpdateClass(ctx, class)
}

func (k Keeper) HasDenom(ctx sdk.Context, id string) bool {
	return k.nk.HasClass(ctx, id)
}

func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom, err error) {
	classes := k.nk.GetClasses(ctx)
	for _, class := range classes {
		var denomMetadata types.DenomMetadata
		if err := k.cdc.Unmarshal(class.Data.GetValue(), &denomMetadata); err != nil {
			return nil, err
		}
		denoms = append(denoms, types.Denom{
			Id:          class.Id,
			Name:        class.Name,
			Schema:      denomMetadata.Schema,
			Creator:     denomMetadata.Creator,
			Symbol:      class.Symbol,
			Description: class.Description,
			PreviewURI:  denomMetadata.PreviewUri,
			Uri:         class.Uri,
			UriHash:     class.UriHash,
		})
	}
	return denoms, nil
}

func (k Keeper) AuthorizeDenomCreator(ctx sdk.Context, id string, creator sdk.AccAddress) error {
	denom, err := k.GetDenomInfo(ctx, id)
	if err != nil {
		return err
	}

	if creator.String() != denom.Creator {
		return sdkerrors.Wrap(types.ErrUnauthorized, creator.String())
	}
	return nil
}

func (k Keeper) HasPermissionToMint(ctx sdk.Context, denomID string, sender sdk.AccAddress) bool {
	denom, err := k.GetDenomInfo(ctx, denomID)
	if err != nil {
		return false
	}

	if sender.String() == denom.Creator {
		return true
	}
	return false
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denomID string) (*types.Denom, error) {
	class, ok := k.nk.GetClass(ctx, denomID)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", denomID)
	}

	var denomMetadata types.DenomMetadata
	if err := k.cdc.Unmarshal(class.Data.GetValue(), &denomMetadata); err != nil {
		return nil, err
	}
	return &types.Denom{
		Id:          class.Id,
		Name:        class.Name,
		Schema:      denomMetadata.Schema,
		Creator:     denomMetadata.Creator,
		Symbol:      class.Symbol,
		Description: class.Description,
		PreviewURI:  class.Uri,
	}, nil
}
