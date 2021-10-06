package onft

import (
	"github.com/OmniFlix/onft/keeper"
	"github.com/OmniFlix/onft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
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