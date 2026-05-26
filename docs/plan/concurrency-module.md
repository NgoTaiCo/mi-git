# Concurrency Module: Goroutines, Channels & Context

> **Meta:** Bắt buộc hoàn thành trước Extension 09. Không thể viết HTTP handler đúng nếu không hiểu goroutine, channel và context.Context. Đây không phải "advanced topic" — đây là nền tảng của mọi Go service.
>
> **Nguyên tắc:** Code skeleton trong file này **không có ruột hàm**. Nó chỉ định hình data contract để bạn tự implement.
>
> **Vị trí trong lộ trình:** Sau Phase 08 (core Mini Git hoàn thành), trước Extension 09 (core boundary).

---

## Concurrency Sprint: Goroutines, Channels & Context

> **Mục tiêu:** Hiểu goroutine, channel, mutex và context.Context đủ để viết HTTP server đúng, test concurrent code, không leak goroutine và propagate cancel qua call chain.
>
> **Mini-git surface:** Object integrity scanner song song, log walker pipeline với channel, context-aware service method chuẩn bị cho HTTP handler.

---

## Module Overview

### Sessions
- Session C1 - Goroutine, WaitGroup và parallel object scanner
- Session C2 - Channel, select và log walker pipeline
- Session C3 - Mutex, context.Context và HTTP-ready service

### Flutter / Dart Bridge
> Go goroutine **không phải** Dart Isolate. Isolate là process riêng, không share heap, phải pass message qua `SendPort`. Goroutine share heap với nhau — hai goroutine có thể đọc/ghi cùng một `map` mà không có compile error, chỉ gây data race lúc runtime. `go test -race` mới bắt được điều đó. Channel **không phải** Dart Stream — channel là blocking pipe kiểu pull, Stream là async event emitter kiểu push. `sync.WaitGroup` gần nhất với `Future.wait(list)` của Dart. `context.Context` gần nhất với `CancelToken` hoặc `ref.onDispose` trong Riverpod: nó lan truyền signal cancel qua toàn bộ call stack và goroutine chain.

### Go Skills Required For This Module
> `go` keyword, `sync.WaitGroup` (Add/Done/Wait), `sync.Mutex` (Lock/Unlock), `sync.RWMutex` (RLock/RUnlock), `make(chan T)` unbuffered, `make(chan T, n)` buffered, channel send/receive (`ch <- val`, `<-ch`), `select` statement, `close(ch)`, `context.Background()`, `context.WithTimeout()`, `context.WithCancel()`, `context.WithDeadline()`, `defer cancel()`, `ctx.Done()`, `ctx.Err()`. `go test -race ./...` không còn optional từ module này.

---

## Session C1: Goroutine, WaitGroup và Parallel Object Scanner

### User Story
> Sau khi Mini Git core xong, bạn muốn thêm `mgit fsck` — quét toàn bộ `.mgit/objects`, verify SHA-1 của từng file. Làm tuần tự mất O(n) với n object file. Làm song song với goroutine pool có thể dùng được nhiều CPU core.
>
> **Mục tiêu học không phải performance.** Mục tiêu là hiểu goroutine chạy thế nào, WaitGroup block thế nào và race condition phát sinh từ đâu.

### CSV Main Task
Concurrency C1 - Goroutine basics, WaitGroup, parallel scanner, race detector

### Acceptance Criteria
- [ ] `mgit fsck` đọc toàn bộ object file trong `.mgit/objects`
- [ ] Spawn goroutine để verify SHA-1 cho mỗi object, có giới hạn số worker
- [ ] Dùng `sync.WaitGroup` để đợi toàn bộ goroutine hoàn thành
- [ ] Collect kết quả không race condition
- [ ] In số object valid và số object corrupt nếu có
- [ ] `go test -race ./...` xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Goroutine rẻ — ~2KB initial stack, không phải OS thread 2MB. Nhưng rẻ không có nghĩa là vô hạn. Spawn 100,000 goroutine cho 100,000 object là dở. Cần worker pool với bounded concurrency."
>
> "Bẫy số 1: append vào shared slice trong goroutine không có lock. `append` không thread-safe. Phải dùng mutex hoặc collect qua channel."
>
> "Bẫy số 2: WaitGroup.Add phải gọi TRƯỚC khi spawn goroutine. Nếu Add gọi bên trong `go func()`, có thể Wait() xong trước Add() chạy."
>
> "Go scheduler là M:N: M goroutine chạy trên N OS thread. GOMAXPROCS mặc định = số CPU. Goroutine yield tại blocking point: channel op, syscall, runtime.Gosched. Không yield trong tight loop → goroutine lấn át scheduler."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/fsck/fsck.go
package fsck

