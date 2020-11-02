package money

import "testing"

func TestIsKnownCurrency(t *testing.T) {
	tests := []struct {
		in  Currency
		out bool
	}{
		{
			in:  "USD",
			out: true,
		},
		{
			in:  "EUR",
			out: true,
		},
		{
			in:  "RUB",
			out: true,
		},
	}
	for _, test := range tests {
		res := IsKnownCurrency(test.in)
		if res != test.out {
			t.Errorf("TestIsKnownCurrency got %t, want %t", res, test.out)
		}
	}
}
