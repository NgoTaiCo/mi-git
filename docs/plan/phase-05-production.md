# Phase 5: Production & Systems Programming (Tuần 9-10)

> **Context:** Ban đã hoc xong 4 phase dau. Bay gio khong con la hoc syntax nua — day la luc code cua ban chay trong production, luc 3AM co nguoi goi dien bao "server crash", va log chi ghi mot dong "error". Phase nay day cho ban viet code song sot duoc trong moi truong thuc te.
>
> **Goal:** Viet code that survives 3AM pages. Observability by design. Graceful everything.

---

## Topic 05.1: Error Handling (Production)

### User Story

> Khach hang (Product Owner) noi: "Toi 3AM, server crash. Log chi ghi 'error'. Toi khong biet loi tu dau ra. Can error trace"

**Context:** Ban dang on-call. He thong dang chay production thi dot nhien tra ve 500 cho 1% request. Log chi co 1 dong `error` — khong stack trace, khong context, khong biet loi tu dau. Ban mat 4 tieng de debug chi vi error khong duoc wrapped dung cach.

### Acceptance Criteria

- [ ] Error wrapping voi `%w` — moi error phai giu nguyen chain tu goc
- [ ] `errors.Is` — kiem tra loi cu the trong chain (vi du: kiem tra co phai `sql.ErrNoRows` khong)
- [ ] `errors.As` — trich xuat custom error type tu chain (vi du: lay `*ValidationError` de lay field details)
- [ ] Custom error types — dinh nghia error co them metadata (field, code, HTTP status)
- [ ] Sentinel errors — error constants dung de so sanh

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Hoi toi on-call o project fintech, co 1 loi khong wrapped suyt nua lam toi mat 4 tieng. Mot goroutine bi fail trong worker pool, loi chi ghi `process failed: error`. Toi phai grep ca codebase moi tim ra goc re cua van de. Tu ngay do, toi co quy tac sat: MOI loi phai wrapped voi context, MOI function phai them thong tin vao chain."
>
> "Van de cot loi o day la: Go khong co stack trace tu dong nhu Java hay Dart. `error` trong Go chi la 1 interface voi method `Error() string`. Neu ban khong chu dong them context, ban se mat trace. Dieu nay khong phai bug cua Go — no la design decision. Go yeu cau ban co y thuc ve error handling."
>
> "Toi se phan ra thanh cac buoc: 1) Hieu error interface va cach wrapping, 2) Dung errors.Is/As dung cach, 3) Thiet ke custom error types cho domain cua minh, 4) Thiet lap sentinel errors cho cac loi thuong gap."

---

### TODO Comments (Code Skeleton)

```go
package errs

import "errors"
import "fmt"

// TODO-[1]: Dinh nghia sentinel errors cho cac loi thuong gap
// SENIOR ASKS: Tai sao dung var thay const cho error? Error co the la const khong?
// HINT: errors.New tra ve pointer — pointer khong the la const

var (
	// TODO: Them sentinel errors cho cac loi: not found, validation, conflict, unauthorized
	// Vi du: ErrNotFound = errors.New("resource not found")
)

// TODO-[2]: Dinh nghia custom error type voi them metadata
// SENIOR ASKS: Khi nao ban can custom type thay vi chi dung fmt.Errorf?
// HINT: Think ve viec ban can trich xuat them thong tin tu error (field, code, status)

type AppError struct {
	// TODO: Them cac field: Code, Message, Field, HTTPStatus
	// SENIOR ASKS: Tai sao struct nay khong co field `Err error`?
	// HINT: Neu co field Err, ban co the wrap error khac ben trong — dieu nay quan trong
}

func (e *AppError) Error() string {
	// TODO: Tra ve string format co day du thong tin
	// SENIOR ASKS: Neu AppError wrap 1 error khac, Error() string nen tra ve gi?
	// HINT: Ban muon thay ca chain hay chi thong tin cua tang nay?
}

// TODO-[3]: Implement Unwrap de errors.Is/As hoat dong
// SENIOR ASKS: errors.Is hoat dong nhu the nao ben duoi? Tai sao can Unwrap?
// HINT: errors.Is doc theo chain qua tung cap Unwrap

func (e *AppError) Unwrap() error {
	// TODO: Tra ve wrapped error
}

// TODO-[4]: Helper functions de tao error nhanh
// SENIOR ASKS: Tai sao khong tao AppError truc tiep ma can helper?
// HINT: Duplication — ban se viet fmt.Errorf(...) o 50 noi khac nhau

func NewNotFound(resource string, id any) error {
	// TODO: Tra ve AppError voi code "NOT_FOUND", HTTP 404
	// Wrap tu ErrNotFound sentinel
}

func NewValidation(field, message string) error {
	// TODO: Tra ve AppError voi code "VALIDATION_ERROR", HTTP 400
}

// TODO-[5]: Ham kiem tra loi cu the trong chain
// SENIOR ASKS: Tai sao khong dung err == ErrNotFound truc tiep?
// HINT: Neu error da duoc wrap, `==` se tra ve false

func IsNotFound(err error) bool {
	// TODO: Dung errors.Is de kiem tra trong ca chain
}

func IsValidation(err error) bool {
	// TODO: Tuong tu cho validation
}

// TODO-[6]: Ham trich xuat AppError tu error bat ky
// SENIOR ASKS: errors.As khac errors.Is cho nao?
// HINT: Is kiem tra "co phai loi nay khong", As la "lay thong tin tu loi nay"

func AsAppError(err error) (*AppError, bool) {
	// TODO: Dung errors.As de trich xuat *AppError
}

// TODO-[7]: Wrap error voi them context
// SENIOR ASKS: Khi nao dung %w, khi nao dung %v trong fmt.Errorf?
// HINT: %w giu error chain — chi dung khi ban muon caller co the unwrap. %v lam mat chain.

func Wrap(err error, context string) error {
	// TODO: Wrap error voi fmt.Errorf("...: %w", ..., err)
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. Tai sao Go khong co exception (try/catch) nhu Dart/Java/Python? Dieu nay buoc ban phai lam gi khac?
2. Khi nao ban dung `%v` thay vi `%w` trong `fmt.Errorf`? Co truong hop nao muon "lam mat" error chain khong?
3. Tai sao `errors.As` can pointer (`&target`) thay vi value? Dieu nay noi len dieu gi ve Go type system?
4. Neu ban co 1 goroutine bi panic, ban co nen recover va wrap thanh error khong? Tai sao?
5. Quy tac cua ban: bao nhieu cap wrap la qua nhieu? Khi nao nen log thay vi wrap?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Co it nhat 4 sentinel errors: ErrNotFound, ErrValidation, ErrConflict, ErrUnauthorized
- [ ] TODO-[2] hoan thanh: AppError struct co Code, Message, Field, HTTPStatus, va Err (wrapped)
- [ ] TODO-[3] hoan thanh: Unwrap() duoc implement, errors.Is hoat dong voi AppError
- [ ] TODO-[4] hoan thanh: Co helper functions NewNotFound, NewValidation, NewConflict
- [ ] TODO-[5] hoan thanh: IsNotFound va IsValidation hoat dong voi wrapped errors
- [ ] TODO-[6] hoan thanh: AsAppError trich xuat duoc tu error chain
- [ ] TODO-[7] hoan thanh: Hieu ro khi nao dung %w vs %v

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: `errors.Is(wrappedErr, ErrNotFound)` tra ve true — vi sao case nay quan trong?
  - **Giai thich:** Day la use case pho bien nhat — ban wrap loi DB thanh loi domain, nhung van muon kiem tra goc
- [ ] Test case: `errors.As(chainErr, &appErr)` tra ve true va appErr co du thong tin — boundary case gi co the fail?
  - **Giai thich:** Neu Unwrap() khong implement dung, As se khong tim thay AppError trong chain
- [ ] Test case: Error() string format phai chua ca wrapped error message — vi sao?
  - **Giai thich:** Log chi doc Error() string — neu khong co du thong tin, debug se kho

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Custom error type nhieu field co lam code cham hon khong? Co dang khong? Khi nao nen dung `fmt.Errorf` don gian thay vi struct?
2. **Neu requirement thay doi:** Neu sau nay can i18n (dich error message sang tieng Viet, Anh, Nhat), thiet ke cua ban co ho tro khong? Can refactor gi?
3. **Architecture decision:** Tai sao toi dat error package o `errs` thay vi `errors`? Co van de gi voi ten `errors`?

---
---

## Topic 05.2: Structured Logging (slog)

### User Story

> Khach hang (Product Owner) noi: "Log phai de Elasticsearch parse duoc. fmt.Printf khong co structure, khong filter duoc"

**Context:** He thong cua ban chay 20 pods tren Kubernetes. Mot pod bi loi nhung ban khong biet pod nao. Log duoc thu thap boi ELK stack (Elasticsearch + Logstash + Kibana) nhung moi log entry la 1 string vo dinh dang kieu `fmt.Printf("User %s login failed from %s", username, ip)`. Khong the filter theo `level=error`, khong the group theo `user_id`, khong the alert khi `http_status=500` vuot nguong.

### Acceptance Criteria

- [ ] Dung `log/slog` (Go 1.21+) thay vi `fmt.Print`/`log.Print`
- [ ] JSON handler de log machine-parseable
- [ ] Log levels: Debug, Info, Warn, Error — dung dung cho tung loai thong tin
- [ ] Attributes: key-value pairs gan voi moi log entry (user_id, request_id, duration)
- [ ] Context: truyen logger qua `context.Context` de moi function cung request co cung attributes

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Toi tung review code co 47 fmt.Println. Khong bao gio. Trong production, log khong phai de doc khi dev — log la du lieu de may parse, filter, va alert. Mot dong log `User login failed` khong co giup gi khi ban can biet: user nao? Luc nao? Tu IP nao? Request ID nao de trace?
>
> "Chuyen tu fmt.Printf sang slog khong chi la thay function — no la thay doi tu duy. Moi log entry phai co structure: {timestamp, level, message, attrs...}. Message la cai con nguoi doc, attrs la cai may parse.
>
> "Toi se phan ra: 1) Setup slog voi JSON handler, 2) Dung levels dung cho, 3) Them attributes co y nghia, 4] Truyen logger qua context de khong phai truyen tham so o moi function."

---

### TODO Comments (Code Skeleton)

```go
package logger

