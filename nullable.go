package nullable

// Nullable is the interface for all nullable field update types.
type Nullable interface {
	// IsSet returns true if the receiver is a removal or modification update.
	IsSet() bool
	// Removed returns whether the receiver is a removal update.
	Removed() bool
	// UnmarshalJSON implements json.Unmarshaler.
	UnmarshalJSON(data []byte) error
	// InterfaceValue returns the (possibly nil) updated value as an interface{}
	// and is used internally for marshalling structs containing Nullables.
	InterfaceValue() interface{}
}
