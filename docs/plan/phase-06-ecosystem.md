# Phase 6: Ecosystem & Architecture (Tuần 11)

> **Context:** Đây là phase chuyển mình từ "viết Go đúng" sang "viết Go production". Bạn đã học stdlib, concurrency, generics, error handling, graceful shutdown — giờ là lúc tích hợp ecosystem thực tế: router chi mạnh mẽ, PostgreSQL với pgx/sqlx, gRPC cho internal service, JWT auth, Flutter bridge, và chiến lược test hiệu quả.
>
> **Senior's rule of thumb:** "90% của production bug không đến từ việc bạn không biết Go syntax — nó đến từ việc bạn chọn sai tool cho đúng bài toán, hoặc integrate chúng sai cách."

---

## Topic 06.1: Router (chi)

### User Story

> **Khách hàng (Product Owner) nói:** "API có 50 endpoints cần grouping, middleware auth. `net/http` quá verbose — mỗi route phải viết handler registration riêng, middleware chain lộn xộn."
>
> **Context:** Bạn đã viết REST API bằng stdlib `net/http` ở Phase 3. API đó có 4-5 endpoint thì ổn, nhưng khi lên 50 endpoints với nhiều nhóm (public, authenticated, admin) thì code routing trở thành spaghetti.

### Acceptance Criteria

- [ ] Chi router được khởi tạo và mount vào `http.Server`
- [ ] Route groups: `/api/v1/public/*`, `/api/v1/auth/*`, `/api/v1/admin/*` với middleware khác nhau
- [ ] Middleware chain: logging, recover (panic handler), auth check — áp dụng ở group level
- [ ] Sub-router hoặc nested group cho resource CRUD (`/api/v1/tasks/*`)
- [ ] Route parameters: `/api/v1/tasks/{id}` — extract `id` trong handler
- [ ] Chi router vẫn là `http.Handler` — compatible với stdlib middleware bạn đã viết ở Phase 3

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Nếu tôi nhận ticket này, điều đầu tiên tôi nghĩ đến là: 'tôi đã có bao nhiêu middleware
> viết bằng stdlib ở Phase 3? Liệu chúng có tái sử dụng được không?' Chi được chọn vì nó
> LÀ stdlib — nó chỉ cung cấp routing tree và middleware chain trên `net/http`. Không
> magical context như Gin, không custom handler signature. Đây là lý do tôi chọn chi
> cho 90% project thực tế."
>
> "Vấn đề cốt lõi ở đây là: routing phải phản ánh domain structure. Public/auth/admin
> là 3 domain concern khác nhau — code phải thể hiện điều đó qua route groups."
>
> "Tôi sẽ phân rã thành các bước:
>  1. Tạo root router và mount vào server
>  2. Tách 3 route groups với middleware riêng
>  3. CRUD sub-router cho tasks resource
>  4. Đảm bảo middleware từ Phase 3 vẫn hoạt động"
>
> "Hồi tôi ở project fintech, chúng tôi có 200+ endpoints. Chi route groups giúp
> code tổ chức theo domain: `/v1/payments/*`, `/v1/settlements/*`, mỗi group có
> middleware auth riêng. Không group = điên cả đầu vì phải gắn middleware từng route."
```

---

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TODO-[1]: Khởi tạo chi router và cấu hình global middleware
// SENIOR ASKS: Tại sao global middleware phải đăng ký TRƯỚC route definitions?
// HINT: Chi middleware stack là LIFO — thứ tự đăng ký quyết định thứ tự thực thi

func setupRouter() chi.Router {
	r := chi.NewRouter()

	// TODO: đăng ký 3 global middleware: RequestID, RealIP, Logger
	// SENIOR ASKS: Recoverer middleware phải đứng thứ mấy trong stack? Tại sao?
	// HINT: Nếu Recoverer không phải outermost, panic sẽ bypass middleware khác

	// TODO-[2]: Tạo 3 route groups với middleware khác nhau
	// SENIOR ASKS: Chi dùng With() hay Group()? Khác nhau gì trong thực tế?
	// HINT: Group() tạo sub-router isolate; With() wrap inline — chọn đúng mục đích

	// Public routes — không cần auth
	r.Route("/api/v1/public", func(r chi.Router) {
		// TODO: register /health, /register, /login
		// HINT: đây là nơi ngưởi dùng chưa có token, đừng đặt auth middleware
	})

	// Authenticated routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		// TODO: đăng ký middleware authJWT — phải đứng trước route handlers
		// SENIOR ASKS: Nếu authJWT fail, handler có được gọi không? Chi xử lý thế nào?
		// HINT: Middleware có thể write response và KHÔNG gọi next.ServeHTTP()

		// TODO: register /tasks CRUD: GET list, POST create, GET /{id}, PUT /{id}, DELETE /{id}
		// HINT: chi.URLParam(r, "id") để lấy route parameter
	})

	// Admin routes — auth + admin role check
	r.Route("/api/v1/admin", func(r chi.Router) {
		// TODO: authJWT + adminOnly middleware chain
		// SENIOR ASKS: adminOnly lấy user info từ đâu sau khi authJWT pass?
		// HINT: chi không có Gin-style context — dùng request.Context() với value
	})

	return r
}

// TODO-[3]: Auth middleware — extract và verify JWT từ header
// SENIOR ASKS: Tại sao middleware nên return http.Handler thay vì chi-specific type?
// HINT: Interface http.Handler cho phép reuse middleware giữa chi và stdlib

type contextKey string

const contextKeyUserID contextKey = "userID"

func authJWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: extract Bearer token từ Authorization header
			// TODO: verify JWT token — nếu fail: w.WriteHeader(401), return
			// TODO: nếu pass: inject userID vào r.Context(), gọi next.ServeHTTP()
			// HINT: đừng gọi next nếu auth fail — đó là cách middleware "dừng" request
		})
	}
}

// TODO-[4]: Admin role check middleware — chạy SAU authJWT
// SENIOR ASKS: adminOnly có nên gọi authJWT lại không, hay giả định authJWT đã chạy?
// HINT: Coupling middleware = bug tiềm ẩn — nhưng ở đây adminOnly KHÔNG THỂ chạy trước auth

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: lấy userID từ context, query DB/role cache kiểm tra admin role
		// TODO: nếu không phải admin: 403 Forbidden
		// HINT: dùng r.Context().Value(contextKeyUserID) — nhớ type assertion
	})
}

// TODO-[5]: Task handlers — mỗi handler nhận http.ResponseWriter + *http.Request
// SENIOR ASKS: Tại sao chi handler signature giống hệt stdlib http.HandlerFunc?
// HINT: Đây là LÝ DO chọn chi — zero vendor lock-in ở handler level

type TaskHandler struct {
	// TODO: inject repository interface — không dùng global variable
	// HINT: dependency injection qua struct field, không dùng global *sql.DB
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: parse query params (page, limit), gọi repository, trả JSON
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: decode JSON body, validate, gọi repository, trả 201 + created task
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: chi.URLParam(r, "id"), parse UUID/int, query DB, trả 200 hoặc 404
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **Chi vs Gin:** Gin dùng `gin.Context` với method `c.Param()`, `c.JSON()`. Chi dùng
   `http.ResponseWriter + *http.Request` thuần. Tưởng tượng 1 năm sau project của bạn
   muốn migrate từ chi sang `net/http` stdlib — migration cost khác biệt thế nào?

2. **Middleware ordering:** Bạn có 4 middleware: Logger, Recoverer, AuthJWT, RateLimit.
   Thứ tự nào đúng? Tại sao RateLimit phải đứng trước hay sau AuthJWT? Điều gì xảy ra
   nếu Logger đứng SAU Recoverer?

3. **Route parameter validation:** `chi.URLParam(r, "id")` trả về string. Bạn nên validate
   format UUID/nguyên ở middleware layer, handler layer, hay repository layer? Mỗi lựa
   chọn có trade-off gì về separation of concerns?

4. **Context value pattern:** `r.Context().Value("userID")` trả `interface{}` — bạn phải
   type assertion. Điều gì xảy ra nếu một handler quên check `ok` trong type assertion?
   Chiến lược nào để biến lỗi runtime này thành compile-time safety?

5. **Chi không có binding/validation như Gin:** Bạn phải tự decode JSON body và validate
   input. Đây là "nhược điểm" hay "đặc điểm"? Khi nào tự viết validation lại tốt hơn
   dùng framework built-in?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: Chi router khởi tạo với 3+ global middleware (RequestID, Logger, Recoverer)
- [ ] TODO-[2] hoàn thành: 3 route groups `/public`, `/auth`, `/admin` với middleware chain riêng
- [ ] TODO-[3] hoàn thành: `authJWT` middleware extract + verify JWT từ Authorization header
- [ ] TODO-[4] hoàn thành: `adminOnly` middleware kiểm tra role từ context
- [ ] TODO-[5] hoàn thành: TaskHandler với 4 methods (List, Create, Get, Delete) đúng REST convention
- [ ] Code chạy được: `go run .` và `curl http://localhost:8080/api/v1/public/health` trả 200
- [ ] Middleware từ Phase 3 (logging, recover) tái sử dụng được mà không sửa logic

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** Request đến `/api/v1/auth/tasks` KHÔNG có header Authorization → 401
  — vì sao case này quan trọng? Đây là first line of defense, test này phải chạy nhanh
- [ ] **Test case:** Request đến `/api/v1/admin/users` với JWT user thường → 403
  — boundary case: auth pass nhưng authorization fail
- [ ] **Test case:** Request đến `/api/v1/auth/tasks/{id}` với id không phải UUID → 400
  — input validation ở handler level
- [ ] **Test case:** Panic trong handler → Recoverer bắt được, trả 500, server không crash
  — đây là LÝ DO có Recoverer middleware, phải test với `httptest.ResponseRecorder`
- [ ] **Test case:** Middleware thứ tự đúng — Logger ghi log request trước hay sau khi
  response được viết? Dùng `httptest` + custom logger để verify

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **Trade-off chi vs stdlib mux:** Bạn mất bao nhiêu dòng code khi dùng stdlib `ServeMux`
   so với chi cho 50 endpoints? Con số đó có đáng để add dependency không? Tiêu chí của
   bạn là gì — LOC, performance, developer experience, hay maintainability?

2. **Nếu requirement thay đổi:** PO yêu cầu "thêm versioning /v2/* với handlers khác /v1".
   Chi hỗ trợ điều này như thế nào? Bạn sẽ tổ chức code ra sao để v1 và v2 có thể coexist?

