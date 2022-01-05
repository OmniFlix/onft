package simulation

import (
	"math/rand"

	"github.com/OmniFlix/onft/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"strings"
)


// genDenoms returns a slice of nft denoms.
func genDenoms(r *rand.Rand, accounts []simtypes.Account) []types.Denom {
	denoms := make([]types.Denom, len(accounts)-1)
	for i := 0; i < len(accounts)-1; i++ {
		denoms[i] = types.Denom{
			Id:          RandID(r, "onftdenom",  10),
			Name:        simtypes.RandStringOfLength(r, 10),
			Symbol:      simtypes.RandStringOfLength(r, 10),
			Schema: "",
			Description: simtypes.RandStringOfLength(r, 10),
			PreviewURI:  simtypes.RandStringOfLength(r, 10),
			Creator: accounts[i].Address.String(),
		}
	}
	return denoms
}

// genONFTCollection returns a slice of collection.
func genONFTCollection(r *rand.Rand, denoms []types.Denom, accounts []simtypes.Account) []types.Collection {
	collections := make([]types.Collection, len(denoms)-1)
	for i := 0; i < len(denoms)-1; i++ {
		onfts := make([]types.ONFT, len(accounts)/4)
		for j := 0; j < len(accounts)/4; j++ {
			owner := accounts[j]
			onfts[j] = types.ONFT{
				Id:           RandID(r, "onft", 10),
				Metadata:     RandMetadata(r),
				Data:         "",
				Owner:        owner.Address.String(),
				Transferable: true,
				Extensible:   true,
				CreatedAt:    simtypes.RandTimestamp(r),
			}
		}
		collections[i] = types.Collection{
			Denom: denoms[i],
			ONFTs: onfts,
		}
	}
	return collections
}

// RandomizedGenState generates a random GenesisState for onft
func RandomizedGenState(simState *module.SimulationState) {
	var collections []types.Collection
	denoms := genDenoms(simState.Rand, simState.Accounts)
	simState.AppParams.GetOrGenerate(
		simState.Cdc, "onft", &collections, simState.Rand,
		func(r *rand.Rand) { collections = genONFTCollection(r, denoms, simState.Accounts)},
	)
	onftGenesis := &types.GenesisState{
		Collections: collections,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(onftGenesis)
}

func RandID(r *rand.Rand, prefix string, n int) string {
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(prefix+id)
}

func RandMetadata(r *rand.Rand) types.Metadata {
	return types.Metadata{
		Name: simtypes.RandStringOfLength(r, 10),
		Description: simtypes.RandStringOfLength(r, 45),
		PreviewURI: simtypes.RandStringOfLength(r, 45),
		MediaURI: simtypes.RandStringOfLength(r, 45),
	}
}
