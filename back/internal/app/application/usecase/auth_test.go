package usecase

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestShouldHashPasswordAndDecode(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	td.CmpNoError(t, err)
	td.CmpNot(t, password, hash)

	// Should success when same password
	err = CompareHashPassword("password", hash)
	td.CmpNoError(t, err)

	// Should fail when different password
	err = CompareHashPassword("other-password", hash)
	td.CmpNotNil(t, err)
	td.Cmp(t, err.Error(), "invalid credentials")
}
