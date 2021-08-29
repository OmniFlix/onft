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

		msg := types.NewMsgCreateDenom(req.ID, req.Name, req.Schema, req.Sender)
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
			metadata.Media = req.MediaURI
		}
		if len(req.PreviewURI) > 0 {
			metadata.Preview = req.PreviewURI
		}
		var onftType types.AssetType
		switch strings.ToLower(req.Type) {
		case "artwork":
			onftType = types.ARTWORK
		case "audio":
			onftType = types.AUDIO
		case "video":
			onftType = types.VIDEO
		default:
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid onft type, valid types are artwork,audio,video")
			return
		}
		transferable := true
		transferability := strings.ToLower(req.Transferable)
		if len(transferability) > 0 && (transferability == "no" || transferability == "false") {
			transferable = false
		}

		msg := types.NewMsgMintONFT(
			req.Denom,
			metadata,
			onftType,
			transferable,
			req.Sender,
			req.Recipient,
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
			metadata.Media = req.MediaURI
		}
		if len(req.PreviewURI) > 0 {
			metadata.Preview = req.PreviewURI
		}
		onftType := strings.ToLower(req.Type)
		if !(len(onftType) > 0 && (onftType == "artwork" || onftType == "audio" || onftType == "video" ||
			onftType == types.DoNotModify)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid onft type, valid types are artwork,audio,video")
			return
		}
		transferable := strings.ToLower(req.Transferable)
		if !(len(transferable) > 0 && (transferable == "no" || transferable == "yes" ||
			transferable == types.DoNotModify)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid option for transferable flag , valid options are yes,no")
			return
		}
		msg := types.NewMsgEditONFT(
			vars[RestParamONFTID],
			vars[RestParamDenom],
			metadata,
			onftType,
			transferable,
			req.Sender,
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
			req.Sender,
			recipient,
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
			req.sender,
			vars[RestParamONFTID],
			vars[RestParamDenom],
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
