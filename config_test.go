package main

import "testing"

func Test_ipcut(t *testing.T) {

	tests := []struct {
		name string
		args string
		want string
	}{
		{"basic", "127.0.0.1:3128", "127.0.0.1"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipcut(tt.args); got != tt.want {
				t.Errorf("ipcut() = %v, want %v", got, tt.want)
			}
		})
	}
}
