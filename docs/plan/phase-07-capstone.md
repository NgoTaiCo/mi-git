# Phase 7: Minigit Capstone (Tuần 12)

> "Git khong phai magic. No chi la mot key-value store co tinh toan, cong them mot DAG va mot chut zlib. Hieu duoc dieu nay, ban co the viet lai Git trong mot buoi chieu." — Senior Staff Engineer

---

## Gioi thieu: Capstone tong hop 12 tuan

Phase 7 la capstone cuoi cung cua lop trinh 12 tuan. Day la noi tat ca kien thuc tu Phase 1 den Phase 6 hoi tu thanh mot san pham hoan chinh: **Minigit** — mot he thong version control thu nho voi 3 thanh phan (CLI + Server + Flutter SDK).

Capstone nay duoc thiet ke theo mo hinh **requirement-simulation**: moi topic la mot ticket thuc te ma senior engineer nhan duoc tu Product Owner. Ban se khong duoc code san — ban se duoc dat cau hoi, goi y, va skeleton code de tu dien.

### San pham dau ra

| Thanh phan | Mo ta | Tech stack |
|---|---|---|
| `minigit-cli` | Local repo, object storage, parallel hashing | Go stdlib, SQLite |
| `minigit-server` | HTTP API, auth, object storage | chi, pgx, PostgreSQL, OTel |
| `minigit-sdk` | Dart client, sync, offline queue | Dart, REST + WebSocket |

### Luong hoc tap khuyen nghj

```
Tuan 12 (2 gio/ngay):
- Ngay 1-2: Topic 07.1 — Git Object Model (hieu noi hoat dong)
- Ngay 3-4: Topic 07.2 — CLI Tool (xay dung commands)
- Ngay 5-6: Topic 07.3 — Server (HTTP API + auth)
- Ngay 7-8: Topic 07.4 — Flutter SDK (Dart client)
- Ngay 9-10: Mini-Project — Integration (ket noi tat ca)
- Ngay 11-12: Testing, polish, demo prep
- Ngay 13-14: Retrospective, portfolio write-up
```

---

## Topic 07.1: Git Object Model

### User Story

> **Product Owner noi:** "Toi muon hieu Git hoat dong ben trong: blob, tree, commit la gi? Tai sao Git 'khong the xoa lich su'? Tai sao hash SHA1 lai quan trong?"
>
> **Context:** Day la nen tang ly thuyet cho toan bo he thong minigit. Neu khong hieu object model, cac topic sau se tro nen vo nghia. Product Owner muon team hieu sau truoc khi viet mot dong code nao.

### Acceptance Criteria

- [ ] Tao duoc **blob object** tu file content, tinh duoc SHA1 hash
- [ ] Tao duoc **tree object** tu danh sach file entries (mode + name + sha)
- [ ] Tao duoc **commit object** tu tree SHA + parent SHA + message + author
- [ ] Giai thich duoc tai sao SHA1 la "content-addressable" — cung content = cung hash
- [ ] Ve duoc DAG (Directed Acyclic Graph) tu mot chuoi commit

---

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan requirement nay:**
> "Neu toi nhan ticket nay, dieu dau tien toi nghi den la: Git khong phai magic. 
> No chi la mot key-value store rat don gian — key la SHA1 hash, value la 
> zlib-compressed content. Khi toi hieu dieu nay lan dau tien, toi ngoi 
> trong quan ca phe va viet lai toan bo Git object model trong 2 tieng.
>
> Voi 3 object types: blob, tree, commit — chung ta co du moi thu can thiet 
> de xay dung mot VCS (Version Control System). Blob luu content, tree luu 
> directory structure, commit luu snapshot va metadata.
>
> Vấn de cot loi o day la: content-addressable storage. Cung mot noi dung 
> luon cho cung mot hash. Day la ly do Git khong can luu "file A version 1, 
> version 2" — no chi luu object A va object B. Dedup tu nhien.
>
> Toi se phan ra thanh cac buoc:
> 1. Hieu SHA1 hash — crypto/sha1 trong Go
> 2. Hieu format cua tung object type — header + content
> 3. Hieu zlib compression — compress/zlib
> 4. Hieu DAG — commit -> tree -> blob(s)"
```

---

### TODO Comments (Code Skeleton)

```go
package object

import (
	"crypto/sha1"
	"compress/zlib"
	// TODO-[1]: Import them cac package can thiet
	// SENIOR ASKS: De tinh SHA1 hash, ban can import package nao tu stdlib?
	// HINT: Go co san crypto/sha1, nhung dung voi — co can import "fmt" khong?
)

// ObjectType dinh nghia cac loai object trong Git
// TODO-[2]: Dinh nghia cac const cho object type
// SENIOR ASKS: Git co may loai object chinh? Blob, Tree, Commit, Tag — 
//   nhung capstone nay ta chi can 3 loai dau. Nen dung string hay iota?
// HINT: Dung string type cho de doc. Dung sai type se anh huong den header format.

type ObjectType string

const (
	TypeBlob   ObjectType = "blob"
	TypeTree   ObjectType = "tree"
	TypeCommit ObjectType = "commit"
)

// GitObject la interface chung cho tat ca object types
// TODO-[3]: Interface nay can nhung method nao?
// SENIOR ASKS: Moi object deu co type va content — nhung serialize thi khac nhau.
//   Interface trong Go la implicit — ban khong can "implements" keyword.
// HINT: Think: Serialize(), Type(), Size() co du khong? Con thieu gi khong?

type GitObject interface {
	GetType() ObjectType
	Serialize() ([]byte, error)
}

// BlobObject luu noi dung file
// TODO-[4]: Cau truc cua BlobObject la gi?
// SENIOR ASKS: Blob chi la "binary large object" — no co metadata khong? 
//   Hay chi don thuan la []byte content?
// HINT: Blob trong Git rat don gian: chi co content. Khong co filename, khong co mode.
//   Filename va mode nam trong Tree object.

type BlobObject struct {
	// TODO: Dien vao day
	// SENIOR ASKS: Blob chi can 1 field duy nhat — field do la gi?
}

func (b *BlobObject) GetType() ObjectType {
	// TODO-[5]: Implement GetType
	// SENIOR ASKS: Method nay tra ve gi? Chi 1 dong code.
	return ""
}

func (b *BlobObject) Serialize() ([]byte, error) {
	// TODO-[6]: Implement Serialize cho Blob
	// SENIOR ASKS: Serialize cua Blob co phuc tap khong? 
	//   Blob chi la content thuan — vay Serialize tra ve gi?
	// HINT: Khong can xu ly gi ca. Tra ve Content la du. Nhung doc tiep Tree roi so sanh.
	return nil, nil
}

// TreeEntry la 1 dong trong tree object
// TODO-[7]: Cau truc cua TreeEntry gom nhung gi?
// SENIOR ASKS: Khi chay `git ls-tree`, ban thay gi? Mode, type, sha, name — dung khong?
// HINT: Git format: "100644 blob <sha>\t<filename>". Vay TreeEntry can mode, name, sha.

type TreeEntry struct {
	// TODO: Dien vao day — 3 field
}

// TreeObject luu cau truc thu muc
// TODO-[8]: TreeObject chua gi?
// SENIOR ASKS: 1 tree co nhieu entries — nhieu file, nhieu subdir. 
//   Vay TreeObject chua collection gi?
// HINT: Slice of TreeEntry. Va entries phai duoc sort theo name.

type TreeObject struct {
	// TODO: Dien vao day
}

func (t *TreeObject) GetType() ObjectType {
	// TODO-[9]: Implement
	return ""
}

func (t *TreeObject) Serialize() ([]byte, error) {
	// TODO-[10]: Implement Serialize cho Tree — DAY LA PHAN QUAN TRONG NHAT
	// SENIOR ASKS: Tree format khac Blob format. No khac nhau the nao?
	//   Moi entry duoc format nhu the nao?
	// HINT: Moi entry: "<mode> <name>\0<20-byte-sha-binary>"
	//   Khong phai hex string — la 20 bytes raw! Dung hex.EncodeToString se sai.
	//   Concatenate tat ca entries lai. Entries phai SORT theo name.
	return nil, nil
}

// CommitObject luu metadata va tro den tree
// TODO-[11]: Cau truc cua CommitObject gom nhung gi?
// SENIOR ASKS: Chay `git cat-file -p <commit-sha>` — ban thay gi?
//   tree, parent, author, committer, message — dung khong?
// HINT: Parent la optional (root commit khong co parent). 
//   Timestamp nen dung time.Time. Author + Committer co the gop chung.

type CommitObject struct {
	// TODO: Dien vao day — tree SHA, []parent SHA, author, message, timestamp
}

func (c *CommitObject) GetType() ObjectType {
	// TODO-[12]: Implement
	return ""
}

func (c *CommitObject) Serialize() ([]byte, error) {
	// TODO-[13]: Implement Serialize cho Commit
	// SENIOR ASKS: Commit format trong Git trong nhu the nao?
	//   Moi field tren 1 dong: "tree <sha>\n", "parent <sha>\n", ...
	// HINT: Format: tree\nparent(s)\nauthor\ncommitter\n\nmessage
	//   Khong co "committer" thi dung author lam committer cung duoc.
	return nil, nil
}

// HashObject tinh SHA1 hash tu object
// TODO-[14]: Ham nay la TRAI TIM cua ca he thong
// SENIOR ASKS: SHA1 hash trong Git duoc tinh tu cai gi?
//   Tu raw content? Tu serialized content? Tu compressed content?
// HINT: Git hash = SHA1("<type> <size>\0<content>")
//   Khong phai hash cua compressed content!
//   Ham nay phai tra ve hash string (hex) VA compressed bytes.

func HashObject(obj GitObject) (hash string, compressed []byte, err error) {
	// TODO-[15]: Implement HashObject
	// STEP 1: Serialize object content
	// STEP 2: Tao header: "<type> <size>\0"
	// STEP 3: Concatenate header + content
	// STEP 4: SHA1 hash
	// STEP 5: Zlib compress
	// STEP 6: Tra ve hash (hex string) va compressed bytes
	//
	// SENIOR ASKS: Tai sao lai can compressed bytes luon? 
	//   Hash chi can de lookup, con compressed de luu vao disk.
	// HINT: fmt.Sprintf("%s %d\0", type, len(content)) — nhung dung string concatenation.
	//   crypto/sha1.Sum(data) tra ve [20]byte. hex.EncodeToString de chuyen thanh string.
	return "", nil, nil
}

// Store luu object vao backing storage
// TODO-[16]: Interface Store de decouple storage implementation
// SENIOR ASKS: Tai sao dung interface? Khong phai luc nao cung luu file?
// HINT: Phase nay dung file storage (.git/objects/xx/xxxx...). 
	//   Nhung sau nay co the dung SQLite hoac PostgreSQL.

type Store interface {
	Put(hash string, data []byte) error
	Get(hash string) ([]byte, error)
	Exists(hash string) bool
}

// FileStore luu object tren filesystem (mimic .git/objects)
// TODO-[17]: Cau truc cua FileStore?
// SENIOR ASKS: Git luu objects o dau? Format path nhu the nao?
	// HINT: .git/objects/xx/xxxxxx... — 2 ky tu dau cua hash = directory name, 
	//   con lai = filename. Tai sao? De tranh 1 directory qua nhieu file.

type FileStore struct {
	// TODO: Dien vao day — chi can 1 field: root path
}

func (fs *FileStore) Put(hash string, data []byte) error {
	// TODO-[18]: Implement Put
	// SENIOR ASKS: Lam the nao de tao path tu hash?
	//   filepath.Join, os.MkdirAll, os.WriteFile
	// HINT: hash[:2] la dir, hash[2:] la filename. filepath.Join(fs.Root, "objects", dir, filename)
	return nil
}

func (fs *FileStore) Get(hash string) ([]byte, error) {
	// TODO-[19]: Implement Get
	// SENIOR ASKS: Nguoc lai cua Put. Doc file roi decompress?
	// HINT: os.ReadFile -> zlib decompress. Ham nay co the tra ve raw object content.
	return nil, nil
}

