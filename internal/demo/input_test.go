package demo

import "testing"

func TestByteToCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   byte
		want string
	}{
		{name: "ctrl-c", in: 3, want: "quit"},
		{name: "space", in: ' ', want: "space"},
		{name: "lowercase", in: 'q', want: "q"},
		{name: "uppercase", in: 'Q', want: "q"},
		{name: "newline ignored", in: '\n', want: ""},
		{name: "symbol ignored", in: '1', want: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := byteToCommand(tc.in); got != tc.want {
				t.Fatalf("byteToCommand(%d) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
