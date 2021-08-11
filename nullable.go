package nullable

// Nullable is the interface for all nullable field types.
type Nullable interface {
	IsSet() bool
	Removed() bool
	UnmarshalJSON(data []byte) error
	// InterfaceValue is used internally for marshalling structs of Nullables.
	InterfaceValue() interface{}
}
