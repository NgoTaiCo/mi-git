package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsRepository(t *testing.T) {
	cases := []struct {
		name      string
		expected  bool
		expectErr bool
		setup     func(baseDir string) error // vì test không thể test các file/thư mục nên phải có 1 hàm tạo file/thư mục để test
	}{
		{
			name:      "testcase 1: .mgit là thư mục hợp lệ",
			expected:  true,
			expectErr: false,
			setup: func(baseDir string) error {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				return os.Mkdir(mgitPath, 0755)
			},
		},
		{
			name:      "testcase 2: .mgit không tồn tại",
			expected:  false,
			expectErr: false,
			setup:     func(baseDir string) error { return nil },
		},
		{
			name:      "testcase 3: .mgit là 1 file",
			expected:  false,
			expectErr: true,
			setup: func(baseDir string) error {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				return os.WriteFile(mgitPath, []byte{}, 0644)
			},
		},
		{
			name:      "testcase 4: symlink",
			expected:  false,
			expectErr: true,
			setup: func(baseDir string) error {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				return os.Symlink(baseDir, mgitPath)
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			// tạo hộp cát cho test nghịch
			baseDir := t.TempDir()

			// tạo môi trường theo từng cấu hình của từng case
			setupErr := testcase.setup(baseDir)
			// test symlink vì nhiều file/folder sẽ alias trỏ tới 1 nơi khác trên máy
			// khi dùng Lstat thì nó sẽ trỏ tới tận nơi gốc nên phải check symlink
			if setupErr != nil {
				// skip là do trên windows đang rất phiền vụ symlink thiếu quyền
				// tạm thời skip với windows trong thời điểm hiện tại
				t.Skipf("setup failed: %v", setupErr)
			}

			// Chạy hàm cần test
			result, err := IsRepository(baseDir)

			if testcase.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != testcase.expected {
				t.Fatalf("got %v, expected %v", result, testcase.expected)
			}
		})
	}
}

func TestFindRoot(t *testing.T) {
	testcases := []struct {
		name        string // tên testcase
		expectedErr bool   // mong đợi err không
		setup       func(baseDir string) (start string, expectedRoot string, err error)
		// baseDir = thư mục tạm để test dựng filesystem giả
		// start = đường dẫn truyền vào FindRoot
		// expectedRoot = repo root mong muốn FindRoot trả về
	}{
		{
			name:        "testcase 1: thư mục hiện tại có .mgit hợp lệ",
			expectedErr: false,
			setup: func(baseDir string) (start string, expectedRoot string, err error) {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir) // lấy được .mgit path
				resultErr := os.Mkdir(mgitPath, 0755)              // tạo 1 folder với .mgit path với việc hứng lỗi nếu có
				// check xem có lỗi không, nil có nghĩa là không lỗi
				// vì ERROR CŨNG LÀ STATE
				if resultErr != nil {
					// trả về rỗng luôn với lỗi
					return "", "", resultErr
				}
				return baseDir, baseDir, nil
			},
		},

		{
			name:        "testcase 2: tìm root bằng cách trỏ lên parent",
			expectedErr: false,
			setup: func(baseDir string) (start string, expectedRoot string, err error) {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				resultErr := os.Mkdir(mgitPath, 0755)
				if resultErr != nil {
					return "", "", resultErr
				}

				// tạo 1 thư mục giả với chuỗi baseDir/a/b/c
				nestedPath := filepath.Join(baseDir, "a", "b", "c")
				// dùng MkdirAll để tạo 1 chuỗi baseDir/a/b/c
				// vì khi mà dùng Mkdir không, nếu 1 thư mục không tồn tại thì nó không tạo luôn
				// os.Mkdir chỉ tạo 1 cấp cuối, parent phải tồn tại.
				// os.MkdirAll tạo toàn bộ parent còn thiếu.
				resultErr = os.MkdirAll(nestedPath, 0755)
				if resultErr != nil {
					return "", "", resultErr
				}
				return nestedPath, baseDir, nil
			},
		},
		{
			name:        "testcase 3: không thấy .mgit",
			expectedErr: true,
			setup: func(baseDir string) (start string, expectedRoot string, err error) {
				nestedPath := filepath.Join(baseDir, "a", "b", "c")

				err = os.MkdirAll(nestedPath, 0755)
				if err != nil {
					return "", "", err
				}

				return nestedPath, "", nil
			},
		},
		{
			name:        "test 4: thấy .mgit nhưng nó là file",
			expectedErr: true,
			setup: func(baseDir string) (start string, expectedRoot string, err error) {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				// tạo 1 file .mgit với path để tạo 1 file tạm để test
				err = os.WriteFile(mgitPath, []byte{}, 0644)
				if err != nil {
					return "", "", err
				}

				childPath := filepath.Join(baseDir, "a")
				err = os.MkdirAll(childPath, 0755)
				if err != nil {
					return "", "", err
				}

				return childPath, "", nil
			},
		},
	}

	// bắt đầu test từng cái
	for _, testcase := range testcases {
		// chạy test ở mỗi testcase vơi t.Run
		t.Run(testcase.name, func(t *testing.T) {
			// tạo TempDir để có 1 thư mục tạm
			baseDir := t.TempDir()

			// chạy testcase để setup mọi thứ như trên
			start, expectedRoot, setupErr := testcase.setup(baseDir)
			if setupErr != nil {
				// setup lỗi và nhảy lỗi ra
				t.Fatalf("Setup failed: %v", setupErr)
			}

			// chạy test hàm findRoot
			repo, err := FindRoot(start)

			// dụa vào testcase hiện tại có mong muốn nhận được lỗi không
			if testcase.expectedErr {
				// không lỗi
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				// stop hẳn luôn vì đúng rồi, vì đang mong có lỗi
				return
			}

			// sau khi check testcase có lỗi mong lỗi không thì check lỗi không
			// err của FindRoot khi find bị lỗi cover trường hợp không mong muốn
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// kiểm tra kiểu của repo
			if repo.Worktree != expectedRoot {
				t.Fatalf("Got worktree %q, expected %q", repo.Worktree, expectedRoot)
			}
			// lấy mong đợi định dạng cuối
			expectedMetaDir := filepath.Join(expectedRoot, DefaultMetaDir)
			// check lỗi dir mong đợi và dir hiện tại coi có trùng không
			if repo.MetaDir != expectedMetaDir {
				t.Fatalf("Got meta dir %q, expected %q", repo.MetaDir, expectedMetaDir)
			}
		})
	}
}

