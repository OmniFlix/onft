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
		if err := k.SetCollection(ctx, c); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	state, err := k.GetCollections(ctx)
	if err != nil {
		panic(err.Error())
	}
	return types.NewGenesisState(state)
}

func DefaultGenesisState() *types.GenesisState {
	return types.NewGenesisState([]types.Collection{})
}