3. **Architecture decision — chi không có binding:** Bạn phải tự decode JSON body. Điều này
   dẫn đến repeated code ở mỗi handler. Chiến lược nào để DRY (Don't Repeat Yourself) mà
   không rơi vào "tự viết mini-framework"? Có nên dùng thư viện validate ngoài không?

4. **Middleware context value — an toàn không?** `context.WithValue` được khuyến cáo chỉ
   dùng cho "request-scoped data". UserID có phải "request-scoped" không? Có pattern nào
   thay thế context value để tránh type assertion risk?
```

---
---

## Topic 06.2: Database (pgx + sqlx)

### User Story

> **Khách hàng (Product Owner) nói:** "Chuyển từ SQLite sang PostgreSQL. Cần connection pool
> tốt, scan struct tự động — không muốn viết `row.Scan(&a, &b, &c)` thủ công cho 20 cột."
>
> **Context:** SQLite ổn cho development, nhưng production cần PostgreSQL: concurrent writes
> tốt hơn, JSONB support, full-text search, row-level security. `database/sql` stdlib abstract
> driver nhưng scan thủ công là pain point với struct có nhiều field.

### Acceptance Criteria

- [ ] `pgx` driver kết nối PostgreSQL với connection pool cấu hình (min/max conns, timeout)
- [ ] `sqlx` query + scan tự động vào struct dùng `db.Select()` và `db.Get()`
- [ ] Struct tags ``db:"column_name"`` map database columns → struct fields
- [ ] Repository pattern: interface định nghĩa methods, implementation dùng `*sqlx.DB`
- [ ] Prepared statement cho query lặp lại (trong loop hoặc hot path)
- [ ] Graceful shutdown đóng DB connection pool

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Chuyển DB là thay đổi infrastructure — nó ảnh hưởng mọi layer của application. Điều
> đầu tiên tôi check: repository interfaces từ Phase 3 có abstract SQLite đủ không? Nếu
> interface thiết kế đúng, đổi implementation từ SQLite sang PostgreSQL chỉ là 1 file mới."
>
> "Vấn đề cốt lõi: connection pool. PostgreSQL có giới hạn max connections (mặc định 100).
> Nếu app của tôi scale lên 10 instances, mỗi instance mở 100 connections = 1000 connections
> > database limit. Tôi phải cấu hình pool size đúng và hiểu connection lifecycle."
>
> "Tôi sẽ phân rã:
>  1. pgx connection string + sslmode config
>  2. Connection pool tuning: SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime
>  3. sqlx struct scanning với db tags
>  4. Repository pattern: interface + pgx implementation
>  5. Graceful shutdown: db.Close() trong shutdown sequence"
>
> "Hồi tôi ở project logistics, chúng tôi migrate từ MySQL sang PostgreSQL. Bug kinh điển
> nhất là quên `rows.Close()` — connection pool bị drain, app treo. Tôi luôn dùng
> `defer rows.Close()` NGAY sau khi check error từ Query. Không bao giờ đặt defer ở cuối
> function — nó phải sát với creation."
```

---

#### TODO Comments (Code Skeleton)

```go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// TODO-[1]: Cấu hình connection pool — SỐNG CÒN trong production
// SENIOR ASKS: Nếu SetMaxOpenConns = 100 và bạn chạy 20 replicas, tổng connection
//              đến PostgreSQL là bao nhiêu? Database mặc định chịu được bao nhiêu?
// HINT: max_connections mặc định PostgreSQL = 100; pool size = max_open_conns + 1 (superuser)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string // TODO: "disable" cho dev, "require" cho production — đừng hardcode

	// Pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	// TODO: xây connection string từ cfg
	// HINT: pgx connection string format: "host=... port=... user=... password=... dbname=... sslmode=..."

	// TODO: mở connection dùng sqlx.Connect với pgx driver
	// HINT: sqlx.Connect = sql.Open + db.Ping — khác gì sql.Open thuần?

	db, err := /* ... */
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	// TODO: cấu hình connection pool — SENIOR ASKS: tại sao phải config pool NGAY sau connect?
	// HINT: Default = 0 unlimited open connections, 2 idle — KHÔNG phù hợp production
	// db.SetMaxOpenConns(cfg.MaxOpenConns)
	// db.SetMaxIdleConns(cfg.MaxIdleConns)
	// db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// TODO-[2]: Repository pattern — interface nhỏ, implementation cụ thể
// SENIOR ASKS: Tại sao interface đặt ở consumer (service layer) thay vì repository package?
// HINT: Go idiom: "accept interfaces, return concrete types" — interface thuộc về ngườI dùng

// TaskRepository định nghĩa contract cho task storage.
// Đặt ở service layer hoặc domain package, không phải repository package.
type TaskRepository interface {
	List(ctx context.Context, userID string, limit, offset int) ([]Task, error)
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

// PostgresTaskRepository là implementation dùng sqlx + PostgreSQL.
type PostgresTaskRepository struct {
	db *sqlx.DB
}

func NewPostgresTaskRepository(db *sqlx.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Task model với db tags — SENIOR ASKS: Tại sao CẦN db tags? Không dùng được không?
// HINT: sqlx dùng db tag để map column name → struct field; không tag = dùng field name

type Task struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description *string   `db:"description" json:"description,omitempty"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TODO-[3]: List tasks với sqlx.Select — auto-scan vào slice
// SENIOR ASKS: sqlx.Select vs Queryx + for-loop — khi nào chọn cái nào?
// HINT: Select tiện cho small-medium result; Queryx cho streaming large result

func (r *PostgresTaskRepository) List(ctx context.Context, userID string, limit, offset int) ([]Task, error) {
	var tasks []Task
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	// TODO: dùng sqlx.SelectContext để scan tự động
	// HINT: named parameter ($1, $2) không phải supported natively bởi sqlx — dùng ? hoặc Rebind
	return tasks, nil
}

// TODO-[4]: GetByID với sqlx.Get — scan 1 row vào struct
// SENIOR ASKS: sqlx.Get trả error gì khi row không tồn tại? Phải handle như thế nào?
// HINT: sql.ErrNoRows — đây là sentinel error từ database/sql, KHÔNG phải lỗi pgx riêng

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	var task Task
	query := `SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`
	// TODO: dùng sqlx.GetContext — xử lý sql.ErrNoRows riêng (trả nil, nil hoặc custom NotFound error)
	return &task, nil
}

// TODO-[5]: Prepared statement trong hot path
// SENIOR ASKS: Khi nào prepared statement thực sự có lợi? Khi nào nó gây hại?
// HINT: PostgreSQL lưu prepared statement per-connection; nếu query chỉ chạy 1 lần
//       thì prepare + execute tốn 2 round-trip > 1 ad-hoc query

