package nup_test

import (
	"encoding/json"
	"fmt"

	nup "github.com/nicheinc/nullable/v2"
)

func Example() {
	type Update struct {
		ID   int                `json:"-"`
		Name nup.Update[string] `json:"name"`
		Flag nup.Update[bool]   `json:"flag"`
	}

	// Update fields are no-ops by default and excluded from JSON.
	out := Update{
		ID:   1,
		Name: nup.Set("Alice"),
	}
	if data, err := nup.MarshalJSON(&out); err == nil {
		fmt.Println("With Flag unset:", string(data))
	}

	// Removal operations are marshalled with explicit null values.
	out.Flag = nup.Remove[bool]()
	if data, err := nup.MarshalJSON(&out); err == nil {
		fmt.Println("With Flag removed:", string(data))
	}

	// Unmarshalling from JSON sets nup types appropriately.
	in := Update{}
	if err := json.Unmarshal([]byte(`{"flag":true}`), &in); err == nil {
		fmt.Println("Name is a", in.Name.Operation())
		if value, isSet := in.Flag.Value(); isSet {
			fmt.Println("Flag is set to", value)
		}
	}

	// Output:
	// With Flag unset: {"name":"Alice"}
	// With Flag removed: {"name":"Alice","flag":null}
	// Name is a no-op
	// Flag is set to true
}
