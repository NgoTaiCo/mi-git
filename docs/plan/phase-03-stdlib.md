# Phase 3: Standard Library & Tooling (Tuan 6-7)

> "Stdlib la vu khi chinh cua Go. Hoc framework truoc ma khong biet net/http, encoding/json, database/sql ben trong hoat dong nhu the nao — giong nhu lai xe ma khong biet may no chay."
> — Senior Go Engineer

---

## Topic 03.1: HTTP Server (net/http)

### User Story

> Khach hang (Product Owner) noi: "Toi can mot REST API cho he thong e-commerce: CRUD products, orders. Chi duoc dung stdlib, khong dung Gin, Echo hay framework nao."
>
> Context: Khach hang muon hieu performance co ban cua Go truoc khi them dependency. Stdlib giup tiet kiem build size, de maintain, va quan trong hon — team se hieu chinh xuc cach Go xu ly HTTP request.

### Acceptance Criteria

- [ ] Server lang nghe duoc tren port 8080, tra loi duoc request
- [ ] Routing dung pattern: GET /products, GET /products/{id}, POST /products, PUT /products/{id}, DELETE /products/{id}
- [ ] Middleware logging in ra method, path, status code, duration cho moi request
- [ ] Co graceful shutdown khi nhan SIGTERM (khong drop request dang xu ly)
- [ ] Response tra ve JSON voi Content-Type: application/json
- [ ] Error response chuan: { "error": "message" } voi HTTP status tuong ung

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Nhieu nguoi cam Gin/Echo ngay khi viet HTTP server trong Go. Nhung toi bat junior hoc stdlib truoc vi 3 ly do:
>
> 1. **Gin/Echo chi la wrapper.** Duoi hood chung cung dung net/http.Handler. Hieu Handler interface = hieu 90% framework.
> 2. **Chi phi giau co.** O project truoc, toi thay team dung Echo chi vi 'quen tay' — app chi co 3 endpoint, stdlib du suc. Build size nhe hon 5MB, startup nhanh hon.
> 3. **Debug de hon.** Khi co bug o middleware chain, neu khong biet http.HandlerFunc hoat dong nhu the nao, ban se ngoi nhin stack trace ma khong biet loi tu dau.
>
> Van de cot loi o day la: **net/http chi co mux (router) rat co ban — khong ho tro method routing hay path parameter truc tiep.** Nen chung ta phai tu xay. Cach toi phan ra:
> - Buoc 1: Hieu Handler interface — chi 1 method: ServeHTTP(ResponseWriter, *Request)
> - Buoc 2: Viet custom mux xu ly routing dung http.ServeMux hoac tu viet
> - Buoc 3: Xay dung middleware chain bang pattern decorator
> - Buoc 4: Them graceful shutdown voi signal handling
>
> Cai nay giong voi van de toi gap o project monitoring API — ban dau team dung Gorilla Mux, sau do chuyen sang chi vi chi co 5 endpoint. Code giam 30%, build nhanh hon."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO-[1]: Dinh nghia Product struct va in-memory store
// SENIOR ASKS: Tai sao khong dung map[string]Product ma dung slice?
// HINT: Nho filter/later query, slice de iterate hon map. Nhung trade-off la gi?

type Product struct {
	// TODO: Dien cac truong: ID, Name, Price, CreatedAt
	// SENIOR ASKS: ID nen dung int hay string? Tai sao?
}

var (
	// TODO: Khai bao store va mutex bao ve
	// SENIOR ASKS: Tai sao can mutex o day? Khong dung thi sao?
	// HINT: Dong thoi doc/ghi vao slice se gay race condition
)

// TODO-[2]: Viet Handler interface implementation
// SENIOR ASKS: http.Handler chi co 1 method. No khac handler function nhu the nao?
// HINT: Function phai wrap thanh Handler thi moi dung voi ServeMux

type AppServer struct {
	// TODO: Luu store hoac dependency
}

func (s *AppServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Routing logic — phan biet method va path
	// SENIOR ASKS: http.ServeMux ho tro method matching khong?
	// HINT: ServeMux khong — phai tu check r.Method hoac dung pattern rieng
}

// TODO-[3]: Viet CRUD handler functions
// SENIOR ASKS: Tai sao signature la (w http.ResponseWriter, r *http.Request)?
// HINT: ResponseWriter viet response, Request chua moi thu tu client

func (s *AppServer) handleListProducts(w http.ResponseWriter, r *http.Request) {
	// TODO: Tra ve JSON danh sach product
	// SENIOR ASKS: json.NewEncoder(w).Encode() vs json.Marshal roi w.Write() — khac nhau gi?
	// HINT: Encoder stream truc tiep — khong can buffer toan bo vao memory
}

func (s *AppServer) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	// TODO: Lay id tu path, tim product, tra ve
	// SENIOR ASKS: Lam sao extract path parameter tu URL khi khong co framework?
	// HINT: strings.TrimPrefix hoac regex hoac path.Split — chon cai nao?
}

func (s *AppServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode JSON body, validate, luu vao store
	// SENIOR ASKS: json.Decoder co can gioi han body size khong?
	// HINT: r.Body khong gioi han mac dinh — http.MaxBytesReader la ai ban
}

// TODO-[4]: Viet middleware logging
// SENIOR ASKS: Middleware trong Go la gi? Phai implement interface dac biet khong?
// HINT: Chi la function nhan Handler, tra ve Handler — decorator pattern

type responseRecorder struct {
	http.ResponseWriter
	// TODO: Luu status code va response size
	// SENIOR ASKS: Tai sao can wrap ResponseWriter? No khong luu status code san a?
	// HINT: ResponseWriter.Write() tu set 200, nhung khong the doc lai duoc
}

func (rr *responseRecorder) WriteHeader(code int) {
	// TODO: Ghi nho status code roi goi original
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Log start time, goi next.ServeHTTP, log duration
		// SENIOR ASKS: Tai sao phai goi next.ServeHTTP chu khong goi handler truc tiep?
		// HINT: Middleware chain — moi lop boc ngoai 1 lop, goi trong de xu ly
	})
}

// TODO-[5]: Graceful shutdown
// SENIOR ASKS: Tai sao can graceful shutdown? Khong co thi sao?
// HINT: Container/K8s se kill process — request dang xu ly bi drop