import (
	"context"
	"log/slog"
	"os"
)

// TODO-[1]: Khoi tao slog voi JSON handler
// SENIOR ASKS: Tai sao JSON handler ma khong phai Text handler trong production?
// HINT: Elasticsearch, Datadog, CloudWatch deu parse JSON de hon plain text

func Init(env string) *slog.Logger {
	// TODO: Neu env == "development", dung TextHandler de doc de
	// TODO: Neu env == "production", dung JSONHandler
	// SENIOR ASKS: Source nao nen quyet dinh env? Flag? Env var?
	// HINT: Environment variable — khong hardcode trong code

	var handler slog.Handler
	// TODO: handler = slog.NewJSONHandler(os.Stdout, opts)
	// opts co them AddSource: true de log ca file:line

	return slog.New(handler)
}

// TODO-[2]: Dung levels dung cho — moi level co y nghia rieng
// SENIOR ASKS: Khi nao dung Debug? Khi nao dung Info? Warn? Error?
// HINT: Debug = chi dev can; Info = flow binh thuong; Warn = bat thuong nhung khong loi; Error = loi can xu ly

// Quy tac levels:
// Debug: SQL queries, request/response body, cache hit/miss
// Info:  Request bat dau/ket thuc, startup, config loaded
// Warn:  Retry, timeout, deprecated API duoc goi
// Error: DB connection fail, validation hard-fail, goroutine panic

// TODO-[3]: Log attributes co y nghia — khong log password, token, PII
// SENIOR ASKS: Cai gi KHONG nen log? Tai sao?
// HINT: GDPR, security audit — password, token, credit card, SSN

// Attributes nen co cho moi HTTP request:
// - request_id (de trace xuyen suot he thong)
// - user_id (de biet ai gay ra)
// - method, path (endpoint nao)
// - status_code (ket qua)
// - duration_ms (performance)
// - error (neu co)

// TODO-[4]: Middleware log HTTP request/response
// SENIOR ASKS: Tai sao middleware la noi tot nhat de log request?
// HINT: No bao quanh handler — co du thong tin truoc va sau khi xu ly

type responseWriter struct {
	// TODO: Wrap http.ResponseWriter de bat status code
	// SENIOR ASKS: Tai sao can wrap? http.ResponseWriter khong expose status?
	// HINT: WriteHeader chi goi 1 lan — minh can ghi lai de log

	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	// TODO: Luu status va goi ResponseWriter.WriteHeader
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// TODO: Luu size va goi ResponseWriter.Write
}

// TODO-[5]: Request logging middleware
// SENIOR ASKS: Log luc bat dau request hay luc ket thuc? Hay ca hai?
// HINT: Chi log luc ket thuc — luc do moi co duration va status code

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Bat dau timer, tao request_id, gan vao context
			// TODO: Goi next.ServeHTTP voi responseWriter wrap
			// TODO: Sau khi xong, log voi duration, status, method, path
		})
	}
}

// TODO-[6]: Lay logger tu context — khong truyen tham so moi function
// SENIOR ASKS: Tai sao dung context thay vi truyen *slog.Logger qua moi ham?
// HINT: Context la request-scoped — no di theo request xuyen suot call stack

type ctxKey struct{}

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	// TODO: Lay logger tu context, neu khong co thi tra ve default
	// SENIOR ASKS: Nen tra ve nil hay default logger khi khong tim thay?
	// HINT: Default (slog.Default()) — de code khong panic khi quen gan
}

// TODO-[7]: Add request-scoped attributes
// SENIOR ASKS: Tai sao khong dung slog.With o global level?
// HINT: slog.With tra ve logger moi — attributes la immutable, thread-safe

