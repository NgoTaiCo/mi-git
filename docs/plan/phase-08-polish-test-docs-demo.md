# Phase 8: Polish, Test, Docs & Demo

> **Meta:** 1 ngày, hardening. Biến project học tập thành artifact demo được. Phase này chỉ có một release mission.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Release Sprint: Hardening + Learning Retro

> **Mục tiêu phase:** Dọn CLI UX, README, learning note, test suite và demo end-to-end từ init đến merge/status.
>
> **Nguồn:** `docs/14 ngày viết Mini Git bằng Go.csv`
>
> **Lịch:** 2026-08-09
>
> **Mini-git surface:** `README.md`, `docs/what-i-learned-about-git.md`, demo script, `go test ./...`

---

## Phase Overview

### Mission
- Day 14 - Polish, test, docs và demo end-to-end

### Flutter / Dart Bridge
> Flutter project “demo được” cần: golden tests, integration tests, README, release notes. Go CLI project cũng vậy: `go test ./...` xanh, demo script tự tạo repo sạch, README có install/run/example, limitation rõ ràng. Đây là ngày nhìn lại toàn bộ package structure và tự hỏi: nếu reviewer clone repo này không có context gì, họ hiểu được goal và architecture không? `docs/architecture.md` là artifact chứng minh bạn thiết kế, không chỉ implement.

### Go Skills Required For This Phase
> `os.MkdirTemp` cho temp repo trong test, `t.Cleanup` để dọn sau test, `os/exec` nếu viết integration test gọi CLI. Table-driven test pattern `[]struct{ name, input, want string }`. `go test -race ./...` bắt buộc.

---

## Day 14 Mission: Polish, test, docs và demo end-to-end

### User Story
> Khách hàng nói: *"Tôi cần dọn project, viết tài liệu và chạy demo hoàn chỉnh."*
>
> Context: Đây là Day 14 trong roadmap Mini Git. Bạn đang build từng lớp Git internals bằng Go, không dùng library Git có sẵn. Output phải là CLI thật, nhưng design phải test được ở package level.

### CSV Main Task
Day 14 - Polish, test, docs và demo end-to-end

### Acceptance Criteria
- [ ] Refactor code cho dễ đọc, chỉ dọn phần phục vụ project
- [ ] Thêm error handling rõ ràng cho CLI commands
- [ ] `mgit init`, `hash-object`, `cat-file`, `add`, `commit`, `log`, `status` chạy được
- [ ] `mgit branch`, `switch`, `checkout` chạy được
- [ ] `mgit diff` chạy được ở mức cơ bản
- [ ] `mgit merge` chạy được ở mức cơ bản hoặc hiểu rõ merge base và conflict
- [ ] README có hướng dẫn cách dùng và giới hạn của `mgit`
- [ ] `what-i-learned-about-git.md` giải thích được Git internals bằng ngôn ngữ của chính mình
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh
- [ ] Demo chạy từ thư mục tạm sạch, không phụ thuộc state cũ

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Đừng gọi là xong chỉ vì happy path chạy. Project học backend cần deliverable: docs, tests, error UX, và khả năng giải thích lại bằng lời của mình. Không giải thích được thì chỉ là gõ theo checklist."
>
> "Tôi sẽ bẻ task thành ba lớp: contract dữ liệu, thao tác filesystem/object store, rồi CLI command. CLI là vỏ ngoài; logic phải test được mà không cần chạy shell."
>
> "Điểm học chính của ngày này: Tôi nên hiểu được gì sau project Sau 14 ngày, cần giải thích được bằng lời của mình: `git add` làm gì, `git commit` tạo object nào, branch có copy source không, HEAD khác branch thế nào, index dùng để làm gì, checkout tác động tới đâu, Git lưu snapshot hay diff, merge base dùng để làm gì và conflict xảy ra vì sao. Demo script gợi ý: ```sh mgit init echo hello > a.txt mgit add a.txt mgit commit -m "first" mgit branch dev mgit switch dev echo dev > a.txt mgit add a.txt mgit commit -m "change on dev" mgit switch main echo main > a.txt mgit add a.txt mgit commit -m "change on main" mgit merge dev mgit status ```"
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/release/checks.go hoặc scripts/demo spec
package release

type DemoStep struct{}
type Gate struct{}

func ValidateCLIUX() error
func RunDemoScript(steps []DemoStep) error
func CheckReleaseGates(gates []Gate) error

// TODO-08-CLI: Không thêm feature mới nếu nó không phục vụ demo/release gate.
// SENIOR ASKS: Đây là hardening day, không phải ngày nhồi thêm scope. Bạn sẽ cắt scope bằng tiêu chí nào?

// TODO-08-A: Chọn package boundary trước khi code.
// SENIOR ASKS: Logic này thuộc CLI, repo, object store, refs, index hay worktree? Vì sao?

// TODO-08-B: Viết test matrix trước khi implement ruột hàm.
// SENIOR ASKS: Happy path nào chưa đủ? Edge case filesystem nào có thể làm bạn tưởng code đúng?

