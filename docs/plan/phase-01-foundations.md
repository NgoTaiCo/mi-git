# Phase 1: Repo Foundations & CLI Bootstrap

> **Meta:** 1 ngày, 1 mini-project nhỏ. Dựng `mgit` executable và metadata repository đầu tiên. Phase này chỉ có một mission, không chia topic giả tạo.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Foundation Sprint: `.mgit` Repository Shell

> **Mục tiêu phase:** Từ zero -> có CLI `mgit init`, repo root discovery, `.mgit/HEAD`, object/ref folders và behavior idempotent.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-07-27
>
> **Mini-git surface:** `.mgit/HEAD`, `.mgit/objects`, `.mgit/refs/heads`, `internals/repo`

---

## Phase Overview

### Mission
- Day 1 - CLI skeleton + `mgit init`

### Flutter / Dart Bridge
> Trong Flutter, `flutter create` scaffold mọi thứ sẵn sàng — framework lifecycle quản lý hộ. Go CLI không có magic lifecycle: `main` là entry point, và bạn tự explicit tạo từng thứ. Filesystem contract ở đây giống project structure nhưng không có framework tạo cho bạn. Package `internal/repo` đóng vai trò giống usecase layer trong Clean Architecture: không biết CLI, không biết output format, chỉ biết filesystem state.

### Go Skills Required For This Phase
> `os.MkdirAll`, `os.WriteFile`, `os.ReadFile`, `os.Stat`, `os.Args` hoặc package `flag`, `filepath.Join`, `fmt.Errorf` với `%w`, Go module (`go mod init`), package visibility (`internal/`).

---

## Day 1 Mission: CLI skeleton + `mgit init`

### User Story
> Khách hàng nói: *"Tôi cần tạo CLI skeleton cho `mgit` và implement lệnh `mgit init`."*
>
> Context: Đây là Day 1 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 1 - CLI skeleton + `mgit init`

### Acceptance Criteria
- [x] Tạo Go module cho project `mgit`
- [x] Tạo `main.go` và command dispatcher bằng `os.Args`
- [x] Tạo package `internals/repo` — (lưu ý: path thực tế là `internals/`, không phải `internal/`)
- [x] Implement `FindRoot` đi ngược từ current directory để tìm `.mgit`
- [x] Implement `Init` tạo `.mgit/objects`, `.mgit/refs/heads`, `.mgit/HEAD`
- [x] Ghi `.mgit/HEAD` với nội dung `ref: refs/heads/main\n` (có trailing newline)
- [x] `Init` idempotent: chạy lại không phá repo, `MkdirAll` không lỗi khi thư mục tồn tại, `WriteFile` ghi đè `HEAD`
- [x] Wire lệnh `mgit init` vào CLI dispatcher trong `main.go` — `case "init":` gọi `repo.Init(os.Getwd())`
- [x] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Nếu tôi nhận requirement này, tôi không nghĩ ngay đến command đẹp. Tôi nghĩ đến filesystem contract: command nào được phép tạo file nào, chạy lại có phá state không, và test repo tạm thời ra sao. Flutter dev hay quen có framework lifecycle; CLI Go không có lifecycle magic, `main` chỉ là cửa vào rồi mọi thứ là explicit."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Git không cần database server. Một repository chỉ là working directory cộng với thư mục metadata chứa object, refs và HEAD."
```

#### TODO Comments (Skeleton / Contract Only)
```go
package main

// run tách khỏi main() để test trực tiếp không cần spawn process
func run(args []string) int
// Trả về: 0 = success, 1 = no args, 2 = unknown command

// TODO-01-CLI: Command dispatcher chỉ parse command và gọi package domain.
// SENIOR ASKS: Tại sao tách `run(args []string) int` khỏi `main()` giúp test CLI dễ hơn?
// HINT: Trong test, bạn truyền args trực tiếp thay vì spawn process.
// STATUS: ✅ Đã implement. Cần thêm case "init" để gọi repo.Init.

package repo // path thực tế: internals/repo

const DefaultMetaDir = ".mgit"

// Repository giữ hai đường dẫn: thư mục làm việc và thư mục metadata
type Repository struct {
    Worktree string // thư mục gốc chứa code
    MetaDir  string // đường dẫn đến .mgit
}

// InitOptions dành cho flag/tuỳ chọn tương lai (hiện tại rỗng)
type InitOptions struct{}

// FindRoot đi ngược từ start lên filesystem đến khi tìm thấy .mgit
func FindRoot(start string) (Repository, error)

// Init tạo cấu trúc .mgit mới hoặc re-init nếu đã tồn tại (idempotent)
func Init(path string, opts InitOptions) (Repository, error)

// IsRepository kiểm tra path có chứa .mgit hợp lệ không
// LƯU Ý: trả về (bool, error) — khác với (bool) trong Git thật
// error xảy ra khi .mgit tồn tại nhưng không phải thư mục (file hoặc symlink)
func IsRepository(path string) (bool, error)

// TODO-01-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?
// STATUS: ✅ Đã quyết định — tất cả filesystem logic nằm trong internals/repo.

// TODO-01-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?
// STATUS: ✅ Đã có TestIsRepository (4 cases), TestFindRoot (4 cases), TestInit (4 cases).

// TODO-01-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.
// STATUS: ✅ main.go chỉ parse args, không biết gì về filesystem.

