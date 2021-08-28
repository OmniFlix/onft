package types

var (
	EventTypeCreateDenom = "create_denom"
	EventTypeMintNFT     = "mint_onft"
	EventTypeTransfer    = "transfer_onft"
	EventTypeEditNFT     = "edit_onft"
	EventTypeBurnNFT     = "burn_onft"

	AttributeValueCategory = ModuleName

	AttributeKeySender     = "sender"
	AttributeKeyRecipient  = "recipient"
	AttributeKeyOwner      = "owner"
	AttributeKeyONFTID     = "onft-id"
	AttributeKeyMediaURI   = "media-uri"
	AttributeKeyPreviewURI = "preview-uri"
	AttributeKeyDenomID      = "denom-id"
)
