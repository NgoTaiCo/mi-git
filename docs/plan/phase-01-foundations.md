# Phase 1: Core Language Foundations (Tuần 1-3)

> **Meta:** 13 topic, 3 tuần, 1 mini-project cuối phase. Mỗi topic tuân theo format Requirement-Simulation: User Story → Senior Thought-Process → TODO Skeleton → Socratic Questions → Checklists.
>
> **Nguyên tắc:** Code skeleton trong file này **KHÔNG chạy được** — chỉ là khung để bạn điền vào. Senior sẽ không bao giờ viết code hoàn chỉnh cho bạn.

---

## Week 1: Go Environment & Variables

> **Mục tiêu tuần:** Từ zero → viết được CLI nhỏ có validate input. Nắm chắc toolchain, biến, zero value, printf, parsing, constants.

---

## Topic 1: Khởi tạo Project Go đầu tiên

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần bạn setup một project Go mới. Team sẽ có nhiều ngườI làm việc trên repo này, nên cần có cách build và run rõ ràng. Tôi muốn thấy một binary chạy được cuối ngày hôm nay."*
>
> Context: Bạn mới join team. Máy bạn đã cài Go. Cần tạo module, viết `main.go` đầu tiên, build thành binary, và chạy.

### Acceptance Criteria
- [ ] Tạo được Go module với `go mod init`
- [ ] Viết `main.go` có `package main` và `func main()`
- [ ] Build thành binary bằng `go build` không warning
- [ ] Binary chạy được trên máy hiện tại
- [ ] Giải thích được sự khác nhau giữa `go run`, `go build`, `go install`
- [ ] Biết dùng `go fmt` để format code

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Nếu tôi nhận ticket này, điều đầu tiên tôi nghĩ đến là: module name đặt sao cho đúng convention.
> Go dùng module path dạng `github.com/username/projectname` — nhưng nếu là internal project,
> có thể dùng `company.local/team/project`. Quan trọng là PHẢI có module path rõ ràng.
>
> "Vấn đề cốt lõi ở đây là: hiểu Go workspace model. Khác Dart/JS, Go không có `node_modules`.
> Dependencies được cache ở `$GOPATH/pkg/mod`. Module file `go.mod` là source of truth.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Tạo thư mục project → 2. `go mod init` → 3. Viết main.go → 4. `go build` → 5. Chạy binary
>
> "Hồi tôi mới học Go, tôi tưởng `go run` và `go build` giống nhau. Sai lầm.
> `go run` compile vào temp dir rồi chạy. `go build` tạo binary ở thư mục hiện tại.
> Trong production, bạn chỉ dùng `go build` hoặc `go install`.
```

#### TODO Comments (Code Skeleton)
```go
// TODO-[1]: Khởi tạo module
// SENIOR ASKS: Module path nên đặt như thế nào? Tại sao không nên đặt tên đơn giản như "myapp"?
// HINT: Nghĩ về uniqueness — 2 project cùng tên "myapp" thì Go phân biệt thế nào?

// $ mkdir ~/playground/go-cli-toolkit && cd ~/playground/go-cli-toolkit
// $ go mod init ???

// TODO-[2]: Viết main.go đầu tiên
// SENIOR ASKS: Tại sao file cần `package main`? Chuyện gì xảy ra nếu đặt `package hello`?
// HINT: Go compiler tìm `func main()` ở đâu để bắt đầu chương trình?

package main

import "fmt"

func main() {
	// TODO-[2a]: In ra "Hello, Go!" với fmt.Println
	// SENIOR ASKS: fmt.Println vs fmt.Printf khác gì? Khi nào dùng cái nào?
	// HINT: Nếu chỉ in string đơn giản, cái nào ít characters hơn?

	// TODO-[2b]: Thêm 1 biến `name := "Go"` rồi in ra dùng fmt.Printf
	// SENIOR ASKS: `%s` là gì? Còn `%v` thì sao?
	// HINT: %s chỉ dùng cho string. %v là "default format" — nhưng có nên lạm dụng không?
}

// TODO-[3]: Build và chạy
// SENIOR ASKS: `go build` tạo binary tên gì theo mặc định? Làm sao đặt tên khác?
// HINT: Thử `go build -o hello`
// SENIOR ASKS: Binary output có dependencies runtime không? Tại sao điều này quan trọng khi deploy?
// HINT: Go compile thành static binary — nhưng có exception nào không?

// TODO-[4]: Khám phá thư mục project
// SENIOR ASKS: Sau khi chạy `go mod init`, file nào được tạo ra? Nội dung của nó là gì?
// HINT: Đọc file go.mod — bạn thấy gì?
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Tại sao Go bắt buộc `package main` cho executable? Nếu bạn quen Java/C#, điều gì tương đương?
2. `go run` compile vào đâu? Tại sao `go run` không phù hợp cho production?
3. Nếu tôi đổi tên file từ `main.go` sang `hello.go`, chương trình có chạy được không? Tại sao?
4. Go binary static linking có nghĩa là gì? So sánh với binary C++ dynamic linking?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Đã tạo module với path đúng convention, `go.mod` tồn tại
- [ ] TODO-[2] hoàn thành: `main.go` in ra string, dùng cả `Println` và `Printf`
- [ ] TODO-[3] hoàn thành: Binary được tạo, chạy được, không phụ thuộc Go runtime
- [ ] TODO-[4] hoàn thành: Giải thích được nội dung `go.mod` và `go.sum` (nếu có)

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `go mod init` với invalid path — Go báo lỗi gì?
- [ ] Test case: Xóa `package main` dòng đầu — compile error là gì?
- [ ] Test case: Đổi `func main()` thành `func mainx()` — linker error là gì?
- [ ] Test case: `go build` trên máy khác cùng OS/arch — binary chạy không?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu team có 5 ngườI cùng code trên module này, `go.mod` có conflict không? Tại sao?
2. Binary Go có thể chạy trên Docker `alpine` image không? Cần gì không?
3. Tại sao `GOPATH` mode (Go <1.11) bị thay thế bởi Go modules? Lợi ích gì?

---

## Topic 2: Variables, Types & Zero Values

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần lưu trữ thông tin ngườI dùng: tên, tuổi, nhiệt độ cơ thể, trạng thái bật/tắt. Hệ thống phải xử lý đúng khi chưa có dữ liệu — đừng để null pointer exception như Java."*
>
> Context: Bạn đang viết module khởi tạo dữ liệu patient. Cần hiểu rõ type system của Go để chọn đúng.

### Acceptance Criteria
- [ ] Khai báo biến bằng cả `var`, `:=`, và `var (...)` block
- [ ] Chọn đúng type cho từng dữ liệu: string, int, float64, bool
- [ ] Giải thích được zero value của mỗi type — và tại sao Go thiết kế vậy
- [ ] Phân biệt rõ `var` vs `:=` — khi nào dùng cái nào
- [ ] Tránh được lỗi redeclaration (gán lại bằng `:=` trong cùng scope)
- [ ] Hiểu type inference: `i := 42` là type gì? `f := 3.14` là type gì?

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Nếu tôi nhận ticket này, điều đầu tiên tôi nghĩ đến là: zero value.
> Hồi tôi mới học Go, tôi tưởng zero value là bug. Sao string lại là ""?
> Sao int lại là 0? Tôi quen Dart/JS với null — thấy "thiếu dữ liệu" phải là null chứ.
> Sai. Go design: variables are always initialized. Không có uninitialized variable.
>
> "Vấn đề cốt lõi ở đây là: hiểu tại sao Go không có 'undefined' như JS hay
> nullable by default như Dart. Đây là design decision — không phải thiếu sót.
> Nó giúp code an toàn hơn, ít null check hơn.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Khai báo biến cho từng field → 2. In ra zero value → 3. Gán giá trị → 4. Test redeclaration
>
> "Hồi tôi ở project fintech, tôi từng gặp bug: developer dùng `var amount float64`
> rồi tính toán mà quên check `amount == 0`. Trong hệ thống đó, 0.0 có nghĩa là
> 'chưa nhập' — nhưng business logic hiểu thành 'free'. Từ đó tôi luôn dùng
> pointer hoặc sentinel value cho optional fields. Nhưng đó là lesson nâng cao —
> ở phase này, chỉ cần nắm zero value là đủ.
```

#### TODO Comments (Code Skeleton)
```go
package main

import "fmt"

