package simulation

import (
	"math/rand"
	"github.com/OmniFlix/onft/types"
	"github.com/OmniFlix/onft/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"strings"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateDenom    = "op_weight_msg_create_denom"
	OpWeightMsgMintONFT       = "op_weight_msg_mint_onft"
	OpWeightMsgEditONFT       = "op_weight_msg_edit_onft"
	OpWeightMsgTransferONFT   = "op_weight_msg_transfer_onft"
	OpWeightMsgBurnONFT       = "op_weight_msg_burn_onft"
	OpWeightMsgTransferDenom = "op_weight_msg_transfer_denom"
	OpWeightMsgUpdateDenom = "op_weight_msg_update_denom"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	var weightCreateDenom, weightMint, weightEdit, weightBurn, weightTransfer, weightUpdateDenom, weightTransferDenom int

	appParams.GetOrGenerate(
		cdc, OpWeightMsgCreateDenom, &weightCreateDenom, nil,
		func(_ *rand.Rand) {
			weightCreateDenom = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintONFT, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditONFT, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferONFT, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgBurnONFT, &weightBurn, nil,
		func(_ *rand.Rand) {
			weightBurn = 10
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferDenom, &weightTransferDenom, nil,
		func(_ *rand.Rand) {
			weightTransferDenom = 10
		},
	)
	appParams.GetOrGenerate(
		cdc, OpWeightMsgUpdateDenom, &weightUpdateDenom, nil,
		func(_ *rand.Rand) {
			weightUpdateDenom = 10
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightCreateDenom,
			SimulateMsgCreateDenom(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMsgMintONFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateMsgEditNFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateMsgTransferONFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightBurn,
			SimulateMsgBurnONFT(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransferDenom,
			SimulateMsgTransferDenom(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightUpdateDenom,
			SimulateMsgUpdateDenom(k, ak, bk),
		),
	}
}


// SimulateMsgCreateDenom simulates create denom msg
func SimulateMsgCreateDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {

		denomName := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		symbol := simtypes.RandStringOfLength(r, 5)
		description := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		previewURI := strings.ToLower(simtypes.RandStringOfLength(r, 10))
		sender, _ := simtypes.RandomAcc(r, accs)


		msg := types.NewMsgCreateDenom(
			symbol,
			denomName,
			"{}",
			description,
			previewURI,
			sender.Address.String(),
		)
		account := ak.GetAccount(ctx, sender.Address)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())
		if spendableCoins.Empty() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateDenom, "unable to create denom"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgCreateDenom,
			Context:         ctx,
			SimAccount:      sender,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgMintONFT simulates a mint onft transaction
func SimulateMsgMintONFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		return simtypes.NewOperationMsg(&types.MsgMintONFT{}, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferONFT simulates the transfer of an NFT
func SimulateMsgTransferONFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		return simtypes.NewOperationMsg(&types.MsgTransferONFT{}, true, "", nil), nil, nil
	}
}

// SimulateMsgEditNFT simulates an edit onft transaction
func SimulateMsgEditNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		return simtypes.NewOperationMsg(&types.MsgEditONFT{}, true, "", nil), nil, nil
	}
}



// SimulateMsgBurnONFT simulates a burn onft transaction
func SimulateMsgBurnONFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		return simtypes.NewOperationMsg(&types.MsgBurnONFT{}, true, "", nil), nil, nil
	}
}

// SimulateMsgTransferDenom simulates a transfer denom transaction
func SimulateMsgTransferDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
) {
		return simtypes.NewOperationMsg(&types.MsgTransferDenom{}, true, "", nil), nil, nil
	}
}

// SimulateMsgUpdateDenom simulates a update denom transaction
func SimulateMsgUpdateDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		return simtypes.NewOperationMsg(&types.MsgUpdateDenom{}, true, "", nil), nil, nil
	}
}

