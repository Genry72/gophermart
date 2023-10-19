package slices

import "testing"

//func TestInslice(t *testing.T) {
//	type args[T comparable] struct {
//		elems []T
//		v     T
//	}
//	type testCase[T comparable] struct {
//		name string
//		args args[T]
//		want bool
//	}
//	tests := []testCase[int, string]{
//		{
//			name: "#1",
//			args: args[string]{[]string{"1"}, "2"},
//			want: false,
//		},
//		{
//			name: "#1",
//			args: args[string]{[]string{"1", "2"}, "2"},
//			want: true,
//		},
//		{
//			name: "#1",
//			args: args[string]{[]string{"1", "2"}, "1"},
//			want: true,
//		},
//		{
//			name: "#1",
//			args: args[int]{[]int{1}, 2},
//			want: false,
//		},
//		{
//			name: "#1",
//			args: args[int]{[]int{1, 2}, 2},
//			want: true,
//		},
//		{
//			name: "#1",
//			args: args[int]{[]int{1, 2}, 1},
//			want: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := Inslice(tt.args.elems, tt.args.v); got != tt.want {
//				t.Errorf("Inslice() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestInslice(t *testing.T) {
	type args[T comparable] struct {
		elems []T
		v     T
	}

	type testCaseString[T comparable] struct {
		name string
		args args[T]
		want bool
	}

	type testCaseInt[T comparable] struct {
		name string
		args args[T]
		want bool
	}

	testsString := []testCaseString[string]{
		{
			name: "#1",
			args: args[string]{elems: []string{"1"}, v: "2"},
			want: false,
		},
		{
			name: "#2",
			args: args[string]{elems: []string{"1", "2"}, v: "2"},
			want: true,
		},
		{
			name: "#3",
			args: args[string]{elems: []string{"1", "2"}, v: "1"},
			want: true,
		},
	}

	testsInt := []testCaseInt[int]{
		{
			name: "#4",
			args: args[int]{elems: []int{1}, v: 2},
			want: false,
		},
		{
			name: "#5",
			args: args[int]{elems: []int{1, 2}, v: 2},
			want: true,
		},
		{
			name: "#6",
			args: args[int]{elems: []int{1, 2}, v: 1},
			want: true,
		},
	}

	for _, tt := range testsString {
		t.Run(tt.name, func(t *testing.T) {
			if got := Inslice(tt.args.elems, tt.args.v); got != tt.want {
				t.Errorf("Inslice() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range testsInt {
		t.Run(tt.name, func(t *testing.T) {
			if got := Inslice(tt.args.elems, tt.args.v); got != tt.want {
				t.Errorf("Inslice() = %v, want %v", got, tt.want)
			}
		})
	}
}