// TODO-01-D: `Init` tạo repo tại target path; `FindRoot` tìm repo đã tồn tại.
// SENIOR ASKS: Vì sao `mgit init` không nên phụ thuộc vào việc current directory đã có `.mgit`?
// STATUS: ✅ Init nhận path explicit, không dùng os.Getwd() ngầm.
```

#### Theory Notes From CSV
- [x] Ghi chú: repo = working directory + metadata folder + object database
  > `Repository.Worktree` = working directory; `Repository.MetaDir` = `.mgit/` = metadata folder; `objects/` bên trong = object database (chưa có nội dung, phase sau mới dùng)

#### Socratic Questions
1. Tại sao `.mgit/HEAD` là file text vẫn đủ biểu diễn current branch?
2. Chạy `mgit init` lần hai nên là success, warning hay error? Vì sao?
3. Repo root discovery khác gì tìm nearest ancestor trong routing tree của Flutter?
4. Nếu `FindRoot` trả zero-value `Repository` cùng error, caller phải xử lý thế nào để không dùng nhầm state rác?

### Output Checklist: Làm sao biết mình xong?
- [x] Domain function `Init()` đã tạo đúng `.mgit` khi test với `t.TempDir()`
- [x] `.mgit/HEAD` chứa `ref: refs/heads/main\n` (verified bởi TestInit block 7)
- [x] `Init` lần 2 không phá repo — TestInit testcase 2 kiểm tra re-init với repo đã có đầy đủ cấu trúc
- [x] Chạy `go build -o mgit.exe . && mgit.exe init` trong thư mục tạm — thấy đúng cấu trúc `.mgit/` — _verified 2026-06-05_

### Test Checklist: Những gì bạn nên tự kiểm tra
- [x] Init trong folder trống tạo đúng `.mgit/objects`, `.mgit/refs/heads`, `.mgit/HEAD` — _TestInit/testcase_1_
- [x] Init lần hai (re-init) không mất HEAD — _TestInit/testcase_2_: setup tạo sẵn cấu trúc đầy đủ, Init ghi đè HEAD với cùng nội dung
- [x] FindRoot từ subfolder `a/b/c` tìm đúng repo root — _TestFindRoot/testcase_2_
- [x] FindRoot ngoài repo trả error rõ, không panic — _TestFindRoot/testcase_3_
- [ ] Init trong subfolder của repo đang tồn tại tạo repo lồng nhau — _(policy chưa quyết định, chưa test; hiện tại Init không ngăn)_
- [x] `.mgit/HEAD` có trailing newline: `ref: refs/heads/main\n` — _TestInit block 7_ so sánh exact string
- [x] Không tạo `.git`; chỉ tạo `.mgit` — hardcoded qua `DefaultMetaDir = ".mgit"`
- [x] `run([]string{})` → exit 1; `run([]string{"xyz"})` → exit 2; không panic — _TestRun_ (3 cases)

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về điều đã học trong ngày

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn đã tách command parsing khỏi filesystem logic chưa?
2. Nếu sau này thêm `mgit status`, code tìm repo root có reuse được không?
3. Có chỗ nào đang dùng absolute path khiến test phụ thuộc máy bạn không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-01-A: CLI Manual Flow
- [x] Chạy được `mgit init` từ terminal trong thư mục tạm — _verified 2026-06-05, output: `Initialized empty repository in .../.mgit`_
- [x] Quan sát được `objects/`, `refs/heads/`, `HEAD` bên trong `.mgit` sau khi chạy
- [x] Error message rõ: lệnh không xác định trả exit 2; `repo.Init` lỗi in ra stderr rồi trả exit 1

### CP-01-B: Test Gate
- [x] Unit test cho package domain: `TestIsRepository`, `TestFindRoot`, `TestInit` — tổng 12 cases
- [x] Test edge case filesystem: `.mgit` là file, symlink (skip Windows), subfolder nested, path không tồn tại
- [x] Integration test CLI: `TestRunInit` (2 cases: init mới + re-init) trong `main_test.go`
- [x] `go test ./...` xanh — 14 PASS, 1 SKIP — _verified 2026-06-05_
- [x] `go fmt ./...` và `go vet ./...` sạch — _verified 2026-06-05_

### CP-01-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Tạo `.git` thay vì `.mgit` và vô tình phá repo thật.
- `mgit init` không idempotent, chạy lần hai ghi đè state đang có.
- Nhét repo discovery vào `main` khiến mọi command sau phải copy/paste logic.

## Progression Rules

### Rule 1: Không qua phase kế tiếp nếu CLI chỉ chạy happy path.
Phải có test và phải nhìn được state thật trong `.mgit`.

### Rule 2: Không nhét logic vào `main`.
Command handler chỉ parse input, gọi package domain, format output.

### Rule 3: Không dùng Git thật làm oracle duy nhất.
Được so sánh để học, nhưng Mini Git format có giới hạn riêng đã ghi trong roadmap.

## Tổng Kết

### Deliverables
- [x] Command `mgit init` chạy được từ CLI — `go build` + thực thi được xác nhận
- [x] Test xanh: `go test ./...` — 14 cases PASS, 1 SKIP (symlink/Windows)
- [x] Tooling xanh: `go fmt ./...` và `go vet ./...` — clean
- [ ] Note học tập tối thiểu 5 dòng cho mission _(bạn tự viết)_
- [ ] Retrospective ghi rõ sai ở đâu và refactor nào cố tình chưa làm _(bạn tự viết)_

> **Việc còn lại trước khi phase này 100% done (chỉ còn phần của bạn):**
> 1. Viết Learning Notes (mục bên dưới) — tối thiểu 5 dòng
> 2. Viết Retrospective — trả lời 3 câu hỏi ở cuối
> 3. Oral Defense CP-01-C — trả lời First-Principles Question không nhìn code

### First-Principles Question
Tại sao `.mgit/HEAD` là file text vẫn đủ biểu diễn current branch?
