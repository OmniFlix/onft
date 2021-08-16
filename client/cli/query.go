package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/OmniFlix/onft/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the oNFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryDenom(),
		GetCmdQueryDenoms(),
		GetCmdQueryCollection(),
		GetCmdQuerySupply(),
		GetCmdQueryONFT(),
	)

	return queryCmd
}

func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use: "supply [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`total supply of a collection of oNFTs.
Example:
$ %s query onft supply [denom]`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var owner sdk.AccAddress

			ownerStr := strings.TrimSpace(viper.GetString(FlagOwner))
			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				Denom: denom,
				Owner: owner,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use: "collection [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get all the oNFTs from a given collection
Example:
$ %s query onft collection <denom>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(context.Background(), &types.QueryCollectionRequest{
				Denom: denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Collection)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denoms",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all denominations of all collections of oNFTs
Example:
$ %s query onft denoms`, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denoms(context.Background(), &types.QueryDenomsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use: "denom [denomID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the denominations by the specified denom name
Example:
$ %s query onft denom <denom>`, version.AppName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Denom(context.Background(), &types.QueryDenomRequest{
				Denom: denom,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Denom)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryONFT() *cobra.Command {
	cmd := &cobra.Command{
		Use: "asset [denomID] [onftID]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a single oNFT from a collection
Example:
$ %s query onft onft <denom> <onftID>`, version.AppName)),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[0])
			if err := types.ValidateDenomID(denom); err != nil {
				return err
			}

			tokenID := strings.TrimSpace(args[1])
			if err := types.ValidateONFTID(tokenID); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.ONFT(context.Background(), &types.QueryONFTRequest{
				Denom: denom,
				Id:    tokenID,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.ONFT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}