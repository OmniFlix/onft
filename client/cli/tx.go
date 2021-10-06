package cli

import (
	"fmt"
	"strings"

	"github.com/OmniFlix/onft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "oNFT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateDenom(),
		GetCmdMintONFT(),
		GetCmdEditONFT(),
		GetCmdTransferONFT(),
		GetCmdBurnONFT(),
	)

	return txCmd
}

func GetCmdCreateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create [symbol]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new denom.
Example:
$ %s tx onft create [symbol] --name=<name> --schema=<schema> --chain-id=<chain-id> --from=<key-name> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			symbol := strings.TrimSpace(args[0])
			denomName, err := cmd.Flags().GetString(FlagDenomName)
			if err != nil {
				return err
			}
			denomName = strings.TrimSpace(denomName)
			schema, err := cmd.Flags().GetString(FlagSchema)
			if err != nil {
				return err
			}
			schema = strings.TrimSpace(schema)
			msg := types.NewMsgCreateDenom(symbol,
				denomName,
				schema,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateDenom)
	_ = cmd.MarkFlagRequired(FlagDenomName)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "mint [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an oNFT.
Example:
$ %s tx onft mint [denom-id] --type <onft-type> --name <onft-name> --description <onft-descritpion> --media-uri=<uri> --preview-uri=<uri> 
--transferable <yes/no> --recipient=<recipient> --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			denomId := strings.ToLower(strings.TrimSpace(args[0]))

			sender := clientCtx.GetFromAddress().String()

			recipient, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}

			recipientStr := strings.TrimSpace(recipient)
			if len(recipientStr) > 0 {
				if _, err = sdk.AccAddressFromBech32(recipientStr); err != nil {
					return err
				}
			} else {
				recipient = sender
			}

			onftMetadata := types.Metadata{}
			onftName, err := cmd.Flags().GetString(FlagONFTName)
			if err != nil {
				return err
			}
			onftName = strings.TrimSpace(onftName)

			onftDescription, err := cmd.Flags().GetString(FlagONFTDescription)
			if err != nil {
				return err
			}
			onftDescription = strings.TrimSpace(onftDescription)

			onftMediaURI, err := cmd.Flags().GetString(FlagONFTMediaURI)
			if err != nil {
				return err
			}
			onftMediaURI = strings.TrimSpace(onftMediaURI)

			onftPreviewURI, err := cmd.Flags().GetString(FlagONFTPreviewURI)
			if err != nil {
				return err
			}
			onftPreviewURI = strings.TrimSpace(onftPreviewURI)

			if len(onftName) > 0 {
				onftMetadata.Name = onftName
			}
			if len(onftDescription) > 0 {
				onftMetadata.Description = onftDescription
			}
			if len(onftMediaURI) > 0 {
				onftMetadata.Media = onftMediaURI
			}
			if len(onftPreviewURI) > 0 {
				onftMetadata.Preview = onftPreviewURI
			}
			var onftType types.AssetType
			onftAssetType, err := cmd.Flags().GetString(FlagONFTType)
			if err != nil {
				return err
			}
			onftAssetType = strings.ToLower(strings.TrimSpace(onftAssetType))
			switch onftAssetType {
			case "artwork":
				onftType = types.ARTWORK
			case "audio":
				onftType = types.AUDIO
			case "video":
				onftType = types.VIDEO
			default:
				return fmt.Errorf("invalid onft type, valid types are artwork,audio,video")
			}
			transferable := true
			transferability, err := cmd.Flags().GetString(FlagONFTType)
			if err != nil {
				return err
			}
			transferability = strings.ToLower(strings.TrimSpace(transferability))
			if len(transferability) > 0 && (transferability == "no" || transferability == "false") {
				transferable = false
			}
			msg := types.NewMsgMintONFT(
				denomId,
				sender,
				recipient,
				onftMetadata,
				onftType,
				transferable,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintONFT)
	_ = cmd.MarkFlagRequired(FlagONFTMediaURI)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEditONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "edit [denom-id] [onft-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the data of an oNFT.
