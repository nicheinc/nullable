package nullable

import "encoding/json"

// Nullable is the interface for all nullable field update types.
type Nullable interface {
	// IsSet returns true if the receiver is a removal or modification update.
	IsSet() bool
	// Removed returns whether the receiver is a removal update.
	Removed() bool
	// InterfaceValue returns the (possibly nil) updated value as an interface{}
	// and is used internally for marshalling structs containing Nullables.
	InterfaceValue() interface{}
	json.Unmarshaler
}
