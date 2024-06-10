package bearer

import (
	"assignment/external"
)

func init() {
	service := NewBearerService()
	external.RegisterService("basic-auth", service)
}
