package cli

import (
	"fmt"
	"github.com/OmniFlix/onft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"
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
		GetCmdUpdateDenom(),
		GetCmdTransferDenom(),
		GetCmdMintONFT(),
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
$ %s tx onft create [symbol] --name=<name> --schema=<schema> --description=<description> --preview-uri=<preview-uri> 
--creation-fee <collection-creation-fee> --chain-id=<chain-id> --from=<key-name> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			symbol := args[0]

			denomName, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			schema, err := cmd.Flags().GetString(FlagSchema)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			previewURI, err := cmd.Flags().GetString(FlagPreviewURI)
			if err != nil {
				return err
			}
			creationFeeStr, err := cmd.Flags().GetString(FlagCreationFee)
			if err != nil {
				return err
			}
			creationFee, err := sdk.ParseCoinNormalized(creationFeeStr)
			if err != nil {
				return fmt.Errorf("failed to parse creation fee: %s", creationFeeStr)
			}

			msg := types.NewMsgCreateDenom(symbol,
				denomName,
				schema,
				description,
				previewURI,
				clientCtx.GetFromAddress().String(),
				creationFee,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCreateDenom)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagCreationFee)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "mint [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an oNFT.
Example:
$ %s tx onft mint [denom-id] \ 
	--name <onft-name> \
	--description <onft-descritpion> \
	--media-uri=<uri> \
	--preview-uri=<uri> \
	--from=<key-name> \
	--chain-id=<chain-id> \
	--fees=<fee>

Additional Flags
    --non-trasferable
    --inextensible
    --nsfw
    --royalty-share="0.05"
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			denomId := args[0]

			sender := clientCtx.GetFromAddress().String()

			recipient, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}

			if len(recipient) > 0 {
				if _, err = sdk.AccAddressFromBech32(recipient); err != nil {
					return err
				}
			} else {
				recipient = sender
			}

			onftMetadata := types.Metadata{}
			onftName, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			onftDescription, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			onftMediaURI, err := cmd.Flags().GetString(FlagMediaURI)
			if err != nil {
				return err
			}

			onftPreviewURI, err := cmd.Flags().GetString(FlagPreviewURI)
			if err != nil {
				return err
			}

			if len(onftName) > 0 {
				onftMetadata.Name = onftName
			}
			if len(onftDescription) > 0 {
				onftMetadata.Description = onftDescription
			}
			if len(onftMediaURI) > 0 {
				onftMetadata.MediaURI = onftMediaURI
			}
			if len(onftPreviewURI) > 0 {
				onftMetadata.PreviewURI = onftPreviewURI
			}
			data, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			transferable := true
			nonTransferable, err := cmd.Flags().GetBool(FlagNonTransferable)
			if err != nil {
				return err
			}
			if nonTransferable {
				transferable = false
			}
			extensible := true
			inExtensible, err := cmd.Flags().GetBool(FlagInExtensible)
			if err != nil {
				return err
			}
			if inExtensible {
				extensible = false
			}
			nsfw := false
			nsfwFlag, err := cmd.Flags().GetBool(FlagNsfw)
			if err != nil {
				return err
			}
			if nsfwFlag {
				nsfw = true
			}
			royaltyShareStr, err := cmd.Flags().GetString(FlagRoyaltyShare)
			if err != nil {
				return err
			}
			royaltyShare := sdk.NewDec(0)
			if len(royaltyShareStr) > 0 {
				royaltyShare, err = sdk.NewDecFromStr(royaltyShareStr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgMintONFT(
				denomId,
				sender,
				recipient,
				onftMetadata,
				data,
				transferable,
				extensible,
				nsfw,
				royaltyShare,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMintONFT)
	_ = cmd.MarkFlagRequired(FlagMediaURI)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdUpdateDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-denom [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the data of Denom.
Example:
$ %s tx onft update-denom [denom-id] --name=<onft-name> --description=<onft-description> 
--preview-uri=<uri> --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			denomId := strings.TrimSpace(args[0])

			denomName, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}

			denomDescription, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			denomPreviewURI, err := cmd.Flags().GetString(FlagPreviewURI)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateDenom(
				denomId,
				denomName,
				denomDescription,
				denomPreviewURI,
				clientCtx.GetFromAddress().String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsUpdateDenom)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdTransferDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer-denom [recipient] [denom-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer a denom to a recipient.
Example:
$ %s tx onft transfer-denom [recipient] [denom-id] --from=<key-name> --chain-id=<chain-id> --fees=<fee>`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			denomId := args[1]

			msg := types.NewMsgTransferDenom(
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
	cmd.Flags().AddFlagSet(FsTransferDenom)
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

			denomId := strings.ToLower(strings.TrimSpace(args[1]))
			onftId := strings.ToLower(strings.TrimSpace(args[2]))

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
			denomId := strings.ToLower(strings.TrimSpace(args[0]))
			onftId := strings.ToLower(strings.TrimSpace(args[1]))

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
