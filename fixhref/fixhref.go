// Package fixref ...
package fixhref

import (
	"fmt"
	"regexp"
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
---
title: 'Design of fixhref program'
author:
- Igor Lesik
---

Program _fixhref_ fixes relative href links in a group of HTML files
located in single root directory.

Package `fixhref` export only one function, it is called `FixHtmlHref`.
In this function we walk through all sub-directories of the root directory
to visit every HTML file.

```{#FixHtmlHref .go .numberLines}
PROSE_BEGIN*/

// Fixes relative hyperlinks in HTML files.
func FixHtmlHref(rootDirPath string) error {
	err := filepath.Walk(rootDirPath,
			func(path string, f os.FileInfo, err error) error {
		return visitFile(path, f, err)
	})
	return err
}
/*PROSE_BEGIN
```

Function `FixHtmlHref` calls function `visitFile` while walking
inside root directory and visiting each file.
Inside function `visitFile`, we first check that there was no
error while walking, next we ignore directory files, at last we
call function `visitHtmlFile` if currently visited file
has .html extention.

```{#visitFile .go .numberLines}
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
/*PROSE_BEGIN
```

Function `visitHtmlFile` opens the file for reading and writing
and then calls `fixHtmlFile` providing the handle for open file
and also `Reader` and `Writer` interfaces.
We open the file for writing because we are going to modify the
file in place, if it has broken links to be fixed.

```{#visitHtmlFile .go .numberLines}
PROSE_BEGIN*/

func visitHtmlFile(path string, finfo os.FileInfo) error {
	fmt.Printf("File %s...\n", path)
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error: can't open %s, %s\n", path, err.Error())
		return err
	}
	defer f.Close()
	return fixHtmlFile(path, finfo, f, bufio.NewReader(f), bufio.NewWriter(f))
}
/*PROSE_BEGIN
```
PROSE_END*/


func fixHtmlFile(
	path string,
	finfo os.FileInfo,
	f *os.File,
	r io.Reader,
	w *bufio.Writer) error {

	htmlDoc, err := html.Parse(r)
	if err != nil {
		fmt.Printf("Error: can't parse %s, %s\n", path, err.Error())
		return err
	}

	numFixes, err := fixHtmlDoc(path, htmlDoc)
	if err != nil {
		return err
	}

	if numFixes > 0 {
		err = writeHtmlFile(htmlDoc, path, f, w)
		if err != nil {
			return err
		}
	}

        fmt.Printf("File %s: %d fixes\n", path, numFixes)

	return nil
}

func writeHtmlFile(htmlDoc *html.Node, path string, f *os.File, w *bufio.Writer) error {
	defer w.Flush()
	f.Seek(0,0)
	_, err := io.WriteString(w, "<!-- Relative links were modified by fixhref -->\n")
	if err != nil {
		fmt.Printf("Error: can't write to %s, %s\n", path, err.Error())
		return err
	}

	return html.Render(w, htmlDoc)
}

func fixHtmlDoc(path string, htmlDoc *html.Node) (uint, error) {

	var numFixesTotal uint = 0

	var findHref func(*html.Node)

	findHref = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for idx, a := range n.Attr {
				if a.Key == "href" {
					numFixes := fixHrefAttr(path, a, &n.Attr[idx])
					numFixesTotal += numFixes
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findHref(c)
		}
	}

	findHref(htmlDoc)

	return numFixesTotal, nil
}

func fixHrefAttr(path string, rdAttr html.Attribute, wrAttr *html.Attribute) uint {
	switch isHrefFixRequired(path, rdAttr.Val) {
	case true:
		return modifyHref(path, rdAttr.Val, &wrAttr.Val)
	case false:
		return 0
	}
	return 0
}

func isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func isRefFileExist(path string, href string) bool {
	linkPath := filepath.Join(filepath.Dir(path), href)
	exist, _ := isFileExist(linkPath)
	return exist
}

func isHrefFixRequired(path string, href string) bool {
	isAbsHref := func(href string) bool {
		match, _ := regexp.MatchString("^[a-zA-Z]+://", href)
		return match
	}

	return !isAbsHref(href) && !isRefFileExist(path, href)
}

func modifyHref(path string, href string, wrHref *string) uint {
	writeHrefFix := func(newHref string, wrHref *string) {
		*wrHref = newHref
	}

	fixFound, newHref := findHrefFix(path, href)
	switch fixFound {
	case true:
		writeHrefFix(newHref, wrHref)
		fmt.Printf("File %s: change '%s' to '%s'\n", path, href, *wrHref)
		return 1
	case false:
		fmt.Printf("File %s: can't fix '%s'\n", path, href)
		return 0
	}
	return 0
}

func findHrefFix(path string, href string) (bool, string) {
	type transform struct {
		pattern string
		replace string
	}

	transforms := [1](transform){transform{"(.*)(group1)(.*)","${1}g1${3}"}}

	for _, tran := range transforms {
		re := regexp.MustCompile(tran.pattern)
		fix := re.ReplaceAllString(href, tran.replace)
		fmt.Printf("try fix: '%s', path: '%s', new: %s\n", fix, path,
			filepath.Join(filepath.Dir(path), fix))
		if isRefFileExist(path, fix) == true {
			return true, fix
		}
	}

	return false, href
}

