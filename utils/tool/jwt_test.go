package tool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(10001)
	assert.Equal(t, nil, err)
	fmt.Println(token)
	parsedToken, err := ParseToken(token)
	assert.Equal(t, nil, err)
	fmt.Printf("%+v\n", parsedToken)
}
