# Extension 10: REST API Adapter

> **Meta:** Backend API đầu tiên cho Mini Git. Phase này chỉ expose core qua HTTP; không copy logic Git vào handler.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Backend Extension Sprint: HTTP Adapter + API Contract

> **Mục tiêu extension:** Xây REST API dùng `net/http` hoặc `chi`, có JSON request/response chuẩn, error mapping rõ, test bằng `httptest`.
>
> **Nguồn:** Sau `extension-09-core-boundary.md`
>
> **Mini-git surface:** REST endpoints cho repos, status, add, commit, log, branches, checkout, diff, merge

---

## Extension Overview

### Mission
- Extension 10 - REST API adapter cho Mini Git core

### Flutter / Dart Bridge
> Trong Flutter, API client không nên biết database schema nội bộ của backend. Tương tự, HTTP handler không nên biết chi tiết object database, tree parser hay refs. Handler chỉ decode JSON, gọi service, encode JSON. Nó giống presentation layer gọi BLoC/usecase.

### Go Skills Required For This Extension
> `net/http`: `http.HandlerFunc`, `http.ServeMux`, `w http.ResponseWriter`, `r *http.Request`, `http.Error`, `http.StatusXxx` constants. JSON: `json.NewDecoder(r.Body).Decode(&req)`, `json.NewEncoder(w).Encode(resp)`, custom error response struct. Context: `r.Context()`, `context.WithTimeout`, `defer cancel()`. Testing: `httptest.NewRecorder()`, `httptest.NewRequest()`, assert status + body. Middleware: `func(http.Handler) http.Handler` pattern. Router: `chi.NewRouter()` hoặc `http.NewServeMux()` tùy chọn.

---

## Mission: Expose Mini Git Core Qua REST

### User Story
> Nhà tuyển dụng hỏi: *"Project CLI này liên quan gì tới backend?"*
>
> Câu trả lời đúng: Cùng core Mini Git được expose qua REST API có test, context timeout, JSON error chuẩn và boundary rõ.

### Main Task
Tạo HTTP server gọi vào core service đã tách ở Extension 09.

### Acceptance Criteria
- [ ] Có API server entrypoint riêng, ví dụ `cmd/mgit-api`
- [ ] Handler không import package CLI
- [ ] Handler chỉ decode request, gọi service, encode response
- [ ] JSON error response nhất quán: code, message, details nếu cần
- [ ] Có route health check
- [ ] Có endpoint tạo/open repo trong data root được cấu hình
- [ ] Có endpoint status, log, commit
- [ ] Endpoint `add` nhận file content và relative path trong request body; server ghi content vào repo workspace tương ứng với repoID — client không gửi raw server path
- [ ] Có endpoint branch/switch/checkout ở mức cơ bản
- [ ] Có endpoint diff/merge nếu core đã pass Phase 07
- [ ] Có middleware logging, recover, request id
- [ ] Có context timeout cho request
- [ ] Có `httptest` cho happy path và error path
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "HTTP không phải nơi viết Git logic. Handler mà tự đọc `.mgit/objects` là hỏng boundary."
>
> "Tôi cần API contract ổn định: request DTO, response DTO, error response. Sau này Flutter client hoặc web UI gọi được mà không cần biết CLI."
>
> "Điểm học chính: backend Go là xử lý boundary, context, error, test, lifecycle. Không phải cứ có `http.ListenAndServe` là backend."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/api/server.go
package api

type Server struct{}
type Config struct{}

func NewServer(cfg Config) (*Server, error)
func (s *Server) Routes() http.Handler

// File: internal/api/dto.go
package api

type ErrorResponse struct{}
type CreateRepoRequest struct{}
type CreateRepoResponse struct{}
type StatusResponse struct{}
type CommitRequest struct{}
type CommitResponse struct{}

// TODO-EXT10-A: Handler gọi core service, không đọc `.mgit` trực tiếp.
// SENIOR ASKS: Nếu handler biết object file path, boundary đang vỡ ở đâu?

// TODO-EXT10-B: Map domain error sang HTTP status.
// SENIOR ASKS: `branch not found`, `dirty worktree`, `conflict` nên map status khác nhau thế nào?

