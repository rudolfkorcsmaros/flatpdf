package main

import (
	"fmt"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/karmdip-mi/go-fitz"
	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func open() string {
	openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
		Title: "Open A File",
		Role:  "OpenFileExample",
		FileFilters: []cfd.FileFilter{
			{
				DisplayName: "PDFs (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
		SelectedFileFilterIndex: 2,
		FileName:                "",
		DefaultExtension:        "pdf",
	})
	if err != nil {
		return ""
	}
	if err := openDialog.Show(); err != nil {
		return ""
	}
	result, err := openDialog.GetResult()
	if err == cfd.ErrorCancelled {
		return ""
	} else if err != nil {
		return ""
	}
	return result
}

func main() {
	fmt.Print("flatPDF by Rudolf Korcsmáros v1.0\n\n")

	file := open()

	if file == "" {
		fmt.Println("Error!")
		return
	}

	doc, err := fitz.New(file)
	if err != nil {
		panic(err)
	}
	folder := strings.TrimSuffix(path.Base(file), filepath.Ext(path.Base(file)))

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	fmt.Println("Detected " + strconv.Itoa(doc.NumPage()) + " pages.")

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)

		resizedImg := resize.Resize(720, 0, img, resize.Bilinear)

		if err != nil {
			panic(err)
		}

		pdf.AddPage()
		pdf.ImageFrom(resizedImg, 0, 0, &gopdf.Rect{W: 595.28, H: 841.89})

		fmt.Println(strconv.Itoa(n+1) + " of " + strconv.Itoa(doc.NumPage()) + " pages done.")
	}

	pdf.WritePdf(folder + "_flat.pdf")

	fmt.Println("Done.")
}
