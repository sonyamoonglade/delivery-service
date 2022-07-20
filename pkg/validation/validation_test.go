package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestName struct {
	name string
	ok   bool
}

type TestPhone struct {
	ph string
	ok bool
}

func TestValidatePhoneNumber(test *testing.T) {
	tt := []TestPhone{
		{ph: "+79128507989", ok: true},
		{ph: "+79128412989", ok: true},
		{ph: "+7912850799", ok: false},
		{ph: "79128501799", ok: false},
		{ph: "+1912850799", ok: false},
	}

	for _, t := range tt {
		actual := ValidatePhoneNumber(t.ph)
		assert.Equal(test, t.ok, actual)
	}
}

func TestValidateUsername(test *testing.T) {

	tt := []TestName{
		{name: "Artem", ok: false},
		{name: "Artem B", ok: false},
		{name: "Artem Bo", ok: false},
		{name: "Ar Bobb", ok: false},
		{name: "Ars Bob", ok: true},
		{name: "As ob", ok: false},
		{name: "Arss Bobs", ok: true},
		{name: "Иван Си", ok: false},
		{name: "Ао Григое", ok: false},
		{name: "Аоы Гри", ok: true},
		{name: "Abы CXш", ok: true},
		{name: "", ok: false},
	}
	for _, t := range tt {
		actual := ValidateUsername(t.name)
		test.Logf("exp - %t, got - %t, v - %s", t.ok, actual, t.name)
		assert.Equal(test, t.ok, actual)

	}
}
