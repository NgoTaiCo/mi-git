# Phase 7: Diff, Merge Base & Three-Way Merge

> **Meta:** 3 ngày, graph và reconciliation. Day 11 diff; Day 12 merge base/fast-forward; Day 13 three-way merge/conflict.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Reconciliation Sprint: Snapshot Diff + Merge

> **Mục tiêu phase:** Implement line diff cơ bản, tìm merge base, fast-forward, three-way merge và conflict marker.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-08-06 -> 2026-08-08
>
> **Mini-git surface:** `mgit diff`, `mgit diff --staged`, `mgit merge`, merge base, conflict markers

---

## Phase Overview

### Missions
- Day 11 - Diff
- Day 12 - Merge base
- Day 13 - Three-way merge và conflict

### Flutter / Dart Bridge
> Flutter `setState` không biết “ai đã đổi gì” — nó chỉ rebuild toàn bộ. Git diff là derived: so sánh blob snapshots để tìm thay đổi cụ thể. Three-way merge dùng base commit giống conflict resolution khi hai Riverpod providers update cùng field: bạn cần biết “state gốc là gì” trước khi quyết định ai thắng. Không có base → không biết ai đã diverge từ điểm tách nhánh. Conflict marker ghi vào working tree (disk), không phải object database — object database immutable.

### Go Skills Required For This Phase
> BFS hoặc DFS bằng queue/stack với `[]ObjectID`, map lookup để deduplicate visited commits. String splitting cho line diff. `strings.Join` cho conflict marker. Graph không có cycle nếu thiết kế đúng nhưng phải guard. `go test -race ./...` cần thiết từ phase này.

---

## Day 11 Mission: Diff derived từ snapshot

### User Story
> Khách hàng nói: *"Tôi cần implement diff line-based đơn giản."*
>
> Context: Đây là Day 11 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 11 - Diff

### Acceptance Criteria
- [ ] Implement `mgit diff` so sánh working directory với index
- [ ] Implement `mgit diff --staged` so sánh index với HEAD
- [ ] Đọc nội dung mới từ working directory hoặc index
- [ ] Implement line diff đơn giản bằng LCS hoặc thuật toán đơn giản hơn
- [ ] Output dạng `--- a/file.txt`, `+++ b/file.txt`, `-old line`, `+new line`
- [ ] Đọc blob cũ từ object database
- [ ] Sửa file sau khi add, chạy `mgit diff` thấy unstaged changes
- [ ] Add file, chạy `mgit diff` không còn thay đổi
- [ ] Chạy `mgit diff --staged` thấy staged changes
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Git không lưu diff làm source of truth. Diff là kết quả tính từ snapshot. Merge cũng không phải “lấy file mới hơn”; nó cần base để biết mỗi bên đã đổi gì từ điểm tách nhánh."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Git thường không lưu diff làm dữ liệu chính. Git lưu snapshot, còn diff được tính ra bằng cách so sánh blob khi cần."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/diff/diff.go
package diff

type Hunk struct{}

func WorkingTreeDiff() ([]Hunk, error)
func StagedDiff() ([]Hunk, error)
func FormatUnified(hunks []Hunk) string

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-07-D11-CLI: Gắn command `diff` và `diff --staged` vào dispatcher.
// SENIOR ASKS: Vì sao diff là derived state, không phải dữ liệu lưu trong object database?