type ObjectReport struct{}
type ScanResult struct{}

func ScanObjects(objectDir string, workers int) (ScanResult, error)
func verifyOne(path string) (ObjectReport, error)

// SENIOR ASKS: Bạn dùng channel hay mutex+slice để collect ObjectReport từ goroutine?
//              Khi nào channel phù hợp hơn? Khi nào mutex+slice phù hợp hơn?
// SENIOR ASKS: Worker pool với workers=4 và 1000 object — số goroutine tối đa tại một thời điểm là bao nhiêu?
//              Kể cả goroutine main và goroutine worker.
```

#### Theory Notes
- [ ] Goroutine stack: bắt đầu ~2-4KB, grow on demand, shrink khi idle — không phải fixed 2MB như OS thread
- [ ] `sync.WaitGroup`: `Add(n)` trước spawn, `Done()` trong goroutine (dùng `defer`), `Wait()` block đến khi counter về 0
- [ ] Race detector: compile với `-race`, detect concurrent write/write và write/read trên cùng memory address
- [ ] GOMAXPROCS: số CPU core mặc định — `runtime.GOMAXPROCS(1)` force single-threaded để reproduce race bugs
- [ ] Worker pool pattern: send tasks vào buffered channel, N goroutine đọc từ channel đó

#### Socratic Questions
1. Tại sao `WaitGroup.Add(1)` PHẢI gọi trong goroutine gốc, không gọi bên trong `go func()`?
2. Goroutine khác OS thread ở điểm nào khiến race condition dễ xảy ra hơn trong Go so với Java?
3. Worker pool dùng buffered channel `jobs := make(chan string, N)`. N nên bằng bao nhiêu: số file, số worker, hay số CPU? Trade-off là gì?
4. `go test -race` bắt được gì mà `go test` bình thường không bắt? Tại sao không mặc định bật?

### Output Checklist
- [ ] `mgit fsck` chạy được, đếm đúng số object valid và corrupt
- [ ] Worker count có thể config (default = runtime.NumCPU())
- [ ] Test cover 0 object, 1 object, nhiều object concurrent

### Test Checklist
- [ ] `ScanObjects` với thư mục rỗng trả ScanResult rỗng, không error
- [ ] Test với 1 object valid và 1 object corrupt (ghi sai SHA-1 header)
- [ ] Test concurrent không race với `go test -race`
- [ ] `go test -race ./...` xanh

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng: goroutine khác Isolate/Thread thế nào, race detector bắt gì, và tại sao GOMAXPROCS quan trọng

### Retrospective
1. Bạn dùng channel hay mutex+slice để collect? Lý do?
2. Worker pool hay goroutine-per-task? Tại sao với object store của mini-git thì cái nào hợp lý hơn?
3. WaitGroup.Add đặt ở đâu trong code của bạn? Có thể đặt trong goroutine không?

---

## Session C2: Channel, Select và Log Walker Pipeline

### User Story
> `mgit log` hiện tại đọc toàn bộ commit rồi mới print — người dùng phải đợi load xong. Bạn sẽ refactor thành streaming pipeline: một goroutine walk DAG và gửi vào channel, caller nhận và print ngay khi có commit đầu tiên.
>
> **Quan trọng hơn pipeline:** bạn học cách goroutine leak xảy ra và cách done channel phòng tránh. Đây là pattern nền tảng — sau này `context.Context` chỉ là abstract layer trên done channel.

### CSV Main Task
Concurrency C2 - Channel pipeline, done pattern, select, goroutine leak prevention

### Acceptance Criteria
- [ ] `WalkLog` gửi commit vào channel theo thứ tự từ HEAD về root
- [ ] Caller đọc từng commit khi nó sẵn sàng, không đợi toàn bộ DAG load
- [ ] Có done channel để cancel sớm (ví dụ `mgit log -n 5` không cần đọc 10,000 commit)
- [ ] Goroutine walk không bị leak khi caller cancel sớm
- [ ] Dùng `select` để handle cả send-to-commits và receive-from-done đồng thời
- [ ] `go test -race ./...` xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Unbuffered channel: sender block cho đến khi receiver nhận. Đây là synchronization point implicit. Buffered(1): sender block chỉ khi buffer đầy — decouples timing nhẹ."
>
> "Goroutine leak là bug phổ biến nhất của Go junior: goroutine gửi vào channel nhưng không ai nhận nữa — goroutine stuck mãi mãi, heap không được GC. Done channel là cách signal goroutine thoát clean."
>
> "`select` là switch cho channel: check nhiều channel cùng lúc, pick case nào ready trước. Case `default` làm select non-blocking. Không có `default`: block cho đến khi ít nhất một case ready."
>
> "Sau này khi học context.Context ở C3, bạn sẽ thấy `ctx.Done()` chính là done channel được wrap thêm timeout/deadline."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/history/walk.go  (file mới trong package đã có từ Phase 04)
// Phase 04 tạo log.go với func Log() trả []commit.Commit — sequential.
// Session C2 thêm walk.go vào cùng package với variant concurrent.
package history

// WalkLog gửi commit vào commits channel theo thứ tự từ HEAD về root.
// Dừng sớm nếu done channel bị close.
// Caller chịu trách nhiệm close(done) nếu muốn cancel.
func WalkLog(
    store object.Store,
    start object.ObjectID,
    commits chan<- commit.Commit,
    done <-chan struct{},
) error

// SENIOR ASKS: Nếu caller close(done) trong khi WalkLog đang block send vào commits,
//              điều gì xảy ra với goroutine WalkLog nếu không có select?
// SENIOR ASKS: Unbuffered vs buffered(1) cho commits — behavior nào khác nhau khi caller xử lý chậm?
//              Khi nào mỗi loại gây vấn đề?
```

