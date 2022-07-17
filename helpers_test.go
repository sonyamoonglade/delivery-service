package tgdelivery

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSixifyOrderId(test *testing.T) {

	tt := []TestId{
		{
			id:       1,
			expected: "000001",
		},
		{
			id:       2712,
			expected: "002712",
		},
		{
			id:       111111,
			expected: "111111",
		},
	}

	for _, t := range tt {
		actual := SixifyOrderId(t.id)
		expected := t.expected
		assert.Equal(test, actual, expected)
	}

}

type TestId struct {
	id       int64
	expected string
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
		assert.Equal(test, actual, t.ok)
	}
}

type TestPhone struct {
	ph string
	ok bool
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
		assert.Equal(test, actual, t.ok)
	}
}

type TestName struct {
	name string
	ok   bool
}
