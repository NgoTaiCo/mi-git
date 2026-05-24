# Phase 2: Goroutines, Channels & Sync Primitives (Tuan 4-5)

> "Trong Go, concurrency khong phai la thu ban hoc. No la thu ban song."  
> — Mot senior engineer nao do, sau lan debug goroutine leak luc 3 gio sang

---

## Tong quan Phase 2

**Muc tieu:** Tu duy CSP (Communicating Sequential Processes), khong bao gio de goroutine leak, chon dung channel hay mutex trong tung tinh huong.

**Du an trong phase nay:** Concurrent Log Analyzer — phan tich log file 10GB song song, khong load het vao RAM.

**Nguon chinh:**
- `go-roadmap-minigit.md` — phase 02-concurrency
- Go.dev blog: "Share Memory By Communicating"
- The Go Programming Language, Chapters 8-9

**Quy tac vang cua phase nay:**
1. `go test -race` phai pass truoc khi commit.
2. Khong dong channel tu receiver — chi sender dong.
3. Context phai duoc truyen xuong, khong bao gio tao `context.Background()` sau main.
4. Moi goroutine phai co co che dung — hoac la `WaitGroup`, hoac la `context.Done()`.

---

## Topic 02.1: Goroutines

### User Story

> Khach hang (Product Owner) noi: "Toi can download 100 file anh tu URL. Tung file mot thi cham qua, can tai song song. Ma khong duoc de chuong trinh bi treo hay loi nhe!"

**Context:** Ban dang viet mot CLI tool de tai batch anh tu cloud storage ve local. Moi file ~1-5MB. Neu tai tuan tu, 100 file * 2 giay = gan 4 phut — qua cham. Product Owner muon xong trong vai chuc giay.

### Acceptance Criteria

- [ ] Tai song song nhieu file cung luc (co the gioi han so goroutine chay dong thoi)
- [ ] Khong bi crash du co URL loi, timeout, hoac server tra 404
- [ ] Chuong trinh biet khi nao TAT CA file da xong — khong exit som
- [ ] Khong de goroutine treo (leak) khi co loi xay ra
- [ ] Co bao cao ket qua: file nao thanh cong, file nao that bai

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Day la case kinh dien cho goroutine. Toi nho hoi o project pipeline automation,
> team viet mot tool download artifacts tu S3 — ban dau tuan tu, 200 file mat 8 phut.
> Sau khi dung goroutine co gioi han, xuong con 25 giay. Nhung ma lan dau trien khai,
> chung toi leak goroutine vi khong doi chung ket thuc — memory cu tang dan.
>
> Van de cot loi o day la 3 cai:
> 1. Tao goroutine de chay song song — don gian, chi can `go` keyword
> 2. Gioi han concurrency — khong de 100 goroutine chay cung luc lam server nghen
> 3. Bao cao khi nao xong — WaitGroup hoac channel signal
>
> Toi se phan ra thanh cac buoc:
> - Buoc 1: Tao ham download co signature ro rang
> - Buoc 2: Tao goroutine cho moi URL
> - Buoc 3: Dung WaitGroup de doi tat ca xong
> - Buoc 4: Xet xem co can gioi han concurrency khong (se hoc o topic 02.4)
>
> Quan trong nhat: khong bao gio viet `go func()` ma khong co co che doi no ket thuc.
> Do la cach nhanh nhat de tao goroutine leak."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// DownloadResult luu ket qua tai moi file
// SENIOR ASKS: Tai sao nen dung struct thay vi chi tra error?
// HINT: Product Owner muon biet file nao thanh cong, file nao loi.
type DownloadResult struct {
	URL      string
	FilePath string
	Err      error
	// TODO-[1]: Them truong phu hop de luu thong tin ket qua
	// SENIOR ASKS: Neu can biet file tai mat bao lau, ban se them gi?
	// HINT: Kieu du lieu nao luu thoi gian trong Go?
}

// downloadFile tai mot file tu URL va luu vao filepath
// SENIOR ASKS: Tai sao tra ve error thay vì tra ve *DownloadResult?
// HINT: Ham helper nen don gian — tra ve gi de caller tu quyet dinh.
func downloadFile(url, filepath string) error {
	// TODO-[2]: Tao HTTP GET request
	// SENIOR ASKS: Nen dung http.Get hay tao http.Request roi dung client?
	// HINT: http.Get tien nhung co han che — nghĩ ve timeout.

	// TODO-[3]: Kiem tra status code
	// SENIOR ASKS: Status code nao la thanh cong? Neu tra 404 thi sao?
	// HINT: Khong phai chi 200 moi la OK — 201, 204 cung co the.

	// TODO-[4]: Tao file output va ghi body vao
	// SENIOR ASKS: Nen dung io.Copy hay doc het vao memory roi ghi?
	// HINT: File co the rat lon — hay tiet kiem RAM.

	return nil // TODO: thay bang logic thuc su
}

// downloadAll tai song song nhieu file
// SENIOR ASKS: Tai sao parameter la slice URL thay vi variadic?
// HINT: Variadic dep nhung khong phai luc nao cung phu hop.
func downloadAll(urls []string, outputDir string) []DownloadResult {
	// TODO-[5]: Khoi tao WaitGroup
	// SENIOR ASKS: WaitGroup de lam gi trong truong hop nay?
	// HINT: main phai biet khi nao TAT CA goroutine da xong.

	// TODO-[6]: Tao channel hoac slice co bao ve de luu ket qua
	// SENIOR ASKS: Nen dung channel hay slice co mutex de nhan ket qua?
	// HINT: O day, don gian la tot nhat. Nhung nghi ve race condition.

	for _, url := range urls {
		// TODO-[7]: Bat dau goroutine cho moi URL
		// SENIOR ASKS: Co vấn de gi neu go func() ben trong loop ma dung bien loop truc tiep?
		// HINT: Day la bug kinh dien cua Go — closure capture.

		// TODO-[8]: Trong goroutine: goi downloadFile, luu ket qua
		// SENIOR ASKS: Can lam gi voi WaitGroup trong moi goroutine?
		// HINT: Co 3 loi goi WaitGroup — Add, Done, Wait. Dung dung ca 3.
	}

	// TODO-[9]: Doi tat ca goroutine ket thuc
	// SENIOR ASKS: Wait() nen goi o dau? Truoc hay sau khi nhan ket qua?
	// HINT: Thứ tu quan trong — nghi ve deadlock.

	return nil // TODO: tra ve ket qua
}

