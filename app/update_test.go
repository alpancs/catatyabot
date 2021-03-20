package app

import "testing"

func Test_cross(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not crossed yet",
			args: args{"*kacang 5.000* dicatat ya bos 👌 #catatan"},
			want: `~kacang 5\.000~ dicatat ya bos 👌 \#catatan`,
		},
		{
			name: "already crossed",
			args: args{"~kacang 5.000~ dicatat ya bos 👌 #catatan"},
			want: `~kacang 5\.000~ dicatat ya bos 👌 \#catatan`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cross(tt.args.text); got != tt.want {
				t.Errorf("cross() = %v, want %v", got, tt.want)
			}
		})
	}
}