// TODO-EXT10-C: Mọi request phải đi qua context.
// SENIOR ASKS: Nếu client disconnect giữa lúc merge, handler và service phản ứng ra sao?
```

#### Suggested API Contract
```text
GET    /healthz
POST   /repos
GET    /repos/{repoID}/status
POST   /repos/{repoID}/add
POST   /repos/{repoID}/commits
GET    /repos/{repoID}/log
GET    /repos/{repoID}/branches
POST   /repos/{repoID}/branches
POST   /repos/{repoID}/switch
POST   /repos/{repoID}/checkout
GET    /repos/{repoID}/diff
POST   /repos/{repoID}/merge
```

#### Theory Notes
- [ ] HTTP handler là adapter, không phải domain — không viết Git logic bên trong handler
- [ ] `net/http` server spawn **một goroutine riêng** cho mỗi request — handler chạy concurrently, service layer phải thread-safe
- [ ] `r.Context()` trả context cancel khi client ngắt kết nối hoặc request timeout — dùng nó, không tự tạo `context.Background()`
- [ ] Pattern handler timeout: `ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()` rồi pass ctx xuống service
- [ ] Không gọi blocking operation trong handler mà không pass context — nếu service ignore ctx, goroutine có thể chạy sau khi client đã disconnect
- [ ] `httptest` giúsp test handler không cần mở port thật — `httptest.NewRecorder()` + `httptest.NewRequest()` là đủ

#### Socratic Questions
1. Vì sao API không nên expose raw filesystem path từ server?
2. Khi core trả conflict, response nên là 409 hay 500? Vì sao?
3. Handler validate gì, core validate gì?
4. Timeout request có đảm bảo không leak goroutine không?

### Output Checklist: Làm sao biết mình xong?
- [ ] `mgit-api` chạy local
- [ ] `GET /healthz` trả OK
- [ ] Có thể tạo repo và gọi status qua HTTP
- [ ] Có thể add/commit/log qua HTTP với core thật
- [ ] Error response nhất quán giữa các endpoint
- [ ] API docs hoặc README có curl examples

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] `httptest` cho health check
- [ ] `httptest` cho create repo
- [ ] `httptest` cho repo không tồn tại
- [ ] `httptest` cho invalid JSON
- [ ] `httptest` cho domain conflict/dirty worktree nếu core hỗ trợ
- [ ] `go test -race ./...`

### Learning Notes / Docs
- [ ] Viết `docs/api.md`
- [ ] Thêm curl demo vào README
- [ ] Ghi rõ limitation: API dùng local filesystem, chưa phải distributed Git server

### Retrospective: Sau khi xong, hãy tự hỏi
1. Handler nào đang phình vì domain logic lẫn vào?
2. Error response nào chưa đủ semantic cho client?
3. Endpoint nào nên đổi tên để bám resource REST hơn?

---

## Extension Checkpoints (BẮT BUỘC)

### CP-EXT10-A: HTTP Contract Gate
- [ ] Route list rõ.
- [ ] Request/response DTO rõ.
- [ ] Error response chuẩn hóa.

### CP-EXT10-B: Test Gate
- [ ] Handler test không phụ thuộc port thật.
- [ ] Error path được test.
- [ ] Race test không phát hiện shared state nguy hiểm.

### CP-EXT10-C: Oral Defense
- [ ] Giải thích được handler, middleware, service, repository khác nhau ở đâu.
- [ ] So sánh được `context.Context` với cancellation pattern trong Dart/Flutter.

## Failure Modes (PHẢI BIẾT)
- Handler tự đọc/ghi `.mgit`.
- Trả toàn bộ Go error string thô cho client.
- Mọi lỗi đều trả 500.
- Không test invalid JSON và not found.

## Progression Rules

### Rule 1: API không được bypass core.
Nếu CLI và API cho kết quả khác nhau, boundary sai hoặc test thiếu.

### Rule 2: Mỗi endpoint phải có error contract.
Không có error contract thì Flutter client sau này không xử lý UX tử tế được.

### Rule 3: Không deploy khi chưa có `httptest`.
Server chạy được trên máy bạn không chứng minh được backend đúng.

## Tổng Kết

### Deliverables
- [ ] `cmd/mgit-api` hoặc API entrypoint tương đương.
- [ ] REST routes gọi core service.
- [ ] API tests với `httptest`.
- [ ] `docs/api.md` có curl examples.

### First-Principles Question
HTTP handler cần biết những gì để xử lý `POST /repos/{id}/commits`, và những gì handler tuyệt đối không nên biết về Git object database?