func WithRequestID(ctx context.Context, reqID string) context.Context {
	log := FromContext(ctx)
	newLog := log.With("request_id", reqID)
	return WithLogger(ctx, newLog)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	// TODO: Tuong tu WithRequestID nhung voi user_id
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. fmt.Printf va slog.Info khac nhau ve performance? Co dang de doi khong?
2. Tai sao log/slog duoc dua vao standard library (Go 1.21) thay vi de third-party? Dieu nay noi len gi?
3. Context value dung interface{} — dieu nay co van de gi? Tai sao van dung duoc?
4. Neu ban co 1 microservice goi 3 service khac, lam sao de request_id di xuyen suot ca 4 service?
5. Khi nao ban nen log voi level Debug nhung khong commit vao production? Cach nao de "tat" debug log ma khong sua code?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Co Init(env) tra ve logger phu hop voi environment
- [ ] TODO-[2] hoan thanh: Hieu va ap dung dung 4 levels (Debug/Info/Warn/Error)
- [ ] TODO-[3] hoan thanh: Biet cai gi khong nen log (PII, password, token)
- [ ] TODO-[4] hoan thanh: responseWriter wrap de bat status code va size
- [ ] TODO-[5] hoan thanh: Middleware log request voi duration, status, method, path
- [ ] TODO-[6] hoan thanh: Logger lay tu context, khong truyen tham so
- [ ] TODO-[7] hoan thanh: Request-scoped attributes (request_id, user_id) duoc gan dung

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Log entry la valid JSON (co the parse bang encoding/json) — vi sao case nay quan trong?
  - **Giai thich:** Neu JSON khong valid, ELK stack se khong parse duoc — log mat
- [ ] Test case: Logger tu context co dung attributes — boundary case gi co the fail?
  - **Giai thich:** Quen gan context, hoac goroutine moi khong ke thua context
- [ ] Test case: Khong log password trong request body — vi sao?
  - **Giai thich:** Security audit — 1 log chua password co the lam ca cong ty mat viec

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** JSON log lon hon plain text ~30%. Co dang khong? Khi nao ban chon plain text?
2. **Neu requirement thay doi:** Neu can log ra nhieu output cung luc (file + stdout + remote), thiet ke cua ban co ho tro khong?
3. **Architecture decision:** Tai sao toi khong dung global `slog.Default()` ma truyen logger qua context? Global co van de gi?

---
---

## Topic 05.3: Graceful Shutdown

### User Story

> Khach hang (Product Owner) noi: "K8s kill pod de deploy version moi. Request dang xu ly bi cat ngang, user thay 502"

**Context:** Ban dang deploy version moi cua service. Kubernetes gui SIGTERM den pod cu, doi 30 giay, roi gui SIGKILL neu pod chua dung. Pod cua ban nhan SIGTERM va lap tuc exit — nhung van con 15 request dang xu ly. Nhung request do bi cat ngang giua chung, user thay 502 Bad Gateway, data co the bi inconsistent.

### Acceptance Criteria

- [ ] Bat tin hieu SIGTERM, SIGINT bang `signal.NotifyContext`
- [ ] Goi `server.Shutdown` voi timeout (khong phai `server.Close`)
- [ ] Cho in-flight requests hoan thanh truoc khi exit
- [ ] Timeout ro rang — khong cho mai mai
- [ ] Health check tra ve unhealthy khi bat dau shutdown de load balancer ngung gui traffic moi

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Hoi toi o project e-commerce, Black Friday, deploy hotfix. Kubernetes kill pod cu, pod moi chua san sang, va 200 request bi 502. Doanh thu mat ~$15k trong 3 phut. Nguyen nhan: code khong handle SIGTERM, pod exit ngay lap tuc.
>
> "Van de cot loi: Go `http.Server` co 2 cach dung — `Close()` va `Shutdown()`. Close() dong ngay lap tuc, cac connection bi dut. Shutdown() doi connections hien tai hoan thanh. Ban muon Shutdown().
>
> "Flow dung: 1) Nhan SIGTERM -> 2) Health check tra ve 503 -> 3) Goi server.Shutdown voi timeout -> 4) Cho active requests xong -> 5) Cleanup resources -> 6) Exit."

---

### TODO Comments (Code Skeleton)

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO-[1]: Tao signal-aware context
// SENIOR ASKS: Tai sao dung signal.NotifyContext thay vi signal.Notify truyen qua channel?
// HINT: NotifyContext ket hop context cancellation + signal handling — sach hon

func setupSignalContext() context.Context {
	// TODO: Tao context bi cancel khi nhan SIGTERM hoac SIGINT
	// SENIOR ASKS: Tai sao lai bat ca SIGTERM lan SIGINT?
	// HINT: SIGINT = Ctrl+C (local dev); SIGTERM = Kubernetes/docker stop
}

// TODO-[2]: Health check handler biet trang thai shutdown
// SENIOR ASKS: Tai sao health check can biet dang shutdown?
// HINT: Load balancer chi ngung gui traffic khi health check fail

type healthHandler struct {
	// TODO: Them truong de biet dang shutdown
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Neu dang shutdown -> tra ve 503 Service Unavailable
	// TODO: Neu binh thuong -> tra ve 200 OK
}

func (h *healthHandler) SetShuttingDown() {
	// TODO: Dat trang thai shutdown
}

// TODO-[3]: Ham start server voi graceful shutdown
// SENIOR ASKS: Tai sao chay server trong goroutine rieng?
// HINT: ListenAndServe block — can chay async de main thread doi signal

func runServer(addr string, handler http.Handler, shutdownTimeout time.Duration) error {
	// TODO: Tao http.Server
	// TODO: Goroutine chay server.ListenAndServe()
	// TODO: Main thread doi signal context bi cancel
	// TODO: Khi nhan signal: SetShuttingDown, goi server.Shutdown voi timeout
	// TODO: Neu Shutdown timeout -> force close

	return nil
}

// TODO-[4]: Shutdown sequence dung thu tu
// SENIOR ASKS: Thu tu cac buoc co quan trong khong? Co the dao thu tu khong?
// HINT: Health check phai fail TRUOC khi shutdown — neu sau, traffic moi van vao

func gracefulShutdown(
	server *http.Server,
	health *healthHandler,
	timeout time.Duration,
) error {
	// Buoc 1: Bao hieu unhealthy
	// Buoc 2: Cho 1 chut de load balancer nhan tin (grace period)
	// Buoc 3: Shutdown server voi timeout
	// Buoc 4: Cleanup (close DB connections, v.v.)
	// Buoc 5: Return

	// TODO: Implement ting buoc
	// SENIOR ASKS: "Cho 1 chut" la bao lau? Con so nay lay tu dau?
	// HINT: Kubernetes terminationGracePeriodSeconds — thuong 30s. Grace period ~5s la du.
}

// TODO-[5]: Long-running request handler
// SENIOR ASKS: Request handler co can thay doi gi de ho tro graceful shutdown?
// HINT: Handler nen kiem tra context.Done() va tra ve early neu bi cancel

func longHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-time.After(10 * time.Second):
		// TODO: Xu ly thanh cong
		w.Write([]byte("done"))
	case <-ctx.Done():
		// TODO: Client hoac server shutdown — tra ve 503
		// SENIOR ASKS: Nen tra ve gi khi context done? Status code nao?
		// HINT: 503 Service Unavailable — khong phai 500 (loi server) hay 499 (client close)
	}
}

// TODO-[6]: Main function noi cac thanh phan
// SENIOR ASKS: Co can WaitGroup khong? Tai sao?
// HINT: Can doi goroutine server ket thuc truoc khi main exit

