package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestId struct {
	id       int64
	expected string
}

type TestOrderId struct {
	text     string
	expected int64
}

type TestAmount struct {
	text     string
	expected int64
}

type TestUsername struct {
	text     string
	expected string
}

func TestSixifyOrderId(test *testing.T) {

	tt := []TestId{
		{id: 1, expected: "000001"},
		{id: 2712, expected: "002712"},
		{id: 111111, expected: "111111"},
	}

	for _, t := range tt {
		actual := SixifyOrderId(t.id)
		expected := t.expected
		assert.Equal(test, expected, actual)
	}

}

func TestExtractOrderId(test *testing.T) {
	tt := []TestOrderId{
		{text: "Заказ #000027", expected: 27},
		{text: "Заказ #000123", expected: 123},
		{text: "Заказ #005123", expected: 5123},
		{text: "Заказ #000000", expected: 0},
		{text: "Заказ #000001", expected: 1},
		{text: "Заказ #123456", expected: 123456},
		{text: "Заказ #999999", expected: 999999},
	}

	for _, t := range tt {
		actual := ExtractOrderId(t.text)
		assert.Equal(test, t.expected, actual)
	}
}

func TestExtractAmount(test *testing.T) {

	tt := []TestAmount{
		{text: "Сумма | 981.0 ₽", expected: 981},
		{text: "Сумма | 1005123.0 ₽", expected: 1005123},
		{text: "Сумма | 1.0 ₽", expected: 1},
		{text: "Сумма | 222.0 ₽", expected: 222},
		{text: "Сумма | 0.0 ₽", expected: 0},
		{text: "Сумма | 99999.0 ₽", expected: 99999},
	}

	for _, t := range tt {
		actual := ExtractAmount(t.text)
		assert.Equal(test, t.expected, actual)
	}
}

func TestExtractUsername(test *testing.T) {

	tt := []TestUsername{
		{text: "Заказчик: Артем Тимофеев", expected: "Артем Тимофеев"},
		{text: "Заказчик: Боб Иванов", expected: "Боб Иванов"},
		{text: "Заказчик: Aleксандр Siмонеnko", expected: "Aleксандр Siмонеnko"},
		{text: "Заказчик: Шишкин Лес", expected: "Шишкин Лес"},
		{text: "Заказчик: Kirill Lexov", expected: "Kirill Lexov"},
		{text: "Заказчик: K L", expected: "K L"},
		{text: "Заказчик: kk Ll", expected: "kk Ll"},
	}

	for _, t := range tt {
		actual := ExtractUsername(t.text)
		assert.Equal(test, t.expected, actual)
	}

}
