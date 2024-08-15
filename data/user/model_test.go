package user

import (
	"fmt"
	"testing"
)

func TestUser_TableName(t *testing.T) {
	u := User{Uid: 12345}
	fmt.Println(u.TableName())
	fmt.Println(u.String())
}
