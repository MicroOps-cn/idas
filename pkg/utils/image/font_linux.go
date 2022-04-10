package image

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"idas/pkg/utils/sets"
)

func loadSystemFonts(ctx context.Context, fontNames sets.Set[string]) (font *truetype.Font, err error) {
	cmd := exec.CommandContext(ctx, "fc-list")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	fonts := map[string]string{}
	scanner := bufio.NewScanner(stdout)
	for {
		if scanner.Scan() {
			fontInfo := bytes.Split(scanner.Bytes(), []byte(":"))
			if len(fontInfo) <= 2 {
				continue
			}
			for _, fontName := range strings.Split(string(fontInfo[1]), ",") {
				if fontNames.Has(strings.TrimSpace(fontName)) {
					fonts[strings.TrimSpace(fontName)] = string(fontInfo[0])
				}
			}
		} else {
			break
		}
	}
	fmt.Println(fonts)
	fmt.Println(fontNames.List())
	for _, fontName := range fontNames.List() {
		fmt.Println(fontName)
		if fontPath, ok := fonts[fontName]; ok {
			fmt.Println(fontName)
			var fontBytes []byte
			fontBytes, err = ioutil.ReadFile(fontPath)
			if err != nil {
				continue
			}
			fmt.Println(fontName)
			font, err = freetype.ParseFont(fontBytes)
			if err != nil {
				continue
			} else {
				return font, nil
			}
		}
	}
	//_ = filepath.Walk("/usr/share/fonts/", func(filePath string, info fs.FileInfo, err error) error {
	//	_, filename := path.Split(filePath)
	//	filename = strings.ToLower(filename)
	//	ext := filepath.Ext(filename)
	//	if fontSuffixSet.Has(ext) {
	//		fmt.Println(filename)
	//	}
	//	if fontSuffixSet.Has(ext) && fontNames.Has(filename) {
	//		fmt.Println(filename)
	//		var fontBytes []byte = nil
	//		fontBytes, err = ioutil.ReadFile(filePath)
	//		if err != nil {
	//			return err
	//		}
	//		font, err = freetype.ParseFont(fontBytes)
	//		if err == nil {
	//			return filepath.SkipDir
	//		}
	//	}
	//	return nil
	//})
	//return nil
	return
}
