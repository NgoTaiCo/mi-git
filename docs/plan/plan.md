# Plan: Mini Git bằng Go - 14 Days / 8 Milestones

## Paradigm Shift
Từ: Checklist 14 ngày chỉ tick task cho xong
Sang: Mô phỏng môi trường làm việc với senior mentor, mỗi milestone đều có mission, contract, test gate và oral defense.

## Source of Truth
- Roadmap CSV: `docs/14 ngày viết Mini Git bằng Go.csv`
- CLI name: `mgit`
- Repo metadata folder: `.mgit`
- Object storage: SHA-1 + zlib
- Index: JSON để dễ học, không giả vờ là binary index của Git thật

## Product Storyline
`mgit` là một Git-internals learning project viết bằng Go. Mục tiêu public repo không phải "clone Git đầy đủ", mà là chứng minh bạn hiểu cách Git lưu snapshot, object, refs, index và merge ở mức core.

Luồng sản phẩm đi từ thấp lên cao:

1. Repository shell: `.mgit`, HEAD, refs folder.
2. Object database: blob/tree/commit object, SHA-1, zlib.
3. History: commit DAG, refs, log.
4. Staging: index, add, status.
5. Navigation: branch, switch, checkout, detached HEAD.
6. Reconciliation: diff, merge base, three-way merge, conflict.
7. Release: docs, tests, demo, retrospective.

## Architecture Evolution
| Milestone | New package surface | State owned | Commands unlocked |
|---|---|---|---|
| 01 | `internal/repo`, CLI dispatcher | `.mgit/HEAD`, folders | `init` |
| 02 | `internal/object` | `.mgit/objects` | `hash-object`, `cat-file` |
| 03 | `internal/tree` | tree objects | `write-tree` |
| 04 | `internal/commit`, `internal/refs`, `internal/history` | commit objects, branch refs | `commit-tree`, `commit`, `log` |
| 05 | `internal/index`, `internal/status` | `.mgit/index` | `add`, `status` |
| 06 | `internal/worktree` + refs hardening | HEAD symbolic/detached, working tree restore | `branch`, `switch`, `checkout` |
| 07 | `internal/diff`, `internal/merge` | merge result, conflict files | `diff`, `merge` |
| 08 | release/demo/docs scripts | public artifact quality | demo flow |

## Backend Extension Evolution
Core Track 01-08 chứng minh bạn hiểu Git internals và Go package design. Extension Track 09-12 chứng minh core đó đủ sạch để expose thành backend service.

| Extension | New package surface | State owned | Backend capability unlocked |
|---|---|---|---|
| C | `internal/fsck`, `internal/object/cache`, `internal/core` | goroutine-safe cache, context-aware service | Goroutine, channel, mutex, context.Context — nền tảng bắt buộc trước HTTP server |
| 09 | `internal/core` service/usecase boundary | core input/output contract, domain errors | CLI và API dùng chung logic |
| 10 | `cmd/mgit-api`, `internal/api` | HTTP routes, DTO, middleware basics | REST API + `httptest` |
| 11 | `internal/store`, `internal/storage` | repo metadata DB, per-repo workspace | multi-repo backend storage |
| 12 | `internal/config`, `internal/auth`, deployment docs | env config, auth boundary, logs, Docker | production-ish public demo |

## Workflow Mỗi Mission
```
a) User Story: Khách hàng đưa yêu cầu trong context Mini Git
b) Acceptance Criteria: Điều kiện chấp nhận lấy từ CSV
c) Senior Guide:
   - Senior Thought-Process: cách bóc requirement
   - TODO Skeleton: chỉ signature/interface/struct rỗng, không có ruột hàm
   - Socratic Questions: hỏi để chống học vẹt
d) Output Checklist: xác nhận command và state .mgit
e) Test Checklist: edge cases bạn phải tự viết test
f) Retrospective: trade-off, lỗi, phần chưa refactor
g) Phase Checkpoints: CP-XX-A/B/C cho manual flow, test gate, oral defense
```

## Format Code Skeleton
```go
package example

type DomainObject struct{}

type Store interface {
	Load(id string) (DomainObject, error)
}

func DoWork(input string) (DomainObject, error)

// TODO: Bạn tự implement ruột hàm trong source code thật.
// SENIOR ASKS: Function này thuộc package nào? Vì sao không đặt trong main?
```