func main() {
	// TODO: Tao server, setup routes, wrap middleware
	// TODO: Lang nghe signal SIGTERM/SIGINT
	// TODO: server.Shutdown(context.WithTimeout) de graceful shutdown
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **http.DefaultServeMux khac gi voi tu tao `http.NewServeMux()`?** Khi nao nen dung cai nao? (Hint: global vs local, testability)

2. **Tai sao Go chon interface voi 1 method (ServeHTTP) lam co che mo rong thay vi class inheritance nhu Java?** Dieu nay anh huong gi den cach viet middleware?

3. **Neu server co 1000 concurrent requests, moi request lam viec 2s — can bao nhieu goroutine?** Go xu ly dieu nay nhu the nao ma khong can thread pool nhu Java?

4. **Middleware chain: `logging(auth(handler)))` — thu tu boc co quan trong khong?** Neu dao nguoc logging va auth thi khac gi? Ve hinh anh stack call di.

5. **`w.Header().Set()` phai goi truoc hay sau `w.WriteHeader()`?** Sau khi viet response body thi set header duoc khong? Tai sao?

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Product struct dinh nghia day du, in-memory store co mutex bao ve
- [ ] TODO-[2] hoan thanh: Co custom handler hoac routing logic hoat dong
- [ ] TODO-[3] hoan thanh: 4 CRUD endpoint tra ve JSON dung chuan, HTTP status chinh xac (201 cho create, 404 cho not found)
- [ ] TODO-[4] hoan thanh: Logging in ra method, path, status code, duration. Co the log ca request body size.
- [ ] TODO-[5] hoan thanh: Server shutdown gracefully khi nhan Ctrl+C, khong drop request dang xu ly
- [ ] Co the chay `curl` test tat ca endpoint thanh cong

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: GET /products tra ve danh sach JSON — **vi sao:** verify basic response format
- [ ] Test case: GET /products/999 tra ve 404 — **boundary case:** ID khong ton tai, khong duoc panic
- [ ] Test case: POST /products voi body khong hop le tra ve 400 — **error path:** client gui JSON sai format
- [ ] Test case: Concurrent POST nhieu requests — **race condition:** verify mutex hoat dong dung
- [ ] Test case: Request co Accept header khac application/json — **content negotiation:** van tra ve JSON

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Stdlib vs framework nhu Gin — khi nao dung stdlib du, khi nao can framework? (Hint: routing phuc tap, validation, binding)
2. **Neu requirement thay doi:** "Them 100 endpoint" — ban se refactor nhu the nao? Co can framework luc do khong?
3. **Architecture decision:** Tai sao toi bat dau voi in-memory store thay voi database ngay? Dieu nay co y nghia gi ve iterative development?

---

## Topic 03.2: JSON (encoding/json)

### User Story

> Khach hang (Product Owner) noi: "API tra ve JSON, nhan JSON tu client. Co field nullable (description co the null), co field an khi empty (tags trong khi empty thi khong xuat hien trong JSON)."
>
> Context: Client la Flutter app — can JSON format chuan, predictable. Khong duoc co field "tags": [] hay "tags": null khi khong co du lieu.

### Acceptance Criteria

- [ ] Struct co struct tags dung dinh dang: `json:"field_name,omitempty"`
- [ ] Unmarshal JSON tu client thanh struct dung — xu ly duoc field thieu, field null
- [ ] Marshal struct thanh JSON — field empty duoc an (omitempty), field null duoc xu ly dung
- [ ] Co custom MarshalJSON/UnmarshalJSON cho it nhat 1 type (vi du: time format custom)
- [ ] Validation sau khi unmarshal: required field, numeric range

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Cai nay nghe de nhung la noi junior hay dinh nhat. Toi nho co lan junior dung `omitempty` voi bool — ket qua la field `is_active` bien mat khi false. Client Flutter bao la server tra thieu field.
>
> Van de cot loi o day la: **Go khong co nullable type natively.** `string` khong the null, `int` khong the null. Zero value ("", 0, false) khac voi null. Nen khi client gui `"price": null`, chung ta phai xu ly nhu the nao?
>
> Cach toi phan ra:
> - Buoc 1: Struct tags co ban — omitempty, chinh ten field
> - Buoc 2: NULL handling — sql.NullString, *string (pointer), hay custom type?
> - Buoc 3: Custom marshal cho type dac biet (timestamp, enum)
> - Buoc 4: Validation sau unmarshal — Go khong co annotation nhu Java, phai tu viet
>
> O project Flutter bridge truoc, toi tung gap van de: Go tra ve `null` cho pointer nhung Flutter parse thanh 0 — khong phan biet duoc 'chua co gia tri' va 'gia tri 0'. Viec nay quan trong cho UX."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// TODO-[1]: Dinh nghia struct voi JSON tags
// SENIOR ASKS: json tag luon co 2 phan: ten field va option. "omitempty" co tac dung gi?
// HINT: Go struct tag la string metadata — reflect package doc no luc compile time

type Product struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	Description *string    `json:"description"` // pointer = nullable
	Tags        []string   `json:"tags,omitempty"` // omitempty = an khi empty
	CreatedAt   CustomTime `json:"created_at"`
	IsActive    bool       `json:"is_active"` // TODO: Co nen omitempty khong?
	// SENIOR ASKS: Neu dung omitempty voi bool, dieu gi xay ra khi IsActive = false?
	// HINT: Go xem false la zero value → omitempty se an field nay → nguy hiem!
}

// TODO-[2]: Custom time type de format JSON
// SENIOR ASKS: Tai sao khong dung time.Time truc tiep ma phai custom?
// HINT: time.Time marshal ra RFC3339 — client co the can format khac

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	// TODO: Format thanh "2006-01-02 15:04:05" thay vi RFC3339
	// SENIOR ASKS: Tai sao Go dung layout string "2006-01-02" thay vi "%Y-%m-%d"?
	// HINT: Day la reference time: Mon Jan 2 15:04:05 MST 2006 = 1-2-3-4-5-6-7
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// TODO: Parse string thanh time.Time
	// SENIOR ASKS: Tai sao receiver la *CustomTime ma khong phai CustomTime?
	// HINT: Can modify receiver — value receiver se modify ban copy
}

// TODO-[3]: NULL handling — so sanh 3 cach
// SENIOR ASKS: Co may cach de xu ly NULL trong JSON voi Go? Uu/nhuoc diem tung cach?
// HINT: (1) Pointer *string (2) sql.NullString (3) json.RawMessage

func demoNullHandling() {
	// JSON: {"name": "Book", "description": null}
	jsonData := `{"name": "Book", "description": null}`
	
	// Cach 1: Pointer
	var p1 struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}
	json.Unmarshal([]byte(jsonData), &p1)
	// TODO: p1.Description bang gi khi JSON la null? Khi JSON khong co field?
	// SENIOR ASKS: Su khac nhau giua null va undefined trong JSON — Go xu ly khac nhau khong?
	
	// Cach 2: sql.NullString
	var p2 struct {
		Name        string         `json:"name"`
		Description sql.NullString `json:"description"`
	}
	// TODO: Unmarshal va kiem tra Valid flag
}

// TODO-[4]: Viet ham decode co validation
// SENIOR ASKS: json.Decoder co validate san khong? Can tu validate them khong?
// HINT: Decoder chi check syntax — business validation (price > 0) phai tu viet

func decodeProduct(data []byte) (Product, error) {
	var p Product
	// TODO: Unmarshal roi validate
	// SENIOR ASKS: Nen unmarshal truoc roi validate, hay co cach nao validate trong luc unmarshal?
	// HINT: Custom Unmarshal co the validate nhung se phuc tap — tach roi la tot hon
	
	// Validate: Name khong rong, Price > 0
	// TODO: Implement validation
	
	return p, nil
}

// TODO-[5]: json.RawMessage cho flexible parsing
// SENIOR ASKS: Khi nao dung json.RawMessage? Cho vi du thuc te.
// HINT: Khi payload co nhieu "shape" khac nhau — vi du: webhook events

type WebhookEvent struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

func processWebhook(data []byte) {
	// TODO: Unmarshal WebhookEvent, roi switch EventType de unmarshal Payload tiep
	// SENIOR ASKS: Tai sao khong unmarshal truc tiep ma phai 2 buoc?
	// HINT: Go can biet struct type truoc khi unmarshal — khong the dynamic
}

