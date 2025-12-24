package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadJSON(t *testing.T) {
	tests := []struct {
		setup func() error
		name  string // description of this test case
		// Named input parameters for target function.
		path    string
		want    map[string]any
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setup()
			if !assert.NoError(t, err) {
				fmt.Printf("An error occurred during test setup: %v", err)
				t.Fail()
			}
			got, gotErr := loadJSON(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("loadJSON() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("loadJSON() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("loadJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
