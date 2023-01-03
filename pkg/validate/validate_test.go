package validate_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wgo-admin/backend/pkg/validate"
)

type Req struct {
	Method string `validate:"oneof=GET POST PUT PATCH DELETE"`
}

func TestValidate(t *testing.T) {
	should := assert.New(t)

	req := &Req{Method: "GET"}

	err := validate.ValidateStruct(req)
	if !should.NoError(err) {
		fmt.Println("err:", err)
	}
	fmt.Println("success")
}

func init() {
	validate.Init()
}
