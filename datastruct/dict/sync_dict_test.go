package dict_test

import (
	"testing"

	"go-redis/datastruct/dict"
)

func TestSyncDict_Clear(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dict.MakeSyncDict()
			s.Clear()
		})
	}
}
