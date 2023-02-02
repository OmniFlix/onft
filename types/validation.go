package types

import "strings"

func IsIBCDenom(denomID string) bool {
	return strings.HasPrefix(denomID, "ibc/")
}
