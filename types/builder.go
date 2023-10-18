package types

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	proto "github.com/cosmos/gogoproto/proto"
)

const (
	Namespace          = "omniflix:"
	KeyMediaFieldValue = "value"
)

var (
	ClassKeyName         = fmt.Sprintf("%s%s", Namespace, "name")
	ClassKeySymbol       = fmt.Sprintf("%s%s", Namespace, "symbol")
	ClassKeyDescription  = fmt.Sprintf("%s%s", Namespace, "description")
	ClassKeyURIHash      = fmt.Sprintf("%s%s", Namespace, "uri_hash")
	ClassKeyCreator      = fmt.Sprintf("%s%s", Namespace, "creator")
	ClassKeySchema       = fmt.Sprintf("%s%s", Namespace, "schema")
	ClassKeyPreviewURI   = fmt.Sprintf("%s%s", Namespace, "preview_uri")
	TokenKeyName         = fmt.Sprintf("%s%s", Namespace, "name")
	TokenKeyURIHash      = fmt.Sprintf("%s%s", Namespace, "uri_hash")
	TokenKeyPreviewURI   = fmt.Sprintf("%s%s", Namespace, "preview_uri")
	TokenKeyTransferable = fmt.Sprintf("%s%s", Namespace, "transferable")
	TokenKeyExtensible   = fmt.Sprintf("%s%s", Namespace, "extensible")
	TokenKeyNSFW         = fmt.Sprintf("%s%s", Namespace, "nsfw")
	TokenKeyRoyaltyShare = fmt.Sprintf("%s%s", Namespace, "royalty_share")
)

type ClassBuilder struct {
	cdc              codec.Codec
	getModuleAddress func(string) sdk.AccAddress
}
type TokenBuilder struct {
	cdc codec.Codec
}
type MediaField struct {
	Value interface{} `json:"value"`
	Mime  string      `json:"mime,omitempty"`
}

func NewClassBuilder(
	cdc codec.Codec,
	getModuleAddress func(string) sdk.AccAddress,
) ClassBuilder {
	return ClassBuilder{
		cdc:              cdc,
		getModuleAddress: getModuleAddress,
	}
}

// BuildMetadata encode class into the metadata format defined by ics721
func (cb ClassBuilder) BuildMetadata(class nft.Class) (string, error) {
	var message proto.Message
	if err := cb.cdc.UnpackAny(class.Data, &message); err != nil {
		return "", err
	}

	metadata, ok := message.(*DenomMetadata)
	if !ok {
		return "", errors.New("unsupported classMetadata")
	}

	kvals := make(map[string]interface{})
	if len(metadata.Data) > 0 {
		err := json.Unmarshal([]byte(metadata.Data), &kvals)
		if err != nil && IsIBCDenom(class.Id) {
			// when classData is not a legal json, there is no need to parse the data
			return base64.RawStdEncoding.EncodeToString([]byte(metadata.Data)), nil
		}
		// note: if metadata.Data is null, it may cause map to be redefined as nil
		if kvals == nil {
			kvals = make(map[string]interface{})
		}
	}
	creator, err := sdk.AccAddressFromBech32(metadata.Creator)
	if err != nil {
		return "", err
	}

	hexCreator := hex.EncodeToString(creator)
	kvals[ClassKeyName] = MediaField{Value: class.Name}
	kvals[ClassKeySymbol] = MediaField{Value: class.Symbol}
	kvals[ClassKeyDescription] = MediaField{Value: class.Description}
	kvals[ClassKeyURIHash] = MediaField{Value: class.UriHash}
	kvals[ClassKeyCreator] = MediaField{Value: hexCreator}
	kvals[ClassKeySchema] = MediaField{Value: metadata.Schema}
	kvals[ClassKeyPreviewURI] = MediaField{Value: metadata.PreviewUri}
	data, err := json.Marshal(kvals)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

// Build create a class from ics721 packetData
func (cb ClassBuilder) Build(classID, classURI, classData string) (nft.Class, error) {
	classDataBz, err := base64.RawStdEncoding.DecodeString(classData)
	if err != nil {
		return nft.Class{}, err
	}

	var (
		name        = ""
		symbol      = ""
		description = ""
		uriHash     = ""
		schema      = ""
		previewURI  = ""
		creator     = cb.getModuleAddress(ModuleName).String()
	)

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(classDataBz, &dataMap); err != nil {
		any, err := codectypes.NewAnyWithValue(&DenomMetadata{
			Creator:     creator,
			Schema:      schema,
			Description: description,
			PreviewUri:  previewURI,
			Data:        string(classDataBz),
		})
		if err != nil {
			return nft.Class{}, err
		}
		return nft.Class{
			Id:          classID,
			Uri:         classURI,
			Name:        name,
			Symbol:      symbol,
			Description: description,
			UriHash:     uriHash,
			Data:        any,
		}, nil
	}
	if v, ok := dataMap[ClassKeyName]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				name = vStr
				delete(dataMap, ClassKeyName)
			}
		}
	}

	if v, ok := dataMap[ClassKeySymbol]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				symbol = vStr
				delete(dataMap, ClassKeySymbol)
			}
		}
	}

	if v, ok := dataMap[ClassKeyDescription]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				description = vStr
				delete(dataMap, ClassKeyDescription)
			}
		}
	}

	if v, ok := dataMap[ClassKeyURIHash]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				uriHash = vStr
				delete(dataMap, ClassKeyURIHash)
			}
		}
	}

	if v, ok := dataMap[ClassKeyCreator]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				creatorAcc, err := sdk.AccAddressFromHexUnsafe(vStr)
				if err != nil {
					return nft.Class{}, err
				}
				creator = creatorAcc.String()
				delete(dataMap, ClassKeyCreator)
			}
		}
	}

	if v, ok := dataMap[ClassKeySchema]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				schema = vStr
				delete(dataMap, ClassKeySchema)
			}
		}
	}

	if v, ok := dataMap[ClassKeyPreviewURI]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				schema = vStr
				delete(dataMap, ClassKeyPreviewURI)
			}
		}
	}

	data := ""
	if len(dataMap) > 0 {
		dataBz, err := json.Marshal(dataMap)
		if err != nil {
			return nft.Class{}, err
		}
		data = string(dataBz)
	}

	denomMeta, err := codectypes.NewAnyWithValue(&DenomMetadata{
		Creator: creator,
		Schema:  schema,
		Data:    data,
	})
	if err != nil {
		return nft.Class{}, err
	}

	return nft.Class{
		Id:          classID,
		Uri:         classURI,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		UriHash:     uriHash,
		Data:        denomMeta,
	}, nil
}

