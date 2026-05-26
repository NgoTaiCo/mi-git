# Phase 5: Index, Add & Status

> **Meta:** 2 ngày, staging area. Day 6 tạo index/add; Day 7 dùng HEAD/index/worktree để tính status.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Staging Sprint: Index + Status Matrix

> **Mục tiêu phase:** Implement `.mgit/index`, `mgit add`, commit từ staged content và status bằng so sánh ba state.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-08-01 -> 2026-08-02
>
> **Mini-git surface:** `.mgit/index`, `mgit add`, `mgit status`, HEAD/index/worktree comparison

---

## Phase Overview

### Missions
- Day 6 - Index / staging area và `mgit add`
- Day 7 - `mgit status`

### Flutter / Dart Bridge
> Trong BLoC, `pendingState` là state chưa được emit ra UI. Index của Git đóng vai trò tương tự: bộ nhớ tạm giữa working directory và commit. `add` không “lưu file” — nó ghi content vào object database và cập nhật map `path→hash` trong index. `commit` flush index thành tree+commit object. `status` là derived state từ 3 nguồn: HEAD, index, working directory — không được lưu sẵn, phải tính lại mỗi lần.

### Go Skills Required For This Phase
> `encoding/json` (`json.Marshal`, `json.Unmarshal`, struct tags), `filepath.Rel` cho path normalization, `os.ReadFile`/`os.WriteFile` cho JSON index, map iteration, `filepath.ToSlash` cho Windows path compat.

---

## Day 6 Mission: Index / staging area và `mgit add`

### User Story
> Khách hàng nói: *"Tôi cần implement staging area bằng file `.mgit/index`."*
>
> Context: Đây là Day 6 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 6 - Index / staging area và `mgit add`

### Acceptance Criteria
- [ ] Implement load/save `.mgit/index`
- [ ] Implement `mgit add <file>` đọc file từ working directory
- [ ] Khi add: tạo blob object và cập nhật `.mgit/index`
- [ ] Support `mgit add` nhiều file trong một lần chạy
- [ ] Update `mgit commit` để tạo tree từ index, không phải toàn bộ working directory
- [ ] Tạo package `internal/index`
- [ ] Thiết kế JSON index map path -> `{ mode, hash }`
- [ ] Sau commit, giữ index đồng bộ với commit mới
- [ ] Tạo `a.txt` và `b.txt`, chỉ `mgit add a.txt`, commit và kiểm tra chỉ `a.txt` được commit
- [ ] `mgit add b.txt`, commit tiếp và kiểm tra commit mới có `b.txt`
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Đây là phase dễ học vẹt nhất. `add` không đánh dấu tên file suông. Nó ghi content thành blob rồi lưu path -> hash vào index. Nếu commit vẫn lấy working tree mới nhất thì ông vừa làm staging area giả."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này `git add` không chỉ đánh dấu file. Nó ghi nội dung file vào object database và cập nhật index. `git commit` commit nội dung đã staged, không nhất thiết là toàn bộ working directory."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/index/index.go
package index

type Entry struct{}
type Index struct{}
type AddOptions struct{}

type Store interface {
	Load() ([]Entry, error)
	Save(entries []Entry) error
}

func Add(paths []string, opts AddOptions) error
func BuildTreeFromIndex(entries []Entry) (object.ObjectID, error)
func NormalizePath(repoRoot string, path string) (string, error)

// File: internal/commit/commit.go
package commit

