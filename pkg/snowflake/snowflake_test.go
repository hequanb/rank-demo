package snowflake

import (
	"fmt"
	"testing"
)

func TestOnGenId(t *testing.T) {
	
	Init("2021-07-01", 01)
	fmt.Println(GenId())
}
