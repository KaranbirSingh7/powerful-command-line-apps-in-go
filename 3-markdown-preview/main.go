package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>{{ .Title }}</title> 
	</head>
	<p>file: {{ .Filename }}</p>
	<body>{{ .Body }}</body>
</html>`
)

// content type represents the HTML content to add into the template
type content struct {
	Title    string
	Body     template.HTML
	Filename string
}

func main() {

	filename := flag.String("file", "", "Markdown file to preview/process")
	skipPreview := flag.Bool("skip-preview", false, "skip previewing html file in browser")
	templateFilename := flag.String("template", "", "Template file to use for header and footer")
	flag.Parse()

	// if flag is not provided, exit
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview, *templateFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// define executable based on Operating System/OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "darwin":
		cName = "open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{
			"/C", "start",
		}
	default:
		return fmt.Errorf("OS not supported")
	}

	// append file to param slaice
	cParams = append(cParams, fname)

	// locate exec in path
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// run open command to open file in browser
	err = exec.Command(cPath, cParams...).Run()

	// wait since we have another function that will delete this temporary file as soon as it in opened
	time.Sleep(2 * time.Second)
	return err
}

func run(filename io, out io.Writer, skipPreview bool, tFname string) error {
	// read file
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// convert markdown -> html header + body + footer
	htmlData, err := parseContent(input, tFname, filename)
	if err != nil {
		return err
	}

	// create a temp file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	// remove temporary file
	defer os.Remove(outName)

	return preview(outName)
}

// saveHTML writes html content to a file
func saveHTML(filename string, htmlData []byte) error {
	return os.WriteFile(filename, htmlData, 0644)
}

// parseContent parse markdown contents and convert it into html
func parseContent(input []byte, tFname string, filename string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// parse the contents of the default template const into new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// if provided by user, then use that template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// set new field
	c := content{
		Title:    "Markdown Preview Tool",
		Body:     template.HTML(body),
		Filename: filepath.Base(filename),
	}

	// render template with our variables
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
