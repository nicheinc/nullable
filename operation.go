package nullable

// Operation represents the operation that an update performs. The OpNoop,
// OpRemove, and OpSet constants are the only valid values of this type.
type Operation byte

const (
	// OpNoop indicates that an update does nothing.
	OpNoop Operation = iota
	// OpRemove indicates that an update removes a field.
	OpRemove
	// OpSet indicates that an update sets a field's value.
	OpSet
)

func (o Operation) String() string {
	switch o {
	case OpNoop:
		return "no-op"
	case OpRemove:
		return "remove"
	default: // Set
		return "set"
	}
}
