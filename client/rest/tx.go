package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/OmniFlix/onft/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/onft/denoms",
		createDenomHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/onft/denoms/{%s}", RestParamDenom),
		updateDenomHandlerFn(cliCtx),
	).Methods("PUT")

	r.HandleFunc(
		fmt.Sprintf("/onft/denoms/{%s}/transfer", RestParamDenom),
		transferDenomHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/onft/onfts/mint"),
		mintONFTHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/onft/onfts/{%s}/{%s}", RestParamDenom, RestParamONFTID),
		editONFTHandlerFn(cliCtx),
	).Methods("PUT")

	r.HandleFunc(
		fmt.Sprintf("/onft/onfts/{%s}/{%s}/transfer", RestParamDenom, RestParamONFTID),
		transferONFTHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/onft/onfts/{%s}/{%s}/burn", RestParamDenom, RestParamONFTID),
		burnONFTHandlerFn(cliCtx),
	).Methods("POST")
}

func createDenomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createDenomReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgCreateDenom(
			req.Symbol,
			req.Name,
			req.Schema,
			req.Description,
			req.PreviewURI,
			req.Sender.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func updateDenomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateDenomReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		vars := mux.Vars(r)

		msg := types.NewMsgUpdateDenom(
			vars[RestParamDenom],
			req.Name,
			req.Description,
			req.PreviewURI,
			req.Sender.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferDenomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferDenomReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		msg := types.NewMsgTransferDenom(
			vars[RestParamDenom],
			req.Sender.String(),
			recipient.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func mintONFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintONFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if req.Recipient.Empty() {
			req.Recipient = req.Sender
		}
		metadata := types.Metadata{}
		if len(req.Name) > 0 {
			metadata.Name = req.Name
		}
		if len(req.Description) > 0 {
			metadata.Description = req.Description
		}
		if len(req.MediaURI) > 0 {
			metadata.MediaURI = req.MediaURI
		}
		if len(req.PreviewURI) > 0 {
			metadata.PreviewURI = req.PreviewURI
		}
		transferable := true
		transferability := strings.ToLower(req.Transferable)
		if len(transferability) > 0 && (transferability == "no" || transferability == "false") {
			transferable = false
		}
		extensible := true
		extensibility := strings.ToLower(req.Extensible)
		if len(extensibility) > 0 && (extensibility == "no" || extensibility == "false") {
			extensible = false
		}
		nsfw := false
		nsfwStr := strings.ToLower(req.Nsfw)
		if len(nsfwStr) > 0 && (nsfwStr == "yes" || nsfwStr == "true") {
			nsfw = true
		}

		msg := types.NewMsgMintONFT(
			req.Denom,
			req.Sender.String(),
			req.Recipient.String(),
			metadata,
			req.Data,
			transferable,
			extensible,
			nsfw,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editONFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editONFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)

		metadata := types.Metadata{}
		if len(req.Name) > 0 {
			metadata.Name = req.Name
		}
		if len(req.Description) > 0 {
			metadata.Description = req.Description
		}
		if len(req.MediaURI) > 0 {
			metadata.MediaURI = req.MediaURI
		}
		if len(req.PreviewURI) > 0 {
			metadata.PreviewURI = req.PreviewURI
		}

		transferable := strings.ToLower(req.Transferable)
		if len(transferable) > 0 && !(transferable == "no" || transferable == "yes" ||
			transferable == types.DoNotModify) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid option for transferable flag , valid options are yes,no")
			return
		}
		extensible := strings.ToLower(req.Extensible)
		if len(extensible) > 0 && !(extensible == "no" || extensible == "yes" ||
			extensible == types.DoNotModify) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid option for extensible flag , valid options are yes,no")
			return
		}
		nsfw := strings.ToLower(req.Nsfw)
		if len(nsfw) > 0 && !(nsfw == "no" || nsfw == "yes" ||
			extensible == types.DoNotModify) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid option for nsfw flag , valid options are yes,no")
			return
		}
		msg := types.NewMsgEditONFT(
			vars[RestParamONFTID],
			vars[RestParamDenom],
			metadata,
			req.Data,
			transferable,
			extensible,
			nsfw,
			req.Sender.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferONFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferONFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		msg := types.NewMsgTransferONFT(
			vars[RestParamONFTID],
			vars[RestParamDenom],
			req.Sender.String(),
			recipient.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func burnONFTHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnONFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)

		// create the message
		msg := types.NewMsgBurnONFT(
			vars[RestParamDenom],
			vars[RestParamONFTID],
			req.Sender.String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
