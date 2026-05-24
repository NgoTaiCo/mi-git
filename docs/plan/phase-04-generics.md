# Phase 4: Generics & Type Safety (Tuần 8)

> **Quote from the trenches:** "Trước Go 1.18, tôi viết 5 hàm Max khác nhau. Sau Go 1.18, tôi viết 1 hàm Max generic... rồi 6 tháng sau refactor lại thành 2 hàm vì over-engineering. Generics là dao hai lưỡi." — Senior Go Engineer

---

## Phase Context

| | |
|---|---|
| **Duration** | 1 week (Tuần 8) |
| **Goal** | Replace `interface{}` where type safety matters; avoid over-engineering with generics |
| **minigit Link** | Refactor `store.go` thành generic `ObjectStore[T Object]` |
| **Prerequisites** | Phase 1-3 hoàn thành, hiểu interface, type system, slice, map |

**Warning — Failure Modes cần tránh:**
- Over-engineering: forcing generics vào code đơn giản
- `comparable` không include slice/map/function — biết trước để không shock
- Generic constraints quá phức tạp → code đọc không nổi

---

## Topic 04.1: Type Parameters & Constraints

### User Story

> **Khách hàng (Product Owner) nói:** "Tôi cần hàm `Max` hoạt động được với cả `int` và `float64`. Không dùng `interface{}` vì mất type safety — tôi đã từng bị runtime panic vì type assert sai ở production."
>
> **Context:** Team đang viết thư viện util chung. Trước đâp dùng `interface{}` + type assertion, đã gây 3 incident trong 6 tháng. Yêu cầu: compile-time type safety, không runtime panic vì type mismatch.

### Acceptance Criteria

- [ ] Generic `Max[T]` hoạt động với `int`, `float64`, `string` (string so sánh lexicographic)
- [ ] Hiểu sự khác biệt `any` vs `comparable` — khi nào dùng cái nào
- [ ] Type inference hoạt động: gọi `Max(3, 5)` không cần viết `Max[int](3, 5)`
- [ ] Không dùng `interface{}` hoặc reflection trong implementation
- [ ] Benchmark: so sánh performance generic vs non-generic (inline function)

### Senior Thought-Process

**Senior nghĩ gì khi nhận requirement này:**

> "Nếu tôi nhận ticket này, điều đầu tiên tôi nghĩ đến là: *bài toán cốt lõi ở đây là ordering* — tức là type `T` phải "so sánh được" bằng `<`, `>`. Không phải type nào cũng làm được điều đó. `any` thì quá rộng — `struct{}` không so sánh được. `comparable` thì chỉ đảm bảo `==` và `!=`, không có `<`.
>
> Hồi tôi ở project fintech, tôi gặp vấn đề này khi viết hàm `Clamp`, `Min`, `Max` cho pricing engine. Cách tôi xử lý là định nghĩa một custom constraint `Ordered` — học từ package `golang.org/x/exp/constraints`. Đây là pattern industry standard.
>
> Tôi sẽ phân rã thành các bước:
> 1. **Constraint design**: T type cần gì? → `<`, `<=`, `>`, `>=`
> 2. **Type parameter declaration**: `func Max[T constraints.Ordered](a, b T) T`
> 3. **Type inference testing**: đảm bảo caller không phải gõ `[T]`
> 4. **Edge case audit**: `NaN` với float? String rỗng?"

### TODO Comments (Code Skeleton)

```go
package generics

import "golang.org/x/exp/constraints"

// TODO-[1]: Định nghĩa constraint Ordered cho phép so sánh < >
// SENIOR ASKS: comparable đã đủ chưa? Nếu chưa, thiếu gì?
// HINT: comparable chỉ đảm bảo == và !=. So sánh thứ tự cần gì khác?

// TODO-[2]: Viết hàm Max[T Ordered]
// SENIOR ASKS: Tại sao không dùng any làm constraint?
// HINT: any cho phép mọi type — kể cả type không so sánh được. Điều gì xảy ra nếu gọi Max(struct{}{}, struct{}{})?

// TODO-[3]: Viết hàm Min[T Ordered]
// SENIOR ASKS: Có nên dùng type parameter cho return type không?
// HINT: Return type cùng là T — nhưng điều gì nếu a và b khác type underlying?

// TODO-[4]: Viết hàm Clamp[T Ordered]
// SENIOR ASKS: Logic Clamp thế nào? Xử lý NaN?
// HINT: Clamp(v, min, max) = Max(min, Min(v, max)). Nhưng thứ tự parameter quan trọng.

// TODO-[5]: Test type inference — gọi không cần [T]
// SENIOR ASKS: Khi nào compiler inference được, khi nào không?
// HINT: Nếu arguments không đủ thông tin để suy luận T — phải gõ explicit.
```

```go
// TODO-[6]: Tự định nghĩa constraint Number (int + float unions)
// SENIOR ASKS: ~ operator có cần ở đây không? Hay chỉ cần union type?
// HINT: Nếu constraint là int | float64 — type int32 có satisfy không?
```

### Socratic Questions

**Câu hỏi để bạn tự suy nghĩ:**

1. **Nếu constraint là `comparable`,** bạn có thể viết `Max` không? `comparable` cho phép `==` nhưng không cho phép `<`. Bạn có thể "hack" để tìm max chỉ dùng `==` không?

2. **Type inference failure:** Khi nào Go compiler KHÔNG thể inference `T`? Thử nghĩ case: `result := Max(3, 3.14)` — compile không?

3. **Performance trade-off:** Generic function có bị chậm hơn function non-generic không? Tại sao? (Hint: monomorphization vs boxing)

4. **Constraint granularity:** Tại sao `constraints.Ordered` include cả `string`? Trong pricing engine, bạn có muốn `Max("apple", "banana")` compile được không?

### Output Checklist

