package util

import "net/url"

func IsURL(link string) error {
	_, err := url.ParseRequestURI(link)

	return err
}
