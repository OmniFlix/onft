package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/OmniFlix/onft/types"
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
$ %s tx onft create [denomSymbol] --name=<name> --schema=<schema> --chain-id=<chain-id> --from=<key-name> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDenom(args[0],
				viper.GetString(FlagDenomName),
				viper.GetString(FlagSchema),
				clientCtx.GetFromAddress(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateDenom)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "mint [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an oNFT.
Example:
$ %s tx onft mint [denomID] --type <onft-type> --name <onft-name> --description <onft-descritpion> --media-uri=<uri> --preview-uri=<uri> 
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

			var recipient = clientCtx.GetFromAddress()

			recipientStr := strings.TrimSpace(viper.GetString(FlagRecipient))
			if len(recipientStr) > 0 {
				recipient, err = sdk.AccAddressFromBech32(recipientStr)
				if err != nil {
					return err
				}
			}

			onftMetadata := &types.Metadata{}
			onftName := viper.GetString(FlagONFTName)
			onftDescription := viper.GetString(FlagONFTDescription)
			onftMediaURI := viper.GetString(FlagONFTMediaURI)
			onftPreviewURI := viper.GetString(FlagONFTPreviewURI)
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
			switch strings.ToLower(viper.GetString(FlagONFTType)) {
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
			transferability := strings.ToLower(viper.GetString(FlagTransferable))
			if len(transferability) > 0 && (transferability == "no" || transferability == "false") {
				transferable = false
			}
			msg := types.NewMsgMintONFT(
				args[0],
				onftMetadata,
				onftType,
				transferable,
				clientCtx.GetFromAddress(),
				recipient,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintONFT)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdEditONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "edit [denomID] [onftID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the data of an oNFT.
Example:
$ %s tx onft edit [denomID] [onftID] --name=<onft-name> --description=<onft-description> --media-uri=<uri>
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

			onftMetadata := &types.Metadata{}
			onftName := viper.GetString(FlagONFTName)
			onftDescription := viper.GetString(FlagONFTDescription)
			onftMediaURI := viper.GetString(FlagONFTMediaURI)
			onftPreviewURI := viper.GetString(FlagONFTPreviewURI)
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
			onftType := strings.ToLower(viper.GetString(FlagONFTType))
			if !(len(onftType) > 0 && (onftType == "artwork" || onftType == "audio" || onftType == "video" ||
				onftType == types.DoNotModify)) {
				return fmt.Errorf("invalid option for type flag , valid options are artwork|audio|video")
			}
			transferable := strings.ToLower(viper.GetString(FlagTransferable))
			if !(len(transferable) > 0 && (transferable == "no" || transferable == "yes" ||
				transferable == types.DoNotModify)) {
				return fmt.Errorf("invalid option for transferable flag , valid options are yes|no")
			}
			msg := types.NewMsgEditONFT(
				args[1],
				args[0],
				onftMetadata,
				onftType,
				transferable,
				clientCtx.GetFromAddress(),
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
		Use: "transfer [recipient] [denomID] [onftID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer an oNFT to a recipient.
Example:
$ %s tx onft transfer [recipient] [denomID] [onftID] --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
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

			msg := types.NewMsgTransferONFT(
				args[2],
				args[1],
				clientCtx.GetFromAddress(),
				recipient,
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
		Use: "burn [denomID] [onftID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn an oNFT.
Example:
$ %s tx onft burn [denomID] [onftID] --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnONFT(clientCtx.GetFromAddress(), args[1], args[0])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
