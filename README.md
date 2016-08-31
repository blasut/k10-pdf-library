# k10 pdf library

This is a small-thrown together program that:
1) Watches a folder for file changes
2) Runs new pdf files through `ghostscript` and outputs the covers as png's
3) Presents the cover images via a web interface to download said pdf's

## Usage
The program uses three flags:
- `-port`: Port used for the web interface [default: 9001]
- `-pdfpath`: Absolute path to where you store your pdf's
- `-thumbpath`: Absolute path to where you want your images

Example: `$ k10-pdf-library -pdfpath /Users/erikvalerius/pdfs/ -port 8888 -thumbpath /Users/erikvalerius/watchtest/thumbs/`

Note that the `pdfpath` and `thumbpath` directories will be exposed via the static file server and should be treated as public directories.

Oh and the paths need trailing slash to work because of lazy string concats. Sry.