func (fs *FileStore) Exists(hash string) bool {
	// TODO-[20]: Implement Exists
	// SENIOR ASKS: Co nen doc toan bo file? Chi can check file ton tai?
	// HINT: os.Stat la du. Khong can doc noi dung.
	return false
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao Git dung SHA1 ma khong phai MD5?** SHA1 co 160 bits (40 hex chars) — du lon de tranh collision trong thuc te. MD5 chi 128 bits va da bi break. Nhung hay nghi: neu SHA1 collision xay ra, Git se bi anh huong the nao?

2. **Tai sao Git la "content-addressable"?** Co nghia la: ban co the tinh hash TU NOI DUNG, khong can ID tu server. Dieu nay co nghia gi cho distributed system? Git khong can server trung tam de tao ID — moi client tu tinh hash duoc. Ban co thay dieu nay quan trong khong?

3. **Tai sao Tree entries phai SORT theo name?** Git yeu cau entries trong tree phai duoc sort alphabetically. Tai sao? Hint: de 2 client tao cung tree hash khi co cung noi dung. Deterministic = dedup hoat dong.

4. **Neu 2 file khac ten nhung cung content, Git luu may lan?** Chi 1 lan! Vi hash tu content, khong tu filename. Filename nam o Tree object. Day la ly do Git "dedup" rat tot. Co thay trade-off khong? (Rename detection phai lam them viec.)

5. **Tai sao CommitObject khong luu parent duoi dang nil ma lai la slice rong?** Khong, root commit khong co parent. Nhung trong Go, nil slice va empty slice khac nhau. Khi serialize, parent field se khong xuat hien. Dieu nay co dung khong? Hay phai viet code dac biet de handle root commit?

6. **Tai sao compressed content duoc luu, nhung hash lai tinh tu uncompressed content?** De khi ban muon doi compression algorithm (gzip -> zstd), hash khong doi. Smart phai khong? Content-addressable nghia la hash chi phu thuoc vao noi dung semantic, khong phu thuoc encoding.

7. **Tai sao Git object header co dang "type size\0" ma khong phai JSON/XML?** Parsable, compact, va khong can external library. 20 nam truoc khi JSON pho bien, Git da chon format don gian nay. Ban se chon gi neu viet lai tu dau?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Import dung cac package (crypto/sha1, compress/zlib, encoding/hex, fmt, io, os, filepath)
- [ ] TODO-[2] hoan thanh: ObjectType duoc dinh nghia bang string const
- [ ] TODO-[3] hoan thanh: GitObject interface co GetType() va Serialize()
- [ ] TODO-[4..6] hoan thanh: BlobObject struct + methods hoat dong
- [ ] TODO-[7..10] hoan thanh: TreeEntry + TreeObject + Serialize (SORT entries!)
- [ ] TODO-[11..13] hoan thanh: CommitObject struct + Serialize
- [ ] TODO-[14..15] hoan thanh: HashObject tinh dung SHA1 va compress dung zlib
- [ ] TODO-[16..20] hoan thanh: Store interface + FileStore implementation
- [ ] **Integration test:** Tao blob -> hash -> luu -> doc lai -> content giong nhau
- [ ] **Integration test:** Tao tree co 2 blobs -> hash -> serialize dung format
- [ ] **Integration test:** Tao commit tro den tree -> hash -> co dung format
- [ ] **SHA1 test:** Cung content cho cung hash (test determinism)

---

### Test Checklist: Nhung gi ban nen tu viet test

```go
// Test case: Blob roundtrip — vi sao case nay quan trong?
// -> Day la "happy path" co ban nhat. Neu case nay fail, moi thu deu fail.

// Test case: Empty blob — boundary case gi co the fail?
// -> "" cung la 1 valid content. SHA1 cua empty string la e69de29... 
//    Neu code khong handle empty, size = 0 co the gay loi.

// Test case: Tree with entries sorted — vi sao sort quan trong?
// -> 2 entries doi cho nhau cho hash khac. Git yeu cau sorted.
//    Day la loi pho bien nhat khi tu implement tree.

// Test case: Tree entry with non-ASCII filename — boundary case?
// -> Unicode filename can phai duoc handle dung. len(name) != rune count.

// Test case: Commit without parent (root commit) — vi sao case nay dac biet?
// -> Root commit khong co "parent" line trong serialized format.
//    Neu code luon ghi "parent \n", root commit se sai format.

// Test case: Commit with 2 parents (merge commit) — vi sao co the fail?
// -> Merge commit co 2+ parent lines. Code phai handle nhieu parent.

// Test case: SHA1 determinism — vi sao quan trong?
// -> Cung input phai cho cung output. Neu khong, dedup khong hoat dong.
//    Dieu nay co nghia: khong dung random, khong dung timestamp trong hash calculation.

// Test case: FileStore Put then Get roundtrip — vi sao quan trong?
// -> Dam bao serialization + compression + storage + retrieval hoat dong dung.
//    Day la integration test quan trong nhat.

// Test case: FileStore Exists for non-existent hash — vi sao can test?
// -> Ham Exists khong duoc panic khi file khong ton tai. Phai tra ve false nhe nhang.
```

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off: Tai sao Git dung zlib ma khong phai gzip/brotli?** Zlib la "deflate + zlib header". Gzip cung la deflate nhung co header khac. Brotli nen tot hon nhung cham hon. Voi object nho (< 1MB), zlib la sweet spot ve speed/compression. Neu ban viet lai 2026, ban co chon zstd khong? Tai sao?

2. **Neu requirement thay doi: ho tro "large files" (> 100MB)?** Git object model co 1 van lon: blob 100MB -> memory 100MB de hash -> crash tren mobile. Git LFS (Large File Storage) giai quyet bang cach khong luu file trong Git ma thay bang pointer. Ban se design "minigit-lfs" nhu the nao?

3. **Architecture decision: tai sao dung interface cho Store?** Ban co the viet truc tiep FileStore khong can interface. Nhung interface cho phep: test voi in-memory store, swap sang SQLite sau nay ma khong doi code. Day la "dependency inversion" — code phu thuoc vao abstraction, khong phai implementation. Ban thay gia tri nay co xung dang voi complexity khong?

---


## Topic 07.2: CLI Tool

### User Story

> **Product Owner noi:** "Xay dung CLI tool cho minigit: khoi tao repo, them file, tao commit, xem lich su. Giong `git init`, `git add`, `git hash-object`, `git cat-file`, `git write-tree`, `git commit`, `git log` nhung don gian hon."
>
> **Context:** Day la thanh phan CLI cua minigit — user se tuong tac truc tiep qua terminal. Tool nay phai: parse command tot, error message ro rang, exit code chuan (0 = success, 1 = error), va co `--help` cho moi command. Muc tieu: thay duoc noi hoat dong cua Git object model qua terminal.

### Acceptance Criteria

- [ ] Command `init` — tao repo moi voi thu muc `.minigit/`
- [ ] Command `hash-object` — doc file tu stdin hoac filepath, tao blob, luu vao object store
- [ ] Command `cat-file` — doc object tu object store theo hash, in ra content
- [ ] Command `write-tree` — quet working directory, tao tree object tu cac file
- [ ] Command `commit` — tao commit object tu tree hien tai + parent commit + message
- [ ] Command `log` — in ra lich su commit tu HEAD tro ve truoc
- [ ] Content-addressable storage: object duoc luu va lookup bang SHA1 hash
- [ ] Error messages ro rang, khong panic trong flow binh thuong
- [ ] `--help` va `help <command>` hoat dong

---

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan ticket nay:**
> "Hoi toi o project logging tool, toi cung viet 1 CLI tuong tu. Dieu dau tien 
> toi nghi: khong duoc dung Cobra ngay. Hoc stdlib flag truoc — hieu no roi 
> moi dung Cobra khi can. Voi capstone nay, flag package la du.
>
> Vấn de cot loi o day la: parse argv -> dispatch command -> thuc thi -> 
> exit dung code. Nhung phan quan trong nhat la quyet dinh kien truc:
> - Moi command la 1 function rieng biet?
> - Hay dung map[string]func? 
> - Hay dung struct voi interface?
>
> Toi chon cach thu 2: map[string]Command, voi Command la interface co 
> Name(), Synopsis(), Run(). Don gian, de test, de them command moi.
>
> Phan kho nhat la `write-tree`: phai di qua filesystem, doc tung file, 
> tao blob, roi build tree object. Phai handle: .minigit ignore, symlink, 
> permissions. Ban dau toi quen sort entries va hash tree khac nhau moi lan — 
> do la loi pho bien nhat."
```

---

### TODO Comments (Code Skeleton)

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	// TODO-[1]: Import cac package can thiet
	// SENIOR ASKS: De doc file, ban can package gi? De tinh SHA1 thi sao?
	// HINT: "bufio" cho stdin, "io" cho interfaces, "strings" cho manipulation.
)

// Cmd la interface cho tat ca commands
// TODO-[2]: Interface Command can gi?
// SENIOR ASKS: Moi command deu co ten, mo ta, va hanh dong. Nhung flag parsing 
//   co nen nam trong interface khong?
// HINT: Don gian thoi: Name() string, Synopsis() string, Run(args []string) error.

type Cmd interface {
	Name() string
	Synopsis() string
	Run(args []string) error
}

// ======= CMD: init =======
// TODO-[3]: CmdInit struct
// SENIOR ASKS: Init command can gi? Chi can tao thu muc .minigit voi 
//   subdirectories nhat dinh. Subdirs nao?
// HINT: .minigit/objects/ (luu objects), .minigit/refs/heads/ (luu branch refs)

type CmdInit struct {
	Path string // repo root path
}

func (c *CmdInit) Name() string     { return "init" }
func (c *CmdInit) Synopsis() string { return "Khoi tao minigit repository moi" }

func (c *CmdInit) Run(args []string) error {
	// TODO-[4]: Implement init
	// STEP 1: Parse flag de lay path (default = current dir)
	// STEP 2: Tao .minigit/ directory
	// STEP 3: Tao .minigit/objects/ directory
	// STEP 4: Tao .minigit/refs/ va .minigit/refs/heads/ directories
	// STEP 5: (Optional) Tao .minigit/HEAD file voi noi dung "ref: refs/heads/main"
	//
	// SENIOR ASKS: Tai sao can HEAD? HEAD la gi trong Git?
	// HINT: HEAD la "con tro" den commit hien tai. Ban dau chua co commit, 
	//   nen HEAD tro den branch "main" (chua ton tai). 
	//   Format: "ref: refs/heads/main" hoac truc tiep SHA1.
	//   Ban co can luu config nhu ten default branch khong?
	return nil
}

// ======= CMD: hash-object =======
// TODO-[5]: CmdHashObject struct
// SENIOR ASKS: Git co 2 mode: luu object vao store, hoac chi in hash.
//   Flag "-w" (write) de phan biet. Ban se xu ly flag nhu the nao?
// HINT: Dung flag.Bool("w", false, "write object to store") trong Run method.
//   Doc tu stdin HOAC tu filepath (os.Stdin vs os.Open).

type CmdHashObject struct {
	Store object.Store // tu topic 07.1
}

func (c *CmdHashObject) Name() string     { return "hash-object" }
func (c *CmdHashObject) Synopsis() string { return "Tinh SHA1 hash cua file, optional luu vao store" }

func (c *CmdHashObject) Run(args []string) error {
	// TODO-[6]: Implement hash-object
	// STEP 1: Parse flags (-w: write to store)
	// STEP 2: Xac dinh input source: stdin hoac filepath tu args con lai
	// STEP 3: Doc toan bo content vao memory ([]byte)
	// STEP 4: Tao BlobObject tu content
	// STEP 5: Hash + (optional compress + store)
	// STEP 6: In hash ra stdout
	//
	// SENIOR ASKS: Doc toan bo file vao memory co van de gi?
	// HINT: Voi file nho (< 10MB) thi OK. Nhung neu file 1GB thi sao?
	//   Topic "large file" se xu ly sau. Hien tai: io.ReadAll la du.
	//   io.ReadAll(reader) tra ve ([]byte, error) — dung no.
	return nil
}

// ======= CMD: cat-file =======
// TODO-[7]: CmdCatFile struct
// SENIOR ASKS: cat-file co mode "-p" (pretty print) va "-t" (show type).
//   Ham nay phai doc object tu store roi in ra.
// HINT: Store.Get tra ve compressed bytes -> decompress -> parse header 
	//   -> determine type -> pretty print.

type CmdCatFile struct {
	Store object.Store
}

func (c *CmdCatFile) Name() string     { return "cat-file" }
func (c *CmdCatFile) Synopsis() string { return "Hien thi noi dung object theo SHA1 hash" }

func (c *CmdCatFile) Run(args []string) error {
	// TODO-[8]: Implement cat-file
	// STEP 1: Parse flags (-p: pretty print content, -t: show type)
	// STEP 2: Lay hash tu args
	// STEP 3: Doc compressed data tu store
	// STEP 4: Decompress bang zlib
	// STEP 5: Parse header: "<type> <size>\0<content>"
	// STEP 6: Tuy flag ma in type hoac content ra stdout
	//
	// SENIOR ASKS: Lam the nao de parse "<type> <size>\0<content>"?
	// HINT: strings.SplitN(data, "\x00", 2) — split tai null byte.
	//   Header = phan dau, content = phan sau.
	//   strings.Fields(header) de tach "blob 123" thanh ["blob", "123"].
	return nil
}

// ======= CMD: write-tree =======
// TODO-[9]: CmdWriteTree struct — DAY LA COMMAND PHUC TAP NHAT
// SENIOR ASKS: write-tree phai quet toan bo working directory, tao blob 
//   cho moi file, roi build tree object. Thu tu lam viec la gi?
// HINT: 1. Walk directory (filepath.WalkDir) 
//   2. Skip .minigit directory
//   3. Voi moi file: doc content -> hash-object (blob) -> tao TreeEntry
//   4. Sort entries theo name
//   5. Serialize tree -> hash -> store
//   6. In tree hash ra stdout

type CmdWriteTree struct {
	Store object.Store
}

func (c *CmdWriteTree) Name() string     { return "write-tree" }
func (c *CmdWriteTree) Synopsis() string { return "Tao tree object tu working directory" }

func (c *CmdWriteTree) Run(args []string) error {
	// TODO-[10]: Implement write-tree
	// STEP 1: Parse flags (neu co, vi du: --ignore-patterns)
	// STEP 2: filepath.WalkDir tu current dir
	// STEP 3: Skip .minigit/ (internal directory)
	// STEP 4: Voi moi file thuong (khong phai dir):
	//   a. Doc content
	//   b. Tao BlobObject
	//   c. HashBlob -> luu vao store (write blob)
	//   d. Tao TreeEntry voi mode="100644", name=relative-path, sha=blob-hash
	// STEP 5: Sort entries by name (quan trong!)
	// STEP 6: Tao TreeObject tu entries
	// STEP 7: HashTree -> luu vao store
	// STEP 8: In tree hash
	//
	// SENIOR ASKS: filepath.WalkDir callback co signature gi?
	// HINT: func(path string, d fs.DirEntry, err error) error
	//   d.IsDir() de check directory. d.Type() de check file mode.
	//   Quan trong: PHAI skip .minigit/ khong se loop vo han!
	return nil
}

// ======= CMD: commit =======
// TODO-[11]: CmdCommit struct
// SENIOR ASKS: Commit can: tree hash, parent commit hash (neu co), 
//   message, author. Tree hash tu dau? Parent tu dau?
// HINT: Tree hash co the tu args (-m "msg" -t <tree-hash>) hoac 
	//   auto-detect tu index (neu ban implement index). 
	//   Parent doc tu HEAD -> resolve ref -> lay commit SHA.
	//   Don gian hoa: tree hash tu args, parent tu HEAD.

type CmdCommit struct {
	Store object.Store
	RepoPath string
}

func (c *CmdCommit) Name() string     { return "commit" }
func (c *CmdCommit) Synopsis() string { return "Tao commit moi" }

func (c *CmdCommit) Run(args []string) error {
	// TODO-[12]: Implement commit
	// STEP 1: Parse flags (-m "commit message", -t <tree-hash>)
	// STEP 2: Neu khong co -t: doc "index" hoac auto-detect (don gian: yeu cau -t)
	// STEP 3: Doc HEAD de lay parent commit hash (neu co)
	// STEP 4: HEAD = "ref: refs/heads/main" -> doc file refs/heads/main -> lay SHA
	// STEP 5: HEAD khong ton tai -> root commit (khong co parent)
	// STEP 6: Tao CommitObject voi tree, parent(s), message, author, timestamp
	// STEP 7: Hash + store commit
	// STEP 8: Update HEAD tro den commit moi (ghi SHA vao refs/heads/main)
	// STEP 9: In commit hash ra stdout
	//
	// SENIOR ASKS: HEAD file co 2 format: "ref: refs/heads/main" hoac 
	//   truc tiep SHA1 string. Ban se parse nhu the nao?
	// HINT: strings.HasPrefix("ref: ") — neu co prefix "ref: " thi la 
	//   symbolic ref, neu khong thi la direct ref (SHA1 truc tiep).
	return nil
}

// ======= CMD: log =======
// TODO-[13]: CmdLog struct
// SENIOR ASKS: Log can in commit tu HEAD tro ve truoc. Day la duyet DAG.
//   Voi single-parent commits, log la linked list traversal.
//   Format output: commit <sha>, Author, Date, Message.
// HINT: Doc HEAD -> lay commit SHA -> doc commit object -> in info -> 
	//   lay parent SHA -> lap lai. Dung for loop, co dieu kien dung 
	//   khi khong con parent (root commit).

type CmdLog struct {
	Store object.Store
	RepoPath string
}

func (c *CmdLog) Name() string     { return "log" }
func (c *CmdLog) Synopsis() string { return "Hien thi lich su commit" }

func (c *CmdLog) Run(args []string) error {
	// TODO-[14]: Implement log
	// STEP 1: Doc HEAD -> resolve thanh commit SHA
	// STEP 2: Vong lap:
	//   a. Doc commit object tu store
	//   b. In commit info (sha, author, date, message)
	//   c. Lay parent SHA -> tro thanh commit SHA tiep theo
	//   d. Neu khong con parent -> dung
	// STEP 3: Handle truong hop HEAD chua ton tai (repo moi, chua commit)
	//
	// SENIOR ASKS: Neu commit co 2 parents (merge commit), log nen di 
	//   theo parent nao? Git di theo "first parent" — ban co nen vay khong?
	// HINT: Ban dau chi support single parent. Merge commit -> lay parent[0].
	return nil
}

// ======= Main =======
// TODO-[15]: Ham main — dispatcher
// SENIOR ASKS: main() phai lam gi? Parse command dau tien, tim command 
//   trong registry, goi Run(). Neu khong tim thay -> help message.
// HINT: os.Args[1] la command name, os.Args[2:] la args cho command.
//   Dung map[string]Cmd de lookup. Default: in usage va exit(1).

func main() {
	// TODO-[16]: Implement main
	// STEP 1: Khoi tao FileStore voi root = .minigit/
	// STEP 2: Tao map[string]Cmd voi tat ca commands
	// STEP 3: Parse os.Args[1] de lay command name
	// STEP 4: Dispatch toi command tuong ung
	// STEP 5: Neu loi: fmt.Fprintln(os.Stderr, err) + os.Exit(1)
	// STEP 6: Neu khong co command: in help + os.Exit(1)
	//
	// SENIOR ASKS: Tai sao in loi ra stderr ma khong phai stdout?
	// HINT: POSIX convention: stdout cho output data, stderr cho diagnostics.
	//   Dieu nay quan cho piping: `minigit log | head` chi pipe stdout.
}

// TODO-[17]: Ham printHelp
// SENIOR ASKS: Help message nen trong nhu the nao? Git co help rat tot.
// HINT: In ten tool, usage, danh sach commands voi synopsis.
//   fmt.Fprintf(os.Stderr, "Usage: %s <command> [args]\n\nCommands:\n", os.Args[0])
//   roi lap qua map commands.
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao dung map[string]Cmd thay vi switch-case dai?** Voi 6 commands, switch-case cung duoc. Nhung khi co 20 commands thi sao? Map cho phep dong goi moi command trong file rieng, test rieng. Ban thay trade-off gi? (Kho debug hon? Runtime error neu key khong ton tai?)

2. **Tai sao `hash-object` doc tu stdin HOAC file?** Git cho phep ca 2. Nhung trong thuc te, khi nao dung stdin? Khi pipe: `echo "hello" | minigit hash-object --stdin`. Tai sao dieu nay quan trong cho Unix philosophy?

3. **Neu `write-tree` gap symlink hoac executable file, mode la gi?** Git dung: 100644 (regular file), 100755 (executable), 120000 (symlink). Ban co can ho tro ca 3 trong capstone? Nen gioi han pham vi nhu the nao?

4. **Tai sao HEAD file la "ref: refs/heads/main" thay vi SHA1 truc tiep?** Neu HEAD luu SHA1 truc tiep, moi lan checkout branch khac phai ghi lai HEAD. Dung "ref:" format, HEAD tro den 1 ref file — khi branch di chuyen (co commit moi), HEAD tu dong "thay doi" vi no chi la pointer. Ban co thay loi ich khong?

5. **`commit` command update HEAD nhu the nao?** Ghi de file refs/heads/main voi SHA1 moi. Neu HEAD la "ref: refs/heads/main", ban doc de biet update file nao. Neu HEAD la SHA1 truc tiep (detached HEAD), ghi de chinh HEAD. Ban se implement detached HEAD khong?

6. **Neu `log` duyet 10000 commits, co van de gi?** Vong lap for don gian se hoat dong, nhung memory? Moi iteration load 1 commit object roi bo — khong co leak. Tuy nhien, neu in ra stdout va user pipe vao `head -5`, chuong trinh se nhan SIGPIPE. Ban co handle signal khong?

7. **Exit code quan trong nhu the nao trong CLI tool?** Exit code 0 = success, 1 = general error, 2 = misuse. Script shell kiem tra `$?` de quyet dinh. Neu tool luon tra ve 0 du loi, CI/CD pipeline se khong biet fail. Ban da tung debug loi do exit code sai chua?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Import dung packages (flag, fmt, os, path/filepath, strings, bufio, io, compress/zlib, v.v.)
- [ ] TODO-[2] hoan thanh: Cmd interface duoc dinh nghia ro rang
- [ ] TODO-[3..4] hoan thanh: CmdInit tao repo structure dung
- [ ] TODO-[5..6] hoan thanh: CmdHashObject tinh hash dung, flag -w hoat dong
- [ ] TODO-[7..8] hoan thanh: CmdCatFile doc object dung, flag -p va -t hoat dong
- [ ] TODO-[9..10] hoan thanh: CmdWriteTree tao tree dung, skip .minigit/, sort entries
- [ ] TODO-[11..12] hoan thanh: CmdCommit tao commit dung, update HEAD
- [ ] TODO-[13..14] hoan thanh: CmdLog in lich su commit tu HEAD
- [ ] TODO-[15..16] hoan thanh: main() dispatch dung, exit code chuan, error ra stderr
- [ ] TODO-[17] hoan thanh: printHelp hoat dong
- [ ] **Integration test:** `init` -> `hash-object` -> `cat-file` roundtrip
- [ ] **Integration test:** `write-tree` -> `cat-file -p` hien thi dung entries
- [ ] **Integration test:** `commit` -> `log` hien thi commit vua tao
- [ ] **Error test:** Command khong ton tai -> exit(1) + help message

---

### Test Checklist: Nhung gi ban nen tu viet test

```go
// Test case: init tao dung directory structure
// -> Check .minigit/objects/ ton tai, .minigit/refs/heads/ ton tai
// -> Case nay co ban — neu fail, moi test sau deu fail

// Test case: hash-object --stdin
// -> Pipe stdin vao hash-object, verify output hash dung
// -> vi du: echo "hello" | minigit hash-object --stdin
//    SHA1 cua "blob 6\0hello\n" = b6fc4c... (tinhh lai cho chinh xac)

// Test case: hash-object -w file.txt
// -> Verify file duoc luu vao .minigit/objects/xx/xxxx...

// Test case: cat-file -p <hash> in dung content
// -> Roundtrip: hash-object -> cat-file = content goc

// Test case: cat-file -t <hash> in dung type
// -> blob, tree, commit — moi type test 1 lan

// Test case: write-tree voi directory trong
// -> Empty tree co hash dac biet: 4b825dc642cb6eb9a060e54bf8d69288fbee4904
//    Day la SHA1 cua "tree 0\0" (empty tree)

// Test case: write-tree voi 2 files
// -> Verify tree co 2 entries, sort theo alphabet

// Test case: commit dau tien (root commit)
// -> Khong co parent, HEAD duoc update

// Test case: commit thu 2 (co parent)
// -> Parent = commit truoc, HEAD duoc update

// Test case: log voi 3 commits
// -> In dung thu tu: commit 3, commit 2, commit 1

// Test case: Unknown command
// -> exit code != 0, help message

// Test case: cat-file voi hash khong ton tai
// -> Error message ro rang, khong panic
```

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off: flag package vs Cobra/Viper?** Flag package: 0 dependency, hoc duoc stdlib, nhug khong co subcommand auto-complete, khong co global flags. Cobra: nhieu feature, ecosystem lon, nhung dependency nang. Voi 6 commands, flag du. Voi 20+ commands, Cobra xung dang. Ban se chon gi neu viet production CLI tool?

2. **Neu requirement thay doi: ho tro "staging area" (index)?** Git co "index" (staging area) de chon file cu the cho commit. Hien tai `write-tree` quet toan bo directory. De ho tro staging, ban can file `.minigit/index` luu danh sach file da "add". `commit` se doc index thay vi scan directory. Dieu nay lam architecture phuc tap ra sao?

3. **Architecture decision: tai sao moi command la struct thay vi function?** Struct cho phep dependency injection (Store, RepoPath). Function khong co state. Voi struct, test de hon vi co the mock Store. Ban thay day la over-engineering cho CLI tool don gian khong?

---


## Topic 07.3: Server

### User Story

> **Product Owner noi:** "Xay dung HTTP API cho minigit server. Users co the dang ky, dang nhap, push objects len server, fetch objects tu server. Auth bang JWT. API phai co structured logging, graceful shutdown, va health checks."
>
> **Context:** Day la thanh phan server cua minigit — cung cap centralized storage cho objects, ho tro collaboration giua nhieu user. Server nay se duoc deploy bang Docker va co the chay tren Fly.io/Railway. Muc tieu: production-ready HTTP API voi observability.

### Acceptance Criteria

- [ ] POST `/api/v1/register` — dang ky user moi (username, password)
- [ ] POST `/api/v1/login` — dang nhap, tra ve JWT token
- [ ] POST `/api/v1/repos` — tao repo moi (can auth)
- [ ] POST `/api/v1/repos/{id}/objects` — upload blob/tree/commit (can auth)
- [ ] GET `/api/v1/repos/{id}/objects/{sha}` — fetch object (can auth)
- [ ] GET `/api/v1/repos/{id}/refs/{name}` — lay commit SHA cua 1 ref (can auth)
- [ ] POST `/api/v1/repos/{id}/refs/{name}` — update ref (can auth)
- [ ] GET `/api/v1/health` — health check (khong can auth)
- [ ] Structured logging (slog) voi JSON handler
- [ ] Graceful shutdown (SIGTERM/SIGINT handling)
- [ ] OpenTelemetry traces cho moi request
- [ ] JWT middleware xac thuc Bearer token

---

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan ticket nay:**
> "Day la capstone server — tich hop tat ca tu Phase 3 den Phase 6. 
> Khi toi nhan ticket API server, day la quy trinh tu duy:
>
> 1. Dau tien: API contract. Viet spec truoc, code sau. 
>    Dung OpenAPI/Swagger neu can, nhung it nhat phai co 
>    document request/response shape.
>
> 2. Thu hai: Auth strategy. JWT stateless = de scale nhung 
>    khong the revoke token (phai doi secret hoac dung blacklist).
>    Voi capstone nay, JWT co secret key don gian la du.
>
> 3. Thu ba: Router + middleware chain. Toi chon `chi` vi no 
>    nhe, compatible voi stdlib http.Handler, va hoat dong 
>    giong middleware stack rat ro rang.
>
> 4. Thu tu: Storage layer. Metadata (users, repos, refs) 
>    -> PostgreSQL. Object content (blob, tree, commit) -> 
>    co the luu trong DB (bytea) hoac S3-compatible. 
>    Voi capstone: PostgreSQL cho ca 2.
>
> 5. Thu nam: Observability. slog cho logging, OTel cho tracing, 
>    Prometheus metrics. Health check endpoint bat buoc.
>
> 6. Cuoi cung: Graceful shutdown. Khi nhan SIGTERM, 
>    http.Server.Shutdown() cho in-flight requests hoan thanh 
>    trong 30s roi moi exit. Khong drop requests bat ngo.
>
> Vấn de cot loi: separation of concerns. Handler chi parse 
> request/tra response. Service layer xu ly business logic. 
> Repository layer xu ly DB. Khong viet SQL trong handler!"
```

---

### TODO Comments (Code Skeleton)

```go
// ======= cmd/server/main.go =======
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// TODO-[1]: Import cac package can thiet
	// SENIOR ASKS: Can nhung package nao cho server?
	// HINT: "github.com/go-chi/chi/v5" cho routing, "github.com/jackc/pgx/v5/pgxpool" 
	//   cho PostgreSQL, "github.com/golang-jwt/jwt/v5" cho JWT.
	//   "log/slog" cho logging, "go.opentelemetry.io/otel" cho tracing.
)

// TODO-[2]: Config struct
// SENIOR ASKS: Server can nhung config gi?
// HINT: Port, DatabaseURL, JWTSecret, ShutdownTimeout, LogLevel.
//   Doc tu env vars. Fail fast neu thieu required config.

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	ShutdownTimeout  time.Duration
}

func loadConfig() *Config {
	// TODO-[3]: Implement config loading
	// SENIOR ASKS: Nen dung env var hay config file?
	// HINT: 12-factor app khuyen dung env vars. 
	//   os.Getenv("PORT"), os.Getenv("DATABASE_URL").
	//   Neu thieu required: log.Fatal hoac return error.
	return nil
}

func main() {
	// TODO-[4]: Setup structured logging (slog)
	// SENIOR ASKS: slog.New voi handler nao cho production?
	// HINT: slog.NewJSONHandler(os.Stdout, nil) cho structured JSON logs.
	//   De dang parse boi log aggregation systems (Datadog, CloudWatch).
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// TODO-[5]: Connect to PostgreSQL
	// SENIOR ASKS: Nen dung *sql.DB hay pgxpool.Pool?
	// HINT: pgxpool.Pool cho async operations va type-rich PostgreSQL.
	//   pgxpool.New(context.Background(), databaseURL) -> (*pgxpool.Pool, error)
	//   Nho defer pool.Close()
	
	// TODO-[6]: Khoi tao repository layer
	// SENIOR ASKS: Tai sao can repository layer? Khong viet SQL trong handler?
	// HINT: Repository pattern: tach DB access ra khoi business logic.
	//   De test: mock repository thay vi mock DB.
	//   UserRepo, RepoRepo, ObjectRepo, RefRepo — moi cai 1 interface.

	// TODO-[7]: Khoi tao service layer
	// SENIOR ASKS: Service layer lam gi khac repository?
	// HINT: Service = business logic. Vi du: "tao repo" = check ten hop le 
	//   + check chua ton tai + insert DB. Repository chi lam "insert DB".

	// TODO-[8]: Khoi tao JWT helper
	// SENIOR ASKS: JWT helper nen chua nhung ham?
	// HINT: GenerateToken(userID string) (string, error) 
	//   va ValidateToken(token string) (*Claims, error)
	//   Claims struct: jwt.RegisteredClaims + UserID string `json:"user_id"`

	// TODO-[9]: Khoi tao router (chi)
	r := chi.NewRouter()
	
	// TODO-[10]: Global middleware
	// SENIOR ASKS: Middleware nao can ap dung global?
	// HINT: Logger (log moi request), Recoverer (recover panic -> 500),
	//   RequestID (gan ID moi request de trace), CORS (neu can).
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// TODO-[11]: Route definitions
	// SENIOR ASKS: Nhung routes nao, phan nao can auth?
	// HINT: Public routes: /api/v1/register, /api/v1/login, /api/v1/health
	//   Authenticated routes: tat ca con lai — dung JWT middleware group.
	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Get("/health", handler.Health)

		// Authenticated
		r.Group(func(r chi.Router) {
			r.Use(jwtMiddleware)
			r.Post("/repos", handler.CreateRepo)
			r.Post("/repos/{id}/objects", handler.UploadObject)
			r.Get("/repos/{id}/objects/{sha}", handler.FetchObject)
			r.Get("/repos/{id}/refs/{name}", handler.GetRef)
			r.Post("/repos/{id}/refs/{name}", handler.UpdateRef)
		})
	})

	// TODO-[12]: Tao HTTP server voi timeout
	// SENIOR ASKS: Tai sao can timeout? DefaultClient khong timeout = nguy hiem.
	// HINT: http.Server{Addr, Handler, ReadTimeout, WriteTimeout, IdleTimeout}
	//   ReadTimeout: 5s, WriteTimeout: 10s — du cho API don gian.

	// TODO-[13]: Graceful shutdown
	// SENIOR ASKS: Graceful shutdown nghia la gi? Tai sao quan trong?
	// HINT: Khi nhan SIGTERM (Docker stop, K8s pod termination), 
	//   khong dung ngay lap tuc. Cho requests dang xu ly hoan thanh.
	//   Pattern: 1 goroutine chay server.ListenAndServe()
	//   1 goroutine cho signal -> server.Shutdown(ctx)

	// TODO-[14]: Start server
	// SENIOR ASKS: server.ListenAndServe() vs server.ListenAndServeTLS()?
	// HINT: Capstone nay khong can TLS (reverse proxy xu ly). 
	//   Production: chi chay HTTP, de Nginx/Traefik xu ly TLS.
}
```

```go
// ======= internal/handler/auth.go =======
// TODO-[15]: Auth handlers
// SENIOR ASKS: Register handler can lam gi? Co nhung buoc validate nao?
// HINT: 1. Parse JSON body (username, password)
//   2. Validate: username >= 3 chars, password >= 8 chars
//   3. Hash password (bcrypt.GenerateFromPassword)
//   4. Luu user vao DB (UserRepo.Create)
//   5. Tra ve 201 Created + user info (khong tra password hash!)

// TODO-[16]: Login handler
// SENIOR ASKS: Login flow nhu the nao?
// HINT: 1. Parse JSON body (username, password)
//   2. Tim user theo username (UserRepo.GetByUsername)
//   3. So sanh password (bcrypt.CompareHashAndPassword)
//   4. Generate JWT token (JWT helper)
//   5. Tra ve 200 OK + {"token": "<jwt>"}

// TODO-[17]: JWT Middleware
// SENIOR ASKS: Middleware hoat dong nhu the nao trong chi?
// HINT: http.HandlerFunc: doc Authorization header -> extract Bearer token 
	//   -> validate token -> gan userID vao request context 
	//   -> goi next handler.
	//   Neu invalid: tra 401 Unauthorized.
```

```go
// ======= internal/handler/object.go =======
// TODO-[18]: Object handlers
// SENIOR ASKS: UploadObject can lam gi?
// HINT: 1. Doc repo ID tu URL param (chi.URLParam)
	//   2. Doc object content tu request body (io.ReadAll)
	//   3. Tinh SHA1 hash (dung lai code tu Topic 07.1)
	//   4. Check object chua ton tai (ObjectRepo.Exists)
	//   5. Luu object vao DB (ObjectRepo.Create)
	//   6. Tra ve 201 + {"sha": "<hash>"}

// TODO-[19]: FetchObject handler
// SENIOR ASKS: FetchObject can lam gi? Response format?
// HINT: 1. Doc repo ID va SHA tu URL params
	//   2. Lookup object trong DB (ObjectRepo.Get)
	//   3. Tra ve 200 + object content (Content-Type: application/octet-stream)
	//   4. Neu khong tim thay: 404 Not Found
```

```go
// ======= internal/repository/object_repo.go =======
// TODO-[20]: Object repository interface + PostgreSQL implementation
// SENIOR ASKS: Interface co nhung method nao?
// HINT: Create(ctx, repoID, sha, objectType, content) error
	//   Get(ctx, repoID, sha) (content, type, error)
	//   Exists(ctx, repoID, sha) (bool, error)

// TODO-[21]: SQL queries
// SENIOR ASKS: Table objects co columns gi? Query nhu the nao?
// HINT: INSERT INTO objects (repo_id, sha, type, content, created_at) 
	//   VALUES ($1, $2, $3, $4, $5) ON CONFLICT (repo_id, sha) DO NOTHING
	//   
	//   SELECT type, content FROM objects WHERE repo_id = $1 AND sha = $2
	//   
	//   content nen dung BYTEA trong PostgreSQL. 
	//   Voi object nho (< 1MB), BYTEA la OK. Voi object lon, can S3.
```

```go
// ======= internal/repository/user_repo.go =======
// TODO-[22]: User repository
// SENIOR ASKS: User table co columns gi? Password luu nhu the nao?
// HINT: id (UUID), username (unique), password_hash (bcrypt), created_at
	//   KHONG BAO GIO luu plain text password!
	//   bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

```go
// ======= internal/service/object_service.go =======
// TODO-[23]: Object service — business logic
// SENIOR ASKS: Service layer can gi khi upload object?
// HINT: 1. Check repo ton tai va user co quyen (RepoRepo.Get + ownership check)
	//   2. Tinh SHA1 tu content de verify (client co the gui hash trong header)
	//   3. Luu object qua repository
	//   4. Tra ve ket qua
	//   
	//   Dieu nay tach biet: handler khong biet DB, service khong biet HTTP.
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao dung repository pattern thay vi viet SQL truc tiep trong handler?** 3 layer (handler -> service -> repository) co ve over-engineering cho capstone. Nhung trong production voi 50+ endpoints, noi ma ban viet SQL trong handler, code se thanh "spaghetti" sau 6 thang. Ban nghi dieu gi xay ra khi doi DB tu PostgreSQL sang MySQL?

2. **JWT token co the revoke khong?** JWT la stateless — server khong luu session. Khi user "dang xuat", token van con valid cho den khi expire. Lam cach nao de revoke? (Blacklist? Short TTL + refresh token? State server-side?) Ban chon cach nao cho capstone, cach nao cho production banking app?

3. **Tai sao health check khong can auth?** Health check duoc goi boi load balancer, monitoring system, K8s liveness probe — khong co token. Neu yeu cau auth, load balancer se nghi server down va stop traffic. Ban co nen tra thong tin chi tiet (DB connection status) trong health check khong?

4. **bcrypt "cost" nen la bao nhieu?** bcrypt.DefaultCost = 10. Moi +1 cost, thoi gian hash x2. Cost 10 ~ 100ms. Cost 14 ~ 1.6s. Trade-off gi? (Bao mat vs user experience khi login.) Banking app nen chon cost bao nhieu?

5. **Object content nen luu trong PostgreSQL (BYTEA) hay S3?** BYTEA: don gian, 1 query lay ca metadata + content, nhung khong scale voi file lon. S3: scale tot, nhung can 2 round-trip (1 DB lookup metadata, 1 S3 fetch). Voi capstone nay, tai sao BYTEA la lua chon hop ly?

6. **Graceful shutdown timeout nen la bao lau?** 30s? 60s? Neu request dang upload 100MB object va mat 45s, ma timeout la 30s, upload se bi drop. Nhung neu timeout qua lau, K8s se force kill sau 30s anyway. Ban tinh timeout dua tren dieu gi?

7. **Tai sao dung `pgx` thay vi `database/sql` voi `lib/pq`?** pgx: native driver, hoat dong tot hon, ho tro type rich (UUID, JSONB, arrays), nhung phai hoc API moi. database/sql: standard, portable, nhung khong ho tro advanced PostgreSQL features. Ban se chon gi khi team chua ai dung pgx?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Import dung packages (chi, pgx, jwt, slog, otel)
- [ ] TODO-[2..3] hoan thanh: Config struct + loadConfig tu env vars
- [ ] TODO-[4] hoan thanh: slog JSON handler setup
- [ ] TODO-[5] hoan thanh: PostgreSQL connection pool
- [ ] TODO-[6] hoan thanh: Repository interfaces + implementations
- [ ] TODO-[7] hoan thanh: Service layer voi business logic
- [ ] TODO-[8] hoan thanh: JWT helper (Generate + Validate)
- [ ] TODO-[9..10] hoan thanh: chi router + global middleware
- [ ] TODO-[11] hoan thanh: Route definitions (public + authenticated)
- [ ] TODO-[12] hoan thanh: HTTP server voi timeout
- [ ] TODO-[13] hoan thanh: Graceful shutdown (SIGTERM/SIGINT)
- [ ] TODO-[14] hoan thanh: Server start
- [ ] TODO-[15] hoan thanh: Register handler
- [ ] TODO-[16] hoan thanh: Login handler
- [ ] TODO-[17] hoan thanh: JWT middleware
- [ ] TODO-[18] hoan thanh: UploadObject handler
- [ ] TODO-[19] hoan thanh: FetchObject handler
- [ ] TODO-[20..21] hoan thanh: Object repository + SQL queries
- [ ] TODO-[22] hoan thanh: User repository
- [ ] TODO-[23] hoan thanh: Object service
- [ ] **Integration test:** Register -> Login -> Upload object -> Fetch object roundtrip
- [ ] **Integration test:** Health check tra ve 200 OK
- [ ] **Integration test:** Unauthorized request -> 401
- [ ] **Integration test:** Graceful shutdown khong drop in-flight requests

---

### Test Checklist: Nhung gi ban nen tu viet test

```go
// Test case: Register thanh cong — vi sao quan trong?
// -> Happy path co ban. User duoc tao voi password hash.
//    Response khong chua password hash!

// Test case: Register voi username da ton tai — boundary case?
// -> 409 Conflict. UNIQUE constraint tren username.
//    Can handle PostgreSQL error code 23505 (unique violation).

// Test case: Register voi password qua ngan — validation?
// -> 400 Bad Request. Password < 8 chars.
//    Test ca empty string, null (trong JSON).

// Test case: Login dung credentials — vi sao quan trong?
// -> Tra ve JWT token. Token phai chua user_id claim.

// Test case: Login sai password — security test?
// -> 401 Unauthorized. Error message khong tiet lo "username ton tai".
//    Thong bao chung chung: "invalid credentials".

// Test case: Upload object khong co auth — auth test?
// -> 401 Unauthorized. Verify middleware hoat dong.

// Test case: Upload object voi token het han — edge case?
// -> 401. exp claim trong JWT < now().

// Test case: Fetch object khong ton tai — error handling?
// -> 404 Not Found. Khong panic, khong 500.

// Test case: Upload duplicate object — idempotency?
// -> ON CONFLICT DO NOTHING -> 200 hoac 201 (tuy design).
//    Hash-based dedup: cung content = cung hash = khong can luu lai.

// Test case: Health check khi DB down — what happens?
// -> Co the tra 503 (DB unreachable) hoac van 200 (server is up).
//    Quyet dinh nay phu thuoc vao monitoring strategy.

// Test case: Graceful shutdown — integration test?
// -> Send SIGTERM, verify server dung sau timeout.
//    Cach test: go routine send request cham, send signal, 
//    verify response van nhan duoc.

// Test case: SQL injection attempt — security?
// -> Thu dung SHA param = "'; DROP TABLE objects; --"
//    Query parameterized ($1) nen khong bi inject.
//    Verify table van con.
```

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off: REST vs gRPC cho internal APIs?** Capstone nay dung REST (JSON) cho tat ca vi Flutter client goi duoc de dang. Nhung neu them service Go noi bo, gRPC (protobuf) se nhanh hon, type-safe hon, nhung can gateway de browser/Flutter goi duoc. Ban se design hybrid architecture nhu the nao?

2. **Neu requirement thay doi: ho tro "git clone" (packfile protocol)?** Git clone khong chi don gian la download tung object — no dung "packfile" de nen nhieu object thanh 1 file, delta-compression de giam dung luong. Implement packfile la 1 project rieng, nhung co the don gian hoa: "thin pack" chi include objects ma client chua co. Ban se design API endpoint nao cho use case nay?

3. **Architecture decision: tai sao khong dung ORM?** Capstone nay viet raw SQL. ORM (GORM) tiet kiem thoi gian nhung: 1) khong biet SQL sinh ra, 2) khong optimize duoc, 3) performance bat ngo. Sau 12 tuan hoc Go, ban thay raw SQL co phai la quyet dinh dung cho junior khong? Khi nao ORM xung dang?