Example:
$ %s tx onft edit [denom-id] [onft-id] --name=<onft-name> --description=<onft-description> --media-uri=<uri>
--preview-uri=<uri> --type=<onft-type> --transferable=<yes|no> --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			denomId, err := cmd.Flags().GetString(args[0])
			if err != nil {
				return err
			}
			denomId = strings.ToLower(strings.TrimSpace(denomId))

			onftId, err := cmd.Flags().GetString(args[1])
			if err != nil {
				return err
			}
			onftId = strings.ToLower(strings.TrimSpace(onftId))

			onftMetadata := types.Metadata{}
			onftName, err := cmd.Flags().GetString(FlagONFTName)
			if err != nil {
				return err
			}
			onftName = strings.TrimSpace(onftName)

			onftDescription, err := cmd.Flags().GetString(FlagONFTDescription)
			if err != nil {
				return err
			}
			onftDescription = strings.TrimSpace(onftDescription)

			onftMediaURI, err := cmd.Flags().GetString(FlagONFTMediaURI)
			if err != nil {
				return err
			}
			onftMediaURI = strings.TrimSpace(onftMediaURI)

			onftPreviewURI, err := cmd.Flags().GetString(FlagONFTPreviewURI)
			if err != nil {
				return err
			}
			onftPreviewURI = strings.TrimSpace(onftPreviewURI)

			if len(onftName) > 0 {
				onftMetadata.Name = onftName
			}
			if len(onftDescription) > 0 {
				onftMetadata.Description = onftDescription
			}
			if len(onftMediaURI) > 0 {
				onftMetadata.Media = onftMediaURI
			}
			if len(onftPreviewURI) > 0 {
				onftMetadata.Preview = onftPreviewURI
			}
			onftType, err := cmd.Flags().GetString(FlagONFTType)
			if err != nil {
				return err
			}
			onftType = strings.ToLower(strings.TrimSpace(onftType))

			if len(onftType) > 0 && !(onftType == "artwork" || onftType == "audio" || onftType == "video" ||
				onftType == types.DoNotModify) {
				return fmt.Errorf("invalid option for type flag , valid options are artwork|audio|video")
			}
			transferable, err := cmd.Flags().GetString(FlagTransferable)
			if err != nil {
				return err
			}
			if !(len(transferable) > 0 && (transferable == "no" || transferable == "yes" ||
				transferable == types.DoNotModify)) {
				return fmt.Errorf("invalid option for transferable flag , valid options are yes|no")
			}
			msg := types.NewMsgEditONFT(
				onftId,
				denomId,
				onftMetadata,
				onftType,
				transferable,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsEditONFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdTransferONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer [recipient] [denom-id] [onft-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer an oNFT to a recipient.
Example:
$ %s tx onft transfer [recipient] [denom-id] [onft-id] --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			denomId, err := cmd.Flags().GetString(args[1])
			if err != nil {
				return err
			}
			denomId = strings.ToLower(strings.TrimSpace(denomId))

			onftId, err := cmd.Flags().GetString(args[2])
			if err != nil {
				return err
			}
			onftId = strings.ToLower(strings.TrimSpace(onftId))

			msg := types.NewMsgTransferONFT(
				onftId,
				denomId,
				clientCtx.GetFromAddress().String(),
				recipient.String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsTransferONFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdBurnONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "burn [denom-id] [onft-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn an oNFT.
Example:
$ %s tx onft burn [denom-id] [onft-id] --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			denomId, err := cmd.Flags().GetString(args[0])
			if err != nil {
				return err
			}
			denomId = strings.ToLower(strings.TrimSpace(denomId))

			onftId, err := cmd.Flags().GetString(args[1])
			if err != nil {
				return err
			}
			onftId = strings.ToLower(strings.TrimSpace(onftId))

			msg := types.NewMsgBurnONFT(denomId, onftId, clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
