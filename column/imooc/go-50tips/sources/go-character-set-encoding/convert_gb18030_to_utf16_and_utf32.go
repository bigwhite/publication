package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

func catFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func gb18030ToUtf16BE(in []byte) ([]byte, error) {
	r := bytes.NewReader(in) //gb18030

	// to rune(the utf8 representation of code point)
	s := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	d := transform.NewReader(s,
		unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()) // to utf16BE, no bom

	out, err := ioutil.ReadAll(d)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func gb18030ToUtf32BE(in []byte) ([]byte, error) {
	r := bytes.NewReader(in) //gb18030

	// to rune(the utf8 representation of code point)
	s := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	d := transform.NewReader(s,
		utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewEncoder()) // to utf32BE, no bom

	out, err := ioutil.ReadAll(d)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func main() {
	src, err := catFile("gb18030.txt")
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}

	// gb18030 to utf-16be
	dst, err := gb18030ToUtf16BE(src)
	if err != nil {
		fmt.Println("convert error:", err)
		return
	}

	fmt.Printf("UTF-16BE(no BOM)编码: ")
	for _, v := range dst {
		fmt.Printf("0x%X ", v)
	}
	fmt.Printf("\n")

	// gb18030 to utf-32be
	dst1, err := gb18030ToUtf32BE(src)
	if err != nil {
		fmt.Println("convert error:", err)
		return
	}

	fmt.Printf("UTF-32BE(no BOM)编码: ")
	for _, v := range dst1 {
		fmt.Printf("0x%X ", v)
	}
	fmt.Printf("\n")
}