func main() {
	// TODO: Demo marshal/unmarshal voi cac truong hop
	// - Product co day du field
	// - Product khong co tags (tags = nil hoac [])
	// - Product co description = null
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`omitempty` voi slice: `[]string{}` vs `nil` — cai nao se xuat hien trong JSON?** Tai sao khac nhau? Dieu nay anh huong gi den Flutter client?

2. **`json:"-"` khac gi voi khong co tag?** Khi nao ban muon field co trong struct nhung khong xuat hien trong JSON?

3. **Pointer receiver vs value receiver cho `UnmarshalJSON`: Neu ban viet `(ct CustomTime)` thay vi `(ct *CustomTime)`, dieu gi xay ra?** Thu nghi ve cach json package goi method nay.

4. **`json.Decoder` co `DisallowUnknownFields()` — tai sao method nay quan trong cho API versioning?** Khi client gui field moi ma server chua biet, ban muon silent ignore hay error?

5. **`json.Marshal` co the panic khong?** Khi nao? Lam sao de tranh?

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Struct tags dung chuan, co vi du omitempty voi slice/string/bool
- [ ] TODO-[2] hoan thanh: CustomTime marshal/unmarshal dung format
- [ ] TODO-[3] hoan thanh: Demo duoc su khac nhau giua pointer null va absent field
- [ ] TODO-[4] hoan thanh: Ham decodeProduct validate name, price > 0
- [ ] TODO-[5] hoan thanh: WebhookEvent parse duoc 2 buoc voi RawMessage
- [ ] In ra console cho thay JSON output ro rang cho tung truong hop

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Marshal Product voi tags nil → khong co field tags trong JSON — **vi sao:** verify omitempty voi nil slice
- [ ] Test case: Unmarshal JSON co description: null → pointer nil — **boundary case:** nullable handling
- [ ] Test case: Unmarshal JSON thieu field → zero value, khong error — **default behavior:** Go khong bat loi thieu field
- [ ] Test case: Unmarshal JSON co unknown field → silent ignore (mac dinh) — **verify behavior:** biet duoc decoder khong strict
- [ ] Test case: Marshal voi DisallowUnknownFields va input co unknown field → error — **strict mode:** API contract enforcement

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Pointer (*string) vs value (string) cho field nullable — chi phi memory, GC pressure, complexity. Khi nao nen dung cai nao?
2. **Neu requirement thay doi:** "API v2 doi ten field 'name' thanh 'title' nhung van phai support v1" — ban se xu ly backward compatibility nhu the nao?
3. **Architecture decision:** Tai sao Go chon struct tags (compile-time string) thay vi annotation nhu Java (runtime reflection)? Loi va hai cua cach tiep can nay?

---

## Topic 03.3: Database (database/sql)

### User Story

> Khach hang (Product Owner) noi: "Luu products vao SQLite. Nhieu request cung luc, khong duoc leak connection. Query phai dung prepared statement."
>
> Context: Hieu ung mang — nhieu user cung truy cap. Connection pool khong quan ly tot = server chet. Day la loi production pho bien nhat toi thay o junior.

### Acceptance Criteria

- [ ] Mo duoc ket noi SQLite voi database/sql (dung mattn/go-sqlite3 driver)
- [ ] Connection pool duoc cau hinh: max open, max idle, max lifetime
- [ ] CRUD products qua prepared statement (Exec, Query, QueryRow)
- [ ] rows.Close() duoc defer ngay sau khi tao
- [ ] Xu ly NULL tu database dung (sql.NullString, sql.NullInt64)
- [ ] Transaction cho operation nhieu buoc (vi du: tao order + cap nhat inventory)

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Loi pho bien nhat toi thay: quen rows.Close(). Hau qua: connection pool can kiet, server khong the query database nua, cascade failure. Hoi toi o project logistics, co 1 service leak connection — sau 2 tieng chay production thi chet. Nguyen nhan? 1 ham helper query ma khong defer rows.Close().
>
> Van de cot loi o day la: **database/sql quan ly connection pool tu dong, nhung khong tu dong giai phong statement hay result set.** Ban phai chu dong quan ly lifecycle.
>
> Cach toi phan ra:
> - Buoc 1: Mo DB connection, cau hinh pool
> - Buoc 2: Prepared statement va lifecycle
> - Buoc 3: Query/Scan rows + rows.Close() pattern
> - Buoc 4: Transaction voi rollback tren error
> - Buoc 5: NULL handling
>
> Quy tac vang toi day junior: **'defer rows.Close() ngay sau khi check error cua Query'** — khong doi, khong quen.
>
> Cai nay giong voi van de toi gap o project e-commerce — tao order phai gom: insert orders, insert order_items, cap nhat inventory. 3 operation, 1 transaction. Rollback neu 1 operation fail."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
)

// TODO-[1]: Cau hinh connection pool
// SENIOR ASKS: Tai sao can SetMaxOpenConns? Khong set thi mac dinh la bao nhieu?
// HINT: Mac dinh = 0 = unlimited — tren SQLite 1 process chi write duoc 1 connection

type DB struct {
	*sql.DB
}

func OpenDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	
	// TODO: Cau hinin pool
	// SENIOR ASKS: SetConnMaxLifetime de lam gi? Tai sao khong de vinh vien?
	// HINT: Connection cu co the bi server dong, stale state, memory leak
	
	// TODO: Ping de verify connection thuc su hoat dong
	// SENIOR ASKS: sql.Open co thuc su mo connection khong?
	// HINT: Khong! Open chi tao object. Ping moi thuc su ket noi.
	
	return &DB{db}, nil
}

// TODO-[2]: Khoi tao schema
// SENIOR ASKS: Nen dung migration tool hay schema init trong code?
// HINT: Code init don gian cho dev/test. Production dung migration tool.

func (db *DB) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(schema)
	return err
}

// TODO-[3]: CRUD voi prepared statement
// SENIOR ASKS: db.Query() va db.QueryRow() khac nhau gi? Khi nao dung cai nao?
// HINT: Query = nhieu rows. QueryRow = 1 row. Exec = khong tra rows.

type Product struct {
	ID          int64
	Name        string
	Price       float64
	Description sql.NullString // TODO: Tai sao dung sql.NullString thay vi *string?
	// SENIOR ASKS: sql.NullString vs *string — khac nhau gi trong thuc te?
	// HINT: NullString ro rang hon (co Valid flag), *string don gian hon nhung de nham
	CreatedAt   time.Time
}

func (db *DB) CreateProduct(p *Product) error {
	// TODO: INSERT voi Exec, lay LastInsertId
	// SENIOR ASKS: Tai sao dung pointer *Product? Co nen tra ve Product (value) khong?
	// HINT: Pointer de modify ID sau insert, va tiet kiem copy neu struct lon
	
	result, err := db.Exec(
		"INSERT INTO products (name, price, description) VALUES (?, ?, ?)",
		p.Name, p.Price, p.Description,
	)
	if err != nil {
		return err
	}
	p.ID, _ = result.LastInsertId()
	return nil
}

func (db *DB) GetProduct(id int64) (Product, error) {
	var p Product
	// TODO: QueryRow + Scan
	// SENIOR ASKS: row.Scan() phai truyen pointer. Neu quen & thi sao?
	// HINT: Khong compile — Scan nhan ...interface{} nhung can addressable value
	
	row := db.QueryRow("SELECT id, name, price, description, created_at FROM products WHERE id = ?", id)
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt)
	// TODO: Xu ly sql.ErrNoRows — tra ve loi gi cho caller?
	// SENIOR ASKS: Nen tra ve sql.ErrNoRows truc tiep hay wrap?
	// HINT: Go convention: ErrNoRows la sentinel error — co the errors.Is() check
	return p, err
}

func (db *DB) ListProducts() ([]Product, error) {
	// TODO: Query + iterate rows
	// SENIOR ASKS: rows.Close() phai defer ngay sau khi check error. Vi sao phai ngay?
	// HINT: Neu co return/error truoc khi defer, rows khong bao gio dong
	
	rows, err := db.Query("SELECT id, name, price, description, created_at FROM products")
	if err != nil {
		return nil, err
	}
	// DONG NAY LA QUY TAC VANG — DE NGAY SAU CHECK ERROR
	defer rows.Close() // TODO: Neu bo dong nay, dieu gi xay ra?
	
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(/* ... */); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	// TODO: rows.Err() sau vong lap — tai sao can?
	// SENIOR ASKS: rows.Next() co the dung vi loi hay chi vi het data?
	// HINT: Next() tra ve false vi loi hoac het — rows.Err() de phan biet
	
	return products, rows.Err()
}

// TODO-[4]: Transaction
// SENIOR ASKS: Tai sao can transaction? Khong co thi co van de gi?
// HINT: Tinh nguyen tu - 1 trong 2 operation fail = data inconsistent

func (db *DB) CreateOrderWithItems(order Order, items []OrderItem) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// TODO: defer rollback. Tai sao rollback ma khong phai commit?
	// SENIOR ASKS: Rollback sau khi da commit co loi khong?
	// HINT: Rollback sau commit la no-op — an toan. Nguoc lai thi khong.
	defer tx.Rollback()
	
	// TODO: INSERT order, INSERT items, UPDATE inventory
	// Moi operation dung tx.Exec() thay vi db.Exec()
	
	return tx.Commit()
}

// TODO-[5]: Context-aware queries
// SENIOR ASKS: QueryContext khac gi Query? Tai sao nen dung?
// HINT: Context co the cancel query khi request timeout — tranh query chay mai