## Output: 14 Ngày Gom Thành 8 Milestone
| Phase | File | Dates | CSV Scope |
|---|---|---|---|
| 01 | [Repo Foundations & CLI Bootstrap](phase-01-foundations.md) | 2026-07-27 | Day 1: CLI skeleton + mgit init |
| 02 | [Object Database & Blob Plumbing](phase-02-object-database.md) | 2026-07-28 | Day 2: Object database |
| 03 | [Tree Snapshots](phase-03-tree-snapshots.md) | 2026-07-29 | Day 3: Tree object |
| 04 | [Commit Objects, Refs & History](phase-04-commit-refs-history.md) | 2026-07-30 -> 2026-07-31 | Day 4: Commit object; Day 5: Refs, HEAD, commit, log |
| 05 | [Index, Add & Status](phase-05-index-status.md) | 2026-08-01 -> 2026-08-02 | Day 6: Index và add; Day 7: Status |
| 06 | [Branch, Switch & Detached HEAD](phase-06-branch-switch-checkout.md) | 2026-08-03 -> 2026-08-05 | Day 8: Branch; Day 9: Switch branch; Day 10: Checkout và detached HEAD |
| 07 | [Diff, Merge Base & Three-Way Merge](phase-07-diff-merge.md) | 2026-08-06 -> 2026-08-08 | Day 11: Diff; Day 12: Merge base; Day 13: Three-way merge |
| 08 | [Polish, Test, Docs & Demo](phase-08-polish-test-docs-demo.md) | 2026-08-09 | Day 14: Polish, test, docs, demo |

## After Core: Backend Extension Track
Extension Track không thuộc 14 ngày gốc. Nó là phần nâng cấp CV sau khi Mini Git core đã pass demo. Không được dùng extension để né việc hiểu object database, index, refs, checkout và merge.

| Extension | File | Scope | Why it matters for CV |
|---|---|---|---|
| C | [Concurrency Foundation](concurrency-module.md) | Goroutine, channel, mutex, context.Context | **BẮT BUỘC** trước Extension 09: Go interview hỏi goroutine/channel, HTTP handler là goroutine per request |
| 09 | [Core Boundary & API-Ready Package Design](extension-09-core-boundary.md) | Tách CLI khỏi core service | Chứng minh biết package design, testability, adapter boundary |
| 10 | [REST API Adapter](extension-10-rest-api.md) | Expose core qua HTTP | Chứng minh `net/http`/middleware/JSON/`httptest` |
| 11 | [Persistence, Repository Layer & Multi-Repo Storage](extension-11-persistence-repository.md) | Metadata DB + per-repo workspace | Chứng minh PostgreSQL/repository pattern/storage design |
| 12 | [Auth, Observability, Docker & Deploy Demo](extension-12-production-demo.md) | Auth, config, logs, Docker, deploy docs | Chứng minh production-ish backend fundamentals |

## Learning Contract
- Không copy code có ruột hàm từ plan. Plan chỉ là khung suy nghĩ.
- Mỗi milestone phải chạy được command thật và test thật.
- Mỗi mission phải có note tối thiểu 5 dòng bằng lời của bạn.
- Không panic trong normal flow. Input sai phải trả error rõ.
- Pass gate mặc định: `go fmt ./...`, `go vet ./...`, `go test ./...` xanh. Phase nào chạm restore/switch/merge phải kiểm tra kỹ filesystem edge cases.

## Public Repo Quality Bar
Đây là tiêu chuẩn để repo đủ đẹp CV, không phải chỉ "code chạy được":

- `README.md` có project goal, scope, install/run, command examples, demo transcript, limitation.
- `docs/what-i-learned-about-git.md` giải thích bằng lời của bạn: object database, tree, commit, refs, index, checkout, merge.
- `docs/retrospective.md` ghi trade-off, bug đã gặp, phần cố tình chưa làm.
- `CHANGELOG.md` hoặc release notes cho từng milestone.
- Tests có table-driven cases cho package domain quan trọng.
- CLI normal errors không panic; message nói rõ command/path/hash/ref nào lỗi.
- Demo chạy trong temp directory sạch, không phụ thuộc state cũ.
- Không claim tương thích Git thật 100%; ghi rõ không support packfile, remote, protocol, rebase, stash, tag, hooks, binary index thật.
- Nếu làm Backend Extension: README phải nói rõ đây là REST API wrapper quanh Mini Git core, không phải GitHub clone hay Git remote server.
- Nếu làm Backend Extension: có `docs/api.md`, `docs/storage-design.md`, `docs/deployment.md` và `docs/runbook.md`.
- Nếu làm Backend Extension: API dùng repo id/server-side workspace, không nhận raw filesystem path từ client.
- `docs/architecture.md` mô tả package structure, data flow, CLI adapter pattern và tại sao core không phụ thuộc CLI.
- `go test -race ./...` xanh, đặc biệt với Phase 6+ (worktree restore, index update, merge commit).