func CommitFromIndex(message string) (object.ObjectID, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-05-D6-CLI: Gắn command `add` vào dispatcher và đổi `commit` sang dùng index.
// SENIOR ASKS: Vì sao `commit` đọc working directory trực tiếp làm staging area vô nghĩa?

// TODO-05-D6-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-05-D6-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-05-D6-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-05-D6-D: Index lưu relative path -> mode/hash, không lưu absolute path.
// SENIOR ASKS: Nếu index lưu `C:\Users\...`, repo còn portable không?
```

#### Theory Notes From CSV
- [ ] Ghi chú: `git add` ghi content vào object database và cập nhật index

#### Socratic Questions
1. Sau `add`, sửa tiếp file thì commit phải lấy bytes nào?
2. HEAD, index, working directory mỗi cái nằm ở đâu trên disk?
3. Vì sao `status` là derived state, không nên lưu sẵn?
4. Index JSON dễ hiểu nhưng khác binary index thật của Git ở trade-off nào?

### Output Checklist: Làm sao biết mình xong?
- [ ] Tạo `a.txt` và `b.txt`, chỉ `mgit add a.txt`, commit và kiểm tra chỉ `a.txt` được commit
- [ ] `mgit add b.txt`, commit tiếp và kiểm tra commit mới có `b.txt`
- [ ] `.mgit/index` là JSON dễ đọc, lưu relative path và blob hash
- [ ] Sửa file sau `add` không làm index tự đổi nếu chưa `add` lại
- [ ] `commit` sau Phase 5 tạo tree từ index, không từ toàn bộ working directory

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Chỉ add `a.txt`, commit không chứa `b.txt`.
- [ ] `add` nhiều file trong một command cập nhật nhiều entry.
- [ ] `add` cùng file hai lần update hash mới, không duplicate entry.
- [ ] Index path là relative, dùng slash ổn định, không phụ thuộc Windows absolute path.
- [ ] Missing file khi `add` trả error rõ, không tạo index entry rác.
- [ ] Sau commit, index vẫn đồng bộ với tree mới.
- [ ] Optional `add <folder>` nếu làm thì bỏ qua `.mgit` và recurse deterministic.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về khác biệt working directory và staging area
- [ ] Support `mgit add <folder>` recursively

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn normalize path trong index ra sao?
2. Status report có tách data khỏi presentation chưa?
3. Có case xóa file tracked chưa được mô hình hóa không?

---

## Day 7 Mission: `mgit status`

### User Story
> Khách hàng nói: *"Tôi cần implement status bằng cách so sánh HEAD, index và working directory."*
>
> Context: Đây là Day 7 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 7 - `mgit status`

### Acceptance Criteria
- [ ] Load tree từ HEAD commit
- [ ] Convert tree thành map `path -> blob hash`
- [ ] Load index thành map `path -> blob hash`
- [ ] Scan working directory thành map `path -> current hash` và bỏ qua `.mgit`
- [ ] So sánh HEAD với index để tìm staged changes
- [ ] Tìm untracked files
- [ ] So sánh index với working directory để tìm unstaged changes
- [ ] Output các nhóm: `Changes to be committed`, `Changes not staged`, `Untracked files`
- [ ] File staged hiện trong `Changes to be committed`
- [ ] File sửa sau khi add hiện trong `Changes not staged`
- [ ] File mới chưa add hiện trong `Untracked files`; repo clean báo working tree clean
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> “Status không phải feature thêm vào. Nó là kết quả tự nhiên khi bạn đã có HEAD, index và working directory. Nếu tôi có ba map `path→hash` từ ba nguồn này, status chỉ là set difference và intersection. Phần khó không phải logic so sánh — phần khó là đảm bảo ba nguồn đó được load đúng từ đúng chỗ trên disk.”
>
> “Vấn đề status hay bị làm sai: commit rồi mà status vẫn thấy 'staged change' vì index không được sync sau commit. Hoặc untracked files bị nhiễm `.mgit` vì không skip đúng prefix. Test matrix của status phức tạp hơn add vì phải test từng combination của 3-state comparison.”
>
> “Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thạo tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell.”
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/status/status.go
package status

type Report struct{}
type Change struct{}

func CompareHEADIndexWorktree() (Report, error)
func FormatReport(report Report) string

// File: internal/index/index.go
package index

type Entry struct{}
type Store interface {
	Load() ([]Entry, error)
	Save(entries []Entry) error
}

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-05-D7-CLI: Gắn command `status` vào dispatcher.
// SENIOR ASKS: Vì sao status report nên là data structure trước, rồi mới format ra text CLI?

// TODO-05-D7-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-05-D7-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-05-D7-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-05-D7-D: Status là derived state từ HEAD, index, working directory.
// SENIOR ASKS: Nếu lưu status vào file, khi nào nó stale?
```

#### Theory Notes From CSV
- [ ] Ghi chú: `git status` là so sánh HEAD, index và working directory

#### Socratic Questions
1. Sau `add`, sửa tiếp file thì commit phải lấy bytes nào?
2. HEAD, index, working directory mỗi cái nằm ở đâu trên disk?
3. Vì sao `status` là derived state, không nên lưu sẵn?
4. Index JSON dễ hiểu nhưng khác binary index thật của Git ở trade-off nào?

### Output Checklist: Làm sao biết mình xong?
- [ ] File staged hiện trong `Changes to be committed`
- [ ] File sửa sau khi add hiện trong `Changes not staged`
- [ ] File mới chưa add hiện trong `Untracked files`; repo clean báo working tree clean
- [ ] File đã tracked rồi bị xóa được báo rõ theo policy bạn chọn
- [ ] `.mgit` không xuất hiện trong `Untracked files`

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Repo clean sau commit: không còn staged, không còn unstaged cho files đã commit.
- [ ] File staged hiển thị trong Changes to be committed.
- [ ] File sửa sau add hiển thị Changes not staged.
- [ ] File mới chưa add hiển thị Untracked files; clean repo báo clean.
- [ ] Repo chưa có commit nhưng đã có index vẫn status được hoặc báo rõ limitation.
- [ ] Empty index + empty worktree báo clean.
- [ ] File staged rồi sửa tiếp xuất hiện đồng thời ở staged và unstaged nếu policy của bạn hỗ trợ; nếu không thì note limitation.
- [ ] `.mgit` và object files không bị scan thành untracked.
- [ ] Status comparison dùng blob hash, không chỉ so sánh timestamp/size.
- [ ] HEAD là detached commit: `status` vẫn tính staged diff đúng không?

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về ba trạng thái của Git

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn normalize path trong index ra sao?
2. Status report có tách data khỏi presentation chưa?
3. Có case xóa file tracked chưa được mô hình hóa không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-05-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, path không tồn tại, index corrupt hoặc HEAD chưa có commit.

### CP-05-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case index JSON, path normalization, delete/modify/untracked và parser tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.

### CP-05-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- `commit` vẫn đọc working directory thay vì index, làm staging area vô nghĩa.
- Lưu absolute path trong index khiến repo không portable và test phụ thuộc máy.
- Status trộn staged/unstaged/untracked, user không biết commit sẽ chứa gì.

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
Sau `add`, sửa tiếp file thì commit phải lấy bytes nào?
