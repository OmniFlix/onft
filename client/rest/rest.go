package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

const (
	RestParamDenom   = "denom"
	RestParamONFTID = "id"
	RestParamOwner   = "owner"
)

type createDenomReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Sender sdk.AccAddress `json:"sender"`
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Schema  string         `json:"schema"`
}

type mintONFTReq struct {
	BaseReq      rest.BaseReq   `json:"base_req"`
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    sdk.AccAddress `json:"recipient"`
	Denom        string         `json:"denom"`
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	MediaURI     string         `json:"media_uri"`
	PreviewURI   string         `json:"preview_uri"`
	Type         string         `json:"type"`
	Transferable string         `json:"transferable"`
}

type editONFTReq struct {
	BaseReq      rest.BaseReq   `json:"base_req"`
	Sender       sdk.AccAddress `json:"sender"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	MediaURI     string         `json:"media_uri"`
	PreviewURI   string         `json:"preview_uri"`
	Type         string         `json:"type"`
	Transferable string         `json:"transferable"`
}

type transferONFTReq struct {
	BaseReq    rest.BaseReq   `json:"base_req"`
	Sender     sdk.AccAddress `json:"sender"`
	Recipient  string         `json:"recipient"`
}

type burnONFTReq struct {
	BaseReq rest.BaseReq    `json:"base_req"`
	sender   sdk.AccAddress `json:"sender"`
}