#### Theory Notes
- [ ] Unbuffered `make(chan T)`: send block cho đến khi receiver sẵn sàng — đây là handshake implicit
- [ ] Buffered `make(chan T, n)`: send chỉ block khi buffer đầy — decouples sender/receiver timing
- [ ] Done channel pattern: `done <-chan struct{}` nhận signal cancel từ ngoài — `struct{}` vì zero size
- [ ] `close(ch)` vs `ch <- struct{}{}`: close broadcast tới MỌI receiver; send chỉ 1 receiver nhận
- [ ] `select` với 2 case ready cùng lúc: Go chọn ngẫu nhiên — không phải FIFO, không phải priority
- [ ] Goroutine leak: block mãi trên channel op, không có cơ chế thoát → GC không thể collect → OOM dần

#### Socratic Questions
1. Tại sao goroutine leak không bị Go runtime tự phát hiện như memory leak trong nhiều GC languages?
2. `close(done)` khác `done <- struct{}{}` như thế nào với nhiều goroutine cùng wait trên done?
3. Select có 2 case ready đồng thời — Go chọn cái nào? Điều này ảnh hưởng thế nào tới fairness?
4. Sau session này, bạn thay `done <-chan struct{}` bằng `ctx context.Context` — method nào của Context map sang done channel?

### Output Checklist
- [ ] `WalkLog` streaming commit qua channel
- [ ] `mgit log -n 5` dừng sau 5 commit, không đọc toàn bộ DAG
- [ ] Test cancel sớm không leak goroutine

### Test Checklist
- [ ] WalkLog với 0 commit — channel close ngay, không gửi gì
- [ ] WalkLog với done close ngay sau khi gọi — goroutine thoát clean
- [ ] WalkLog bình thường — nhận đủ N commit rồi caller cancel
- [ ] `go test -race ./...` xanh

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng: goroutine leak là gì, tại sao runtime không tự phát hiện, cách done channel phòng tránh

### Retrospective
1. Bạn dùng `close(done)` hay `done <- struct{}{}`? Trade-off với nhiều goroutine cùng wait?
2. Buffered(1) hay unbuffered cho commits? Bạn chọn gì và vì sao?
3. Nếu WalkLog return error giữa chừng, channel commits đã được close chưa? Caller biết bằng cách nào?

---

## Session C3: Mutex, context.Context và HTTP-Ready Service

### User Story
> Trước khi viết HTTP handler (Extension 10), service layer phải đáp ứng 3 điều kiện:
> 1. Object lookup cần cache — cache phải goroutine-safe (nhiều handler request cùng lúc)
> 2. Mọi operation phải có thể bị cancel/timeout bởi HTTP request lifecycle
> 3. Service phải safe to call concurrently vì mỗi HTTP request là một goroutine riêng
>
> Session này là cầu nối giữa concurrency fundamentals và backend. Sau đây bạn viết handler mà biết mình đang làm gì.

