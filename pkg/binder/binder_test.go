package binder

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBindOK(t *testing.T) {

	type MockInput struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}

	//imitate receiving json body from http request
	inp := MockInput{
		Name: "Mock",
		Age:  18,
	}
	bits, err := json.Marshal(inp)
	require.NoError(t, err)
	r := bytes.NewReader(bits)

	var out MockInput

	err = Bind(r, &out)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestBindErr(t *testing.T) {

	type MockInput struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"required"`
	}

	//imitate receiving json body from http request
	inp := MockInput{
		Name: "Mock", //Not passing age...
	}
	bits, err := json.Marshal(inp)
	require.NoError(t, err)
	r := bytes.NewReader(bits)

	var out MockInput

	err = Bind(r, &out)
	require.NotNil(t, err)

	text := err.Error()
	t.Log(text)
	c0 := strings.Contains(text, "validation error")
	c1 := strings.Contains(text, "Field validation for 'Age'")
	c2 := strings.Contains(text, "'required' tag")

	require.True(t, c0)
	require.True(t, c1)
	require.True(t, c2)
}