func main() {
	// TODO: Init logger, config
	// TODO: Setup signal context
	// TODO: Tao handler va health check
	// TODO: Run server voi graceful shutdown
	// TODO: Log khi shutdown hoan tat
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. `server.Close()` va `server.Shutdown()` khac nhau cho nao? Khi nao dung Close?
2. Kubernetes terminationGracePeriodSeconds la 30 giay. Dieu gi xay ra neu server.Shutdown cua ban can 45 giay?
3. Tai sao can "grace period" giua SetShuttingDown va Shutdown? Khong the shutdown ngay sao?
4. `r.Context()` trong HTTP handler co y nghia gi khi server shutdown? No bi cancel khi nao?
5. Neu ban co 1 background worker (goroutine) chay doc queue, no co nen dung khi nhan SIGTERM khong? Ngay lap tuc hay sau khi xong job hien tai?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Signal context duoc setup dung
- [ ] TODO-[2] hoan thanh: Health check tra ve 503 khi dang shutdown
- [ ] TODO-[3] hoan thanh: Server chay trong goroutine, main thread doi signal
- [ ] TODO-[4] hoan thanh: Shutdown sequence dung thu tu (health -> grace -> shutdown -> cleanup)
- [ ] TODO-[5] hoan thanh: Handler kiem tra r.Context().Done()
- [ ] TODO-[6] hoan thanh: Main function noi cac thanh phan, doi goroutine ket thuc

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Send SIGTERM -> server tra ve 503 health -> dung sau khi requests xong — vi sao case nay quan trong?
  - **Giai thich:** Day la flow day du — neu thieu buoc nao, graceful shutdown khong hoat dong
- [ ] Test case: Request dang xu ly khi nhan SIGTERM van hoan thanh — boundary case gi co the fail?
  - **Giai thich:** Neu timeout qua ngan, request bi cat ngang; neu qua dai, K8s gui SIGKILL
- [ ] Test case: Khong nhan traffic moi sau khi health check fail — vi sao?
  - **Giai thich:** Neu van nhan traffic, se co them requests can drain, co the khong kip truoc SIGKILL

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Grace period 5 giay co the lam deploy cham. Co nen bo qua khong? Risk la gi?
2. **Neu requirement thay doi:** Neu can zero-downtime deploy (khong request nao bi 503), thiet ke can thay doi gi? (Hint: rolling deploy + readiness probe)
3. **Architecture decision:** Tai sao khong dung `server.Close()` ngay? Co truong hop nao Close la dung khong?

---
---

## Topic 05.4: Observability

### User Story

> Khach hang (Product Owner) noi: "Production gap loi nhung khong biet bottleneck o dau. Can metrics va health check"

**Context:** Service cua ban bi cham luc gio cao diem. Ban khong biet vi sao — DB? CPU? Memory? Goroutine leak? Ban khong co visibility vao he thong dang chay. Ban can metrics de hieu performance, health checks de K8s biet khi nao restart pod, va alerts khi co van de.

### Acceptance Criteria

- [ ] Prometheus metrics: counter (requests total), histogram (request duration), gauge (active connections)
- [ ] HTTP middleware expose metrics tai `/metrics`
- [ ] Liveness probe: K8s biet pod co song khong (neu deadlock -> restart)
- [ ] Readiness probe: K8s biet pod san sang nhan traffic khong (neu DB down -> ngung traffic)
- [ ] Startup probe: K8s doi pod khoi dong xong truoc khi check liveness/readiness

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "'Khong biet bottleneck o dau' — cau noi nay toi nghe nhieu hon ca 'it works on my machine'. Trong production, ban khong the chay debugger. Ban can telemetry: metrics, logs, traces. Trong 3 cai do, metrics la cai re nhat trien khai nhung co impact lon nhat.
>
> "Prometheus la standard de facto. No pull metrics tu endpoint `/metrics` cua ban. Ban dung client library de dinh nghia counters, histograms, gauges — roi middleware tu dong update chung.
>
> "3 loai health check khac nhau: Liveness (pod co song khong?), Readiness (pod san sang chua?), Startup (pod khoi dong xong chua?). Nham lan 3 cai nay la nguyen nhan pho bien nhat cua incident o K8s."

---

### TODO Comments (Code Skeleton)

```go
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO-[1]: Dinh nghia cac metrics
// SENIOR ASKS: Tai sao dung promauto thay vi prometheus.NewCounter?
// HINT: promauto register voi default registry — it boilerplate hon

var (
	// Counter: dem so requests — khong giam bao gio
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Histogram: do thoi gian xu ly request
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Gauge: so requests dang xu ly (active)
	activeRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
	)
)

// TODO-[2]: Middleware update metrics
// SENIOR ASKS: Tai sao middleware la noi tot nhat de collect metrics?
// HINT: No bao quanh moi request — co day du thong tin truoc va sau

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Tang activeRequests
		// TODO: Bat dau timer
		// TODO: Goi next handler
		// TODO: Ghi lai status code
		// TODO: Giam activeRequests
		// TODO: Record requestDuration va requestsTotal
	})
}

// TODO-[3]: Expose /metrics endpoint
// SENIOR ASKS: Tai sao /metrics khong nen qua middleware logging va metrics?
// HINT: Khong muon log request den /metrics — no se spam logs va metrics

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// TODO-[4]: 3 loai health check
// SENIOR ASKS: Liveness va Readiness khac nhau cho nao? Cho vi du thuc te.
// HINT: Liveness = pod co song khong (deadlock -> restart); Readiness = pod san sang chua (DB down -> remove traffic)

// Liveness: don gian — chi can process khong deadlock
// Tra ve 200 neu server dang chay, 500 neu co deadlock/goroutine leak
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Don gian — chi tra ve 200
	// SENIOR ASKS: Co nen check DB connection trong liveness?
	// HINT: KHONG! Neu DB down, liveness fail -> K8s restart pod. Nhung DB down khong fix duoc bang restart!
}

// Readiness: check dependencies (DB, cache, external services)
// Tra ve 200 neu san sang nhan traffic, 503 neu chua san sang
func ReadinessHandler(checks map[string]HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check tat ca dependencies
		// TODO: Neu 1 cai fail -> tra ve 503 + chi tiet loi
		// SENIOR ASKS: Nen tra ve chi tiet loi (loai dependency nao fail) khong?
		// HINT: Co — giup ops biet nhanh van de o dau
	}
}

// HealthChecker interface cho cac dependency
type HealthChecker interface {
	Check(ctx context.Context) error
}

// TODO-[5]: DB health check implementation
// SENIOR ASKS: Nen check DB nhu the nao? SELECT 1? Hay ping?
// HINT: SELECT 1 — don gian, nhanh, du de biet connection con song

type DBChecker struct {
	DB *sql.DB
}

func (d *DBChecker) Check(ctx context.Context) error {
	// TODO: Thuc hien SELECT 1 voi timeout
	// SENIOR ASKS: Tai sao can timeout rieng cho health check?
	// HINT: Khong muon health check treo mai — K8s co timeout rieng
}

// TODO-[6]: Startup probe
// SENIOR ASKS: Startup probe khac liveness cho nao? Khi nao can?
// HINT: App khoi dong lau — khong muon K8s tuong no deadlock roi restart

func StartupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Tra ve 200 neu app khoi dong xong (DB connected, cache warmed, v.v.)
	}
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. Tai sao metrics dung pull model (Prometheus den lay) thay vi push model (app gui di)? Uu/nhuoc diem?
2. Histogram va Summary khac nhau cho nao? Khi nao dung cai nao?
3. Liveness probe fail -> K8s restart pod. Tai sao khong nen check external dependency trong liveness?
4. Neu readiness probe cua tat ca pods cung fail, dieu gi xay ra voi service? Giai phap?
5. Metrics co anh huong den performance khong? Can bo qua metrics de toi uu performance khong?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: 3 loai metrics duoc dinh nghia (counter, histogram, gauge)
- [ ] TODO-[2] hoan thanh: Middleware tu dong update metrics cho moi request
- [ ] TODO-[3] hoan thanh: `/metrics` endpoint expose duoc Prometheus scrape
- [ ] TODO-[4] hoan thanh: Liveness va Readiness handlers phan biet ro
- [ ] TODO-[5] hoan thanh: DB health check voi timeout
- [ ] TODO-[6] hoan thanh: Startup probe implementation

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: `/metrics` tra ve format Prometheus text — vi sao case nay quan trong?
  - **Giai thich:** Prometheus chi parse duoc format chuan — neu sai format, metrics khong duoc collect
- [ ] Test case: Readiness fail khi DB xuong — boundary case gi co the fail?
  - **Giai thich:** DB connection pool het -> readiness fail -> traffic bi chuyen di
- [ ] Test case: Counter khong bi giam khi request that bai — vi sao?
  - **Giai thich:** Counter chi tang — khong bao gio giam. Neu can do "hien tai", dung gauge.

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Prometheus client co overhead. Ban co nen tu viet metrics thay vi dung library khong? Risk la gi?
2. **Neu requirement thay doi:** Neu can distributed trace (trace request xuyen suot nhieu service), ban can them gi? (Hint: OpenTelemetry)
3. **Architecture decision:** Tai sao toi dung pull model thay vi push? Truong hop nao push la dung hon?

---
---

## Topic 05.5: Configuration

### User Story

> Khach hang (Product Owner) noi: "Config khac nhau o dev/staging/prod. Khong hardcode, validate khi startup"

**Context:** Developer moi clone repo, chay `go run .` va app crash vi thieu `DATABASE_URL`. Khong ai biet can set nhung env var nao. Tren production, app chay duoc 10 phut roi crash vi `TIMEOUT` la string "30s" nhung code doc nhu integer. Config nam rac roi khap noi: 1 phan trong code, 1 phan trong env, 1 phan trong file JSON.

### Acceptance Criteria

- [ ] Config doc tu env vars — khong hardcode secrets
- [ ] Config struct voi tag (duration, default values)
- [ ] Validate tat ca config khi startup — fail fast neu thieu/thuong
- [ ] Ho tro nhieu environment (dev/staging/prod) khong can sua code
- [ ] Khong commit secrets vao git — dung .env file cho local dev

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "12 Factor App principle #3: Store config in environment. Hoi toi review 1 PR co API key hardcode: `const apiKey = "sk-live-..."`. Commit len GitHub, 5 phut sau key bi revoke vi security scan cua GitHub phat hien. Tu do, toi co 0-tolerance voi secrets trong code.
>
> "Van de config khong chi la 'doc tu env' — no la ca quy trinh: 1) Dinh nghia struct ro rang, 2) Parse tu env voi type safety, 3) Validate khi startup, 4) Fail fast neu sai. Neu app chay duoc voi config sai, no se crash luc 2AM khi gap edge case.
>
> "Toi thuong dung approach: Config struct + `envconfig` hoac tu parse. Khong dung global var — truyen Config vao cac service can no."

---

### TODO Comments (Code Skeleton)

```go
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO-[1]: Dinh nghia Config struct
// SENIOR ASKS: Tai sao dung struct thay vi global vars hoac map[string]string?
// HINT: Type safety — compiler bao loi neu ban goi sai ten field; IDE autocomplete hoat dong

type Config struct {
	// Server
	Port         string        `env:"PORT" default:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"30s"`

	// Database
	DatabaseURL string `env:"DATABASE_URL" required:"true"`
	MaxDBConns  int    `env:"MAX_DB_CONNS" default:"10"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" default:"info"` // debug, info, warn, error

	// Security
	JWTSecret    string `env:"JWT_SECRET" required:"true"`
	BcryptCost   int    `env:"BCRYPT_COST" default:"10"`
	AllowedHosts string `env:"ALLOWED_HOSTS" default:"*"`

	// Environment
	Env string `env:"ENV" default:"development"` // development, staging, production
}

// TODO-[2]: Ham load config tu environment
// SENIOR ASKS: Tai sao khong dung library nhu caarlos0/env?
// HINT: Co the dung — nhung truoc het hieu cach no hoat dong bang cach tu viet

func Load() (*Config, error) {
	cfg := &Config{}

	// TODO: Parse tung field tu env var
	// SENIOR ASKS: Nen dung reflection de parse tu dong hay viet manual?
	// HINT: Reflection tien nhung cham, kho debug. Manual nhieu code nhung ro rang.

	// Parse server config
	cfg.Port = getEnv("PORT", "8080")
	rawReadTimeout := getEnv("READ_TIMEOUT", "30s")
	// TODO: Parse rawReadTimeout thanh time.Duration
	// SENIOR ASKS: time.ParseDuration ho tro nhung don vi gi?
	// HINT: "ns", "us", "ms", "s", "m", "h" — khong ho tro "d" (days)

	// Parse database config
	cfg.DatabaseURL = getEnv("DATABASE_URL", "")
	// TODO: Parse MAX_DB_CONNS thanh int, co default

	// TODO: Parse cac field con lai

	return cfg, nil
}

func getEnv(key, fallback string) string {
	// TODO: Doc tu os.Getenv, tra ve fallback neu khong co
}

// TODO-[3]: Validate config — fail fast
// SENIOR ASKS: Tai sao validate khi startup thay vi luc dung?
// HINT: Phat hien loi som — khong de app chay 1 luc roi crash vi config khong hop le

func (c *Config) Validate() error {
	var errs []string

	// TODO: Check required fields
	if c.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET is required")
	}

	// TODO: Check valid values
	// SENIOR ASKS: Nen check JWTSecret du manh khong?
	// HINT: Co — it nhat 32 chars cho HS256
	if len(c.JWTSecret) < 32 {
		errs = append(errs, "JWT_SECRET must be at least 32 characters")
	}

	// TODO: Check LogLevel hop le
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		errs = append(errs, fmt.Sprintf("LOG_LEVEL must be one of: debug, info, warn, error, got: %s", c.LogLevel))
	}

	// TODO: Check BcryptCost hop le (4-31)
	if c.BcryptCost < 4 || c.BcryptCost > 31 {
		errs = append(errs, fmt.Sprintf("BCRYPT_COST must be between 4 and 31, got: %d", c.BcryptCost))
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n- %s", strings.Join(errs, "\n- "))
	}
	return nil
}