func (db *DB) SearchProducts(ctx context.Context, keyword string) ([]Product, error) {
	// TODO: Dung QueryContext thay vi Query
	// SENIOR ASKS: Context nay tu dau ra? Client request → server handler → database
	// HINT: http.Request co r.Context() — truyen xuong
	rows, err := db.QueryContext(ctx, "SELECT * FROM products WHERE name LIKE ?", "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// ...
}

func main() {
	// TODO: Open DB, init schema, demo CRUD
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`sql.NullString` va `*string` — khi unmarshal tu database, neu gia tri la NULL, chung khac nhau ra sao?** Cai nao de xu ly hon trong template/API response?

2. **`rows.Close()` la no-op neu rows da het (Next() tra ve false) — vay tai sao van phai defer?** (Hint: Next() dung vi error, khong phai het data)

3. **`db.Query()` vs `db.Prepare()` + `stmt.Query()` — khi nao prepared statement co loi?** (Hint: 1 query 1 lan → overhead. 100 lan → loi vi cache plan)

4. **Transaction isolation level: SQLite mac dinh la gi?** Dieu gi xay ra neu 2 goroutine cung `BEGIN` va cung `UPDATE` 1 row?

5. **`sql.NullString` khong unmarshal duoc JSON truc tiep — ban se viet custom Unmarshal nhu the nao?**

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: DB mo duoc, pool duoc cau hinh, Ping thanh cong
- [ ] TODO-[2] hoan thanh: Schema tao bang products thanh cong
- [ ] TODO-[3] hoan thanh: 4 CRUD operation hoat dong, rows.Close() defer ngay sau check error
- [ ] TODO-[4] hoan thanh: Transaction rollback tren error, commit khi thanh cong
- [ ] TODO-[5] hoan thanh: SearchProducts dung context, cancel duoc
- [ ] Chay `go test -race` qua — khong co race condition

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Create product roi Get product — verify round-trip — **vi sao:** dam bao serialize/deserialize dung
- [ ] Test case: Get product voi ID khong ton tai — verify ErrNoRows — **boundary case:** khong panic
- [ ] Test case: Concurrent CreateProduct 100 goroutine — verify khong race — **stress test:** SQLite + mutex
- [ ] Test case: Transaction rollback khi buoc 2 fail — verify data khong bi partial commit — **integrity:** core transaction benefit
- [ ] Test case: Query voi context da cancel — verify query dung som — **context propagation:** resource cleanup

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** In-memory store (sync.RWMutex) vs SQLite — khi nao in-memory du, khi nao can persistence?
2. **Neu requirement thay doi:** "Chuyen sang PostgreSQL" — code database/sql co gi thay doi? Layer nao can refactor?
3. **Architecture decision:** Tai sao database/sql dung interface-based design (driver model) thay vi implementation cu the? Dieu nay giup gi khi switch database?

---

## Topic 03.4: File I/O (os, io, bufio)

### User Story

> Khach hang (Product Owner) noi: "Upload file CSV 1GB, doc tung dong xu ly. Khong load het vao memory. Co the ghi tam file neu can."
>
> Context: Server chi co 512MB RAM. File 1GB khong the doc toan bo. Phai dung streaming I/O — doc tung chunk, xu ly, giai phong.

### Acceptance Criteria

- [ ] Doc file lon (>100MB) ma khong vuot qua 50MB memory usage
- [ ] Dung bufio.Scanner hoac bufio.Reader doc tung dong
- [ ] Ghi ket qua ra file tam dung os.CreateTemp — tu dong xoa sau khi dung
- [ ] Copy file tu reader sang writer dung io.Copy
- [ ] Gioi han memory voi io.LimitReader khi can

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Toi tung thay junior dung ioutil.ReadAll de doc file upload. File 2GB, server 1GB RAM — ket qua? OOM kill. Go runtime khong bao loi de nhin, process bien mat. Log chi ghi 'killed'.
>
> Van de cot loi o day la: **Go co nhieu lop I/O abstraction (io.Reader, bufio.Reader, os.File) — dung sai layer = performance kem hoac memory leak.**
>
> Cach toi phan ra:
> - Buoc 1: Phan biet io.Reader (interface) vs os.File (concrete) vs bufio.Reader (buffered)
> - Buoc 2: Streaming read — bufio.Scanner cho text, bufio.Reader cho binary
> - Buoc 3: Temp files — tao, dung, xoa (defer os.Remove)
> - Buoc 4: io.Copy de chuyen data giua reader/writer khong can buffer trung gian
> - Buoc 5: Memory gioi han — io.LimitReader, http.MaxBytesReader
>
> O project log processor truoc, toi phai xu ly file log 5GB moi ngay. Scanner voi bufio 64KB doc tung dong, xu ly, ghi ket qua ra file moi. Memory on dinh ~100MB."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TODO-[1]: Doc file lon tung dong
// SENIOR ASKS: os.Open tra ve gi? *os.File implement interface nao?
// HINT: *os.File implement io.Reader, io.Writer, io.Closer — nhieu interface

func processFileLineByLine(filename string, processor func(string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close() // TODO: Tai sao defer f.Close() ma khong phai defer f.Close() ngay?
	// SENIOR ASKS: Neu bo defer f.Close(), dieu gi xay ra?
	// HINT: File descriptor leak — OS gioi han FD moi process (thuong 1024-65536)
	
	// Cach 1: bufio.Scanner — don gian, nhung gioi han 64KB/dong
	scanner := bufio.NewScanner(f)
	// TODO: Tang buffer neu dong dai hon 64KB
	// SENIOR ASKS: Scanner mac dinh buffer bao nhieu? Lam sao tang?
	// HINT: scanner.Buffer(make([]byte, 1024), 1024*1024) — max 1MB/dong
	
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if err := processor(line); err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}
	}
	// TODO: scanner.Err() sau vong lap — tai sao can?
	// SENIOR ASKS: scanner.Scan() co the tra ve false vi loi hay het file?
	// HINT: Giong rows.Next() — Err() de phan biet
	return scanner.Err()
}

// TODO-[2]: bufio.Reader cho kiem soat cao hon
// SENIOR ASKS: bufio.Reader khac bufio.Scanner nhu the nao? Khi nao dung cai nao?
// HINT: Reader cho bat ky delimiter (khong chi newline), co Peek, ReadString, ReadLine

func processFileWithReader(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	reader := bufio.NewReader(f)
	// TODO: Doc tung dong bang ReadString('\n')
	// SENIOR ASKS: ReadString('\n') khac gi ReadLine? Luc nao dung cai nao?
	// HINT: ReadLine tra ve slice (co the partial), ReadString tra ve string day du
	
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			// Xu ly line
			fmt.Print(line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO-[3]: Ghi file tam — temp file pattern
// SENIOR ASKS: Tai sao dung os.CreateTemp thay vi tu chon ten file?
// HINT: Ten duy nhat, tu dong xoa duoc, khong conflict

func processToTemp(inputPath string) (string, error) {
	// Tao temp file trong thu muc tam
	tmpFile, err := os.CreateTemp("", "processed-*.csv")
	if err != nil {
		return "", err
	}
	// TODO: defer xoa temp file. Tai sao defer ma khong xoa ngay?
	// SENIOR ASKS: Khi nao temp file duoc xoa? Neu process bi kill thi sao?
	// HINT: Defer chay khi function return. Kill -9 thi khong chay defer!
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	// TODO: Mo input, copy sang tmpFile qua io.Copy
	// SENIOR ASKS: io.Copy co buffer san khong? Buffer bao nhieu?
	// HINT: io.Copy dung 32KB buffer mac dinh — co the customize
	
	return tmpFile.Name(), nil
}

// TODO-[4]: Gioi han memory — LimitReader
// SENIOR ASKS: io.LimitReader de lam gi? Khi nao can?
// HINT: Gioi han so byte doc tu reader — security (khong bi DoS bang upload khong lo)

func safeRead(reader io.Reader, maxBytes int64) ([]byte, error) {
	// TODO: Dung LimitReader de gioi han
	// SENIOR ASKS: LimitReader vuot qua thi tra ve gi? EOF hay error?
	// HINT: Tra ve EOF khi dat limit — khong phai error! Can check size sau.
	limited := io.LimitReader(reader, maxBytes)
	return io.ReadAll(limited)
}

// TODO-[5]: Multi-reader: doc file nen gzip + xu ly
// SENIOR ASKS: gzip.NewReader nhan gi? Tai sao no nhan io.Reader?
// HINT: io.Reader la abstraction — gzip, bufio, file, network deu la reader

func processGzipFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	// TODO: Tao gzip reader tu file reader
	// SENIOR ASKS: gzip.NewReader co phai doc toan bo file khong?
	// HINT: Khong! No stream tung chunk — memory nho
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzReader.Close()
	
	scanner := bufio.NewScanner(gzReader)
	// ... xu ly
	return nil
}

// TODO-[6]: io.TeeReader — doc dong thoi ghi ra 2 noi
// SENIOR ASKS: TeeReader dung de lam gi? Cho vi du thuc te.
// HINT: Doc file + tinh hash dong thoi — 1 lan doc, 2 tac vu

func calculateHashAndProcess(reader io.Reader) ([]byte, error) {
	// TODO: Su dung io.TeeReader de copy sang hasher trong luc doc
	return nil, nil
}

func main() {
	// TODO: Demo processFileLineByLine voi file CSV
	// TODO: Demo memory usage — file lon ma RAM khong tang
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`bufio.Scanner` co buffer 64KB. Neu 1 dong dai 1MB, Scanner lam gi?** Giai phap la gi? Tai sao khong tang buffer len 1GB luon?

2. **`io.Copy(dst, src)` — src la network connection, dst la file. Khi src bi dong dot ngot, Copy tra ve gi?** Ban xu ly nhu the nao?

3. **os.CreateTemp tao file o dau?** Co the chi dinh thu muc khac khong? Tai sao nen dung thu muc tam thay vi thu muc ung dung?

4. **`defer f.Close()` trong vong lap `for _, f := range files` — co van de gi?** (Hint: defer chay khi function return, khong phai khi iteration ket thuc)

5. **`io.ReadAll` vs `bufio.Scanner` — cai nao memory-efficient hon?** Khi nao nen dung `io.ReadAll`? (Hint: khi biet chac data nho)

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: processFileLineByLine doc file lon tung dong, memory < 50MB
- [ ] TODO-[2] hoan thanh: processFileWithReader dung bufio.Reader
- [ ] TODO-[3] hoan thanh: Temp file tao, ghi, defer xoa
- [ ] TODO-[4] hoan thanh: safeRead gioi han memory voi LimitReader
- [ ] TODO-[5] hoan thanh: Doc file gzip ma khong giai nen truoc
- [ ] TODO-[6] hoan thanh: Hieu io.TeeReader pattern

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: Process file 10MB, verify memory khong vuot 20MB — **stress test:** dung runtime.ReadMemStats
- [ ] Test case: File trong (0 byte) — khong hang — **boundary case:** Scanner voi file trong
- [ ] Test case: File 1 dong dai 100KB (vuot mac dinh 64KB) — **edge case:** scanner buffer overflow
- [ ] Test case: Temp file tu xoa sau khi function return — **resource cleanup:** verify khong con file tam
- [ ] Test case: LimitReader voi data 2MB, limit 1MB — chi doc duoc 1MB — **limit enforcement:** verify io.EOF

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** bufio.Scanner (don gian) vs bufio.Reader (linh hoat) — khi nao dung cai nao? Co rule of thumb khong?
2. **Neu requirement thay doi:** "File CSV 10GB, can xu ly song song 4 goroutine" — ban se chia file nhu the nao? (Hint: seek offsets, record boundaries)
3. **Architecture decision:** Tai sao Go chon `io.Reader` (1 method Pull-based) thay vi Push-based callback? Loi ich cua interface nho (1 method) so voi interface lon?

---

## Topic 03.5: Testing & Benchmarking

### User Story

> Khach hang (Product Owner) noi: "Viet test cho toan bo API. Khach hang can report performance tren 1000 concurrent users. Chi so: latency p99 < 100ms, throughput > 5000 req/s."
>
> Context: Khach hang la enterprise — can bao dam performance truoc khi sign contract. Test khong chi la 'pass' ma phai co so lieu cu the.

### Acceptance Criteria

- [ ] Table-driven tests cho CRUD API — 1 ham test, nhieu case
- [ ] Subtests dung t.Run() — tung case chay rieng, fail rieng
- [ ] httptest.Server de test HTTP handler ma khong can bind port that
- [ ] Benchmark cho ham tinh toan — do thoi gian va allocations
- [ ] Memory profile voi pprof — tim allocations khong can thiet
- [ ] CPU profile — xem ham nao ton thoi gian nhat

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Test trong Go khac Java/Python. Table-driven la pattern bat buoc — khong phai option. Toi tung review code junior viet 10 ham test, moi ham 1 case, copy-paste 90%. Refactor thanh 1 ham table-driven: 40 dong thay vi 400 dong.
>
> Van de cot loi o day la: **Testing trong Go la first-class citizen — khong can framework ngoai. nhung phai biet cach viet test 'dung Go way'.**
>
> Cach toi phan ra:
> - Buoc 1: Table-driven test pattern — slice cau truc, vong lap, t.Run()
> - Buoc 2: httptest — test HTTP ma khong can chay server that
> - Buoc 3: Benchmark + memory allocation tracking
> - Buoc 4: pprof — CPU va memory profile
> - Buoc 5: Race detection voi -race flag
>
> Quy tac toi day junior: **'Moi ham export phai co test. Moi test phai co error case. Moi benchmark phai co -benchmem.'**
>
> Cai nay giong voi van de toi gap o project payment gateway — can chung minh API xu ly duoc 10k req/s. Benchmark + pprof giup tim duoc bottleneck o JSON marshal."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO-[1]: Table-driven test cho tinh toan
// SENIOR ASKS: Tai sao table-driven? Khong viet 10 ham Test rieng duoc khong?
// HINT: Table-driven: 1 ham, nhieu case, de them case, khong copy-paste

func CalculateDiscount(price float64, tier string) (float64, error) {
	// TODO: Implement — gia su bronze=0%, silver=5%, gold=10%
	// error neu price < 0 hoac tier khong hop le
	return 0, nil
}

func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name    string      // TODO: Tai sao co field name?
		price   float64     // input
		tier    string      // input
		want    float64     // expected output
		wantErr bool        // TODO: Tai sao khong so sanh error message?
		// SENIOR ASKS: Nen so sanh error type, error message, hay chi check != nil?
		// HINT: Type cho sentinel, message cho user-facing, nil cho generic
	}{
		{"bronze no discount", 100.0, "bronze", 100.0, false},
		{"silver 5%", 100.0, "silver", 95.0, false},
		{"gold 10%", 100.0, "gold", 90.0, false},
		{"negative price", -10.0, "gold", 0, true},
		{"invalid tier", 100.0, "platinum", 0, true},
		{"free price", 0.0, "gold", 0.0, false}, // TODO: Day la boundary case
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // TODO: t.Run() khac gi khong dung?
			// SENIOR ASKS: t.Run() cho phep chay subtest song song khong?
			// HINT: t.Parallel() — nhung phai dam bao test khong race nhau
			got, err := CalculateDiscount(tt.price, tt.tier)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				// TODO: So sanh float64 bang == — co an toan khong?
				// SENIOR ASKS: Tai sao so sanh float truc tiep co the fail? Giai phap?
				t.Errorf("CalculateDiscount() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO-[2]: httptest — test HTTP handler khong can chay server that
// SENIOR ASKS: httptest.Server khac gi server that? Port la gi?
// HINT: Server chay tren localhost:0 (random port) hoac ResponseRecorder (khong can port)

func TestHandler(t *testing.T) {
	// Cach 1: ResponseRecorder — khong can mang that
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()
	
	// TODO: Goi handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Product{
			{ID: 1, Name: "Book"},
		})
	})
	handler.ServeHTTP(rr, req)
	
	// TODO: Verify status code, content-type, body
	// SENIOR ASKS: rr.Result() khac gi rr.Code, rr.Body?
	// HINT: Result() tra ve *http.Response (full), Code la int, Body la *bytes.Buffer
	
	if rr.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("want Content-Type application/json, got %s", ct)
	}
	
	// TODO: Verify JSON body — cach nao tot nhat?
	// SENIOR ASKS: json.Unmarshal roi so sanh struct, hay so sanh string?
	// HINT: Unmarshal + struct so sanh linh hoat hon (bo qua field khong quan trong)
}

