# Phase 4: Commit Objects, Refs & History

> **Meta:** 2 ngày, biến snapshot thành lịch sử. Day 4 tạo commit object; Day 5 nối commit vào refs, HEAD và log.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## History Sprint: Commit DAG + Refs

> **Mục tiêu phase:** Tạo commit từ tree, lưu parent, update branch hiện tại và đọc log theo parent chain.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-07-30 -> 2026-07-31
>
> **Mini-git surface:** `mgit commit-tree`, `mgit commit`, `mgit log`, `.mgit/refs/heads/main`

---

## Phase Overview

### Missions
- Day 4 - Commit object và `commit-tree`
- Day 5 - Refs, HEAD, `commit` và `log`

### Flutter / Dart Bridge
> BLoC lưu state history bằng stream events theo thứ tự thời gian. Commit DAG của Git khác: không phải linear list mà là graph hướng ngược — mỗi commit pointer về parent. Branch không phải “version của code” mà là pointer mỏng tới một commit hash. HEAD là pointer-của-pointer: symbolic ref trỏ tới branch, branch trỏ tới commit. Detached HEAD là khi HEAD bỏ qua bước branch pointer — sẽ build ở Phase 6.

### Go Skills Required For This Phase
> `strings.TrimPrefix`, `strings.Split`, `bufio.Scanner` cho multi-line parse, `time.Time` format, multiple return error pattern, `errors.New`/`fmt.Errorf`. File path resolve cho refs. Struct với optional fields (`[]ObjectID` cho parents).

---

## Day 4 Mission: Commit object và `commit-tree`

### User Story
> Khách hàng nói: *"Tôi cần implement commit object tối giản."*
>
> Context: Đây là Day 4 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 4 - Commit object và `commit-tree`

### Acceptance Criteria
- [ ] Implement commit format gồm `tree`, optional `parent`, `author`, timestamp và message
- [ ] Implement `mgit commit-tree <tree-hash> -m "message"`
- [ ] Lưu commit object vào object database
- [ ] Update `cat-file` để hiển thị commit object
- [ ] Implement parser cho commit object
- [ ] Cho phép parent optional khi tạo commit object
- [ ] Chạy `mgit cat-file <commit-hash>` và thấy tree, author, timestamp, message
- [ ] Chạy `mgit write-tree`, rồi `mgit commit-tree <tree-hash> -m "first commit"`
- [ ] `commit-tree` chỉ tạo commit object, chưa update branch hiện tại
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Commit không phải “save file”. Commit là metadata trỏ đến root tree và parent. Branch cũng không phải copy code; branch là pointer file. Đây là lúc cai tư duy OOP object graph trong RAM."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Commit là snapshot metadata. Commit trỏ tới tree. Commit sau có thể trỏ tới commit trước bằng parent. Lịch sử Git là một graph/DAG."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/commit/commit.go
package commit

type Commit struct{}
type Options struct{}

