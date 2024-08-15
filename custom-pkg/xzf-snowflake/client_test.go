package xzf_snowflake

import (
	"fmt"
	"testing"
)

func TestGenID(t *testing.T) {
	if err := Init(StartTime, MachineId); err != nil {
		return
	}
	fmt.Println(GenID("AA"))
}
