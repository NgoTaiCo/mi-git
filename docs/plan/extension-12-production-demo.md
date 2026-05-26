# Extension 12: Auth, Observability, Docker & Deploy Demo

> **Meta:** Production-ish backend capstone. Phase này không biến Mini Git thành GitHub; nó làm API đủ sạch để demo với nhà tuyển dụng.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Backend Extension Sprint: Production-ish API

> **Mục tiêu extension:** Thêm auth cơ bản, middleware, structured logging, config, Docker, deploy demo và runbook để chứng minh bạn hiểu backend vận hành.
>
> **Nguồn:** Sau `extension-11-persistence-repository.md`
>
> **Mini-git surface:** JWT/session boundary, middleware, logs, config, Dockerfile, docker-compose, deploy docs, API demo

---

## Extension Overview

### Mission
- Extension 12 - Production-ish backend demo cho Mini Git API

### Flutter / Dart Bridge
> Flutter app production cần flavor config, auth token, error UX, logging. Backend Go cũng vậy: config qua env, middleware, auth, request id, structured log, graceful shutdown. Chạy được local chưa đủ; phải vận hành được và debug được.

### Go Skills Required For This Extension
> `os.Getenv`, `os.LookupEnv`, manual env parse hoặc `envconfig` library. `log/slog`: `slog.New`, `slog.With`, structured key/value logging, `slog.Error`. Signal: `os.Signal`, `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`. Graceful shutdown: `http.Server.Shutdown(ctx)`. JWT: `golang-jwt/jwt` HMAC sign/verify, custom claims struct. Dockerfile: multi-stage `FROM golang:1.22-alpine AS builder` → `FROM alpine:latest`, copy binary only. `docker-compose.yml`: services api + postgres, health check, env file.

---

## Mission: Public Backend Demo Sẵn Sàng Cho CV

### User Story
> Nhà tuyển dụng hỏi: *"Nếu đưa API này lên môi trường thật, em handle auth, config, log, deploy và debug lỗi thế nào?"*
>
> Câu trả lời đúng: Có middleware, auth boundary, config rõ, Docker, runbook, deploy notes, health check và demo script.

### Main Task
Hardening backend API để public repo nhìn như một backend project có thể review nghiêm túc.

### Acceptance Criteria
- [ ] Có config qua environment variables
- [ ] Có graceful shutdown
- [ ] Có middleware request id, logging, recover, CORS nếu cần
- [ ] Có auth boundary cho repo owner: JWT hoặc token đơn giản có ghi rõ limitation
- [ ] Không log secret/token
- [ ] Có structured logging
- [ ] Có health check và readiness nếu dùng DB
- [ ] Có Dockerfile multi-stage
- [ ] Có `docker-compose` cho API + database
- [ ] Có migration command hoặc migration docs
- [ ] Có deploy guide Fly.io / Render / Railway / Cloud Run
- [ ] Có API demo script hoặc curl collection
- [ ] README ghi rõ scope, architecture, limitations, demo
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Production-ish không có nghĩa là nhồi Kubernetes. Nó nghĩa là service biết config, log, shutdown, auth, test và deploy tối thiểu."
>
> "Tôi cần thấy bạn biết cắt scope: Mini Git API không phải GitHub clone. Nó là backend wrapper quanh Git-like core để chứng minh Go backend fundamentals."
>
> "Điểm học chính: vận hành backend là quản lý failure mode, không phải chỉ trả JSON happy path."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/config/config.go
package config

type Config struct{}

func Load() (Config, error)

// File: internal/auth/auth.go
package auth

type Principal struct{}
type Verifier interface {
	Verify(token string) (Principal, error)
}

// File: internal/api/middleware.go
package api

type Middleware func(http.Handler) http.Handler

func RequestID(next http.Handler) http.Handler
func Recover(next http.Handler) http.Handler
func Auth(next http.Handler) http.Handler

// TODO-EXT12-A: Config không hard-code secret, port, database URL.
// SENIOR ASKS: Vì sao hard-code config làm deploy và test đều bẩn?