func (r *PostgresTaskRepository) Create(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO tasks (id, user_id, title, description, status, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :status, :created_at, :updated_at)
	`
	// TODO: dùng sqlx.NamedQuery hoặc prepare + exec với struct
	// HINT: sqlx.Named giúp map struct field → named parameters — rất tiện cho INSERT phức tạp
	return nil
}

// TODO-[6]: Graceful shutdown — đóng DB connection
// SENIOR ASKS: Tại sao KHÔNG NÊN dùng defer db.Close() trong NewPostgresDB?
// HINT: DB lifecycle = application lifecycle; đóng ở main() shutdown sequence, không phải constructor

func (r *PostgresTaskRepository) Close() error {
	// TODO: gọi db.Close() — đợi in-flight queries xong rồi close connections
	return nil
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **Connection pool math:** App của bạn chạy 5 replicas, mỗi replica cấu hình MaxOpenConns=25.
   PostgreSQL max_connections=100. Điều gì xảy ra khi traffic spike và tất cả replicas đồng
   loạt mở max connections? Tính toán: có đủ connections không? Nếu không, pool behavior
   của `database/sql` sẽ là gì — block hay fail?

2. **sqlx.Select vs raw Query:** `sqlx.Select` load toàn bộ result vào memory. Tưởng tượng
   bảng `tasks` có 10 triệu rows, user query không có limit. Điều gì xảy ra? Giải pháp nào
   để vừa tiện lợi của sqlx, vừa tránh OOM?

3. **Prepared statement lifecycle:** `sqlx.DB` tự động prepare statement khi dùng `Get/Select`.
   Statement được lưu per-connection. Khi connection bị đóng (do idle timeout hay pool evict),
   prepared statement bị mất. Điều này có ý nghĩa gì với connection pool tuning?

4. **Nullable columns:** `description` trong Task struct là `*string` (pointer) chứ không phải
   `string`. Tại sao? Điều gì xảy ra nếu dùng `string` và database column có NULL value?

5. **Repository interface scope:** TaskRepository có 5 methods. Một service chỉ cần `GetByID`
   và `Create` — nó nên phụ thuộc vào `TaskRepository` đầy đủ hay một interface nhỏ hơn?
   Principle nào của Go/SOLID áp dụng ở đây?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: `NewPostgresDB` mở connection, cấu hình pool (MaxOpen/MaxIdle/Lifetime)
- [ ] TODO-[2] hoàn thành: `TaskRepository` interface ở domain/service layer, `PostgresTaskRepository` implement
- [ ] TODO-[3] hoàn thành: `List` dùng `sqlx.SelectContext` scan auto vào `[]Task`
- [ ] TODO-[4] hoàn thành: `GetByID` dùng `sqlx.GetContext`, xử lý `sql.ErrNoRows` đúng
- [ ] TODO-[5] hoàn thành: `Create` dùng `sqlx.NamedQuery` hoặc prepared statement
- [ ] TODO-[6] hoàn thành: Graceful shutdown gọi `db.Close()` trong shutdown sequence
- [ ] `docker-compose up postgres` chạy, app kết nối thành công, CRUD hoạt động

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** Connection pool max connections — verify `SetMaxOpenConns` được apply
  — dùng `db.Stats()` để kiểm tra số open connections không vượt limit
- [ ] **Test case:** `GetByID` với ID không tồn tại → `sql.ErrNoRows` được wrap thành domain error
  — boundary case quan trọng: caller phải phân biệt "not found" vs "database error"
- [ ] **Test case:** `List` với empty result → trả `[]Task{}` (empty slice) chứ không phải `nil`
  — nil slice và empty slice khác nhau khi JSON encode: `null` vs `[]`
- [ ] **Test case:** Graceful shutdown — sau `db.Close()`, new query trả lỗi
  — đảm bảo shutdown sequence đúng: stop accepting requests → đợi in-flight → close DB
- [ ] **Test case:** Context cancellation — `List` với `ctx.Done()` đang active → query dừng sớm
  — đây là lý do mọi DB method phải nhận `context.Context`

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **Trade-off sqlx vs raw database/sql:** sqlx thêm dependency để tiết kiệm ~5 dòng
   `row.Scan()` mỗi query. Với 20 queries trong codebase, bạn tiết kiệm ~100 dòng nhưng
   thêm 1 dependency. Bạn có chắc project cần sqlx không, hay raw `database/sql` đủ?
   Threshold nào để quyết định "đáng" add dependency?

2. **Nếu requirement thay đổi:** PO yêu cầu "support cả SQLite và PostgreSQL — user tự
   chọn ở config". Architecture của bạn hỗ trợ điều này không? Cần thay đổi gì ở
   repository layer? Có cần build tag hay compile-time switch không?

3. **Connection pool tuning — "set and forget?"** Bạn đặt MaxOpenConns=25, MaxIdleConns=10.
   Sau 6 tháng, traffic tăng 3x. Ai phát hiện pool đang bottleneck? Bạn thêm observability
   gì để monitor connection pool health? `db.Stats()` có metrics nào hữu ích?

4. **sqlx.NamedQuery — tiện nhưng ẩn complexity:** NamedQuery đằng sau có thể chạy
   regex parse query, construct positional params. Điều này có performance cost không?
   Khi nào bạn nên dùng positional params ($1, $2) thay vì named params?
```

---
---

## Topic 06.3: gRPC & Protocol Buffers

### User Story

> **Khách hàng (Product Owner) nói:** "Microservice internal cần giao tiếp hiệu quả.
> REST JSON quá chậm và nặng — serialize mất 5ms, payload JSON gấp 3 lần binary.
> Cần gRPC cho service-to-service, nhưng web client vẫn cần REST."
>
> **Context:** Hệ thống có 3 services: API Gateway, Task Service, Notification Service.
> Giao tiếp hiện tại qua REST JSON. Khi traffic cao, JSON serialization/deserialization
> trở thành bottleneck. Cần binary protocol cho internal, nhưng browser/client vẫn gọi REST.

### Acceptance Criteria

- [ ] `.proto` file định nghĩa TaskService với 4 RPC: CreateTask, GetTask, ListTasks, UpdateTask
- [ ] gRPC server chạy trên port riêng (50051), implement TaskService interface
- [ ] REST Gateway translate HTTP/JSON → gRPC (dùng grpc-gateway)
- [ ] Client streaming: UploadTaskAttachments — upload nhiều file chunks qua gRPC stream
- [ ] Context propagation: deadline/timeout từ REST gateway → gRPC service
- [ ] TLS/mTLS cho gRPC production (self-signed cert cho dev)

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "gRPC cho internal, REST cho external. Đây là rule tôi luôn tuân theo từ năm 2018.
> Lý do: browser không gọi gRPC native (trừ gRPC-Web, nhưng đó là story khác). External
> clients cần REST/JSON. Internal services giao tiếp binary hiệu quả hơn."
>
> "Vấn đề cốt lõi: không phải 'gRPC nhanh hơn REST' — mà là 'binary serialization +
> HTTP/2 multiplexing + strongly typed contract' = reliability cao hơn cho microservices.
> Protocol Buffers là contract-first: đổi .proto = đổi cả server lẫn client."
>
> "Tôi sẽ phân rã:
>  1. Viết .proto: message types + service definitions
>  2. Generate Go code: protoc-gen-go + protoc-gen-go-grpc
>  3. Implement gRPC server (TaskServiceServer interface)
>  4. Setup REST gateway: HTTP → gRPC translation
>  5. Streaming: client-side upload chunks
>  6. Context propagation: timeout/deadline pass-through"
>
> "Hồi tôi ở project IoT platform, chúng tôi có 40+ microservices. gRPC + protobuf
> giảm cross-service latency từ ~15ms xuống ~3ms. Nhưng bug kinh điển là quên update
> .proto version — client v1.2 gọi server v1.1, field mới bị silently ignore.
> Tôi luôn dùng version trong proto package và validate compatibility."
```

---

#### TODO Comments (Code Skeleton)

```protobuf
// TODO-[1]: Viết .proto file định nghĩa service và message types
// SENIOR ASKS: Tại sao nên viết .proto TRƯỚC khi viết Go code? Đây là workflow gì?
// HINT: Contract-first development — API contract = single source of truth

syntax = "proto3";

package taskmanager.v1;
option go_package = "github.com/youruser/taskmanager/proto/taskmanager/v1;taskv1";

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto"; // cho REST gateway annotations

// Task message — tương ứng với Task struct trong DB
message Task {
  string id = 1;
  string user_id = 2;
  string title = 3;
  // TODO: field number — tại sao KHÔNG ĐƯỢC đổi number sau khi đã ship?
  // HINT: Protobuf binary encoding dùng field number làm key; đổi = break compatibility
  string description = 4;
  TaskStatus status = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
}

enum TaskStatus {
  TASK_STATUS_UNSPECIFIED = 0; // proto3: first enum value MUST be 0
  TASK_STATUS_PENDING = 1;
  TASK_STATUS_IN_PROGRESS = 2;
  TASK_STATUS_DONE = 3;
}

// Request/Response messages — SENIOR ASKS: tại sao nên có message riêng cho mỗi RPC?
// HINT: Request/Response decoupling cho phép thêm field mà không ảnh hưởng RPC khác

message CreateTaskRequest {
  string user_id = 1;
  string title = 2;
  string description = 3;
}

message CreateTaskResponse {
  Task task = 1;
}

message GetTaskRequest {
  string id = 1;
}

message GetTaskResponse {
  Task task = 1;
}

message ListTasksRequest {
  string user_id = 1;
  int32 page_size = 2;
  string page_token = 3; // SENIOR ASKS: tại sao dùng page_token thay vì offset?
  // HINT: page_token cho phía server encode cursor/offset; client không cần biết logic
}

message ListTasksResponse {
  repeated Task tasks = 1;
  string next_page_token = 2;
}

// TODO-[2]: Service definition với REST gateway annotations
// SENIOR ASKS: Annotations trong .proto có ảnh hưởng đến pure gRPC client không?
// HINT: Không — annotations chỉ đọc bởi grpc-gateway plugin, gRPC client bỏ qua

service TaskService {
  rpc CreateTask(CreateTaskRequest) returns (CreateTaskResponse) {
    option (google.api.http) = {
      post: "/v1/tasks"
      body: "*"
    };
  }

  rpc GetTask(GetTaskRequest) returns (GetTaskResponse) {
    option (google.api.http) = {
      get: "/v1/tasks/{id}"
    };
  }

  rpc ListTasks(ListTasksRequest) returns (ListTasksResponse) {
    option (google.api.http) = {
      get: "/v1/tasks"
    };
  }

  // TODO-[3]: Client streaming RPC — upload file chunks
  // SENIOR ASKS: Tại sao streaming phù hợp cho upload file hơn unary RPC?
  // HINT: Unary = load toàn bộ file vào memory; streaming = chunk-by-chunk, O(chunk) memory

  rpc UploadTaskAttachments(stream UploadTaskAttachmentRequest) returns (UploadTaskAttachmentResponse);
}

message UploadTaskAttachmentRequest {
  oneof data {
    string task_id = 1;
    bytes chunk = 2;
  }
}

message UploadTaskAttachmentResponse {
  repeated string attachment_ids = 1;
}
```

```go
package grpcserver

import (
	"context"
	"io"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	taskv1 "github.com/youruser/taskmanager/proto/taskmanager/v1"
)

// TODO-[4]: Implement gRPC server — generated interface
// SENIOR ASKS: protoc generate interface rất lớn. Bạn PHẢI implement tất cả methods?
// HINT: Go compiler bắt buộc implement all interface methods — nhưng có thể embed UnimplementedTaskServiceServer

type TaskServer struct {
	taskv1.UnimplementedTaskServiceServer // embed để forward compatibility
	repo TaskRepository // TODO: inject repository — cùng interface từ Topic 06.2
}

func NewTaskServer(repo TaskRepository) *TaskServer {
	return &TaskServer{repo: repo}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	// TODO: validate input, gọi repo.Create, map domain model → proto message
	// HINT: nên có layer mapper riêng giữa domain model và proto message — không viết inline
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}

func (s *TaskServer) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {
	// TODO: gọi repo.GetByID, xử lý not found → codes.NotFound
	// HINT: map repository errors → gRPC status codes: NotFound, InvalidArgument, Internal
	return nil, status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}

func (s *TaskServer) ListTasks(ctx context.Context, req *taskv1.ListTasksRequest) (*taskv1.ListTasksResponse, error) {
	// TODO: parse pagination, gọi repo.List, map result
	return nil, status.Errorf(codes.Unimplemented, "method ListTasks not implemented")
}

// TODO-[5]: Client streaming — UploadTaskAttachments
// SENIOR ASKS: Server phải handle gì khi client disconnect giữa chừng (network error)?
// HINT: ctx.Done() sẽ signal; server phải cleanup partial upload (temp files, etc.)

func (s *TaskServer) UploadTaskAttachments(stream taskv1.TaskService_UploadTaskAttachmentsServer) error {
	var taskID string
	var chunks []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// TODO: process complete upload — lưu file, trả response
			return stream.SendAndClose(&taskv1.UploadTaskAttachmentResponse{
				// AttachmentIds: ...
			})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}

		switch data := req.Data.(type) {
		case *taskv1.UploadTaskAttachmentRequest_TaskId:
			taskID = data.TaskId
		case *taskv1.UploadTaskAttachmentRequest_Chunk:
			chunks = append(chunks, data.Chunk...)
			// TODO: nếu file lớn, không nên accumulate vào memory — write ra temp file
		}

		// TODO: kiểm tra ctx.Done() — client có thể cancel giữa chừng
		select {
		case <-stream.Context().Done():
			return status.Errorf(codes.Canceled, "client cancelled")
		default:
		}
	}
}

// TODO-[6]: Khởi động gRPC server + REST gateway cùng lúc
// SENIOR ASKS: Hai server chạy trên 2 port khác nhau. Điều gì xảy ra khi graceful shutdown?
// HINT: Shutdown sequence: 1) stop gateway accepting new HTTP requests 2) drain in-flight
//       3) stop gRPC server 4) wait for RPC completion

