package keeper

import (
	"fmt"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"

	"github.com/OmniFlix/onft/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
	nk       nftkeeper.Keeper
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
		nk:       nftkeeper.NewKeeper(storeKey, cdc, ak, bk),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("OmniFlix/%s", types.ModuleName))
}

// NFTkeeper returns a cosmos-sdk nftkeeper.Keeper.
func (k Keeper) NFTkeeper() nftkeeper.Keeper {
	return k.nk
}
