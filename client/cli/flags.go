package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagName         = "name"
	FlagDescription  = "description"
	FlagMediaURI     = "media-uri"
	FlagPreviewURI   = "preview-uri"
	FlagData         = "data"
	FlagTransferable = "transferable"
	FlagExtensible   = "extensible"
	FlagRecipient    = "recipient"
	FlagOwner        = "owner"
	FlagDenomID      = "denom-id"
	FlagSchema       = "schema"
)

var (
	FsCreateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateDenom   = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferDenom = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintONFT      = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditONFT      = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferONFT  = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsCreateDenom.String(FlagSchema, "", "Denom schema")
	FsCreateDenom.String(FlagName, "", "Name of the denom")
	FsCreateDenom.String(FlagDescription, "", "Description for denom")
	FsCreateDenom.String(FlagPreviewURI, "", "Preview image uri for denom")

	FsUpdateDenom.String(FlagName, "[do-not-modify]", "Name of the denom")
	FsUpdateDenom.String(FlagDescription, "[do-not-modify]", "Description for denom")
	FsUpdateDenom.String(FlagPreviewURI, "[do-not-modify]", "Preview image uri for denom")

	FsTransferDenom.String(FlagRecipient, "", "recipient of the denom")

	FsMintONFT.String(FlagMediaURI, "", "Media uri of onft")
	FsMintONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsMintONFT.String(FlagPreviewURI, "", "Preview uri of onft")
	FsMintONFT.String(FlagName, "", "Name of onft")
	FsMintONFT.String(FlagDescription, "", "Description of onft")
	FsMintONFT.String(FlagData, "", "custom data of onft")

	FsMintONFT.String(FlagTransferable, "yes", "transferability of onft (yes | no)")
	FsMintONFT.String(FlagExtensible, "yes", "extensisbility of onft (yes | no)")

	FsEditONFT.String(FlagMediaURI, "[do-not-modify]", "Media uri of onft")
	FsEditONFT.String(FlagPreviewURI, "[do-not-modify]", "Preview uri of onft")
	FsEditONFT.String(FlagName, "[do-not-modify]", "Name of nft")
	FsEditONFT.String(FlagDescription, "[do-not-modify]", "Description of onft")
	FsEditONFT.String(FlagTransferable, "[do-not-modify]", "transferability of onft")
	FsEditONFT.String(FlagData, "[do-not-modify]", "custom data of onft")
	FsEditONFT.String(FlagExtensible, "[do-not-modify]", "extensibility of onft (yes | no)")

	FsTransferONFT.String(FlagRecipient, "", "Receiver of the onft. default value is sender address of transaction")
	FsQuerySupply.String(FlagOwner, "", "The owner of a nft")
	FsQueryOwner.String(FlagDenomID, "", "id of the denom")
}
