package test

func StrToPtr(s string) *string {
	return &s
}

func IntToPtr(i int) *int {
	return &i
}

func FloatToPtr(f float64) *float64 {
	return &f
}

func BoolToPtr(b bool) *bool {
	return &b
}