// TODO-[4]: IsProduction helper
// SENIOR ASKS: Tai sao can helper nay? Khong the so sanh truc tiep?
// HINT: De doc + khong sai chinh ta. "production" vs "prod" vs "prd"

func (c *Config) IsProduction() bool {
	// TODO: Tra ve true neu Env == "production"
}

func (c *Config) IsDevelopment() bool {
	// TODO: Tra ve true neu Env == "development"
}

// TODO-[5]: Khong commit secrets — .env.example
// SENIOR ASKS: Tai sao .env vao .gitignore nhung .env.example thi khong?
// HINT: .env chua secrets thuc; .env.example chi chua ten key khong co value

// File: .env.example
// PORT=8080
// DATABASE_URL=postgresql://localhost:5432/mydb
// JWT_SECRET=your-secret-key-here
// LOG_LEVEL=info
// ENV=development

// TODO-[6]: Sample .env cho local dev (khong commit)
// File: .env (trong .gitignore)
// PORT=8080
// DATABASE_URL=postgresql://dev:dev@localhost:5432/devdb
// JWT_SECRET=dev-secret-not-for-production-32chars
// LOG_LEVEL=debug
// ENV=development
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. Tai sao 12 Factor App khuyen cao "store config in env" thay vi file config? File co uu diem gi khong?
2. `time.ParseDuration("30s")` va `time.ParseDuration("30")` khac nhau cho nao? Dieu nay anh huong gi den UX cua config?
3. Tai sao fail fast la tot? Co truong hop nao ban muon "co gang chay" du config thieu khong?
4. `os.Getenv("PORT")` tra ve string — con so? Ban co nen chuyen het sang int/float khong?
5. Neu 1 config chi dung tren production (vi du: S3 bucket), ban xu ly nhu the nao cho local dev?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Config struct co day du field voi type dung
- [ ] TODO-[2] hoan thanh: Ham Load() doc tu env voi default values
- [ ] TODO-[3] hoan thanh: Validate() kiem tra required + valid values
- [ ] TODO-[4] hoan thanh: IsProduction(), IsDevelopment() helpers
- [ ] TODO-[5] hoan thanh: .env.example khong chua secrets, duoc commit
- [ ] TODO-[6] hoan thanh: .env thuc duoc cho vao .gitignore

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Config thieu required field -> tra ve error — vi sao case nay quan trong?
  - **Giai thich:** Dieu nay ngan app chay voi config khong hoan chinh — crash som hon crash muon