---


## Topic 07.4: Flutter SDK

### User Story

> **Product Owner noi:** "Xay dung Dart SDK cho minigit — Flutter client co the dang nhap, sync repo, xem file tree, doc blob content. SDK phai ho tro offline: queue cac thao tac khi mat mang, sync lai khi co mang."
>
> **Context:** Day la thanh phan SDK ma Flutter developer se dung. SDK la 1 Dart package (co the publish len pub.dev) cung cap high-level API cho minigit server. Muc tieu: developer chi goi `MinigitClient.init(...)`, `client.sync()`, `client.readFile(path)` — moi thu phuc tap duoc an di.

### Acceptance Criteria

- [ ] Dart class `MinigitClient` voi config: baseURL, token
- [ ] `login(username, password)` -> JWT token, luu vao `flutter_secure_storage`
- [ ] `createRepo(name)` -> tao repo tren server
- [ ] `pushObject(content)` -> upload blob, tra ve SHA1
- [ ] `fetchObject(sha)` -> download blob content
- [ ] `writeTree(entries)` -> tao tree object tren server (hoac local)
- [ ] `commit(treeSha, message)` -> tao commit
- [ ] `log()` -> lay lich su commit
- [ ] **Offline support:** Queue cac thao tac khi mat mang
- [ ] **Local cache:** Cache objects tren device (SQLite hoac Hive)
- [ ] **Sync:** Tu dong sync khi co mang lai

