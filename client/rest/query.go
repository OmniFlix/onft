package rest

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/OmniFlix/onft/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/onft/denoms/{%s}/supply", RestParamDenom),
		querySupply(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/onft/collections/{%s}", RestParamDenom),
		queryCollection(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		"/onft/denoms",
		queryDenoms(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/onft/denoms/{%s}", RestParamDenom),
		queryDenom(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/onft/asset/{%s}/{%s}", RestParamDenom, RestParamONFTID),
		queryONFT(cliCtx, queryRoute),
	).Methods("GET")
}

func querySupply(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := strings.TrimSpace(mux.Vars(r)[RestParamDenom])
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		var owner sdk.AccAddress
		ownerStr := r.FormValue(RestParamOwner)
		if len(ownerStr) > 0 {
			ownerAddress, err := sdk.AccAddressFromBech32(strings.TrimSpace(ownerStr))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = ownerAddress
		}
		params := types.NewQuerySupplyParams(denom, owner)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySupply), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		out := binary.LittleEndian.Uint64(res)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, out)
	}
}

func queryCollection(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := mux.Vars(r)[RestParamDenom]
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryCollectionParams(denom)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCollection), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDenom(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		denom := mux.Vars(r)[RestParamDenom]
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryDenomParams(denom)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDenom), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDenoms(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDenoms), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryONFT(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		denom := vars[RestParamDenom]
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		onftID := vars[RestParamONFTID]
		if err := types.ValidateONFTID(onftID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryONFTParams(denom, onftID)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryONFT), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
