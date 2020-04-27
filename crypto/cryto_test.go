package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	pass := HashPassword("123456")
	if len(pass) != 32 {
		t.Fail()
	}

	t.Log(pass)

}

func TestGenerateOpenIdForImiFromUserId(t *testing.T) {

	pass1 := GenerateOpenIdForImiFromUserId(1)
	pass2 := GenerateOpenIdForImiFromUserId(2)

	if len(pass1) != 32 {
		t.Fail()
	}

	t.Log(pass1)
	t.Log(pass2)
	assert.NotEqual(t, pass1, pass2)
}

func TestHash(t *testing.T) {

	pass1 := HashDeviceId("00000000-3a3d-215a-ffff-fffff69f0d9b")

	if len(pass1) != 32 {
		t.Fail()
	}

	t.Log(pass1)
}
