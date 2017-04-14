package main

import (
	"bufio"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please provide the file path")
		return
	}

	file := os.Args[1]
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("file not exist")
		return
	}

	goArguments := []string{"test", "-c", "-o", "/dev/null"}
	base := path.Base(file)
	tempDir, err := ioutil.TempDir("", "goflycheck_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	pkg, err := build.ImportDir(path.Dir(file), build.AllowBinary)
	tempFileName := path.Join(tempDir, base)
	goArguments = append(goArguments, tempFileName)
	if err == nil {
		var files []string
		files = append(files, pkg.GoFiles...)
		files = append(files, pkg.CgoFiles...)
		for _, f := range files {
			if f == base {
				continue
			}
			goArguments = append(goArguments, path.Join(tempDir, f))
			copyFileContents(path.Join(path.Dir(file), f), path.Join(tempDir, f))
		}
	}

	tempFile, err := os.Create(tempFileName)
	if err != nil {
		fmt.Println("os create error", err)
	}
	defer func() {
		tempFile.Close()
	}()

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 1000)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		tempFile.Write(buf[:n])
	}

	cmd := exec.Command("go", goArguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	fmt.Println(string(out))
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
