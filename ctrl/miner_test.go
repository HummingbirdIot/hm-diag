package ctrl

import (
	"fmt"
	"testing"
)

func TestParseSnapshotStateResult(t *testing.T) {
	r, err := parseSnapshotStateResult("time=,file=,state=pending")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
}
