package hash

import (
	"testing"
)

func TestComputeSHA256(t *testing.T) {

	tests := []struct {
		name string
		data string
		key  string
		want string
	}{
		{
			name: "short data",
			data: "d1",
			key:  "secretkey",
			want: "a1d9be656e1e582467e8470e61fd5949c5cea6663b1d1fb06cba3770926dcc61",
		},
		{
			name: "long data",
			data: "d1h72dnawiu7ft3irianfjbdkfhvidsjgnrsgkn/f;iobhdgohnzot8g74-wbks,vnlsdghIUEnwekgl9t8yoh4bng,nb.ms>?M",
			key:  "secretkey",
			want: "c7f46f0e33ef1c878ddac249279abf462c7a41eeffc28520a98373d3d3ca4e8a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeSHA256([]byte(tt.data), tt.key); got != tt.want {
				t.Errorf("ComputeSHA256 = %v, want %v", got, tt.want)
			}
		})
	}
}