// TODO-[3]: Benchmark
// SENIOR ASKS: Benchmark khac Test nhu the nao? Go chay khac nhau ra sao?
// HINT: go test -bench=. — chay nhieu lan, tinh thoi gian trung binh

func BenchmarkCalculateDiscount(b *testing.B) {
	// TODO: b.N la gi? Tai sao khong hardcode so lan chay?
	// SENIOR ASKS: Go quyet dinh b.N nhu the nao?
	// HINT: Go tang b.N cho den khi du do tin cay — thuong > 1s total
	for i := 0; i < b.N; i++ {
		CalculateDiscount(99.99, "gold")
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	p := Product{ID: 1, Name: "Test Product", Price: 29.99}
	b.ResetTimer() // TODO: Tai sao ResetTimer?
	// SENIOR ASKS: Neu khong ResetTimer, thoi gian setup co tinh vao benchmark khong?
	// HINT: Co! ResetTimer xoa timer truoc vong lap chinh
	b.ReportAllocs() // TODO: ReportAllocs de lam gi?
	// SENIOR ASKS: -benchmem flag khac gi b.ReportAllocs()?
	// HINT: -benchmem hien allocs cho moi benchmark, ReportAllocs chi cho 1
	
	for i := 0; i < b.N; i++ {
		json.Marshal(p)
	}
}

// TODO-[4]: Benchmark so sanh 2 implementation
// SENIOR ASKS: Tai sao can so sanh? Khi nao nen viet 2 benchmark?
// HINT: Khi refactor — dam bao version moi khong cham hon

func BenchmarkStringConcatenation(b *testing.B) {
	b.Run("Plus", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = "Hello" + " " + "World"
		}
	})
	b.Run("Sprintf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s %s", "Hello", "World")
		}
	})
	b.Run("Builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			sb.WriteString("Hello")
			sb.WriteString(" ")
			sb.WriteString("World")
			_ = sb.String()
		}
	})
}

