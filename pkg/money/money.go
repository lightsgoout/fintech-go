package money

import "github.com/shopspring/decimal"

type Numeric struct {
	value decimal.Decimal
}

func (n Numeric) Sub(a Numeric) Numeric {
	return Numeric{
		value: n.value.Sub(a.value),
	}
}

func (n Numeric) Add(a Numeric) Numeric {
	return Numeric{
		value: n.value.Add(a.value),
	}
}

func (n Numeric) LessThan(a Numeric) bool {
	return n.value.LessThan(a.value)
}

func (n Numeric) String() string {
	return n.value.String()
}

func NewNumericFromInt64(value int64) Numeric {
	return Numeric{
		value: decimal.NewFromInt(value),
	}
}

func NewNumericFromStringMust(value string) Numeric {
	return Numeric{
		value: decimal.RequireFromString(value),
	}
}
