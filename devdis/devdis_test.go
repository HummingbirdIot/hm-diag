package devdis

import (
	"testing"
)

func setup() func() {
	Init()
	go dis.Register()
	return func() {
	}
}

func TestRegister(t *testing.T) {
	teardown := setup()
	defer teardown()

}
