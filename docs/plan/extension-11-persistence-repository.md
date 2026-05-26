# Extension 11: Persistence, Repository Layer & Multi-Repo Storage

> **Meta:** Backend hardening. Phase này biến API từ demo một repo thành service quản lý nhiều repo có metadata, storage policy và test integration.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.

---

## Backend Extension Sprint: Storage + Repository Pattern

> **Mục tiêu extension:** Thiết kế persistence cho nhiều repository: metadata nằm ở database, object/index/refs vẫn nằm trong filesystem per repo để giữ đúng tinh thần Git internals.
>
> **Nguồn:** Sau `extension-10-rest-api.md`
>
> **Mini-git surface:** repository metadata, data root, per-repo workspace, migrations, repository pattern, integration tests

---

## Extension Overview

### Mission
- Extension 11 - Persistence và multi-repo backend storage

### Flutter / Dart Bridge
> Trong Flutter, repository pattern che nguồn dữ liệu: local cache, REST, SQLite. Trong backend Go, repository layer che database metadata. Nhưng đừng nhầm: Git object storage của project vẫn là filesystem để học Git internal; database chỉ quản lý metadata như repo id, owner, name, path, created_at.

### Go Skills Required For This Extension
> `pgx/v5`: `pgxpool.New(ctx, connStr)`, `pgxpool.Config{MaxConns: N}`, `pool.QueryRow(ctx, sql, args...)`, `row.Scan(&field1, &field2)`, `pool.Exec(ctx, sql, args...)`, `pool.Begin(ctx)` — transaction, `tx.Commit(ctx)`, `defer tx.Rollback(ctx)` (safe after Commit). Error handling: `errors.Is(err, pgx.ErrNoRows)`. Migration: raw SQL file + run at startup hoặc migration tool. Interface: define `RepoStore` interface tại consumer side, implementation nhận `*pgxpool.Pool`.

---

## Mission: Multi-Repo Persistence

### User Story
> Nhà tuyển dụng hỏi: *"API này chạy được cho nhiều user/repo không, hay chỉ là wrapper quanh một folder local?"*
>
> Câu trả lời đúng: Backend có metadata store, data root tách biệt, repository layer test được, và mỗi repo có workspace riêng.

### Main Task
Thêm persistence layer cho repo metadata và storage policy cho nhiều repo.

### Acceptance Criteria
- [ ] Có model metadata cho repo: id, name, owner/user, path, created_at, updated_at
- [ ] Có repository interface nhỏ cho metadata CRUD
- [ ] Có implementation database bằng `pgx` hoặc `database/sql`
- [ ] Có migration SQL
- [ ] Object database vẫn nằm trong filesystem per repo dưới data root
- [ ] API không nhận raw filesystem path từ client
- [ ] Repo id map tới workspace path qua storage policy server-side
- [ ] Có lock hoặc policy rõ cho operation ghi cùng repo
- [ ] Có integration test cho metadata repository
- [ ] Có cleanup strategy cho test data
- [ ] `go fmt ./...`, `go vet ./...`, `go test ./...` chạy xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Đừng nhét blob/tree/commit vào Postgres chỉ để khoe database. Mục tiêu học Git internal là object database trên filesystem. DB chỉ nên quản lý metadata backend."
>
> "Tôi cần data ownership rõ: database giữ metadata, filesystem giữ Git-like state. API dùng repo id, không expose path."
>
> "Điểm học chính: repository pattern trong Go là interface nhỏ, test được, không phải folder `repository` chứa mọi thứ."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/store/repository.go
package store

type RepoID string
type RepoMetadata struct{}
type CreateRepoParams struct{}

type RepoStore interface {
	Create(params CreateRepoParams) (RepoMetadata, error)
	FindByID(id RepoID) (RepoMetadata, error)
	ListByOwner(owner string) ([]RepoMetadata, error)
	Delete(id RepoID) error
}

// File: internal/storage/workspace.go
package storage

type WorkspaceResolver interface {
	Resolve(repoID string) (string, error)
}

// TODO-EXT11-A: DB lưu metadata, filesystem lưu object database.
// SENIOR ASKS: Vì sao đưa blob content vào Postgres ở phase này làm mất mục tiêu học Git internal?

// TODO-EXT11-B: Client không được gửi raw server path.
// SENIOR ASKS: Nếu API nhận `C:\Users\...\repo`, security boundary vỡ thế nào?

