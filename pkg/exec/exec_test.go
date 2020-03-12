package exec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	c := Command("echo", "Hello, world!").Run()
	require.NoError(t, c.Error())
	require.Equal(t, 0, c.ProcessState.ExitCode())
	require.Equal(t, "Hello, world!\n", c.Out())
	require.Equal(t, "", c.Err())

	c = Command("bash", "-c", "echo hello 1>&2").Run()
	require.Equal(t, "", c.Out())
	require.Equal(t, "hello\n", c.Err())
}

// func TestOutput(t *testing.T) {
// 	stdout, stderr, c := Command("echo", "Hello, world!").Output()
// 	require.NoError(t, c.Error())
// 	require.Equal(t, "Hello, world!\n", stdout)
// 	require.Equal(t, "", stderr)
// 	require.Empty(t, stderr)
// }

// func TestOutputBytes(t *testing.T) {
// 	stdout, stderr, c := Command("echo", "Hello, world!").OutputBytes()
// 	require.NoError(t, c.Error())
// 	require.Equal(t, []byte("Hello, world!\n"), stdout)
// 	require.Equal(t, []byte(""), stderr)
// 	require.Empty(t, stderr)
// }

func TestRunFg(t *testing.T) {
	c := Command("echo", "Hello, world!").Run()
	require.NoError(t, c.Error())
}

func TestBake(t *testing.T) {
	echo := Command("echo", "Hello, World").Bake()
	c := echo("again!").Run()
	require.NoError(t, c.Error())
	require.Equal(t, 0, c.ProcessState.ExitCode())

	echo = echo("again").Bake()
	c = echo("and again!!").Run()
	require.NoError(t, c.Error())
	require.Equal(t, "Hello, World again and again!!\n", c.Out())
	require.Equal(t, "", c.Err())
	require.Empty(t, c.Err())
}

func TestIn(t *testing.T) {
	c := Command("pwd").Dir("/").Run()
	require.NoError(t, c.Error())
	require.Equal(t, "/\n", c.Out())

	c = Command("pwd").Dir("/etc").Run()
	require.NoError(t, c.Error())
	require.Equal(t, "/etc\n", c.Out())
}

func TestInBake(t *testing.T) {
	pwd := Command("pwd").Bake()

	c := pwd().Dir("/").Run()
	require.NoError(t, c.Error())
	require.Equal(t, "/\n", c.Out())

	c = pwd().Dir("/etc").Run()
	require.NoError(t, c.Error())
	require.Equal(t, "/etc\n", c.Out())
}

func dirlist(dir string) (dirs []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path == dir || !info.IsDir() {
			return nil
		}

		dirs = append(dirs, path)
		return filepath.SkipDir
	})
	return
}

func TestWithTempdir(t *testing.T) {
	var tdir string
	dir, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, WithTempdir(dir, "testing", func() error {
		tdir, err = os.Getwd()
		require.NoError(t, err)
		require.NotEqual(t, dir, tdir)

		require.True(t, strings.HasPrefix(tdir, dir))

		listdirs, err := dirlist(dir)
		require.NoError(t, err)
		require.Contains(t, listdirs, tdir)

		return nil
	}))

	currdir, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, dir, currdir)

	listdirs, err := dirlist(dir)
	require.NoError(t, err)
	require.NotContains(t, listdirs, tdir)

	require.NoError(t, WithTempdir("", "", func() error {
		tdir, err = os.Getwd()
		require.NoError(t, err)
		require.NotEqual(t, dir, tdir)

		require.True(t, strings.HasPrefix(tdir, os.TempDir()))

		listdirs, err = dirlist(os.TempDir())
		require.NoError(t, err)
		require.Contains(t, listdirs, tdir)

		return nil
	}))

	currdir, err = os.Getwd()
	require.NoError(t, err)
	require.Equal(t, dir, currdir)

	listdirs, err = dirlist(dir)
	require.NoError(t, err)
	require.NotContains(t, listdirs, tdir)

}

// func TestPipe(t *testing.T) {
// 	echo := Command("echo", "aabba\naba\nabbaa")
// 	grep := Command("grep", "abba")
// 	grep = echo.Pipe(grep)

// 	require.NoError(t, grep.Error())
// 	require.Equal(t, "aabba\nabbaa\n", grep.Out())

// 	echo = Command("echo", "aabba\naba\nabbaa")
// 	grep1 := Command("grep", "abba")
// 	grep2 := Command("grep", "baa")
// 	grep2 = grep2.Pipe(echo.Pipe(grep1))

// 	require.NoError(t, grep2.Error())
// 	require.Equal(t, "", grep2.Err())
// 	require.Equal(t, "abbaa", grep2.Out())
// }