func RunServer(grpcAddr, gatewayAddr string, repo TaskRepository) error {
	// TODO: tạo gRPC server với TLS interceptor
	grpcServer := grpc.NewServer(/* TODO: interceptor */)
	taskv1.RegisterTaskServiceServer(grpcServer, NewTaskServer(repo))

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	// TODO: chạy gRPC server trong goroutine

	// TODO: tạo REST gateway (grpc-gateway) mux, register handlers, chạy HTTP server

	// TODO: graceful shutdown: nhận signal → shutdown gateway → grpcServer.GracefulStop()

	return nil
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **Protobuf field evolution:** Bạn thêm field `priority` (int32) vào Task message với
   field number 8. Client cũ (không biết field này) gọi server mới. Server cũ phản hồi
   client mới. Điều gì xảy ra trong mỗi case? Protobuf backward/forward compatibility
   hoạt động thế nào? Khi nào compatibility BỊ PHÁ VỠ?

2. **gRPC vs REST — "nhanh hơn" đúng không?** Benchmark thực tế: gRPC (protobuf + HTTP/2)
   vs REST (JSON + HTTP/1.1) trong 3 scenario: (a) small payload < 1KB, (b) large payload
   1MB, (c) high-frequency ping 1000 req/s. Khi nào gRPC thực sự có lợi? Khi nào
   difference không đáng kể?

3. **REST Gateway — double hop penalty:** HTTP request → gateway → gRPC service.
   Mỗi hop có serialization cost. Điều này có nghĩa là REST gateway LUÔN chậm hơn
   native REST handler không? Tính latency penalty của extra serialization hop.
   Khi nào penalty này chấp nhận được?

4. **Streaming error handling:** Client streaming RPC — nếu stream.Recv() trả error
   ở giữa upload (sau khi đã nhận 50% file), server đã lưu partial data ở đâu đó.
   Chiến lược cleanup nào đảm bảo không để rác (orphan temp files) trong storage?

5. **Proto versioning:** Bạn có .proto package `taskmanager.v1`. Sau 1 năm, cần breaking
   change. Bạn chọn `taskmanager.v2` hay vẫn dùng v1? Nếu chọn v2, bao lâu thì deprecated
   v1? Ai quyết định khi nào shutdown v1 — technical team hay business team?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: `.proto` file định nghĩa Task, enum TaskStatus, request/response messages
- [ ] TODO-[2] hoàn thành: TaskService với 4 RPC + REST gateway annotations
- [ ] TODO-[3] hoàn thành: Client streaming UploadTaskAttachments với oneof (metadata + chunk)
- [ ] TODO-[4] hoàn thành: Go gRPC server implement 4 unary RPC methods
- [ ] TODO-[5] hoàn thành: UploadTaskAttachments stream handler xử lý chunks + EOF + cancel
- [ ] TODO-[6] hoàn thành: Dual server — gRPC (50051) + REST gateway (8080) cùng chạy
- [ ] `protoc` generate code thành công, `go build` pass, `grpcurl` test RPC thành công

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** CreateTask với title rỗng → codes.InvalidArgument
  — protobuf validation ở server level, trước khi chạm DB
- [ ] **Test case:** GetTask với ID không tồn tại → codes.NotFound
  — gRPC status code mapping từ repository error
- [ ] **Test case:** Client streaming upload 3 chunks → file complete, đúng content
  — verify: tổng size = sum(chunks), hash match
- [ ] **Test case:** Client cancel stream giữa chừng → server cleanup temp file
  — dùng `context.WithCancel` + `cancel()` giữa stream sends
- [ ] **Test case:** REST gateway call `POST /v1/tasks` (HTTP/JSON) → translated thành
  gRPC CreateTask → trả HTTP 200 + JSON response
  — verify: gateway chuyển đổi đúng HTTP status ↔ gRPC status code

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **gRPC complexity cost:** Bạn thêm protobuf compiler, code generation, 2 server processes,
   gateway layer. Nếu hệ thống chỉ có 2 services giao tiếp, gRPC có "overkill" không?
   Ngưỡng nào (số services, traffic volume, latency requirement) để gRPC "đáng"?

2. **Nếu requirement thay đổi:** PO yêu cầu "real-time notification từ server → client".
   gRPC server streaming (Server-Side Streaming RPC) khác WebSocket/SSE thế nào?
   Khi nào chọn gRPC streaming, khi nào chọn SSE? Browser compatibility ảnh hưởng
   quyết định không?

3. **Schema evolution — "who owns the contract?"** .proto file là API contract giữa
   teams/services. Nếu 3 teams đều sửa .proto cùng lúc, conflict giải quyết thế nào?
   Bạn có proto registry (Buf Schema Registry) hay dùng git merge? Kinh nghiệm thực
   tế: "protobuf conflict" là một trong những merge conflict đau đầu nhất — tại sao?

4. **Context propagation — "deadline across services":** REST gateway nhận request với
   30s timeout. Nó gọi gRPC service A, A gọi gRPC service B. Deadline/timeout nên
   propagate như thế nào qua các hop? Nếu service B chậm 25s, service A còn bao nhiêu
   time? Đây là "distributed timeout" problem — cách giải quyết?
```

---
---

## Topic 06.4: Authentication (JWT)

### User Story

> **Khách hàng (Product Owner) nói:** "Mobile app cần login, token lưu local. Server cần
> verify mà không query DB mỗi request — hệ thống có 10K concurrent users, query DB auth
> mỗi request sẽ kill database."
>
> **Context:** Hệ thống cần 2 flow: (1) User đăng ký/đăng nhập → server trả JWT token,
> (2) Mọi request sau đó gửi JWT trong header → server verify locally (không query DB).
> Password phải được hash trước khi lưu — không bao giờ lưu plaintext.

### Acceptance Criteria

- [ ] JWT generate: tạo access token (15 phút) và refresh token (7 ngày) khi login
- [ ] JWT verify: middleware verify token signature + expiry không cần DB
- [ ] Bcrypt hash: password hash trước khi lưu DB, verify khi login
- [ ] Token refresh flow: refresh token → new access token + new refresh token (rotation)
- [ ] Logout: blacklist refresh token (Redis/cache hoặc DB table)
- [ ] Secure defaults: HS256/RS256, secret rotation strategy

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Stateless auth = JWT = scalability. Nhưng JWT không phải silver bullet — nó có trade-off
> rất lớn: không thể revoke immediately (vì stateless). Cái PO muốn 'verify không query DB'
> chính là stateless benefit, nhưng cũng là limitation."
>
> "Vấn đề cốt lõi: tôi cần cân bằng giữa performance (stateless) và security (revocation).
> Giải pháp thực tế: access token stateless (15 phút) + refresh token stateful (revoke được).
> Nếu user logout, access token vẫn valid đến hết 15 phút — đây là accepted trade-off."
>
> "Tôi sẽ phân rã:
>  1. Bcrypt hash password — đừng bao giờ tự viết hash function
>  2. JWT generate: access (short) + refresh (long) với different secrets
>  3. JWT verify: parse → verify signature → check expiry → extract claims
>  4. Middleware: extract token từ Authorization header
>  5. Refresh flow: verify refresh token → issue new pair → blacklist old refresh
>  6. Logout: add refresh token to blacklist (Redis/DB)"
>
> "Hồi tôi audit 1 project startup, họ lưu JWT secret trong code repo (commit history),
> không có expiry, dùng HS256 với secret là 'mysecret123'. Bị breach trong 1 tuần.
> Tôi luôn: (1) secret từ env var, (2) expiry < 1 giờ cho access token, (3) RS256 cho
> production (private key sign, public key verify — public key có thể distribute an toàn)."
```

---

#### TODO Comments (Code Skeleton)

```go
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO-[1]: Bcrypt password hashing — KHÔNG BAO GIỜ tự viết hash
// SENIOR ASKS: bcrypt.DefaultCost = 10. Nên dùng bao nhiêu cho production? Tại sao không
//              dùng cost cao nhất có thể?
// HINT: Cost tỷ lệ thuận với thờI gian hash; cost quá cao = DoS vector (login chậm =
//       kẻ tấn công dùng CPU của bạn miễn phí). Thường dùng 10-12.

func HashPassword(password string) (string, error) {
	// TODO: bcrypt.GenerateFromPassword với cost phù hợp
	// HINT: output là byte slice chứa salt + hash — lưu trực tiếp vào DB
	return "", nil
}

func CheckPassword(password, hash string) bool {
	// TODO: bcrypt.CompareHashAndPassword — trả error nếu không match
	// HINT: timing attack safe — function này LUÔN mất cùng thờI gian dù match hay không
	return false
}

// TODO-[2]: JWT claims structure — SENIOR ASKS: Tại sao KHÔNG NÊN lưu sensitive data trong JWT?
// HINT: JWT payload base64-encoded, không encrypted — anyone decode được (dù không sửa được)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair chứa cả access và refresh token
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TODO-[3]: JWT generate — access + refresh với thờI hạn khác nhau
// SENIOR ASKS: Tại sao access token ngắn (15 phút) mà refresh token dài (7 ngày)?
// HINT: Access token stateless = không revoke được; nến bị steal, thờI gian exploit window
//       phải ngắn. Refresh token lưu DB nên revoke được — cho phép window dài hơn.

func GenerateTokenPair(userID, email, role string, accessSecret, refreshSecret []byte) (*TokenPair, error) {
	now := time.Now()

	// TODO: tạo access token — expiry 15 phút, sign với accessSecret
	accessClaims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}
	// TODO: jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)

	// TODO: tạo refresh token — expiry 7 ngày, sign với refreshSecret (KHÁC accessSecret!)
	// HINT: nếu attacker lấy được accessSecret, refresh token vẫn an toàn (khác secret)

	return nil, fmt.Errorf("not implemented")
}

// TODO-[4]: JWT verify — middleware dùng trong HTTP handler chain
// SENIOR ASKS: Hàm verify nhận secret làm parameter thay vì global variable — tại sao?
// HINT: Secret rotation — bạn có thể có nhiều secrets valid cùng lúc trong transition period

func VerifyAccessToken(tokenString string, secret []byte) (*Claims, error) {
	// TODO: jwt.Parse với custom Claims, keyfunc trả secret
	// HINT: KeyFunc phải check signing method — tránh "none" algorithm attack

	// TODO: kiểm tra token.Valid, Claims.ExpiresAt
	// TODO: trả *Claims nếu valid, error nếu invalid/expired

	return nil, fmt.Errorf("not implemented")
}

// TODO-[5]: Refresh token flow — verify refresh, issue new pair, blacklist old
// SENIOR ASKS: Tại sao refresh LUÔN trả cả access token MỚI và refresh token MỚI?
// HINT: Token rotation — mỗi refresh = invalidate old refresh token; nếu attacker steal
//       refresh token, legitimate user sẽ detect (refresh token không work) → forced re-login

func RefreshToken(refreshToken string, refreshSecret []byte, tokenStore TokenStore) (*TokenPair, error) {
	// TODO: verify refresh token signature
	// TODO: check refresh token KHÔNG trong blacklist (gọi tokenStore.IsBlacklisted)
	// TODO: nếu valid: generate new pair, blacklist old refresh token
	// TODO: nếu blacklisted: có thể là token reuse attack → invalidate toàn bộ user session
	return nil, fmt.Errorf("not implemented")
}

// TokenStore là interface cho refresh token blacklist — có thể implement bằng Redis hoặc DB
type TokenStore interface {
	Blacklist(ctx context.Context, tokenID string, expiry time.Duration) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

// TODO-[6]: Auth middleware — tích hợp vào chi router
// SENIOR ASKS: Middleware nên extract claims và inject vào context hay return claims trực tiếp?
// HINT: Context injection cho phép downstream handlers access claims; nhưng cần helper function
//       để tránh repeated type assertion

func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: extract Authorization: Bearer <token> header
			// TODO: VerifyAccessToken, nếu fail → 401
			// TODO: inject claims vào r.Context(), gọi next.ServeHTTP
			// HINT: viết helper GetUserIDFromContext(ctx) để downstream handler dùng
		})
	}
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **JWT stateless vs stateful trade-off:** "Verify không query DB" là benefit của stateless.
   Nhưng khi user đổi password, tất cả JWT cũ vẫn valid đến hết expiry. Bạn chấp nhận
   risk này (15 phút window) hay thêm "token version" check (query DB)? Khi nào mỗi
   approach phù hợp?

2. **HS256 vs RS256:** HS256 dùng symmetric secret (cùng key sign + verify). RS256 dùng
   private key sign, public key verify. Với microservices (5 services cần verify token),
   RS256 cho phép distribute public key mà không expose signing capability. Với monolith
   (1 server), HS256 đơn giản hơn. Bạn chọn gì? Tại sao?

3. **Refresh token reuse detection:** User refresh → new pair issued. Attacker dùng old
   refresh token → phát hiện reuse. Nhưng làm sao phân biệt "reuse attack" với "legitimate
   user retry do network timeout"? False positive = lock out legitimate user. Chiến lược
   nào giảm false positive?

4. **Token storage ở client:** Mobile app lưu access token ở đâu? Keychain/Keystore (secure)
   hay SharedPreferences (insecure)? Refresh token lưu khác biệt thế nào? iOS Keychain
   vs Android EncryptedSharedPreferences — trade-off giữa security và convenience?

5. **Secret rotation:** Bạn cần rotate JWT secret (compromise suspicion, hoặc policy 90 ngày).
   Làm sao rotate mà không logout toàn bộ user đang active? Strategy "overlap period"
   hoạt động thế nào?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: `HashPassword` và `CheckPassword` dùng bcrypt với cost hợp lý
- [ ] TODO-[2] hoàn thành: Claims struct với UserID, Email, Role, RegisteredClaims
- [ ] TODO-[3] hoàn thành: `GenerateTokenPair` tạo access (15m) + refresh (7d) với different secrets
- [ ] TODO-[4] hoàn thành: `VerifyAccessToken` parse + verify signature + check expiry
- [ ] TODO-[5] hoàn thành: `RefreshToken` verify refresh + check blacklist + rotation
- [ ] TODO-[6] hoàn thành: Auth middleware extract token → verify → inject claims vào context
- [ ] Login flow: register → hash password → login → trả token pair → access resource bằng token
- [ ] Logout flow: refresh token bị blacklist, không thể dùng refresh nữa

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** `HashPassword` với password "hello" → hash khác nhau mỗi lần gọi
  — bcrypt tự động generate random salt; verify vẫn pass dù hash khác nhau
- [ ] **Test case:** `CheckPassword` với password sai → false (KHÔNG panic, KHÔNG error)
  — timing attack safe: sai password mất cùng thờI gian với đúng password
- [ ] **Test case:** Access token expired → verify trả error "token expired"
  — dùng `jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))` trong test
- [ ] **Test case:** Access token tampered (sửa payload) → verify trả error "invalid signature"
  — sửa 1 ký tự trong token string rồi verify
- [ ] **Test case:** Refresh token reuse → second use bị phát hiện, toàn bộ session bị invalidate
  — test blacklist check và potential cascade logout
- [ ] **Test case:** Auth middleware KHÔNG có Authorization header → 401, handler KHÔNG được gọi
  — verify middleware short-circuit behavior

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **JWT vs Session:** "Stateless = scalable" là argument pro-JWT. Nhưng session với Redis
   cũng scalable (shared session store). So sánh: JWT (stateless) vs Session + Redis
   (stateful) trong 3 khía cạnh: scalability, security (revocation), complexity.
   Khi nào session + Redis tốt hơn JWT? Bạn có chắc JWT là đúng choice không?

2. **Nếu requirement thay đổi:** "Support OAuth2 login (Google, Apple) thêm vào". JWT
   system hiện tại integrate với OAuth2 flow thế nào? Bạn dùng JWT của riêng app hay
   dùng token từ Google? "Sign-in with Apple" trả identity token — cách verify và
   link với local user account?

3. **Bcrypt cost tuning:** Bạn chọn cost=10. Sau 2 năm, CPU nhanh hơn, bcrypt cost=10
   hash nhanh gấp đôi (dễ bị brute force hơn). Bạn tăng cost cho user mới — nhưng
   password hash cũ của user cũ vẫn ở cost=10. Chiến lược "hash upgrade on login"
   hoạt động thế nào? Có cách nào force re-hash không?

4. **Token size:** JWT chứa claims có thể lớn (userID, email, role, permissions...).
   Mỗi request gửi JWT trong Authorization header. Nếu token > 4KB, một số proxy/server
   từ chối header. Bạn có giới hạn claims size không? "Thin JWT" (chỉ userID) + "fat
   session" (full profile trong Redis) có phải better pattern?
```

