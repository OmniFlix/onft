package helpers

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "onft-app"
)

type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