- [ ] TODO-[1] hoàn thành: Constraint `Ordered` (hoặc dùng từ `x/exp/constraints`) được define rõ
- [ ] TODO-[2] hoàn thành: `Max[T Ordered]` hoạt động với int, float64, string
- [ ] TODO-[3] hoàn thành: `Min[T Ordered]` hoạt động tương tự
- [ ] TODO-[4] hoàn thành: `Clamp[T Ordered]` xử lý đúng boundary
- [ ] TODO-[5] hoàn thành: Demonstrate type inference — ít nhất 3 ví dụ
- [ ] TODO-[6] hoàn thành: Custom `Number` constraint với type union

### Test Checklist

- [ ] **Test case:** `Max(3, 5) == 5` (int, basic case) — vì sao case này quan trọng? → Sanity check
- [ ] **Test case:** `Max(-1.5, -3.2) == -1.5` (float64, negative) — boundary với negative numbers
- [ ] **Test case:** `Max("apple", "banana") == "banana"` (string, lexicographic) — chứng minh Ordered include string
- [ ] **Test case:** `Clamp(10, 0, 5) == 5` (clamp to max) — boundary case
- [ ] **Test case:** `Clamp(-5, 0, 10) == 0` (clamp to min) — boundary case
- [ ] **Test case:** Gọi `Max` với 2 type khác nhau (e.g., int và int64) — có compile không? Tại sao đây là behavior mong muốn?

### Retrospective

**Sau khi xong, hãy tự hỏi:**

1. **Nếu tôi bắt buộc chỉ dùng `comparable` (không có `<`),** tôi có thể implement `Max` không? Cách nào? Trade-off là gì?

2. **Nếu requirement thay đổi:** "Max cũng phải hoạt động với `time.Time`" — `time.Time` có satisfy `Ordered` không? Nếu không, approach nào?

3. **Over-engineering check:** Viết `Max[T]` có quá mức cho 1 project chỉ dùng `int` không? Khi nào YAGNI, khi nào future-proof?

---

## Topic 04.2: Generic Types (Struct)

### User Story

> **Khách hàng (Product Owner) nói:** "Xây dựng Cache generic: key là string, value có thể là bất kỳ type nào. Có TTL expiration. Trước đây team dùng `map[string]interface{}` rồi type-assert khắp nơi — 40% bug report là panic từ việc này."
>
> **Context:** Service có 15 loại entity khác nhau cần cache. Hiện tại mỗi entity có 1 map riêng → 15 maps, code duplicate. Yêu cầu: 1 cache type-safe, hỗ trợ TTL tự động expire, benchmark chứng minh performance tốt hơn hoặc bằng `map[string]interface{}`.

### Acceptance Criteria

- [ ] `Cache[K comparable, V any]` hoạt động type-safe — không cần type assertion khi Get
- [ ] Hỗ trợ TTL expiration: item tự động bị loại bỏ sau duration
- [ ] Thread-safe: concurrent access không panic/race
- [ ] Benchmark: so sánh `Cache[string, *User]` vs `map[string]interface{}` + type assertion
- [ ] API: `Set(key K, value V, ttl time.Duration)`, `Get(key K) (V, bool)`, `Delete(key K)`, `Len() int`

### Senior Thought-Process

**Senior nghĩ gì khi nhận requirement này:**

> "Generic cache vs `map[string]interface{}` — speed khác nhau bao nhiêu? Đây là câu hỏi đầu tiên tôi nghĩ đến. Thực tế: generic cache *nhanh hơn* vì không cần boxing/unboxing qua interface{} và không cần type assertion.
>
> Ở project trước, tôi viết cache layer cho API gateway. Học được bài học đau: `interface{}` cache có vẻ tiện nhưng mỗi lần Get là 1 lần type assertion — không chỉ chậm mà còn dễ crash. Một junior từng ghi `cache.Set("user_123", user)` rồi ở chỗ khác `cache.Get("user_123").(*Admin)` — panic vì wrong type assertion.
>
> Với generic cache, lỗi đó sẽ là compile error: `Cache[string, *User]` không cho phép `.(*Admin)`.
>
> Tôi sẽ phân rã:
> 1. **Storage**: `map[K]*cacheEntry[V]` — underlying map dùng generic key và value
> 2. **TTL**: mỗi entry có `expiresAt time.Time`, lazy eviction (kiểm tra khi Get) + background cleanup
> 3. **Thread-safety**: `sync.RWMutex` — read-heavy workload cần RLock
> 4. **Benchmark**: table so sánh generic vs interface{} với cùng workload"

### TODO Comments (Code Skeleton)

```go
package cache

import (
	"sync"
	"time"
)

// TODO-[1]: Định nghĩa cacheEntry[V any] — struct lưu value + expiration
// SENIOR ASKS: Tại sao dùng pointer *cacheEntry thay vì value cacheEntry?
// HINT: Pointer giảm memory khi copy? Hay có lý do khác?

type cacheEntry[V any] struct {
	// TODO: value V và expiresAt time.Time
}

// TODO-[2]: Định nghĩa Cache[K comparable, V any] struct
// SENIOR ASKS: K comparable — tại sao K phải comparable?
// HINT: underlying storage là map[K]... — map key type cần satisfy điều kiện gì?

type Cache[K comparable, V any] struct {
	// TODO: mu sync.RWMutex, data map[K]*cacheEntry[V]
}

// TODO-[3]: Viết constructor New[K, V]()
// SENIOR ASKS: Khi nào nên dùng constructor function cho generic type?
// HINT: Generic type literal: Cache[string, int]{} — đã đủ rõ chưa?

// TODO-[4]: Implement Set(key K, value V, ttl time.Duration)
// SENIOR ASKS: TTL = 0 nghĩa là gì? Không expire hay expire ngay?
// HINT: Định nghĩa behavior rõ: ttl <= 0 = no expiration? Hay reject?

// TODO-[5]: Implement Get(key K) (V, bool)
// SENIOR ASKS: Làm sao trả về "zero value" của V khi không tìm thấy?
// HINT: var zero V — generic zero value. Đừng dùng nil vì V có thể là value type.

// TODO-[6]: Lazy eviction trong Get — nếu entry expired, xóa và trả not-found
// SENIOR ASKS: Tại sao LAZY eviction thay vì background goroutine?
// HINT: Goroutine cleanup có thể gây goroutine leak nếu cache bị bỏ quên. Lazy = simple.

// TODO-[7]: Implement Delete và Len
// SENIOR ASKS: Len đếm active entries hay tất cả entries (kể cả expired)?
// HINT: Behavior nào surprise ít hơn cho caller?
```