// TODO-EXT12-B: Middleware không được nuốt lỗi domain.
// SENIOR ASKS: Recover middleware xử lý panic khác gì error business bình thường?

// TODO-EXT12-C: Demo phải chạy từ môi trường sạch.
// SENIOR ASKS: Demo public chứng minh gì nếu nó chỉ chạy trên máy local đã setup sẵn?
```

#### Theory Notes
- [ ] Ghi chú: JWT/token là authentication, không thay thế authorization
- [ ] Ghi chú: structured logging khác `fmt.Println` ở đâu
- [ ] Ghi chú: graceful shutdown bảo vệ request đang chạy như thế nào
- [ ] Ghi chú: Docker multi-stage giảm artifact production ra sao

#### Socratic Questions
1. Recover middleware có nên biến mọi panic thành 500 rồi bỏ qua không?
2. Request id giúp debug flow qua middleware/service/repository như thế nào?
3. Secret nên lấy từ đâu khi chạy local, Docker, deploy?
4. Readiness check khác health check ở đâu nếu database down?

### Output Checklist: Làm sao biết mình xong?
- [ ] API chạy bằng `go run`
- [ ] API chạy bằng Docker
- [ ] API + DB chạy bằng `docker-compose`
- [ ] Health/readiness hoạt động
- [ ] Auth chặn request thiếu token ở endpoint cần bảo vệ
- [ ] Demo script tạo repo, commit, log qua HTTP
- [ ] README có kiến trúc, cách chạy, API examples, deploy notes, limitation

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Middleware request id
- [ ] Middleware auth success/failure
- [ ] Invalid/missing token
- [ ] Config thiếu env bắt buộc
- [ ] Readiness khi DB unavailable
- [ ] Graceful shutdown nếu có test được ở mức package
- [ ] `go test -race ./...`

### Learning Notes / Docs
- [ ] Viết `docs/deployment.md`
- [ ] Viết `docs/runbook.md`
- [ ] Cập nhật `docs/what-i-learned-about-backend-go.md`
- [ ] Cập nhật README phần CV story: Git internals + Go backend API

### Retrospective: Sau khi xong, hãy tự hỏi
1. Limitation nào phải nói thẳng thay vì giấu?
2. Log hiện tại có đủ để debug request lỗi không?
3. Auth hiện tại chứng minh backend thinking hay chỉ là token check tượng trưng?

---

## Extension Checkpoints (BẮT BUỘC)

### CP-EXT12-A: Runtime Gate
- [ ] Config qua env.
- [ ] Graceful shutdown.
- [ ] Health/readiness.

### CP-EXT12-B: Security/Observability Gate
- [ ] Auth boundary rõ.
- [ ] Không log secret.
- [ ] Request log có request id/status/duration.

### CP-EXT12-C: Public Demo Gate
- [ ] Docker/compose chạy được từ repo sạch.
- [ ] Demo HTTP flow chạy được.
- [ ] README đủ để reviewer chạy theo.

## Failure Modes (PHẢI BIẾT)
- Hard-code secret, DB URL, port.
- Middleware auth chỉ tồn tại cho có nhưng không test.
- Dockerfile build được nhưng container không chạy do path/config sai.
- Deploy docs nói chung chung, không có command cụ thể.
- README claim quá đà như "Git server" trong khi chưa có remote protocol.

## Progression Rules

### Rule 1: Production-ish nghĩa là debug được.
Nếu request lỗi mà log không giúp truy ra repo/user/endpoint/status, chưa đạt.

### Rule 2: Security phải có limitation rõ.
Token demo được, nhưng phải ghi rõ chưa có full user management nếu chưa làm.

### Rule 3: Public demo phải reproducible.
Reviewer clone repo, chạy docs, thấy được kết quả.

## Tổng Kết

### Deliverables
- [ ] Config/env setup.
- [ ] Middleware auth/logging/recover.
- [ ] Dockerfile + docker-compose.
- [ ] Deployment docs + runbook.
- [ ] HTTP demo script.
- [ ] Public README story.

### First-Principles Question
Một backend API "production-ish" khác gì với một HTTP server chỉ chạy được happy path trên máy local của bạn?