---
---

## Topic 06.5: Flutter Bridge

### User Story

> **Khách hàng (Product Owner) nói:** "Flutter app cần gọi API Go, real-time notification,
> upload ảnh. API cần contract rõ ràng để mobile team integrate dễ dàng."
>
> **Context:** Mobile team đang viết Flutter app cần giao tiếp với Go backend. Các tính năng:
> (1) Gọi REST API để CRUD tasks, (2) Real-time notification khi task được assign,
> (3) Upload ảnh đính kèm cho task. Mobile team cần biết chính xác request/response format.

### Acceptance Criteria

- [ ] JSON API contract: OpenAPI/Swagger spec hoặc document rõ request/response cho mỗi endpoint
- [ ] REST API endpoints cho Flutter: login, task CRUD, file upload
- [ ] WebSocket server: real-time notification khi task created/updated
- [ ] File upload: multipart form data, lưu file (local/S3), trả URL
- [ ] Error response format chuẩn: `{ "error": "code", "message": "...", "details": {} }`
- [ ] CORS config cho Flutter web (nếu chạy trên browser)

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Mobile + backend = contract-first. Nếu không có contract rõ ràng, 2 team sẽ liên tục
> 'anh ơi API trả gì?' 'em ơi field này tên gì?'. Tôi luôn viết OpenAPI spec hoặc ít
> nhất là document JSON format rõ ràng trước khi implement."
>
> "Vấn đề cốt lõi: real-time trong mobile. HTTP request/response không real-time được.
> WebSocket hoặc SSE là cần thiết. Nhưng mobile có quirk: app background, connection drop,
> battery optimization tắt network. WebSocket cần reconnection logic ở cả client lẫn server."
>
> "Tôi sẽ phân rã:
>  1. API contract document: JSON format cho mỗi endpoint
>  2. REST endpoints: login + task CRUD (reuse từ Topic 06.1, 06.2, 06.4)
>  3. WebSocket: upgrade HTTP connection, broadcast task events
>  4. File upload: multipart parser, temp storage, permanent storage
>  5. Error format chuẩn để Flutter parse dễ dàng
>  6. CORS: config cho Flutter web dev"
>
> "Hồi tôi làm app food delivery, mobile team liên tục complaint vì backend đổi field
> name không báo trước. Tôi bắt đầu dùng contract test: mỗi PR phải pass contract test
> giữa backend response và Dart model. Không pass = không merge. Đây là 'API contract
> testing' và nó cứu chúng tôi hàng tuần."
```

---

#### TODO Comments (Code Skeleton)

```go
package flutterbridge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// =============================================================================
// API CONTRACT (Document này = single source of truth cho mobile team)
// =============================================================================
// POST /api/v1/auth/register    → {email, password} → {user, access_token, refresh_token}
// POST /api/v1/auth/login       → {email, password} → {user, access_token, refresh_token}
// GET  /api/v1/tasks            → ?page=&limit=     → {tasks[], total, next_page}
// POST /api/v1/tasks            → {title, desc}     → {task}
// GET  /api/v1/tasks/{id}       →                   → {task}
// PUT  /api/v1/tasks/{id}       → {title, desc, status} → {task}
// POST /api/v1/tasks/{id}/attach→ multipart form    → {attachment_url}
// WS   /api/v1/ws               →                   → real-time events

// TODO-[1]: Error response format chuẩn — mobile team parse dễ dàng
// SENIOR ASKS: Tại sao error response PHẢI có structure chuẩn, không phải plain text?
// HINT: Flutter JSON decoder cần predictable format; plain text = parse error ở client

type APIError struct {
	Code    string                 `json:"code"`              // "INVALID_INPUT", "NOT_FOUND", "UNAUTHORIZED"
	Message string                 `json:"message"`           // human-readable
	Details map[string]interface{} `json:"details,omitempty"` // extra context cho debug
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	// TODO: set Content-Type: application/json, write status code, encode APIError
}

// TODO-[2]: REST API endpoints cho Flutter
// SENIOR ASKS: Tại sao KHÔNG nên để Flutter trực tiếp gọi gRPC? gRPC có vấn đề gì ở mobile?
// HINT: gRPC-Web support chưa native trong Flutter (cần plugin), binary proto khó debug,
//       proxy/VPN corporate có thể block HTTP/2 traffic

// TaskHandler đã có từ Topic 06.1 — TODO: thêm method UploadAttachment

type TaskHandler struct {
	repo            TaskRepository
	authSecret      []byte
	wsHub           *WebSocketHub
	uploadMaxSize   int64
	uploadStorePath string
}

// ListTasks — GET /api/v1/tasks
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: extract userID từ context (auth middleware đã inject)
	// TODO: parse query params: page, limit (với default values)
	// TODO: gọi repo.List, trả JSON với pagination info
}

// UploadAttachment — POST /api/v1/tasks/{id}/attach
func (h *TaskHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	// TODO: parse multipart form, enforce max file size
	// SENIOR ASKS: Tại sao KHÔNG NÊN đọc toàn bộ file vào memory?
	// HINT: r.ParseMultipartForm(maxMemory) — file > maxMemory được write ra temp disk file

	// TODO: validate file type (whitelist: jpg, png, pdf)
	// TODO: generate unique filename, lưu vào uploadStorePath (hoặc S3)
	// TODO: trả JSON với attachment URL
}

// TODO-[3]: WebSocket server — real-time notification
// SENIOR ASKS: WebSocket hoạt động thế nào ở tầng protocol? Khác HTTP request/response ở điểm gì?
// HINT: WebSocket bắt đầu bằng HTTP handshake (Upgrade header), sau đó là persistent
//       bidirectional connection — không còn request/response cycle

// WebSocketHub quản lý tất cả active connections và broadcast messages
type WebSocketHub struct {
	// TODO: clients map[*websocket.Conn]bool — thread-safe access cần mutex
	// TODO: broadcast channel — goroutine nhận message và gửi đến all clients
	// TODO: register/unregister channels — goroutine quản lý add/remove clients
}

func NewWebSocketHub() *WebSocketHub {
	// TODO: khởi tạo hub với channels và goroutine loop
	return nil
}

// Message struct gửi qua WebSocket — SENIOR ASKS: Tại sao nên có "type" field?
// HINT: Flutter client cần dispatch message theo type: "task.created", "task.updated",
//       "notification.new" — giống event-driven architecture

type WSMessage struct {
	Type      string          `json:"type"`       // "task.created", "task.assigned", "ping"
	Payload   json.RawMessage `json:"payload"`    // flexible payload theo type
	Timestamp time.Time       `json:"timestamp"`  // client kiểm tra stale message
}

func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// TODO: upgrade HTTP → WebSocket dùng websocket.Upgrader
	// TODO: check origin (CORS cho WebSocket) — SENIOR ASKS: CheckOrigin return true an toàn không?
	// HINT: CheckOrigin nil = allow all origins = NGUY HIỂM production; chỉ allow known origins

	// TODO: register client vào hub
	// TODO: goroutine đọc message từ client (heartbeat/ping-pong handling)
	// TODO: khi client disconnect → unregister, close connection
}

// Broadcast gửi message đến tất cả connected clients
func (h *WebSocketHub) Broadcast(msg WSMessage) {
	// TODO: gửi vào broadcast channel — goroutine loop sẽ fan-out đến clients
}

// TODO-[4]: CORS config cho Flutter web dev
// SENIOR ASKS: CORS chỉ cần cho web (browser) hay cả mobile app? Tại sao?
// HINT: Mobile app (Flutter iOS/Android) không bị CORS restriction — CORS là browser security
//       feature. Nhưng Flutter web và development cần CORS.

func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: check Origin header có trong allowedOrigins không
			// TODO: set Access-Control-Allow-Origin, Allow-Methods, Allow-Headers
			// TODO: handle preflight OPTIONS request
		})
	}
}

// TODO-[5]: WebSocket heartbeat / ping-pong
// SENIOR ASKS: Tại sao WebSocket cần ping-pong? HTTP không cần mà?
// HINT: WebSocket = persistent connection; proxy/firewall có thể tự đóng idle connection.
//       Ping-pong giữ connection alive và phát hiện zombie connection.