---

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan ticket nay:**
> "Flutter SDK la cau noi giua Go backend va Flutter UI. Khi toi viet SDK 
> cho team mobile, toi luon nghi ve 3 viec:
>
> 1. API contract phai on dinh. Backend thay doi response shape -> 
>    mobile app tren store bi crash. Dung versioned API (/api/v1/) 
>    va viet contract tests.
>
> 2. Network la unreliable. 3G/4G/WiFi mat bat cu luc nao. 
>    Moi network call phai co: timeout, retry, offline queue. 
>    User khong bao gio duoc mat data vi mat mang.
>
> 3. Token storage phai an toan. JWT token khong duoc luu 
>    vao SharedPreferences (plain text). flutter_secure_storage 
>    dung Keychain (iOS) va Keystore (Android).
>
> Voi offline support, toi chon pattern: 
> - 'Command Queue' — moi thao tac (push, commit) la 1 command.
> - Commands duoc luu vao local DB (SQLite) voi status: pending, synced, failed.
> - Khi co network: dequeue va execute.
> - Khi mat network: enqueue, khong block UI.
>
> Vấn de cot loi: offline-first architecture. App phai hoat dong 
> du khong co mang, sync khi co mang. Day la requirement quan trong 
> nhat cho mobile app."
```

---

### TODO Comments (Code Skeleton)

```dart
// ======= lib/minigit_sdk.dart =======
// TODO-[1]: Export public API
// SENIOR ASKS: Dart package nen expose nhung class nao?
// HINT: MinigitClient, SyncStatus, exceptions. 
//   Khong expose internal classes (network layer, cache layer).