## Per-Milestone Artifact Matrix
| Milestone | Must ship | Test focus | Docs you write |
|---|---|---|---|
| 01 | `mgit init` | repo init idempotent, root discovery | note repo = working dir + metadata |
| 02 | blob object store | hash vector, zlib, corrupt object | content-addressable storage |
| 03 | tree object | deterministic sort, `.mgit` ignored | blob vs tree, snapshot |
| 04 | commit/ref/log | parent chain, HEAD/ref update order | commit DAG, branch pointer |
| 05 | index/add/status | staged vs unstaged vs untracked | index vs working directory |
| 06 | branch/switch/checkout | dirty worktree, detached HEAD | HEAD symbolic vs detached |
| 07 | diff/merge/conflict | merge base, fast-forward, conflict marker | snapshot-vs-diff, three-way merge |
| 08 | release demo | full command flow, docs examples | final Git internals essay |

## Backend Extension Artifact Matrix
| Extension | Must ship | Test focus | Docs you write |
|---|---|---|---|
| C | goroutine-safe cache + context-aware service | concurrent cache, cancel propagation, `go test -race` | learning notes: goroutine vs Isolate, done channel, context.Context |
| 09 | core service boundary | core tests without CLI process, domain error mapping | `docs/architecture.md` |
| 10 | REST API adapter | `httptest`, invalid JSON, not found, conflict mapping | `docs/api.md` |
| 11 | repo metadata store + workspace resolver | DB integration, storage resolver, multi-repo isolation | `docs/storage-design.md` |
| 12 | auth/logging/config/Docker/deploy | middleware, config, readiness, race test | `docs/deployment.md`, `docs/runbook.md` |

## Final Demo Flow
```sh
mgit init
echo hello > a.txt
mgit add a.txt
mgit commit -m "first"
mgit branch dev
mgit switch dev
echo dev > a.txt
mgit add a.txt
mgit commit -m "change on dev"
mgit switch main
echo main > a.txt
mgit add a.txt
mgit commit -m "change on main"
mgit merge dev
mgit status
```

## Learning Velocity Guide
Ghi chú cho Flutter dev 4 năm: biết cái nào cần học sâu, cái nào có thể đi nhanh hơn.

**KHÔNG bỏ qua (ROI cao nhất):**
- Phase 02: `io.Reader`, `[]byte`, `defer` — I/O foundation của Go
- Phase 04: struct design + file parsing + error wrapping — pattern bạn sẽ dùng mọi nơi
- Phase 05: `encoding/json`, `map`, `filepath.Rel` — backend data handling cơ bản
- Phase 06: filesystem atomicity, temp dir test pattern — test technique chứng minh bạn biết test
- **Concurrency Module**: goroutine + channel + context — **PHẢI xong trước Extension 09**, Go interview không bỏ qua
- Extension 10: `net/http`, middleware, `httptest` — đây là piece CV chính của track backend

**Có thể đi nhanh hơn:**
- Phase 03: khi đã hiểu blob (Phase 02), tree chỉ là recursion thêm — 1 ngày thay vì 2
- Phase 01: CLI scaffold đơn giản, Flutter dev quen OOP sẽ không bí
- Extension 09: nếu package design đã clean từ đầu theo plan, đây là refactor nhỏ, không phải rewrite

**Dấu hiệu bạn đang học vẹt:**
- Pass test nhưng không giải thích được tại sao goroutine không cần mutex ở chỗ đó
- Copy skeleton rồi implement mà không đọc Socratic Questions
- Xong phase nhưng learning notes trống

## Backend Extension Demo Flow
```sh
mgit-api

curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/repos
curl http://localhost:8080/repos/<repo-id>/status
curl -X POST http://localhost:8080/repos/<repo-id>/add
curl -X POST http://localhost:8080/repos/<repo-id>/commits
curl http://localhost:8080/repos/<repo-id>/log
curl http://localhost:8080/repos/<repo-id>/diff
```
