# Plan: Go Senior Learning Guide - Requirement-Simulation Style

## Paradigm Shift
Từ: Sách hướng dẫn đầy code mẫu
Sang: Mô phỏng môi trường làm việc với senior mentor (Socratic)

## Workflow Mỗi Topic
```
a) User Story: Khách hàng đưa yêu cầu (requirement simulation)
b) Acceptance Criteria: Điều kiện chấp nhận
   
c) Senior Guide: (THOUGHT-PROCESS + CÂU HỎI GỢI Ý)
   - Senior Thought-Process: "Nếu tôi nhận requirement này, tôi nghĩ..."
   - TODO Comments: Code skeleton với TODO comments
   - Socratic Questions: Chỉ hỏi, không trả lờI
   - Hints: Gợi ý từng bước nhỏ nếu stuck
   
d) Checklist Output: Xác nhận hoàn thành phần nào

e) Test Checklist: AI gợi ý các test scenarios, user tự viết test
   - "Nên test trường hợp này..."
   - "Boundary case này có thể fail..."
   
f) Retrospective: "Tại sao chọn cách này?"
   - Trade-off analysis
   - "Nếu requirement thay đổi thì sao?"
```

## Format Code
```go
// ┌─────────────────────────────────────────────────────────────┐
// │ USER STORY: Khách hàng muốn...                              │
// │ ACCEPTANCE CRITERIA:                                         │
// │   1. ...                                                     │
// │   2. ...                                                     │
// └─────────────────────────────────────────────────────────────┘

// ┌─────────────────────────────────────────────────────────────┐
// │ SENIOR THOUGHT-PROCESS:                                      │
// │ "Nếu tôi nhận requirement này, tôi sẽ nghĩ về..."           │
// │ "Vấn đề chính ở đây là..."                                   │
// │ "Cần phân rã thành các phần..."                              │
// └─────────────────────────────────────────────────────────────┘

package main

// TODO: Hãy nghĩ xem, package main cần import gì?
// HINT: Bạn cần đọc input từ command line... đâu là package phù hợp?

// TODO: Định nghĩa hàm main - đây là entry point
// SENIOR ASKS: Một CLI tool thường cần xử lý gì đầu tiên?

func main() {
    // TODO: Parse command-line arguments
    // SENIOR ASKS: os.Args hay flag package? Khi nào dùng cái nào?
    // HINT: Nếu chỉ có 1-2 argument đơn giản... nếu phức tạp hơn...
    
    // TODO: Validate input
    // SENIOR ASKS: Khách hàng nói "xử lý lỗi input" - bạn nghĩ cần validate những gì?
    
    // TODO: Business logic
    // SENIOR ASKS: Công thức chuyển đổi là gì? Bạn có thể viết riêng thành hàm không?
    
    // TODO: Output result
    // SENIOR ASKS: fmt.Printf với format nào để đẹp?
}

// ┌─────────────────────────────────────────────────────────────┐
// │ OUTPUT CHECKLIST: Khi nào phần này hoàn thành?               │
// │ [ ] TODO 1: Parse arguments — có thể đọc được input         │
// │ [ ] TODO 2: Validate — chương trình không panic với input   │
// │ [ ] TODO 3: Convert — kết quả đúng với test cases           │
// │ [ ] TODO 4: Output — format đẹp, dễ đọc                     │
// └─────────────────────────────────────────────────────────────┘

// ┌─────────────────────────────────────────────────────────────┐
// │ TEST CHECKLIST: Những gì bạn nên tự viết test để kiểm tra    │
// │ 1. Test happy path: 100°C → 212°F                           │
// │ 2. Test negative: -40°C → -40°F                             │
// │ 3. Test zero: 0°C → 32°F                                    │
// │ 4. Test invalid input: "abc" → error, không panic           │
// │ 5. Test missing argument: Không có input → error msg rõ ràng │
// └─────────────────────────────────────────────────────────────┘

// ┌─────────────────────────────────────────────────────────────┐
// │ RETROSPECTIVE:                                               │
// │ 1. Tại sao bạn chọn os.Args hay flag? Trade-off là gì?       │
// │ 2. Nếu khách hàng thêm 10 đơn vị nữa, code bạn có dễ mở    │
// │    rộng không?                                               │
// │ 3. Bạn validate ở đâu: trong main hay tách ra hàm riêng?     │
// │    Tại sao?                                                  │
// └─────────────────────────────────────────────────────────────┘
```

## Output: 7 file MD riêng biệt
1. phase-01-foundations.md — Tuần 1-3
2. phase-02-concurrency.md — Tuần 4-5
3. phase-03-stdlib.md — Tuần 6-7
4. phase-04-generics.md — Tuần 8
5. phase-05-production.md — Tuần 9-10
6. phase-06-ecosystem.md — Tuần 11
7. phase-07-capstone.md — Tuần 12
