package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	archFlag := flag.String("arch", "", "override detected architecture (x86_64, i686, aarch64)")
	threadFlag := flag.String("thread", "", "override thread model (posix, win32)")
	excFlag := flag.String("exception", "", "override exception handling (seh, dwarf, sjlj)")
	crtFlag := flag.String("crt", "", "override CRT (ucrt, msvcrt)")
	jsonFlag := flag.Bool("json", false, "output as JSON")
	offlineFlag := flag.Bool("offline", false, "skip network fetch, use embedded snapshot only")
	listFlag := flag.Bool("list", false, "list all matching builds")
	verFlag := flag.Bool("version", false, "show version")
	flag.Parse()

	if *verFlag {
		fmt.Println("mingw-chooser", version)
		os.Exit(0)
	}

	_ = archFlag
	_ = threadFlag
	_ = excFlag
	_ = crtFlag
	_ = jsonFlag
	_ = offlineFlag
	_ = listFlag
	fmt.Println("TODO: implement")
}
