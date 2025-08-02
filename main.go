package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	var jackSrcFiles, cmpFile string
	flag.StringVar(&jackSrcFiles, "s", "", "source file in jack extension (e.g. Add.jack or a Directory with multiple jack files)")
	flag.StringVar(&cmpFile, "c", "", "compare file in xml extension (e.g. Add.xml)")
	flag.Parse()
	if jackSrcFiles == "" {
		fmt.Println("No source file provided")
		flag.Usage()
		os.Exit(1)
	}
	// check if the source is a directory
	srcStat, err := os.Stat(jackSrcFiles)
	if err != nil {
		fmt.Println("Error getting source file status", err)
		os.Exit(1)
	}

	jackFiles := []*os.File{}
	defer func() {
		for _, f := range jackFiles {
			f.Close()
		}
	}()

	if srcStat.IsDir() {
		// get all jack files in the directory
		rawJackFiles, err := filepath.Glob(filepath.Join(jackSrcFiles, "*.jack"))
		if err != nil {
			fmt.Printf("Error listing jack files %s: %s\n", jackSrcFiles, err)
			os.Exit(1)
		}
		for _, jackFile := range rawJackFiles {
			srcF, err := os.Open(jackFile)
			if err != nil {
				fmt.Printf("Error opening jack file %s: %s\n", jackFile, err)
				os.Exit(1)
			}
			jackFiles = append(jackFiles, srcF)
		}
	} else {
		srcF, err := os.Open(jackSrcFiles)
		if err != nil {
			fmt.Printf("Error opening jack file %s: %s\n", jackSrcFiles, err)
			os.Exit(1)
		}
		jackFiles = append(jackFiles, srcF)
	}

	if len(jackFiles) == 0 {
		fmt.Printf("No jack files found in %s\n", jackSrcFiles)
		os.Exit(1)
	}

	// run the process in parallel
	wg := sync.WaitGroup{}
	for _, jackFile := range jackFiles {
		wg.Add(1)
		go func(jackFile *os.File) {
			defer wg.Done()
			processJackFile(jackFile)
		}(jackFile)
	}
	wg.Wait()

	fmt.Printf("Analysis complete for %d files ✅\n", len(jackFiles))
}

func processJackFile(jackFile *os.File) {
	jackFileContent, err := io.ReadAll(jackFile)
	if err != nil {
		fmt.Printf("Error reading jack file %s: %s\n", jackFile.Name(), err)
		os.Exit(1)
	}
	tokenizer, err := NewTokenizer(string(jackFileContent))
	if err != nil {
		printError(jackFile.Name(), err)
		os.Exit(1)
	}

	tokensFile := bytes.Buffer{}
	tokens := []Token{}

	for {
		token, err := tokenizer.Advance()
		if err != nil {
			if err == errNoMoreTokens {
				break
			}
			fmt.Printf("Error advancing tokenizer %s: %s\n", jackFile.Name(), err)
			os.Exit(1)
		}
		tokens = append(tokens, token)
	}

	fmt.Fprintf(&tokensFile, "<tokens>\n")
	for _, token := range tokens {
		fmt.Fprintf(&tokensFile, "%s\n", token.Tag())
	}
	fmt.Fprintf(&tokensFile, "</tokens>\n")

	// create a tokens file with *T.xml
	if err := os.WriteFile(strings.Replace(jackFile.Name(), ".jack", "T.xml", 1), tokensFile.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing tokens file %s: %s\n", jackFile.Name(), err)
		os.Exit(1)
	}

	// create a string buffer instead of a file
	xmlBuffer := bytes.Buffer{}
	ce := NewCompilationEngine(tokenizer, &xmlBuffer)
	if err := ce.ProcessClass(); err != nil {
		printError(jackFile.Name(), err)
		os.Exit(1)
	}

	xmlFile := strings.Replace(jackFile.Name(), ".jack", ".xml", 1)
	xmlFileContent := xmlBuffer.String()
	xmlFileContent = FormatXML(xmlFileContent, "", "  ")
	if err := os.WriteFile(xmlFile, []byte(xmlFileContent), 0644); err != nil {
		fmt.Printf("Error writing xml file %s: %s\n", xmlFile, err)
		os.Exit(1)
	}
	println("Compilation engine complete for", jackFile.Name(), " ✅")
}

func printError(fileName string, err error) {
	if e, ok := err.(*AnalyzerError); ok {
		fmt.Printf("Error in file %s:%d -> %s\n\t%s\n", fileName, e.LineNum, e.Err, e.Line)
		if s := e.Stack; s != "" {
			println("--------------------------------")
			fmt.Println(s)
		}
	} else {
		fmt.Printf("Error in file %s: %s\n", fileName, err)
	}
}
