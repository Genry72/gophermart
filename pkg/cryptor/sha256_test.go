package cryptor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash256(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "#1",
			args: args{
				src: "mypass",
			},
			want: "ea71c25a7a602246b4c39824b855678894a96f43bb9b71319c39700a1e045222",
			err:  nil,
		},
		{
			name: "#2",
			args: args{
				src: "",
			},
			want: "",
			err:  ErrEmptyPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sha256(tt.args.src)
			if err != nil {
				assert.ErrorIs(t, err, tt.err)
				return
			}

			assert.Equalf(t, tt.want, got, "want: %s actual: %s", tt.want, got)
		})
	}
}