### CSV Main Task
Concurrency C3 - Mutex, RWMutex, context.Context, goroutine-safe service method

### Acceptance Criteria
- [ ] In-memory object cache implement với `sync.RWMutex` (nhiều reader, ít writer)
- [ ] Service methods nhận `context.Context` làm tham số đầu tiên, tên `ctx`
- [ ] Long operation check `ctx.Done()` để cancel sớm với `ctx.Err()`
- [ ] Test cancel propagation: ctx timeout/cancel → operation dừng sớm, trả error rõ
- [ ] Không dùng `ctx.Value` để truyền business data (chỉ dùng cho request-scoped metadata như trace id)
- [ ] Mọi `context.WithTimeout` và `context.WithCancel` đều có `defer cancel()`
- [ ] `go test -race ./...` xanh

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "context.Context là interface 4 method. Quan trọng nhất: Done() trả chan struct{} đóng khi cancel hoặc timeout. Sau khi học channel ở C2, ctx.Done() không còn magic — nó chỉ là done channel được Go standard library wrap thêm timeout/deadline."
>
> "Quy tắc: context là tham số đầu tiên, tên là ctx, không bao giờ nil (truyền context.Background() thay vì nil). Không nhét vào struct — context là per-call, không phải per-service lifetime."
>
> "RWMutex: nhiều goroutine có thể RLock cùng lúc (concurrent reads), nhưng Lock exclusive (single writer blocks all). Dùng khi read >> write."
>
> "context.Value anti-pattern: không dùng để truyền user ID, config hay business data. Chỉ dùng cho cross-cutting system data: request id, trace id, auth principal. Lý do: type-unsafe, hidden dependency, khó test."
```

#### TODO Comments (Skeleton / Contract Only)
```go
// File: internal/object/cache.go
package object

type Cache struct{}

func NewCache() *Cache
func (c *Cache) Get(id ObjectID) (Object, bool)
func (c *Cache) Set(id ObjectID, obj Object)

// SENIOR ASKS: Tại sao dùng *Cache (pointer receiver) thay vì Cache (value receiver)?
//              Nếu dùng value receiver và copy Cache, điều gì xảy ra với mutex bên trong?

// File: internal/core/service.go
// NOTE: Đây là DRAFT để luyện context.Context propagation pattern.
// Extension 09 sẽ redesign package này thành Service INTERFACE với domain errors,
// input/output structs và full boundary. Xem extension-09-core-boundary.md.
package core

type Service struct{}
type StatusReport struct{}

func NewService(store object.Store, cache *object.Cache) *Service
func (s *Service) ReadObject(ctx context.Context, id object.ObjectID) (object.Object, error)
func (s *Service) GetStatus(ctx context.Context, repoPath string) (StatusReport, error)
func (s *Service) AddFile(ctx context.Context, repoPath string, relPath string, content []byte) error

