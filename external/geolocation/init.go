package geolocation

import (
	"assignment/external"
)

func init() {
	service := NewGeolocationService()
	external.RegisterService("geolocation", service)
}
