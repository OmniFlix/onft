package types

var (
	EventTypeCreateDenom  = "create_denom"
	EventTypeMintONFT     = "mint_onft"
	EventTypeTransferONFT = "transfer_onft"
	EventTypeEditONFT     = "edit_onft"
	EventTypeBurnONFT     = "burn_onft"

	AttributeValueCategory = ModuleName

	AttributeKeySender      = "sender"
	AttributeKeyRecipient   = "recipient"
	AttributeKeyCreator     = "creator"
	AttributeKeyOwner       = "owner"
	AttributeKeyDenomID     = "denom-id"
	AttributeKeyDenomSymbol = "denom-symbol"
	AttributeKeyDenomName   = "denom-name"
	AttributeKeyONFTID      = "onft-id"
	AttributeKeyMediaURI    = "media-uri"
	AttributeKeyPreviewURI  = "preview-uri"
)