```go
// TODO-[8]: Benchmark — Cache[string, int] vs map[string]interface{}
// SENIOR ASKS: Benchmark nên đo gì? Set, Get, hay cả hai?
// HINT: Dùng b.N loop, testing.B. So sánh cả latency và allocation (benchmem).
```

### Socratic Questions

**Câu hỏi để bạn tự suy nghĩ:**

1. **Nếu K là `string`,** tại sao cần `comparable`? String đã có `==` rồi mà? → `comparable` là tập hợp các type có `==` — string satisfy tự động. Nhưng constraint này báo cho ngườI đọc biết điều gì?

2. **`V any` vs `V comparable`:** Value type có cần `comparable` không? Trong cache, bạn có bao giờ so sánh 2 value không? Nếu không, tại sao `any` là đúng?

3. **Memory layout:** `Cache[string, int]` và `Cache[string, float64]` — Go compiler tạo 2 type riêng biệt (monomorphization) hay 1 type chung? Điều này ảnh hưởng binary size thế nào?

4. **Thiết kế alternative:** Thay vì `Cache[K comparable, V any]`, bạn có thể dùng `Cache[V any]` với key luôn là string? Trade-off? Khi nào K generic có giá trị thực sự?

### Output Checklist

- [ ] TODO-[1] hoàn thành: `cacheEntry[V any]` struct với value + expiresAt
- [ ] TODO-[2] hoàn thành: `Cache[K comparable, V any]` struct với RWMutex + map
- [ ] TODO-[3] hoàn thành: Constructor function hoạt động đúng
- [ ] TODO-[4] hoàn thành: Set với TTL behavior rõ ràng (documented)
- [ ] TODO-[5] hoàn thành: Get trả zero value + bool khi miss
- [ ] TODO-[6] hoàn thành: Lazy eviction hoạt động trong Get
- [ ] TODO-[7] hoàn thành: Delete và Len đếm đúng
- [ ] TODO-[8] hoàn thành: Benchmark có so sánh với `map[string]interface{}`

### Test Checklist

- [ ] **Test case:** Set + Get cùng key trả đúng value — sanity check
- [ ] **Test case:** Get key không tồn tại trả `(zero, false)` — miss handling
- [ ] **Test case:** Item hết hạn sau TTL → Get trả miss — core TTL functionality
- [ ] **Test case:** Set với TTL=0 → item không expire (hoặc reject, tùy design) — boundary
- [ ] **Test case:** Concurrent Set + Get từ 10 goroutines — race detection (`go test -race`)
- [ ] **Test case:** Delete key → Get sau đó trả miss — deletion works
- [ ] **Test case:** Len sau khi một số item expired — counts active entries only
- [ ] **Test case:** Cache với value type là struct (không pointer) — đảm bảo `V any` hoạt động

### Retrospective

**Sau khi xong, hãy tự hỏi:**

1. **Trade-off RWMutex vs Mutex:** `RWMutex` tốt cho read-heavy, nhưng overhead cao hơn `Mutex` ở low contention. Bạn có benchmark để chứng minh RWMutex là đúng choice không?

2. **Nếu requirement thay đổi:** "Cache cần hỗ trợ LRU eviction khi đầy" — bạn sẽ modify struct thế nào? Generic type có giúp hay cản trở?

3. **Binary size concern:** Mỗi instantiation `Cache[string, int]`, `Cache[string, User]`, `Cache[int, Order]` tạo ra code riêng. Trong microservice deploy bằng container, điều này có đáng lo không?

---

## Topic 04.3: Type Sets & ~ Operator

### User Story

