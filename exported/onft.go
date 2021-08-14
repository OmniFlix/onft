package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type ONFT interface {
	GetID() string
	GetOwner() sdk.AccAddress
	GetName() string
	GetDescription() string
	GetMediaURI() string
	GetPreviewURI() string
	GetType() string
	IsTransferable() bool
	GetCreatedTime() time.Time
}

