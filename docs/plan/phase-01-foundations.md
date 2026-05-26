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
> **Mini-git surface:** `.mgit/HEAD`, `.mgit/objects`, `.mgit/refs/heads`, `internal/repo`

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
- [ ] Tạo Go module cho project `mgit`
- [ ] Tạo `main.go` và command dispatcher bằng `os.Args` hoặc package `flag`
- [ ] Tạo package `internal/repo`
- [ ] Implement hàm tìm repo root bằng cách đi ngược từ current directory để tìm `.mgit`
- [ ] Implement `mgit init` tạo `.mgit/objects`, `.mgit/refs/heads`, `.mgit/HEAD`
- [ ] Ghi `.mgit/HEAD` với nội dung `ref: refs/heads/main`
- [ ] In message thành công và đảm bảo chạy lại `mgit init` không phá repo
- [ ] Chạy `mgit init` trong folder trống và thấy `.mgit` được tạo
- [ ] Kiểm tra `.mgit/HEAD` trỏ tới `refs/heads/main`
- [ ] Chạy lại `mgit init` và xác nhận repo không bị hỏng
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

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

func run(args []string) int

// TODO-01-CLI: Command dispatcher chỉ parse command và gọi package domain.
// SENIOR ASKS: Tại sao tách `run(args []string) int` khỏi `main()` giúp test CLI dễ hơn?
// HINT: Trong test, bạn truyền args trực tiếp thay vì spawn process.

package repo

type Repository struct{}
type InitOptions struct{}

func FindRoot(start string) (Repository, error)
func Init(path string, opts InitOptions) (Repository, error)
func IsRepository(path string) bool

// TODO-01-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-01-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-01-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-01-D: `Init` tạo repo tại target path; `FindRoot` tìm repo đã tồn tại.
// SENIOR ASKS: Vì sao `mgit init` không nên phụ thuộc vào việc current directory đã có `.mgit`?
```

#### Theory Notes From CSV
- [ ] Ghi chú: repo = working directory + metadata folder + object database

#### Socratic Questions
1. Tại sao `.mgit/HEAD` là file text vẫn đủ biểu diễn current branch?
2. Chạy `mgit init` lần hai nên là success, warning hay error? Vì sao?
3. Repo root discovery khác gì tìm nearest ancestor trong routing tree của Flutter?
4. Nếu `FindRoot` trả zero-value `Repository` cùng error, caller phải xử lý thế nào để không dùng nhầm state rác?

### Output Checklist: Làm sao biết mình xong?
- [ ] Chạy `mgit init` trong folder trống và thấy `.mgit` được tạo
- [ ] Kiểm tra `.mgit/HEAD` trỏ tới `refs/heads/main`
- [ ] Chạy lại `mgit init` và xác nhận repo không bị hỏng

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Init trong folder trống tạo đúng `.mgit/objects`, `.mgit/refs/heads`, `.mgit/HEAD`.
- [ ] Init lần hai không làm mất HEAD hiện có.
- [ ] FindRoot từ subfolder tìm đúng repo root.
- [ ] FindRoot ngoài repo trả error rõ, không panic.
- [ ] Init trong subfolder của repo đã tồn tại không tạo repo lồng nhau nếu đó là policy bạn chọn; nếu cho phép thì phải ghi rõ lý do.
- [ ] `.mgit/HEAD` có newline cuối file hoặc không? Chọn một chuẩn và test ổn định.
- [ ] Không tạo `.git`; chỉ tạo `.mgit`.
- [ ] `run([]string{})` và `run([]string{"unknown"})` trả exit code/message rõ, không panic.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về điều đã học trong ngày

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn đã tách command parsing khỏi filesystem logic chưa?
2. Nếu sau này thêm `mgit status`, code tìm repo root có reuse được không?
3. Có chỗ nào đang dùng absolute path khiến test phụ thuộc máy bạn không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-01-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, command không tồn tại, hoặc path không tạo được.

### CP-01-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case filesystem hoặc parser tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.

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
- [ ] Command chính của phase chạy được.
- [ ] Test xanh: `go test ./...`.
- [ ] Tooling xanh: `go fmt ./...` và `go vet ./...`.
- [ ] Note học tập tối thiểu 5 dòng cho mission.
- [ ] Retrospective ghi rõ sai ở đâu và refactor nào cố tình chưa làm.

### First-Principles Question
Tại sao `.mgit/HEAD` là file text vẫn đủ biểu diễn current branch?