func TestInit(t *testing.T) {
	testcases := []struct {
		name        string
		expectedErr bool
		setup       func(baseDir string) (targetPath string, err error)
	}{
		{
			name:        "testcase 1: init trong folder trống",
			expectedErr: false,
			setup: func(baseDir string) (string, error) {
				return baseDir, nil
			},
		},
		{
			name:        "testcase 2: init lần 2 cho 1 repo đã tồn tại mgit",
			expectedErr: false,
			setup: func(baseDir string) (string, error) {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				objectsPath := filepath.Join(mgitPath, "objects")
				refsHeadsPath := filepath.Join(mgitPath, "refs", "heads")
				headPath := filepath.Join(mgitPath, "HEAD")

				err := os.MkdirAll(objectsPath, 0755)
				if err != nil {
					return "", err
				}

				err = os.MkdirAll(refsHeadsPath, 0755)
				if err != nil {
					return "", err
				}

				err = os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644)
				if err != nil {
					return "", err
				}

				return baseDir, nil
			},
		},
		{
			name:        "testcase 3: .mgit là 1 file",
			expectedErr: true,
			setup: func(baseDir string) (string, error) {
				mgitPath := filepath.Join(baseDir, DefaultMetaDir)
				err := os.WriteFile(mgitPath, []byte{}, 0644)
				if err != nil {
					return "", err
				}

				return baseDir, nil
			},
		},
		{
			name:        "test 4: target path không tồn tại",
			expectedErr: true,
			setup: func(baseDir string) (string, error) {
				return filepath.Join(baseDir, "missing"), nil
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			baseDir := t.TempDir()

			targetPath, setupErr := testcase.setup(baseDir)
			if setupErr != nil {
				t.Fatalf("setup failed: %v", setupErr)
			}

			repo, err := Init(targetPath, InitOptions{})

			if testcase.expectedErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			//block 1
			if repo.Worktree != targetPath {
				t.Fatalf("got worktree %q, expected %q", repo.Worktree, targetPath)
			}

			//block 2
			expectedMetaDir := filepath.Join(targetPath, DefaultMetaDir)

			if repo.MetaDir != expectedMetaDir {
				t.Fatalf("got meta dỉ %q, expected %q", repo.MetaDir, expectedMetaDir)
			}

			//block 3
			info, err := os.Stat(expectedMetaDir)
			if err != nil {
				t.Fatalf("expected .mgit directory: %v", err)
			}
			if !info.IsDir() {
				t.Fatalf("expected .mgit to be directory")
			}

			//block 4
			objectPath := filepath.Join(expectedMetaDir, "objects")

			info, err = os.Stat(objectPath)
			if err != nil {
				t.Fatalf("expected objects directory: %v", err)
			}
			if !info.IsDir() {
				t.Fatalf("expected objects to be directory")
			}

			//block 5
			refsHeadsPath := filepath.Join(expectedMetaDir, "refs", "heads")

			info, err = os.Stat(refsHeadsPath)
			if err != nil {
				t.Fatalf("expected refs/heads directory: %v", err)
			}
			if !info.IsDir() {
				t.Fatalf("expected refs/heads to be directory")
			}

			//block 6
			headPath := filepath.Join(expectedMetaDir, "HEAD")

			info, err = os.Stat(headPath)
			if err != nil {
				t.Fatalf("expected HEAD file: %v", err)
			}
			if info.IsDir() {
				t.Fatalf("expected HEAD to be file")
			}

			//block 7

			content, err := os.ReadFile(headPath)
			if err != nil {
				t.Fatalf("read HEAD: %v", err)
			}

			if string(content) != "ref: refs/heads/main\n" {
				t.Fatalf("got HEAD %q, expected %q", string(content), "ref: refs/heads/main\n")
			}
		})
	}
}
