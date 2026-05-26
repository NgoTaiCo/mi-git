# Phase 3: Tree Snapshots

> **Meta:** 1 ngày, snapshot thư mục. Tách content bytes khỏi path/name. Phase này chỉ có một mission: tree object cho working directory.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Snapshot Sprint: Blob vs Tree

> **Mục tiêu phase:** Viết được `write-tree`, tree entry text format, recursive traversal bỏ qua `.mgit`, sort ổn định để hash repeatable.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-07-29
>
> **Mini-git surface:** `mgit write-tree`, tree entry `100644 blob <hash> file.txt`, tree entry `040000 tree <hash> src`

---

## Phase Overview

### Mission
- Day 3 - Tree object: biểu diễn thư mục

### Flutter / Dart Bridge
> Widget tree trong Flutter là live object graph trong RAM — mutable, có lifecycle, rebuild theo state. Tree object trong Mini Git là immutable snapshot trên disk: ghi một lần, không bao giờ sửa. `write-tree` giống như serialize toàn bộ widget tree thành immutable JSON tại một thời điểm — nhưng từng leaf (blob) và node (tree) đều có SHA-1 riêng. Đổi tên widget → hash của tree-parent thay đổi dù blob content giống.

### Go Skills Required For This Phase
> `os.ReadDir`, `filepath.Walk` hoặc `os.DirFS`, `sort.Slice`, `strings.Builder`, multi-return function, recursive function call trong Go. Byte slice format/parse không dùng regex.

---

## Day 3 Mission: Tree object và directory snapshot

### User Story
> Khách hàng nói: *"Tôi cần implement tree object để biểu diễn snapshot thư mục."*
>
> Context: Đây là Day 3 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 3 - Tree object: biểu diễn thư mục

### Acceptance Criteria
- [ ] Implement tree entry format text: `100644 blob <hash> file.txt`
- [ ] Implement tree entry format text: `040000 tree <hash> src`
- [ ] Duyệt working directory recursively và bỏ qua `.mgit`
- [ ] Với file: tạo blob object và lấy blob hash
- [ ] Root directory tạo ra root tree hash
- [ ] Update `cat-file` để hiển thị tree object ở dạng text dễ đọc
- [ ] Sort tree entries theo tên để hash ổn định
- [ ] Implement `mgit write-tree`
- [ ] Với folder: tạo tree object và trả về tree hash
- [ ] Đổi tên file hoặc sửa nội dung và xác nhận tree hash thay đổi
- [ ] Tạo folder có nhiều file/subfolder, chạy `mgit write-tree`, rồi `mgit cat-file <tree-hash>`
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Blob không biết tên file. Nhớ câu này. Nếu ông nhét path vào blob, toàn bộ model Git đổ. Tree là nơi nối name/mode/type/hash lại thành snapshot."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Blob không lưu tên file. Tên file nằm trong tree. Commit không trỏ trực tiếp tới từng file, mà trỏ tới root tree."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/tree/tree.go
package tree

type Entry struct{}
type Tree struct{}
type WriteOptions struct{}

func WriteTree(repoRoot string, opts WriteOptions) (object.ObjectID, error)
func BuildTreeFromDirectory(repoRoot string, dir string, opts WriteOptions) (Tree, error)
func FormatTree(tree Tree) []byte
func ParseTree(content []byte) (Tree, error)
func FlattenTree(store object.Store, id object.ObjectID) (map[string]object.ObjectID, error)

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-03-CLI: Gắn command `write-tree` vào dispatcher và update `cat-file` để đọc tree.
// SENIOR ASKS: Vì sao `cat-file` nên parse theo object type thay vì assume mọi object là blob?

// TODO-03-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-03-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-03-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-03-D: Tree hash phải deterministic.
// SENIOR ASKS: Filesystem traversal order có được OS đảm bảo không? Nếu không sort thì object ID còn đáng tin không?
```

#### Theory Notes From CSV
- [ ] Ghi chú: blob không lưu tên file; tên file nằm trong tree

#### Socratic Questions
1. Đổi tên file nhưng giữ nguyên content thì blob hash và tree hash thay đổi thế nào?
2. Filesystem traversal order có deterministic không? Nếu không sort thì test nào flaky?
3. Vì sao phải bỏ qua `.mgit` khi build tree?
4. Tree object giống hay khác gì widget tree snapshot trong Flutter?

### Output Checklist: Làm sao biết mình xong?
- [ ] Đổi tên file hoặc sửa nội dung và xác nhận tree hash thay đổi
- [ ] Tạo folder có nhiều file/subfolder, chạy `mgit write-tree`, rồi `mgit cat-file <tree-hash>`
- [ ] `cat-file <tree-hash>` hiển thị entry theo format text dễ đọc: mode, type, hash, name
- [ ] `.mgit` không xuất hiện trong tree output
- [ ] Phase này chưa tạo commit/index; tree chỉ là snapshot object

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Folder có nhiều file/subfolder tạo tree đọc lại được.
- [ ] Hai lần `write-tree` cùng content ra cùng hash.
- [ ] Đổi tên file làm tree hash đổi nhưng blob hash cũ vẫn giữ.
- [ ] `.mgit` không xuất hiện trong tree output.
- [ ] Tree entries được sort theo tên để hash ổn định.
- [ ] File và folder cùng tên prefix, ví dụ `app` và `app.txt`, vẫn sort/format ổn định.
- [ ] Empty directory được xử lý theo policy rõ ràng: bỏ qua hay tạo tree rỗng, phải ghi trong note.
- [ ] Path lưu trong tree là relative path/name, không phải absolute path từ máy bạn.
- [ ] Filename có space vẫn parse được nếu format của bạn định nghĩa rõ; nếu chưa support thì phải error rõ.
- [ ] `ParseTree` phát hiện dòng entry thiếu field hoặc hash sai length.
- [ ] `cat-file` với tree object không in raw compressed bytes và không assume blob.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về blob, tree và snapshot thư mục

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn đang lưu path relative hay absolute?
2. Mode/type/hash/name có format đủ rõ để parse lại không?
3. Có case Unicode filename nào cần ghi nhận cho phase sau không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-03-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với repo thiếu, unreadable file/folder, object hash sai format, hoặc tree parse lỗi.

### CP-03-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case filesystem, sort order hoặc parser tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.

### CP-03-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Không sort tree entry nên cùng input sinh hash khác nhau theo filesystem order.
- Đưa `.mgit` vào tree và tự snapshot metadata repository.
- Gắn path/name vào blob khiến model blob-tree sai từ gốc.

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
Đổi tên file nhưng giữ nguyên content thì blob hash và tree hash thay đổi thế nào?
