package main

import (
	"fmt"

	"github.com/wreed4/go-exec/pkg/exec"
)

func main() {
	c := exec.Command("echo", "Hello, World!").Run()
	fmt.Println(c.IsError())
	if c.IsError() {
		fmt.Println(c.Error())
		fmt.Println(c.ProcessState.ExitCode())
	}
	fmt.Printf("out: '%s'\n", c.Out())
	fmt.Printf("err: '%s'\n\n", c.Err())

	c = exec.Command("echo", "Hello, World!").RunFg()
	fmt.Println(c.IsError())

	exec.Command("pwd").RunFg()
	exec.WithTempdir("", "", func() error {
		exec.Command("pwd").RunFg()
		return nil
	})
	exec.Command("pwd").RunFg()

	p := exec.Pipe(
		exec.Command("echo", "aabba\naba\nabbaa"),
		exec.Command("cat"),
		exec.Command("grep", "abb"),
	)
	p = p.Pipe(
		exec.Command("grep", "aa$"),
	)
	fmt.Println(p, "\n\n-------------")

	c = p.RunFg()
	if c.IsError() {
		fmt.Println("result:", c)
		fmt.Println(c.Error())
		fmt.Println(c.ProcessState.ExitCode())
	}
	fmt.Printf("out: '%s'\n", c.Out())
	fmt.Printf("err: '%s'\n\n", c.Err())
}