- [ ] Test case: Config voi gia tri khong hop le -> tra ve error — boundary case gi co the fail?
  - **Giai thich:** Duration sai format ("30" thay vi "30s"), int parse fail ("abc" -> 0)
- [ ] Test case: Config day du -> tra ve *Config hop le — vi sao?
  - **Giai thich:** Dam bao happy path van hoat dong, default values duoc ap dung dung

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Env vars co gioi han (khong ho tro list, object). Khi nao ban can file config (YAML/JSON) thay vi env?
2. **Neu requirement thay doi:** Neu can hot-reload config khong restart app, thiet ke can thay doi gi?
3. **Architecture decision:** Tai sao toi validate trong Config.Validate() thay vi trong tung constructor cua service? Uu/nhuoc?

---
---

## Topic 05.6: Security

### User Story

> Khach hang (Product Owner) noi: "Audit tim thay password plain-text trong code. SQL injection qua search API. Can fix"

**Context:** Security audit vua xong va co 3 findings nghiem trong: 1) Password duoc luu plain-text trong database, 2) API search dung string concatenation thay vi parameterized query, cho phep SQL injection, 3) Khong co input validation, 1 API nhan user ID dang string va dung truc tiep vao query khong co check. Them vao do, frontend khong the goi API vi CORS policy block.

### Acceptance Criteria

- [ ] Password hashing voi bcrypt — khong bao gio luu plain-text
- [ ] SQL parameterized queries — khong bao gio concatenate user input vao SQL
- [ ] Input validation cho moi API endpoint — reject early, tra ve 400 voi chi tiet
- [ ] CORS configuration — chi cho phep trusted origins
- [ ] Khong expose internal details trong error response (stack trace, SQL query)

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Security audit findings la nhu toa do hieu suong — no cho biet noi nao yeu nhung khong phai la tat ca. Neu audit tim duoc 3 lo hong, thuc te co the co 10. Toi luon gia dinh: neu co 1 SQL injection, co the co 10 cho khac. Neu co 1 password plain-text, co the co ca API key plain-text.
>
> "3 nguyen tac vang cua security trong Go backend:
> 1. Khong bao gio tin vao input tu client
> 2. Khong bao gio luu password plain-text — dung bcrypt
> 3. Khong bao gio concatenate string vao SQL — dung parameterized query
>
> "Them vao do: CORS khong phai security feature — no chi la browser policy. Dung CORS de bao ve API la hieu lam. Nhung van can cau hinh dung de frontend co the goi duoc."

---

### TODO Comments (Code Skeleton)