// ======= lib/src/client.dart =======
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'cache.dart';
import 'sync_queue.dart';

// TODO-[2]: MinigitClient class
// SENIOR ASKS: Client nen la singleton hay instance thuong?
// HINT: Instance thuong cho de test va de mock. 
//   Singleton kho test, kho mock, nhung de dung (MinigitClient.instance.xxx).
//   Voi SDK, instance pattern linh hoat hon.

class MinigitClient {
  final String baseUrl;
  final http.Client _httpClient;
  final FlutterSecureStorage _secureStorage;
  final ObjectCache _cache;
  final SyncQueue _syncQueue;

  // TODO-[3]: Constructor
  // SENIOR ASKS: Constructor co nen require token ngay khong?
// HINT: Khong. Token co the null (chua login). 
//   Config chi can baseUrl. Token duoc luu/load tu secure storage.

  MinigitClient({
    required this.baseUrl,
    http.Client? httpClient,
    FlutterSecureStorage? secureStorage,
  })  : _httpClient = httpClient ?? http.Client(),
        _secureStorage = secureStorage ?? const FlutterSecureStorage(),
        _cache = ObjectCache(),
        _syncQueue = SyncQueue();

  // ======= AUTH =======

  // TODO-[4]: Login method
  // SENIOR ASKS: Login can lam gi? Token luu o dau?
  // HINT: 1. POST /api/v1/login voi body {"username": ..., "password": ...}
  //   2. Parse response -> lay token
  //   3. Luu token vao flutter_secure_storage (key: "minigit_token")
  //   4. Tra ve token hoac true/false
  //   5. Handle error: 401 = wrong credentials, 500 = server error

  Future<String> login(String username, String password) async {
    // TODO: Implement
    return '';
  }

  // TODO-[5]: Logout method
  // SENIOR ASKS: Logout chi can xoa token?
  // HINT: 1. Xoa token khoi secure storage
  //   2. (Optional) POST /api/v1/logout de invalidate server-side
  //   3. Clear cache

  Future<void> logout() async {
    // TODO: Implement
  }

  // TODO-[6]: Internal: get stored token
  // SENIOR ASKS: Moi request can token — lay tu dau?
  // HINT: _secureStorage.read(key: "minigit_token"). Cache trong memory?
  //   Co the luu token trong _currentToken de khong doc storage moi lan.

  Future<String?> _getToken() async {
    // TODO: Implement
    return null;
  }

  // TODO-[7]: Internal: authenticated HTTP headers
  // SENIOR ASKS: Header cho authenticated request?
  // HINT: 'Authorization': 'Bearer <token>', 'Content-Type': 'application/json'

  Future<Map<String, String>> _authHeaders() async {
    // TODO: Implement
    return {};
  }

  // ======= OBJECT OPERATIONS =======

  // TODO-[8]: Push blob object
  // SENIOR ASKS: Upload object nhu the nao? Content body la gi?
  // HINT: POST /api/v1/repos/{repoId}/objects
  //   Body: raw bytes (khong phai JSON! vi object la binary)
  //   Headers: Content-Type: application/octet-stream
  //   Response: {"sha": "<hash>"}
  //   Nho check cache truoc: neu object da co trong cache, khong can upload.

  Future<String> pushObject(String repoId, List<int> content) async {
    // TODO: Implement
    // STEP 1: Check network availability
    // STEP 2: Neu offline -> enqueue to sync queue -> return pending SHA
    // STEP 3: Neu online -> POST to server
    // STEP 4: Parse SHA from response
    // STEP 5: Cache locally
    // STEP 6: Return SHA
    return '';
  }

  // TODO-[9]: Fetch blob object
  // SENIOR ASKS: Download object nhu the nao? Response format?
  // HINT: GET /api/v1/repos/{repoId}/objects/{sha}
  //   Response: raw bytes (not JSON)
  //   Check cache truoc: neu da co trong cache, return ngay.
  //   Neu khong: download -> cache -> return.

  Future<List<int>> fetchObject(String repoId, String sha) async {
    // TODO: Implement
    // STEP 1: Check local cache
    // STEP 2: Neu co -> return cached content
    // STEP 3: Neu khong -> GET from server
    // STEP 4: Cache locally
    // STEP 5: Return content
    return [];
  }

  // ======= TREE OPERATIONS =======

  // TODO-[10]: Write tree
  // SENIOR ASKS: Tree co the tao local hay phai goi server?
  // HINT: Tree format la deterministic — co the tao local.
  //   1. Tao tree object local (dung lai logic tu Topic 07.1)
  //   2. Upload tung blob
  //   3. Upload tree object
  //   Hoac: server co endpoint POST /api/v1/repos/{id}/trees

  Future<String> writeTree(String repoId, List<TreeEntry> entries) async {
    // TODO: Implement
    return '';
  }

  // ======= COMMIT OPERATIONS =======

  // TODO-[11]: Create commit
  // SENIOR ASKS: Commit can nhung gi? Parent tu dau lay?
  // HINT: 1. Lay HEAD commit hien tai (fetch ref "main")
    //   2. Tao commit object voi tree SHA, parent, message
    //   3. Upload commit object
    //   4. Update HEAD ref tro den commit moi

  Future<String> commit(String repoId, String treeSha, String message) async {
    // TODO: Implement
    return '';
  }

  // TODO-[12]: Get log
  // SENIOR ASKS: Log lay tu server hay local?
  // HINT: Tu server: GET /api/v1/repos/{id}/commits?ref=main
  //   Tra ve list commits tu moi nhat den cu nhat.
  //   Neu offline: co the lay tu local cache (neu da sync truoc do).

  Future<List<CommitInfo>> log(String repoId, {String ref = 'main'}) async {
    // TODO: Implement
    return [];
  }

  // ======= SYNC =======

  // TODO-[13]: Sync — DAY LA PHUC TAP NHAT
  // SENIOR ASKS: Sync hoat dong nhu the nao?
  // HINT: 1. Lay pending commands tu sync queue
    //   2. Thuc thi tung command tren server
    //   3. Neu thanh cong: mark synced
    //   4. Neu that bai: mark failed, retry sau
    //   5. Kiem tra conflicts (2 devices edit cung file)
    //   6. Pull changes tu server ve local

  Future<SyncResult> sync(String repoId) async {
    // TODO: Implement
    return SyncResult();
  }

  // TODO-[14]: Dispose
  // SENIOR ASKS: Tai sao can dispose?
  // HINT: http.Client giu connections open. Can close khi khong dung.
  //   Dieu nay quan trong de tranh memory leak va socket exhaustion.

  void dispose() {
    _httpClient.close();
  }
}

// ======= lib/src/models.dart =======
// TODO-[15]: Data models
// SENIOR ASKS: Nhung model nao can?
// HINT: CommitInfo (sha, author, message, timestamp), 
//   TreeEntry (mode, name, sha), SyncStatus enum

class CommitInfo {
  final String sha;
  final String author;
  final String message;
  final DateTime timestamp;

  CommitInfo({required this.sha, required this.author, required this.message, required this.timestamp});

  factory CommitInfo.fromJson(Map<String, dynamic> json) {
    // TODO: Implement fromJson
    throw UnimplementedError();
  }
}

class TreeEntry {
  final String mode;
  final String name;
  final String sha;

  TreeEntry({required this.mode, required this.name, required this.sha});
}

enum SyncStatus { pending, syncing, synced, failed }

class SyncResult {
  final SyncStatus status;
  final int pushed;
  final int pulled;
  final List<String> errors;

  SyncResult({this.status = SyncStatus.synced, this.pushed = 0, this.pulled = 0, this.errors = const []});
}

// ======= lib/src/cache.dart =======
// TODO-[16]: ObjectCache — local storage
// SENIOR ASKS: Cache nen dung gi? SQLite? Hive? SharedPreferences?
// HINT: SQLite cho object storage (co the luu binary content).
//   Hoac Hive (key-value, nhanh, nhung khong to cho binary lon).
//   Voi capstone: SQLite la hop ly nhat.

class ObjectCache {
  // TODO: Implement SQLite-based cache
  // get(sha) -> content hoac null
  // put(sha, content) -> luu vao SQLite
  // exists(sha) -> bool
  // clear() -> xoa tat ca
}

// ======= lib/src/sync_queue.dart =======
// TODO-[17]: SyncQueue — offline command queue
// SENIOR ASKS: Queue nen luu o dau? SQLite table?
// HINT: Table: sync_queue (id, repo_id, command_type, payload, status, created_at)
//   Command types: push_object, write_tree, create_commit, update_ref
//   enqueue() -> INSERT
//   dequeue() -> SELECT pending ORDER BY created_at
//   markSynced(id) -> UPDATE status = 'synced'
//   markFailed(id, error) -> UPDATE status = 'failed', error = ?

class SyncQueue {
  // TODO: Implement
  // enqueue(repoId, commandType, payload)
  // dequeue(repoId) -> list pending commands
  // markStatus(id, status)
}

// ======= lib/src/exceptions.dart =======
// TODO-[18]: Custom exceptions
// SENIOR ASKS: Tai sao can custom exceptions thay vi dung Exception don thuan?
// HINT: De caller co the catch cu the. Vi du: AuthException, NetworkException, 
//   ServerException. UI co the hien thi message phu hop.

class MinigitException implements Exception {
  final String message;
  MinigitException(this.message);
}

class AuthException extends MinigitException {
  AuthException(super.message);
}

class NetworkException extends MinigitException {
  NetworkException(super.message);
}

class ServerException extends MinigitException {
  final int statusCode;
  ServerException(this.statusCode, String message) : super(message);
}
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Tai sao `pushObject` co offline queue nhung `fetchObject` khong?** Fetch la "read" — co the fail nhe nhang, user se thu lai. Push la "write" — user muon luu lai, khong chap nhan mat data. Viec "enqueue writes, ignore read failures" la pattern pho bien trong mobile offline-first. Ban co dong y khong?

