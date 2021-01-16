# Go Exec

Conveniently shell out from Golang.  Many libraries for executing commands in
go are about more fine-grained control of subprocesses.  This library is not
that.  This library is for conveniently mimicking the command line from go in
more scripting and command-line tooling applications. 

This introduces a type `Cmd` that wraps the built-in `Command` type.


## Examples

Run a command and capture output:
```golang
c := exec.Command("echo", "Hello, World!").Run()
if c.IsError() {
    fmt.Println(c.Error())
    fmt.Println(c.ProcessState.ExitCode())
}
fmt.Printf("out: '%s'\n", c.Out())
fmt.Printf("err: '%s'\n\n", c.Err())
```

Run a command in the foreground:
```golang
c := exec.Command("echo", "Hello, World!").RunFg()
fmt.Println(c.IsError())
```

Change directories in a scope:
```golang
exec.Command("pwd").RunFg()
exec.WithTempdir("", "", func() error {
    exec.Command("pwd").RunFg()
    return nil
})
exec.Command("pwd").RunFg()
```

Piping commands one into another
```golang
p := exec.Pipe(
    exec.Command("echo", "aabba\naba\nabbaa"),
    exec.Command("cat"),
    exec.Command("grep", "abb"),
)
p = p.Pipe(
    exec.Command("grep", "aa$"),
)

c = p.RunFg()
```

Baking commands:
```golang
echo := exec.Command("echo").Bake()
echo = echo("again").Bake()
echo("and again!!").RunFg()
```
