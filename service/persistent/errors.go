package persistent

import "errors"

var errAccountDoesNotExist = errors.New("account does not exist")
var errIncompatibleCurrency = errors.New("incompatible currency")
var errInsufficientFunds = errors.New("insufficient funds")