func BuildCommit(tree object.ObjectID, parent []object.ObjectID, opts Options) ([]byte, error)
func ParseCommit(content []byte) (Commit, error)
func CommitTree(store object.Store, tree object.ObjectID, opts Options) (object.ObjectID, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-04-D4-CLI: Gắn command `commit-tree` vào dispatcher và update `cat-file` để đọc commit.
// SENIOR ASKS: Vì sao `commit-tree` không được update `.mgit/refs/heads/main`?

// TODO-04-D4-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-04-D4-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-04-D4-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-04-D4-D: Inject author/time qua Options để test hash deterministic.
// SENIOR ASKS: Nếu timestamp lấy trực tiếp từ time.Now trong BuildCommit, test hash ổn định kiểu gì?
```

#### Theory Notes From CSV
- [ ] Ghi chú: commit là snapshot metadata trỏ tới root tree và parent

#### Socratic Questions
1. `mgit log` bắt đầu từ hash nào và lấy parent ở đâu?
2. HEAD symbolic ref khác commit hash trực tiếp thế nào?
3. Update branch trước khi ghi commit object có rủi ro gì?
4. Commit DAG khác linked list ở điểm nào khi bắt đầu có merge?

### Output Checklist: Làm sao biết mình xong?
- [ ] Chạy `mgit cat-file <commit-hash>` và thấy tree, author, timestamp, message
- [ ] Chạy `mgit write-tree`, rồi `mgit commit-tree <tree-hash> -m "first commit"`
- [ ] `commit-tree` trả commit hash nhưng không ghi `.mgit/refs/heads/main`
- [ ] Commit object có header `commit <size>\0<content>` khi lưu qua object store

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Viết test tạo commit từ tree
- [ ] Commit đầu tiên không có parent parse được.
- [ ] Commit thứ hai có parent là commit đầu.
- [ ] Message nhiều dòng parse lại đúng, không bị cắt mất dòng sau.
- [ ] Timestamp/author deterministic trong test nhờ truyền qua `Options`.
- [ ] Commit thiếu dòng `tree` hoặc tree hash sai length trả error rõ.
- [ ] `cat-file` với commit không in raw bytes mù mờ, mà hiển thị tree/parent/author/message.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về commit object và lịch sử dạng DAG

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tách refs khỏi object store không?
2. Commit parser chịu được message nhiều dòng không?
3. Nếu sau này detached HEAD xuất hiện, API refs hiện tại có vỡ không?

---

## Day 5 Mission: Refs, HEAD, `commit` và `log`

### User Story
> Khách hàng nói: *"Tôi cần implement commit thật sự vào branch hiện tại và xem lịch sử bằng log."*
>
> Context: Đây là Day 5 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 5 - Refs, HEAD, `commit` và `log`

### Acceptance Criteria
- [ ] Resolve HEAD symbolic ref: `ref: refs/heads/main`
- [ ] Implement lấy current commit từ current branch
- [ ] Implement `mgit commit -m "message"` tạo tree, lấy parent, tạo commit và update branch
- [ ] Implement `mgit log` đi ngược theo parent chain
- [ ] Hiển thị log gồm commit hash, author, date và message
- [ ] Tạo package `internal/refs`
- [ ] Implement đọc `.mgit/HEAD`
- [ ] Update `refs/heads/main` sau khi commit
- [ ] Init repo, tạo file, chạy `mgit commit -m "first"`
- [ ] Sửa file, commit lần 2, chạy `mgit log` thấy 2 commit mới nhất trước
- [ ] Phase này `mgit commit` vẫn tạo tree từ working directory; Phase 5 mới đổi sang commit từ index
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Commit không phải “save file”. Commit là metadata trỏ đến root tree và parent. Branch cũng không phải copy code; branch là pointer file. Đây là lúc cai tư duy OOP object graph trong RAM."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Branch chỉ là một file chứa commit hash. HEAD thường là symbolic ref trỏ tới branch hiện tại."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/commit/commit.go
package commit

type Commit struct{}
type Options struct{}

func BuildCommit(tree object.ObjectID, parent []object.ObjectID, opts Options) ([]byte, error)
func ParseCommit(content []byte) (Commit, error)
func CommitTree(store object.Store, tree object.ObjectID, opts Options) (object.ObjectID, error)

// NOTE-DAY5: BuildCommit, ParseCommit, CommitTree đã implement ở Day 4.
// Day 5 tập trung vào refs.Store, HEAD resolution và history Log. Không cần implement lại commit package.

// File: internal/refs/refs.go
package refs

type Store interface {
	ReadHEAD() (string, error)
	CurrentBranch() (string, error)
	ResolveHEAD() (object.ObjectID, error)
	UpdateCurrentBranch(id object.ObjectID) error
}

// File: internal/history/log.go
package history

func Log(start object.ObjectID) ([]commit.Commit, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-04-D5-CLI: Gắn command `commit` và `log` vào dispatcher.
// SENIOR ASKS: Vì sao `commit` nên update branch chỉ sau khi object commit đã ghi thành công?

// TODO-04-D5-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-04-D5-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-04-D5-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-04-D5-D: HEAD ở phase này là symbolic ref, nhưng API phải chuẩn bị cho detached HEAD sau này.
// SENIOR ASKS: Nếu `.mgit/HEAD` không bắt đầu bằng `ref:`, ResolveHEAD nên trả gì?
```

#### Theory Notes From CSV
- [ ] Ghi chú: branch là file chứa commit hash, HEAD thường là symbolic ref

#### Socratic Questions
1. `mgit log` bắt đầu từ hash nào và lấy parent ở đâu?
2. HEAD symbolic ref khác commit hash trực tiếp thế nào?
3. Update branch trước khi ghi commit object có rủi ro gì?
4. Commit DAG khác linked list ở điểm nào khi bắt đầu có merge?

### Output Checklist: Làm sao biết mình xong?
- [ ] Init repo, tạo file, chạy `mgit commit -m "first"`
- [ ] Sửa file, commit lần 2, chạy `mgit log` thấy 2 commit mới nhất trước
- [ ] `.mgit/refs/heads/main` chứa hash của commit mới nhất
- [ ] `.mgit/HEAD` vẫn là symbolic ref `ref: refs/heads/main`
- [ ] `mgit commit-tree` và `mgit commit` khác nhau rõ: một cái tạo object, một cái update branch

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Commit đầu tiên không có parent parse được.
- [ ] Commit thứ hai có parent là commit đầu.
- [ ] `log` hiển thị newest-first.
- [ ] HEAD trỏ branch hiện tại và branch file được update sau commit.
- [ ] Branch file chưa tồn tại thì commit đầu tiên vẫn tạo/update được ref đúng.
- [ ] `log` trên repo chưa có commit trả message rõ, không panic.
- [ ] Parent hash trong commit thứ hai đúng bằng hash commit thứ nhất.
- [ ] Nếu ghi commit object fail, branch ref không được update.
- [ ] HEAD malformed, ví dụ `ref refs/heads/main`, trả error rõ.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về refs, HEAD và parent chain

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tách refs khỏi object store không?
2. Commit parser chịu được message nhiều dòng không?
3. Nếu sau này detached HEAD xuất hiện, API refs hiện tại có vỡ không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-04-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, tree/commit hash sai format, HEAD malformed, hoặc ref không tồn tại.

### CP-04-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case parser, refs filesystem hoặc parent chain tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.

### CP-04-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Commit không lưu parent nên lịch sử không đi ngược được.
- Update branch trước khi ghi commit object xong, tạo ref trỏ tới object chưa tồn tại.
- Coi HEAD luôn là symbolic ref, làm phase detached HEAD sau này vỡ.

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
- [ ] Note học tập tối thiểu 5 dòng cho mỗi mission.
- [ ] Retrospective ghi rõ sai ở đâu và refactor nào cố tình chưa làm.

### First-Principles Question
`mgit log` bắt đầu từ hash nào và lấy parent ở đâu?
