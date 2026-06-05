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
	// tạo 1 biến giữ đường dẫn hiện tại bằng đầu vào là start
	currentPath := start

	// bắt đầu vòng lặp để đi tìm root với
	for {
		// kiểm tra currentPath hiện tại có phải là 1 repo mgit không
		// trả về 2 biến để kiểm tra
		// isMGitRepo là biến check có phải repo mgit không, err dùng để check có lỗi gì không
		isMGitRepo, err := IsRepository(currentPath)

		// err != nil có nghĩa là err chứa lỗi
		if err != nil {
			// trả về Repo kèm lỗi không phải mgit repo
			return Repository{}, fmt.Errorf("not mgit repo at %v: %w", currentPath, err)
		}

		// nếu biến check current path là repo
		if isMGitRepo {
			// trả về 1 struct repo với đủ thông tin
			return Repository{
				Worktree: currentPath,
				MetaDir:  filepath.Join(currentPath, DefaultMetaDir),
			}, nil
		}

		// nếu không lỗi lẫn không phải repo thì lên cha tìm tiếp
		parent := filepath.Dir(currentPath)

		// nếu trỏ lên mà vẫn bằng currenPath, có nghĩa là đã tới root
		// thì báo lỗi
		if parent == currentPath {
			return Repository{}, fmt.Errorf("repository not found")
		}

		// sau cùng là sẽ đặt lại currentPath là parent để bắt đầu 1 loop mới
		currentPath = parent
	}
}

// Init khởi tạo một repository mới tại path đã cho.
// Nếu repo đã tồn tại (.mgit là thư mục), thực hiện re-init (idempotent).
// Trả về lỗi nếu: path không tồn tại, hoặc .mgit đã tồn tại nhưng là file/symlink.
func Init(path string, opts InitOptions) (Repository, error) {
	// Bước 1: kiểm tra target path có tồn tại và là thư mục không
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Repository{}, fmt.Errorf("path %q không tồn tại", path)
		}
		return Repository{}, fmt.Errorf("stat %q: %w", path, err)
	}
	if !info.IsDir() {
		return Repository{}, fmt.Errorf("path %q không phải thư mục", path)
	}

	metaDir := filepath.Join(path, DefaultMetaDir)

	// Bước 2: kiểm tra nếu .mgit đã tồn tại nhưng KHÔNG phải thư mục (file hoặc symlink)
	// Dùng Lstat để không đi xuyên qua symlink — tránh bị lừa bởi symlink trỏ tới dir
	mgitInfo, err := os.Lstat(metaDir)
	if err == nil && !mgitInfo.IsDir() {
		return Repository{}, fmt.Errorf("%q đã tồn tại nhưng không phải thư mục", metaDir)
	}

	// Bước 3: tạo cây thư mục bên trong .mgit
	// MkdirAll không báo lỗi nếu thư mục đã tồn tại → re-init an toàn
	if err := os.MkdirAll(filepath.Join(metaDir, "objects"), 0755); err != nil {
		return Repository{}, fmt.Errorf("tạo objects: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(metaDir, "refs", "heads"), 0755); err != nil {
		return Repository{}, fmt.Errorf("tạo refs/heads: %w", err)
	}

	// Bước 4: tạo (hoặc ghi đè) file HEAD trỏ về nhánh mặc định
	headPath := filepath.Join(metaDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return Repository{}, fmt.Errorf("tạo HEAD: %w", err)
	}

	return Repository{
		Worktree: path,
		MetaDir:  metaDir,
	}, nil
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