// TODO-[5]: Parallel test
// SENIOR ASKS: t.Parallel() de lam gi? Khi nao nen dung, khi nao khong?
// HINT: Song phuong thu gio chay test, nhung chi dung khi test doc lap

func TestHandlerParallel(t *testing.T) {
	handler := setupHandler()
	
	t.Run("concurrent requests", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			t.Run(fmt.Sprintf("req-%d", i), func(t *testing.T) {
				t.Parallel() // TODO: Tai sao dat trong t.Run long nhau?
				// SENIOR ASKS: t.Parallel() o ngoai vong lap duoc khong?
				// HINT: Khong! Parallel phai o trong subtest
				req := httptest.NewRequest(http.MethodGet, "/products", nil)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				if rr.Code != http.StatusOK {
					t.Errorf("got %d", rr.Code)
				}
			})
		}
	})
}

// TODO-[6]: Test helper
// SENIOR ASKS: Tai sao can test helper? Khac gi voi helper binh thuong?
// HINT: t.Helper() — loai bo helper khoi stack trace, trace chi hien test that

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status: got %d, want %d", got, want)
	}
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`t.Fatal()` vs `t.Errorf()` — khac nhau gi? Khi nao dung cai nao?** (Hint: Fatal dung test ngay, Error ghi loi roi tiep tuc)

2. **`httptest.NewRecorder()` khong co actual network call. Vay no test duoc gi, khong test duoc gi?** (Hint: khong test duoc timeout, real network error)

3. **Benchmark chay `b.N` lan. Neu ham co side effect (vi du: ghi file), lam sao benchmark con chinh xac?** (Hint: ResetTimer khong reset side effect)

4. **`-race` flag phat hien race condition nhu the nao?** Co false positive khong? Chi phi performance khi chay -race?

5. **pprof CPU profile: `go test -cpuprofile=cpu.prof -bench=.`. File `.prof` doc bang gi?** `go tool pprof` co nhung lenh nao hru ich?

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Table-driven test voi >= 5 case, co error case
- [ ] TODO-[2] hoan thanh: Test HTTP handler voi ResponseRecorder, verify status + body
- [ ] TODO-[3] hoan thanh: Benchmark CalculateDiscount + JSONMarshal co ReportAllocs
- [ ] TODO-[4] hoan thanh: Sub-benchmark so sanh 3 cach concatenate string
- [ ] TODO-[5] hoan thanh: Parallel test 100 concurrent requests
- [ ] TODO-[6] hoan thanh: Test helper voi t.Helper()
- [ ] `go test -race` pass — khong co race condition
- [ ] `go test -bench=. -benchmem` chay duoc, co report allocations

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: POST /products voi body JSON hop le — status 201 — **happy path:** verify creation
- [ ] Test case: POST /products voi body khong phai JSON — status 400 — **error path:** invalid input
- [ ] Test case: GET /products/999 — status 404 — **boundary case:** not found
- [ ] Benchmark: So sanh json.Marshal vs json.NewEncoder — **performance:** cai nao allocations it hon?
- [ ] Test case: Concurrent POST + GET — verify khong race — **concurrency:** -race flag

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Test quality vs test quantity — 100 test nhanh hay 20 test chat luong cao? "Code coverage 100%" co y nghia gi?
2. **Neu requirement thay doi:** "Integration test voi real database" — ban se setup nhu the nao? (Hint: test database, migration rollback)
3. **Architecture decision:** Tai sao Go chon convention `TestXxx` + `BenchmarkXxx` thay vi annotation nhu `@Test`? Loi ich cua convention over configuration?

---

## Topic 03.6: Go Modules & Build

### User Story

> Khach hang (Product Owner) noi: "Setup CI/CD: build binary cho Linux, macOS, Windows. Embed static files vao binary. Binary phai co version info khi chay --version."
>
> Context: Ung dung can deploy da nen tang — dev team dung macOS, production chay Linux, co khach dung Windows. Static files (template, config) can embed de chi deliver 1 file.

### Acceptance Criteria

- [ ] Go module duoc khoi tao voi go mod init, quan ly dependency
- [ ] Cross-compile: GOOS/GOARCH cho linux/amd64, darwin/amd64, windows/amd64
- [ ] go:embed embed static files vao binary
- [ ] ldflags truyen version, commit hash, build time vao binary
- [ ] Binary in ra version info khi chay flag -version
- [ ] Build script (Makefile hoac shell script) tu dong hoa build da nen tang

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Hoi toi co junior build binary tren Mac, copy len Linux server — khong chay duoc. 'Bad executable format'. Co ban khong biet cross-compile. Go ho tro dieu nay cuc ky de chi voi 2 bien moi truong.
>
> Van de cot loi o day la: **Go build la single static binary — nhung phai biet cach build dung target va embed dung cach.**
>
> Cach toi phan ra:
> - Buoc 1: Go modules — khoi tao, tidy, vendor (neu can)
> - Buoc 2: Cross-compile voi GOOS/GOARCH
> - Buoc 3: go:embed cho static assets
> - Buoc 4: ldflags cho build-time variables
> - Buoc 5: Build script tu dong hoa
>
> O project truoc, toi setup CI build 6 target (3 OS x 2 arch) trong 30 giay. Go khong can toolchain rieng cho tung OS — cross-compile san."

#### TODO Comments (Code Skeleton)