func NewTokenBuilder(cdc codec.Codec) TokenBuilder {
	return TokenBuilder{
		cdc: cdc,
	}
}

// BuildMetadata encode nft into the metadata format defined by ics721
func (tb TokenBuilder) BuildMetadata(token nft.NFT) (string, error) {
	var message proto.Message
	if err := tb.cdc.UnpackAny(token.Data, &message); err != nil {
		return "", err
	}

	nftMetadata, ok := message.(*ONFTMetadata)
	if !ok {
		return "", errors.New("unsupport classMetadata")
	}
	kvals := make(map[string]interface{})
	if len(nftMetadata.Data) > 0 {
		err := json.Unmarshal([]byte(nftMetadata.Data), &kvals)
		if err != nil && IsIBCDenom(token.ClassId) {
			// when nftMetadata is not a legal json, there is no need to parse the data
			return base64.RawStdEncoding.EncodeToString([]byte(nftMetadata.Data)), nil
		}
		// note: if nftMetadata.Data is null, it may cause map to be redefined as nil
		if kvals == nil {
			kvals = make(map[string]interface{})
		}
	}
	kvals[TokenKeyName] = MediaField{Value: nftMetadata.Name}
	kvals[TokenKeyURIHash] = MediaField{Value: token.UriHash}
	data, err := json.Marshal(kvals)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

// Build create a nft from ics721 packet data
func (tb TokenBuilder) Build(classId, tokenId, tokenURI, tokenData string) (nft.NFT, error) {
	tokenDataBz, err := base64.RawStdEncoding.DecodeString(tokenData)
	if err != nil {
		return nft.NFT{}, err
	}

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(tokenDataBz, &dataMap); err != nil {
		metadata, err := codectypes.NewAnyWithValue(&ONFTMetadata{
			Data: string(tokenDataBz),
		})
		if err != nil {
			return nft.NFT{}, err
		}

		return nft.NFT{
			ClassId: classId,
			Id:      tokenId,
			Uri:     tokenURI,
			Data:    metadata,
		}, nil
	}

	var (
		name    string
		uriHash string
	)
	if v, ok := dataMap[TokenKeyName]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				name = vStr
				delete(dataMap, TokenKeyName)
			}
		}
	}

	if v, ok := dataMap[TokenKeyURIHash]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if vStr, ok := vMap[KeyMediaFieldValue].(string); ok {
				uriHash = vStr
				delete(dataMap, TokenKeyURIHash)
			}
		}
	}

	data := ""
	if len(dataMap) > 0 {
		dataBz, err := json.Marshal(dataMap)
		if err != nil {
			return nft.NFT{}, err
		}
		data = string(dataBz)
	}

	metadata, err := codectypes.NewAnyWithValue(&ONFTMetadata{
		Name: name,
		Data: data,
	})
	if err != nil {
		return nft.NFT{}, err
	}

	return nft.NFT{
		ClassId: classId,
		Id:      tokenId,
		Uri:     tokenURI,
		UriHash: uriHash,
		Data:    metadata,
	}, nil
}
