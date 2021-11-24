package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

func NewGenesisState(collections []Collection) *GenesisState {
	return &GenesisState{
		Collections: collections,
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, c := range data.Collections {
		if err := ValidateDenomID(c.Denom.Id); err != nil {
			return err
		}

		for _, nft := range c.ONFTs {
			if nft.GetOwner().Empty() {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing onft owner")
			}

			if err := ValidateONFTID(nft.GetID()); err != nil {
				return err
			}

			if err := ValidateMediaURI(nft.GetMediaURI()); err != nil {
				return err
			}
		}
	}
	return nil
}