// TODO-07-D11-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-07-D11-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-07-D11-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-07-D11-D: Đọc old/new content từ đúng source: HEAD, index hoặc working directory.
// SENIOR ASKS: `diff` và `diff --staged` khác nhau ở cặp snapshot nào?
```

#### Theory Notes From CSV
- [ ] Ghi chú: Git lưu snapshot; diff được tính bằng cách so sánh blob khi cần

#### Socratic Questions
1. Nếu chỉ so sánh ours/theirs mà không có base, case nào không biết ai đổi?
2. Fast-forward khác merge commit thật ở object nào được tạo?
3. Conflict marker phải ghi vào working tree hay object database? Vì sao?
4. Line diff LCS có trade-off gì so với thuật toán đơn giản hơn?

### Output Checklist: Làm sao biết mình xong?
- [ ] Sửa file sau khi add, chạy `mgit diff` thấy unstaged changes
- [ ] Add file, chạy `mgit diff` không còn thay đổi
- [ ] Chạy `mgit diff --staged` thấy staged changes
- [ ] Output không xuất hiện object hash nội bộ thay cho content diff
- [ ] File binary hoặc unreadable được báo limitation rõ

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] `diff` thấy unstaged change sau khi sửa file.
- [ ] `diff --staged` thấy staged change sau add.
- [ ] File chỉ đổi trong working directory không xuất hiện ở `diff --staged`.
- [ ] File staged rồi sửa tiếp xuất hiện ở cả diff staged và diff unstaged nếu policy hỗ trợ.
- [ ] New file, deleted file, empty file có output rõ hoặc limitation rõ.
- [ ] Diff dùng blob content, không chỉ so sánh hash rồi in "changed".

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về snapshot-vs-diff

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tách graph logic khỏi file merge logic không?
2. Merge result có báo rõ conflict paths không?
3. Log có đọc được commit nhiều parent chưa?

---

## Day 12 Mission: Merge base và fast-forward

### User Story
> Khách hàng nói: *"Tôi cần implement phần graph logic của merge: tìm common ancestor."*
>
> Context: Đây là Day 12 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 12 - Merge base

### Acceptance Criteria
- [ ] Implement `mgit merge <branch>` skeleton
- [ ] Resolve target branch commit là theirs
- [ ] Đi ngược parent chain của theirs để tìm ancestor chung gần nhất
- [ ] In/debug được ours, theirs và merge base
- [ ] Đi ngược parent chain của ours để lấy ancestor set
- [ ] Resolve current commit là ours
- [ ] Xử lý fast-forward: nếu ours là ancestor của theirs, update current branch tới theirs và checkout tree
- [ ] Tạo branch `dev` từ `main`, commit trên `dev`, quay về `main`
- [ ] Chạy `mgit merge dev`; nếu `main` chưa đổi, merge fast-forward thành công
- [ ] In/debug được merge base
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Git không lưu diff làm source of truth. Diff là kết quả tính từ snapshot. Merge cũng không phải “lấy file mới hơn”; nó cần base để biết mỗi bên đã đổi gì từ điểm tách nhánh."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Merge không chỉ so sánh hai commit. Git cần tìm tổ tiên chung, gọi là merge base, để biết mỗi bên đã thay đổi gì kể từ điểm tách nhánh."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/merge/merge.go
package merge

type Plan struct{}
type Result struct{}

func FindBase(ours object.ObjectID, theirs object.ObjectID) (object.ObjectID, error)
func PlanMerge(branch string) (Plan, error)
func ApplyMerge(plan Plan) (Result, error)
func IsAncestor(ancestor object.ObjectID, tip object.ObjectID) (bool, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-07-D12-CLI: Gắn `merge <branch>` skeleton và fast-forward path.
// SENIOR ASKS: Fast-forward update ref khác gì tạo merge commit?

// TODO-07-D12-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-07-D12-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-07-D12-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-07-D12-D: Graph logic không được phụ thuộc vào working tree.
// SENIOR ASKS: Merge base là quan hệ giữa commit objects hay giữa file hiện tại?
```

#### Theory Notes From CSV
- [ ] Ghi chú: merge cần tìm tổ tiên chung để biết mỗi bên đổi gì từ điểm tách nhánh

#### Socratic Questions
1. Vì sao fast-forward không cần tạo merge commit object mới?
2. BFS và DFS đều tìm được ancestor — nhưng với DAG lớn, cái nào cho merge base “gần nhất” đáng tin hơn?
3. Nếu DAG có merge commit với 2 parents, `FindBase` traverse bao nhiêu node? Có loop vô hạn không nếu không mark visited?
4. `IsAncestor(A, B)` trả `true` nghĩa là gì theo graph direction? A là parent của B hay ngược lại?

### Output Checklist: Làm sao biết mình xong?
- [ ] Tạo branch `dev` từ `main`, commit trên `dev`, quay về `main`
- [ ] Chạy `mgit merge dev`; nếu `main` chưa đổi, merge fast-forward thành công
- [ ] In/debug được merge base
- [ ] Fast-forward update current branch và restore tree/index theo target commit
- [ ] Không tạo merge commit trong fast-forward

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Merge base đúng trong lịch sử branch tách từ main.
- [ ] Fast-forward update branch mà không tạo merge commit.
- [ ] Branch target không tồn tại trả error rõ.
- [ ] Không có common ancestor trả error rõ hoặc limitation rõ.
- [ ] Ours bằng theirs thì merge báo already up to date.
- [ ] Graph có merge commit nhiều parent vẫn không loop vô hạn.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về merge base và common ancestor

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tách graph logic khỏi file merge logic không?
2. Merge result có báo rõ conflict paths không?
3. Log có đọc được commit nhiều parent chưa?

---

## Day 13 Mission: Three-way merge và conflict

### User Story
> Khách hàng nói: *"Tôi cần implement merge đơn giản cho file text."*
>
> Context: Đây là Day 13 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 13 - Three-way merge và conflict