// SENIOR ASKS: Tại sao ctx là tham số, không phải field trong Service struct?
//              Service được tạo một lần, nhưng mỗi HTTP request có lifetime khác nhau.
// SENIOR ASKS: `select { case <-ctx.Done(): return nil, ctx.Err() }` — ctx.Err() trả gì
//              khi timeout? Khi cancel? Handler phân biệt hai trường hợp này để làm gì?
```

#### Theory Notes
- [ ] `sync.Mutex`: Lock/Unlock — chỉ một goroutine hold lock tại một lúc; dùng khi data có cả read và write
- [ ] `sync.RWMutex`: RLock/RUnlock cho concurrent reads, Lock/Unlock cho exclusive write — dùng khi read >> write
- [ ] `sync.Mutex` trong struct: KHÔNG được copy (go vet cảnh báo) — luôn dùng pointer receiver hoặc truyền pointer
- [ ] `context.Context` interface: `Deadline()`, `Done()`, `Err()`, `Value(key)` — Done() là chan struct{} close khi cancel/timeout
- [ ] `context.WithTimeout(parent, d)`: cancel tự động sau duration d, hoặc khi parent cancel — cái nào trước thì thắng
- [ ] `defer cancel()`: bắt buộc sau mọi `context.WithTimeout`/`WithCancel` để giải phóng goroutine timer nội bộ
- [ ] Convention: `ctx context.Context` là tham số đầu tiên mọi function có I/O hoặc blocking op — không phải field trong struct

#### Socratic Questions
1. `context.WithTimeout(parent, 5*time.Second)` — nếu parent bị cancel sau 1 giây, child timeout sau mấy giây thực tế?
2. `sync.Mutex` bị copy qua value receiver — go vet cảnh báo gì? Điều gì thực sự xảy ra lúc runtime?
3. `ctx.Err()` trả `context.Canceled` vs `context.DeadlineExceeded` — HTTP handler xử lý hai case này khác nhau thế nào?
4. Tại sao `context.Value` nên dùng custom key type (không phải `string`) làm key?

### Output Checklist
- [ ] `Cache` thread-safe với nhiều goroutine đọc đồng thời (RLock) và ít goroutine ghi (Lock)
- [ ] `ReadObject` nhận ctx, dừng sớm nếu ctx expired
- [ ] Service method signatures chuẩn: `ctx` đầu tiên, `error` cuối
- [ ] Mọi `context.WithTimeout/Cancel` đều có `defer cancel()` ngay bên dưới

### Test Checklist
- [ ] Cache.Get/Set với 10 goroutine concurrent không race
- [ ] ReadObject với canceled context trả error ngay, không block
- [ ] ReadObject với valid context trả object đúng
- [ ] `go test -race ./...` xanh

### Learning Notes / Docs
- [ ] Viết ít nhất 5 dòng: context.Context là gì, tại sao là tham số không phải field, khi nào dùng WithTimeout vs WithCancel vs WithDeadline

### Retrospective
1. Service struct của bạn có field nào không goroutine-safe không? List ra và giải thích cách fix.
2. Cache dùng `sync.Map` (built-in) thay `map + RWMutex` — trade-off performance và readability là gì?
3. HTTP handler sẽ gọi `s.GetStatus(r.Context(), repoPath)` — `r.Context()` đến từ đâu và cancel khi nào?

---

## Module Checkpoints (BẮT BUỘC)

### CP-C-A: Concepts Check
- [ ] Chạy được `mgit fsck` với worker pool, output rõ số object valid/corrupt.
- [ ] `mgit log -n 5` dừng sau 5 commit, không load toàn bộ DAG.
- [ ] Service method nhận ctx, trả error khi ctx cancel.

### CP-C-B: Test Gate
- [ ] `go test -race ./...` xanh toàn bộ module — không một race nào.
- [ ] Có test concurrent cache (N goroutine đọc/ghi đồng thời).
- [ ] Có test cancel propagation: ctx timeout → operation return sớm với error rõ.

### CP-C-C: Oral Defense
- [ ] Giải thích được goroutine leak — tại sao xảy ra, bằng chứng từ code của bạn không bị leak.
- [ ] Giải thích được: tại sao ctx là tham số không phải field trong struct.
- [ ] Giải thích được: ctx.Done() liên hệ thế nào với done channel ở Session C2.

## Failure Modes (PHẢI BIẾT)
- WaitGroup.Add gọi bên trong goroutine → Wait() có thể return trước goroutine chạy.
- Goroutine gửi vào channel khi không còn receiver → leak mãi mãi.
- Mutex bị copy qua value receiver → hai bản mutex độc lập, không protect gì.
- `defer cancel()` bị quên → goroutine timer nội bộ của context không được giải phóng → resource leak.
- `ctx.Value` để truyền user ID hoặc config → hidden dependency, khó test, type-unsafe.

## Progression Rules

### Rule 1: Không qua Extension 09 nếu chưa xanh `go test -race ./...`.
Race condition trong core service sẽ gây bug ngẫu nhiên trong HTTP handler — rất khó debug sau.

### Rule 2: Không bỏ qua Session C2 done channel.
`context.Context` ở C3 sẽ không có ý nghĩa nếu bạn không hiểu done channel là gì.

### Rule 3: Không copy code skeleton — tự viết ruột hàm.
Goroutine/channel/mutex là thứ cần viết tay mới thấm. Copy không giúp gì.

## First-Principles Question
> HTTP server của Go spawn một goroutine mới cho **mỗi** request. Khi 100 user đồng thời gửi request, có 100 goroutine cùng gọi vào `s.GetStatus()`. Nếu `GetStatus` đọc một `map` trong Service struct mà không có mutex, race condition xảy ra.
>
> **Câu hỏi:** `s.GetStatus` được gọi từ 100 goroutine concurrently. Function cần: (1) đọc object từ cache, (2) đọc filesystem `.mgit/objects`, (3) ghi vào cache nếu miss. Bạn phải protect gì bằng mutex? Bạn không cần protect gì? Filesystem read có cần lock không?
>
> Trả lời bằng lời trước khi qua Extension 09.
