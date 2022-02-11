package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

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

		msg := types.NewMsgMintONFT(
			req.Denom,
			req.Sender.String(),
			req.Recipient.String(),
			metadata,
			req.Data,
			req.Transferable,
			req.Extensible,
			req.Nsfw,
			req.RoyaltyShare,
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