2. **Conflict resolution khi 2 devices edit cung file?** Device A push file X v1. Device B (offline) edit file X -> push khi co mang. Server da co v1, B gui v1' (khac A). Lam the nao xu ly? Git strategy: merge/rebase. Mobile don gian hon: "last write wins" hoac "manual merge". Ban chon gi?

3. **Tai sao dung `flutter_secure_storage` thay vi SharedPreferences?** SharedPreferences luu plain text — bat ky app nao co root access deu doc duoc. flutter_secure_storage dung iOS Keychain va Android Keystore — hardware-backed encryption tren device moi. JWT token = key to the kingdom. Khong bao gio luu no plain text.

4. **http.Client nen reuse hay tao moi moi request?** Reuse! http.Client giu connection pool (HTTP keep-alive). Tao moi moi lan = tao TCP connection moi = cham 3x-5x. Nhung reuse co van de: can close khi dispose. Memory leak neu khong close.

5. **Cache invalidation strategy nao phu hop?** "There are only two hard things in Computer Science: cache invalidation and naming things." Voi minigit SDK, cache invalidation don gian vi objects la immutable (cung hash = cung content khong doi). Vay khong can invalidate! Chi can eviction (LRU khi day storage). Day la loi ich cua content-addressable storage — ban co thay khong?

6. **Sync queue: xu ly theo thu tu hay song song?** Neu push 3 objects A, B, C — B phu thuoc A (B la tree chua A). Vay phai theo thu tu! Sync queue phai xu ly sequentially, khong song song. Co cach nao song song duoc khong? (Topological sort cua dependency graph.)

7. **`dispose()` quan trong nhu the nao trong Flutter?** Khi widget bi destroy, resources phai duoc giai phong. Neu khong dispose http.Client, memory leak tich tu. Trong Flutter: dung StatefulWidget, goi dispose() trong `dispose()` method. Hoac dung Provider/BLoC de quan ly lifecycle. Ban se chon pattern nao?

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: Public API exports dung
- [ ] TODO-[2..3] hoan thanh: MinigitClient class + constructor
- [ ] TODO-[4] hoan thanh: Login method hoat dong, luu token vao secure storage
- [ ] TODO-[5] hoan thanh: Logout method xoa token
- [ ] TODO-[6] hoan thanh: _getToken doc tu secure storage
- [ ] TODO-[7] hoan thanh: _authHeaders tra ve dung headers
- [ ] TODO-[8] hoan thanh: pushObject upload thanh cong, co offline queue
- [ ] TODO-[9] hoan thanh: fetchObject download + cache thanh cong
- [ ] TODO-[10] hoan thanh: writeTree tao tree thanh cong
- [ ] TODO-[11] hoan thanh: commit tao commit + update HEAD
- [ ] TODO-[12] hoan thanh: log lay lich su commit
- [ ] TODO-[13] hoan thanh: sync hoat dong (enqueue + dequeue + execute)
- [ ] TODO-[14] hoan thanh: dispose() giai phong resources
- [ ] TODO-[15] hoan thanh: Data models (CommitInfo, TreeEntry, SyncStatus)
- [ ] TODO-[16] hoan thanh: ObjectCache bang SQLite
- [ ] TODO-[17] hoan thanh: SyncQueue bang SQLite
- [ ] TODO-[18] hoan thanh: Custom exceptions
- [ ] **Integration test:** Login -> Push object -> Fetch object roundtrip
- [ ] **Integration test:** Offline push -> enqueue -> online -> sync success
- [ ] **Integration test:** Cache hit (fetch 2 lan, lan 2 khong goi network)

---

### Test Checklist: Nhung gi ban nen tu viet test

```dart
// Test case: Login thanh cong tra ve token
// -> Mock HTTP server tra ve 200 + {"token": "abc123"}
//   Verify token duoc luu vao secure storage.

// Test case: Login sai password
// -> Mock 401. Verify AuthException duoc throw.

// Test case: Push object khi online
// -> Mock 201 + {"sha": "<hash>"}. Verify network call duoc thuc hien.

// Test case: Push object khi offline
// -> Mock no network. Verify object duoc enqueue, khong throw exception.
//   User van thay "pending" status.

// Test case: Fetch object cache miss
// -> Mock network call tra ve content. Verify content duoc cache.

// Test case: Fetch object cache hit
// -> Khong mock network. Content tra ve tu cache. Verify no HTTP call.

// Test case: Sync queue xu ly dung thu tu
// -> Enqueue 3 commands. Mock network success. 
//   Verify chung duoc execute theo dung thu tu.

// Test case: Sync voi 1 command failed
// -> Mock command 2 fail. Verify command 1 synced, command 2 failed,
//   command 3 van duoc xu ly (khong bi block boi command 2).

// Test case: Network timeout
// -> Mock request cham > 10s. Verify NetworkException duoc throw.
//   Khong block UI vo han.

// Test case: Dispose giai phong http client
// -> Goi dispose(). Verify khong the goi request sau do.

// Test case: Token het han -> auto refresh?
// -> (Optional) 401 tu server -> tu dong refresh token -> retry request.
//   Day la advanced feature, khong bat buoc cho capstone.
```

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off: REST vs WebSocket cho sync real-time?** Capstone nay dung REST cho sync (polling hoac manual). WebSocket cho real-time events ("co commit moi!"). WebSocket giu connection open = tiep pin hon. Ban se chon gi cho app can battery efficient vs app can real-time?

2. **Neu requirement thay doi: ho tro "partial sync" (chi sync 1 directory)?** Hien tai sync la "full sync". Nhung voi repo 10K files, full sync tren mobile cham. Partial sync: chi download objects trong 1 path. Server can endpoint moi: GET /api/v1/repos/{id}/tree?path=src/. Ban se design nhu the nao?

3. **Architecture decision: tai sao SDK khong expose HTTP client?** SDK cung cap high-level API (`pushObject`, `commit`) thay vi `post('/objects', ...)`. Dieu nay goi la "information hiding" — encapsulation. Lo ich: backend doi API endpoint -> chi doi SDK, khong doi app code. Chi phi: SDK team phai maintain. Worth it?

---


## Mini-Project: Minigit Full System

### User Story

> **Product Owner noi:** "Tich hop tat ca: CLI + Server + Flutter SDK. Day la capstone cuoi cung. Mot developer co the: khoi tao repo local, commit, push len server, roi dung Flutter app de xem repo va sync. He thong phai co architecture dung, test xanh, va documentation day du."
>
> **Context:** Day la "graduation project" cua khoa hoc 12 tuan. Moi thanh phan (CLI, Server, SDK) da duoc viet o cac topic truoc. Bay gio can ket hop chung thanh 1 he thong hoan chinh co the demo. Day la co hoi de chung minh ban hieu cach cac thanh phan tuong tac voi nhau.

### Acceptance Criteria

- [ ] **E2E Workflow 1: Local-only** — `init` -> `hash-object` -> `write-tree` -> `commit` -> `log` (chi dung CLI, khong can server)
- [ ] **E2E Workflow 2: Push len server** — `init` -> `commit` -> push objects len server (qua HTTP API) -> verify tren server
- [ ] **E2E Workflow 3: Flutter sync** — Login Flutter app -> tao repo -> sync objects -> xem file tree -> xem commit history
- [ ] **E2E Workflow 4: Multi-device** — CLI push file A -> Flutter app fetch file A -> xem noi dung giong nhau
- [ ] Architecture diagram co du 3 thanh phan va interactions
- [ ] All tests pass: `go test ./...` cho server, `flutter test` cho SDK
- [ ] Docker Compose chay toan bo he thong (server + PostgreSQL) bang 1 lenh
- [ ] README documentation cho tung thanh phan va integration
- [ ] Demo script: step-by-step commands de chay E2E workflow

---

### Senior Thought-Process

```markdown
**Senior nghi gi khi nhan capstone nay:**
> "Capstone la luc moi thu ket noi. Hoi toi tot nghiep bootcamp, 
> toi cung lam project tuong tu — va toi hoc duoc nhieu hon ca 
> 11 tuan truoc cong lai. Vi luc nay, ban khong hoc syntax nua, 
> ban hoc cach ket noi cac thanh phan thanh he thong.
>
> Quy trinh tu duy cua toi khi integrate:
> 1. Xac dinh 'integration points' — noi nao 2 thanh phan noi chuyen?
>    CLI <-> Server: HTTP API
>    Server <-> DB: SQL queries  
>    SDK <-> Server: HTTP API (giong CLI nhung qua Dart)
>
> 2. Bat dau tu 'happy path' don gian nhat. Dung co gang 
>    handle moi edge case ngay. Lam duoc 1 luong E2E roi moi polish.
>
> 3. Dockerize som. Neu khong dockerize, ban se mat 2 tieng 
>    moi lan setup DB, migration, env vars. Docker Compose 
>    chay tat ca trong 30 giay.
>
> 4. Viet integration tests. Khong phai unit test — integration test. 
>    1 test chay ca server + DB, verify API tra ve dung.
>
> Vấn de cot loi: system thinking. Khong chi code duoc 
> function tot — ma phai hieu cach function do anh huong 
> den toan bo he thong."
```

---

### TODO Comments (Code Skeleton)

```go
// ======= E2E Test: Workflow 1 (Local CLI) =======
// TODO-[1]: Test script cho local CLI workflow
// SENIOR ASKS: Test E2E nen viet nhu the nao?
// HINT: Dung Go test voi temp directory. Moi test step:
//   1. Tao temp dir
//   2. Cd vao temp dir
//   3. Chay tung command (go run . init, go run . hash-object, ...)
//   4. Verify output (hash dung, file ton tai, ...)
//   5. Cleanup temp dir

func TestE2E_LocalWorkflow(t *testing.T) {
	// TODO: Implement
	// STEP 1: t.TempDir() de tao temp directory
	// STEP 2: Chay CmdInit
	// STEP 3: Ghi file "hello.txt" voi noi dung "hello world"
	// STEP 4: Chay CmdHashObject voi file
	// STEP 5: Verify hash la SHA1 dung
	// STEP 6: Chay CmdWriteTree
	// STEP 7: Verify tree hash
	// STEP 8: Chay CmdCommit voi tree hash
	// STEP 9: Verify commit co parent = nil (root commit)
	// STEP 10: Chay CmdLog, verify co 1 commit
}

// ======= E2E Test: Workflow 2 (CLI + Server) =======
// TODO-[2]: Test script cho push-to-server workflow
// SENIOR ASKS: Integration test voi server can gi?
// HINT: httptest.Server de mock server, HOAC chay server that trong goroutine.
//   Voi capstone: dung httptest.Server de test HTTP interactions.
//   Voi integration: chay server that + PostgreSQL (testcontainers hoac local).

func TestE2E_PushToServer(t *testing.T) {
	// TODO: Implement
	// STEP 1: Start test server (httptest.NewServer)
	// STEP 2: Register user -> login -> lay token
	// STEP 3: Tao repo -> lay repo ID
	// STEP 4: Tinh SHA1 cua "hello" local
	// STEP 5: POST /api/v1/repos/{id}/objects voi content
	// STEP 6: Verify response 201 + SHA dung
	// STEP 7: GET /api/v1/repos/{id}/objects/{sha}
	// STEP 8: Verify content goc = "hello"
}

// ======= E2E Test: Workflow 3 (Flutter SDK + Server) =======
// TODO-[3]: Test script cho Flutter SDK workflow
// SENIOR ASKS: Test Dart code nhu the nao?
// HINT: flutter test hoac dart test. Mock HTTP bang mockito hoac http.MockClient.
//   Test: login -> push object -> fetch object -> verify roundtrip.

// ======= docker-compose.yml =======
// TODO-[4]: Docker Compose cho toan bo he thong
// SENIOR ASKS: Services nao can? Config nhu the nao?
// HINT: 2 services: postgres va server.
//   postgres: image: postgres:16, env POSTGRES_USER/PASSWORD/DB
//   server: build: ./server/, ports: "8080:8080", depends_on: postgres
//   environment: DATABASE_URL=postgresql://user:pass@postgres/db

// version: '3.8'
// services:
//   postgres:
//     image: postgres:16-alpine
//     environment:
//       POSTGRES_USER: minigit
//       POSTGRES_PASSWORD: minigit
//       POSTGRES_DB: minigit
//     ports:
//       - "5432:5432"
//     volumes:
//       - postgres_data:/var/lib/postgresql/data
//   
//   server:
//     build: ./cmd/server/
//     ports:
//       - "8080:8080"
//     environment:
//       PORT: 8080
//       DATABASE_URL: postgres://minigit:minigit@postgres:5432/minigit?sslmode=disable
//       JWT_SECRET: capstone-secret-key-change-in-production
//     depends_on:
//       - postgres
// 
// volumes:
//   postgres_data:

// ======= Makefile =======
// TODO-[5]: Makefile cho common commands
// SENIOR ASKS: Nhung commands nao thuong dung?
// HINT: make dev (docker compose up), make test (go test), 
//   make migrate (chay migration), make seed (insert test data),
//   make build (build binary), make clean.

// .PHONY: dev test migrate build clean
// dev:
// 	docker-compose up --build
// test:
// 	go test ./... -race -count=1
// migrate:
// 	go run ./cmd/migrate up
// build:
// 	go build -o bin/minigit-server ./cmd/server/
// 	go build -o bin/minigit-cli ./cmd/cli/
// clean:
// 	rm -rf bin/
// 	docker-compose down -v

// ======= README.md =======
// TODO-[6]: README cho toan bo project
// SENIOR ASKS: README nen co nhung phan nao?
// HINT: 
//   # Minigit
//   ## Overview (1 doan mo ta)
//   ## Architecture (diagram hoac mo ta)
//   ## Quick Start (docker compose up -> test commands)
//   ## API Endpoints (list + example curl)
//   ## Project Structure (tree cac package)
//   ## Testing (go test, flutter test)
//   ## Tech Stack (list libraries)
//   ## Known Limitations (staging area, packfile, LFS)
//   ## Future Work (WebSocket real-time, gRPC, etc.)
```

---

### Socratic Questions

**Cau hoi de ban tu suy nghi:**