```go
package main

import (
	"_ "embed" // TODO: Tai sao can import voi underscore?
	// SENIOR ASKS: embed package khong can dung truc tiep — vi sao van phai import?
	// HINT: Go compiler can nhin thay import de enable compile-time feature
	"fmt"
	"runtime"
)

// TODO-[1]: Build-time variables qua ldflags
// SENIOR ASKS: Tai sao dung var ma khong phai const?
// HINT: ldflags chi co the set gia tri cho variable tai build time, khong phai const

var (
	Version   = "dev"     // duoc ghi de boi -ldflags khi build
	Commit    = "unknown" // git commit hash
	BuildTime = "unknown" // thoi gian build
)

// TODO-[2]: go:embed static files
// SENIOR ASKS: //go:embed la comment hay directive? Tai sao bat dau bang //go:?
// HINT: La compiler directive — special comment format ma Go toolchain doc

//go:embed templates/*.html
var templatesFS embed.FS // TODO: Tai sao dung embed.FS ma khong phai string?
// SENIOR ASKS: embed.FS khac string hay []byte nhu the nao?
// HINT: embed.FS implement fs.FS — dung voi http.FS, template.ParseFS

//go:embed static/css/*
var cssFS embed.FS

// TODO-[3]: Serve embedded files qua HTTP
// SENIOR ASKS: embed.FS khac os.FS — lam sao dung voi http.FileServer?
// HINT: http.FS(embed.FS) convert sang http.FileSystem

func setupStaticHandler() http.Handler {
	// TODO: Tao handler serve file tu embed.FS
	// SENIOR ASKS: fs.Sub() de lam gi? Khi nao can?
	// HINT: Khi chi muon expose 1 subfolder cua embed.FS
	fsys, _ := fs.Sub(templatesFS, "templates")
	return http.FileServer(http.FS(fsys))
}

// TODO-[4]: Version info endpoint
func versionHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Tra ve JSON: version, commit, build_time, go_version, os/arch
	// SENIOR ASKS: runtime.Version() tra ve gi? runtime.GOOS/GOARCH?
	// HINT: Go version (go1.21.0), target OS va architecture luc build
	info := map[string]string{
		"version":    Version,
		"commit":     Commit,
		"build_time": BuildTime,
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
	}
	json.NewEncoder(w).Encode(info)
}

// TODO-[5]: Makefile build script
// SENIOR ASKS: Tai sao can Makefile? Khong dung shell script duoc khong?
// HINT: Makefile co dependency tracking, parallel build, convention

// Makefile skeleton (khong phai Go code, viet rieng):
/*
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.Version=$(VERSION) \
           -X main.Commit=$(COMMIT) \
           -X main.BuildTime=$(DATE) \
           -s -w # strip debug info, giam binary size

.PHONY: build build-all clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/app .

build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/app-linux-amd64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/app-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/app-windows-amd64.exe .

# TODO: -s -w giam binary size bao nhieu? Tai sao?
# SENIOR ASKS: Strip symbol table co anh huong gi den debug khong?
# HINT: Khong con stack trace symbol — nhung co the dung external symbol file
*/

func main() {
	// TODO: Khoi tao server, serve embedded files, in version
	fmt.Printf("App %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`go mod tidy` lam nhung gi chinh xac?** Co bao gio no xoa dependency ma ban dang dung khong? (Hint: no xoa unused, them missing — neu import ma khong dung function, van giu)

2. **Cross-compile: `GOOS=windows GOARCH=amd64 go build` — co can cai dat Windows SDK khong?** Tai sao Go co the lam dieu nay ma C++ khong? (Hint: Go toolchain co san runtime cho moi OS/arch)

3. **`go:embed` co gioi han kich thuoc file khong?** File 100MB co nen embed khong? Binary size tang co van de gi? (Hint: startup time, memory)

4. **`-ldflags "-s -w"` giam binary size ~30%. Nhung `go version -m binary` khong hoat dong sau do — tai sao?** Ban se chon trade-off nhu the nao cho production? (Hint: build voi -s -w cho release, giu lai cho staging de debug)

5. **`go env GOPROXY` la gi? Tai sao trong enterprise can set proxy rieng?** (Hint: firewall, audit dependency, speed)

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: 3 bien Version, Commit, BuildTime duoc khai bao
- [ ] TODO-[2] hoan thanh: go:embed directive dung, embed.FS khai bao dung
- [ ] TODO-[3] hoan thanh: Static file serve qua HTTP tu embedded FS
- [ ] TODO-[4] hoan thanh: Version info endpoint tra ve JSON day du
- [ ] TODO-[5] hoan thanh: Makefile build duoc 3 target, ldflags set version
- [ ] Chay `make build` tao binary co version info khi chay
- [ ] Binary size truoc/sau `-s -w` — biet duoc chenh lech

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: GET /version tra ve JSON voi version = "dev" khi chay `go test` — **default behavior:** khong ldflags trong test
- [ ] Test case: Embedded template file co the doc duoc — **verify:** fs.ReadFile duoc
- [ ] Test case: Build voi ldflags, verify binary chay `--version` in ra dung — **integration:** test build pipeline
- [ ] Test case: go:embed voi file khong ton tai — build fail — **error case:** verify compile-time check

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** go:embed vs doc file tu disk — khi nao embed, khi nao de file rieng? (Hint: single binary deploy vs dynamic update)
2. **Neu requirement thay doi:** "Can hot-reload template khong can rebuild binary" — ban se refactor nhu the nao?
3. **Architecture decision:** Tai sao Go chon model "single static binary" thay vi runtime nhu Java/Python? Loi va hai cua single binary deployment?

---

## Mini-Project: File/Text Toolkit API

### User Story

> Khach hang (Product Owner) noi: "API cho phep upload file, word count, grep, format text. Co SQLite backend, chi dung stdlib. Tao binary chay duoc, co test va benchmark."
>
> Context: Day la project tong hop Phase 3 — ket hop net/http, encoding/json, database/sql, os/io, testing, go modules. Stdlib only — khong dung framework, chi driver SQLite.

### Acceptance Criteria

- [ ] Upload file qua HTTP POST — luu vao SQLite hoac filesystem
- [ ] Word count endpoint: POST /api/wordcount voi text body → tra ve so tu, so dong, so ky tu
- [ ] Grep endpoint: POST /api/grep voi {text, pattern} → tra ve cac dong match
- [ ] Format endpoint: POST /api/format voi text → tra ve trimmed, normalized whitespace
- [ ] SQLite backend: danh sach file da upload, metadata (size, upload time, word count)
- [ ] Middleware: logging, recover (khong panic crash server)
- [ ] Test: table-driven cho handlers, integration test voi SQLite
- [ ] Benchmark: word count function voi input 1MB
- [ ] Build: Makefile, cross-compile, version info, go:embed cho HTML UI (neu co)

### Senior Thought-Process

#### Senior Thought-Process

**Senior nghi gi khi nhan requirement nay:**

> "Day la project tong hop het Phase 3. Khong hoc them concept moi — nhung phai ket hop 6 topic da hoc thanh 1 he thong hoan chinh. Ky nang quan trong: boc tach layer, viet interface, quan ly dependency.
>
> Van de cot loi o day la: **Biet tung phan khong du — phai biet rap noi lai.** API layer (http) → business layer (word count/grep) → storage layer (sqlite). Moi layer chi biet layer ke tiep qua interface.
>
> Cach toi phan ra project nay:
> - Layer 1: HTTP handlers — doc request, goi service, tra response
> - Layer 2: Service — business logic (wordcount, grep, format)
> - Layer 3: Repository — database access (CRUD file metadata)
> - Layer 4: Models — struct chung giua cac layer
>
> Goi y kien truc tu project that:
> ```
> main.go         → setup server, inject dependencies
> handler/        → HTTP handlers, chi goi service
> service/        → business logic, khong biet HTTP hay DB cu the
> repository/     → SQLite queries, khong biet HTTP
> model/          → struct chung
> ```
>
> Dieu quan trong: **handler khong goi db truc tiep. Service khong biet request/response HTTP.** Nhu vay moi test tung layer rieng duoc."

#### TODO Comments (Code Skeleton)

```go
// ============================================================
// FILE: model/file.go
// ============================================================
package model

import "time"

// TODO-[1]: Dinh nghia model
// SENIOR ASKS: Tai sao tao package model rieng? Khong de trong main duoc khong?
// HINT: De share giua handler/service/repository — khong tao dependency cycle

type FileRecord struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Size      int64     `json:"size" db:"size"`
	WordCount int       `json:"word_count" db:"word_count"`
	Content   string    `json:"-" db:"content"` // TODO: Tai sao json:"-"?
	// SENIOR ASKS: json:"-" co y nghia gi? Khac gi voi bo trong?
	// HINT: An field hoan toan khi marshal — khong xuat hien trong JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type WordCountResult struct {
	Words     int `json:"words"`
	Lines     int `json:"lines"`
	Chars     int `json:"chars"`
	Bytes     int `json:"bytes"`
}

type GrepRequest struct {
	Text    string `json:"text"`
	Pattern string `json:"pattern"`
}

type GrepResult struct {
	Matches []string `json:"matches"`
	Count   int      `json:"count"`
}

// ============================================================
// FILE: service/text.go
// ============================================================
package service

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// TODO-[2]: Word count logic
// SENIOR ASKS: Lam sao dem tu chinh xac? "Hello   World" la 2 tu hay 5?
// HINT: strings.Fields() split theo whitespace — xu ly duoc nhieu space

func WordCount(text string) model.WordCountResult {
	// TODO: Dem words, lines, chars, bytes
	// SENIOR ASKS: utf8.RuneCountInString() khac len() nhu the nao?
	// HINT: len = bytes, RuneCount = characters (quan trong cho tieng Viet!)
	
	lines := strings.Count(text, "\n") + 1
	if text == "" {
		lines = 0
	}
	words := len(strings.Fields(text))
	
	return model.WordCountResult{
		Words: words,
		Lines: lines,
		Chars: utf8.RuneCountInString(text),
		Bytes: len(text),
	}
}

// TODO-[3]: Grep logic
// SENIOR ASKS: regexp.Compile co the panic khong? Nen xu ly nhu the nao?
// HINT: Compile tra ve error — khong panic. MustCompile moi panic.