```go
package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TODO-[1]: Password hashing voi bcrypt
// SENIOR ASKS: Tai sao bcrypt ma khong phai MD5, SHA-256, hay tu viet?
// HINT: bcrypt duoc thiet ke cho password — no cham (co tinh), co salt tu dong, chong rainbow table

// SENIOR ASKS: bcrypt.DefaultCost = 10. Co nen tang len 12, 14?
// HINT: Cost 10 = ~100ms. Cost 12 = ~400ms. Can bang giua security va UX.

func HashPassword(password string) (string, error) {
	// TODO: Dung bcrypt.GenerateFromPassword voi cost tu config
	// SENIOR ASKS: Tai sao tra ve string thay vi []byte?
	// HINT: De luu vao DB (TEXT/VARCHAR) va so sanh
}

func CheckPassword(password, hash string) bool {
	// TODO: Dung bcrypt.CompareHashAndPassword
	// SENIOR ASKS: Ham nay tra ve bool — tai sao khong tra ve error?
	// HINT: De dung trong handler: if !CheckPassword(...) { return 401 }
	//       Nhung van log error ben trong
}

// TODO-[2]: SQL parameterized queries — KHONG BAO GIO concatenate
// SENIOR ASKS: Cach nao sau day an toan? Tai sao?
//   db.Query("SELECT * FROM users WHERE name = '" + name + "'")
//   db.Query("SELECT * FROM users WHERE name = $1", name)
// HINT: Cai 1 cho phep SQL injection. Cai 2 parameterized — DB xu ly escaping.

// BAD — Dung lam:
// func SearchUsers(db *sql.DB, query string) ([]User, error) {
//     sql := "SELECT * FROM users WHERE name LIKE '%" + query + "%'"
//     return db.Query(sql) // SQL INJECTION!
// }

// GOOD:
// func SearchUsers(db *sql.DB, query string) ([]User, error) {
//     sql := "SELECT * FROM users WHERE name LIKE $1"
//     return db.Query(sql, "%"+query+"%")
// }

// TODO-[3]: Input validation struct
// SENIOR ASKS: Tai sao validation nen o struct level thay vi tung handler?
// HINT: De tai su dung, de test, khong de quen validate o cho nao do

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Name     string `json:"name" validate:"required,max=100"`
	Age      int    `json:"age" validate:"min=13,max=150"`
}

// Validator co the dung go-playground/validator hoac tu viet
// Truoc het, tu viet de hieu logic:

func ValidateCreateUser(req *CreateUserRequest) map[string]string {
	errors := make(map[string]string)

	// TODO: Validate email format (co @, co domain)
	// TODO: Validate password length (8-72 — bcrypt gioi han 72 bytes)
	// TODO: Validate name khong rong, khong qua dai
	// TODO: Validate age hop le

	return errors
}

// TODO-[4]: CORS middleware
// SENIOR ASKS: CORS la gi? No bao ve server khoi attack gi?
// HINT: Cross-Origin Resource Sharing — browser policy, khong phai security feature chong attack
//       No ngan frontend domain A goi API domain B neu khong duoc phep

func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// TODO: Check origin co trong allowedOrigins khong
			// SENIOR ASKS: Neu allowedOrigins = ["*"] — co an toan khong?
			// HINT: Khong! * cho phep moi origin. Production nen liet ke cu the.

			// TODO: Set headers Access-Control-Allow-Origin, Methods, Headers
			// TODO: Handle preflight request (OPTIONS method)
			// TODO: Goi next handler
		})
	}
}

// TODO-[5]: Security headers
// SENIOR ASKS: Nhung headers nao nen them cho moi response?
// HINT: X-Content-Type-Options, X-Frame-Options, Content-Security-Policy, v.v.

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Them cac security headers
		// w.Header().Set("X-Content-Type-Options", "nosniff")
		// w.Header().Set("X-Frame-Options", "DENY")
		// w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// TODO-[6]: Khong expose internal details trong error
// SENIOR ASKS: Neu DB query loi — ban tra ve gi cho client?
// HINT: Log chi tiet (de debug), tra ve generic message (de khong lo thong tin)

func SanitizeError(err error, isDev bool) string {
	if isDev {
		// TODO: Development — co the tra ve chi tiet de debug
		return err.Error()
	}
	// TODO: Production — chi tra ve generic message
	// "Internal server error" — khong de lo SQL query, stack trace, DB host
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. bcrypt co gioi han 72 bytes cho password. Neu user nhap password 100 ky tu, dieu gi xay ra? Ban co nen reject khong?
2. Parameterized query chong SQL injection nhu the nao? Ben duoi DB lam gi voi parameter?
3. CORS chi la browser policy — vay con cach nao khac de goi API tu origin khac ma khong qua CORS?
4. `X-Content-Type-Options: nosniff` chong lai attack gi? Khong co no thi sao?
5. Neu 1 API tra ve error "pq: duplicate key value violates unique constraint 'users_email_key'" — co van de gi? Ban fix nhu the nao?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: HashPassword va CheckPassword hoat dong voi bcrypt
- [ ] TODO-[2] hoan thanh: Hieu va ap dung parameterized queries (khong concatenate)
- [ ] TODO-[3] hoan thanh: ValidateCreateUser kiem tra day du cac truong
- [ ] TODO-[4] hoan thanh: CORS middleware chi cho phep trusted origins
- [ ] TODO-[5] hoan thanh: Security headers duoc them cho moi response
- [ ] TODO-[6] hoan thanh: Error response khong lo internal details trong production

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: HashPassword tra ve hash khac nhau moi lan — vi sao case nay quan trong?
  - **Giai thich:** bcrypt tu dong them salt — hash giong nhau 2 lan = khong co salt = khong an toan
- [ ] Test case: CheckPassword voi password sai -> false — boundary case gi co the fail?
  - **Giai thich:** Timing attack — ham so sanh co tra ve ngay khi khac length khong?
- [ ] Test case: SQL query voi input chua quote — van an toan — vi sao?
  - **Giai thich:** Day la test chong SQL injection — dau ' va " deu phai duoc xu ly an toan

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** bcrypt cham (100ms+). Co nen dung goroutine de hash khong? Risk la gi?
2. **Neu requirement thay doi:** Neu can rate limiting cho login API (chong brute force), ban them no o dau?
3. **Architecture decision:** Tai sao toi khong dung `database/sql` context-aware trong skeleton? Khi nao context quan trong?

---
---

## Mini-Project: Production HTTP Server

### User Story

> Khach hang (Product Owner) noi: "Tat ca topics tren tich hop vao 1 server production-ready"

**Context:** Ban da hoc 6 topics rieng le. Bay gio la luc tich hop tat ca vao 1 server hoan chinh: structured logging voi slog, Prometheus metrics, graceful shutdown, config tu env, security (bcrypt, CORS, input validation). Server phai co health checks, khong expose internal errors, va dung het nhung gi ban da hoc trong Phase 5.

### Acceptance Criteria

- [ ] Server dung `slog` voi JSON handler — log structured, khong fmt.Printf
- [ ] Prometheus metrics tai `/metrics` — request count, duration, active requests
- [ ] Graceful shutdown voi signal handling — doi in-flight requests
- [ ] Config tu env vars — validate khi startup, fail fast
- [ ] Security: bcrypt cho password, parameterized SQL, input validation, CORS
- [ ] Health checks: `/health/live` (liveness) va `/health/ready` (readiness)
- [ ] Error handling: custom error types, wrapping, khong expose internal details
- [ ] Khong secrets trong code — dung .env cho local

---

### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Day la luc ban chung minh ban hieu cach cac thanh phan ket hop voi nhau. Production server khong chi la 'viet code chay duoc' — no la su can bang giua: reliability (graceful shutdown, error handling), observability (logs, metrics, health), security (bcrypt, validation, CORS), va operability (config, deployment).
>
> "Hoi toi review 1 server junior viet — code chay duoc nhung: 1) dung fmt.Printf, 2) khong co graceful shutdown, 3) password luu plain-text, 4) config hardcode, 5) metrics = 0. Production code khong phai ve chay duoc — no la ve chay duoc Lien Tuc, Duoc Bao Ve, va Co The Quan Sat.
>
> "Toi se huong dan ban thiet ke tung thanh phan va cach noi chung lai. Khong viet code hoan chinh — ban phai tu noi cac manh lai."

---

### TODO Comments (Code Skeleton)

```go
// ===== main.go =====
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"
)

// TODO-[1]: Khoi tao tat ca dependencies theo dung thu tu
// SENIOR ASKS: Thu tu khoi tao co quan trong khong? Config truoc hay Logger truoc?
// HINT: Config la goc — moi thu khac phu thuoc config. Config -> Logger -> DB -> Server

func main() {
	// Buoc 1: Load config
	// cfg, err := config.Load()
	// if err != nil { log.Fatalf("failed to load config: %v", err) }
	// if err := cfg.Validate(); err != nil { log.Fatalf("invalid config: %v", err) }

	// Buoc 2: Init logger
	// logger := logger.Init(cfg.Env)
	// slog.SetDefault(logger)

	// Buoc 3: Init DB connection pool
	// db, err := sql.Open("pgx", cfg.DatabaseURL)
	// ... configure pool: SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime

	// Buoc 4: Wire HTTP handlers
	// mux := http.NewServeMux()
	// wireRoutes(mux, db, cfg)

	// Buoc 5: Wrap middleware (theo dung thu tu!)
	// handler := security.SecurityHeadersMiddleware(mux)
	// handler = security.CORSMiddleware(cfg.AllowedHosts)(handler)
	// handler = logger.RequestLogger(logger)(handler)
	// handler = metrics.MetricsMiddleware(handler)
	// handler = recovery.RecoverMiddleware(logger)(handler) // catch panic

	// Buoc 6: Start server voi graceful shutdown
	// server := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	// runServer(server, cfg.ShutdownTimeout)
}

// TODO-[2]: Wire routes
// SENIOR ASKS: Tai sao tach wireRoutes ra ham rieng?
// HINT: De test — ban co the goi wireRoutes trong test ma khong can chay main

func wireRoutes(mux *http.ServeMux, db *sql.DB, cfg *config.Config) {
	// TODO: Register cac handler
	// mux.Handle("/health/live", ...)   // liveness
	// mux.Handle("/health/ready", ...)  // readiness
	// mux.Handle("/metrics", ...)       // Prometheus
	// mux.Handle("/api/users", ...)     // API endpoints
	// SENIOR ASKS: Tai sao /metrics khong qua logging middleware?
	// HINT: No se tao vong lap — metrics middleware log metrics, log middleware log request to metrics
}

// TODO-[3]: Graceful shutdown hoan chinh
// SENIOR ASKS: Tai sao can 2 goroutine (1 chay server, 1 cho signal)?
// HINT: ListenAndServe block — phai chay async. Main thread doi signal de dieu khien.

