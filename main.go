package main

import "os"

func run(args []string) int {
	if len(args) == 0 {
		return 1
	}

	switch args[0] {
	case "help":
		return 0
	case "unknown":
		return 2
	default:
		return 1
	}
}

// TODO-01-CLI: Command dispatcher chỉ parse command và gọi package domain.
// SENIOR ASKS: Tại sao tách `run(args []string) int` khỏi `main()` giúp test CLI dễ hơn?
// HINT: Trong test, bạn truyền args trực tiếp thay vì spawn process.

// lý do dùng os.Args thay cho flag vì các yếu tố sau
// 1. args đơn giản
// 2. script rất nhanh
// nhưng nên dùng flag nếu
// cần các flag và type conversion như -v -h,... nhiều option và dễ kiểm soát hơn
func main() {
	os.Exit(run(os.Args[1:]))
}
