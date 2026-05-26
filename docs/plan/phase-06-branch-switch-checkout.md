# Phase 6: Branch, Switch & Detached HEAD

> **Meta:** 3 ngày, navigation state. Day 8 tạo/list branch; Day 9 switch branch; Day 10 checkout commit và detached HEAD.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Navigation Sprint: Moving HEAD Safely

> **Mục tiêu phase:** Tạo/list branch, switch branch cập nhật HEAD/index/worktree, checkout commit và biểu diễn detached HEAD rõ ràng.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-08-03 -> 2026-08-05
>
> **Mini-git surface:** `mgit branch`, `mgit switch`, `mgit checkout`, detached HEAD

---

## Phase Overview

### Missions
- Day 8 - Branch
- Day 9 - Switch branch
- Day 10 - Checkout và detached HEAD

### Flutter / Dart Bridge
> Flutter Navigator push/pop là additive: thêm route lên stack, pop quay lại. Switch branch trong Git không phải push/pop — nó là snapshot restoration: xóa tracked files của branch cũ, ghi files từ tree mới, update index, update HEAD. Mọi bước đều filesystem I/O có thể fail. Detached HEAD = HEAD trỏ thẳng commit hash không qua branch, giống Navigator không có named route. Điểm quan trọng: không có “undo” nếu restore làm mất file chưa commit.

### Go Skills Required For This Phase
> `os.Remove`, `os.RemoveAll` cho tracked file cleanup, `filepath.Walk` cho restore, `strings.HasPrefix` cho ref parsing, `os.ReadDir` cho branch listing. Atomic write pattern: ghi temp file rồi rename. `go test -race ./...` bắt buộc từ phase này.

---

## Day 8 Mission: Branch pointer và branch listing

### User Story
> Khách hàng nói: *"Tôi cần implement branch creation và branch listing."*
>
> Context: Đây là Day 8 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 8 - Branch

### Acceptance Criteria
- [ ] Implement `mgit branch <name>`
- [ ] Tạo branch mới bằng cách ghi current commit hash vào `.mgit/refs/heads/<name>`
- [ ] Validate branch name cơ bản
- [ ] Không cho tạo branch trùng tên
- [ ] Implement `mgit branch` không argument để list branches trong `.mgit/refs/heads`
- [ ] Đánh dấu branch hiện tại bằng `*`
- [ ] Chạy `mgit branch` thấy `main`
- [ ] Chạy `mgit branch dev`, kiểm tra `.mgit/refs/heads/dev` tồn tại
- [ ] Chạy `mgit branch` thấy `main` và `dev`, branch hiện tại có dấu `*`
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Switch branch không phải đổi biến `currentBranch`. Nó restore snapshot vào disk. Động tác này nguy hiểm vì ghi đè working tree. Nếu không check dirty state, CLI của ông sẽ ăn mất bài làm của chính ông."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Branch cực kỳ nhẹ. Tạo branch không copy source code, chỉ tạo một pointer mới tới commit hiện tại."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/refs/refs.go
package refs

type Branch struct{}
type HEADState struct{}

func CreateBranch(name string) error
func ListBranches() ([]Branch, error)
func SetHEADToBranch(name string) error
func SetHEADDetached(id object.ObjectID) error

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-06-D8-CLI: Gắn command `branch` vào dispatcher.
// SENIOR ASKS: Vì sao tạo branch chỉ ghi file ref, không copy working directory?

