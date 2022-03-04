package error

import "errors"

const (
	ParserError       = "couldn't parser body"
	PathNotFoundError = "couldn't find this path"
	AddressNotFound   = "address not found"
	ResponseError     = "couldn't answer to request"
)

var ErrAddressNotFound = errors.New(AddressNotFound)
