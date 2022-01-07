package simulation

import (
	"fmt"
	"github.com/OmniFlix/onft/keeper"
	"github.com/OmniFlix/onft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"
	"strings"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateDenom   = "op_weight_msg_create_denom"
	OpWeightMsgMintONFT      = "op_weight_msg_mint_onft"
	OpWeightMsgEditONFT      = "op_weight_msg_edit_onft"
	OpWeightMsgTransferONFT  = "op_weight_msg_transfer_onft"
	OpWeightMsgBurnONFT      = "op_weight_msg_burn_onft"
	OpWeightMsgTransferDenom = "op_weight_msg_transfer_denom"
	OpWeightMsgUpdateDenom   = "op_weight_msg_update_denom"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	var weightCreateDenom, weightMint, weightEdit, weightTransfer, weightBurn, weightUpdateDenom, weightTransferDenom int

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
		denomId := RandID(r, "onftdenom", 10)
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
		msg.Id = denomId
		denom, _ := k.GetDenom(ctx, msg.Id)
		if denom.Size() != 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateDenom, "denom exist"), nil, nil
		}
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
		denom, err  := getRandomDenom(ctx, k, r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMintONFT, err.Error()), nil, err
		}
		randomRecipient, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintONFT(
			denom.Id,
			denom.Creator,
			randomRecipient.Address.String(),
			RandMetadata(r),
			"{}",
			true,
			true,
		)
		onftId := RandID(r, "onft", 10)
		msg.Id = onftId
		minter, _ := sdk.AccAddressFromBech32(denom.Creator)
		account := ak.GetAccount(ctx, minter)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		sender, found := simtypes.FindAccount(accs, minter)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMintONFT, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgMintONFT,
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

// SimulateMsgEditNFT simulates the edit of an NFT
func SimulateMsgEditNFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("account invalid")
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditONFT, err.Error()), nil, err
		}

		msg := types.NewMsgEditONFT(
			nftID,
			denom,
			RandMetadata(r),
			"{}",
			"[do-not-modify]",
			"[do-not-modify]",
			ownerAddr.String(),
		)

		account := ak.GetAccount(ctx, ownerAddr)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", ownerAddr)
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgEditONFT, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgEditONFT,
			Context:         ctx,
			SimAccount:      ownerAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgTransferONFT simulates the transfer of an nft
func SimulateMsgTransferONFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferONFT, err.Error()), nil, err
		}
		if ownerAddr.Empty() {
			err = fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}

		recipientAccount, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferONFT(
			nftID,
			denom,
			ownerAddr.String(),                // sender
			recipientAccount.Address.String(), // recipient
		)
		account := ak.GetAccount(ctx, ownerAddr)

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferONFT, err.Error()), nil, err
		}

		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgTransferONFT,
			Context:         ctx,
			SimAccount:      ownerAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}


// SimulateMsgBurnONFT simulates a burn onft transaction
func SimulateMsgBurnONFT(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		ownerAddr, denom, nftID := getRandomNFTFromOwner(ctx, k, r)
		if ownerAddr.Empty() {
			err = fmt.Errorf("invalid account")
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurnONFT, err.Error()), nil, err
		}

		msg := types.NewMsgBurnONFT(denom, nftID, ownerAddr.String())

		account := ak.GetAccount(ctx, ownerAddr)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		ownerAccount, found := simtypes.FindAccount(accs, ownerAddr)
		if !found {
			err = fmt.Errorf("account %s not found", msg.Sender)
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgBurnONFT, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgBurnONFT,
			Context:         ctx,
			SimAccount:      ownerAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgTransferDenom simulates a transfer denom transaction
func SimulateMsgTransferDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		denom, err  := getRandomDenom(ctx, k, r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, err.Error()), nil, err
		}
		account := ak.GetAccount(ctx, creator)
		ownerAccount, found := simtypes.FindAccount(accs, account.GetAddress())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferDenom, "creator not found"), nil, nil
		}

		recipient, _ := simtypes.RandomAcc(r, accs)
		msg := types.NewMsgTransferDenom(
			denom.Id,
			denom.Creator,
			recipient.Address.String(),
		)

		spendableCoins := bk.SpendableCoins(ctx, ownerAccount.Address)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgTransferDenom,
			Context:         ctx,
			SimAccount:      ownerAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgUpdateDenom simulates a update denom transaction
func SimulateMsgUpdateDenom(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (
		opMsg simtypes.OperationMsg, fOps []simtypes.FutureOperation, err error,
	) {
		denom, err  := getRandomDenom(ctx, k, r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDenom, err.Error()), nil, err
		}

		creator, err := sdk.AccAddressFromBech32(denom.Creator)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDenom, err.Error()), nil, err
		}
		account := ak.GetAccount(ctx, creator)
		ownerAccount, found := simtypes.FindAccount(accs, account.GetAddress())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDenom, "creator not found"), nil, nil
		}
		msg := types.NewMsgUpdateDenom(
			denom.Id,
			simtypes.RandStringOfLength(r, 10),
			simtypes.RandStringOfLength(r, 45),
			simtypes.RandStringOfLength(r, 45),
			ownerAccount.Address.String(),
		)

		spendableCoins := bk.SpendableCoins(ctx, ownerAccount.Address)


		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         types.TypeMsgUpdateDenom,
			Context:         ctx,
			SimAccount:      ownerAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func getRandomNFTFromOwner(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) (address sdk.AccAddress, denomID, tokenID string) {
	owners := k.GetOwners(ctx)

	ownersLen := len(owners)
	if ownersLen == 0 {
		return nil, "", ""
	}

	// get random owner
	i := r.Intn(ownersLen)
	owner := owners[i]

	idCollectionsLen := len(owner.IDCollections)
	if idCollectionsLen == 0 {
		return nil, "", ""
	}

	// get random collection from owner's balance
	i = r.Intn(idCollectionsLen)
	idCollection := owner.IDCollections[i] // nfts IDs
	denomID = idCollection.DenomId

	idsLen := len(idCollection.OnftIds)
	if idsLen == 0 {
		return nil, "", ""
	}

	// get random nft from collection
	i = r.Intn(idsLen)
	tokenID = idCollection.OnftIds[i]

	ownerAddress, _ := sdk.AccAddressFromBech32(owner.Address)
	return ownerAddress, denomID, tokenID
}


func getRandomDenom(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) (types.Denom, error) {
	var denoms = []string{denomId1, denomId2}
	i := r.Intn(len(denoms))
	denom, _ := k.GetDenom(ctx, denoms[i])
	if denom.Size() == 0 {
		return types.Denom{}, fmt.Errorf("no denoms created")
	}
	return denom, nil
}