// TODO-06-D8-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-06-D8-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-06-D8-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-06-D8-D: Validate branch name trước khi ghi file.
// SENIOR ASKS: Nếu branch name chứa `../`, nó có thể ghi ra ngoài `.mgit/refs/heads` không?
```

#### Theory Notes From CSV
- [ ] Ghi chú: tạo branch chỉ tạo pointer mới, không copy source code

#### Socratic Questions
1. Switch branch phải thay đổi HEAD, index và working tree theo thứ tự nào?
2. Tạo branch có copy source code không? Bằng chứng nằm ở file nào?
3. Detached HEAD nguy hiểm vì mất cái gì: object, branch pointer hay working file?
4. Dirty working tree nên block, warn hay overwrite? Vì sao?

### Output Checklist: Làm sao biết mình xong?
- [ ] Chạy `mgit branch` thấy `main`
- [ ] Chạy `mgit branch dev`, kiểm tra `.mgit/refs/heads/dev` tồn tại
- [ ] Chạy `mgit branch` thấy `main` và `dev`, branch hiện tại có dấu `*`
- [ ] Branch file chứa đúng current commit hash, không chứa snapshot content
- [ ] Tạo branch trùng tên trả error rõ

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] `branch dev` tạo `.mgit/refs/heads/dev` bằng current commit.
- [ ] `branch` list được current branch có dấu `*`.
- [ ] Branch name rỗng, chứa slash nguy hiểm, hoặc `..` bị reject.
- [ ] Repo chưa có commit thì `branch dev` trả error rõ hoặc tạo branch rỗng theo policy đã ghi.
- [ ] Detached HEAD khi list branch không giả vờ có current branch.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về branch pointer

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tái sử dụng status clean check trước switch không?
2. Restore tree có xóa tracked files cũ không?
3. API refs có biểu diễn được symbolic và detached HEAD bằng type rõ không?

---

## Day 9 Mission: Switch branch và restore snapshot

### User Story
> Khách hàng nói: *"Tôi cần implement `mgit switch <branch>`."*
>
> Context: Đây là Day 9 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 9 - Switch branch

### Acceptance Criteria
- [ ] Load commit mà branch trỏ tới
- [ ] Restore tree ra working directory
- [ ] Update index theo tree mới
- [ ] Phiên bản đơn giản: xóa tracked files cũ rồi ghi file từ tree mới
- [ ] Báo lỗi hoặc warning đơn giản nếu working directory chưa clean
- [ ] Resolve branch name thành `refs/heads/<branch>`
- [ ] Update `.mgit/HEAD` thành `ref: refs/heads/<branch>`
- [ ] Load tree của commit đó
- [ ] Switch về `main` và thấy working directory quay về nội dung của `main`
- [ ] Switch lại `dev` và thấy nội dung của `dev`
- [ ] Tạo branch `dev`, switch sang `dev`, commit thay đổi trên `dev`
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Switch branch không phải đổi biến `currentBranch`. Nó restore snapshot vào disk. Động tác này nguy hiểm vì ghi đè working tree. Nếu không check dirty state, CLI của ông sẽ ăn mất bài làm của chính ông."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Switch/checkout không chỉ đổi HEAD. Nó còn cập nhật working tree và index theo commit mà branch mới trỏ tới."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/refs/refs.go
package refs

type Branch struct{}
type HEADState struct{}

func CreateBranch(name string) error
func ListBranches() ([]Branch, error)
func SetHEADToBranch(name string) error
func SetHEADDetached(id object.ObjectID) error

// File: internal/worktree/worktree.go
package worktree

func RestoreTree(id object.ObjectID) error
func IsClean() (bool, error)
func UpdateIndexToTree(id object.ObjectID) error

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-06-D9-CLI: Gắn command `switch` vào dispatcher.
// SENIOR ASKS: Vì sao switch phải update HEAD, index và working tree chứ không chỉ HEAD?

// TODO-06-D9-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-06-D9-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-06-D9-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-06-D9-D: Check working tree clean trước khi restore.
// SENIOR ASKS: Nếu dirty file bị overwrite, user mất dữ liệu ở bước nào?
```

#### Theory Notes From CSV
- [ ] Ghi chú: switch/checkout cập nhật HEAD, working tree và index

#### Socratic Questions
1. Switch branch phải thay đổi HEAD, index và working tree theo thứ tự nào?
2. Tạo branch có copy source code không? Bằng chứng nằm ở file nào?
3. Detached HEAD nguy hiểm vì mất cái gì: object, branch pointer hay working file?
4. Dirty working tree nên block, warn hay overwrite? Vì sao?

### Output Checklist: Làm sao biết mình xong?
- [ ] Switch về `main` và thấy working directory quay về nội dung của `main`
- [ ] Switch lại `dev` và thấy nội dung của `dev`
- [ ] Tạo branch `dev`, switch sang `dev`, commit thay đổi trên `dev`
- [ ] `.mgit/HEAD` đổi thành `ref: refs/heads/<branch>`
- [ ] `.mgit/index` đồng bộ với tree của branch mới

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] `switch dev` restore đúng nội dung của dev.
- [ ] Switch branch không tồn tại trả error rõ.
- [ ] Dirty working tree bị block hoặc warning theo policy đã ghi.
- [ ] Tracked file chỉ có ở branch cũ bị xóa khỏi working tree khi switch.
- [ ] Index sau switch phản ánh tree mới, không giữ hash cũ.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về việc switch thay đổi working tree

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tái sử dụng status clean check trước switch không?
2. Restore tree có xóa tracked files cũ không?
3. API refs có biểu diễn được symbolic và detached HEAD bằng type rõ không?

---

## Day 10 Mission: Checkout commit và detached HEAD

### User Story
> Khách hàng nói: *"Tôi cần implement checkout trực tiếp một commit hash."*
>
> Context: Đây là Day 10 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 10 - Checkout commit và detached HEAD

