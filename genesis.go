package onft

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/OmniFlix/onft/keeper"
	"github.com/OmniFlix/onft/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	for _, c := range data.Collections {
		if err := k.SetDenom(ctx, c.Denom); err != nil {
			panic(err)
		}
		if err := k.SetCollection(ctx, c); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(k.GetCollections(ctx))
}

func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState([]types.Collection{})
}

func ValidateGenesis(data types.GenesisState) error {
	for _, c := range data.Collections {
		if err := types.ValidateDenomID(c.Denom.Id); err != nil {
			return err
		}
		if !utf8.ValidString(c.Denom.Name) {
			return sdkerrors.Wrap(types.ErrInvalidDenom, "denom name is invalid")
		}

		for _, onft := range c.ONFTs {
			if onft.GetOwner().Empty() {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner")
			}

			if err := types.ValidateONFTID(onft.GetID()); err != nil {
				return err
			}

			if err := types.ValidateMediaURI(onft.GetMediaURI()); err != nil {
				return err
			}
		}
	}
	return nil
}