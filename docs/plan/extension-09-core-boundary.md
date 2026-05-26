# Extension 09: Core Boundary & API-Ready Package Design

> **Meta:** Sau 14 ngày core. Đây không phải phase thêm feature Git mới; đây là phase tách ranh giới package để cùng một core dùng được cho CLI và REST API.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Backend Extension Sprint: Tách Core Khỏi Adapter

> **Mục tiêu extension:** Biến `mgit` từ CLI project thành core library có boundary rõ: CLI chỉ parse input/output, backend API sau này gọi cùng service/domain logic.
>
> **Nguồn:** Sau `docs/plan/concurrency-module.md`
>
> **Ghi chú bàn giao:** Concurrency Module (Session C3) đã tạo `internal/core/service.go` dạng concrete struct làm bài tập context propagation. Extension 09 REDESIGN file đó: thay struct bằng Service interface, thêm domain errors (Init/Add/Commit/Status), tách input/output struct, xóa direct object access ra khỏi core boundary.
>
> **Mini-git surface:** package core/service boundary, domain errors, testable use cases, no CLI dependency inside core

---

## Extension Overview

### Mission
- Extension 09 - Core package boundary cho CLI và API dùng chung

### Flutter / Dart Bridge
> Trong Flutter, Widget không nên tự gọi database rồi tự xử lý business rule. Widget gọi BLoC/usecase; usecase gọi repository. Ở đây CLI và REST handler giống Widget: chỉ nhận input và format output. Core package mới là usecase/domain. Đừng bê kiểu OOP dày class của Dart sang Go; Go cần package nhỏ, interface nhỏ, data contract rõ.
### Go Skills Required For This Extension
> `errors.Is`, `errors.As`, `errors.New`, `fmt.Errorf("%w", err)` error wrapping, custom sentinel error (exported var), custom error type (exported struct + Error() string), `internal/` package visibility rule, interface defined at consumer side (Go idiom), package-level function vs method receiver, `go vet ./...` detect locking mistake.
---

## Mission: Core Boundary Cho CLI Và Backend

### User Story
> Nhà tuyển dụng hỏi: *"Nếu mai anh cần expose project này thành REST API, em có phải copy logic từ CLI không?"*
>
> Câu trả lời đúng: Không. CLI và API chỉ là adapter. Logic `init/add/commit/status/branch/merge` nằm trong core service, test được không cần shell.

### Main Task
Tách package boundary để core Mini Git có thể được gọi bởi CLI hiện tại và REST API ở extension sau.

### Acceptance Criteria
- [ ] `cmd/mgit` hoặc `internal/cli` chỉ parse args, gọi core service, format output
- [ ] Core package không import CLI package
- [ ] Core package không phụ thuộc `os.Args`, stdout, stderr
- [ ] Service method nhận input struct rõ ràng, trả output struct hoặc error rõ ràng
- [ ] Domain error có thể phân loại bằng `errors.Is` hoặc typed error
- [ ] Repository path/data root được truyền qua option, không hard-code working directory toàn cục
- [ ] Test package core gọi trực tiếp service, không cần spawn CLI process
- [ ] Existing CLI demo ở Phase 08 vẫn chạy
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Backend API không cứu được kiến trúc nếu core đang dính vào CLI. Tôi cần nhìn thấy boundary trước khi nhìn thấy HTTP."
>
> "Tôi sẽ tách input/output của từng use case thành struct. CLI adapter biến flags thành input struct. Sau này HTTP adapter biến JSON thành input struct. Core không biết request, response, stdout là gì."
>
> "Điểm học chính: package boundary là cách Go thay thế cho class hierarchy rườm rà. Nếu core sạch, backend API chỉ là adapter mỏng."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/core/service.go
package core

type RepositoryPath string
type ObjectID string

type InitInput struct{}
type InitOutput struct{}
type AddInput struct{}
type AddOutput struct{}
type CommitInput struct{}
type CommitOutput struct{}
type StatusInput struct{}
type StatusOutput struct{}

type Service interface {
	Init(input InitInput) (InitOutput, error)
	Add(input AddInput) (AddOutput, error)
	Commit(input CommitInput) (CommitOutput, error)
	Status(input StatusInput) (StatusOutput, error)
}

