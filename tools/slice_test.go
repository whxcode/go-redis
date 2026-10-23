package tools

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	ages := []int{10, 20}

	agesStr := Map[int, string](ages, func(age int, i int) string {
		return strconv.FormatInt(int64(age), 10) + "呵呵呵"
	})

	if len(agesStr) != len(ages) {
		t.Errorf("Expected length %d, got %d", len(ages), len(agesStr))
		return
	}
}