// TODO-08-C: Implement từng hàm nhỏ, không nhét toàn bộ vào command handler.
// SENIOR ASKS: Nếu đổi CLI flag, package domain có phải sửa không? Nếu có thì boundary đang sai.

// TODO-08-D: Demo script phải tự tạo repo tạm, chạy command, assert output/state, rồi cleanup.
// SENIOR ASKS: Nếu demo chỉ chạy trên folder local đang có sẵn `.mgit`, nó chứng minh được gì?
```

#### Theory Notes From CSV
- [ ] Tự viết lại learning note của ngày này bằng lời của bạn.

#### Socratic Questions
1. Demo end-to-end chứng minh invariant nào của object database, refs và index?
2. README nên nói rõ giới hạn nào để không giả vờ tương thích Git thật?
3. Nếu test chỉ cover happy path, phase nào có nguy cơ hồi quy cao nhất?
4. Bạn giải thích `add`, `commit`, branch, HEAD, index, checkout, snapshot, merge base bằng lời của mình ra sao?

### Output Checklist: Làm sao biết mình xong?
- [ ] `mgit init`, `hash-object`, `cat-file`, `add`, `commit`, `log`, `status` chạy được
- [ ] `mgit branch`, `switch`, `checkout` chạy được
- [ ] `mgit diff` chạy được ở mức cơ bản
- [ ] `mgit merge` chạy được ở mức cơ bản hoặc hiểu rõ merge base và conflict
- [ ] README có hướng dẫn cách dùng và giới hạn của `mgit`
- [ ] `what-i-learned-about-git.md` giải thích được Git internals bằng ngôn ngữ của chính mình
- [ ] Demo script tự tạo repo mới và chạy flow từ `init` tới `merge/status`
- [ ] Release notes/changelog ghi rõ scope làm được và limitation chưa làm
- [ ] CLI error UX nhất quán: command sai, repo thiếu, path/hash/branch sai đều có message rõ
- [ ] `docs/architecture.md` giải thích package structure, data flow và CLI adapter pattern

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Viết test cho object storage
- [ ] Viết test cho commit/log
- [ ] Viết test cho branch/switch
- [ ] Viết test cho index
- [ ] `go test ./...` xanh.
- [ ] Object/index/commit/log/branch/switch test đều có case lỗi.
- [ ] Demo chạy trong thư mục tạm mới, không phụ thuộc state cũ.
- [ ] README command examples khớp CLI thật.
- [ ] `go vet ./...` xanh.
- [ ] Test suite có case lỗi cho parser object/tree/commit, index corrupt, dirty worktree, conflict merge.
- [ ] Demo sau khi chạy không để lại `.mgit` hoặc file rác ngoài thư mục tạm.
- [ ] README nói rõ không support packfile, remote, rebase, stash, tag, hooks, binary index thật.
- [ ] `go test -race ./...` xanh.

### Learning Notes / Docs
- [ ] Viết `README.md` giải thích cách build và dùng `mgit`
- [ ] Viết `docs/what-i-learned-about-git.md`
- [ ] Tạo demo script end-to-end với init/add/commit/branch/switch/merge/status
- [ ] Viết bài tổng kết: add, commit, branch, HEAD, index, checkout, snapshot, merge base, conflict
- [ ] Viết test cho merge conflict nếu kịp

### Retrospective: Sau khi xong, hãy tự hỏi
1. Phần nào cố tình chưa refactor và vì sao?
2. Bug nào làm bạn hiểu Git sâu hơn?
3. Nếu làm lại 14 ngày, phase nào bạn sẽ test sớm hơn?

---

## Phase Checkpoints (BẮT BUỘC)

### CP-08-A: CLI Manual Flow
- [ ] Chạy được toàn bộ command chính của phase trong thư mục tạm.
- [ ] Quan sát được file thay đổi trong `.mgit`, không chỉ nhìn output CLI.
- [ ] Error message rõ với input sai, repo thiếu, hash/branch/path không tồn tại và conflict unresolved.

### CP-08-B: Test Gate
- [ ] Có unit test cho package domain chính.
- [ ] Có test edge case filesystem, parser, refs, index, checkout và merge tương ứng toàn project.
- [ ] Chạy `go test ./...` xanh trước khi qua phase tiếp theo.
- [ ] `docs/architecture.md` mô tả package dependency và adapter design.

### CP-08-C: Oral Defense
- [ ] Trả lời được First-Principles Question cuối file mà không nhìn code.
- [ ] So sánh được concept Go/Mini Git với Dart/Flutter bằng lời của bạn.

## Failure Modes (PHẢI BIẾT)
- Demo chỉ chạy nhờ state cũ trong thư mục local, không chạy được từ repo sạch.
- README mô tả command không khớp CLI thật.
- Test chỉ cover happy path, không bắt lỗi parser/filesystem quan trọng.

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
- [ ] Note học tập tổng kết bằng lời của bạn, không copy từ plan.
- [ ] Retrospective ghi rõ sai ở đâu và refactor nào cố tình chưa làm.

### First-Principles Question
Demo end-to-end chứng minh invariant nào của object database, refs và index?