func (h *WebSocketHub) run() {
	// TODO: goroutine vĩnh viễn loop:
	//   - select case register: add client
	//   - select case unregister: remove client, close conn
	//   - select case broadcast: iterate clients, write message
	//   - ticker: gửi ping đến idle clients, nếu không pong thì unregister
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **API versioning:** Backend phát triển nhanh, API thay đổi. Mobile app update chậm
   (user không update app ngay). Nếu backend đổi response format, app cũ crash.
   Versioning strategy nào: URL version (/v1/, /v2/), header version (Accept: application/vnd.v2+json),
   hay parameter version (?api-version=2)? Trade-off mỗi approach trong context Flutter
   (có thể hardcode version trong app)?

2. **WebSocket vs SSE:** Server-Sent Events (SSE) là HTTP-based, chỉ server → client
   (unidirectional). WebSocket bidirectional. Flutter task notification chỉ cần server
   push (không cần client → server qua WS). Tại sao KHÔNG dùng SSE? Khi nào SSE đơn
   giản hơn WebSocket?

3. **File upload strategy:** Upload ảnh 5MB qua multipart form đơn giản. Nhưng nếu upload
   video 500MB trên mobile network (3G/4G yếu), multipart form có vấn đề gì? Resumable
   upload (chunked upload) giải quyết gì? Bạn có cần design cho resumable upload ngay
   từ đầu hay YAGNI?

4. **WebSocket authentication:** HTTP request có Authorization header. WebSocket upgrade
   request cũng có header. Nhưng sau khi kết nối established, mọi WS message không có
   HTTP header. Làm sao associate WS connection với authenticated user? Query parameter
   token (ws://.../?token=xxx)? Message-based auth? Trade-off?

5. **Offline-first mobile:** Mobile app cần work offline (tạo task khi không có mạng,
   sync khi có mạng). Backend API design hỗ trợ offline-first thế nào? Conflict resolution
   khi user tạo task offline, sync lên, nhưng server đã có task mới hơn từ device khác?
   Pattern "last-write-wins" vs "CRDT" — bạn chọn gì cho use case task manager?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: `APIError` struct + `WriteError` helper, dùng ở tất cả handlers
- [ ] TODO-[2] hoàn thành: Task CRUD REST endpoints + file upload endpoint
- [ ] TODO-[3] hoàn thành: WebSocket server với `WebSocketHub`, register/unregister, broadcast
- [ ] TODO-[4] hoàn thành: CORS middleware cho Flutter web development
- [ ] TODO-[5] hoàn thành: WebSocket ping-pong heartbeat, zombie connection detection
- [ ] API contract document (Markdown/OpenAPI) hoàn chỉnh cho mobile team
- [ ] Flutter app (hoặc Postman) gọi API thành công: login → create task → upload file → WS notification

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** Upload file vượt max size → 413 Payload Too Large
  — boundary case: file đúng bằng max size vs max size + 1 byte
- [ ] **Test case:** WebSocket connection với invalid token → connection bị từ chối ngay
  — auth ở WebSocket upgrade phase
- [ ] **Test case:** Broadcast message khi 0 client connected → không panic, không block
  — edge case: hub chạy nhưng chưa có client nào
- [ ] **Test case:** WebSocket client disconnect đột ngột (không gửi close frame) → hub
  detect timeout qua ping-pong, unregister đúng
  — zombie connection detection; dùng `websocket.WriteControl` với deadline
- [ ] **Test case:** CORS preflight request → 204 No Content với đúng headers
  — `OPTIONS` request từ browser trước actual request

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **WebSocket scalability:** 1 server instance giữ WebSocket connections trong memory
   (map[*conn]bool). Khi scale lên 3 server instances (load balancer), client A ở server 1
   cần nhận notification về task của client B ở server 2. Broadcast cross-server giải
   quyết thế nào? Redis Pub/Sub? NATS? Shared memory không work vì processes tách biệt.

2. **Nếu requirement thay đổi:** "Support real-time collaborative editing — nhiều user
   edit cùng task đồng thờI". WebSocket + broadcast đơn giản không đủ — cần operational
   transform (OT) hay CRDT. Bạn có cần thay đổi API contract? WebSocket payload format
   thay đổi thế nào? Có thư viện nào hỗ trợ không?

3. **Mobile battery optimization:** WebSocket persistent connection tiêu hao pin.
   Mobile OS (iOS/Android) tắt background network sau vài phút. Chiến lược nào để
   notification vẫn đến khi app background? Firebase Cloud Messaging (FCM) + silent
   push notification → wake app → re-establish WebSocket? Hay đơn giản dùng FCM thay
   cho WebSocket cho push notification?

4. **File upload — local vs cloud:** Bạn lưu file ở local disk (`uploadStorePath`).
   Khi chạy trên cloud (Kubernetes, multiple pods), local disk của pod 1 không thấy
   file từ pod 2. Object storage (S3, MinIO) giải quyết vấn đề này. Presigned URL
   pattern: server generate S3 presigned URL → Flutter upload trực tiếp lên S3 (bypass
   server) có lợi gì về scalability? Trade-off?
```

---
---

## Topic 06.6: Testing Strategy

### User Story

> **Khách hàng (Product Owner) nói:** "Test coverage thấp, integration test chạy 10 phút.
> Cần chiến lược test hiệu quả — unit test nhanh, integration test đúng chỗ, không test
> thừa."
>
> **Context:** Codebase hiện có vài unit test cơ bản nhưng thiếu integration test.
   Mỗi lần chạy integration test cần khởi động PostgreSQL thật, chạy migration, seed
   data — mất 10 phút. Developer không chạy integration test locally nữa, bug production tăng.

### Acceptance Criteria

- [ ] Unit test: repository (mock DB), handlers (httptest), auth (table-driven)
- [ ] Integration test: real PostgreSQL dùng testcontainers-go, chạy < 2 phút
- [ ] Contract test: verify API response format match Flutter Dart model
- [ ] Test pyramid: 70% unit, 20% integration, 10% e2e — đúng tỷ lệ
- [ ] Parallel test: `t.Parallel()` cho unit test độc lập, giảm tổng thờI gian
- [ ] Mock strategy: interface-based mock, không mock external library

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "'Test coverage thấp' không phải root cause — nó là symptom. Root cause thường là:
> (1) test chậm nên dev không chạy, (2) test flaky nên dev không tin, (3) test khó viết
> nên dev không viết. Tôi phải giải quyết cả 3."
>
> "Vấn đề cốt lõi: integration test chạy 10 phút vì khởi động PostgreSQL thật.
> Testcontainers giải quyết điều này — nó chạy PostgreSQL trong Docker container,
> tự cleanup sau test. Nhưng testcontainers vẫn chậm hơn in-memory test.
> Giải pháp: phân biệt rõ unit (fast, in-memory) vs integration (slow, real DB)."
>
> "Tôi sẽ phân rã:
>  1. Unit test: mock repository interfaces, test business logic
>  2. Handler test: httptest.ResponseRecorder, không cần real server
>  3. Integration test: testcontainers-go + PostgreSQL, test real queries
>  4. Contract test: verify JSON schema match Flutter model
>  5. Parallel test: t.Parallel() cho speed
>  6. Makefile: `make test-unit`, `make test-integration`, `make test-all`"
>
> "Hồi tôi ở project e-commerce, integration test chạy 45 phút. Chúng tôi refactor
> thành 3 tiers: unit test chạy < 30 giây (dev chạy liên tục), integration test chạy
> < 3 phút (chạy pre-commit), e2e test chạy 15 phút (chạy CI pipeline). Developer
> happiness tăng vọt, bug production giảm 60%."
```

---

#### TODO Comments (Code Skeleton)

```go
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO-[1]: Mock repository — interface-based, không dùng external mock library phức tạp
// SENIOR ASKS: Tại sao mock ở repository interface thay vì mock sqlx.DB trực tiếp?
// HINT: Mock repository = mock business dependency; mock sqlx = mock implementation detail.
//       Test handler không cần biết handler dùng PostgreSQL hay SQLite.

// MockTaskRepository implement TaskRepository interface cho test
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) List(ctx context.Context, userID string, limit, offset int) ([]Task, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]Task), args.Error(1)
}

func (m *MockTaskRepository) Create(ctx context.Context, task *Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Task), args.Error(1)
}

// TODO-[2]: Unit test handler với httptest — không cần real server hay real DB
// SENIOR ASKS: httptest.ResponseRecorder vs khởi động server trên port thật — khác biệt gì?
// HINT: httptest bypass network stack, không cần port, chạy song song an toàn, nhanh hơn 100x

func TestTaskHandler_List(t *testing.T) {
	t.Parallel() // TODO: đánh dấu parallel cho unit test độc lập

	tests := []struct {
		name       string
		setupMock  func(*MockTaskRepository)
		wantStatus int
		wantTasks  int
		wantErr    string
	}{
		{
			name: "success - returns tasks",
			setupMock: func(m *MockTaskRepository) {
				m.On("List", mock.Anything, "user-1", 20, 0).Return([]Task{
					{ID: "1", Title: "Task 1", UserID: "user-1"},
					{ID: "2", Title: "Task 2", UserID: "user-1"},
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantTasks:  2,
		},
		{
			name: "empty list",
			setupMock: func(m *MockTaskRepository) {
				m.On("List", mock.Anything, "user-1", 20, 0).Return([]Task{}, nil)
			},
			wantStatus: http.StatusOK,
			wantTasks:  0,
		},
		{
			name: "database error",
			setupMock: func(m *MockTaskRepository) {
				m.On("List", mock.Anything, "user-1", 20, 0).Return([]Task{}, errors.New("db down"))
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // TODO: subtest cũng parallel — SENIOR ASKS: Khi nào KHÔNG NÊN parallel?
			// HINT: Không parallel khi test share mutable state (global var, database)

			repo := new(MockTaskRepository)
			tt.setupMock(repo)

			handler := &TaskHandler{repo: repo}

			// TODO: tạo request với userID trong context (simulated auth)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, "user-1"))

			rec := httptest.NewRecorder()
			handler.List(rec, req)

			// TODO: assert status code, response body
			assert.Equal(t, tt.wantStatus, rec.Code)
			// TODO: nếu success, decode JSON và assert task count
			// TODO: verify mock expectations: repo.AssertExpectations(t)
		})
	}
}

// TODO-[3]: Integration test với testcontainers-go
// SENIOR ASKS: Tại sao integration test cần real DB thay vì in-memory SQLite?
// HINT: PostgreSQL behavior khác SQLite (JSONB, concurrency, NULL handling, type system).
//       Test trên SQLite không bắt được bug PostgreSQL-specific.

func TestTaskRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// TODO: khởi động PostgreSQL container dùng testcontainers-go
	// HINT: tc-postgres module trong testcontainers-go — container tự cleanup sau test

	// TODO: chạy migration (schema create)
	// TODO: seed test data
	// TODO: chạy CRUD operations, assert kết quả trên real DB
	// TODO: cleanup: testcontainers tự remove container
}

// TODO-[4]: Contract test — verify API response JSON schema
// SENIOR ASKS: Contract test khác integration test ở điểm gì? Test cái gì?
// HINT: Contract test không test business logic; nó test "response format match expected
//       schema" — đảm bảo Flutter app không bị breaking change khi backend update

func TestAPIContract_TaskResponse(t *testing.T) {
	// TODO: tạo task, gọi API, decode response
	// TODO: verify JSON có tất cả required fields: id, user_id, title, status, created_at
	// TODO: verify field types: id là string, created_at là RFC3339 timestamp
	// TODO: verify KHÔNG CÓ unexpected fields (backward compatible check)
}

// TODO-[5]: Benchmark test — handler performance
// SENIOR ASKS: Benchmark test khi nào cần viết? Mục đích là gì?
// HINT: Viết benchmark khi có performance requirement rõ ràng, hoặc khi optimize
//       suspected bottleneck. Dùng để verify optimization có hiệu quả (before vs after).

func BenchmarkTaskHandler_List(b *testing.B) {
	repo := new(MockTaskRepository)
	repo.On("List", mock.Anything, "user-1", 20, 0).Return([]Task{
		{ID: "1", Title: "Task 1"},
	}, nil)

	handler := &TaskHandler{repo: repo}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, "user-1"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		handler.List(rec, req)
	}
}

// TODO-[6]: Test utilities — helpers giảm boilerplate
// SENIOR ASKS: DRY trong test code có quan trọng không? Hay "test code nên explicit và repeat"?
// HINT: "A little duplication is better than wrong abstraction" — nhưng helpers như
//       setupMockRepo, assertJSONResponse, newAuthenticatedRequest giảm noise đáng kể

func newAuthenticatedRequest(t *testing.T, method, url string, body io.Reader, userID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, url, body)
	req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, userID))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, status int, v interface{}) {
	t.Helper()
	require.Equal(t, status, rec.Code)
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), v))
}
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **Test pyramid vs ice cream cone:** Test pyramid nói 70% unit, 20% integration, 10% e2e.
   "Ice cream cone anti-pattern" là unit ít, integration/e2e nhiều. Codebase của bạn
   hiện giống hình dạng nào? Làm sao chuyển từ ice cream cone về pyramid khi legacy
   code không dễ unit test?

2. **Mocking philosophy:** "Mockist" (mock mọi dependency) vs "Classicist" (chỉ mock
   external systems: DB, API, file system). Bạn thuộc trường phái nào? Trong Go,
   interface-based design hỗ trợ cả hai — nhưng quá nhiều mock = test brittleness
   (thay đổi implementation = test fail dù behavior không đổi). Cân bằng ở đâu?

3. **Flaky test:** Integration test đôi khi fail ngẫu nhiên (container chưa ready,
   race condition). Developer bắt đầu "rerun until pass" — đây là dấu hiệu gì?
   Chiến lược nào để loại bỏ flakiness trong testcontainers-based tests? Retry logic?
   Health check? Wait-for-it pattern?

4. **Test naming convention:** Bạn dùng `TestFunctionName` hay `TestScenario_Description`?
   Ví dụ: `TestCreateTask` vs `TestCreateTask_Success` vs `TestCreateTask_ValidInput`.
   Table-driven test với `tt.name` — convention nào giúp test failure message dễ đọc
   nhất khi chạy `go test -v`?

5. **Coverage target:** PO yêu cầu "80% coverage". Điều này có ý nghĩa không? Bạn có
   thể đạt 80% coverage mà không test edge cases không? Chiến lược "coverage threshold"
   trong CI có giúp quality không, hay chỉ khuyến khích "test đủ để pass threshold"?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: `MockTaskRepository` implement `TaskRepository` interface
- [ ] TODO-[2] hoàn thành: Handler unit test table-driven với `httptest`, parallel execution
- [ ] TODO-[3] hoàn thành: Integration test với `testcontainers-go`, PostgreSQL real
- [ ] TODO-[4] hoàn thành: Contract test verify API JSON response schema
- [ ] TODO-[5] hoàn thành: Benchmark test cho handler hot path
- [ ] TODO-[6] hoàn thành: Test helpers giảm boilerplate (`newAuthenticatedRequest`, `assertJSONResponse`)
- [ ] `make test-unit` chạy < 30 giây, `make test-integration` chạy < 2 phút
- [ ] CI pipeline chạy unit + integration test, fail nếu coverage < 60%

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** Mock không được gọi (handler return sớm do validation) → verify
  `repo.AssertNotCalled` — đảm bảo không query DB cho invalid input
- [ ] **Test case:** WebSocket hub broadcast với concurrent register/unregister → race-free
  — dùng `go test -race` để detect data race trong hub goroutines
- [ ] **Test case:** Integration test context cancellation → DB query dừng, connection
  trả về pool — verify không leak connection khi client cancel request
- [ ] **Test case:** Contract test: response JSON KHÔNG có field `password_hash` — security
  leak test; verify sensitive fields không bao giờ xuất hiện trong API response
- [ ] **Test case:** Parallel subtest isolation — 2 subtest chạy parallel, mỗi test có
  mock repo riêng, không interfere — verify `t.Parallel()` hoạt động đúng

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **Test cost-benefit:** Viết unit test cho getter/setter (boilerplate code) có giá trị
   không? "Test everything" vs "Test behavior có risk". Bạn phân biệt "code đơn giản"
   (không cần test) và "code phức tạp" (cần test) bằng tiêu chí nào? Cyclomatic
   complexity threshold? Business criticality?

2. **Nếu requirement thay đổi:** "Thêm feature flag system — feature bật/tắt runtime".
   Test strategy thay đổi thế nào? Bạn cần test cả 2 path (feature on/off)? Số test
   cases tăng gấp đôi? Có pattern nào giảm test multiplication không?

3. **Testcontainers production parity:** testcontainers chạy PostgreSQL 15, nhưng
   production chạy PostgreSQL 14.5 (RDS). Minor version difference có thể gây bug
   không? "Test trên dev khác production" là eternal problem — chiến lược nào giảm
   thiểu risk? IaC (Terraform) để ensure infrastructure consistency?

4. **Mock maintenance cost:** Mỗi lần thêm method vào TaskRepository interface, tất cả
   mock implementations phải update. Trong codebase lớn (10+ mock), đây là "interface
   change ripple effect". Có tool nào auto-generate mock từ interface (mockery, mockgen)?
   Trade-off của auto-generated mock vs hand-written mock?
```

---
---

## Mini-Project: Task Manager API

### User Story

> **Khách hàng (Product Owner) nói:** "Cần một Task Manager API hoàn chỉnh: chi router,
> PostgreSQL với sqlx, JWT authentication, gRPC service, WebSocket real-time, đầy đủ
> test. Đây là capstone integration của toàn bộ Phase 6."
>
> **Context:** Bạn đã học 6 topics riêng lẻ. Mini-project này yêu cầu integrate tất cả
> thành 1 hệ thống hoàn chỉnh chạy được. Đây là "vertical slice" — từ HTTP request
> đến database, qua auth, WebSocket, gRPC gateway.

### Acceptance Criteria

- [ ] Chi router với route groups (public, auth, admin)
- [ ] PostgreSQL repository với sqlx struct scanning
- [ ] JWT auth: register, login, refresh, logout
- [ ] gRPC service + REST gateway chạy song song
- [ ] WebSocket real-time: task created/updated events
- [ ] File upload: multipart, save to disk, return URL
- [ ] Test: unit (mock), integration (testcontainers), contract
- [ ] Docker Compose: PostgreSQL + app chạy trong 1 command
- [ ] API documentation (Markdown hoặc OpenAPI) cho Flutter team

---

### Senior Thought-Process

```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Integration project = bài test tư duy systems design. Không phải 'viết code nhiều',
> mà là 'kết nối đúng các thành phần'. Lỗi phổ biến nhất ở integration: lifecycle
> management — ai khởi tạo trước, ai đóng sau."
>
> "Vấn đề cốt lõi: dependency startup order. PostgreSQL phải ready trước app.
> gRPC server và REST gateway phải start cùng lúc. WebSocket hub chạy background.
> Graceful shutdown phải reverse order: stop accepting → đợi in-flight → close WS →
> stop gateway → stop gRPC → close DB connections."
>
> "Tôi sẽ phân rã:
>  1. Docker Compose: PostgreSQL container
>  2. App bootstrap: config → DB → repository → handlers → router → servers (HTTP + gRPC)
>  3. Wire chi router với auth middleware + route groups
>  4. gRPC server + REST gateway dual server
>  5. WebSocket hub goroutine
>  6. Graceful shutdown với signal handling
>  7. Test: unit + integration + contract"
>
> "Hồi tôi làm project production đầu tiên, tôi mất 2 ngày debug vì startup order sai.
> App khởi động trước PostgreSQL (do Docker Compose không đợi healthy), migration fail,
> app crash, restart loop. Tôi học được: dùng health check wait script hoặc retry logic
> trong app. Đừng assume dependency đã ready khi app start."
```

---

#### TODO Comments (Code Skeleton)

```go
package main

// TODO-[1]: Application config — tập trung, validate ở startup
// SENIOR ASKS: Tại sao config nên là struct read-only sau khi khởi tạo?
// HINT: Config mutation sau startup = behavior không predictable; fail-fast tại startup

type Config struct {
	HTTPAddr        string        `env:"HTTP_ADDR" envDefault:":8080"`
	GRPCAddr        string        `env:"GRPC_ADDR" envDefault:":50051"`
	DatabaseURL     string        `env:"DATABASE_URL"`
	JWTAccessSecret string        `env:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret string       `env:"JWT_REFRESH_SECRET"`
	UploadMaxSize   int64         `env:"UPLOAD_MAX_SIZE" envDefault:"10485760"` // 10MB
	UploadPath      string        `env:"UPLOAD_PATH" envDefault:"./uploads"`
}

// TODO-[2]: Application struct — dependency container
// SENIOR ASKS: Tại sao dùng struct chứa dependencies thay vì global variables?
// HINT: Dependency injection qua struct = testable (swap implementation), no global state,
//       explicit dependency graph

type Application struct {
	config   Config
	db       *sqlx.DB
	repo     TaskRepository
	taskHandler   *TaskHandler
	authHandler   *AuthHandler
	grpcServer    *grpc.Server
	wsHub         *WebSocketHub
	httpServer    *http.Server
}

func NewApplication(cfg Config) (*Application, error) {
	// TODO: connect DB, run migrations
	// TODO: create repository
	// TODO: create handlers
	// TODO: create WebSocket hub
	// TODO: setup chi router, gRPC server, REST gateway
	return nil, nil
}

// TODO-[3]: Bootstrap — khởi động tất cả components
// SENIOR ASKS: Thứ tự khởi động nào đúng? DB → handlers → servers? Hay servers → handlers?
// HINT: Dependency order: config → infrastructure (DB) → domain (repository) →
//       application (handlers) → transport (HTTP/gRPC/WS servers)

func (app *Application) Run() error {
	// TODO: start gRPC server trong goroutine
	// TODO: start HTTP server (chi router + REST gateway) trong goroutine
	// TODO: start WebSocket hub goroutine (background broadcast loop)

	// TODO: wait for OS signal (SIGTERM, SIGINT)
	// TODO: graceful shutdown với timeout context

	return nil
}

// TODO-[4]: Graceful shutdown — REVERSE startup order
// SENIOR ASKS: Tại sao shutdown phải ngược startup order?
// HINT: Nếu đóng DB trước HTTP server, in-flight requests sẽ fail với DB error.
//       Đúng: stop accepting → drain in-flight → đóng resources

func (app *Application) Shutdown(ctx context.Context) error {
	// TODO: 1) shutdown HTTP server (stop accepting new connections)
	// TODO: 2) graceful stop gRPC server
	// TODO: 3) close WebSocket hub (close all client connections)
	// TODO: 4) close DB connection pool
	// TODO: 5) wait for goroutines to finish với context timeout
	return nil
}

// TODO-[5]: Router setup — tích hợp tất cả từ các topics
// SENIOR ASKS: Nếu bạn có 20 handler methods, setupRouter() sẽ rất dài. Cách tổ chức?
// HINT: Tách thành setupPublicRoutes, setupAuthRoutes, setupAdminRoutes — mỗi function
//       register routes cho 1 group

func (app *Application) setupRouter() chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORSMiddleware([]string{"http://localhost:*"})) // dev only

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", app.authHandler.Register)
		r.Post("/auth/login", app.authHandler.Login)
		r.Post("/auth/refresh", app.authHandler.Refresh)
	})

	// Authenticated routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(AuthMiddleware([]byte(app.config.JWTAccessSecret)))

		r.Get("/tasks", app.taskHandler.List)
		r.Post("/tasks", app.taskHandler.Create)
		r.Get("/tasks/{id}", app.taskHandler.Get)
		r.Put("/tasks/{id}", app.taskHandler.Update)
		r.Delete("/tasks/{id}", app.taskHandler.Delete)
		r.Post("/tasks/{id}/attach", app.taskHandler.UploadAttachment)

		r.Get("/ws", app.wsHub.HandleWebSocket) // WebSocket upgrade
	})

	return r
}