func Grep(text, pattern string) (model.GrepResult, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return model.GrepResult{}, err // TODO: Wrap error?
	}
	
	scanner := bufio.NewScanner(strings.NewReader(text))
	var matches []string
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			matches = append(matches, line)
		}
	}
	
	return model.GrepResult{
		Matches: matches,
		Count:   len(matches),
	}, scanner.Err()
}

// TODO-[4]: Format text
// SENIOR ASKS: "Normalize whitespace" co nghia gi cu the?
// HINT: strings.Fields() + strings.Join() = normalize spaces

func FormatText(text string) string {
	// TODO: Trim space, replace multiple spaces with single, normalize newlines
	return text // TODO: Implement
}

// ============================================================
// FILE: repository/file.go
// ============================================================
package repository

import (
	"context"
	"database/sql"
)

// TODO-[5]: Repository pattern
// SENIOR ASKS: Tai sao dung interface? Khong dung truc tiep *sql.DB duoc khong?
// HINT: De mock trong test — service test khong can database that

type FileRepository interface {
	Create(ctx context.Context, file *model.FileRecord) error
	GetByID(ctx context.Context, id int64) (model.FileRecord, error)
	List(ctx context.Context) ([]model.FileRecord, error)
	Delete(ctx context.Context, id int64) error
}

type sqliteRepo struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) Create(ctx context.Context, file *model.FileRecord) error {
	// TODO: INSERT va scan LastInsertId
	// SENIOR ASKS: RETURNING clause khac gi LastInsertId?
	// HINT: RETURNING (SQLite 3.35+) lay gia tri sau insert 1 cach — LastInsertId 2 cach
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO files (name, size, word_count, content) VALUES (?, ?, ?, ?)",
		file.Name, file.Size, file.WordCount, file.Content,
	)
	if err != nil {
		return err
	}
	file.ID, _ = result.LastInsertId()
	return nil
}

// TODO: Implement GetByID, List, Delete

// ============================================================
// FILE: handler/handler.go
// ============================================================
package handler

type Handler struct {
	fileRepo FileRepository // interface!
	svc      *service.TextService
}

func NewHandler(repo FileRepository, svc *service.TextService) *Handler {
	return &Handler{fileRepo: repo, svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// TODO: Register routes
	// SENIOR ASKS: Tai sao khong dung default mux ma truyen vao?
	// HINT: De test — tao mux rieng cho test, khong anh huong global
	mux.HandleFunc("POST /api/upload", h.handleUpload)
	mux.HandleFunc("POST /api/wordcount", h.handleWordCount)
	mux.HandleFunc("POST /api/grep", h.handleGrep)
	mux.HandleFunc("POST /api/format", h.handleFormat)
	mux.HandleFunc("GET /api/files", h.handleListFiles)
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse multipart form, lay file, luu content
	// SENIOR ASKS: r.ParseMultipartForm gioi han memory nhu the nao?
	// HINT: Tham so dau tien = max memory truoc khi ghi tam file
	
	// Parse form → doc file → word count → luu vao DB
	// TODO: Return JSON with file metadata
}

// TODO: Implement cac handler con lai

// ============================================================
// FILE: main.go
// ============================================================
package main

func main() {
	// TODO: Open DB, init schema
	// TODO: Create repository, service, handler
	// TODO: Setup routes voi middleware (logging, recover)
	// TODO: Start server voi graceful shutdown
}
```

#### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Repository pattern: `FileRepository` la interface co 4 method. Service phu thuoc vao interface hay implementation?** Dieu nay khac gi voi service truc tiep dung `*sql.DB`?

2. **Upload file 100MB qua HTTP — `r.ParseMultipartForm(32 << 20)` nghia la gi?** File lon hon 32MB se xay ra gi? (Hint: ghi tam file)

3. **Word count tieng Viet: "Xin chào thế giới" co may tu?** `strings.Fields()` dem dung khong? Tai sao can `utf8.RuneCountInString()` thay vi `len()`?

4. **Regex trong grep: Neu pattern la `(?i)hello`, co nghia la gi?** Co van de gi security khi cho client tu do nhap regex? (Hint: ReDoS — catastrophic backtracking)

5. **Project structure: `handler/` → `service/` → `repository/` — dependency huong xuong. Dieu gi xay ra neu repository import model tu handler?** (Hint: circular dependency — Go khong cho phep)

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Model struct dinh nghia day du voi JSON/DB tags
- [ ] TODO-[2] hoan thanh: WordCount dem dung words, lines, chars, bytes — ho tro Unicode
- [ ] TODO-[3] hoan thanh: Grep voi regex hoat dong, xu ly loi compile pattern
- [ ] TODO-[4] hoan thanh: FormatText normalize whitespace
- [ ] TODO-[5] hoan thanh: Repository interface + SQLite implementation
- [ ] Upload file hoat dong, tra ve metadata JSON
- [ ] `go test ./...` xanh — tat ca test pass
- [ ] `go test -bench=. -benchmem` co benchmark word count
- [ ] `go test -race` khong phat hien race condition
- [ ] Makefile build duoc binary cho >= 2 OS
- [ ] Binary co version info, co the serve embedded static files

### Test Checklist: Nhung gi ban nen tu viet test

- [ ] Test case: WordCount voi text trong → 0 tu, 0 dong — **boundary case:** empty input
- [ ] Test case: WordCount voi text tieng Viet "Xin chào" → 2 tu, 8 chars — **unicode handling:** verify RuneCount
- [ ] Test case: Grep voi pattern khong hop le → error — **error path:** invalid regex
- [ ] Test case: Upload file → verify DB co record — **integration test:** round-trip
- [ ] Test case: HTTP handler voi body khong phai JSON → 400 — **error path:** bad request
- [ ] Benchmark: WordCount voi 1MB text — do memory allocations — **performance baseline**
- [ ] Test case: Recover middleware — handler panic → server khong crash — **resilience:** verify 500 response

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Repository interface cho phep mock de test — nhung tang 1 lop abstraction. Khi nao project nho khong can interface? (Hint: < 3 method, khong can test → khong can interface)
2. **Neu requirement thay doi:** "Ho tro 1000 concurrent upload" — bottleneck se o dau? (Hint: SQLite chi 1 writer, HTTP multipart parsing memory)
3. **Architecture decision:** Tai sao chia 3 layer (handler/service/repo) ma khong gop lai? Khi nao nen gop, khi nao nen tach? So sanh voi MVC pattern ban co the biet.
4. **Bai hoc lon nhat tu project nay:** Ban hoc duoc gi ve cach to chuc Go project? Dieu gi khac voi Dart/Flutter project?

---

## Tong Ket Phase 3

### Cac Chu De Da Hoc

| Topic | Package Chinh | Ky Nang Cot Loi |
|-------|--------------|-----------------|
| HTTP Server | net/http | Handler interface, middleware chain, graceful shutdown |
| JSON | encoding/json | Struct tags, omitempty, NULL handling, custom marshal |
| Database | database/sql | Connection pool, prepared statement, rows.Close(), transaction |
| File I/O | os, io, bufio | Streaming, buffered reader, temp files, LimitReader |
| Testing | testing | Table-driven tests, httptest, benchmark, pprof |
| Build | go modules, go:embed, ldflags | Cross-compile, embed assets, build automation |

### Cac Quy Tac Vang Phase 3

1. **rows.Close() defer NGAY sau check error** — khong quen, khong doi
2. **Table-driven tests cho moi ham co nhieu case** — khong copy-paste test
3. **Stdlib first, framework sau** — biet net/http truoc khi dung Gin
4. **Single static binary** — Go build output de deploy
5. **Interface cho layer boundary** — handler/service/repo qua interface

### Checkpoint Cuoi Phase

Truoc khi vao Phase 4, dam bao:
- [ ] Viet duoc REST API CRUD chi dung stdlib
- [ ] Marshal/Unmarshal JSON voi struct tags, custom type
- [ ] Query SQLite voi prepared statement, transaction, rows.Close()
- [ ] Doc file lon (>100MB) ma khong vuot memory
- [ ] Viet table-driven test, benchmark, chay pprof
- [ ] Build binary da nen tang voi go:embed va ldflags
- [ ] Mini-project File/Text Toolkit chay duoc, test xanh, co benchmark

---

*"Stdlib la nen tang. Framework la cong cu. Xay nen tang vung, dung cong cu moi hieu qua."
— Senior Go Engineer*