func main() {
	// TODO-[1]: Khai báo biến bằng var — zero value demonstration
	// SENIOR ASKS: Không gán giá trị, in ra biến — bạn thấy gì? Tại sao Go làm vậy?
	// HINT: Go spec: "When storage is initialized, it is given a zero value." Đọc spec có làm sao đâu.

	var patientName string   // TODO: In ra. Zero value là gì?
	var patientAge int       // TODO: In ra. Zero value là gì?
	var temperature float64  // TODO: In ra. Zero value là gì?
	var isActive bool        // TODO: In ra. Zero value là gì?

	// TODO-[2]: Khai báo bằng := — short declaration
	// SENIOR ASKS: `name := "Alice"` — compiler infer type gì? Nếu tôi muốn int32 thay vì int, làm sao?
	// HINT: := luôn chọn type "mặc định": int, float64, string. Muốn type khác → phải dùng var.

	// TODO-[3]: Khai báo nhiều biến bằng var block
	// SENIOR ASKS: Tại sao var block gọn hơn khai báo riêng lẻ? Khi nào nên dùng?
	// HINT: Nhìn code đẹp hơn, nhóm biến có liên quan lại với nhau.

	// var (
	//     name string = "Alice"  // TODO: explicit type + value
	//     age int    = 30        // TODO: có thể bỏ `int` không?
	//     // ...
	// )

	// TODO-[4]: Redeclaration trap
	// SENIOR ASKS: Đoạn code sau compile được không? Nếu không, lỗi gì?
	// HINT: := có quy tắc "at least one new variable" — nhưng cẩn thận, nó dễ gây bug shadowing.

	// x := 10
	// x := 20  // TODO: Compile error? Tại sao?

	// x := 10
	// x, y := 20, 30  // TODO: Compile được không? Tại sao? x bây giờ là bao nhiêu?

	// TODO-[5]: Type inference quiz
	// SENIOR ASKS: Mỗi biến sau có type gì? Giải thích.
	// HINT: Go spec quy định rõ type của untyped constant khi gán.

	// a := 42          // TODO: type gì?
	// b := 3.14        // TODO: type gì?
	// c := 3.0         // TODO: type gì? float64 hay int?
	// d := 'A'         // TODO: type gì? rune hay string?

	// TODO-[6]: Mixed type assignment
	// SENIOR ASKS: int + float64 — Go tự ép kiểu không? Tại sao?
	// HINT: Go KHÔNG tự động promote type như C/Java. Đây là design decision — nghĩ xem tại sao.

	// i := 5
	// f := 3.14
	// sum := i + f  // TODO: Compile? Fix thế nào?
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Tại sao Go không tự động convert `int` → `float64` khi cộng? C lại làm được — điều gì Go trade off?
2. Zero value của slice là gì? Có phải `nil` không? Và `nil` slice khác gì empty slice `[]T{}`?
3. Khi nào nên dùng `var name string` thay vì `name := ""`? Có khác biệt thực sự không?
4. `:=` trong if block có thể "shadow" biến ngoài không? Tìm ví dụ và giải thích tại sao điều này nguy hiểm.
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Đã in zero value của 4 types cơ bản, giải thích được
- [ ] TODO-[2] hoàn thành: Dùng `:=` đúng, biết type được infer
- [ ] TODO-[3] hoàn thành: Dùng var block gọn gàng
- [ ] TODO-[4] hoàn thành: Thử redeclaration, hiểu shadowing rule
- [ ] TODO-[5] hoàn thành: Xác định đúng type của mỗi inference
- [ ] TODO-[6] hoàn thành: Giải thích tại sao Go không tự promote type

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: Khai báo `var x int` rồi in ra — zero value là 0, không phải null/undefined
- [ ] Test case: `s := ""` vs `var s string` — kết quả giống nhau, nhưng use case khác nhau
- [ ] Test case: `x, y := 1, 2; x, y = y, x` — swap hoạt động không?
- [ ] Test case: `x := 10; if true { x := 20; fmt.Println(x) }; fmt.Println(x)` — shadowing output?
- [ ] Boundary case: `var i int8 = 128` — compile error vì overflow?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu bạn đến từ Dart/JS, zero value thay thế `null` trong Go như thế nào? Khi nào vẫn cần `nil`?
2. Nếu requirement thay đổi: "tuổi có thể không biết" — bạn dùng type gì? `int` không đủ nữa?
3. Tại sao Go chọn `int` là 64-bit trên 64-bit arch nhưng 32-bit trên 32-bit? Khi nào nên dùng `int64` explicit?

---

## Topic 3: Formatted Output (printf)

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần in báo cáo nhiệt độ ra console. Phải đẹp, có căn cột, 2 chữ số thập phân. NgườI dùng cuối nhìn phải hiểu ngay — đừng để họ thấy debug output."*
>
> Context: Bạn đang viết tool diagnostic cho healthcare app. Output phải professional.

### Acceptance Criteria
- [ ] Dùng đúng verb: `%s`, `%d`, `%f`, `%.2f`, `%T`, `%v`, `%#v`
- [ ] Căn cột: `%10s`, `%-10s`, `%10.2f`
- [ ] Format số: 2 chữ số thập phân, leading zeros nếu cần
- [ ] Tạo được bảng đẹp với `fmt.Printf`
- [ ] Phân biệt `Printf` (format) vs `Println` (simple) vs `Sprintf` (return string)
- [ ] Không dùng `%v` cho production output — chỉ dùng cho debug

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Formatted output nghe đơn giản, nhưng đây là UX. NgườI dùng cuối nhìn thấy output
> của tool bạn — nếu lộn xộn, họ nghĩ tool buggy. Professional output = professional tool.
>
> "Vấn đề cốt lõi ở đây là: chọn đúng verb cho đúng type, căn cột để dễ đọc,
> và QUAN TRỌNG NHẤT — phân biệt output cho user vs output cho developer debug.
>
> "Tôi sẽ phân rã thành các bước:
> 1. In từng giá trị với verb phù hợp → 2. Căn cột → 3. Format số → 4. Tạo table
>
> "Hồi tôi ở project trước, có junior dùng `%v` cho mọi thứ. Khi struct có pointer,
> `%v` in ra memory address — ngườI dùng nhìn thấy `0xc000012345` và panic.
> Từ đó tôi có rule: production output → explicit verbs, never %v.
```

#### TODO Comments (Code Skeleton)
```go
package main

import "fmt"

func main() {
	// TODO-[1]: Cơ bản — in từng type với verb đúng
	// SENIOR ASKS: `fmt.Printf("%s", 42)` — chuyện gì xảy ra? Compile error hay runtime error?
	// HINT: Go kiểm tra type tại compile time cho Printf. Thử đi, bạn sẽ ngạc nhiên.

	name := "Alice"
	age := 30
	temp := 36.5789
	active := true

	// TODO: Dùng fmt.Printf với đúng verb cho từng biến
	// fmt.Printf("Name: ???\n", name)
	// fmt.Printf("Age: ???\n", age)
	// fmt.Printf("Temp: ???\n", temp)   // muốn 2 chữ số thập phân
	// fmt.Printf("Active: ???\n", active)

	// TODO-[2]: Căn cột — tạo bảng
	// SENIOR ASKS: `%-10s` vs `%10s` khác gì? Khi nào căn trái, khi nào căn phải?
	// HINT: Text thường căn trái, số thường căn phải. Nghĩ xem tại sao?

	// Output mong muốn:
	// Name       Age    Temp    Status
	// Alice      30     36.58   true
	// Bob        25     37.10   false
	// Charlotte  28     36.42   true

	// TODO: Viết code tạo bảng trên với fmt.Printf

	// TODO-[3]: Debug vs Production output
	// SENIOR ASKS: Tại sao không dùng %v cho production? Thử %v với struct xem sao.
	// HINT: %v là "default format" — nó không control được. Production cần control.

	type Patient struct {
		Name string
		Age  int
	}
	p := Patient{Name: "Alice", Age: 30}
	// TODO: In p bằng %v, %+v, %#v — khác nhau gì? Cái nào phù hợp cho debug? Cái nào cho user?

	// TODO-[4]: fmt.Sprintf — tạo string thay vì in ra stdout
	// SENIOR ASKS: Khi nào cần Sprintf thay vì Printf? Đưa 2 use case.
	// HINT: Nếu bạn cần log string, hoặc trả về error message — không muốn in ra stdout.

	// TODO-[5]: Zero-padding và width
	// SENIOR ASKS: `%04d` nghĩa là gì? Dùng khi nào?
	// HINT: ID "0001" trông professional hơn "1" trong báo cáo.

	// TODO: In ID patient dạng "PT-0001", "PT-0042"
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. `fmt.Printf` kiểm tra type mismatch tại compile time hay runtime? Tại sao điều này quan trọng?
2. `%T` in ra gì? Khi nào bạn thực sự cần nó trong production code?
3. Tại sao `%v` với struct chứa slice pointer lại in ra memory address? Có cách nào custom không?
4. `fmt.Sprintf` allocate memory — trong hot loop có vấn đề performance không? Có alternative không?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Dùng đúng verb cho từng type, không dùng `%v` generic
- [ ] TODO-[2] hoàn thành: Bảng output thẳng hàng, dễ đọc
- [ ] TODO-[3] hoàn thành: Phân biệt được `%v`/`%+v`/`%#v`, biết khi nào dùng cái nào
- [ ] TODO-[4] hoàn thành: Dùng `Sprintf` để tạo string không in ra stdout
- [ ] TODO-[5] hoàn thành: Format số với zero-padding

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `%s` với `int` — compile error (nếu static analysis) hay runtime fmt failure
- [ ] Test case: `%.2f` với `3.1` → output là `3.10` hay `3.1`? (hint: width includes decimal)
- [ ] Test case: `%10s` với string dài hơn 10 chars — bị truncate hay vẫn hiện full?
- [ ] Boundary case: `Printf` với nhiều arguments hơn format verbs — chuyện gì xảy ra?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu bảng có column width dynamic (không biết trước độ dài string), `fmt.Printf` còn đủ không?
2. `fmt.Printf` dùng reflection không? Nếu có, performance impact là gì?
3. Nếu requirement thêm: output phải support Unicode (tiếng Việt) — `%10s` vẫn căn đúng không?

---

## Topic 4: Numbers & String Parsing

### User Story
> Khách hàng (Product Owner) nói: *"Tool cần đọc input từ ngườI dùng là số nhiệt độ. Tính toán xong in ra kết quả. Nếu user nhập chữ thay vì số, báo lỗi rõ ràng — đừng để chương trình crash với stack trace."*
>
> Context: Bạn đang viết phần core của `convert` command. Input đến từ CLI args, không phải từ database. Mọi input từ bên ngoài đều là "bẩn" — phải validate.

### Acceptance Criteria
- [ ] Parse string → number: `strconv.Atoi`, `strconv.ParseFloat`, `strconv.ParseInt`
- [ ] Handle error — không `panic`, không bỏ qua lỗi
- [ ] Validate input: reject empty, reject non-numeric, reject negative nếu không hợp lệ
- [ ] Dùng error message rõ ràng — ngườI dùng cuối hiểu được
- [ ] Explicit type conversion giữa numeric types
- [ ] Hiểu truncation vs rounding: `int(2.9)` = 2, không phải 3

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Parse input từ user là một trong những nguồn bug nhiều nhất trong CLI tool.
> NgườI dùng nhập gì cũng có thể: chữ, số âm, số quá lớn, empty string, UTF-8 emoji.
> Code của bạn PHẢI xử lý tất cả.
>
> "Vấn đề cốt lõi ở đây là: Go có 2 giá trị return — `(result, error)`.
> Khác Dart/JS, Go không throw exception. Bạn phải check error mỗi lần parse.
> Đây là một phần của Go error handling philosophy: explicit beats implicit.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Parse string → int → 2. Parse string → float → 3. Validate → 4. Error handling
>
> "Hồi tôi ở project healthcare, có 1 lần developer dùng `strconv.Atoi` parse user input
> rồi dùng kết quả mà không check error. Khi user nhập "abc", `Atoi` trả về `0, error`.
> Developer bỏ qua error → hệ thống ghi nhận nhiệt độ 0°C → bác sĩ panic vì tưởng bệnh nhân
> hạ thân nhiệt. Một dòng `if err != nil` bị thiếu → hậu quả nghiêm trọng.
```

#### TODO Comments (Code Skeleton)
```go
package main

import (
	// TODO-[1]: Import đúng packages
	// SENIOR ASKS: strconv có những hàm nào cho parsing? Khác gì với fmt.Sscanf?
	// HINT: strconv chuyên parse string↔number. Đơn giản, nhanh, explicit.
	"fmt"
	"strconv"
)

// parseTemperature nhận input string, trả về float64 và error
// SENIOR ASKS: Tại sao return (float64, error) thay vì chỉ float64?
// HINT: Go không có exception. Error là value, phải trả về explicit.
func parseTemperature(input string) (float64, error) {
	// TODO-[2]: Parse string → float64
	// SENIOR ASKS: ParseFloat cần bitSize 32 hay 64? Khi nào dùng cái nào?
	// HINT: Nếu bạn lưu vào float64, dùng bitSize 64. Consistency quan trọng.

	// TODO-[3]: Validate parsed value
	// SENIOR ASKS: Nhiệt độ cơ thể hợp lệ là bao nhiêu? Làm sao reject giá trị vô lý?
	// HINT: Domain validation — không chỉ check parse được, còn phải check hợp lý.

	// TODO-[4]: Trả về error rõ ràng
	// SENIOR ASKS: fmt.Errorf vs errors.New — khi nào dùng cái nào?
	// HINT: errors.New cho static message. fmt.Errorf khi cần embed value vào message.

	return 0, fmt.Errorf("TODO: implement")
}

func main() {
	// TODO-[5]: Test với nhiều input
	// SENIOR ASKS: Bạn cần test bao nhiêu case để confident? Liệt kê.
	// HINT: Happy path + error path. Ít nhất 6 case: valid, empty, non-numeric, negative, too large, decimal.

	testInputs := []string{"36.5", "abc", "", "-5", "999", "36,5"}
	for _, input := range testInputs {
		// TODO: Gọi parseTemperature, handle error, in kết quả đẹp
	}

	// TODO-[6]: Explicit conversion demo
	// SENIOR ASKS: int(2.9) = ? int(-2.9) = ? Tại sao không phải round?
	// HINT: Go truncate toward zero. Không phải round, không phải floor.

	// TODO-[7]: strconv.Atoi vs strconv.ParseInt
	// SENIOR ASKS: Atoi trả về int. ParseInt trả về int64. Khi nào dùng cái nào?
	// HINT: Atoi gọn hơn nhưng ít control. ParseInt cho phép chỉ định base và bitSize.
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Tại sao Go dùng `(value, error)` thay vì exceptions? Lợi và hại của cách tiếp cận này?
2. `strconv.ParseFloat("36,5", 64)` — dấu phẩy thay vì chấm. Kết quả là gì? Làm sao handle locale?
3. Nếu input là `" 36.5 "` (có space), ParseFloat hoạt động không? Nếu không, fix thế nào?
4. `int(float64(maxInt))` — chuyện gì xảy ra? Tại sao đây là silent bug nguy hiểm?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Import đúng packages (fmt, strconv, errors)
- [ ] TODO-[2] hoàn thành: ParseFloat đúng với bitSize 64
- [ ] TODO-[3] hoàn thành: Domain validation (range check)
- [ ] TODO-[4] hoàn thành: Error message rõ ràng, dùng fmt.Errorf/errors.New đúng
- [ ] TODO-[5] hoàn thành: Test với >= 6 input khác nhau, không panic
- [ ] TODO-[6] hoàn thành: Hiểu truncation behavior
- [ ] TODO-[7] hoàn thành: Phân biệt Atoi vs ParseInt

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `"36.5"` → 36.5 (happy path)
- [ ] Test case: `"abc"` → error (non-numeric)
- [ ] Test case: `""` → error (empty string)
- [ ] Test case: `"-5"` → reject hoặc accept tùy domain (negative)
- [ ] Test case: `"999"` → reject nếu out of reasonable range
- [ ] Test case: `"36,5"` → error (European decimal separator)
- [ ] Test case: `" 36.5 "` → có space — cần TrimSpace trước khi parse?
- [ ] Boundary case: `"1e309"` → overflow, error gì?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu bạn cần support cả dấu phẩy (European format), thiết kế hàm parse thay đổi thế nào?
2. `strconv.ParseFloat` dùng algorithm gì? Có precision issue không?
3. Tại sao Go không có `try/catch`? Bạn có thấy `(value, error)` verbose không? Cách nào giảm boilerplate?

---

## Topic 5: Constants & iota

### User Story
> Khách hàng (Product Owner) nói: *"Hệ thống đơn hàng của tôi có các trạng tháI: Pending, Confirmed, Shipped, Delivered, Cancelled. Tôi không muốn dùng magic number 0, 1, 2 trong code. Cần type-safe — đừng để ai đó truyền số 999 vào hàm nhận status."*
>
> Context: Bạn đang thiết kế domain model cho e-commerce module. Constants cần readable, type-safe, và có thể convert to/from string.

### Acceptance Criteria
- [ ] Dùng `iota` tạo enum-like constants với named type
- [ ] Type-safe: không thể gán `int` vào `OrderStatus` trực tiếp
- [ ] Có `String()` method để in ra tên readable ("Pending" thay vì "0")
- [ ] Phân biệt typed vs untyped constant
- [ ] Biết giới hạn: const không thể là slice, map, struct mutable
- [ ] Áp dụng pattern: loại bỏ magic numbers khỏi codebase

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Magic numbers là kẻ thù của maintainable code. `if status == 2` — 2 là gì?
> Shipped? Delivered? Phải mở comment hoặc document mới biết. Hỏng rồi.
>
> "Vấn đề cốt lõi ở đây là: Go không có enum kiểu C#/Java. Nhưng Go có
> constant + iota + named type — kết hợp lại thì mạnh hơn enum của nhiều ngôn ngữ.
> Type-safe, zero-cost, compile-time checked.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Tạo named type → 2. iota constants → 3. String() method → 4. Type safety demo
>
> "Hồi tôi refactor 1 codebase cũ, tôi thay 30 chỗ magic number bằng constants.
> Catch được 2 bug: một chỗ `if status == 3` nhưng status chỉ có 0-2 (out of range),
> và một chỗ `status = 0` nhưng 0 không phải valid initial state. Constants + type
> giúp compiler bắt được những lỗi này.
```

#### TODO Comments (Code Skeleton)
```go
package main

import "fmt"

// TODO-[1]: Tạo named type cho OrderStatus
// SENIOR ASKS: Tại sao không dùng `type OrderStatus string` mà dùng `int`?
// HINT: `int` cho phép so sánh `<`, `>=` — hữu ích cho priority check. String thì không.
// Nhưng `string` cũng có lợi: dễ serialize. Trade-off là gì?

type OrderStatus int

// TODO-[2]: Dùng iota tạo constants
// SENIOR ASKS: iota bắt đầu từ 0 — nếu tôi muốn bắt đầu từ 1, làm sao?
// HINT: `iota + 1` hoặc `iota(1)` — cái nào đúng syntax?

const (
	// TODO: Pending = iota
	// TODO: Confirmed
	// TODO: Shipped
	// TODO: Delivered
	// TODO: Cancelled
)

// TODO-[3]: String() method — làm cho fmt.Println in ra tên thay vì số
// SENIOR ASKS: Đây là implementation của interface gì trong Go?
// HINT: fmt.Stringer — một trong những interface quan trọng nhất trong Go.
// Go dùng duck typing: nếu type có String() string → nó là Stringer.

// func (s OrderStatus) String() string {
//     switch s {
//     case Pending: return "Pending"
//     // TODO: Các case còn lại
//     default: return fmt.Sprintf("OrderStatus(%d)", s)
//     }
// }

// TODO-[4]: Type safety demo
// SENIOR ASKS: Đoạn code sau compile được không? Nếu không, lỗi gì?
// HINT: Named type và underlying type khác nhau — đây chính là type safety.

// var s OrderStatus = 2        // TODO: Compile? Tại sao?
// var i int = 2
// var s2 OrderStatus = i      // TODO: Compile? Tại sao?
// var s3 OrderStatus = OrderStatus(i)  // TODO: Compile? Tại sao?

// TODO-[5]: Untyped constant — đặc tính độc đáo của Go
// SENIOR ASKS: const MaxSize = 100 — typed hay untyped? Nó gán vào int8 được không?
// HINT: Untyped constant "adapt" vào context. Nhưng nếu vượt range → compile error.

// const MaxSize = 100
// var tiny int8 = MaxSize  // TODO: Compile? Tại sao?
// const Huge = 1000
// var tiny2 int8 = Huge    // TODO: Compile? Tại sao?

// TODO-[6]: Bitflag pattern với iota
// SENIOR ASKS: `1 << iota` nghĩa là gì? Dùng khi nào?
// HINT: Bitwise shift tạo powers of 2: 1, 2, 4, 8... Dùng OR để combine flags.

// type Permission int
// const (
//     Read Permission = 1 << iota
//     Write
//     Execute
// )

func main() {
	// TODO: Tạo vài OrderStatus, in ra, demo type safety
	// TODO: Test bitflag combination
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. `iota` reset về 0 khi nào? Nếu tôi có 2 const block liên tiếp, iota thứ 2 bắt đầu từ mấy?
2. Tại sao Go không cho `const` là slice hoặc map? Nghĩ về compile-time vs runtime.
3. `String()` method của fmt.Stringer — tại sao Go không đặt tên là `ToString()` như Java?
4. Nếu tôi gán `OrderStatus(999)` — compile được không? Chuyện gì xảy ra ở `String()`?
5. `const Pi = 3.14159` — tại sao có thể gán vào cả `float32` lẫn `float64`?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Named type OrderStatus được tạo
- [ ] TODO-[2] hoàn thành: 5 constants dùng iota, bắt đầu từ 0
- [ ] TODO-[3] hoàn thành: String() method implement fmt.Stringer
- [ ] TODO-[4] hoàn thành: Demo type safety — `int` không gán trực tiếp vào `OrderStatus`
- [ ] TODO-[5] hoàn thành: Giải thích typed vs untyped constant
- [ ] TODO-[6] hoàn thành: Bitflag pattern với Permission

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `fmt.Println(Pending)` in ra "Pending" chứ không phải "0"
- [ ] Test case: `if status == Pending` hoạt động đúng
- [ ] Test case: `var x int = 2; var s OrderStatus = x` — compile error
- [ ] Test case: `OrderStatus(999).String()` — không panic, trả về default format
- [ ] Test case: `Read | Write` = 3, kiểm tra bit combination
- [ ] Boundary case: `iota` trong const block thứ 2 — reset về 0 hay tiếp tục?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu requirement thêm 10 status nữa — `String()` method dài ra. Có cách nào auto-generate không?
2. So sánh Go iota enum với Dart enum class — cái nào type-safe hơn? Cái nào flexible hơn?
3. Tại sao constants trong Go "zero-cost"? Có runtime overhead không?
4. Khi nào nên dùng `const` vs khi nào nên dùng `var` cho config values?

---

## Topic 6: CLI Tool — convert

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần tool command-line chuyển đổI đơn vị nhiệt độ. Chạy như thế này: `convert --from=celsius --to=fahrenheit 100` → output `100°C = 212.00°F`. Nếu nhập sai, báo lỗi rõ, không được crash."*
>
> Context: Đây là mini-project tổng hợp Week 1. Bạn phải dùng mọi thứ đã học: variables, parsing, formatted output, constants, error handling.

### Acceptance Criteria
- [ ] Parse CLI arguments (`os.Args` hoặc `flag` package)
- [ ] Validate: unit hợp lệ (celsius, fahrenheit, kelvin), value là số
- [ ] Convert đúng công thức
- [ ] Output đẹp: `100°C = 212.00°F`
- [ ] Error message rõ ràng, không panic
- [ ] Help text khi không đủ arguments hoặc `--help`
- [ ] Exit code đúng: 0 cho success, 1 cho error

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Đây là ticket đầu tiên có thực — không phải exercise lý thuyết nữa.
> Tôi sẽ approach như 1 ticket thực tế: đọc requirement → phân rã → implement → test.
>
> "Vấn đề cốt lõi ở đây là: tách logic thành functions rõ ràng.
> Parse → Validate → Compute → Render. Không gộp tất cả vào main().
> Đây là pattern sẽ lặp lại trong mọi CLI tool bạn viết.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Parse args → 2. Validate input → 3. Convert → 4. Format output → 5. Handle errors
>
> "Hồi tôi viết CLI tool đầu tiên, tôi gộp tất cả vào main(). 200 dòng,
> nested if 4 cấp, không test được. Senior review xong bảo: 'Phân rã đi em.'
> Tôi refactor thành 4 functions — code giống nhau nhưng đọc được, test được.
> Lesson: function nhỏ > function lớn. Luôn luôn.
```

#### TODO Comments (Code Skeleton)
```go
package main

import (
	"fmt"
	"os"
	"strconv"
	// TODO-[1]: Import thêm packages cần thiết
	// SENIOR ASKS: strings package có hàm nào hữu ích cho CLI argument parsing?
	// HINT: strings.ToLower cho case-insensitive unit, strings.HasPrefix cho flag parsing.
)

// TODO-[2]: Tạo named type cho TemperatureUnit
// SENIOR ASKS: Tại sao dùng named type thay vì string thường?
// HINT: Type safety — không thể truyền sai unit vào hàm.

// type TemperatureUnit int
// const (
//     Celsius TemperatureUnit = iota
//     Fahrenheit
//     Kelvin
// )

// TODO-[3]: Parse unit từ string
// SENIOR ASKS: Làm sao parse "celsius", "C", "c" đều thành Celsius?
// HINT: Normalize: toLower + switch hoặc map lookup.

// func parseUnit(s string) (TemperatureUnit, error) { ... }

// TODO-[4]: Convert temperature
// SENIOR ASKS: Công thức C→F? F→C? C→K? K→C? Viết từng hàm riêng hay 1 hàm switch?
// HINT: 1 hàm switch gọn hơn. Nhưng mỗi conversion riêng thì test dễ hơn. Trade-off?

// func convert(value float64, from, to TemperatureUnit) (float64, error) { ... }

// TODO-[5]: Format output
// SENIOR ASKS: In ra `100°C = 212.00°F` — dùng Printf thế nào?
// HINT: Cần symbol cho mỗi unit. Thêm method `Symbol()` cho TemperatureUnit?

// func formatResult(value float64, from, to TemperatureUnit, result float64) string { ... }

// TODO-[6]: Main flow
// SENIOR ASKS: Cấu trúc main() nên như thế nào? Có nên gộp tất cả vào main không?
// HINT: KHÔNG. main() chỉ nên: parse args → gọi logic → handle error → exit.

func main() {
	// TODO: Kiểm tra số lượng arguments
	// TODO: Parse flags hoặc positional args
	// TODO: Parse temperature value
	// TODO: Parse unit from/to
	// TODO: Convert
	// TODO: Print result hoặc error
	// TODO: Exit với code đúng
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. `os.Args[0]` là gì? Nếu tôi chạy `go run . convert 100`, os.Args có mấy phần tử?
2. `flag` package vs parse `os.Args` thủ công — khi nào dùng cái nào?
3. Tại sao `main()` nên ngắn? Hàm main dài 100 dòng có vấn đề gì?
4. `os.Exit(1)` khác gì với `return` trong main?
5. Nếu tôi muốn thêm command `convert area` (m2 → ft2) — refactor thế nào?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Import đúng packages
- [ ] TODO-[2] hoàn thành: TemperatureUnit named type + iota constants
- [ ] TODO-[3] hoàn thành: Parse unit từ string, case-insensitive
- [ ] TODO-[4] hoàn thành: Convert với đầy đủ công thức
- [ ] TODO-[5] hoàn thành: Format output đẹp, 2 chữ số thập phân
- [ ] TODO-[6] hoàn thành: Main flow rõ ràng, không gộp logic

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `convert 100 c f` → `100°C = 212.00°F`
- [ ] Test case: `convert 32 f c` → `32°F = 0.00°C`
- [ ] Test case: `convert 0 c k` → `0°C = 273.15K`
- [ ] Test case: `convert abc c f` → error message rõ ràng
- [ ] Test case: `convert 100 c x` → "unknown unit x"
- [ ] Test case: Không đủ args → help text
- [ ] Test case: `--help` → usage instructions
- [ ] Boundary case: `-40 c f` → `-40°C = -40.00°F` (intersection point)
- [ ] Boundary case: `convert` với absolute zero — có reject không?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu cần thêm 5 đơn vị nữa (rankine, delisle...) — code của bạn dễ extend không?
2. `os.Args` parsing và convert logic đang ở cùng package. Khi nào nên tách ra package riêng?
3. Tool này chạy trên Windows với `;` separator có vấn đề gì không?
4. Nếu convert cần precision cao (financial calculation) — `float64` có đủ không?

---

## Week 2: Control Flow & Error Handling

> **Mục tiêu tuần:** Code chạy ổn định trong điều kiện xấu. Nắm chắc if/switch, loops, error handling philosophy của Go.

---

## Topic 7: If/Switch Decision Design

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần classify đơn hàng theo giá trị: low risk (< $100), medium risk ($100-$1000), high risk (> $1000). Logic phải đọc được trong 1 màn hình — không được nested if quá sâu."*
>
> Context: Bạn đang viết risk assessment module cho fintech app. Code sẽ được audit — phải clean, readable.

### Acceptance Criteria
- [ ] Dùng guard clause để giảm nesting
- [ ] Switch nhiều case rõ ràng
- [ ] Không có nested if quá 2 cấp
- [ ] Logic đọc được trong 1 màn hình (không scroll)
- [ ] Biết 4 dạng switch: expression, no-expression, multiple values, type switch

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Classify theo giá trị — nghe đơn giản. Nhưng junior thường viết nested if 5 cấp.
> 'Arrowhead antipattern' — nhìn code như mũi tên chỉ sang phải. Không được.
>
> "Vấn đề cốt lõi ở đây là: guard clause. Return sớm, giảm nesting.
> Mỗi guard check 1 điều kiện: nil check → range check → business logic.
> Code chính ở cuối function, không lồng trong if nào.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Viết version xấu (nested if) → 2. Refactor bằng guard clause → 3. Dùng switch
>
> "Hồi tôi review code, tôi thường đếm indentation levels. Quá 2 cấp = refactor.
> Không phải vì 'đẹp' — mà vì brain của chúng ta chỉ hold được ~3-4 levels nesting.
> Hơn nữa là quên context, dẫn đến bug.
```

#### TODO Comments (Code Skeleton)
```go
package main

import "fmt"

// Order represents a purchase order
type Order struct {
	ID    string
	Total float64
	Items []string
}

// TODO-[1]: Viết classifyRisk — VERSION XẤU (nested if)
// SENIOR ASKS: Viết version với nested if 4 cấp. Nhìn code có dễ đọc không?
// HINT: Không — nhưng hãy viết để thấy sự khủng khiếp, rồi refactor.

// func classifyRiskBad(order *Order) string { ... }

// TODO-[2]: Refactor bằng guard clause
// SENIOR ASKS: Mỗi guard clause return sớm. Code chính còn lại gì?
// HINT: Sau khi tất cả guards return, code chính ở top level — 0 indentation.

// func classifyRiskGood(order *Order) string { ... }

// TODO-[3]: Dùng switch không expression cho range check
// SENIOR ASKS: Switch không có expression (= switch true) — khi nào dùng?
// HINT: Khi có nhiều range check. Đọc gọn hơn if-else chain.

// func classifyRiskSwitch(total float64) string { ... }

// TODO-[4]: Switch với multiple values
// SENIOR ASKS: `case "pending", "processing":` — chạy như thế nào?
// HINT: OR logic — match bất kỳ value nào trong list.

// TODO-[5]: Type switch (bonus)
// SENIOR ASKS: Type switch dùng khi nào? Gặp trong thực tế không?
// HINT: JSON unmarshaling, interface{} handling. Sẽ gặp nhiều ở Phase 3-4.

func main() {
	// TODO: Tạo vài Order test, gọi classify, in kết quả
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Guard clause tối đa nên có mấy? Nếu function có 10 guards, có vấn đề không?
2. Go switch không tự fallthrough — tại sao? C/Java thì khác — trade-off là gì?
3. Khi nào chọn `if-else` thay vì `switch`? Khi nào ngược lại?
4. `switch` trong Go có thể match trên type không? Khác gì với type assertion?
5. Nếu requirement thêm: "medium risk chia thành medium-low và medium-high" — code thay đổi thế nào?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Viết version xấu để thấy vấn đề
- [ ] TODO-[2] hoàn thành: Refactor bằng guard clause, nesting ≤ 1
- [ ] TODO-[3] hoàn thành: Dùng switch no-expression cho range
- [ ] TODO-[4] hoàn thành: Switch multiple values
- [ ] TODO-[5] hoàn thành: Biết type switch (không cần implement sâu)

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `nil` order → guard catches it
- [ ] Test case: Total = -100 → invalid input
- [ ] Test case: Total = 50 → low risk
- [ ] Test case: Total = 500 → medium risk
- [ ] Test case: Total = 5000 → high risk
- [ ] Test case: Total = 100 (boundary) → medium risk
- [ ] Test case: Total = 1000 (boundary) → medium risk
- [ ] Test case: Empty items → still classifies based on total

### Retrospective: Sau khi xong, hãy tự hỏi
1. Rule "indentation ≤ 2 levels" áp dụng cho mọi ngôn ngữ hay chỉ Go?
2. Guard clause có phải lúc nào cũng return error? Có khi return value bình thường không?
3. Nếu logic phức tạp hơn (machine learning classification), guard clause còn đủ không?

---

## Topic 8: Loops & range

### User Story
> Khách hàng (Product Owner) nói: *"Tôi có danh sách đơn hàng. Cần: tính tổng giá trị, filter đơn hàng > $1000, và tìm đơn hàng đầu tiên có status 'shipped'."*
>
> Context: Bạn đang viết reporting module. Cần iterate qua slice — đúng cách, không off-by-one.

### Acceptance Criteria
- [ ] 3 cách viết for: classic (C-style), condition-only (while-like), range
- [ ] Dùng `range` đúng với slice và map
- [ ] Không có off-by-one bug
- [ ] Biết `break` và `continue` — khi nào dùng
- [ ] Hiểu: `for` là vòng lặp DUY NHẤT trong Go — không có while, do-while, foreach
- [ ] Biết `range` copy value — trap khi dùng với pointer

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Iterate qua collection — nghe đơn giản. Nhưng Go chỉ có `for` — không có while,
> foreach, forEach như JS. range là multi-purpose: slice, map, string, channel.
> Nhưng range có trap: nó copy value, không phải reference.
>
> "Vấn đề cốt lõi ở đây là: chọn đúng loop style cho từng task.
> Classic for khi cần index control. Range khi chỉ cần iterate.
> Condition-only khi không biết trước số lần lặp.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Sum (range) → 2. Filter (classic for or range) → 3. Find first (range + break)
>
> "Hồi tôi gặp 1 bug production: range qua slice struct, lấy địa chỉ của range variable.
> Tất cả pointers trỏ cùng 1 memory address! Vì range copy value vào cùng 1 variable.
> Tôi mất 3 tiếng debug. Lesson: `for i := range items { ptr := &items[i] }` —
> lấy địa chỉ của slice element, không phải range variable.
```

#### TODO Comments (Code Skeleton)
```go
package main

import "fmt"

type Order struct {
	ID     string
	Total  float64
	Status string
}

// TODO-[1]: Tính tổng giá trị đơn hàng bằng range
// SENIOR ASKS: range trả về gì với slice? (index, value). Nếu chỉ cần value, làm sao bỏ index?
// HINT: `for _, order := range orders` — blank identifier `_` để bỏ qua index.

// func sumOrders(orders []Order) float64 { ... }

// TODO-[2]: Filter đơn hàng > $1000
// SENIOR ASKS: Trả về slice mới hay modify slice cũ? Khi nào chọn cái nào?
// HINT: Trả về slice mới (pure function) — predictable, testable.

// func filterHighValue(orders []Order) []Order { ... }

// TODO-[3]: Tìm đơn hàng đầu tiên có status 'shipped'
// SENIOR ASKS: Trả về `Order` hay `*Order`? Nếu không tìm thấy thì sao?
// HINT: Trả về `(Order, bool)` — found flag. Hoặc `(*Order, error)`. Cân nhắc?

// func findFirstShipped(orders []Order) (Order, bool) { ... }

// TODO-[4]: 3 dạng for
// SENIOR ASKS: Viết lại TODO-[1] bằng classic for. Viết countdown bằng condition-only.
// HINT: `for i := 0; i < len(orders); i++` và `for n > 0 { n-- }`

// TODO-[5]: range trap — pointer bug
// SENIOR ASKS: Đoạn code sau có bug gì? Fix thế nào?
// HINT: `ptr := &order` lấy địa chỉ của RANGE VARIABLE — mỗi iteration cùng address.

// var pointers []*Order
// for _, order := range orders {
//     pointers = append(pointers, &order)  // BUG!
// }

func main() {
	orders := []Order{
		{ID: "A1", Total: 50, Status: "pending"},
		{ID: "A2", Total: 1500, Status: "shipped"},
		{ID: "A3", Total: 750, Status: "shipped"},
	}
	// TODO: Gọi các function trên, in kết quả
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. `for i := range slice` — `i` có phải luôn bắt đầu từ 0 không? Có thể bắt đầu từ giữa slice không?
2. `range` qua map — thứ tự có guaranteed không? Code phụ thuộc vào thứ tự map có ổn không?
3. `range` qua string — trả về byte hay rune? `for i, c := range "Hello, 世界"` — `c` type gì?
4. Tại sao Go không có `while` keyword? `for condition` đã đủ chưa?
5. `break` trong nested loop — break khỏi loop nào? Có label trong Go không?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Sum bằng range, dùng blank identifier
- [ ] TODO-[2] hoàn thành: Filter trả về slice mới
- [ ] TODO-[3] hoàn thành: Find first với found flag
- [ ] TODO-[4] hoàn thành: 3 dạng for đều viết được
- [ ] TODO-[5] hoàn thành: Hiểu range pointer trap và cách fix

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: Empty slice → sum = 0, filter = empty, find = not found
- [ ] Test case: 1 element → sum = element value
- [ ] Test case: All elements match filter → return copy of original
- [ ] Test case: No element matches filter → return empty slice (not nil?)
- [ ] Test case: Find first — trả về đúng element đầu tiên match
- [ ] Boundary case: Large slice — performance của range vs classic for?
- [ ] Trap case: Range pointer bug — verify addresses are different after fix

### Retrospective: Sau khi xong, hãy tự hỏi
1. `range` copy value — performance impact với struct lớn? Có cách nào range by reference không?
2. Filter bằng loop vs functional-style (Map/Filter/Reduce) — Go prefer cái nào? Tại sao?
3. `for i := 0; i < len(s); i++` vs `for i := range s` — cái nào nhanh hơn? Tại sao?

---

## Topic 9: Error Handling

### User Story
> Khách hàng (Product Owner) nói: *"Khi parse file CSV mà có dòng lỗi, báo lỗi rõ ràng cho tôi biết dòng nào lỗi, vì sao lỗi. Không được dừng chương trình — skip dòng lỗi, continue xử lý các dòng còn lại."*
>
> Context: Bạn đang viết CSV import module. File có thể hàng nghìn dòng — không thể dừng vì 1 dòng lỗi.

### Acceptance Criteria
- [ ] Error return: `(result, error)` pattern
- [ ] Error message rõ ràng: chứa context (dòng nào, field nào)
- [ ] Wrap error: `fmt.Errorf("...: %w", err)` để preserve error chain
- [ ] Không `panic` trong normal flow — chỉ panic cho unrecoverable
- [ ] Dùng `errors.New` cho static message, `fmt.Errorf` cho dynamic
- [ ] Phân biệt: error vs log vs return code

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Error handling là một trong những điểm khác biệt lớn nhất của Go.
> Không try/catch. Không exceptions. Mọi error là value.
> Nghe verbose, nhưng explicit error paths giúp code predictable.
>
> "Vấn đề cốt lõi ở đây là: error wrapping.
> Khi parse dòng 42, field 'price' fail → error nên là:
> 'line 42: invalid price "abc": strconv.ParseFloat: parsing \"abc\"'. 
> Từng layer wrap thêm context.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Parse line → 2. Wrap error với context → 3. Aggregate errors → 4. Continue processing
>
> "Hồi tôi debug 1 production issue, tôi thấy log: 'parse error'. Thế parse cái gì?
> Dòng nào? File nào? Không biết. Mất 2 giờ tìm. Từ đó tôi luôn wrap error
> với context: fmt.Errorf("file %s line %d: %w", filename, lineNum, err).
> Đây gọi là 'error annotation' — mỗi layer thêm thông tin.
```

#### TODO Comments (Code Skeleton)
```go
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseResult chứa kết quả parse và các lỗi gặp phải
type ParseResult struct {
	ValidLines []LineItem
	Errors     []error
}

type LineItem struct {
	LineNo int
	Name   string
	Price  float64
}

// TODO-[1]: Parse 1 dòng CSV
// SENIOR ASKS: Split string bằng gì? strings.Split? csv package?
// HINT: Đơn giản thì strings.Split. Nhưng production dùng encoding/csv — nó handle quoted fields.

// func parseLine(line string, lineNo int) (LineItem, error) { ... }

// TODO-[2]: Error wrapping với context
// SENIOR ASKS: `fmt.Errorf("line %d: %w", lineNo, err)` — `%w` làm gì?
// HINT: `%w` wrap error. Sau này dùng errors.Is/errors.Unwrap để check.

// TODO-[3]: Parse file với continue-on-error
// SENIOR ASKS: Lỗi từng dòng collect vào đâu? Return ngay hay aggregate?
// HINT: Aggregate — trả về cả valid lines và errors. Caller quyết định.

// func parseFile(lines []string) ParseResult { ... }

// TODO-[4]: Phân biệt errors.New vs fmt.Errorf
// SENIOR ASKS: Khi nào errors.New? Khi nào fmt.Errorf?
// HINT: errors.New cho static string (định nghĩa sentinel errors).
// fmt.Errorf khi cần embed runtime values vào message.

// TODO-[5]: Sentinel error
// SENIOR ASKS: Sentinel error là gì? Tại sao cần?
// HINT: Package-level var error = errors.New("..."). Caller dùng errors.Is để check.

// var ErrInvalidFormat = errors.New("invalid CSV format")

func main() {
	csvLines := []string{
		"Apple,1.50",
		"Banana,invalid",
		"Orange,2.00",
		"", // empty line
	}
	// TODO: Parse, handle errors, print results
}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. `fmt.Errorf("%w", err)` vs `fmt.Errorf("%v", err)` — khác gì? Khi nào dùng `%w`?
2. `errors.Is(err, ErrNotFound)` dùng khi nào? So sánh với `err == ErrNotFound`?
3. Tại sao Go không có `try/catch`? Nếu bạn thấy `(value, error)` verbose — cách nào giảm?
4. `panic` khác `error` như thế nào? Khi nào được phép panic?
5. Nếu 1 function gọi 5 hàm, mỗi hàm trả error — bạn check error sau mỗi lần gọi. Có pattern nào gọn hơn không?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Parse line thành struct
- [ ] TODO-[2] hoàn thành: Wrap error với context (line number, field)
- [ ] TODO-[3] hoàn thành: Continue-on-error, aggregate errors
- [ ] TODO-[4] hoàn thành: Phân biệt errors.New vs fmt.Errorf
- [ ] TODO-[5] hoàn thành: Define sentinel error

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: Valid line → parsed correctly
- [ ] Test case: Non-numeric price → error with line number context
- [ ] Test case: Wrong field count → error
- [ ] Test case: Empty line → skip or error?
- [ ] Test case: All lines invalid → no valid lines, all errors collected
- [ ] Test case: errors.Is(wrappedErr, ErrInvalidFormat) → true
- [ ] Boundary case: Very long line — memory usage?

### Retrospective: Sau khi xong, hãy tự hỏi
1. `fmt.Errorf("...: %w", err)` — `%w` dùng reflection không? Performance cost?
2. Nếu requirement đổi: "dừng ngay khi gặp lỗi" — code thay đổi bao nhiêu %?
3. Error message nên bằng tiếng Anh hay tiếng Việt? Convention trong team?
4. `ParseResult` struct chứa `[]error` — có nên dùng custom error type thay vì `[]error` không?

---

## Topic 10-12: Spec-first Development

### User Story
> Khách hàng (Product Owner) nói: *"Trước khi code, viết spec cho tôi. Tôi phải approve spec trước. Sau đó code phải match spec 100%. Nếu không match, không pass review."*
>
> Context: Bạn đang viết feature `stats` — tính toán thống kê trên slice số (min, max, avg, sum). Requirement rõ ràng — phải có spec trước.

### Acceptance Criteria
- [ ] Có file `spec.md` mô tả: inputs, outputs, errors, examples
- [ ] Code match spec: function signatures, error cases, output format
- [ ] Test matrix ≥ 12 cases: happy path + error path + boundary
- [ ] README có cách chạy và ví dụ
- [ ] Spec được review và approved trước khi code (simulated)

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Spec-first là habit của senior engineer. Không phải vì process — mà vì
> spec giúp bạn nghĩ rõ requirement trước khi code. Code là phần dễ.
> Hiểu đúng requirement mới khó.
>
> "Vấn đề cốt lõi ở đây là: định nghĩa contract rõ ràng.
> Input gì? Output gì? Error khi nào? Edge cases nào?
> Khi bạn viết spec, bạn thấy những chỗ requirement mơ hồ — clarify trước khi code.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Viết spec.md → 2. Review spec (tự review) → 3. Implement → 4. Test matrix → 5. Verify
>
> "Hồi tôi viết 1 feature mà không viết spec trước. Code xong, demo cho PO.
> PO nói: 'Tôi muốn khi empty slice thì trả lỗi, không phải return 0.'
> Tôi phải refactor 50% code. Nếu tôi viết spec trước, clarify chỗ này ngay
> từ đầu — đỡ mất 2 giờ. Lesson: 15 phút viết spec → tiết kiệm hours debugging.
```

#### TODO Comments (Code Skeleton)
```go
// ============================================
// FILE: spec.md (VIẾT TRƯỚC KHI CODE)
// ============================================

// # Stats Specification
//
// ## Overview
// Package `stats` cung cấp functions tính toán thống kê cơ bản trên slice số.
//
// ## Functions
//
// ### Sum(numbers []float64) (float64, error)
// - Input: slice of float64
// - Output: tổng các số
// - Error: nếu slice empty → `ErrEmptySlice`
// - Example: Sum([]float64{1, 2, 3}) → 6, nil
//
// ### Average(numbers []float64) (float64, error)
// - Input: slice of float64
// - Output: trung bình cộng
// - Error: nếu slice empty → `ErrEmptySlice`
// - Example: Average([]float64{1, 2, 3}) → 2, nil
//
// ### MinMax(numbers []float64) (min, max float64, err error)
// - Input: slice of float64
// - Output: min và max
// - Error: nếu slice empty → `ErrEmptySlice`
// - Example: MinMax([]float64{3, 1, 4}) → 1, 4, nil
//
// ## Sentinel Errors
// - `ErrEmptySlice = errors.New("slice is empty")`
//
// ## Non-functional Requirements
// - Không dùng global state
// - Pure functions (input → output, không side effect)
// - Time complexity: O(n)
// - Không panic
//
// ## Test Matrix
// | Case | Input | Expected | Notes |
// |---|---|---|---|
// | T1 | [1,2,3] | Sum=6 | Happy path |
// | T2 | [] | ErrEmptySlice | Empty input |
// | T3 | [5] | Sum=5, Min=5, Max=5 | Single element |
// | T4 | [-1,-2,-3] | Sum=-6 | Negative numbers |
// | T5 | [0,0,0] | Sum=0 | All zeros |
// | T6 | [1.5, 2.5] | Sum=4.0 | Floats |
// | T7 | [1e308, 1e308] | Inf | Overflow |
// | ... | ... | ... | ... |

// ============================================
// FILE: stats.go (CODE SAU KHI SPEC APPROVED)
// ============================================

package stats

import "errors"

// TODO-[1]: Định nghĩa sentinel error
// SENIOR ASKS: Tại sao sentinel error ở package level?
// HINT: Caller dùng errors.Is để check. Phải exported để caller access.

// var ErrEmptySlice = errors.New("slice is empty")

// TODO-[2]: Implement Sum
// SENIOR ASKS: Time complexity? Có cần check nil slice không?
// HINT: len(nilSlice) == 0, nên check len là đủ. Không cần nil check riêng.

// func Sum(numbers []float64) (float64, error) { ... }

// TODO-[3]: Implement Average
// SENIOR ASKS: float64 division — có precision issue không? Khi nào?
// HINT: Luôn có. Nhưng với stats cơ bản, float64 đủ. Financial thì dùng decimal.

// func Average(numbers []float64) (float64, error) { ... }

// TODO-[4]: Implement MinMax
// SENIOR ASKS: Named returns — dùng hay không? Khi nào?
// HINT: Named returns gọn khi nhiều return values. Nhưng dễ gây confusion. Cân nhắc.

// func MinMax(numbers []float64) (min, max float64, err error) { ... }

// TODO-[5]: Viết test matrix
// SENIOR ASKS: Table-driven test trong Go viết như thế nào?
// HINT: Slice of structs: []struct{name string; input []float64; wantSum float64; wantErr bool}{...}
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Tại sao viết spec trước code giúp catch bug sớm? Đưa 2 ví dụ cụ thể.
2. Spec nên chi tiết đến mức nào? Có nên định nghĩa implementation detail không?
3. "Code match spec 100%" — nếu bạn phát hiện spec sai trong lúc code, làm gì?
4. Test matrix ≥ 12 cases — có quá nhiều không? Khi nào ít hơn được?
5. Pure function là gì? Tại sao Go encourage pure functions? Nhược điểm?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Sentinel error defined và exported
- [ ] TODO-[2] hoàn thành: Sum implement + test pass
- [ ] TODO-[3] hoàn thành: Average implement + test pass
- [ ] TODO-[4] hoàn thành: MinMax implement + test pass
- [ ] TODO-[5] hoàn thành: Test matrix ≥ 12 cases all pass
- [ ] spec.md tồn tại và được review

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: T1 — Happy path [1,2,3] → Sum=6, Avg=2, Min=1, Max=3
- [ ] Test case: T2 — Empty slice → ErrEmptySlice
- [ ] Test case: T3 — Single element → all values = element
- [ ] Test case: T4 — Negative numbers → correct negative sum
- [ ] Test case: T5 — All zeros → Sum=0
- [ ] Test case: T6 — Float precision → 1.5 + 2.5 = 4.0
- [ ] Test case: T7 — Overflow → +Inf (không panic)
- [ ] Test case: T8 — Large slice → O(n) performance
- [ ] Test case: T9 — Nil slice → treated as empty
- [ ] Test case: T10 — Very large numbers → no precision loss (relative)
- [ ] Test case: T11 — Mixed positive/negative → correct result
- [ ] Test case: T12 — Infinity in input → propagate correctly

### Retrospective: Sau khi xong, hãy tự hỏi
1. Nếu requirement thêm median và mode — spec thay đổi bao nhiêu %?
2. Table-driven test trong Go — ưu/nhược điểm so với test framework như JUnit?
3. `float64` cho financial calculation — tại sao không nên? Alternative?
4. Spec.md nên lưu ở đâu? Cùng repo? Hay wiki? Lợi ích của mỗi cách?

---

## Week 3: Foundation CLI Toolkit (Mini-project)

> **Mục tiêu tuần:** Mô phỏng mini lifecycle đi làm. Ship 1 project có: 3 commands, README, test, không panic.

---

## Topic 13: CLI Toolkit — convert + stats + inspect

### User Story
> Khách hàng (Product Owner) nói: *"Tôi cần bộ tool CLI hoàn chỉnh gồm 3 commands: `convert` chuyển đổI đơn vị, `stats` tính thống kê trên dãy số, `inspect` phân tích file text (số dòng, số từ). Phải có README, hướng dẫn chạy, và không được crash dù input có bẩn."*
>
> Context: Đây là deliverable cuối Phase 1 — Foundation CLI Toolkit. Bạn sẽ áp dụng tất cả kiến thức 2 tuần trước vào project thực.

### Acceptance Criteria
- [ ] 3 commands chạy được: `convert`, `stats`, `inspect`
- [ ] `convert`: giống Topic 6, nhưng robust hơn
- [ ] `stats`: tính min/max/avg/sum từ CLI args hoặc stdin
- [ ] `inspect`: đọc file, đếm lines, words, characters
- [ ] README với cách cài đặt, chạy, ví dụ
- [ ] `go test ./...` pass
- [ ] Không panic với input bất kỳ
- [ ] Help text cho mỗi command
- [ ] Exit code 0 cho success, 1 cho error

### Senior Guide

#### Senior Thought-Process
```markdown
**Senior nghĩ gì khi nhận requirement này:**
> "Đây là mini-project đầu tiên — tôi sẽ approach như 1 project thực.
> Không code ngay. Phân rã trước.
>
> "Vấn đề cốt lõi ở đây là: cấu trúc project.
> Không gộp tất cả vào 1 file main.go. Tách commands ra riêng.
> Package structure: main → commands → internal logic.
>
> "Tôi sẽ phân rã thành các bước:
> 1. Thiết kế architecture → 2. Implement từng command → 3. Error handling
> → 4. Test → 5. README
>
> "Architecture note của tôi cho project này:
> ```
> go-cli-toolkit/
> ├── main.go              # entry point, command dispatch
> ├── go.mod
> ├── README.md
> ├── commands/
> │   ├── convert.go       # convert command
> │   ├── stats.go         # stats command
> │   └── inspect.go       # inspect command
> └── internal/
>     ├── convert/
>     │   └── convert.go   # conversion logic (testable)
>     ├── stats/
>     │   └── stats.go     # stats logic (từ Topic 10-12)
>     └── inspect/
>         └── inspect.go   # file analysis logic
> ```
>
> "Pattern: commands/ chứa CLI parsing (os.Args). internal/ chứa pure logic.
> Logic trong internal không phụ thuộc CLI — test được mà không cần mock CLI args.
> Đây là pattern sẽ dùng lại ở mọi project.
>
> "Hồi tôi làm project tương tự, tôi gộp tất cả vào main.go. 300 dòng.
> Không test được. Không reuse được. Senior review bảo: 'commands phải tách ra.'
> Tôi refactor mất 1 ngày. Lesson: phân rã từ đầu, không đợi review mới refactor.
```

#### TODO Comments (Code Skeleton)

```go
// ============================================
// FILE: main.go
// ============================================
package main

import (
	"fmt"
	"os"
	// TODO-[1]: Import packages commands
	// SENIOR ASKS: Tại sao không gộp tất cả commands vào main?
	// HINT: Separation of concerns. main chỉ dispatch. Mỗi command self-contained.
)

// TODO-[2]: Command dispatch
// SENIOR ASKS: `os.Args[1]` là command name. Làm sao dispatch đúng?
// HINT: Switch hoặc map[string]func. Map gọn hơn nếu nhiều commands.

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Dispatch based on os.Args[1]
	// switch os.Args[1] {
	// case "convert": commands.RunConvert(os.Args[2:])
	// case "stats":   commands.RunStats(os.Args[2:])
	// case "inspect": commands.RunInspect(os.Args[2:])
	// default: printUsage(); os.Exit(1)
	// }
}

func printUsage() {
	// TODO: In help text đẹp, liệt kê tất cả commands
}
```

```go
// ============================================
// FILE: commands/convert.go
// ============================================
package commands

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO-[3]: RunConvert — entry point cho command convert
// SENIOR ASKS: Dùng `flag` package hay parse os.Args thủ công?
// HINT: `flag` package tốt hơn: tự động --help, validate, type conversion.

// func RunConvert(args []string) {
//     TODO: Define flags: --from, --to
//     TODO: Parse args
//     TODO: Validate
//     TODO: Call internal/convert.Convert
//     TODO: Print result or error
//     TODO: Exit with correct code
// }

// ============================================
// FILE: internal/convert/convert.go
// ============================================
package convert

import "fmt"

// TemperatureUnit type + constants
// ... (từ Topic 6)

// Convert function — pure logic, no CLI dependency
// TODO-[4]: Tại sao function này không dùng os.Exit? Không in ra stdout?
// SENIOR ASKS: Làm sao test function này mà không capture stdout?
// HINT: Trả về (float64, error). Caller (commands/) quyết định in/exit.

// func Convert(value float64, from, to TemperatureUnit) (float64, error) { ... }
```

```go
// ============================================
// FILE: commands/stats.go
// ============================================
package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO-[5]: RunStats — tính stats từ args hoặc stdin
// SENIOR ASKS: Làm sao đọc từ stdin? `bufio.Scanner` — dùng thế nào?
// HINT: `scanner := bufio.NewScanner(os.Stdin); scanner.Scan()` đọc từng dòng.

// func RunStats(args []string) { ... }

// ============================================
// FILE: internal/stats/stats.go
// ============================================
package stats

// ... (từ Topic 10-12, đã implement)
// Sum, Average, MinMax — pure functions
```

```go
// ============================================
// FILE: commands/inspect.go
// ============================================
package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TODO-[6]: RunInspect — đọc file, đếm lines/words/chars
// SENIOR ASKS: Đếm "words" như thế nào? strings.Fields — hàm này làm gì?
// HINT: strings.Fields split bằng whitespace (space, tab, newline).

// func RunInspect(args []string) { ... }

// ============================================
// FILE: internal/inspect/inspect.go
// ============================================
package inspect

// TODO-[7]: InspectResult struct + Inspect function
// SENIOR ASKS: Tại sao cần struct thay vì trả về 3 giá trị riêng lẻ?
// HINT: Extensible — sau này thêm "line count" không đổi signature.

// type Result struct {
//     Lines int
//     Words int
//     Chars int
// }
// func Inspect(r io.Reader) (Result, error) { ... }
```

#### Socratic Questions
```markdown
**Câu hỏi để bạn tự suy nghĩ:**
1. Tại sao `internal/` package không thể import từ bên ngoài module?
2. `commands/` vs `internal/` — ranh giới ở đâu? Khi nào code ở commands, khi nào ở internal?
3. `flag` package vs `os.Args` manual parsing — khi nào nên dùng thư viện thứ 3 như `cobra`?
4. `bufio.Scanner` có giới hạn dòng không? Nếu file có dòng > 64K — chuyện gì xảy ra?
5. Nếu requirement thêm command thứ 4 — code thay đổi bao nhiêu %? Điều này nói lên điều gì về design?
```

### Output Checklist: Làm sao biết mình xong?
- [ ] TODO-[1] hoàn thành: Import structure rõ ràng
- [ ] TODO-[2] hoàn thành: Command dispatch hoạt động
- [ ] TODO-[3] hoàn thành: RunConvert với flag package
- [ ] TODO-[4] hoàn thành: Convert pure function ở internal/
- [ ] TODO-[5] hoàn thành: RunStats đọc từ args và stdin
- [ ] TODO-[6] hoàn thành: RunInspect đọc file
- [ ] TODO-[7] hoàn thành: InspectResult struct + Inspect function
- [ ] README.md hoàn chỉnh với examples
- [ ] `go test ./...` xanh

### Test Checklist: Những gì bạn nên tự kiểm tra
- [ ] Test case: `tool convert --from=celsius --to=fahrenheit 100` → 212.00
- [ ] Test case: `tool convert` → help text, exit 1
- [ ] Test case: `tool stats 1 2 3 4 5` → sum=15, avg=3, min=1, max=5
- [ ] Test case: `tool stats` (no args) → read from stdin
- [ ] Test case: `tool inspect file.txt` → lines, words, chars
- [ ] Test case: `tool inspect nonexistent.txt` → error, exit 1
- [ ] Test case: `tool unknown-cmd` → unknown command error
- [ ] Test case: `tool --help` → usage for all commands
- [ ] Test case: `tool convert abc c f` → invalid number error
- [ ] Test case: `tool convert 100 c unknown` → unknown unit error
- [ ] Test case: Large file inspect — performance acceptable?
- [ ] Test case: Unicode file inspect — word count đúng?

### Retrospective: Sau khi xong, hãy tự hỏi
1. Pattern `commands/` + `internal/` — bạn sẽ dùng lại ở project nào? Có trường hợp không phù hợp?
2. Nếu cần share logic giữa commands (ví dụ: common argument parsing) — extract ở đâu?
3. `go test ./...` chỉ chạy unit test. Integration test cho CLI tool viết như thế nào?
4. Binary size của tool này bao nhiêu? `go build -ldflags="-s -w"` giảm bao nhiêu?
5. Sau 3 tuần, bạn thấy Go khác Dart/JS nhất ở điểm nào? Điều gì bạn thích nhất? Điều gì khó chịu nhất?

---

## Appendix: Go Toolchain Cheatsheet

```
# Module
$ go mod init <module-path>    # Khởi tạo module mới
$ go mod tidy                  # Dọn dẹp dependencies
$ go mod download              # Tải dependencies

# Build & Run
$ go run <file.go>             # Compile + chạy (temp binary)
$ go run .                     # Chạy package hiện tại
$ go build                     # Build binary (tên từ module)
$ go build -o <name>           # Build với tên custom
$ go install                   # Build + install vào $GOPATH/bin

# Format & Quality
$ go fmt ./...                 # Format toàn bộ project
$ go vet ./...                 # Static analysis
$ go lint ./...                # (cần cài golint)

# Test
$ go test ./...                # Chạy tất cả tests
$ go test -v ./...             # Verbose output
$ go test -race ./...          # Race condition detection
$ go test -cover ./...         # Coverage report

# Documentation
$ go doc <package>             # Xem documentation
$ go doc <package.Symbol>      # Xem doc của symbol
```

## Appendix: Zero Values Reference

| Type | Zero Value |
|---|---|
| bool | false |
| int, int8, int16, int32, int64 | 0 |
| uint, uint8, uint16, uint32, uint64 | 0 |
| float32, float64 | 0.0 |
| string | "" (empty string) |
| pointer | nil |
| slice | nil |
| map | nil (không thể write trước khi make) |
| channel | nil |
| interface | nil |
| struct | các field đều zero value |

## Appendix: Format Verbs Quick Reference

| Verb | Type | Output |
|---|---|---|
| %s | string | raw string |
| %d | integer | decimal |
| %f | float | decimal notation |
| %.2f | float | 2 decimal places |
| %e | float | scientific notation |
| %v | any | default format |
| %+v | struct | default + field names |
| %#v | any | Go syntax representation |
| %T | any | type of value |
| %t | bool | true/false |
| %10s | string | right-aligned, width 10 |
| %-10s | string | left-aligned, width 10 |
| %010d | int | zero-padded, width 10 |
| %q | string | double-quoted |
| %x | string/int | hex |
| %w | error | wrap error (fmt.Errorf only) |

---

> **Phase 1 Complete.** Bạn đã có nền tảng vững: toolchain, variables, types, constants, control flow, error handling, và 1 CLI toolkit chạy được. Sẵn sàng cho Phase 2: Concurrency & Memory Model.
>
> **Next Phase Preview:** Goroutines, channels, select, sync primitives — và tại sao "Share Memory By Communicating" là philosophy cốt lõi của Go concurrency.
