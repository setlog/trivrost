// +build ignore

// Modified version of https://github.com/tv42/becky

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var assetNoDev = "package main\n\ntype asset struct {\n\tName    string\n\tContent string\n\tetag    string\n}\n"

var (
	flagVar           = flag.String("var", "", "variable name to use, \"_\" to ignore (default: file basename without extension)")
	flagWrap          = flag.String("wrap", "", "wrapper function or type (default: filename extension)")
	flagLib           = flag.Bool("lib", true, "generate asset_*.gen.go files defining the asset type")
	flagIgnoreMissing = flag.Bool("ignore-missing", false, "write empty data if file is missing")
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [OPTS] FILE..\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Creates files FILE.gen.go and asset_*.gen.go\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if flag.NArg() > 1 && *flagVar != "" && *flagVar != "_" {
		log.Fatal("cannot combine -var with multiple files")
	}

	packages := map[string]*build.Package{}

	for _, filename := range flag.Args() {
		dir, base := filepath.Split(filename)
		if dir == "" {
			dir = "."
		}

		pkg, err := getPkg(packages, dir)
		if err != nil {
			log.Fatal(err)
		}

		variable := *flagVar
		if variable == "" {
			variable = strings.SplitN(base, ".", 2)[0]
		}

		wrap := *flagWrap
		if wrap == "" {
			wrap = filepath.Ext(base)
			if wrap == "" {
				log.Fatalf("files without extension need -wrap: %s", filename)
			}

			wrap = wrap[1:]
		}

		if err := process(filename, pkg.Name, variable, wrap); err != nil {
			log.Fatal(err)
		}

	}
}

// autogen writes a warning that the file has been generated automatically.
func autogen(w io.Writer) error {
	// broken into parts here so grep won't find it
	const warning = "// AUTOMATICALLY " + "GENERATED FILE. DO NOT EDIT.\n\n"
	_, err := io.WriteString(w, warning)
	return err
}

func process(filename, pkg, variable, wrap string) error {
	var src io.Reader
	srcFile, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) && *flagIgnoreMissing {
			log.Println("treating missing file", filename, "as empty.")
			src = bytes.NewReader([]byte{})
		} else {
			return err
		}
	} else {
		defer srcFile.Close()
		src = srcFile
	}

	tmp, err := ioutil.TempFile(filepath.Dir(filename), ".tmp.asset-")
	if err != nil {
		return err
	}
	defer func() {
		if tmp != nil {
			_ = os.Remove(tmp.Name())
		}
	}()
	defer tmp.Close()

	in := bufio.NewReader(src)
	out := bufio.NewWriter(tmp)

	if err := autogen(out); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "package %s\n\nfunc init() {\n", pkg); err != nil {
		return err
	}
	if err := embed(variable, wrap, filepath.Base(filename), in, out); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "}\n"); err != nil {
		return err
	}
	if err := out.Flush(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	gen := filename + ".gen.go"
	if err := os.Rename(tmp.Name(), gen); err != nil {
		return err
	}
	tmp = nil
	return nil
}

func embed(variable, wrap, filename string, in io.Reader, out io.Writer) error {
	_, err := fmt.Fprintf(out, "\t%s = %s(\"\" +\n",
		variable, wrap)
	if err != nil {
		return err
	}
	buf := make([]byte, 1*1024*1024)
	eof := false
	for !eof {
		n, err := in.Read(buf)
		switch err {
		case io.EOF:
			eof = true
		case nil:

		default:
			return err
		}
		if n == 0 {
			continue
		}
		s := string(buf[:n])
		s = strconv.QuoteToASCII(s)
		s = "\t\t" + s + " +\n"
		if _, err := io.WriteString(out, s); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, "\t\t\"\")\n"); err != nil {
		return err
	}
	return nil
}

func getPkg(packages map[string]*build.Package, dir string) (*build.Package, error) {
	if pkg, found := packages[dir]; found {
		return pkg, nil
	}

	pkg, err := loadPkg(dir)
	if err != nil {
		return nil, err
	}

	packages[dir] = pkg
	return pkg, nil
}

func loadPkg(dir string) (*build.Package, error) {

	if !filepath.IsAbs(dir) {
		if abs, err := filepath.Abs(dir); err == nil {
			dir = abs
		}
	}

	pkg, err := build.ImportDir(dir, 0)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