1. **Khi nao integration test, khi nao unit test?** Unit test: nhanh, isolate, test 1 function. Integration test: cham, test he thong cong tac. Rule of thumb: business logic quan trong = unit test, API contract = integration test. Voi capstone nay, ban se viet bao nhieu integration test?

2. **Tai sao Docker Compose la "bat buoc" cho capstone?** Neu khong dung Docker, moi lan chay integration test ban phai: start PostgreSQL manually, tao database, chay migration, set env vars, nho stop sau khi xong. Docker Compose lam dieu nay trong 30 giay va cleanup tu dong. "Works on my machine" khong du — phai "works on any machine."

3. **Test data (fixtures) nen setup nhu the nao?** Moi test chay xong phai de lai DB trong trang thai sach (hoac moi test chay tren DB rieng). Cach pho bien: 1) Transaction rollback sau moi test, 2) TRUNCATE tables truoc moi test, 3) Docker container moi cho moi test suite. Ban chon cach nao? Tai sao?

4. **Makefile co con can thiet trong thoi dai CI/CD?** Make rat tot cho local dev: 1 lenh duy nhat `make test` thay vi nho 5 buoc. Nhung trong CI/CD (GitHub Actions), script YAML cung lam duoc. Makefile la "developer experience" — no khong phai requirement, nhung khong co no thi team se mat thoi gian.

5. **README quan trong den muc nao?** Hoi toi nhan 1 project khong co README, toi mat 2 ngay de hieu cach chay. Voi README tot: 30 phut. README la "face" cua project — recruiter se doc no truoc khi doc code. Ban viet README cho nguoi chua biet gi ve project.

---

### Output Checklist: Lam sao biet minh xong?

- [ ] TODO-[1] hoan thanh: E2E test local CLI workflow pass
- [ ] TODO-[2] hoan thanh: E2E test push-to-server workflow pass
- [ ] TODO-[3] hoan thanh: E2E test Flutter SDK workflow pass
- [ ] TODO-[4] hoan thanh: docker-compose.yml chay duoc (`docker-compose up` -> server hoat dong)
- [ ] TODO-[5] hoan thanh: Makefile co dev, test, migrate, build, clean
- [ ] TODO-[6] hoan thanh: README day du (overview, architecture, quick start, API, structure, testing, stack, limitations)
- [ ] **Integration test:** `docker-compose up` -> curl health -> 200 OK
- [ ] **Integration test:** Full CLI -> Server -> SDK roundtrip
- [ ] **Code quality:** `go test -race ./...` pass (zero races)
- [ ] **Code quality:** `go vet ./...` pass
- [ ] **Code quality:** `gofmt` applied
- [ ] **Demo:** Co the demo 5 phut voi 1 terminal va 1 Flutter app

---

### Test Checklist: Nhung gi ban nen tu viet test

```go
// Test case: Full E2E — init local, commit, push server, fetch SDK
// -> Day la "holy grail" test. Chay ca 3 thanh phan.
//    Kiem tra data integrity: content khong bi corrupt qua network.

// Test case: Concurrent uploads from 2 clients
// -> 2 goroutine push cung luc. Server phai xu ly dung, 
//    khong bi race condition. go test -race kiem tra.

// Test case: Server restart giua upload va fetch
// -> Upload thanh cong -> restart server -> fetch van thanh cong.
//    Verify persistence hoat dong.

// Test case: Large object (10MB) upload + download
// -> Verify khong bi memory issue. Streaming hoat dong.

// Test case: Unauthenticated access to protected endpoint
// -> Tat ca protected endpoints phai tra 401 khi khong co token.
//    Test automation: lap qua tat ca routes, goi khong token.

// Test case: Invalid JWT token
// -> Token sai format, token het han, token khong dung signature.
//    Tat ca phai tra 401.

// Test case: Docker Compose startup
// -> `docker-compose up --build` thanh cong trong < 60s.
//    Health check tra 200. Cac service noi chuyen duoc.

// Test case: Graceful shutdown du requests dang xu ly
// -> Goi SIGTERM trong luc upload -> upload van hoan thanh.
```

---

### Retrospective: Sau khi xong, hay tu hoi

1. **Trade-off: monorepo vs polyrepo?** Capstone nay dung monorepo (1 repo chua ca CLI + Server + SDK). Uber dung monorepo de share code. Microservice team thuong dung polyrepo. Ban thay loi/lo gi cua moi cach? Voi team 3 nguoi, chon gi?

2. **Neu requirement thay doi: ho tro 1000 users concurrent?** Hien tai server dung 1 instance. De scale: horizontal scaling (nhieu instances + load balancer), DB read replicas, Redis cache. Ban se identify bottleneck dau tien? (DB connection pool? Memory? CPU?)

3. **Architecture decision: tai sao khong dung microservices?** Capstone nay la monolith (1 server). Microservices: scale independently, independent deploy, nhung distributed system complexity (network latency, partial failures, data consistency). Sau 12 tuan hoc, ban thay microservices co phu hop voi junior khong?

---


---

## Architecture Overview

> "Architecture la nhung quyet dinh kho thay doi sau nay. Chon dung luc dau, ban se tiet kiem 100 gio refactoring sau nay. Chon sai, ban se danh 2 tuan sua lai."

