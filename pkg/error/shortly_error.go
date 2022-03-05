package error

import "errors"

const (
	ParserError       = "couldn't parser body"
	PathNotFoundError = "couldn't find this path"
	AddressNotFound   = "address not found"
	ResponseError     = "couldn't answer to request"
	InvalidURLError   = "URL is not valid"
)

var ErrAddressNotFound = errors.New(AddressNotFound)

const (
	ParserErrCode = 13001
	InvalidParams = 13002
)
