package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultMetaDir = ".mgit"

type Repository struct {
	Worktree string
	MetaDir  string
}

type InitOptions struct{}

func FindRoot(start string) (Repository, error) {
	currentPath := start

	for {
		isMGitRepo, err := IsRepository(currentPath)

		if err != nil {
			return Repository{}, fmt.Errorf("not mgit repo at %v: %w", currentPath, err)
		}
		if isMGitRepo {
			return Repository{
				Worktree: currentPath,
				MetaDir:  filepath.Join(currentPath, DefaultMetaDir),
			}, nil
		}

		parent := filepath.Dir(currentPath)

		if parent == currentPath {
			return Repository{}, fmt.Errorf("repository not found")
		}

		currentPath = parent
	}
}

func Init(path string, opts InitOptions) (Repository, error) {
	return Repository{}, fmt.Errorf("")
}

func IsRepository(path string) (bool, error) {
	// lấy đường dẫn hiện tại + .mgit
	mgitPath := filepath.Join(path, DefaultMetaDir)
	// lấy thông tin
	// Lstat lấy thẳng tới symlink luôn tránh việc symlink che mất thông tin của file/folder hiện tại
	info, err := os.Lstat(mgitPath)

	// xử lý lỗi nếu có
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat %q: %w", mgitPath, err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%q exist but is not a directory", mgitPath)
	}
	return true, nil
}

// TODO-01-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-01-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-01-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-01-D: `Init` tạo repo tại target path; `FindRoot` tìm repo đã tồn tại.
// SENIOR ASKS: Vì sao `mgit init` không nên phụ thuộc vào việc current directory đã có `.mgit`?