func runServer(server *http.Server, shutdownTimeout time.Duration, logger *slog.Logger) {
	// TODO: Goroutine chay server
	// TODO: Main thread doi signal
	// TODO: Health check set unhealthy
	// TODO: Grace period
	// TODO: server.Shutdown voi timeout
	// TODO: Close DB
	// TODO: Log "server stopped"
}

// ===== handlers/user.go =====
package handlers

// TODO-[4]: User handler voi day du validation + error handling
// SENIOR ASKS: Handler nen nhe hay nang? Business logic o dau?
// HINT: Handler chi parse input + goi service + tra ve response. Logic o service layer.

type UserHandler struct {
	// TODO: Dependencies: UserService, Logger
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// Parse request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// TODO: Tra ve 400 voi loi parse JSON
		// Khong lo detail cua error parsing (co the chua internal info)
	}

	// Validate input
	if validationErrs := security.ValidateCreateUser(&req); len(validationErrs) > 0 {
		// TODO: Tra ve 400 voi danh sach loi cu the
	}

	// Goi service
	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		// TODO: Dung errors.Is/errors.As de xac dinh loi -> tra ve status code phu hop
		// ErrConflict -> 409, ErrValidation -> 400, ErrNotFound -> 404, default -> 500
		// Log chi tiet (de debug), tra ve generic message (de khong lo info)
		logger.ErrorContext(ctx, "create user failed", "error", err)
		// Tra ve: {"error": "Internal server error"} (production) hoac chi tiet (dev)
	}

	// Tra ve response
	// TODO: Tra ve 201 Created voi user JSON
}

// TODO-[5]: Error response chuan
// SENIOR ASKS: Tai sao can dinh dang error response chuan?
// HINT: Frontend can parse duoc. Ban muon {error, code, details} thay vi string vo dinh dang.

type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, code string, message string, details map[string]string) {
	// TODO: Set Content-Type, status code, json.Encode ErrorResponse
}

// ===== service/user.go =====
package service

// TODO-[6]: Service layer voi business logic
// SENIOR ASKS: Tai sao tach service rieng khoi handler?
// HINT: De test (khong can HTTP), de reuse (co the dung tu CLI khac), de thay doi (thay DB khong anh huong API)

type UserService struct {
	// TODO: Dependencies: UserRepository, Config
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	// TODO: Hash password bang bcrypt
	// TODO: Goi repository de luu user
	// TODO: Wrap loi repository bang domain error (ErrConflict neu duplicate email)
}

// ===== repository/user.go =====
package repository

// TODO-[7]: Repository layer voi SQL parameterized
// SENIOR ASKS: Tai sao tach repository rieng khoi service?
// HINT: De thay doi DB (PostgreSQL -> MySQL), de mock trong test

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	// TODO: Dung parameterized query: INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3)
	// KHONG BAO GIO: "INSERT ... VALUES ('" + u.Email + "', ..."
	// SENIOR ASKS: pq.ErrUniqueViolation — lam sao biet loi nay de wrap thanh ErrConflict?
	// HINT: errors.As de trich xuat *pq.Error, kiem tra Code == "23505"
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. Middleware duoc goi theo thu tu nao? Neu dao thu tu CORS va auth middleware, dieu gi xay ra?
2. Tai sao tach 3 layers: handler, service, repository? 1 layer co duoc khong? Khi nao nen tach, khi nao nen gop?
3. `server.Shutdown` dong tat ca connections — dieu gi xay ra neu 1 request dang giu transaction DB?
4. Request ID duoc tao o dau? Di nhu the nao xuyen suot handler -> service -> repository -> DB query?
5. Khi nao ban nen dung goroutine trong handler? Khi nao khong nen?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Main khoi tao theo dung thu tu: Config -> Logger -> DB -> Server
- [ ] TODO-[2] hoan thanh: Routes duoc wire dung, /metrics khong qua logging
- [ ] TODO-[3] hoan thanh: Graceful shutdown hoan chinh voi health -> grace -> shutdown
- [ ] TODO-[4] hoan thanh: Handler co validation + error handling + khong lo internal details
- [ ] TODO-[5] hoan thanh: Error response chuan co format {error, code, details}
- [ ] TODO-[6] hoan thanh: Service layer co business logic, bcrypt, wrap errors
- [ ] TODO-[7] hoan thanh: Repository dung parameterized queries, khong SQL injection

---

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Graceful shutdown — request dang xu ly van hoan thanh — vi sao case nay quan trong?
  - **Giai thich:** Day la "proof" cua production-ready — khong mat request nao
- [ ] Test case: SQL injection attempt — input `' OR '1'='1` khong anh huong query — boundary case gi?
  - **Giai thich:** Dieu nay chung minh parameterized queries hoat dong — cong ty khong bi hack
- [ ] Test case: Password khong bao gio tra ve trong response — vi sao?
  - **Giai thich:** Security — du la hash cung khong nen lo. Chi luu, khong doc.
- [ ] Test case: Metrics duoc collect dung sau N requests — kiem dem co chinh xac?
  - **Giai thich:** Dieu nay anh huong alerting — neu count sai, ban bi false positive/negative
- [ ] Test case: Config thieu required field -> app khong start — fail fast
  - **Giai thich:** Phat hien loi som hon la chay 1 luc roi crash

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Tach 3 layers (handler/service/repository) nhieu boilerplate. Khi nao 2 layers du? Khi nao can 4?
2. **Neu requirement thay doi:** Neu can ho tro 10,000 concurrent requests, thiet ke can thay doi gi? (Hint: connection pool, rate limiting, caching)
3. **Architecture decision:** Tai sao toi khong dung framework (Gin, Echo) ma dung stdlib? Luc nao nen chuyen sang framework?
4. **Neu deploy:** Dockerfile cua ban se trong nhu the nao? Multi-stage build? Base image nao?
5. **Observability day du:** Log + metrics + health checks — con thieu gi de co "production-ready" that su? (Hint: distributed tracing, alerting rules, runbooks)

---

## Phase 5 Summary

Ban da hoc 6 topics cot loi cua production systems programming:

| Topic | Core Skill | "3AM Moment" |
|-------|-----------|--------------|
| Error Handling | Wrapping + custom types | Tim duoc goc loi trong 5 phut thay vi 4 tieng |
| Structured Logging | slog + JSON + context | Filter log theo request_id, user_id |
| Graceful Shutdown | Signal + Shutdown + drain | Khong mat request nao khi deploy |
| Observability | Prometheus + health checks | Biet bottleneck o dau truoc khi user phan nan |
| Configuration | Env vars + validate + fail fast | Khong mat 2 tieng vi config sai |
| Security | Bcrypt + SQL params + validation | Khong bi hack vi SQL injection |

**Mini-project tich hop tat ca vao 1 server production-ready.** Server nay la nen tang cho moi backend ban viet sau nay.

---

> **Senior's Final Words:**
>
> "Production code khong phai ve code 'chay duoc'. No la ve code 'chay duoc mai mai, co the quan sat, co the bao ve, va khi loi xay ra — ban biet chinh xac loi tu dau trong 5 phut.'
>
> "Hoi toi, junior developer gioi nhat ma toi tung mentor la nguoi luon hoi: 'Neu cai nay loi luc 3AM, toi co tim duoc nguyen nhan khong?' Neu cau tra loi la 'khong' — refactor lai."
>
> "Phase 5 xong roi. Ban da san sang cho Phase 6: Ecosystem & Architecture. Do la luc ban hoc khi nao roi khoi stdlib, khi nao dung framework, va cach thiet ke system thuc su lon."