> **Khách hàng (Product Owner) nói:** "Tôi cần hàm `Sum` hoạt động với cả `int`, `int64`, `uint`. Nhưng KHÔNG hoạt động với `string`. Team từng dùng `interface{}` rồi type-switch, nhưng mỗi lần thêm numeric type là phải sửa 3 chỗ."
>
> **Context:** Pricing engine tính tổng order value. Dùng `int64` cho số tiền (xu's). Nhưng 1 microservice cũ dùng `int` (legacy), 1 service khác dùng `uint` (counting). Cần 1 hàm `Sum` hoạt động với cả 3 mà vẫn type-safe.

### Acceptance Criteria

- [ ] Custom constraint `Number` include `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`
- [ ] `~` operator được dùng để cho phép named types có underlying type là numeric
- [ ] `Sum[T Number](values []T) T` hoạt động với slice bất kỳ numeric type
- [ ] `Sum[string]` phải là COMPILE ERROR — string không satisfy Number
- [ ] `Sum` hoạt động với custom named type: `type Money int64; Sum([]Money{10, 20})`

### Senior Thought-Process

**Senior nghĩ gì khi nhận requirement này:**

> "Type sets và `~` operator — đây là phần 'nâng cao' của generics mà nhiều ngườI học nhưng ít ngườI thực sự hiểu. Câu hỏi cốt lõi: tại sao cần `~int` thay vì chỉ `int`?
>
> Hồi tôi refactor pricing engine, tôi có `type Money int64`. Nếu constraint chỉ là `int64`, thì `Money` không satisfy — vì `Money` là named type, không phải `int64` trực tiếp. `~int64` nói: *và cả các type có underlying type là int64*. Đây là điểm then chốt.
>
> Tôi phân rã:
> 1. **Type set declaration**: `Number` = union tất cả numeric types + `~` cho mỗi cái
> 2. **Sum function**: iterate + accumulate, return zero value nếu slice rỗng
> 3. **Verify compile-time rejection**: thử `Sum("hello")` → phải fail compile
> 4. **Named type test**: `type Money int64; Sum([]Money{...})` → phải work
>
> Cái này giống vấn đề tôi gặp ở project X: cần viết hàm aggregate (Sum, Avg, Min, Max) hoạt động với mọi numeric type kể cả custom. `~` operator là chìa khóa."

### TODO Comments (Code Skeleton)

```go
package numeric

// TODO-[1]: Định nghĩa Number constraint với ~ operator
// SENIOR ASKS: Viết ~int64 | ~int32 — khác gì với int64 | int32 (không có ~)?
// HINT: Named type type Money int64 — có satisfy int64 constraint không? Có satisfy ~int64 không?

// Number là constraint cho mọi numeric type, kể cả named types dựa trên chúng.
type Number interface {
	// TODO: union với ~ cho mỗi type
}

// TODO-[2]: Viết Sum[T Number](values []T) T
// SENIOR ASKS: Làm sao init accumulator với zero value của T?
// HINT: var sum T — đủ chưa? Tại sao không cần gán explicit?

// TODO-[3]: Xử lý slice rỗng — return zero value
// SENIOR ASKS: Zero value của T generic lấy bằng cách nào?
// HINT: var result T → zero value. Khác gì với new(T) hay &T{}?

// TODO-[4]: Viết Avg[T Number] — average = Sum / len
// SENIOR ASKS: Tổng chia cho int — type mismatch?
// HINT: len(values) là int. T / int → compile? Cần convert len sang T không?
```

```go
// TODO-[5]: Verify named type works: type Money int64; Sum([]Money{10, 20, 30})
// SENIOR ASKS: Nếu bỏ ~ operator đi, điều gì xảy ra?
// HINT: Compile error: Money does not satisfy Number. Underlying type != type itself.

// TODO-[6]: Verify string REJECTED at compile time
// SENIOR ASKS: Điều gì xảy ra nếu thử Sum[string]([]string{"a", "b"})?
// HINT: Nên là compile error — đây chính là value proposition của generics.
```

### Socratic Questions

**Câu hỏi để bạn tự suy nghĩ:**

1. **`~` operator:** `~int64` bao gồm những type nào? `int64`, `Money` (type Money int64), `type Price int64` — cả 3 satisfy không? Type `int32` satisfy `~int64` không?

2. **Type set vs interface:** Trước Go 1.18, interface là method set. Sau Go 1.18, interface cũng có thể là type set. Điều này có nghĩa là: `interface{ ~int }` là valid. Nhưng bạn có thể dùng nó như một regular interface không? (e.g., `var x interface{ ~int } = 5`)

3. **Union type quá rộng:** `Number` constraint bao gồm cả `float64`. `Sum([]float64{0.1, 0.2})` kết quả là gì? Generic giải quyết được floating point precision issue không?

4. **Không dùng ~:** Nếu bạn chỉ dùng `int | int64 | uint` (không có `~`), việc dùng `type Money int64` bị ảnh hưởng thế nào? Đây là lý do chính để dùng `~` trong production code.

### Output Checklist

- [ ] TODO-[1] hoàn thành: `Number` constraint với `~` cho tất cả numeric types
- [ ] TODO-[2] hoàn thành: `Sum[T Number]` accumulate đúng
- [ ] TODO-[3] hoàn thành: Xử lý slice rỗng → zero value
- [ ] TODO-[4] hoàn thành: `Avg[T Number]` hoạt động đúng
- [ ] TODO-[5] hoàn thành: Named type `Money int64` hoạt động với `Sum`
- [ ] TODO-[6] hoàn thành: Chứng minh `Sum[string]` bị reject at compile time

### Test Checklist

- [ ] **Test case:** `Sum([]int{1, 2, 3}) == 6` — basic int
- [ ] **Test case:** `Sum([]int64{1, 2, 3}) == 6` — int64, verify type preservation
- [ ] **Test case:** `Sum([]uint{1, 2, 3}) == 6` — unsigned type
- [ ] **Test case:** `Sum([]float64{1.5, 2.5}) == 4.0` — floating point
- [ ] **Test case:** `Sum([]Money{10, 20})` với `type Money int64` — named type, ~ operator test
- [ ] **Test case:** `Sum([]T{})` (empty slice) → zero value — boundary
- [ ] **Test case:** `Avg([]int{1, 2, 3, 4}) == 2` (hoặc 2.5?) — watch out: integer division!

### Retrospective

**Sau khi xong, hãy tự hỏi:**

1. **~ operator trade-off:** Dùng `~` cho phép named types nhưng có thể cho phép *quá nhiều*. Nếu `type UserID int64` satisfy `~int64`, và bạn vô tình `Sum([]UserID{1, 2})` — điều này có ý nghĩa business không?

2. **Nếu requirement thay đổi:** "Cũng cần Sum cho `complex64` và `complex128`" — modify constraint thế nào? `complex` có operator `<` không? Nếu không, bạn vẫn có thể dùng cùng constraint không?

3. **Generic vs code generation:** Một approach khác là dùng `go generate` để tạo `SumInt`, `SumInt64`, `SumFloat64`. So sánh 2 approach: generic vs code generation. Khi nào cái nào thắng?

---

## Mini-Project: Generic Object Store

### User Story

> **Khách hàng (Tech Lead) nói:** "Refactor storage layer để type-safe. Trước dùng `interface{}` phải type-assert khắp nơi — 40% bug là từ đó. Tất cả entities (User, Repo, Commit, Blob) dùng chung 1 storage interface nhưng mỗi cái cần type-safe CRUD."
>
> **Context:** Storage layer hiện tại:
> ```go
> // CŨ — bỏ đi
> type Store interface {
>     Save(key string, value interface{}) error
>     Get(key string) (interface{}, error)  // caller phải type-assert
>     Delete(key string) error
> }
> // MỖI lần Get: val, _ := store.Get("user_123"); user := val.(*User)  // PANIC nếu sai!
> ```
> Yêu cầu: `Store[T]` generic, type-safe CRUD, benchmark chứng minh không chậm hơn `interface{}` version.

### Acceptance Criteria

- [ ] `Store[T]` generic với type-safe `Save(key string, value T)`, `Get(key string) (T, bool)`, `Delete(key string)`, `List() []T`
- [ ] CRUD hoạt động với ít nhất 3 entity types: `User`, `Repo`, `Commit`
- [ ] Không có type assertion nào trong caller code
- [ ] Error handling rõ ràng: key not found, save failure
- [ ] Benchmark: `Store[User]` vs `interface{}` store — memory và speed
- [ ] Thread-safe với concurrent access
- [ ] Refactor từ `interface{}` version sang generic có "migration path" rõ ràng

### Senior Thought-Process

**Senior nghĩ gì khi nhận requirement này:**

> "Refactor storage layer — đây là ticket kiểu 'thay đổi foundation'. Nếu làm sai, 15 service phía trên bị ảnh hưởng. Điều tôi nghĩ đến đầu tiên: *không phải viết mới từ đầu, mà là migration từ cũ sang mới*.
>
> Ở project minigit, tôi từng refactor `store.go` từ `interface{}` sang generic. Bài học:
> 1. **Không xóa interface cũ ngay** — wrap nó bằng generic adapter
> 2. **Migration từng bước**: 1 entity type trước, test pass, rồi mới sang entity tiếp theo
> 3. **Generic type không phải silver bullet**: nếu entity có behavior khác nhau nhiều, generic có thể ép buộc quá mức
>
> Phân rã bài toán:
> 1. **Core Store[T]**: interface + in-memory implementation (map[string]T)
> 2. **Entity types**: User, Repo, Commit — mỗi cái 1 file riêng
> 3. **Migration adapter**: `type UserStore = Store[*User]` — type alias giảm boilerplate
> 4. **Benchmark**: so sánh với cũ, đảm bảo không regression
> 5. **Race test**: `go test -race` với concurrent access"

### TODO Comments (Code Skeleton)

```go
package store

// ============================================================
// STEP 1: Define generic Store interface
// ============================================================

// TODO-[1]: Định nghĩa Store[T any] interface
// SENIOR ASKS: Tại sao T any thay vì T comparable?
// HINT: Key là string riêng — T là value type, không phải key. Value có cần comparable không?

type Store[T any] interface {
	// TODO: Save, Get, Delete, List methods
}

// ============================================================
// STEP 2: In-memory implementation
// ============================================================

// TODO-[2]: Định nghĩa MemoryStore[T any]
// SENIOR ASKS: Tại sao dùng map[string]T thay vì map[string]interface{}?
// HINT: Map value type là T trực tiếp → không boxing, không type assertion.

type MemoryStore[T any] struct {
	// TODO: mu sync.RWMutex, data map[string]T
}

// TODO-[3]: Constructor NewMemoryStore[T any]()
// SENIOR ASKS: Function generic syntax: func NewMemoryStore[T any]() *MemoryStore[T] — T cần khai báo ở đâu?
// HINT: T phải được declare TRƯỚC function name: func Name[T Constraint](params)

// TODO-[4]: Implement Save(key string, value T) error
// SENIOR ASKS: Override behavior hay append-only?
// HINT: Document behavior rõ: set (upsert) hay chỉ insert nếu chưa tồn tại?

// TODO-[5]: Implement Get(key string) (T, bool)
// SENIOR ASKS: Trả về T khi not found — zero value problem.
// HINT: var zero T; return zero, false. Đảm bảo T có thể là pointer hoặc value.

// TODO-[6]: Implement Delete(key string) error
// SENIOR ASKS: Delete key không tồn tại → error hay no-op?
// HINT: Go convention: Delete map key không tồn tại = no-op. Nên match hay khác?

// TODO-[7]: Implement List() []T
// SENIOR ASKS: Thứ tự items trong List có quan trọng không?
// HINT: Map iteration order = random. Nếu cần deterministic, cần gì thêm?
```

```go
// ============================================================
// STEP 3: Entity types
// ============================================================

// TODO-[8]: Định nghĩa User, Repo, Commit structs
// SENIOR ASKS: Nên dùng value type hay pointer type cho T?
// HINT: Store[*User] vs Store[User] — pointer cho phép nil (deleted?), value thì copy.

// ============================================================
// STEP 4: Type aliases để giảm verbosity
// ============================================================

// TODO-[9]: Tạo type alias cho từng entity store
// SENIOR ASKS: type UserStore = Store[*User] — cái này có cần generic parameter không?
// HINT: Type alias instantiate generic type — không cần re-declare T.

// ============================================================
// STEP 5: Migration from interface{} store
// ============================================================

// TODO-[10]: Viết hàm chuyển đổi từ interface{} store sang generic
// SENIOR ASKS: Nếu legacy store chứa mixed types (User và Repo cùng 1 map), migration thế nào?
// HINT: Không thể mixed trong Store[T] — mỗi T = 1 store instance. Đây là feature hay limitation?
```

```go
// ============================================================
// STEP 6: Benchmark
// ============================================================

// TODO-[11]: Benchmark MemoryStore[User] vs map[string]interface{} + type assertion
// SENIOR ASKS: Setup benchmark fair — cả 2 cùng workload, cùng goroutine count.
// HINT: Dùng b.Run() với sub-benchmark. Đo cả time và alloc (benchmem).
```

### Socratic Questions

**Câu hỏi để bạn tự suy nghĩ:**

1. **`Store[T any]` vs `Store[T comparable]`:** Value type T có cần comparable không? Nếu không dùng comparable, bạn mất khả năng gì? (Hint: không thể check duplicate bằng value comparison)

2. **Pointer vs Value:** `Store[*User]` và `Store[User]` — nếu `T = *User`, zero value của T là `nil`. Nếu `T = User`, zero value là struct rỗng. `Get` trả về `(nil, false)` vs `(User{}, false)` — cái nào caller-friendly hơn?

3. **Single Store vs Multiple Store:** Trước refactor: 1 store chứa mọi thứ. Sau refactor: `userStore`, `repoStore`, `commitStore` — 3 instances. Trade-off về memory? Khởi tạo? Quản lý lifecycle?

4. **Interface Segregation:** `Store[T]` có 4 methods (Save, Get, Delete, List). Nếu 1 caller chỉ cần `Get`, có nên tách thành `Getter[T]`, `Saver[T]` không? Generic interface có nên nhỏ hay lớn?

5. **Binary Size Impact:** `MemoryStore[*User]`, `MemoryStore[*Repo]`, `MemoryStore[*Commit]` — Go compiler tạo 3 bản monomorphized code. Trong microservice nhỏ, điều này có vấn đề không?

### Output Checklist

- [ ] TODO-[1] hoàn thành: `Store[T any]` interface với 4 methods
- [ ] TODO-[2] hoàn thành: `MemoryStore[T any]` struct với mutex + map
- [ ] TODO-[3] hoàn thành: Generic constructor hoạt động
- [ ] TODO-[4] hoàn thành: Save với documented behavior (upsert)
- [ ] TODO-[5] hoàn thành: Get trả zero value + bool
- [ ] TODO-[6] hoàn thành: Delete hoạt động đúng
- [ ] TODO-[7] hoàn thành: List trả đúng items
- [ ] TODO-[8] hoàn thành: 3 entity types (User, Repo, Commit) được define
- [ ] TODO-[9] hoàn thành: Type alias giảm verbosity
- [ ] TODO-[10] hoàn thành: Migration strategy documented
- [ ] TODO-[11] hoàn thành: Benchmark có so sánh với `interface{}`

### Test Checklist

- [ ] **Test case:** Save + Get round-trip cho User — type-safe, không assertion
- [ ] **Test case:** Get key không tồn tại → `(zero, false)` — miss path
- [ ] **Test case:** Save ghi đè (upsert) → Get trả value mới — overwrite behavior
- [ ] **Test case:** Delete → Get sau đó trả miss — deletion
- [ ] **Test case:** List trả đúng số lượng items — listing
- [ ] **Test case:** Concurrent Save + Get từ 10 goroutines — `go test -race` pass
- [ ] **Test case:** Save nil pointer `(*User)(nil)` — boundary, store xử lý thế nào?
- [ ] **Test case:** MemoryStore với 3 different T types (User, Repo, Commit) — cùng lúc, không interfere
- [ ] **Test case:** Benchmark: generic store không chậm hơn interface{} > 10%

### Retrospective

**Sau khi xong, hãy tự hỏi:**

1. **Nếu entity count tăng từ 3 lên 30:** Type alias cho mỗi entity store (`UserStore`, `RepoStore`...) có scalable không? Có cách nào reduce boilerplate mà vẫn type-safe?

2. **Persistence layer:** `MemoryStore` chỉ lưu in-memory. Nếu thêm `DiskStore[T]` hoặc `RedisStore[T]` — generic interface `Store[T]` có giúp swap implementation không? Có cần generic method signatures phức tạp hơn không?

3. **Query capability:** `List()` trả tất cả. Nếu cần `ListBy(predicate func(T) bool)` — generic function parameter syntax là gì? Điều này mở ra khả năng gì?

4. **Nếu rollback về `interface{}`:** Bạn có thể wrap `Store[T]` thành non-generic store không? Điều này có quan trọng trong migration strategy không?

---

## Weekly Integration & Review

### Day Structure (Tuần 8)

```
Day 1-2: Topic 04.1 — Type Parameters & Constraints
  - Đọc: go.dev/doc/tutorial/generics
  - Code: Max, Min, Clamp generic
  - Benchmark: generic vs non-generic

Day 3-4: Topic 04.2 — Generic Types (Struct)
  - Code: Cache[K, V] với TTL
  - Test: race detection, benchmark
  - So sánh: generic cache vs map[string]interface{}

Day 5: Topic 04.3 — Type Sets & ~ Operator
  - Code: Number constraint, Sum, Avg
  - Test: named types, compile-time rejection

Day 6-7: Mini-Project — Generic Object Store
  - Refactor: từ interface{} sang Store[T]
  - Benchmark: cuối cùng
  - Retrospective + Week 8 review
```

### Common Pitfalls (Bug Diary Template)

```markdown
## Generics Bug Diary — Tuần 8

### Pitfall 1: comparable không đủ cho ordering
- Lỗi: dùng `comparable` làm constraint cho `Max` → compile error vì thiếu `<`
- Fix: dùng `constraints.Ordered` hoặc tự định nghĩa

### Pitfall 2: Quên ~ operator với named types
- Lỗi: `Number` = `int64 | int32` → `type Money int64` không satisfy
- Fix: thêm `~`: `~int64 | ~int32`

### Pitfall 3: Generic zero value
- Lỗi: trả về `nil` cho generic T → compile error nếu T là value type
- Fix: `var zero T; return zero, false`

### Pitfall 4: Type inference failure
- Lỗi: `result := Max(3, 3.14)` → cannot infer T
- Fix: `Max[float64](3, 3.14)` hoặc ép kiểu trước

### Pitfall 5: Over-engineering với generics
- Lỗi: viết `Process[T, U, V](...)` cho code đơn giản
- Fix: YAGNI — nếu chỉ dùng 1 type, không cần generic

### Pitfall 6: Generic và binary size
- Lỗi: instantiate 50 generic types → binary phình to
- Fix: đo bằng `go build -ldflags="-s -w"` và `ls -la`, giới hạn instantiation
```

### Self-Quiz (Không mở notes)

1. `any` khác `comparable` chỗ nào? Cho ví dụ type satisfy `any` nhưng không satisfy `comparable`.
2. `~int64` include những type nào? `int64`? `type MyInt int64`? `int32`?
3. Tại sao `Max(3, 3.14)` không compile? Fix bằng cách nào?
4. Generic zero value: `var t T` vs `new(T)` — khác nhau? Cái nào dùng cho return value?
5. `Cache[K comparable, V any]` — tại sao K phải comparable còn V thì any?
6. Type inference: khi nào compiler tự infer `T`, khi nào phải gõ explicit?
7. Monomorphization là gì? Ảnh hưởng đến binary size thế nào?
8. Generics thay thế `interface{}` trong mọi trường hợp? Nếu không, ngoại lệ nào?

### Week 8 Done Criteria

- [ ] 3 topic hoàn thành với code chạy được
- [ ] Mini-project Store[T] refactor xong, benchmark có
- [ ] `go test -race ./...` pass toàn bộ
- [ ] Tự trả lờI được 6/8 self-quiz
- [ ] Retrospective viết cho từng topic
- [ ] Bug diary có ít nhất 2 entries từ kinh nghiệm thực

### Link to minigit

```markdown
## minigit Integration

Generics trong Phase 4 trực tiếp áp dụng cho:
- `ObjectStore[T Object]` — storage layer type-safe
- `Cache[SHA, *Blob]` — object caching với TTL
- `Sum[Number]` — aggregate operations cho packfile

Refactor path:
1. Phase 3: `store.go` dùng `interface{}`
2. Phase 4: refactor thành `ObjectStore[T Object]`
3. Phase 5-7: tất cả entities dùng Store[T] — không còn type assertion
```

---

## References

- **Primary:** https://go.dev/doc/tutorial/generics — Go official generics tutorial
- **Spec:** https://go.dev/ref/spec#Type_parameter_lists — Type parameters spec
- **Blog:** https://go.dev/blog/intro-generics — Introducing Generics (Go blog)
- **Blog:** https://go.dev/blog/when-generics — When to use generics
- **Package:** `golang.org/x/exp/constraints` — Standard constraints (Ordered, Signed, Unsigned, Integer, Float, Complex)
- **Book:** "The Go Programming Language" (Donovan & Kernighan) — Chapter nâng cao về type system

---

## Appendix A: Generics Decision Framework

> **Công cụ để quyết định "có nên dùng generics không" — dùng ở mọi code review.**

### Flowchart: Khi nào dùng Generic?

```
Bắt đầu: Bạn đang viết function/type cho 1 bài toán
│
├─ Chỉ dùng cho 1 type duy nhất? → KHÔNG dùng generic
│   Ví dụ: ParseUser(input string) *User → không cần [T]
│
├─ Dùng cho 2-3 type, mỗi type logic giống hệt nhau? → CÂN NHẮC generic
│   Ví dụ: Max(int), Max(float64) → Max[T Ordered] ✅
│
├─ Dùng cho 5+ type, logic giống nhau? → DÙNG generic
│   Ví dụ: Cache[string, User], Cache[string, Repo] → Cache[K, V] ✅
│
├─ Logic phụ thuộc vào behavior của type (method)? → Interface, không phải generic
│   Ví dụ: Save(reader io.Reader) → interface{} không, generic không, interface ✅
│
├─ Cần "container" type-safe (list, cache, store)? → DÙNG generic
│   Ví dụ: []T, map[K]V → List[T], Cache[K,V], Store[T] ✅
│
└─ Code chỉ dùng 1 type NHƯNG "để sau có thể mở rộng"? → KHÔNG dùng generic (YAGNI)
    Ví dụ: "sau này có thể cần Max cho custom type" → viết Max(int) trước
```

### Generic vs Interface vs Reflection — So sánh

| Tình huống | Generic | Interface | Reflection |
|---|---|---|---|
| Container type-safe (Cache, Store) | **✅ Tốt nhất** | ⚠️ Cần type assertion | ❌ Chậm, không type-safe |
| Algorithm trên nhiều numeric types | **✅ Tốt nhất** | ❌ Không phù hợp | ⚠️ Được nhưng chậm |
| Behavior-based (Save, Close) | ❌ Không phù hợp | **✅ Tốt nhất** | ⚠️ Overkill |
| Serialization (JSON, XML) | ⚠️ Không giúp nhiều | ✅ `json.Marshaler` | ❌ Tránh nếu có thể |
| Dependency Injection | ❌ Không | **✅ Tốt nhất** | ❌ Không |
| Unit Test Mock | ❌ Khó mock | **✅ Dễ mock** | ❌ Không |
| Binary Size nhạy cảm | ⚠️ Monomorphization phình to | ✅ Không ảnh hưởng | ✅ Không ảnh hưởng |
| Compile Time nhạy cảm | ⚠️ Chậm hơn | ✅ Nhanh | ✅ Nhanh |

### Những câu hỏi trước khi viết generic

1. **"Có thực sự cần không?"** — Nếu hiện tại chỉ dùng 1 type, đừng viết generic. Refactor sau khi cần.
2. **"Constraint có quá phức tạp không?"** — Nếu constraint dài hơn function signature, reconsider.
3. **"NgườI đọc có hiểu không?"** — `func Process[T any](input T) T` dễ hiểu. `func Transform[T any, U comparable, V ~int | ~float64](in T, mapper func(T) U) (V, error)` thì không.

---

## Appendix B: Benchmark Template (Copy & Điền)

> **Senior nói:** "Không benchmark = không có bằng chứng. Mọi quyết định dùng generics thay vì interface{} phải có số đo."

```go
package benchmark

import (
	"testing"
	"time"
)

// TODO-[B1]: Điền type cần benchmark
// SENIOR ASKS: Chọn entity đại diện — struct size nào? Pointer hay value?

type User struct {
	ID       string
	Username string
	Email    string
	Created  time.Time
	// TODO: thêm fields để struct có size thực tế (~100 bytes)
}

// TODO-[B2]: Benchmark Generic Store
// SENIOR ASKS: Pre-allocate map với make để benchmark fair?
// HINT: Không pre-allocate = benchmark allocation, không phải operation.

func BenchmarkMemoryStore_Get(b *testing.B) {
	store := NewMemoryStore[*User]()
	// Pre-populate
	for i := 0; i < 1000; i++ {
		// TODO: Save 1000 items
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// TODO: Get từ store
		}
	})
}

// TODO-[B3]: Benchmark interface{} Store + type assertion
// SENIOR ASKS: Đảm bảo assertion không bị optimize away.
// HINT: Gán kết quả assertion vào biến package-level để compiler không optimize.

func BenchmarkInterfaceStore_Get(b *testing.B) {
	// TODO: Tạo map[string]interface{} store
	// TODO: Pre-populate với *User values
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// TODO: Get + type assertion .(*User)
		}
	})
}

// TODO-[B4]: Run với benchmem
// Command: go test -bench=. -benchmem -count=5
// SENIOR ASKS: Đọc 4 cột: ns/op, B/op, allocs/op, và so sánh.
// HINT: Generic thường thắng ở allocs/op (0 vs 1) nhưng không chênh lệch ns/op nhiều.
```

### Đọc kết quả benchmark

```
// Output mẫu:
// BenchmarkMemoryStore_Get-8     50000000    22.3 ns/op    0 B/op    0 allocs/op
// BenchmarkInterfaceStore_Get-8  30000000    38.1 ns/op   16 B/op    1 allocs/op
//
// Giải thích:
// - 50M vs 30M ops/sec: Generic nhanh hơn ~60%
// - 0 vs 16 B/op: Generic không allocate, interface{} allocate 16 bytes mỗi lần assertion
// - 0 vs 1 allocs/op: Generic 0 GC pressure, interface{} có
```

---

## Appendix C: War Story — The `interface{}` Incident

> **Câu chuyện có thật từ production — để bạn hiểu tại sao generics quan trọng.**

### Bối cảnh

Hồi 2021, team tôi maintain 1 microservice caching product catalog. Cache layer dùng `map[string]interface{}`:

```go
// Code cũ — production incident source
type Cache struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}
```

### Incident #1: Wrong Type Assertion (P1 Outage)

```go
// File: product_service.go
func GetProductPrice(cache *Cache, productID string) (float64, error) {
    val, found := cache.Get("price_" + productID)
    if !found {
        return 0, errors.New("not found")
    }
    // ⚠️ Từng là *Product, nhưng junior refactor thành Price struct
    price := val.(*Product).Price  // PANIC: val is *Price, not *Product
    return price, nil
}
```

**Hậu quả:** Service crash 15 phút, 2000 requests failed, rollback.

### Incident #2: Nil Interface Trap

```go
// File: user_service.go  
func GetUser(cache *Cache, userID string) *User {
    val, found := cache.Get("user_" + userID)
    if !found || val == nil {
        return nil
    }
    // ⚠️ val != nil nhưng val.(*User) có thể panic
    // Vì interface value = (type=*User, value=nil) → val == nil = false!
    return val.(*User)
}
```

**Hậu quả:** Null pointer panic ở caller, log không có context để debug.

### Giải pháp cuối cùng (2022, sau Go 1.18)

```go
// Refactor sang generic — zero incidents từ đó đến nay
type Cache[V any] struct {
    mu   sync.RWMutex
    data map[string]V  // V = *Product, *User, *Order — type-safe at compile time
}

func (c *Cache[V]) Get(key string) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok  // Không cần assertion — V đã đúng type
}

// Usage:
var productCache = NewCache[*Product]()
priceCache := NewCache[*Price]()  // Separate cache, separate type

// Compile error nếu dùng sai type:
// productCache.Set("p1", &Price{})  // ERROR: cannot use *Price as *Product
```

### Bài học rút ra

| Vấn đề | `interface{}` | Generic |
|---|---|---|
| Type safety | Runtime (panic) | Compile time |
| Refactor safety | Dễ miss caller | Compiler bắt tất cả |
| Performance | Boxing + assertion | Direct, zero overhead |
| Readability | "Cái gì trong đây?" | "Cache của *Product" |
| On-call 3am | Debug interface{} | Không cần debug |

> **Senior kết luận:** "Generic không phải để code 'cool'. Nó để bạn ngủ ngon hơn, ít page hơn, và junior không thể vô tình crash production bằng type assertion sai."

---

## Appendix D: Generics Anti-Patterns (Đừng làm)

### Anti-Pattern 1: Generic cho 1 type

```go
// ❌ SAI: Chỉ dùng cho User, không cần generic
func ProcessUser[T *User](user T) T { return user }

// ✅ ĐÚNG: Đơn giản, không generic
func ProcessUser(user *User) *User { return user }
```

### Anti-Pattern 2: Constraint quá phức tạp

```go
// ❌ SAI: Đọc không nổi, maintain không nổi
func Transform[T any, U interface {
    ~int | ~int64 | ~uint | ~float64
    comparable
}, V ~string | ~[]byte](in T, fn func(T) U) (V, error)

// ✅ ĐÚNG: Tách thành named constraint, signature gọn
func Transform[T any](in T, fn func(T) Number) (string, error)
```

### Anti-Pattern 3: Thay thế interface bằng generic

```go
// ❌ SAI: Generic không thay thế behavior interface
func SaveTo[T any](writer T, data []byte) error
// T cần có Write method — nhưng constraint nào đảm bảo điều đó?

// ✅ ĐÚNG: Interface cho behavior
func SaveTo(writer io.Writer, data []byte) error
```

### Anti-Pattern 4: Generic "just in case"

```go
// ❌ SAI: YAGNI — chưa cần nhưng "để sau"
func Calculate[T Number](a, b T) T { return a + b }
// Chỉ dùng Calculate(1, 2) trong codebase

// ✅ ĐÚNG: Viết đơn giản, refactor khi cần
func Calculate(a, b int) int { return a + b }
```

### Anti-Pattern 5: Quên zero value của generic type

```go
// ❌ SAI: return nil cho generic T — compile error nếu T là int
func MaybeGet[T any](found bool, val T) T {
    if !found { return nil }  // ERROR: cannot use nil as T
    return val
}

// ✅ ĐÚNG: Dùng var zero value
func MaybeGet[T any](found bool, val T) T {
    if !found { var zero T; return zero }
    return val
}
```

---

> *"Generics giải quyết vấn đề type safety. Nhưng nó không giải quyết vấn đề thiết kế. Nếu interface của bạn đã lộn xộn, generics chỉ giúp bạn lộn xộn *type-safely* mà thôi."* — Senior's closing thought for Phase 4.
