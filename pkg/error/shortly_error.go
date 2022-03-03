package error

import "errors"

const (
	ParserError       = "couldn't parser body"
	PathNotFoundError = "couldn't find this path"
	AddressNotFound   = "address not found"
)

var AddressNotFoundErr = errors.New(AddressNotFound)
