package main

import (
	"os"
	"path/filepath"
	"testing"
)

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

// TestRunInit kiểm tra lệnh "init" trong CLI dispatcher.
// Dùng t.Chdir để thay đổi working directory an toàn — Go 1.24+ tự khôi phục sau test.
func TestRunInit(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T) string // trả về dir để chdir vào
		want  int
	}{
		{
			name: "init trong folder trống trả về 0",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			want: 0,
		},
		{
			name: "re-init khi .mgit đã tồn tại vẫn thành công",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				// dựng sẵn cấu trúc .mgit đầy đủ để mô phỏng repo có sẵn
				mgit := filepath.Join(dir, ".mgit")
				_ = os.MkdirAll(filepath.Join(mgit, "objects"), 0755)
				_ = os.MkdirAll(filepath.Join(mgit, "refs", "heads"), 0755)
				_ = os.WriteFile(filepath.Join(mgit, "HEAD"), []byte("ref: refs/heads/main\n"), 0644)
				return dir
			},
			want: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setup(t)
			// t.Chdir thay đổi cwd cho sub-test này và tự khôi phục sau khi test xong
			t.Chdir(dir)

			got := run([]string{"init"})
			if got != tc.want {
				t.Fatalf("run([init]) = %d, want %d", got, tc.want)
			}

			// kiểm tra .mgit thực sự được tạo trên filesystem
			info, err := os.Stat(filepath.Join(dir, ".mgit"))
			if err != nil {
				t.Fatalf("expected .mgit directory: %v", err)
			}
			if !info.IsDir() {
				t.Fatal(".mgit phải là thư mục")
			}

			// kiểm tra HEAD chứa đúng nội dung
			headContent, err := os.ReadFile(filepath.Join(dir, ".mgit", "HEAD"))
			if err != nil {
				t.Fatalf("read HEAD: %v", err)
			}
			if string(headContent) != "ref: refs/heads/main\n" {
				t.Fatalf("HEAD = %q, want %q", string(headContent), "ref: refs/heads/main\n")
			}
		})
	}
}