### Acceptance Criteria
- [ ] Implement `mgit checkout <commit-hash>`
- [ ] Restore tree từ commit ra working directory
- [ ] Update index theo tree của commit
- [ ] Ghi `.mgit/HEAD` trực tiếp thành commit hash
- [ ] Khi HEAD không bắt đầu bằng `ref:`, coi là detached HEAD
- [ ] Update commit logic để cảnh báo hoặc xử lý detached HEAD đơn giản
- [ ] Hiển thị warning khi vào detached HEAD
- [ ] Validate commit object tồn tại trước khi checkout
- [ ] `mgit branch` hiểu trạng thái detached HEAD hoặc báo rõ không có current branch
- [ ] Có ít nhất 2 commit, checkout commit cũ và thấy working directory quay về nội dung cũ
- [ ] `.mgit/HEAD` chứa commit hash trực tiếp
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Switch branch không phải đổi biến `currentBranch`. Nó restore snapshot vào disk. Động tác này nguy hiểm vì ghi đè working tree. Nếu không check dirty state, CLI của ông sẽ ăn mất bài làm của chính ông."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Detached HEAD nghĩa là HEAD trỏ trực tiếp tới commit, không thông qua branch. Nếu commit tiếp trong detached HEAD, commit đó có thể không được branch nào giữ lại."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/refs/refs.go
package refs

type Branch struct{}
type HEADState struct{}

func CreateBranch(name string) error
func ListBranches() ([]Branch, error)
func SetHEADToBranch(name string) error
func SetHEADDetached(id object.ObjectID) error

// TODO-DAY10: HEADState phải biểu diễn được hai trạng thái riêng biệt:
// 1. Symbolic: “ref: refs/heads/main” → cần type indicator + branch name
// 2. Detached: raw commit hash → cần type indicator + ObjectID
// SENIOR ASKS: Không có sealed class như Dart trong Go. Dùng union type kiểu nào? Interface hay struct với field?

// File: internal/worktree/worktree.go
package worktree

func RestoreTree(id object.ObjectID) error
func IsClean() (bool, error)
func UpdateIndexToTree(id object.ObjectID) error

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-06-D10-CLI: Gắn command `checkout` vào dispatcher và cảnh báo detached HEAD.
// SENIOR ASKS: Detached HEAD mất branch pointer, không mất object. Giải thích bằng file `.mgit/HEAD`.

// TODO-06-D10-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-06-D10-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-06-D10-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-06-D10-D: Validate commit object trước khi ghi HEAD detached.
// SENIOR ASKS: Nếu ghi HEAD trước rồi phát hiện hash không phải commit, repo rơi vào state gì?
```

#### Theory Notes From CSV
- [ ] Ghi chú: detached HEAD trỏ trực tiếp tới commit, không qua branch

#### Socratic Questions
1. Switch branch phải thay đổi HEAD, index và working tree theo thứ tự nào?
2. Tạo branch có copy source code không? Bằng chứng nằm ở file nào?
3. Detached HEAD nguy hiểm vì mất cái gì: object, branch pointer hay working file?
4. Dirty working tree nên block, warn hay overwrite? Vì sao?

### Output Checklist: Làm sao biết mình xong?
- [ ] `mgit branch` hiểu trạng thái detached HEAD hoặc báo rõ không có current branch
- [ ] Có ít nhất 2 commit, checkout commit cũ và thấy working directory quay về nội dung cũ
- [ ] `.mgit/HEAD` chứa commit hash trực tiếp
- [ ] Checkout hash không tồn tại không đổi HEAD
- [ ] Warning detached HEAD nói rõ commit sau đó có thể không được branch giữ

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] `checkout <old>` ghi HEAD trực tiếp bằng commit hash và cảnh báo detached.
- [ ] Checkout object không phải commit bị reject.
- [ ] Dirty working tree bị block hoặc warning theo policy đã ghi.
- [ ] `branch` trong detached HEAD báo rõ detached state.
- [ ] Commit trong detached HEAD được cảnh báo hoặc xử lý theo policy đã ghi.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về rủi ro commit trong detached HEAD

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn có tái sử dụng status clean check trước switch không?
2. Restore tree có xóa tracked files cũ không?
3. API refs có biểu diễn được symbolic và detached HEAD bằng type rõ không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-06-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, branch không tồn tại, dirty worktree, hash sai format hoặc object không phải commit.

### CP-06-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case refs filesystem, dirty worktree, restore tree và detached HEAD tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.
- [ ] Chạy `go test -race ./...` xanh — phase này có filesystem restore multi-step có thể race.

### CP-06-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Tạo branch bằng copy working directory thay vì pointer file.
- Switch chỉ đổi HEAD mà không restore index và working tree.
- Checkout detached HEAD mà không cảnh báo, commit mới dễ không có branch giữ.

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
Switch branch phải thay đổi HEAD, index và working tree theo thứ tự nào?
