// Package fixref ...
package fixhref

import (
	"fmt"
	"path/filepath"
	"os"
	"bufio"
        "io"
	"golang.org/x/net/html"
)


// Fixes relative hyperlinks in HTML files.
func FixHtmlHref(rootDirPath string) error {
	/*cache*/
	err := filepath.Walk(rootDirPath,
			func(path string, f os.FileInfo, err error) error {
		return visitFile(path, f, err /*,cache*/)
	})
	return err
}


func visitFile(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visit: %s\n", path)
	if f.IsDir() == false && filepath.Ext(path) == ".html" {
		return visitHtmlFile(path, f)
	}
	return nil
}

func visitHtmlFile(path string, finfo os.FileInfo) error {
	fmt.Printf("Visit: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error: can't open %s\n", path)
		return err
	}
	defer f.Close()
	return fixHtmlFile(path, finfo, f, bufio.NewReader(f), bufio.NewWriter(f))
}


func fixHtmlFile(path string, finfo os.FileInfo, f *os.File, r io.Reader, w *bufio.Writer) error {
	htmlDoc, err := html.Parse(r)
	if err != nil {
		fmt.Printf("Error: can't parse %s\n", path)
		return err
	}
	defer w.Flush()
	f.Seek(0,0)
	_, err = io.WriteString(w, "<-- Modified by fixhref -->\n")
	if err != nil {
		fmt.Printf("Error: can't parse %s\n", path)
		return err
	}

	return html.Render(w, htmlDoc)
}
