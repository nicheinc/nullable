package nullable_test

import (
	"encoding/json"
	"fmt"

	"github.com/nicheinc/nullable"
)

func Example() {
	type Update struct {
		ID   int             `json:"-"`
		Name nullable.String `json:"name"`
		Flag nullable.Bool   `json:"flag"`
	}

	// Fields are unset by default.
	out := Update{
		ID:   1,
		Name: nullable.NewString("Alice"),
	}
	if data, err := nullable.MarshalJSON(&out); err == nil {
		fmt.Println("Flag unset:", string(data))
	}

	// Fields can be explicitly nulled to remove them.
	out.Flag.SetPtr(nil)
	if data, err := nullable.MarshalJSON(&out); err == nil {
		fmt.Println("Flag removed:", string(data))
	}

	// Unmarshalling JSON sets nullable fields appropriately.
	in := Update{}
	if err := json.Unmarshal([]byte(`{"flag":true}`), &in); err == nil {
		fmt.Println("Name is set:", in.Name.IsSet())
		fmt.Println("Flag value:", *in.Flag.Value())
	}

	// Output:
	// Flag unset: {"name":"Alice"}
	// Flag removed: {"name":"Alice","flag":null}
	// Name is set: false
	// Flag value: true
}
