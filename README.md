# flatpdf: A GO cli tool to flatten PDF files

`flatpdf` is a cli tool to flatten PDF files to make them read-only. It is not for flattening out form fields, annotations, or other interactive elements alone but for flattening the entire PDF file. It converts a given PDF file into a new PDF file with all the interactive elements removed. The new PDF file is read-only and cannot be edited.

`flatpdf` is made possible by the following packages:
- [go-pdfium](https://github.com/klippa-app/go-pdfium)
- [pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Installation
1. Install [GO](https://go.dev/doc/install)
2. Compile `flatpdf`
```bash
git clone https://github.com/blucin/flatpdf.git
cd flatpdf
CGO_ENABLED=0 go build
```
3. Add `flatpdf` to your PATH (Read more: [go docs](https://go.dev/doc/tutorial/compile-install))

## Usage
```
Usage:
  flatpdf flat [pdf files] [flags]

Flags:
  -d, --density int     image density (default 600)
  -h, --help            help for flat
  -i, --input strings   Input PDF file(s)
  -o, --output string   output directory relative to the current path. Does not add _flat suffix on save
 so it will override any input files present in the output directory
  -q, --quality int     image quality (0-100) for jpeg encoding (default 99)
```

## Example
1. Flatten a single pdf
```bash
flatpdf flat -i filename.pdf  #output: filename_flat.pdf
```

2. Flatten multiple pdfs with specified output directory.
```bash
flatpdf flat -i filename1.pdf filename2.pdf -o outDir #output: outDir/filename1.pdf outDir/filename2.pdf
```

> [!CAUTION]
> Make sure you don't run example 2 with any input files present in the `outDir`.
> `flatpdf` won't add any suffix and might lead to corruption of the input file.
>
> Example misusage of the command:
> ```bash
> flatpdf flat -i filename1.pdf outDir/filename2.pdf -o outDir
> ```