// TODO-EXT11-C: Operation ghi phải nghĩ tới concurrency.
// SENIOR ASKS: Hai request commit cùng repo cùng lúc có thể làm hỏng ref/index ra sao?
```

#### Theory Notes
- [ ] Database giữ metadata (id, name, path, owner, timestamps) — filesystem giữ object database: đây là data ownership rõ ràng, không phải "DB cho sang"
- [ ] Connection pool: tạo **một** `*pgxpool.Pool` khi khởi động service, dùng chung cho mọi request — không tạo pool per request
- [ ] Row scan pattern:
  ```go
  row := pool.QueryRow(ctx, "SELECT id, name, path FROM repos WHERE id=$1", id)
  var r RepoMetadata
  if err := row.Scan(&r.ID, &r.Name, &r.Path); err != nil {
      if errors.Is(err, pgx.ErrNoRows) { return nil, ErrNotFound }
      return nil, fmt.Errorf("findByID: %w", err)
  }
  ```
- [ ] Transaction pattern:
  ```go
  tx, err := pool.Begin(ctx)
  if err != nil { return err }
  defer tx.Rollback(ctx) // an toàn dùng defer dù sau commit, vì Rollback sau Commit là no-op
  if _, err = tx.Exec(ctx, "INSERT INTO repos ...", args...); err != nil {
      return fmt.Errorf("insert repo: %w", err)
  }
  return tx.Commit(ctx)
  ```
- [ ] Migration: raw SQL file chạy khi start (`CREATE TABLE IF NOT EXISTS`) — hoặc dùng `golang-migrate`
- [ ] Repository pattern trong Go: interface nhỏ defined tại consumer side, implementation nhận pool qua constructor — không phải abstract class

#### Socratic Questions
1. Vì sao repo id an toàn hơn raw path trong API?
2. Metadata transaction có bảo vệ được object file write không?
3. Operation nào cần lock theo repo?
4. Integration test database khác unit test ở điểm nào?

### Output Checklist: Làm sao biết mình xong?
- [ ] API tạo được nhiều repo
- [ ] Mỗi repo có workspace riêng
- [ ] Metadata repo lưu ở database
- [ ] Object/index/refs lưu ở filesystem per repo
- [ ] Không endpoint nào expose raw server path
- [ ] Có migration và hướng dẫn chạy local database

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Unit test storage path resolver
- [ ] Integration test RepoStore create/find/list
- [ ] Test repo id không tồn tại
- [ ] Test duplicate repo name theo owner nếu có constraint
- [ ] Test concurrent write policy ở mức tối thiểu
- [ ] `go test -race ./...`

### Learning Notes / Docs
- [ ] Viết `docs/storage-design.md`
- [ ] Ghi rõ DB schema và filesystem layout
- [ ] Ghi rõ trade-off: không lưu Git objects trong DB

### Retrospective: Sau khi xong, hãy tự hỏi
1. Interface nào quá to?
2. Query nào đang leak domain rule vào handler?
3. Storage path policy có test đủ case Windows/Linux chưa?

---

## Extension Checkpoints (BẮT BUỘC)

### CP-EXT11-A: Persistence Gate
- [ ] Migration chạy được.
- [ ] Repo metadata CRUD chạy được.
- [ ] API dùng repo id để resolve workspace.

### CP-EXT11-B: Storage Gate
- [ ] Object database vẫn là filesystem per repo.
- [ ] Không expose raw path.
- [ ] Có strategy cho cleanup test data.

### CP-EXT11-C: Oral Defense
- [ ] Giải thích được vì sao DB không thay thế object database ở project này.
- [ ] Giải thích được connection pool, migration, repository pattern bằng lời của mình.

## Failure Modes (PHẢI BIẾT)
- Dùng Postgres để lưu mọi blob chỉ vì muốn "có DB".
- API nhận path tùy ý từ client.
- Không có migration, schema tạo thủ công.
- Repository interface quá rộng, khó fake/test.

## Progression Rules

### Rule 1: Database phục vụ backend metadata, không phá Git learning goal.
Object model vẫn phải nhìn thấy được trong `.mgit`.

### Rule 2: Client không điều khiển filesystem path.
Server quyết định workspace từ repo id.

### Rule 3: Integration test phải có setup/teardown rõ.
Test database bẩn là nguồn flaky rất nhanh.

## Tổng Kết

### Deliverables
- [ ] Repo metadata store.
- [ ] Migration SQL.
- [ ] Multi-repo workspace resolver.
- [ ] Integration tests.
- [ ] Storage design doc.

### First-Principles Question
Trong backend Mini Git, dữ liệu nào là metadata nên nằm trong database, dữ liệu nào là Git-like state nên nằm trong filesystem, và vì sao?
