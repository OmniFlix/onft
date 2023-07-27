package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

const (
	ModuleName = "onft"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	PrefixONFT        = []byte{0x01}
	PrefixOwners      = []byte{0x02}
	PrefixCollection  = []byte{0x03}
	PrefixDenom       = []byte{0x04}
	PrefixDenomSymbol = []byte{0x05}
	PrefixCreator     = []byte{0x06}

	delimiter = []byte("/")
)

func SplitKeyOwner(key []byte) (address sdk.AccAddress, denom, id string, err error) {
	key = key[len(PrefixOwners)+len(delimiter):]
	keys := bytes.Split(key, delimiter)
	if len(keys) != 3 {
		return address, denom, id, errors.New("wrong KeyOwner")
	}

	address, _ = sdk.AccAddressFromBech32(string(keys[0]))
	denom = string(keys[1])
	id = string(keys[2])
	return
}

func SplitKeyDenom(key []byte) (denomID, tokenID string, err error) {
	keys := bytes.Split(key, delimiter)
	if len(keys) != 2 {
		return denomID, tokenID, errors.New("wrong KeyOwner")
	}

	denomID = string(keys[0])
	tokenID = string(keys[1])
	return
}

func KeyOwner(address sdk.AccAddress, denomID, onftID string) []byte {
	key := append(PrefixOwners, delimiter...)
	if address != nil {
		key = append(key, []byte(address.String())...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if address != nil && len(denomID) > 0 && len(onftID) > 0 {
		key = append(key, []byte(onftID)...)
	}
	return key
}

func KeyONFT(denomID, onftID string) []byte {
	key := append(PrefixONFT, delimiter...)
	if len(denomID) > 0 {
		key = append(key, []byte(denomID)...)
		key = append(key, delimiter...)
	}

	if len(denomID) > 0 && len(onftID) > 0 {
		key = append(key, []byte(onftID)...)
	}
	return key
}

func KeyCollection(denomID string) []byte {
	key := append(PrefixCollection, delimiter...)
	return append(key, []byte(denomID)...)
}

func KeyDenomID(id string) []byte {
	key := append(PrefixDenom, delimiter...)
	return append(key, []byte(id)...)
}

func KeyDenomCreator(address sdk.AccAddress, denomId string) []byte {
	key := append(PrefixCreator, delimiter...)
	if address != nil {
		key = append(key, []byte(address)...)
		key = append(key, delimiter...)
	}
	if address != nil && len(denomId) > 0 {
		key = append(key, []byte(denomId)...)
		key = append(key, delimiter...)
	}
	return key
}

func KeyDenomSymbol(symbol string) []byte {
	key := append(PrefixDenomSymbol, delimiter...)
	return append(key, []byte(symbol)...)
}

func MustMarshalSupply(cdc codec.BinaryCodec, supply uint64) []byte {
	supplyWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&supplyWrap)
}

func MustUnMarshalSupply(cdc codec.BinaryCodec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

func MustMarshalONFTID(cdc codec.BinaryCodec, onftID string) []byte {
	onftIDWrap := gogotypes.StringValue{Value: onftID}
	return cdc.MustMarshal(&onftIDWrap)
}

func MustUnMarshalONFTID(cdc codec.BinaryCodec, value []byte) string {
	var onftIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &onftIDWrap)
	return onftIDWrap.Value
}

func MustMarshalDenomID(cdc codec.BinaryCodec, denomID string) []byte {
	denomIDWrap := gogotypes.StringValue{Value: denomID}
	return cdc.MustMarshal(&denomIDWrap)
}

func MustUnMarshalDenomID(cdc codec.BinaryCodec, value []byte) string {
	var denomIDWrap gogotypes.StringValue
	cdc.MustUnmarshal(value, &denomIDWrap)
	return denomIDWrap.Value
}