func main() {
	urls := []string{
		"https://example.com/img1.jpg",
		"https://example.com/img2.jpg",
		"https://example.com/img3.jpg",
		// ... 100 URL
	}

	// TODO-[10]: Gọi downloadAll va in ket qua
	// SENIOR ASKS: Nen xu ly ket qua nhu the nao cho Product Owner de nhin?
	// HINT: Tong hop: thanh cong bao nhieu, that bai bao nhieu, loi cu the la gi.
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Neu ban goi `go func()` 1000 lan thi co van de gi khong?** He thong cua ban co du tai nguyen de chay 1000 goroutine dong thoi khong? Goroutine nhe hon thread nhung khong phai la mien phi.

2. **Khi goroutine chay, neu main() return truoc khi goroutine xong, dieu gi xay ra?** Go co cho goroutine chay den cung khong, hay kill chung? Ban da bao gio thay log bi mat vi main exit som chua?

3. **Tai sao trong vong `for`, `go func(){ fmt.Println(i) }()` lai in ra gia tri cuoi cung thay vi moi gia tri khac nhau?** Day la bug Go pho bien nhat. Giai phap la gi? Co may cach fix?

4. **Neu mot goroutine bi panic, cac goroutine khac co bi anh huong khong?** Process co crash khong? Lam the nao de bao ve chuong trinh khoi panic trong goroutine?

### Output Checklist

- [ ] TODO-[1] hoan thanh: `DownloadResult` co du cac truong de luu thong tin ket qua (URL, duong dan, loi, co the thoi gian tai)
- [ ] TODO-[2] hoan thanh: Tao HTTP request dung cach, xu ly duoc loi ket noi
- [ ] TODO-[3] hoan thanh: Kiem tra status code, tra loi ro rang cho cac status khac 200
- [ ] TODO-[4] hoan thanh: Ghi file dung cach, khong do RAM (dung `io.Copy`)
- [ ] TODO-[5] hoan thanh: `sync.WaitGroup` duoc khoi tao dung
- [ ] TODO-[6] hoan thanh: Co co che luu ket qua khong bi race condition
- [ ] TODO-[7] hoan thanh: Goroutine duoc tao dung, khong bi closure capture bug
- [ ] TODO-[8] hoan thanh: Moi goroutine goi `wg.Done()` khi xong
- [ ] TODO-[9] hoan thanh: `wg.Wait()` goi dung cho, khong deadlock
- [ ] TODO-[10] hoan thanh: Output bao cao tong hop ro rang cho Product Owner

### Test Checklist

- [ ] Test case: Tai 3 file nho song song — ca 3 thanh cong, kiem tra file ton tai tren dia  
  *Vi sao quan trong: Day la happy path co ban, phai chay duoc dau tien.*
- [ ] Test case: URL khong ton tai (404) — khong crash, bao cao loi dung  
  *Vi sao quan trong: Server khong phai luc nao cung online.*
- [ ] Test case: 1 URL trong danh sach loi, cac URL khac van thanh cong  
  *Vi sao quan trong: Loi mot file khong duoc lam dung tat ca — fault tolerance.*
- [ ] Test case: Goi `downloadAll` voi slice rong  
  *Boundary case: Edge case hay quen, co the gay panic neu khong xu ly.*
- [ ] Test case: Chay `go test -race` va khong co race condition  
  *Vi sao quan trong: Day la quy tac vang cua phase 2 — zero race.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Tai sao khong phai luc nao song song cung nhanh hon? Neu file rat nho
   (1KB), goroutine co hieu qua khong? Chi phi khoi tao goroutine so voi thoi gian download?

2. **Neu requirement thay doi:** Product Owner muon tai 10,000 file thay vi 100.
   Ban se thay doi gi? Co can gioi han so goroutine chay dong thoi khong?

3. **Architecture decision:** Tai sao toi chon WaitGroup thay vi channel cho bai nay?
   Khi nao WaitGroup phu hop hon, khi nao channel phu hop hon? Cau tra loi cua ban
   co thay doi neu ban can streaming ket qua thay vi doi het xong moi in?

4. **Bai hoc thuc te:** "Hoi toi o project pipeline automation, toi tung co 1 goroutine
   leak vi quen goi wg.Done() trong 1 nhanh cua if-else. Memory tang 200MB/ngay.
   Bay gio toi luon viet `defer wg.Done()` ngay sau `wg.Add(1)` — khong bao gio quen nua."
```

---

## Topic 02.2: Channels (Buffered & Unbuffered)

### User Story

> Khach hang (Product Owner) noi: "Cac worker download xong phai gui ket qua ve cho main bao cao. Khong duoc dung shared variable — toi da bi race condition roi, khong muon thay nua."

**Context:** Sau khi dung goroutine, team phat hien co bug race condition khi nhieu goroutine cung ghi vao 1 slice ket qua. Product Owner muon mot co che "communication" sach se, khong share memory.

### Acceptance Criteria

- [ ] Dung channel de truyen ket qua tu worker ve main
- [ ] Hieu ro su khac nhau giua buffered va unbuffered channel
- [ ] Khong bi block vo han (deadlock) khi worker xong ma channel khong duoc doc
- [ ] Dong channel dung cach — chi sender dong, khong dong tu receiver
- [ ] Dung `range` over channel de doc het ket qua

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Channel la cai hay nhat cua Go. Nhung ma dung sai thi la con ac mong.
> Hoi toi interview mot ban junior, code dung channel nhung bi deadlock vi
> quen dong channel. Kho khan hon nua, ban ay con dong channel tu receiver —
> may man la khong panic, nhung khong phai luc nao cung may man vay.
>
> Van de cot loi o day: Do Not Communicate by Sharing Memory.
> Instead, Share Memory by Communicating. — day la cau noi noi tieng cua Rob Pike.
>
> Toi se phan ra:
> 1. Tao channel de nhan ket qua tu workers
> 2. Worker gui ket qua vao channel khi xong
> 3. Main doc tu channel — dung range hoac select
> 4. Dong channel khi khong con gi de gui — chi sender dong
>
> Quan trong: Unbuffered channel = synchronous — gui va nhan phai gap nhau.
> Buffered channel = bat dong bo — co the gui len den N phan tu ma khong block.
> Khi nao dung cai nao? Toi thuong dung buffered channel khi khong muon
> worker bi block cho den khi main doc."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"fmt"
	"sync"
)

// Result chua ket qua cong viec cua worker
type Result struct {
	ID    int
	Value string
	Err   error
}

// worker thuc hien cong viec va gui ket qua qua channel
// SENIOR ASKS: Tai sao parameter la `chan<- Result` thay vi `chan Result`?
// HINT: Go ho tro directional channel — chi gui, chi nhan, hoac ca hai.
func worker(id int, task string, resultChan chan<- Result) {
	// TODO-[1]: Thuc hien cong viec (simulate bang time.Sleep)

	// TODO-[2]: Tao Result va gui vao channel
	// SENIOR ASKS: Nen gui ket qua truoc hay sau khi xu ly loi?
	// HINT: Ca thanh cong va that bai deu la ket qua can bao cao.

	// TODO-[3]: Neu co loi, van gui Result voi Err duoc set
}

// workerPool chay nhieu worker va thu thap ket qua qua channel
// SENIOR ASKS: Tai sao ham nay tra ve []Result thay vi dung channel?
// HINT: API cua ham nen ro rang — caller muon du lieu, khong muon xu ly channel.
func workerPool(tasks []string) []Result {
	// TODO-[4]: Tao buffered channel voi kich thuoc phu hop
	// SENIOR ASKS: Nen dung buffered hay unbuffered? Kich thuoc buffer la bao nhieu?
	// HINT: Buffered = worker khong bi block. Kich thuoc = so worker thi du.

	// TODO-[5]: Tao WaitGroup de doi workers xong
	// SENIOR ASKS: Neu da co channel de nhan ket qua, tai sao van can WaitGroup?
	// HINT: Channel va WaitGroup phuc vu 2 muc dich khac nhau.

	// TODO-[6]: Khoi dong workers trong goroutine
	for i, task := range tasks {
		// SENIOR ASKS: Closure capture bug co the xay ra o day khong?
		// HINT: Chi can 1 cau lenh `defer wg.Done()` de khong bao gio quen.
	}

	// TODO-[7]: Goroutine rieng de dong channel khi tat ca workers xong
	// SENIOR ASKS: Tai sao phai dong channel trong goroutine rieng?
	// HINT: Neu dong channel truoc khi workers xong, dieu gi xay ra?

	// TODO-[8]: Doc ket qua tu channel — dung range
	// SENIOR ASKS: Tai sao range over channel tu dong dung khi channel dong?
	// HINT: Day la co che quan trong de biet khi nao het du lieu.

	return nil // TODO: tra ve slice ket qua
}

// demoUnbuffered cho thay su khac nhau giua buffered va unbuffered
func demoUnbuffered() {
	// TODO-[9]: Tao unbuffered channel va demo synchronous send/receive
	// SENIOR ASKS: Goi send truoc hay receive truoc trong unbuffered channel?
	// HINT: Unbuffered = ca hai phai gap nhau — nhu bat tay.
}

// demoBuffered cho thay buffered channel hoat dong nhu the nao
func demoBuffered() {
	// TODO-[10]: Tao buffered channel size 2, gui 3 gia tri
	// SENIOR ASKS: Send thu 3 co block khong? Tai sao?
	// HINT: Buffered chi block khi buffer day.
}

func main() {
	tasks := []string{"task-A", "task-B", "task-C", "task-D", "task-E"}

	// TODO-[11]: Goi workerPool va in ket qua

	// TODO-[12]: (Tuy chon) Demo buffered vs unbuffered
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Neu ban tao unbuffered channel, goi send tu goroutine A, va goroutine B chua san sang nhan — dieu gi xay ra voi A?** A co tiep tuc chay duoc khong? Dieu nay co nghia la gi cho performance?

2. **Tai sao chi sender moi duoc phep dong channel?** Dieu gi xay ra neu 2 goroutine cung dong 1 channel? Co cach nao de tranh khong?

3. **Khi ban dung `range` tren channel, no dung khi nao?** Neu channel khong bao gio dong, `range` se lam gi? Day co phai la cach tao leak goroutine khong?

4. **Buffered channel size bao nhieu la "du"?** Size 1 vs size 100 khac nhau gi? Co bao gio ban muon channel co size 0 khong?

### Output Checklist

- [ ] TODO-[1] hoan thanh: Worker thuc hien cong viec va chuan bi ket qua
- [ ] TODO-[2] hoan thanh: Ket qua duoc gui vao channel dung cach
- [ ] TODO-[3] hoan thanh: Ca loi cung duoc gui qua channel (khong bo qua)
- [ ] TODO-[4] hoan thanh: Buffered channel duoc tao voi kich thuoc hop ly
- [ ] TODO-[5] hoan thanh: WaitGroup dong bo giua workers va dong channel
- [ ] TODO-[6] hoan thanh: Workers chay trong goroutine, khong bi closure capture bug
- [ ] TODO-[7] hoan thanh: Channel duoc dong trong goroutine rieng, dung thoi diem
- [ ] TODO-[8] hoan thanh: Main doc ket qua bang range over channel
- [ ] TODO-[9] hoan thanh: Demo unbuffered channel chay duoc
- [ ] TODO-[10] hoan thanh: Demo buffered channel chay duoc
- [ ] TODO-[11] hoan thanh: Ket qua cuoi cung duoc in ro rang

### Test Checklist

- [ ] Test case: 3 workers chay song song, tat ca thanh cong, ket qua day du  
  *Vi sao quan trong: Day la happy path — phai chac chan truoc khi them complexity.*
- [ ] Test case: Worker bi loi nhung van gui ket qua ve channel  
  *Vi sao quan trong: Khong duoc lang quen loi — loi cung la ket qua.*
- [ ] Test case: Chay `go test -race` khong phat hien race  
  *Vi sao quan trong: Day la muc tieu chinh — no shared memory, no race.*
- [ ] Test case: Dong channel tu receiver — chung minh dieu nay co the gay panic  
  *Boundary case: Hieu tai sao quy tac "chi sender dong" ton tai.*
- [ ] Test case: Tao deadlock co chu y va chung minh cach phat hien  
  *Vi sao quan trong: Biet deadlock khi thay — debug de hon neu biet pattern.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Ban da dung buffered channel. Neu doi thanh unbuffered,
   performance co thay doi khong? Khi nao unbuffered lai tot hon buffered?
   (Goi y: unbuffered dam bao backpressure tu dong.)

2. **Neu requirement thay doi:** Product Owner muon khong chi nhan ket qua,
   ma con muon "progress bar" — cap nhat real-time moi khi 1 worker xong.
   Ban se thay doi architecture nhu the nao? Channel van du khong, hay can them co che?

3. **Architecture decision:** Tai sao khong dung 1 shared slice co mutex thay vi channel?
   Go proverb noi "Share memory by communicating", nhung mutex van ton tai trong stdlib.
   Khi nao mutex la lua chon tot hon? (Cau tra loi: khi du lieu duoc share nhieu,
   communication chi lam phuc tap them.)

4. **Bai hoc thuc te:** "Toi tung thay mot codebase co 20 channel truyen qua lai
   giua cac goroutine — debug 1 bug mat 2 ngay vi khong biet du lieu di dau.
   Channel la con dao hai luoi: dung 1-2 channel thi sach, dung 10 channel thi ma.
   KISS principle van ap dung cho concurrency."
```

---

## Topic 02.3: Select & Timeout

### User Story

> Khach hang (Product Owner) noi: "Toi can timeout neu download qua 5 giay. Khong de treo chuong trinh — toi khong bao gio muon thay man hinh dung yen khong phan hoi."

**Context:** Product Owner tung dung mot tool download ma bi treo vi server khong phan hoi. Khong co timeout, chuong trinh doi mai. Bay gio requirement ro rang: moi download toi da 5 giay, het 5 giay thi bao loi va tiep tuc.

### Acceptance Criteria

- [ ] Dung `select` de lang nghe ket qua tu nhieu channel
- [ ] Timeout 5 giay cho moi download — khong de treo
- [ ] Dung `time.After` hoac `context.WithTimeout` de implement timeout
- [ ] Xu ly timeout gracefully — khong panic, bao cao loi ro rang
- [ ] Dung `default` case khi can non-blocking select

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Select la 'switch' cho channel. No la construct quyen luc nhat trong Go concurrency.
> Hoi toi viet mot health-check service, can ping 10 services moi 30 giay.
> Dung select + timeout, service nao khong phan hoi trong 3 giay thi mark DOWN.
> Ma khong can goroutine rieng cho timeout — chi can time.After.
>
> Van de cot loi:
> 1. Select cho phep doi nhieu channel — cai nao san sang truoc thi xu ly truoc
> 2. time.After tra ve channel ma se nhan gia tri sau duration
> 3. Ket hop select + time.After = timeout tu nhien
>
> Toi se phan ra:
> 1. Tao ham download voi timeout bang select + time.After
> 2. Trong select: 1 case nhan ket qua, 1 case timeout
> 3. Xet default case cho non-blocking operations
> 4. For-select loop cho continuous listening
>
> Mot dieu can nho: select chon ngau nhien neu nhieu case san sang.
> Dieu nay co nghia la khong co 'priority' — neu can priority, phai viet logic phu."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// DownloadResult la ket qua tai file
type DownloadResult struct {
	URL     string
	Success bool
	Err     error
}

// downloadWithTimeout tai file voi timeout xac dinh
// SENIOR ASKS: Tai sao dung select trong ham nay ma khong phai WaitGroup?
// HINT: WaitGroup doi completion, select cho phep 'hoac ket qua, hoac timeout'.
func downloadWithTimeout(url string, timeout time.Duration) DownloadResult {
	resultChan := make(chan DownloadResult, 1)

	// TODO-[1]: Khoi dong goroutine de thuc hien download
	// SENIOR ASKS: Tai sao resultChan phai la buffered (size 1)?
	// HINT: Neu unbuffered, goroutine co the block mai neu timeout xay ra truoc.

	// TODO-[2]: Dung select de cho ket qua HOAC timeout
	// SENIOR ASKS: select co preference giua cac case khong?
	// HINT: Select chon ngau nhien neu nhieu case san sang cung luc.

	// TODO-[3]: Case 1: Nhan ket qua thanh cong

	// TODO-[4]: Case 2: Timeout
	// SENIOR ASKS: time.After co memory leak khong? Khi nao?
	// HINT: Day la cau hoi interview pho bien — nghĩ ve goroutine cua time.After.

	// TODO-[5]: Tra ve ket qua phu hop
	return DownloadResult{URL: url, Success: false}
}

// downloadMultiple voi timeout cho moi file
func downloadMultiple(urls []string, timeout time.Duration) []DownloadResult {
	// TODO-[6]: Tao slice hoac channel de luu ket qua

	// TODO-[7]: Khoi dong goroutine cho moi URL
	// SENIOR ASKS: Co nen gioi han so goroutine chay dong thoi khong?
	// HINT: Neu co 1000 URL, ban co muon 1000 goroutine cung luc?

	// TODO-[8]: Thu thap ket qua — dung select hoac WaitGroup

	return nil
}

// tickerDemo cho thay cach dung time.Ticker voi select
// SENIOR ASKS: Ticker va Timer khac nhau gi?
// HINT: Ticker lap lai, Timer chi 1 lan. Rat de nham.
func tickerDemo() {
	// TODO-[9]: Tao ticker moi 1 giay
	// SENIOR ASKS: Tai sao phai Stop() ticker khi xong?
	// HINT: Leak goroutine — ticker chay mai neu khong stop.

	// TODO-[10]: Dung for-select de xu ly tick
}

// nonBlockingDemo cho thay default case trong select
func nonBlockingDemo() {
	ch := make(chan int)

	// TODO-[11]: Non-blocking receive bang select voi default
	// SENIOR ASKS: Khi nao default case duoc chon?
	// HINT: Default chi chon khi KHONG case nao san sang.
}

func main() {
	urls := []string{
		"https://example.com/slow-file.jpg",   // server chậm
		"https://example.com/fast-file.jpg",   // server nhanh
		"https://example.com/broken-url",      // khong ton tai
	}

	// TODO-[12]: Goi downloadMultiple voi timeout 5 giay
	// TODO-[13]: In ket qua cuoi cung

	// TODO-[14]: (Tuy chon) Chay tickerDemo va nonBlockingDemo
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Neu ban dung `time.After(5 * time.Second)` trong 1 vong lap, moi lan lap co tao 1 goroutine moi khong?** Neu co, dieu do co nghia la gi cho memory? Lam the nao de tranh?

2. **Trong `select`, neu 2 case cung san sang cung luc, Go chon case nao?** Co phai luon case dau tien khong? Tai sao viec nay quan trong khi viet code production?

3. **Khi nao ban dung `default` case trong select?** No co phai la "else" cua select khong? Hay co muc dich rieng?

4. **Tai sao `for { select { ... } }` lai la pattern pho bien cho goroutine "daemon"?** Goroutine nay se chay mai mai — lam sao de dung no khi can?

### Output Checklist

- [ ] TODO-[1] hoan thanh: Goroutine download chay trong nen, buffered channel duoc dung
- [ ] TODO-[2] hoan thanh: Select voi nhieu case (ket qua + timeout)
- [ ] TODO-[3] hoan thanh: Case nhan ket qua xu ly dung
- [ ] TODO-[4] hoan thanh: Case timeout xu ly dung, khong memory leak
- [ ] TODO-[5] hoan thanh: Tra ve ket qua co y nghia cho ca 2 truong hop
- [ ] TODO-[6] hoan thanh: Co che luu tru ket qua nhieu download
- [ ] TODO-[7] hoan thanh: Goroutine duoc quan ly, khong leak
- [ ] TODO-[8] hoan thanh: Tat ca ket qua duoc thu thap day du
- [ ] TODO-[9] hoan thanh: Ticker duoc tao va stop dung cach
- [ ] TODO-[10] hoan thanh: For-select voi ticker chay dung
- [ ] TODO-[11] hoan thanh: Non-blocking select demo chay duoc
- [ ] TODO-[12] hoan thanh: Download voi timeout hoat dong
- [ ] TODO-[13] hoan thanh: Bao cao ket qua ro rang

### Test Checklist

- [ ] Test case: Download thanh cong truoc khi timeout — ket qua dung  
  *Vi sao quan trong: Happy path — timeout khong duoc anh huong.*
- [ ] Test case: Server qua cham, timeout sau 5 giay — khong treo  
  *Vi sao quan trong: Day la muc tieu chinh cua requirement.*
- [ ] Test case: Ket hop ca thanh cong va timeout trong cung 1 batch  
  *Vi sao quan trong: Fault isolation — 1 download loi khong anh huong cai khac.*
- [ ] Test case: Chay 100 lan de kiem tra memory leak tu time.After  
  *Boundary case: Day la bug hay gap voi time.After trong loop.*
- [ ] Test case: `go test -race` khong phat hien race  
  *Vi sao quan trong: Quy tac vang cua phase 2.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** time.After trong loop co memory leak. Neu ban can timeout
   lap di lap lai, ban se dung gi? (Goi y: time.NewTimer + Reset)

2. **Neu requirement thay doi:** Product Owner muon "progressive timeout" —
   download file nho thi 5 giay, file lon thi 30 giay. Ban se thay doi
   signature cua ham nhu the nao? Context co ho tro dieu nay khong?

3. **Architecture decision:** Ban da dung select + time.After. Cach khac la
   dung `context.WithTimeout`. So sanh 2 cach — khi nao dung cai nao?
   (Goi y: context tot hon khi can propagate cancellation xuong nhieu tang.)

4. **Bai hoc thuc te:** "Mot lan toi debug mot service bi OOM. Nguyen nhan:
   `time.After` trong 1 vong lap chay 1000 lan/phut, moi lan tao 1 goroutine
   ma khong bao gio duoc garbage collected. Fix: dung `time.NewTimer` va
   `Reset()`. Bay gio toi luon de y den time.After trong loop."
```

---

## Topic 02.4: Sync Package

### User Story

> Khach hang (Product Owner) noi: "Nhieu goroutine cung ghi vao mot slice ket qua. Du lieu bi race condition — cung 1 index bi ghi de, so lieu khong khop."

**Context:** Ban dang xay dung mot aggregator service — nhieu worker goroutine cung ghi vao shared data structure. Phat hien: counter sai, slice thieu phan tu, doi khi panic "concurrent map write".

### Acceptance Criteria

- [ ] `sync.Mutex` bao ve shared data — khong co race condition
- [ ] `sync.RWMutex` cho read-heavy workloads — doc song song, ghi doc quyen
- [ ] `sync.WaitGroup` doi tat ca goroutine hoan thanh
- [ ] Phan biet khi nao dung mutex, khi nao dung channel
- [ ] Chay `go test -race` pass 100%

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Mutex hay channel? Cau hoi kinh dien. Toi thuong chon khi:
> - Dung channel khi: communicate between goroutines, pipeline pattern, fan-out/fan-in
> - Dung mutex khi: bao ve 1 shared state don gian (counter, cache, config)
>
> Hoi toi review code cua mot ban, thay ban dung channel de bao ve 1 counter.
> Code phuc tap: phai tao goroutine lang nghe channel, xu ly message, tra ve ket qua.
> Toi refactor thanh sync/atomic.AddInt64 — code giam 80%, nhanh hon 10x.
> Nhung roi ban ay lai dung atomic cho 1 struct phuc tap — khong duoc,
> atomic chi cho primitive types.
>
> Van de cot loi:
> 1. Mutex = mutual exclusion. 1 goroutine duoc phep ghi tai 1 thoi diem.
> 2. RWMutex = nhieu reader cung luc, hoac 1 writer.
> 3. WaitGroup = dem so goroutine dang chay, block den khi het.
>
> Rule of thumb: neu shared state phuc tap (struct, slice, map) → mutex.
> Neu flow du lieu quan trong hon → channel."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"fmt"
	"sync"
)

// SafeCounter la counter thread-safe dung Mutex
type SafeCounter struct {
	// TODO-[1]: Them mutex va value
	// SENIOR ASKS: Nen dung Mutex hay RWMutex cho counter don gian?
	// HINT: Counter chi co increment — khong can read lock.
}

// Inc tang counter len 1 — thread-safe
func (c *SafeCounter) Inc() {
	// TODO-[2]: Lock, tang, Unlock
	// SENIOR ASKS: Tai sao phai dung defer Unlock?
	// HINT: Panic giua Lock va Unlock = deadlock neu khong defer.
}

// Value tra ve gia tri hien tai — thread-safe
func (c *SafeCounter) Value() int {
	// TODO-[3]: Lock de doc
	// SENIOR ASKS: Nen dung Lock hay RLock cho read operation?
	// HINT: Co the co nhieu goroutine doc cung luc.
}

// ResultCollector thu thap ket qua tu nhieu workers — dung Mutex
type ResultCollector struct {
	mu      sync.Mutex
	results []string
}

// Add them ket qua vao collector
func (rc *ResultCollector) Add(result string) {
	// TODO-[4]: Lock va append
	// SENIOR ASKS: Tai sao append slice khong phai la atomic operation?
	// HINT: Append co the allocate memory moi — khong thread-safe.
}

// All tra ve copy cua results
func (rc *ResultCollector) All() []string {
	// TODO-[5]: Lock va tra ve copy
	// SENIOR ASKS: Tai sao phai tra ve copy, khong tra ve slice truc tiep?
	// HINT: Caller co the modify slice — race condition ngam.
}

// Cache la read-heavy cache dung RWMutex
type Cache struct {
	// TODO-[6]: Them RWMutex va data storage
	// SENIOR ASKS: Map co thread-safe khong?
	// HINT: Khong. Concurrent map read/write = panic hoac race.
}

// Set luu gia tri vao cache
func (c *Cache) Set(key, value string) {
	// TODO-[7]: Lock cho write
}

// Get doc gia tri tu cache
func (c *Cache) Get(key string) (string, bool) {
	// TODO-[8]: RLock cho read
	// SENIOR ASKS: RWMutex co overhead so voi Mutex khong?
	// HINT: Co, nhung tot hon cho read-heavy workload.
}

func demoMutex() {
	var wg sync.WaitGroup
	counter := SafeCounter{}

	// TODO-[9]: Tao 100 goroutine, moi goroutine tang counter 1000 lan
	// SENIOR ASKS: Neu khong dung Mutex, ket qua mong doi la gi?
	// HINT: Race condition — ket qua < 100,000.

	wg.Wait()
	fmt.Println("Counter:", counter.Value()) // Nen la 100,000
}

func demoRWMutex() {
	cache := Cache{}
	cache.Set("key1", "value1")

	var wg sync.WaitGroup

	// TODO-[10]: 10 goroutine doc (Get), 2 goroutine ghi (Set)
	// SENIOR ASKS: RWMutex cho phep bao nhieu reader cung luc?
	// HINT: Khong gioi han — mien la khong co writer.

	wg.Wait()
}

func main() {
	// TODO-[11]: Chay demoMutex
	// TODO-[12]: Chay demoRWMutex
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao `defer mu.Unlock()` lai la best practice thay vi goi `Unlock()` truc tiep?** Co truong hop nao `defer` khong phu hop khong?

2. **Neu ban quen goi `Unlock()` trong 1 nhanh cua if-else, dieu gi xay ra?** Lam sao de khong bao gio quen?

3. **RWMutex co cho phep "upgrade" tu read lock sang write lock khong?** Neu co writer dang cho, reader moi co duoc phep vao khong?

4. **Khi nao `sync.Once` huu ich hon mutex?** Cho vi du thuc te ma Once la lua chon hoan hao.

### Output Checklist

- [ ] TODO-[1] hoan thanh: SafeCounter co mutex va value field
- [ ] TODO-[2] hoan thanh: Inc() thread-safe voi Lock/Unlock
- [ ] TODO-[3] hoan thanh: Value() thread-safe, co the dung RLock
- [ ] TODO-[4] hoan thanh: ResultCollector.Add() thread-safe
- [ ] TODO-[5] hoan thanh: ResultCollector.All() tra ve copy
- [ ] TODO-[6] hoan thanh: Cache co RWMutex va map
- [ ] TODO-[7] hoan thanh: Cache.Set() thread-safe voi Lock
- [ ] TODO-[8] hoan thanh: Cache.Get() thread-safe voi RLock
- [ ] TODO-[9] hoan thanh: demoMutex chay ra dung ket qua 100,000
- [ ] TODO-[10] hoan thanh: demoRWMutex voi concurrent readers/writers
- [ ] TODO-[11] hoan thanh: demoMutex chay thanh cong
- [ ] TODO-[12] hoan thanh: demoRWMutex chay thanh cong

### Test Checklist

- [ ] Test case: 100 goroutine dong thoi Inc() — ket qua dung 100,000  
  *Vi sao quan trong: Day la proof mutex hoat dong — khong mutex thi < 100,000.*
- [ ] Test case: Concurrent read/write tren Cache — khong panic  
  *Vi sao quan trong: Go map khong thread-safe — test phai chung minh cache bao ve duoc.*
- [ ] Test case: Chay `go test -race` — zero race  
  *Vi sao quan trong: Day la muc tieu chinh cua topic.*
- [ ] Test case: ResultCollector.All() tra ve copy — chung minh caller khong the modify  
  *Boundary case: Defensive copy — tranh race condition ngam.*
- [ ] Test case: So sanh Mutex vs RWMutex voi benchmark  
  *Vi sao quan trong: Hieu trade-off — RWMutex khong phai luc nao cung tot hon.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Ban da dung Mutex cho counter. sync/atomic co phai la choice
   tot hon khong? So sanh performance bang benchmark. Khi nao atomic,
   khi nao mutex, khi nao channel?

2. **Neu requirement thay doi:** Product Owner muon them "Delete" vao Cache.
   Ban se dung Lock hay RLock cho Delete? Tai sao?

3. **Architecture decision:** sync.Map co san trong stdlib. Tai sao toi lai khong
   dung no ma tu viet Cache? sync.Map co un diem gi so voi map+mutex?
   (Goi y: sync.Map tot cho specific use case — nhieu core, read-heavy, key types khac nhau.
   Nhung khong phai silver bullet.)

4. **Bai hoc thuc te:** "Debug 1 race condition trong production mat 6 tieng.
   Nguyen nhan: 1 goroutine doc map trong khi goroutine khac dang delete.
   Tu do toi luon chay `go test -race` truoc moi commit. Va luon nho:
   'Shared data + goroutine = mutex hoac channel. Khong co exception.'"
```

---

## Topic 02.5: Context

### User Story

> Khach hang (Product Owner) noi: "Khi user cancel request (dong browser), tat ca goroutine dang chay phai dung. Toi khong muon server van xu ly request da bi huy."

**Context:** Ban dang viet mot HTTP API server. Khi client dong ket noi (vi du: user refresh browser), cac goroutine dang xu ly request do van tiep tuc chay — lang phi tai nguyen. Can co co che de "lan truyen" tin hieu huy bo tu tren xuong duoi.

### Acceptance Criteria

- [ ] Context duoc truyen qua tat ca cac ham xu ly
- [ ] Khi context bi cancel, tat ca goroutine con dung lai
- [ ] Timeout tu dong huy bo sau duration xac dinh
- [ ] Goroutine lang nghe `ctx.Done()` de biet khi nao dung
- [ ] Khong bao gio luu context trong struct — chi truyen qua parameter

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Context la API kho hieu nhat nhung quan trong nhat trong Go.
> Hoi toi moi hoc, toi khong hieu tai sao moi ham deu co `ctx context.Context`
> la parameter dau tien. Trong thay lam — nhung sau khi debug 1 goroutine leak
> vi khong dung context, toi hieu: no la lifeline cua moi request.
>
> Quy tac vang ve context:
> 1. Chi main() duoc tao context.Background(). Khong bao gio tao trong ham khac.
> 2. Context la parameter dau tien cua ham: `func Do(ctx context.Context, ...)`
> 3. Khong bao gio luu context trong struct. Truyen qua parameter.
> 4. Derive context: WithCancel, WithTimeout, WithDeadline — de propagate xuong.
> 5. Luon kiem tra `ctx.Err()` de biet ly do cancel.
>
> Van de cot loi: context tao 1 "cancellation tree". Khi root bi cancel,
> tat ca children deu nhan duoc tin hieu. Dieu nay rat manh cho:
> - HTTP request bi client huy
> - Timeout tu dong
> - Graceful shutdown khi nhan SIGTERM"
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// processTask xu ly 1 cong viec, co the bi huy bo
// SENIOR ASKS: Tai sao context phai la parameter dau tien?
// HINT: Go convention — context.Context luon la parameter dau tien.
func processTask(ctx context.Context, taskID int) error {
	// TODO-[1]: Tao ticker de simulate cong viec
	// SENIOR ASKS: Tai sao dung ticker thay vi time.Sleep trong vong lap?
	// HINT: time.Sleep khong the kiem tra ctx.Done() giua chung.

	// TODO-[2]: Vong lap xu ly cong viec
	for i := 0; i < 10; i++ {
		// TODO-[3]: Kiem tra ctx.Done() — neu bi huy thi tra ve som
		// SENIOR ASKS: Nen dung select de kiem tra ctx.Done() hay goi truc tiep?
		// HINT: Select cho phep ket hop ctx.Done() voi cong viec thuc su.

		// TODO-[4]: Thuc hien 1 phan cong viec
		fmt.Printf("Task %d: step %d/10\n", taskID, i+1)

		// TODO-[5]: Sleep 1 giay de simulate processing time
	}

	return nil
}

// processWithTimeout chay task voi timeout
func processWithTimeout(taskID int, timeout time.Duration) error {
	// TODO-[6]: Tao context voi timeout
	// SENIOR ASKS: Tai sao phai defer cancel()?
	// HINT: Memory leak neu khong cancel — timer goroutine chay mai.

	// TODO-[7]: Truyen context xuong processTask
	return processTask(ctx, taskID)
}

// processMultiple chay nhieu task va co the huy bo tat ca
func processMultiple(ctx context.Context, taskIDs []int) {
	// TODO-[8]: Derive cancelable context tu parent
	// SENIOR ASKS: Nen dung WithCancel hay WithTimeout?
	// HINT: Cho phep caller huy bat cu luc nao — WithCancel.

	var wg sync.WaitGroup

	for _, id := range taskIDs {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			// TODO-[9]: Truyen derived context xuong task
			// SENIOR ASKS: Nen truyen ctx hay childCtx?
			// HINT: childCtx — neu parent cancel, child cung cancel.

			// TODO-[10]: Xu ly loi tra ve
		}(id)
	}

	wg.Wait()
}

// gracefulShutdown xu ly SIGTERM va huy bo context
// SENIOR ASKS: Tai sao can graceful shutdown?
// HINT: Process khong duoc dung ngay lap tuc — can doi request dang xu ly.
func gracefulShutdown() {
	// TODO-[11]: Lang nghe OS signal (SIGTERM, SIGINT)
	// SENIOR ASKS: Package nao dung de bat OS signal trong Go?
	// HINT: os/signal — can chu y toi Notify va NotifyContext.

	// TODO-[12]: Khi nhan signal, cancel context

	// TODO-[13]: Doi tat ca goroutine hoan thanh trong khoang thoi gian
	// SENIOR ASKS: Nen cho bao lau truoc khi force exit?
	// HINT: Timeout — khong cho mai, nhung cho du de cleanup.
}

func main() {
	// TODO-[14]: Demo processWithTimeout — task chay het truoc khi timeout
	// TODO-[15]: Demo processWithTimeout — task bi timeout giua chung
	// TODO-[16]: Demo cancel giua chung — tao context, chay task, roi cancel()
	// TODO-[17]: Demo graceful shutdown (neu co the)
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao khong duoc luu context trong struct?** Neu ban co 1 struct `Service` ma nhieu ham dung, tai sao khong luu `ctx context.Context` trong struct de khoi truyen qua parameter?

2. **`ctx.Done()` tra ve channel. Khi nao channel nay duoc dong?** Neu ban goi `<-ctx.Done()` ma khong kiem tra `ctx.Err()`, ban mat di thong tin gi?

3. **Neu ban tao `context.WithTimeout` va khong goi `cancel()`, co van de gi?** Dieu gi xay ra voi goroutine quan ly timer?

4. **`context.WithValue` duoc dung de lam gi?** Tai sao nhieu senior khuyen KHONG NEN dung no de truyen "optional parameters"? Cho vi du dung va sai.

### Output Checklist

- [ ] TODO-[1] hoan thanh: Ticker duoc dung de simulate cong viec
- [ ] TODO-[2] hoan thanh: Vong lap xu ly cong viec co kiem tra cancellation
- [ ] TODO-[3] hoan thanh: ctx.Done() duoc kiem tra trong moi buoc lap
- [ ] TODO-[4] hoan thanh: Cong viec duoc thuc hien tung buoc
- [ ] TODO-[5] hoan thanh: Processing time duoc simulate dung
- [ ] TODO-[6] hoan thanh: Context voi timeout duoc tao, cancel duoc defer
- [ ] TODO-[7] hoan thanh: Timeout duoc propagate xuong processTask
- [ ] TODO-[8] hoan thanh: Cancelable context duoc derive tu parent
- [ ] TODO-[9] hoan thanh: Derived context duoc truyen vao workers
- [ ] TODO-[10] hoan thanh: Loi tra ve duoc xu ly dung (context.Canceled, context.DeadlineExceeded)
- [ ] TODO-[11] hoan thanh: Lang nghe OS signal dung cach
- [ ] TODO-[12] hoan thanh: Context bi cancel khi nhan signal
- [ ] TODO-[13] hoan thanh: Co timeout cho graceful shutdown

### Test Checklist

- [ ] Test case: Task chay het ma khong bi cancel — tra ve nil  
  *Vi sao quan trong: Happy path — context khong duoc can thiep vao luong binh thuong.*
- [ ] Test case: Timeout xay ra truoc khi task xong — tra ve DeadlineExceeded  
  *Vi sao quan trong: Day la muc tieu chinh — phan biet timeout vs manual cancel.*
- [ ] Test case: Manual cancel — tra ve Canceled  
  *Vi sao quan trong: Phan biet ly do cancel de log dung.*
- [ ] Test case: Cancel parent context — tat ca child tasks dung  
  *Vi sao quan trong: Cancellation propagation la core feature cua context.*
- [ ] Test case: `go test -race` pass  
  *Vi sao quan trong: Quy tac vang cua phase 2.*
- [ ] Test case: Khong luu context trong struct — chung minh qua code review  
  *Vi sao quan trong: Day la quy tac quan trong nhat ve context.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Context parameter them 1 doi so vao moi ham. Co worth it khong?
   Khi nao ban se KHONG dung context? (Goi y: background jobs khong can cancel
   co the dung context.Background().)

2. **Neu requirement thay doi:** Product Owner muon "partial result" — khi task
   bi cancel, van tra ve nhung gi da lam duoc. Ban se thay doi ham signature
   va xu ly context nhu the nao?

3. **Architecture decision:** `http.Request` co `r.Context()` roi. Tai sao Go team
   lai chon design nay thay vi truyen context qua parameter? Uu diem la gi?
   (Goi y: Backward compatibility — them field vao Request khong anh huong API.)

4. **Bai hoc thuc te:** "Toi tung thay code nhu sau:
   ```go
   func (s *Service) Process() { // khong co context parameter!
       <-s.ctx.Done() // luu context trong struct!
   }
   ```
   Day la anti-pattern nghiem trong. Context phai la parameter dau tien,
   khong phai field cua struct. Ly do: 1 Service instance duoc dung cho
   nhieu request khac nhau — moi request co context rieng."
```

---

## Topic 02.6: Atomic Operations

### User Story

> Khach hang (Product Owner) noi: "Toi can dem so request da xu ly. Counter don gian nhung phai thread-safe. Khong muon dung mutex vi no cham — day chi la increment mot so."

**Context:** Ban dang viet mot metrics collector — can dem so request, so error, so cache hit. Counter don gian nhung duoc goi rat nhieu lan. Product Owner muon overhead thap nhat co the.

### Acceptance Criteria

- [ ] `sync/atomic` duoc dung de increment/decrement counter
- [ ] Counter thread-safe khong can mutex
- [ ] Benchmark chung minh atomic nhanh hon mutex cho use case nay
- [ ] Hieu han che cua atomic — chi cho primitive types
- [ ] Biet khi nao atomic, khi nao mutex, khi nao channel

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Day la bai toan 'counter don gian' — don gian den muc ban junior deu nghi
> 'chi can bien int roi ++'. Nhung bien int + ++ + nhieu goroutine = race condition.
>
> Toi co 3 lua chon:
> 1. sync.Mutex — an toan nhung cham nhat (context switch, lock contention)
> 2. sync/atomic — nhanh nhat nhung chi cho int32, int64, uint32, uint64, pointer
> 3. Channel — qua heavy cho counter don gian
>
> Kinh nghiem: neu chi la ++/--/read 1 gia tri primitive → atomic.
> Neu la struct, slice, map → mutex.
> Neu la pipeline → channel.
>
> Toi nho lan benchmark: atomic.AddInt64 = ~10ns, mutex.Lock/Unlock = ~50ns.
> 5x khac biet. Nhung neu code phuc tap vi atomic, thi mutex don gian hon.
> 'Premature optimization is the root of all evil' — nhung neu counter la
> hot path (duoc goi trieu lan/giay), atomic la choice dung."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// AtomicCounter dung sync/atomic — nhanh, don gian
type AtomicCounter struct {
	// TODO-[1]: Khai bao value dung atomic
	// SENIOR ASKS: Nen dung int32 hay int64?
	// HINT: Tren 64-bit system, int64 nhanh hon int32. Nhung phai consistent.
	value int64
}

// Inc tang counter len 1
func (c *AtomicCounter) Inc() {
	// TODO-[2]: Dung atomic.AddInt64
	// SENIOR ASKS: Tai sao khong viet c.value++?
	// HINT: c.value++ khong phai la atomic operation — race condition.
}

// Value tra ve gia tri hien tai
func (c *AtomicCounter) Value() int64 {
	// TODO-[3]: Dung atomic.LoadInt64
	// SENIOR ASKS: Tai sao khong doc truc tiep c.value?
	// HINT: Tren mot so architecture, doc truc tiep co the doc gia tri cu.
}

// MutexCounter dung sync.Mutex — cham hon nhung linh hoat
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// BenchmarkAtomic — benchmark atomic counter
func BenchmarkAtomic(b *testing.B) {
	counter := AtomicCounter{}

	// TODO-[4]: Viet benchmark — b.N goroutine, moi goroutine Inc()
	// SENIOR ASKS: Nen dung b.RunParallel khong?
	// HINT: testing.B co RunParallel de benchmark concurrent access.
}

// BenchmarkMutex — benchmark mutex counter
func BenchmarkMutex(b *testing.B) {
	counter := MutexCounter{}

	// TODO-[5]: Tuong tu benchmark mutex
}

func main() {
	// TODO-[6]: Chay counter voi 1000 goroutine, verify ket qua
	// TODO-[7]: Chay benchmark va so sanh
	// TODO-[8]: Demo atomic.StoreInt64 va atomic.CompareAndSwapInt64
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **`atomic.AddInt64` va `c.value++` nhin giong nhau. Khac nhau o cap do nao?** CPU thuc thi chung nhu the nao? Tai sao 1 cai thread-safe, 1 cai khong?

2. **`atomic.Value` co the dung de lam gi ma khong the lam voi atomic.AddInt64?** Cho vi du thuc te.

3. **Khi nao ban KHONG NEN dung atomic du no co ve "du"?** Cho vi du counter phuc tap can nhieu buoc xu ly.

4. **`sync/atomic` va `unsafe` co lien quan gi?** Tai sao atomic an toan ma unsafe thi khong?

### Output Checklist

- [ ] TODO-[1] hoan thanh: AtomicCounter co value int64
- [ ] TODO-[2] hoan thanh: Inc() dung atomic.AddInt64
- [ ] TODO-[3] hoan thanh: Value() dung atomic.LoadInt64
- [ ] TODO-[4] hoan thanh: BenchmarkAtomic viet dung
- [ ] TODO-[5] hoan thanh: BenchmarkMutex viet dung
- [ ] TODO-[6] hoan thanh: Counter chay dung voi 1000 goroutine
- [ ] TODO-[7] hoan thanh: Benchmark chay duoc, co ket qua so sanh
- [ ] TODO-[8] hoan thanh: Demo StoreInt64 va CompareAndSwapInt64

### Test Checklist

- [ ] Test case: 1000 goroutine dong thoi Inc() — ket qua dung 1,000,000  
  *Vi sao quan trong: Chung minh atomic hoat dong — khong atomic thi sai.*
- [ ] Test case: Benchmark — atomic nhanh hon mutex (thuong 3-5x)  
  *Vi sao quan trong: Day la ly do chon atomic cho counter.*
- [ ] Test case: `go test -race` pass  
  *Vi sao quan trong: Quy tac vang — atomic phai thread-safe.*
- [ ] Test case: CompareAndSwap — chung minh hieu semantic "test-and-set"  
  *Boundary case: Day la primitive cho nhieu lock-free algorithms.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Atomic nhanh hon mutex 3-5x. Nhung neu ban can them logic
   vao Inc() (vi du: kiem tra max value), ban co the van dung atomic khong?
   Khi nao ban phai chuyen sang mutex?

2. **Neu requirement thay doi:** Product Owner muon them "gauge" — mot gia tri
   co the tang/giam (khong chi tang). Atomic co ho tro khong? Ban se viet nhu the nao?

3. **Architecture decision:** `expvar` package trong stdlib dung atomic ben duoi.
   Ban co nen dung expvar thay vi tu viet counter? Uu/nhuoc diem?

4. **Bai hoc thuc te:** "Toi tung thay mot ban dung atomic cho 1 struct:
   ```go
   type Stats struct { Total, Errors int }
   var stats atomic.Value // luu *Stats
   ```
   Day la pattern hop le — atomic.Value luu pointer den immutable struct.
   Khi update: tao Stats moi, atomic.Store. Khi read: atomic.Load.
   Read-heavy + occasional write = pattern nay rat hieu qua."
```

---

## Mini-Project: Concurrent Log Analyzer

### User Story

> Khach hang (Product Owner) noi: "Phan tich log file 10GB, dem HTTP status codes, tinh avg response time. Phai xu ly song song, khong load het file vao RAM. Server cua toi chi co 4GB RAM."

**Context:** Day la bai tap interview thuc te ma toi hay ra cho candidate. File log 10GB, moi dong co format: `IP - - [timestamp] "METHOD URL HTTP/1.1" STATUS SIZE "REFERER" "USER_AGENT" RESPONSE_TIME`. Can:
- Dem so lan xuat hien cua moi HTTP status code (200, 404, 500, ...)
- Tinh thoi gian phan hoi trung binh
- Xu ly song song de nhanh
- Khong load het file vao RAM (chi doc tung phan)
- Co graceful shutdown khi nhan SIGTERM

### Acceptance Criteria

- [ ] **Worker pool** — gioi han so goroutine chay dong thoi (configurable, mac dinh = so CPU cores)
- [ ] **Bounded concurrency** — khong de goroutine vuot qua gioi han
- [ ] **Graceful shutdown** — khi nhan SIGTERM, doi workers hoan thanh trong 30 giay roi exit
- [ ] **Khong OOM** — chi doc file theo chunks (vi du: 64KB/lan), khong `ioutil.ReadAll`
- [ ] **Thread-safe aggregation** — counter va avg calculator phai dung mutex/atomic
- [ ] **Bao cao cuoi cung** — in status code counts, avg response time, tong so dong da xu ly
- [ ] **Progress reporting** — in ra man hinh tien do xu ly moi 1 giay

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**

> "Day la bai tap interview thuc te toi hay ra cho candidate. No test du tat ca
> cac kien thuc cua phase 2: goroutine, channel, select, WaitGroup, mutex,
> context, va ca graceful shutdown.
>
> Hoi toi phong van mot ban senior tu Google, ban ay giai quyet trong 45 phut.
> Mot ban junior khac thi mat 2 tieng va code con memory leak. Su khac biet
> khong phai o syntax — ma o viec hieu duoc concurrency patterns.
>
> Architecture toi de xuat:
> 1. 1 goroutine DOC file → gui lines vao channel (producer)
> 2. N goroutine WORKER → nhan lines tu channel, parse, gui ket qua vao channel khac
> 3. 1 goroutine AGGREGATOR → nhan ket qua, cap nhat counter + avg
> 4. Context + signal handler cho graceful shutdown
> 5. Ticker cho progress reporting
>
> Quan trong nhat: backpressure. Neu producer doc nhanh hon worker xu ly,
> channel se day → producer bi block tu dong. Dieu nay ngan OOM.
> Day la ly do buffered channel + worker pool manh me.
>
> Code structure:
> - main: setup, khoi dong cac goroutine, doi completion
> - parseLine: extract status code va response time tu 1 dong log
> - processChunk: worker function — nhan lines, parse, gui ket qua
> - aggregateResult: cap nhat counter va avg
> - reportProgress: in tien do moi giay
>
> Luu y: sync.Map KHONG phai la choice tot o day vi chung ta can iterate
> qua status codes de in bao cao. Map + RWMutex phu hop hon."
```

### TODO Comments (Code Skeleton)

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ParsedLine chua ket qua parse 1 dong log
type ParsedLine struct {
	StatusCode   int
	ResponseTime int64 // milliseconds
}

// Aggregator thu thap va tinh toan ket qua
type Aggregator struct {
	// TODO-[1]: Them mutex bao ve counters
	// SENIOR ASKS: RWMutex hay Mutex? Tai sao?
	// HINT: Chung ta doc nhieu (workers ghi) va in bao cao 1 lan cuoi (doc).

	// TODO-[2]: Them map dem status codes
	// SENIOR ASKS: Map nao? sync.Map hay map + mutex?
	// HINT: sync.Map khong ho tro iterate de in bao cao de dang.

	// TODO-[3]: Them atomic counter cho tong response time va so dong
	// SENIOR ASKS: Tai sao tong response time dung atomic ma khong dung mutex?
	// HINT: int64 + atomic don gian hon mutex.
}

// AddLine cap nhet aggregator voi 1 dong da parse
func (a *Aggregator) AddLine(p ParsedLine) {
	// TODO-[4]: Lock, tang status code count, unlock
	// TODO-[5]: Atomic add response time va line count
}

// Report in bao cao cuoi cung
func (a *Aggregator) Report() {
	// TODO-[6]: Lock de doc, in status counts, tinh avg response time
	// SENIOR ASKS: Tinh avg = total / count — dieu gi xay ra neu count = 0?
	// HINT: Division by zero — luon kiem tra truoc khi chia.
}

// parseLine parse 1 dong log, tra ve ParsedLine va error
// SENIOR ASKS: Ham nay co can nhan context khong?
// HINT: Parse la CPU-bound, khong can timeout. Khong can context.
func parseLine(line string) (ParsedLine, error) {
	// TODO-[7]: Parse status code va response time tu dong log
	// Format: IP - - [timestamp] "METHOD URL HTTP/1.1" STATUS SIZE "REFERER" "AGENT" RT
	// SENIOR ASKS: Nen dung strings.Split, regexp, hay manual scan?
	// HINT: File 10GB — regexp cham. Manual scan hoac Split nhanh hon.

	return ParsedLine{}, nil
}

// processWorker la worker doc tu lineChan, parse, gui ket qua
// SENIOR ASKS: Tai sao khong can WaitGroup trong ham nay?
// HINT: WaitGroup goi o caller, khong o worker.
func processWorker(ctx context.Context, id int, lineChan <-chan string,
	agg *Aggregator, wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO-[8]: Vong for-select lang nghe lineChan va ctx.Done()
	// SENIOR ASKS: Nen dung `for line := range lineChan` hay `for { select { ... } }`?
	// HINT: Range tu dong dung khi channel dong. Nhung range khong lang nghe ctx.Done().

	// TODO-[9]: Parse line va add vao aggregator
	// TODO-[10]: Xu ly loi parse — bo qua hay bao cao?
}

// readFile doc file va gui lines vao channel
// SENIOR ASKS: Ham nay nen return gi?
// HINT: Khong return slice — gui qua channel. Return error neu co.
func readFile(ctx context.Context, filepath string, lineChan chan<- string) error {
	// TODO-[11]: Mo file — dung os.Open
	// SENIOR ASKS: Nen dung bufio.Scanner hay bufio.Reader?
	// HINT: Scanner de dung nhung gioi han 64KB/dong. Reader linh hoat hon.

	// TODO-[12]: Doc tung dong, gui vao channel
	// SENIOR ASKS: Nen kiem tra ctx.Done() giua moi dong khong?
	// HINT: Neu file rat lon va can shutdown nhanh — co, kiem tra moi N dong.

	// TODO-[13]: Dong channel khi doc xong
	// SENIOR ASKS: Ai co trach nhiem dong channel?
	// HINT: Sender dong — day la readFile goroutine.

	return nil
}

// reportProgress in tien do moi interval
func reportProgress(ctx context.Context, agg *Aggregator, interval time.Duration) {
	// TODO-[14]: Tao ticker, in progress moi interval
	// SENIOR ASKS: Ticker co nen Stop() khong?
	// HINT: Co — tranh leak goroutine, du goroutine nay se ket thuc cung chuong trinh.

	// TODO-[15]: Lang nghe ctx.Done() de dung reporting
}

func main() {
	// TODO-[16]: Parse command-line arguments: file path, so workers
	// SENIOR ASKS: Nen dung bao nhieu workers?
	// HINT: runtime.NumCPU() la diem bat dau tot. Nhung co the tune.

	// TODO-[17]: Tao context voi signal handler
	// SENIOR ASKS: context.WithCancel hay NotifyContext?
	// HINT: Go 1.16+ co signal.NotifyContext — tien loi hon.

	// TODO-[18]: Tao buffered channel cho lines
	// SENIOR ASKS: Kich thuoc buffer la bao nhieu?
	// HINT: So worker * 2 hoac so CPU cores — khong can qua lon.

	// TODO-[19]: Khoi tao aggregator

	// TODO-[20]: Khoi dong worker pool voi WaitGroup

	// TODO-[21]: Khoi dong goroutine doc file

	// TODO-[22]: Khoi dong goroutine report progress

	// TODO-[23]: Doi workers xong (wg.Wait())

	// TODO-[24]: In bao cao cuoi cung
}
```

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao lai dung buffered channel giua file reader va workers?** Neu dung unbuffered, dieu gi xay ra voi toc do doc file? Co loi ich gi khi producer bi block tu dong?

2. **Khi user nhan Ctrl+C, worker dang xu ly 1 dong log co bi ngat giua chung khong?** Neu co, dong log do co bi mat khong? Lam sao de chi mat toi da 1 dong, khong phai 1000 dong?

3. **Tinh trung binh bang `total / count` co van de gi voi so nguyen khong?** Neu `total = 1000, count = 3`, ket qua la bao nhieu? Lam sao de chinh xac hon?

4. **Neu file log co 100 trieu dong, ban co can gioi han so goroutine khong?** Tai sao khong de moi dong la 1 goroutine? Go runtime quan ly duoc khong?

### Output Checklist

- [ ] TODO-[1] hoan thanh: Mutex/RWMutex bao ve counters
- [ ] TODO-[2] hoan thanh: Map dem status codes thread-safe
- [ ] TODO-[3] hoan thanh: Atomic counter cho total response time va line count
- [ ] TODO-[4] hoan thanh: AddLine() thread-safe
- [ ] TODO-[5] hoan thanh: Atomic operations cho response time
- [ ] TODO-[6] hoan thanh: Report() in bao cao day du, xu ly division by zero
- [ ] TODO-[7] hoan thanh: parseLine() hoat dong dung voi log format
- [ ] TODO-[8] hoan thanh: processWorker() lang nghe ca lineChan va ctx.Done()
- [ ] TODO-[9] hoan thanh: Worker parse line va add vao aggregator
- [ ] TODO-[10] hoan thanh: Error handling cho parse — khong crash
- [ ] TODO-[11] hoan thanh: Doc file dung bufio, khong ReadAll
- [ ] TODO-[12] hoan thanh: Kiem tra context khi doc file
- [ ] TODO-[13] hoan thanh: Channel dong dung cach boi sender
- [ ] TODO-[14] hoan thanh: Progress reporting chay moi interval
- [ ] TODO-[15] hoan thanh: Progress dung khi context cancel
- [ ] TODO-[16] hoan thanh: CLI arguments duoc parse
- [ ] TODO-[17] hoan thanh: Signal handler + context hoat dong
- [ ] TODO-[18] hoan thanh: Buffered channel voi kich thuoc hop ly
- [ ] TODO-[19] hoan thanh: Aggregator duoc khoi tao
- [ ] TODO-[20] hoan thanh: Worker pool chay dung so luong
- [ ] TODO-[21] hoan thanh: File reader goroutine hoat dong
- [ ] TODO-[22] hoan thanh: Progress reporter goroutine hoat dong
- [ ] TODO-[23] hoan thanh: Chuong trinh doi workers xong truoc khi exit
- [ ] TODO-[24] hoan thanh: Bao cao cuoi cung duoc in

### Test Checklist

- [ ] Test case: File nho (10 dong) — ket qua dem va avg chinh xac  
  *Vi sao quan trong: Happy path — phai dung truoc khi test file lon.*
- [ ] Test case: File 1GB — khong OOM, hoan thanh trong thoi gian hop ly  
  *Vi sao quan trong: Day la muc tieu chinh — khong load het vao RAM.*
- [ ] Test case: File co dong loi (malformed) — bo qua dong loi, khong crash  
  *Vi sao quan trong: Log file thuc te khong phai luc nao cung sach.*
- [ ] Test case: Nhan Ctrl+C giua chung — graceful shutdown, khong panic  
  *Vi sao quan trong: Production code phai xu ly signal dung.*
- [ ] Test case: Chay `go test -race` — zero race  
  *Vi sao quan trong: Quy tac vang cua phase 2.*
- [ ] Test case: Benchmark — so sanh 1 worker vs N workers  
  *Vi sao quan trong: Chung minh concurrency co y nghia.*

### Retrospective

```markdown
### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off:** Ban da dung worker pool voi channel. Cach khac la dung
   `sync/atomic` cho counter va goroutine khong gioi han. So sanh 2 cach
   ve complexity, performance, va maintainability.

2. **Neu requirement thay doi:** Product Owner muoc them "top 10 slowest URLs"
   va "requests per minute histogram". Ban se thay doi data structures va
   concurrency model nhu the nao? Map + mutex van du khong?

3. **Architecture decision:** Tai sao ban khong dung `io.Reader` truc tiep
   ma lai dung channel lam trung gian? Dieu gi xay ra neu bo channel va
   workers truc tiep doc tu 1 shared reader?

4. **Bai hoc thuc te:** "Toi tung deploy 1 log processor tuong tu len production.
   Lan dau chay voi file 50GB, no OOM sau 10 phut. Nguyen nhan: toi dung
   `strings.Split(line, " ")` cho moi dong — tao qua nhieu allocations.
   Fix: dung `fmt.Sscanf` hoac manual parsing. Lesson: concurrency giup nhanh,
   nhung allocation profile moi la key de khong OOM. Hay chay `go tool pprof`
   sau khi benchmark."

5. **Cau hoi bonus:** "Neu ban duoc yeu cau them 'exactly-once processing'
   (khong duoc xu ly 1 dong 2 lan du co worker crash), ban se thay doi
   architecture nhu the nao? Day la boundary giua concurrency va distributed
   systems."
```

---

## Tong ket Phase 2

### Cac kien thuc cot loi da hoc

1. **Goroutines** — chay song song, nhe, nhung phai co co che doi
2. **Channels** — communicate, khong share memory; buffered vs unbuffered
3. **Select** — timeout, non-blocking, random selection
4. **Sync Package** — Mutex, RWMutex, WaitGroup — bao ve shared state
5. **Context** — cancellation tree, timeout, graceful shutdown
6. **Atomic** — nhanh, don gian, primitive-only counters

### Checklist phai pass truoc khi sang Phase 3

- [ ] `go test -race` pass 100% tren tat ca code trong phase nay
- [ ] Mini-project Concurrent Log Analyzer chay duoc, khong OOM, co graceful shutdown
- [ ] Giai thich duoc su khac nhau: channel vs mutex vs atomic — khi nao dung cai nao
- [ ] Viet duoc code pattern: worker pool, fan-out/fan-in, pipeline
- [ ] Hieu duoc: goroutine leak patterns va cach tranh

### Nguon hoc them

- [Go.dev blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [The Go Programming Language, Chapters 8-9](https://www.gopl.io/)
- [Go 101: Concurrency](https://go101.org/article/concurrent-common-mistakes.html)
- [Go by Example: Concurrency Patterns](https://gobyexample.com/)

---

*"Concurrency khong phai la toi uu hoa — no la suy nghi khac ve van de.
Neu ban co the giai quyet van de tuan tu, hay lam vay truoc. Chi dung
concurrency khi no tao ra su khac biet co y nghia."*

*— Loi khuyen cuoi cung tu Senior Developer*
