package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagONFTName         = "name"
	FlagONFTDescription  = "description"
	FlagONFTMediaURI     = "media-uri"
	FlagONFTPreviewURI   = "preview-uri"
	FlagData             = "data"
	FlagTransferable     = "transferable"
	FlagExtensible       = "extensible"
	FlagRecipient        = "recipient"
	FlagOwner            = "owner"
	FlagDenomName        = "name"
	FlagDenomID          = "denom-id"
	FlagSchema           = "schema"
	FlagDenomDescription = "description"
	FlagDenomPreviewURI  = "preview-uri"
)

var (
	FsCreateDenom  = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintONFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditONFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferONFT = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply  = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner   = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateDenom.String(FlagSchema, "", "Denom schema")
	FsCreateDenom.String(FlagDenomName, "", "Name of the denom")
	FsCreateDenom.String(FlagDenomDescription, "", "Description for denom")
	FsCreateDenom.String(FlagDenomPreviewURI, "", "Preview image uri for denom")

	FsMintONFT.String(FlagONFTMediaURI, "", "Media uri of onft")
	FsMintONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsMintONFT.String(FlagONFTPreviewURI, "", "Preview uri of onft")
	FsMintONFT.String(FlagONFTName, "", "Name of onft")
	FsMintONFT.String(FlagONFTDescription, "", "Description of onft")
	FsMintONFT.String(FlagData, "", "custom data of onft")

	FsMintONFT.Bool(FlagTransferable, true, "transferability of onft (true | false)")
	FsMintONFT.Bool(FlagExtensible, false, "extensisbility of onft (true | false)")

	FsEditONFT.String(FlagONFTMediaURI, "[do-not-modify]", "Media uri of onft")
	FsEditONFT.String(FlagONFTPreviewURI, "[do-not-modify]", "Preview uri of onft")
	FsEditONFT.String(FlagONFTName, "[do-not-modify]", "Name of nft")
	FsEditONFT.String(FlagONFTDescription, "[do-not-modify]", "Description of onft")
	FsEditONFT.String(FlagTransferable, "[do-not-modify]", "transferability of onft")
	FsEditONFT.String(FlagData, "[do-not-modify]", "custom data of onft")

	FsTransferONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")
}
