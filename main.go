package main

import (
	"fmt"
	"os"

	"minigit/internal/repo"
)

func run(args []string) int {
	if len(args) == 0 {
		return 1
	}

	switch args[0] {
	case "help":
		// lệnh help thành công
		return 0

	case "init":
		// lấy thư mục hiện tại làm target của init
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: không lấy được current directory: %v\n", err)
			return 1
		}
		// gọi domain function, không viết filesystem logic ở đây
		r, err := repo.Init(cwd, repo.InitOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Printf("Initialized empty repository in %s\n", r.MetaDir)
		return 0

	default:
		// mọi lệnh không xác định đều trả về exit code 2
		return 2
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
