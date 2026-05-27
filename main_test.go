package main

import "testing"

// dùng *testing.T để lấy state giá trị nhất quán đúng vị trí nhớ nhờ pointer
func TestRun(t *testing.T) {
	// tạo 1 table giá trị input, và giá trị expected
	// []struct là tạo slice các đối tượng annonymous dùng nhanh để test
	cases := []struct {
		name string   // tên của giá trị input
		args []string // args truyền vào từ CLI
		want int      // expected
	}{
		{
			name: "no args returns error code",
			args: []string{},
			want: 1,
		},
		{
			name: "help return success",
			args: []string{"help"},
			want: 0,
		},
		{
			name: "unknown command returns error",
			args: []string{"unknown"},
			want: 2,
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			got := run(testcase.args)
			if got != testcase.want {
				t.Fatalf("run(%v) = %d, want %d", testcase.args, got, testcase.want)
			}
		})
	}
}
