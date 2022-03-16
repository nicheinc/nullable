package nullable

// Operation represents the operation that an update performs. The Noop, Remove,
// and Set constants are the only valid values of this type.
type Operation byte

const (
	// Noop indicates that an update does nothing.
	Noop Operation = iota
	// Remove indicates that an update removes a field.
	Remove
	// Set indicates that an update sets a field's value.
	Set
)

func (o Operation) String() string {
	switch o {
	case Noop:
		return "no-op"
	case Remove:
		return "remove"
	default: // Set
		return "set"
	}
}