// TODO-EXT09-A: CLI chỉ được map args -> input struct -> output text.
// SENIOR ASKS: Nếu đổi từ CLI sang REST JSON, core package có phải sửa không?

// TODO-EXT09-B: Domain error phải phân loại được.
// SENIOR ASKS: API cần trả 404, 409, 400; core error của bạn đang chứa đủ semantic chưa?

// TODO-EXT09-C: Không để core đọc stdout/stderr.
// SENIOR ASKS: Vì sao print trong domain logic làm package khó test và khó expose HTTP?
```

#### Theory Notes
- [ ] Ghi chú: adapter khác domain logic ở đâu
- [ ] Ghi chú: vì sao Go interface nên nhỏ và được định nghĩa ở phía consumer khi hợp lý
- [ ] Ghi chú: khác biệt giữa package boundary trong Go và Clean Architecture layer trong Flutter

#### Socratic Questions
1. Nếu command handler gọi `os.Exit` trong lúc đang xử lý commit, API adapter sẽ tái sử dụng kiểu gì?
2. Error `branch not found` nên là string thường hay typed/sentinel error? Vì sao?
3. Input struct giúp bạn tránh phụ thuộc CLI flag như thế nào?
4. Có package nào đang import ngược từ core sang CLI không?

### Output Checklist: Làm sao biết mình xong?
- [ ] CLI chạy như cũ nhưng logic đã gọi qua core service
- [ ] Unit test có thể gọi `Service.Commit` hoặc use case tương đương trực tiếp
- [ ] Không có `fmt.Println` trong package domain/core
- [ ] Không có `os.Args` trong package domain/core
- [ ] Error domain đủ semantic để adapter CLI/API format riêng

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test core service với temp repo path
- [ ] Test CLI adapter với args sai nhưng không làm domain panic
- [ ] Test domain error mapping ở mức table-driven
- [ ] Chạy `go test -race ./...` để bắt shared state nếu có cache/global

### Learning Notes / Docs
- [ ] Cập nhật `docs/architecture.md` giải thích CLI adapter, future API adapter và core service
- [ ] Cập nhật README mục "Architecture"
- [ ] Ghi một đoạn "Why this is API-ready"

### Retrospective: Sau khi xong, hãy tự hỏi
1. Boundary nào ban đầu đặt sai?
2. Function nào đang nhận quá nhiều thứ vì bạn chưa hiểu use case?
3. Package nào nên nhỏ hơn thay vì gom vào một package "utils"?

---

## Extension Checkpoints (BẮT BUỘC)

### CP-EXT09-A: Boundary Gate
- [ ] Core không import CLI.
- [ ] CLI không chứa business rule Git.
- [ ] Core service test được bằng Go test thường.

### CP-EXT09-B: Error Gate
- [ ] Error có semantic rõ cho not found, invalid input, conflict, dirty worktree.
- [ ] CLI format error thân thiện.
- [ ] Future API có thể map error sang HTTP status.

### CP-EXT09-C: Oral Defense
- [ ] Giải thích được adapter vs core bằng ví dụ CLI và REST.
- [ ] So sánh được boundary này với BLoC/usecase/repository trong Flutter.

## Failure Modes (PHẢI BIẾT)
- Copy logic từ CLI sang API sau này.
- Tạo interface to kiểu `GitService` có 30 method nhưng không có contract rõ.
- Domain logic print trực tiếp ra terminal.
- Dùng global working directory làm test flaky.

## Progression Rules

### Rule 1: Chưa có boundary thì chưa viết HTTP.
REST API chỉ expose core sạch. Không dùng HTTP để che kiến trúc yếu.

### Rule 2: Adapter không quyết định business rule.
Adapter chỉ parse, validate bề mặt, gọi core, format response.

### Rule 3: Core không biết transport.
Core không biết CLI, HTTP, JSON, JWT, Docker.

## Tổng Kết

### Deliverables
- [ ] Core service/usecase boundary.
- [ ] CLI adapter dùng core.
- [ ] Test xanh: `go test ./...`.
- [ ] Architecture note giải thích API-ready design.

### First-Principles Question
Nếu cùng một hành động `commit` được gọi từ CLI và REST API, phần nào là domain invariant bắt buộc nằm trong core, và phần nào chỉ là transport concern của adapter?
