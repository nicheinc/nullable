package nullable_test

import (
	"encoding/json"
	"fmt"

	"github.com/nicheinc/nullable/v2"
)

func Example() {
	type Update struct {
		ID   int                     `json:"-"`
		Name nullable.Update[string] `json:"name"`
		Flag nullable.Update[bool]   `json:"flag"`
	}

	// Update fields are no-ops by default and excluded from JSON.
	out := Update{
		ID:   1,
		Name: nullable.NewSet("Alice"),
	}
	if data, err := nullable.MarshalJSON(&out); err == nil {
		fmt.Println("With Flag unset:", string(data))
	}

	// Removal operations are marshalled with explicit null values.
	out.Flag = nullable.NewRemove[bool]()
	if data, err := nullable.MarshalJSON(&out); err == nil {
		fmt.Println("With Flag removed:", string(data))
	}

	// Unmarshalling from JSON sets nullable update fields appropriately.
	in := Update{}
	if err := json.Unmarshal([]byte(`{"flag":true}`), &in); err == nil {
		fmt.Println("Name is a", in.Name.Operation())
		if value, ok := in.Flag.Value(); ok {
			fmt.Println("Flag is set to", value)
		}
	}

	// Output:
	// With Flag unset: {"name":"Alice"}
	// With Flag removed: {"name":"Alice","flag":null}
	// Name is a no-op
	// Flag is set to true
}