### System Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MINIGIT CAPSTONE SYSTEM                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐      HTTP/REST      ┌─────────────────────────────────┐   │
│  │              │  ┌──────────────┐   │         minigit-server           │   │
│  │  minigit-cli │──┤  REST API    │───┤  ┌──────────┐  ┌──────────────┐ │   │
│  │  (Go binary) │  │  /api/v1/*   │   │  │  Router  │  │   Services   │ │   │
│  │              │  └──────────────┘   │  │  (chi)   │  │   (business) │ │   │
│  │ • init       │                     │  └────┬─────┘  └──────┬───────┘ │   │
│  │ • hash-object│                     │       │               │         │   │
│  │ • cat-file   │      WS (optional)  │  ┌────▼─────┐  ┌──────▼───────┐ │   │
│  │ • write-tree │  ┌──────────────┐   │  │Middleware│  │ Repositories │ │   │
│  │ • commit     │──┤  WebSocket   │───┤  │• JWT    │  │  • User      │ │   │
│  │ • log        │  │  /events     │   │  │• Logging│  │  • Repo      │ │   │
│  └──────────────┘  └──────────────┘   │  │• Metrics│  │  • Object    │ │   │
│                                        │  └────┬─────┘  │  • Ref       │ │   │
│  ┌──────────────┐                     │       │          └──────┬───────┘ │   │
│  │              │      HTTP/REST      │  ┌────▼─────────────────▼──────┐  │   │
│  │ minigit-sdk  │  ┌──────────────┐   │  │       PostgreSQL           │  │   │
│  │ (Dart pkg)   │──┤  REST API    │───┤  │  users, repos, objects,    │  │   │
│  │              │  │  + WS client │   │  │  refs tables               │  │   │
│  │ • login      │  └──────────────┘   │  └──────────────────────────┘  │   │
│  │ • pushObject│                     │                                 │   │
│  │ • fetchObject│                    │  ┌──────────────────────────┐   │   │
│  │ • sync       │                     │  │     Observability       │   │   │
│  └──────────────┘                     │  │ • slog (structured)     │   │   │
│                                       │  │ • OpenTelemetry traces  │   │   │
│         Flutter App                   │  │ • Prometheus metrics    │   │   │
│    ┌──────────────┐                   │  └──────────────────────────┘   │   │
│    │  UI Screens  │                   │                                 │   │
│    │ • Repo list  │                   └─────────────────────────────────┘   │
│    │ • File tree│                                                           │
│    │ • Commit log│                                                          │
│    │ • Sync status│                                                         │
│    └──────────────┘                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Component Details

#### 1. minigit-cli

| Aspect | Detail |
|--------|--------|
| **Language** | Go 1.21+ |
| **Storage** | FileSystem (`.minigit/objects/`) + SQLite (index, optional) |
| **Parallelism** | Worker pool cho blob hashing (goroutines + channels) |
| **Entry point** | `cmd/minigit/main.go` |
| **Key packages** | `object`, `store`, `command` |
| **Binary size** | ~5-10MB (static linked) |

```markdown
**Senior ghi chu ve CLI:**
> "CLI la thanh phan don gian nhat nhung day la noi logic Git thuc su song. 
> Server chi la 'thin wrapper' quanh object model. Khi debug issue, 
> toi luon bat dau tu CLI de xem object model co dung khong.
>
> Quyet dinh kien truc quan trong:
> - Storage: FileSystem cho objects (mimic .git/objects/)
> - No SQLite cho capstone co ban (chi can file storage)
> - Worker pool: 4-8 workers cho hashing parallel, bounded concurrency
> - Khong can caching — binary chay 1 lan roi exit"
```

#### 2. minigit-server

| Aspect | Detail |
|--------|--------|
| **Language** | Go 1.21+ |
| **Router** | chi/v5 (lightweight, stdlib-compatible) |
| **Database** | PostgreSQL 16 (pgx/v5 pool) |
| **Auth** | JWT (golang-jwt/jwt/v5), bcrypt password hashing |
| **Logging** | log/slog (JSON handler) |
| **Tracing** | OpenTelemetry (otlptracehttp) |
| **Metrics** | Prometheus (promhttp.Handler) |
| **Graceful shutdown** | 30s timeout |
| **Entry point** | `cmd/server/main.go` |
| **Key packages** | `handler`, `service`, `repository`, `middleware`, `auth` |

```markdown
**Senior ghi chu ve Server:**
> "Server la thanh phan production. Moi quyet dinh o day deu co hau qua 
> khi chay that. Toi da tung debug server crash luc 3AM vi goroutine 
> leak — tu do toi luon check `go test -race` va viet graceful shutdown.
>
> Quyet dinh kien truc quan trong:
> - chi over stdlib: routing + middleware chain ro rang
> - pgx over database/sql: native PostgreSQL, type-rich
> - Repository pattern: de test, de swap DB
> - Service layer: business logic khong nam trong handler
> - JWT + bcrypt: auth standard, khong tu viet crypto
> - slog JSON: structured logs cho log aggregation
> - OTel: distributed tracing de debug trong microservices"
```

#### 3. minigit-sdk

| Aspect | Detail |
|--------|--------|
| **Language** | Dart 3+ |
| **HTTP** | `package:http` |
| **Secure storage** | flutter_secure_storage |
| **Local cache** | SQLite (sqflite) hoac Hive |
| **Offline queue** | SQLite table |
| **Serialization** | dart:convert (JSON) |
| **Entry point** | `lib/minigit_sdk.dart` |
| **Key files** | `client.dart`, `cache.dart`, `sync_queue.dart`, `models.dart` |

```markdown
**Senior ghi chu ve SDK:**
> "SDK la 'face' cua he thong doi voi Flutter developer. API phai clean, 
> consistent, va handle duoc offline. Toi da tung dung SDK ma 
> throw generic Exception — khong biet la network loi hay auth loi. 
> Rat kho debug. Custom exceptions la bat buoc.
>
> Quyet dinh kien truc quan trong:
> - Instance pattern (khong singleton): de test, de mock
> - flutter_secure_storage: JWT khong duoc plain text
> - Offline queue: commands luu SQLite, sync khi co mang
> - Cache: content-addressable = immutable = khong can invalidate
> - dispose(): giai phong http.Client de tranh leak"
```

### Data Flow (E2E)

```
1. User viet file "hello.txt" trong working directory
   |
2. minigit-cli write-tree: quet directory -> tao blobs -> tao tree
   |
3. minigit-cli commit: tao commit -> update HEAD -> in commit SHA
   |
4. minigit-cli push: POST /api/v1/repos/{id}/objects (tung object)
   |
5. Server nhan request -> JWT verify -> insert DB (ON CONFLICT skip)
   |
6. User mo Flutter app -> login -> sync -> fetch objects
   |
7. minigit-sdk: GET /api/v1/repos/{id}/objects/{sha} -> cache local
   |
8. Flutter UI hien thi file tree tu tree object + blob content
```

---

## API Contract & DB Schema

### API Endpoints

#### Authentication (Public)

```http
### Register
POST /api/v1/register
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password_123"
}

### Response 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "created_at": "2026-06-01T12:00:00Z"
}

### Response 409 Conflict
{
  "error": "username already exists"
}

### Response 400 Bad Request
{
  "error": "password must be at least 8 characters"
}
```

```http
### Login
POST /api/v1/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "secure_password_123"
}

### Response 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}

### Response 401 Unauthorized
{
  "error": "invalid credentials"
}
```

#### Repository (Authenticated)

```http
### Create Repo
POST /api/v1/repos
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "name": "my-project"
}

### Response 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "my-project",
  "owner_id": "550e8400-e29b-41d4-a716-446655440000",
  "default_branch": "main",
  "created_at": "2026-06-01T12:00:00Z"
}

### Response 409 Conflict
{
  "error": "repo name already exists for this user"
}
```

```http
### List Repos
GET /api/v1/repos
Authorization: Bearer <jwt>

### Response 200 OK
{
  "repos": [
    {
      "id": "...",
      "name": "my-project",
      "default_branch": "main",
      "created_at": "2026-06-01T12:00:00Z"
    }
  ]
}
```

#### Objects (Authenticated)

```http
### Upload Object
POST /api/v1/repos/{repo_id}/objects
Authorization: Bearer <jwt>
Content-Type: application/octet-stream

<raw bytes: blob/tree/commit content>

### Response 201 Created
{
  "sha": "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
  "type": "blob",
  "size": 0
}

### Response 409 Conflict (duplicate — idempotent)
{
  "sha": "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
  "message": "object already exists"
}
```

```http
### Fetch Object
GET /api/v1/repos/{repo_id}/objects/{sha}
Authorization: Bearer <jwt>

### Response 200 OK
Content-Type: application/octet-stream

<raw bytes>

### Response 404 Not Found
{
  "error": "object not found"
}
```

#### Refs (Authenticated)

```http
### Get Ref
GET /api/v1/repos/{repo_id}/refs/{name}
Authorization: Bearer <jwt>

### Response 200 OK
{
  "name": "main",
  "commit_sha": "aabbccdd11223344556677889900aabbccddeeff"
}

### Response 404 Not Found
{
  "error": "ref not found"
}
```

```http
### Update Ref
POST /api/v1/repos/{repo_id}/refs/{name}
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "commit_sha": "aabbccdd11223344556677889900aabbccddeeff"
}

### Response 200 OK
{
  "name": "main",
  "commit_sha": "aabbccdd11223344556677889900aabbccddeeff"
}
```

#### Health (Public)

```http
### Health Check
GET /api/v1/health

### Response 200 OK
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2026-06-01T12:00:00Z"
}
```

### DB Schema (PostgreSQL)

```sql
-- ======= Users =======
-- SENIOR ASKS: Tai sao id dung UUID thay vi serial/bigint?
-- HINT: UUID khong expose thu tu tao account, de shard cross-region,
--   va khong bi "enumerate" attack (khong doan duoc user #1, #2, #3).

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

-- ======= Repos =======
-- SENIOR ASKS: Moi user co nhieu repos. Ownership quan trong.
-- HINT: owner_id FK den users. name unique per user (chua handle o DB level).

CREATE TABLE repos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    default_branch VARCHAR(50) DEFAULT 'main',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(owner_id, name)
);

CREATE INDEX idx_repos_owner ON repos(owner_id);

-- ======= Objects =======
-- SENIOR ASKS: Tai sao (repo_id, sha) la composite PK?
-- HINT: SHA co the trung nhau giua repos (cung content). 
--   Nhung chung ta khong muon 1 repo co 2 object cung SHA.
--   content: BYTEA luu compressed object content. 
--   Voi object > 10MB, can reconsider (S3 thay BYTEA).

CREATE TABLE objects (
    repo_id UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    sha CHAR(40) NOT NULL,  -- SHA1 hex string = 40 chars
    type VARCHAR(10) NOT NULL CHECK (type IN ('blob', 'tree', 'commit')),
    size INTEGER NOT NULL,  -- uncompressed size in bytes
    content BYTEA NOT NULL, -- zlib-compressed content
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (repo_id, sha)
);

CREATE INDEX idx_objects_repo ON objects(repo_id);
CREATE INDEX idx_objects_sha ON objects(sha);

-- ======= Refs =======
-- SENIOR ASKS: Tai sao (repo_id, name) la composite PK?
-- HINT: 1 repo co nhieu refs (main, develop, feature/x). 
--   Ten ref chi unique trong 1 repo.

CREATE TABLE refs (
    repo_id UUID NOT NULL REFERENCES repos(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    commit_sha CHAR(40) NOT NULL REFERENCES objects(repo_id, sha),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (repo_id, name)
);

-- ======= Sync Queue (cho Flutter SDK offline support) =======
-- SENIOR ASKS: Table nay co o server hay chi o client?
-- HINT: Chi o client (SQLite tren device). Khong can tren server.
--   Nhung neu muon "pending changes" view cross-device, can server table.
--   Capstone: chi client-side.

-- CLIENT-SIDE ONLY (SQLite):
-- CREATE TABLE sync_queue (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     repo_id TEXT NOT NULL,
--     command_type TEXT NOT NULL CHECK (command_type IN ('push_object', 'write_tree', 'commit', 'update_ref')),
--     payload TEXT NOT NULL, -- JSON
--     status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'syncing', 'synced', 'failed')),
--     error_message TEXT,
--     created_at INTEGER DEFAULT (unixepoch()),
--     updated_at INTEGER DEFAULT (unixepoch())
-- );
```

---

### Socratic Questions (Architecture & API)

**Cau hoi de ban tu suy nghi:**

1. **Tai sao API dung `application/octet-stream` cho object thay vi JSON base64?** Octet-stream = raw bytes, khong encoding overhead. JSON base64 = +33% size (1MB -> 1.33MB) + CPU encode/decode. Git objects la binary — khong phai text. Raw bytes la efficient nhat. Ban co thay trade-off gi? (Khong de doc khi debug? Dung hex dump neu can.)

2. **Tai sao objects table dung composite PK (repo_id, sha) thay vi 1 id UUID rieng?** Git philosophy: object duoc identify by content hash, khong can synthetic ID. Composite PK = 1 lookup query (SELECT ... WHERE repo_id=$1 AND sha=$2). UUID ID = 2 lookups (id -> sha, sha -> content). Don gian hon, nhanh hon.

3. **Tai sao ref.commit_sha FK den objects(repo_id, sha)?** Referential integrity: khong the tao ref tro den commit khong ton tai. PostgreSQL dam bao dieu nay. Neu xoa commit, ref cung bi xoa (ON DELETE CASCADE). Ban co muon CASCADE khong? (Trong Git real, ref co the "dangling" — tro den object da bi garbage collect.)

4. **`repos` co `default_branch` column — tai sao?** Khi tao repo moi, HEAD can biet branch mac dinh la gi. Git dung "master" hoac "main". Column nay cho phep customize. No cung dong vai tro khi Flutter app "clone" repo — biet ref nao de fetch dau tien.

5. **Health check endpoint nen tra thong tin gi?** Chi "status: healthy" hay chi tiet hon (DB connection, disk space, memory)? Detailed health check = nhieu thong tin de debug, nhung cung expose internals. Ban se chon gi? (Chi tiet cho internal monitoring, don gian cho public.)

---

### Output Checklist: Architecture & API

- [ ] System diagram du 3 thanh phan + interactions
- [ ] API contract day du: request/response shape cho tat ca endpoints
- [ ] Status codes chuan: 200, 201, 400, 401, 404, 409
- [ ] DB schema: 4 tables (users, repos, objects, refs) + indexes + constraints
- [ ] Composite PK design dung (repo_id, sha) va (repo_id, name)
- [ ] Referential integrity: FK constraints + ON DELETE CASCADE
- [ ] Content-Type: application/octet-stream cho objects
- [ ] Error response format consistent (JSON voi "error" field)

---

## Decision Heuristics

> "Khong co quyet dinh dung/sai tuyet doi — chi co quyet dinh phu hop voi context. Day la nhung 'rules of thumb' toi dung sau 12 nam."

### Heuristic Table

| Tinh huong | Quyet dinh | Ly do |
|---|---|---|
| Co the viet bang stdlib trong <200 dong? | **Dung stdlib** | It dependency = it van de bao tri, it security audit, nho binary |
| Can routing + middleware? | **chi hoac httprouter** | Nhe, compatible stdlib, khong magic. Tranh Gin/Echo khi hoc |
| Can struct scanning tu SQL? | **sqlx** | Giam boilerplate 50% ma van viet raw SQL. Khong dung ORM khi chua biet SQL |
| Public API cho browsers/clients? | **REST + JSON** | Universal, de debug, khong can proto compiler |
| Internal service communication? | **gRPC** | Nhanh hon, type-safe, streaming. Dung REST gateway cho external |
| Luu password? | **bcrypt, never roll your own** | Crypto la chuyen gia lam. bcrypt da duoc audit. Argon2 neu can memory-hard |
| Real-time den Flutter? | **WebSocket hoac SSE** | WebSocket: 2-way, stateful. SSE: 1-way, simpler, auto-reconnect. Chon SSE neu chi push |
| Error tu external lib? | **Wrap voi `%w`, log structured** | `fmt.Errorf("query failed: %w", err)` de dung `errors.Is`. Slog cho machine-parseable |
| 3rd party dependency? | **Check license + maintenance + vendor neu critical** | License incompatible = phap ly risk. Unmaintained = security risk. Critical dep = vendor (copy vao repo) |
| Config cho app? | **Env vars, fail fast** | 12-factor app. Khoi dong = validate config, thieu = log.Fatal. Khong dung default "secret" |
| Test: unit hay integration? | **Unit cho logic, integration cho API contract** | Unit: nhanh, nhieu case. Integration: dam bao components noi chuyen duoc |
| Goroutine cho parallel work? | **Worker pool voi bounded concurrency** | Khong `go func()` vo han — goroutine leak. Buffered channel + fixed workers |
| Logging trong production? | **slog JSON, khong fmt.Printf** | Structured logs = query duoc ("level=error AND component=auth"). fmt.Printf = parse bang regex |
| Database timeout? | **Context voi timeout moi query** | `ctx, cancel := context.WithTimeout(5s)`. Khong de query chay vo han |
| Deploy container? | **Multi-stage Dockerfile** | Stage 1: build. Stage 2: runtime (distroless/alpine). Final image < 20MB |
| Graceful shutdown? | **Signal handling + drain timeout** | SIGTERM -> stop nhan requests -> cho in-flight hoan thanh -> exit. Khong drop requests |

### Quy trinh ra quyet dinh (Decision Flowchart)

```
Bat dau: Co van de can giai quyet
    |
    v
Can viet tu dau hay co san? 
    |
    +-- Co thu vien tot, license OK? --> Dung thu vien
    |
    +-- Co the viet stdlib <200 dong? --> Viet stdlib
    |
    +-- Phuc tap, can maintain? --> Chon thu vien community (check stars/issues)
    |
    v
Can performance cao? 
    |
    +-- Latency critical? --> Benchmark truoc khi quyet dinh
    |
    +-- Thong thuong? --> Chon de hieu, de maintain
    |
    v
Production hay prototype?
    |
    +-- Prototype/MVP? --> Chon nhanh nhat, refactor sau
    |
    +-- Production? --> Chon bao mat, observe, maintain
    |
    v
Ghi lai quyet dinh: ADR (Architecture Decision Record)
    |
    v
Done: Implement + test + document
```

### Senior Stories: Nhung quyet dinh "dau don"

```markdown
**Story 1: "Toi da chon ORM vi luoi viet SQL"**
> Nam thu 3 di lam, toi chon GORM cho 1 microservice. 6 thang sau, 
> 1 query sinh ra 47 JOIN lam DB CPU 100%. Toi khong biet SQL nao 
> duoc sinh ra vi GORM abstraction an no. Mat 2 ngay rewrite bang 
> raw SQL. Bai hoc: hoc SQL truoc khi dung ORM. sqlx la sweet spot.

**Story 2: "Goroutine leak lam server crash"**  
> Toi viet `go process(item)` trong 1 loop khong gioi han. 2 tieng 
> sau, server OOM vi 2 trieu goroutine. Fix: worker pool voi 
> channel co buffer. Bai hoc: moi goroutine phai co "exit strategy".

**Story 3: "fmt.Printf trong production"**
> Toi debug bang fmt.Printf va quen xoa. Production log 500GB/ngay 
> vi log moi request detail. Logging system cost $2000/thang. 
> Fix: log/slog voi level control, INFO cho production, DEBUG chi 
> khi can. Bai hoc: khong bao gio fmt.Printf trong production.

**Story 4: "JWT secret trong source code"**
> 1 junior commit JWT secret vao GitHub. Bot scan tim thay trong 
> 30 phut. Phai rotate secret, logout all users, write incident report. 
> Bai hoc: secret trong env var, khong bao gio trong code. 
> Dung .env.example (khong co gia tri that) cho local dev.
```

---

## Tong ket: Tu Zero den Capstone

### Nhung gi ban da hoc qua 12 tuan

| Phase | Noi dung chinh | Artefact |
|-------|---------------|----------|
| **Phase 1** | Syntax, type system, control flow | CLI toolkit (convert, stats, inspect) |
| **Phase 2** | Goroutines, channels, sync | Concurrent log analyzer |
| **Phase 3** | Stdlib: HTTP, JSON, DB, testing | REST API with SQLite |
| **Phase 4** | Generics, type safety | Generic cache, functional helpers |
| **Phase 5** | Production: logging, graceful shutdown, OTel | Production-hardened API |
| **Phase 6** | Ecosystem: chi, pgx, gRPC, JWT, testing | Refactored API + contract |
| **Phase 7** | Capstone: tich hop tat ca | **Minigit: CLI + Server + Flutter SDK** |

### Nhung nguyen tac "bat bien"

```markdown
1. **Hieu truoc, viet sau.** Khong viet code khi chua hieu problem.
2. **Test song hanh.** Viet test cung luc hoac truoc code (TDD).
3. **Fail fast.** Validate input, config, preconditions ngay — khong de loi lan xa.
4. **Explicit over implicit.** Go khong co implicit. Error phai handle. Type phai convert.
5. **Simplicity matters.** Code don gian hon > code "thong minh".
6. **No magic.** Hieu moi dong code ban viet. Khong copy-paste ma khong hieu.
7. **Observability.** Production khong co printf. Structured logs, traces, metrics.
8. **Security by default.** Khong bao gio roll your own crypto. Bcrypt, JWT standard.
9. **Graceful everything.** Shutdown, errors, degradation — tat ca phai graceful.
10. **Context everywhere.** Moi goroutine, moi query, moi request phai co context voi timeout.
```

### Con duong phia truoc

Capstone ket thuc, nhung hanh trinh hoc Go moi bat dau. Nhung huong di tiep theo:

- **Distributed Systems:** NATS, Redis, microservices patterns
- **Advanced Go:** unsafe, cgo, WASM (tinygo)
- **DevOps:** Kubernetes operator, GitHub Actions, CI/CD pipeline
- **Open Source:** Contribute to Go ecosystem, viet blog, mentor nguoi moi
- **Specialization:** Database engine, networking, security, cloud-native

> "The best time to plant a tree was 20 years ago. The second best time is now."
> 
> 12 tuan qua, ban da trong 1 cay. Hay cham soc no — code moi ngay, doc source code cac project lon, va dung so sai. Senior cung tung la junior, cung tung viet code xau, cung tung debug den 3AM. Dieu khac biet la ho khong bo cuoc.

---

*File nay duoc viet theo format REQUIREMENT-SIMULATION: khong co code hoan chinh, chi co skeleton, hints, va cau hoi de tu suy nghi. Moi TODO can duoc implement boi nguoi hoc — do la cach hoc tot nhat.*

*Staff Engineer / Go Specialist — 12+ years shipping production systems.*
