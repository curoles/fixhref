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

// To extract the prose:
// awk '/prose_begin/{flag=1;next}/prose_end/{flag=0}flag' file
// https://stackoverflow.com/questions/38972736/how-to-select-lines-between-two-patterns/38972737#38972737

/*PROSE_BEGIN
# Design of fixhref program

here
is my
prose

link [visitFile](#visitFile)

### Function FixHtmlHref {#FixHtmlHref}

```{#code1 .go .numberLines}
PROSE_END
PROSE_BEGIN*/

// Fixes relative hyperlinks in HTML files.
func FixHtmlHref(rootDirPath string) error {
	/*cache*/
	err := filepath.Walk(rootDirPath,
			func(path string, f os.FileInfo, err error) error {
		return visitFile(path, f, err /*,cache*/)
	})
	return err
}
/*
```

bla bla 

```{#visitFile .go .numberLines}
PROSE_END
PROSE_BEGIN*/

func visitFile(path string, f os.FileInfo, err error) error {
	switch {
	case err != nil:
		return err
	case f.IsDir() == true:
		return nil
	case filepath.Ext(path) == ".html":
		return visitHtmlFile(path, f)
	default:
		return nil
        }
	return nil
}
/*
```
PROSE_END*/

func visitHtmlFile(path string, finfo os.FileInfo) error {
	fmt.Printf("File %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		//fmt.Printf("Error: can't open %s\n", path)
		return err
	}
	defer f.Close()
	return fixHtmlFile(path, finfo, f, bufio.NewReader(f), bufio.NewWriter(f))
}


func fixHtmlFile(
	path string,
	finfo os.FileInfo,
	f *os.File,
	r io.Reader,
	w *bufio.Writer) error {

	htmlDoc, err := html.Parse(r)
	if err != nil {
		//fmt.Printf("Error: can't parse %s\n", path)
		return err
	}

	numFixes, err := fixHtmlDoc(htmlDoc)
	if err != nil {
		return err
	}

	if numFixes > 0 {
		err = writeHtmlFile(htmlDoc, path, f, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeHtmlFile(htmlDoc *html.Node, path string, f *os.File, w *bufio.Writer) error {
	defer w.Flush()
	f.Seek(0,0)
	_, err := io.WriteString(w, "<!-- Modified by fixhref -->\n")
	if err != nil {
		//fmt.Printf("Error: can't write to %s\n", path)
		return err
	}

	return html.Render(w, htmlDoc)
}

func fixHtmlDoc(htmlDoc *html.Node) (uint, error) {

	var numFixes uint = 0

	var findHref func(*html.Node)

	findHref = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for idx, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					n.Attr[idx].Val = a.Val + "fix"
					numFixes += 1
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findHref(c)
		}
	}

	findHref(htmlDoc)

	return numFixes, nil
}
