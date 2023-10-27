# flatpdf: A GO cli tool to flatten PDF files
Note: Requires [imagemagick](https://imagemagick.org/) as a dependency, see [Installation](#installation)

`flatpdf` is a cli tool to flatten PDF files to make them read-only. It is not for flattening out form fields, annotations, or other interactive elements alone but for flattening the entire PDF file. It converts a given PDF file into a new PDF file with all the interactive elements removed. The new PDF file is read-only and cannot be edited. `flatpdf` supports go threads to flatten multiple pdfs more efficiently.

`flatpdf` is made possible by the following packages:
- [imagick](https://github.com/gographics/imagick)
- [pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Installation
1. Install [GO](https://go.dev/doc/install)
2. Install [imagemagick](https://imagemagick.org/) for your OS (Make sure you get `libmagick-dev` with it if on *nix)
e.g.
```bash
yay imagemagick #AUR
```
3. Confirm [Installation](https://stackoverflow.com/questions/40063438/wheres-libmagickwand-dev-i-installed-them-all-and-how-can-i-include-it-to-a-c) on *nix
4. Compile `flatpdf`
```bash
go build
```
5. Add `flatpdf` to your PATH (Read more: [go docs](https://go.dev/doc/tutorial/compile-install))

## Usage
```
flatpdf is a multithreaded pdf flattener to make them read-only.
Pass pdf files as arguments to flat, new pdf files will be generated
with the same name but with the suffix '_flat'. Use the -h flag to 
see all available options.

Usage:
  flatpdf [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  flat        Flats the pdf files passed as arguments
  help        Help about any command
  version     Print the version number of flatpdf
```

- flat command flags:

```
Usage:
  flatpdf flat [pdf files] [flags]

Flags:
  -d, --density int     image density (default 600)
  -h, --help            help for flat
  -o, --output string   output directory
  -q, --quality int     image quality (0-100) (default 99)
  -t, --threads int     number of threads (default 0)
```

> If you try to overwrite an already flatten pdf e.g. **`flatpdf filename.pdf filename_flat.pdf`** then instead of overwriting, flatpdf will append pages to it from itself. Expected behaviour will be added later.

## Example
1. Flatten a single pdf
```bash
flatpdf filename.pdf  #output: filename_flat.pdf
```
2. Flatten multiple pdfs with go threads
```bash
flatpdf filename1.pdf filename2.pdf ... filenameN.pdf -t 8
```
3. Flatten multiple pdfs without threads
```bash
flatpdf filename1.pdf filename2.pdf #output: filename1_flat.pdf filename2_flat.pdf
```