// TODO-[6]: Docker Compose setup cho development
// docker-compose.yml:
// version: "3.8"
// services:
//   postgres:
//     image: postgres:15-alpine
//     environment:
//       POSTGRES_USER: taskmanager
//       POSTGRES_PASSWORD: taskmanager
//       POSTGRES_DB: taskmanager
//     ports:
//       - "5432:5432"
//     volumes:
//       - postgres_data:/var/lib/postgresql/data
//     healthcheck:
//       test: ["CMD-SHELL", "pg_isready -U taskmanager"]
//       interval: 5s
//       timeout: 5s
//       retries: 5
//   app:
//     build: .
//     ports:
//       - "8080:8080"
//       - "50051:50051"
//     environment:
//       DATABASE_URL: postgres://taskmanager:taskmanager@postgres:5432/taskmanager?sslmode=disable
//     depends_on:
//       postgres:
//         condition: service_healthy
// volumes:
//   postgres_data:

// TODO-[7]: Test strategy cho mini-project
// SENIOR ASKS: Mini-project có nhiều components — bạn test integration như thế nào?
// HINT: 3-tier testing:
//   1. Unit: handler + mock repo (chạy < 30s)
//   2. Integration: real DB với testcontainers, test repository + API end-to-end (chạy < 2m)
//   3. Contract: verify JSON response schema không đổi (chạy < 10s)

