package utils

import "testing"

func TestString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{"success", args{"test"}, func() *string {
			s := "test"

			return &s
		}()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.s); *got != *tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