### Acceptance Criteria
- [ ] Load tree của base, ours và theirs thành map `path -> blob hash`
- [ ] Với mỗi file, lấy content từ base, ours và theirs
- [ ] Rule: nếu `base == ours` và theirs thay đổi thì lấy theirs
- [ ] Rule: nếu `ours == theirs` thì giữ ours
- [ ] Rule: nếu cả ours và theirs thay đổi khác nhau thì tạo conflict
- [ ] Nếu không conflict: ghi merged result, update index và tạo merge commit có 2 parent
- [ ] Update commit format để support nhiều dòng `parent`
- [ ] Rule: nếu `base == theirs` và ours thay đổi thì giữ ours
- [ ] Ghi conflict marker với `<<<<<<< HEAD`, `=======`, `>>>>>>> branch-name`
- [ ] Nếu có conflict: ghi file conflict ra working directory, không tạo merge commit, báo user cần resolve
- [ ] Merge case không conflict hoạt động và tạo merge commit 2 parent
- [ ] `log` hiển thị được merge commit
- [ ] `status` sau conflict cho thấy file cần xử lý
- [ ] Merge case conflict tạo marker đúng
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Git không lưu diff làm source of truth. Diff là kết quả tính từ snapshot. Merge cũng không phải “lấy file mới hơn”; nó cần base để biết mỗi bên đã đổi gì từ điểm tách nhánh."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Three-way merge dùng base, ours và theirs. Conflict xảy ra khi cả hai nhánh cùng thay đổi một vùng nội dung theo cách khác nhau và Git không thể tự quyết định."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/merge/merge.go
package merge

type Plan struct{}
type Result struct{}
type Conflict struct{}

func FindBase(ours object.ObjectID, theirs object.ObjectID) (object.ObjectID, error)
func PlanMerge(branch string) (Plan, error)
func ApplyMerge(plan Plan) (Result, error)
func ThreeWayMerge(base []byte, ours []byte, theirs []byte) (Result, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-07-D13-CLI: Hoàn thiện `merge <branch>` non-fast-forward.
// SENIOR ASKS: Vì sao có conflict thì không được tạo merge commit?

// TODO-07-D13-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-07-D13-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-07-D13-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-07-D13-D: Conflict marker ghi vào working tree, không ghi object database.
// SENIOR ASKS: Object database immutable; vậy file conflict chưa resolved nằm ở đâu?
```

#### Theory Notes From CSV
- [ ] Ghi chú: three-way merge dùng base, ours và theirs

#### Socratic Questions
1. Conflict marker ghi vào working tree — nhưng object database immutable. Vậy file conflict chưa resolved nằm ở đâu xét về repo state?
2. Sau khi có conflict, `mgit status` phải thêm state gì để báo “repo đang trong unresolved merge state”?
3. Merge không conflict tạo merge commit 2 parent. Nếu `mgit log` chỉ follow parent đầu tiên, lịch sử của branch nào bị mất?
4. File-level three-way merge (đủ cho CV demo) khác line-level LCS ở trade-off nào? Khi nào cần nâng lên line-level?

### Output Checklist: Làm sao biết mình xong?
- [ ] Merge case không conflict hoạt động và tạo merge commit 2 parent
- [ ] `log` hiển thị được merge commit
- [ ] `status` sau conflict cho thấy file cần xử lý
- [ ] Merge case conflict tạo marker đúng
- [ ] Conflict không update current branch và không tạo merge commit
- [ ] Merge commit có đúng 2 dòng parent khi không conflict

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Conflict ghi marker và không tạo merge commit.
- [ ] base==ours lấy theirs.
- [ ] base==theirs giữ ours.
- [ ] ours==theirs giữ ours không conflict.
- [ ] ours và theirs cùng đổi khác nhau tạo conflict marker đúng.
- [ ] Merge không conflict update index và tạo commit 2 parent.
- [ ] `log` parse được commit nhiều parent.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về conflict và vì sao Git không tự quyết định được

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tách graph logic khỏi file merge logic không?
2. Merge result có báo rõ conflict paths không?
3. Log có đọc được commit nhiều parent chưa?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-07-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, branch không tồn tại, merge base không tìm được hoặc conflict chưa resolved.

### CP-07-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case diff output, graph traversal, fast-forward, conflict marker và merge commit nhiều parent.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.
- [ ] Chạy `go test -race ./...` xanh — phase này có graph traversal và filesystem state.

### CP-07-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Thiết kế diff như source of truth thay vì derived từ snapshot.
- Tìm merge base bằng so sánh hai tip commit, bỏ qua ancestor chain.
- Có conflict nhưng vẫn tạo merge commit, làm lịch sử nói dối.

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
Nếu chỉ so sánh ours/theirs mà không có base, case nào không biết ai đổi?
