# Phase 2: Object Database & Blob Plumbing

> **Meta:** 1 ngày, 1 plumbing layer. Biến file content thành object addressable bằng hash. Phase này chỉ có một mission: blob object database.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Object Sprint: Blob, SHA-1, Zlib

> **Mục tiêu phase:** Hiểu Git object không phải record database, mà là immutable bytes có header, hash và compressed storage path.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-07-28
>
> **Mini-git surface:** `internal/object`, `mgit hash-object`, `mgit cat-file`, `.mgit/objects/ab/cdef...`

---

## Phase Overview

### Mission
- Day 2 - Object database: `hash-object` và `cat-file`

### Flutter / Dart Bridge
> Trong Dart, `hashCode` là identity hint — hai object khác nhau có thể share `hashCode`. SHA-1 trong Go ngược lại: đây là deterministic content digest. Hai blob có cùng SHA-1 thì guaranteed cùng bytes — đây là content-addressable storage. Gần nhất với Dart là nếu bạn hash toàn bộ serialized state để detect duplicate event, nhưng SHA-1 có guarantee cryptographic, không chỉ là collision hint. Byte slice `[]byte` khác `List<int>` của Dart: không có hidden copy, mutable, và phải cẩn thận với `io.Reader` pipeline.

### Go Skills Required For This Phase
> `crypto/sha1`, `compress/zlib` (`zlib.NewWriter`, `zlib.NewReader`), `bytes.Buffer`, `io.Reader`/`io.Writer`, byte slice `[]byte`, `fmt.Sprintf` với `%x`, `os.MkdirAll` + `os.Create`, `defer` cho resource cleanup.

---

## Day 2 Mission: Object database, `hash-object`, `cat-file`

### User Story
> Khách hàng nói: *"Tôi cần implement object database tối giản với blob object."*
>
> Context: Đây là Day 2 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 2 - Object database: `hash-object` và `cat-file`

### Acceptance Criteria
- [ ] Tạo package `internal/object`
- [ ] Implement object format `blob <size>\0<content>`
- [ ] Tính SHA-1 trên toàn bộ object bytes gồm header và content
- [ ] Nén object bằng zlib trước khi lưu
- [ ] Lưu object vào `.mgit/objects/ab/cdef...` theo 2 ký tự đầu hash
- [ ] Implement `mgit hash-object <file>`
- [ ] Implement `mgit cat-file <hash>` đọc object, giải nén zlib, parse header và in content
- [ ] Tạo `hello.txt`, chạy `mgit hash-object hello.txt`, kiểm tra object được lưu
- [ ] Chạy `mgit cat-file <hash>` và thấy nội dung file gốc
- [ ] Thay đổi nội dung file và xác nhận hash thay đổi
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Ở đây cấm nghĩ “ID là field trong model” kiểu app backend. Object ID của Git là kết quả của bytes. Sai một byte header là hash khác. Nếu ông hash content thô rồi gọi là Git-like object storage thì đó là học vẹt."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau ngày này Git là content-addressable storage. Object ID không phải ID ngẫu nhiên, mà được tính từ chính nội dung object."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/object/object.go
package object

type ObjectType string
type ObjectID string
type Header struct{}
type StoreOptions struct{}

type Store interface {
	WriteObject(kind ObjectType, content []byte) (ObjectID, error)
	ReadObject(id ObjectID) (ObjectType, []byte, error)
}

func FormatObject(kind ObjectType, content []byte) []byte
func ParseObject(raw []byte) (ObjectType, []byte, error)
func HashObject(raw []byte) ObjectID
func ObjectPath(root string, id ObjectID) (string, error)
func NewStore(repoRoot string, opts StoreOptions) Store

// File: main.go hoặc internal/cli package
package main

func run(args []string) int

// TODO-02-CLI: Gắn command `hash-object` và `cat-file` vào dispatcher có từ Phase 1.
// SENIOR ASKS: CLI nên đọc file và gọi object.Store ở đâu để test package object không phụ thuộc stdout?

// TODO-02-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-02-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-02-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-02-D: Hash pipeline phải là format object -> hash raw object bytes -> zlib compress -> write file.
// SENIOR ASKS: Nếu bạn nén trước rồi mới hash, ObjectID còn đại diện cho nội dung object không?
```

#### Theory Notes From CSV
- [ ] Ghi chú: vì sao object ID được tính từ nội dung object

#### Socratic Questions
1. Vì sao hash phải tính trên `blob <size>\0<content>`, không chỉ file content?
2. Zlib nằm trước hay sau SHA-1 trong pipeline? Vì sao?
3. Nếu content chứa byte zero, parser header phải dừng ở đâu?
4. Dart `hashCode` khác content-addressable hash ở điểm chết người nào?

### Output Checklist: Làm sao biết mình xong?
- [ ] Tạo `hello.txt`, chạy `mgit hash-object hello.txt`, kiểm tra object được lưu
- [ ] Chạy `mgit cat-file <hash>` và thấy nội dung file gốc
- [ ] Thay đổi nội dung file và xác nhận hash thay đổi
- [ ] Object nằm đúng path `.mgit/objects/<2-char-prefix>/<38-char-suffix>`
- [ ] Không tạo tree/commit/index trong phase này

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Viết test nhỏ cho hash và read object
- [ ] Hash cùng content tạo cùng ObjectID.
- [ ] Đổi một byte content làm hash đổi.
- [ ] `cat-file` đọc lại đúng bytes gốc.
- [ ] Hash invalid length hoặc object thiếu file trả error rõ.
- [ ] Known vector: content `hello\n` với object format `blob 6\0hello\n` tạo hash `ce013625030ba8dba906f756967f9e9ca394464a`.
- [ ] `ParseObject` phát hiện header thiếu NUL byte.
- [ ] `ParseObject` phát hiện size trong header không khớp content length.
- [ ] Object content có byte zero vẫn round-trip đúng vì parser chỉ tách ở NUL đầu tiên sau header.
- [ ] Object file corrupt hoặc zlib invalid trả error rõ, không panic.
- [ ] `hash-object` ngoài repo trả error rõ vì không tìm thấy `.mgit`.

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng note về content-addressable storage

### Retrospective: Sau khi xong, hãy tự hỏi
1. Bạn parse bytes hay convert string bừa bãi?
2. Object store có phụ thuộc current working directory không?
3. Error có wrap đủ context path/hash để debug không?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-02-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, file input không tồn tại, hash sai format, hoặc object không tồn tại.

### CP-02-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case filesystem, zlib hoặc parser tương ứng phase.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.

### CP-02-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Hash file content thô thay vì hash object bytes có header.
- Nén zlib trước khi hash, làm ObjectID lệch khỏi model Git.
- Parse object header bằng string split mơ hồ, hỏng khi content chứa byte zero.

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
Vì sao hash phải tính trên `blob <size>\0<content>`, không chỉ file content?