// Makefile targets:
// test-unit: go test -short ./...
// test-integration: go test -run Integration ./...
// test-contract: go test -run Contract ./...
// test-all: test-unit test-integration test-contract
```

---

#### Socratic Questions

```markdown
**Câu hỏi để bạn tự suy nghĩ:**

1. **Monolith vs Microservices:** Mini-project của bạn là monolith (1 codebase, 1 deploy).
   Nhưng bạn đã thêm gRPC — "microservice communication protocol". Có phải gRPC chỉ dùng
   cho microservices? Monolith dùng gRPC có "over-engineering" không? Khi nào monolith
   nên có internal gRPC interface (ngay cả khi chưa tách service)?

2. **Server startup failure:** App khởi động, kết nối PostgreSQL fail. Bạn chọn:
   (a) crash immediately (fail-fast), (b) retry với backoff, (c) chạy nhưng mark unhealthy.
   Docker Compose/Kubernetes sẽ handle mỗi lựa chọn thế nào? Kinh nghiệm production:
   fail-fast với container orchestration = auto-restart healthy hơn "limp mode".

3. **WebSocket trong cùng HTTP server:** `/api/v1/ws` dùng chi router, upgrade trong
   handler. Nhưng REST gateway cũng chạy trên cùng port 8080. Có xung đột không?
   grpc-gateway có handle WebSocket không? Nếu không, routing như thế nào để WS request
   đi đúng handler mà không qua gateway?

4. **Testing a system vs testing components:** Bạn có unit test cho handler, integration
   test cho repository, contract test cho API. Nhưng ai test "tất cả cùng lúc"?
   End-to-end test (Selenium/Postman collection) test toàn bộ flow. Nhược điểm của e2e
   test là gì? "Ice cream cone" anti-pattern — bạn đang đi về hướng đó không?

5. **Config management:** Config struct đọc từ env var. Trong development, bạn phải set
   5-6 env var mỗi lần chạy. `.env` file giải quyết vấn đề này. Nhưng `.env` file có
   vấn đề gì? (Không phải production-grade, secret trong file, không version-controlled).
   Chiến lược nào: `.env` cho dev, env var cho staging/production, secret manager cho
   sensitive config? Kubernetes Secret/ConfigMap?
```

---

### Output Checklist: Làm sao biết mình xong?

- [ ] TODO-[1] hoàn thành: Config struct với validation, đọc từ env var
- [ ] TODO-[2] hoàn thành: Application struct là dependency container — không global variables
- [ ] TODO-[3] hoàn thành: Bootstrap khởi động DB → handlers → servers (HTTP + gRPC + WS)
- [ ] TODO-[4] hoàn thành: Graceful shutdown reverse order với timeout context
- [ ] TODO-[5] hoàn thành: Chi router với route groups, auth middleware, all CRUD endpoints
- [ ] TODO-[6] hoàn thành: Docker Compose: `docker compose up` chạy PostgreSQL + app
- [ ] TODO-[7] hoàn thành: 3-tier test (unit/integration/contract) đều pass
- [ ] `curl http://localhost:8080/health` trả 200
- [ ] Register → Login → Create Task → Get Task → Upload File → WebSocket notification
- [ ] gRPC `grpcurl` CreateTask → GetTask thành công
- [ ] `make test-all` pass trong < 3 phút

---

### Test Checklist: Những gì bạn nên tự viết test

- [ ] **Test case:** End-to-end: register → login → create task → list tasks → logout
  — verify toàn bộ flow, dùng real HTTP request (httptest server)
- [ ] **Test case:** Graceful shutdown: gửi SIGTERM sau 5 request concurrent → tất cả
  request hoàn thành, server đóng đúng sequence
- [ ] **Test case:** gRPC REST gateway: HTTP request đến gateway → forward đến gRPC
  → trả HTTP response đúng — verify dual-stack hoạt động
- [ ] **Test case:** WebSocket broadcast: 3 client connect, 1 task created → cả 3 client
  nhận WS message — verify fan-out hoạt động
- [ ] **Test case:** File upload → file tồn tại trên disk → download URL trả đúng content
  — verify upload end-to-end, không chỉ "upload success"

---

### Retrospective: Sau khi xong, hãy tự hỏi

```markdown
1. **Architecture bloat:** Bạn thêm chi, sqlx, pgx, gRPC, protobuf, gorilla/websocket,
   jwt library, testify. 7+ external dependencies. Mỗi dependency = risk (security CVE,
   maintenance abandonment, breaking changes). Bạn đánh giá dependency risk thế nào
   trước khi add? Tiêu chí: GitHub stars, release frequency, issue response time,
   Go module graph size?

2. **"I built it, but do I understand it?"** Mini-project chạy được không đồng nghĩa
   hiểu từng component. Bạn có thể giải thích từng dòng trong chi middleware chain
   không? Khi sqlx.Select scan struct, điều gì xảy ra reflection-level? gRPC server
   goroutine model: mỗi RPC = 1 goroutine? WebSocket read/write goroutine separation?
   Hiểu internals = debug được production issues.

3. **Production readiness checklist:** Code chạy local chưa đủ. Production cần:
   structured logging (slog), metrics (Prometheus), distributed tracing (OpenTelemetry),
   health checks (/health, /ready), rate limiting, request timeout, CORS config đúng,
   TLS certificate, secret management, database migration strategy, backup/restore plan.
   Mini-project của bạn còn thiếu gì? Phase nào trong roadmap cover những cái còn thiếu?

4. **"What would break first?"** Nếu hệ thống này đi production với 1000 concurrent users,
   component nào fail đầu tiên? PostgreSQL connection pool? WebSocket hub memory?
   File upload disk I/O? gRPC message size? Đây là "bottleneck analysis" — kỹ năng
   cần cho system design interview và production debugging.

5. **From here to capstone:** Mini-project này là warm-up cho Week 12 capstone
   (Go REST API có DB, auth, Docker, deploy, Flutter gọi được). Bạn sẽ tái sử dụng
   bao nhiêu % code từ mini-project này? Architecture pattern nào transferable?
   Topic nào cần học thêm trước khi đến Week 12?
```

---

## Appendix: Phase 6 Decision Heuristics

```markdown
| Tình huống | Quyết định | Lý do |
|---|---|---|
| Cần routing + middleware | chi hoặc httprouter | Compatible stdlib, không lock-in |
| Struct scan từ SQL | sqlx | Biết SQL trước khi dùng ORM |
| Internal service communication | gRPC + protobuf | Binary, typed contract, streaming |
| Public API cho browser/mobile | REST + JSON | CORS, debug-friendly, universal |
| Lưu password | bcrypt (golang.org/x/crypto) | Never roll your own crypto |
| Token auth | JWT access (stateless) + refresh (stateful) | Balance performance vs revocation |
| Real-time từ server | WebSocket hoặc SSE | WebSocket = bidirectional, SSE = simpler |
| Test database | testcontainers-go | Production parity, auto cleanup |
| Mock dependency | Interface-based mock | Go idiomatic, không cần magic |

## Key Takeaways (Week 11)

1. **Chi = stdlib++:** Chi không thay thế `net/http`, nó bổ sung routing tree và
   middleware chain. Zero vendor lock-in là lý do chọn chi.

2. **pgx + sqlx = raw SQL power:** Biết SQL trước khi dùng ORM. sqlx struct scanning
   tiết kiệm boilerplate nhưng bạn vẫn viết SQL query — điều này tốt cho performance
   và debugging.

3. **gRPC cho internal, REST cho external:** Không dùng gRPC cho public browser API.
   REST gateway cho phép maintain 1 service definition (.proto) phục vụ cả 2 use case.

4. **JWT trade-off:** Stateless auth = scalable nhưng không immediately revocable.
   Access token ngắn + refresh token có blacklist = practical balance.

5. **Contract-first cho mobile:** API contract (OpenAPI hoặc document) là communication
   bridge giữa backend và mobile team. Không có contract = integration hell.

6. **Test pyramid:** 70% unit (fast), 20% integration (real DB), 10% e2e. testcontainers
   cho integration test production-parity mà không cần shared dev database.
```

---

*"Học ecosystem không phải để biết nhiều thư viện — mà để biết KHI NÀO dùng thư viện nào,
KẾT NỐI chúng thế nào, và TRADE-OFF của mỗi lựa chọn."* — Senior's closing note.
