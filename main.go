package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/schollz/progressbar/v3"
)

const (
	alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	nums  = "1234567890"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func alphaOnly(s string) bool {
	for _, char := range s {
		if !strings.Contains(alpha, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func numOnly(s string) bool {
	for _, char := range s {
		if !strings.Contains(nums, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func parse(wordlist string, output string, minLen string, maxLen string, phrase string, specialChar bool, replace bool, number bool) {
	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Error opening wordlist")
		os.Exit(1)
	}
	defer file.Close()

	count, err := lineCounter(file)
	if err != nil {
		fmt.Println("Error Reading wordlist")
		os.Exit(1)
	}

	file, err = os.Open(wordlist)
	if err != nil {
		fmt.Println("Error opening wordlist")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if maxLen == "0" {
		maxLen = "65536"
	}

	minLenint, err := strconv.Atoi(minLen)
	if err != nil {
		fmt.Println("Please put an integer for minlen")
		os.Exit(1)
	}

	maxLenint, err := strconv.Atoi(maxLen)
	if err != nil {
		fmt.Println("Please put an interger for maxlen")
		os.Exit(1)
	}

	outfile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Unable to create out file")
		os.Exit(1)
	}

	bar := progressbar.NewOptions(count,
		progressbar.OptionSetDescription("Parsing "+wordlist),
		progressbar.OptionSetWidth(20))

	found := 0

	if specialChar && number {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) <= maxLenint && len(line) >= minLenint && strings.Contains(strings.ToLower(line), strings.ToLower(phrase)) && alphaOnly(line) == false && strings.ContainsAny(line, "0123456789") {
				if _, err := outfile.Write([]byte(line + "\n")); err != nil {
					fmt.Println("Error Writing to out file")
					os.Exit(1)
				}
				found += 1
			}
			bar.Add(1)
		}
	} else if specialChar {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) <= maxLenint && len(line) >= minLenint && strings.Contains(strings.ToLower(line), strings.ToLower(phrase)) && alphaOnly(line) == false {
				if _, err := outfile.Write([]byte(line + "\n")); err != nil {
					fmt.Println("Error Writing to out file")
					os.Exit(1)
				}
				found += 1
			}
			bar.Add(1)
		}
	} else if number {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) <= maxLenint && len(line) >= minLenint && strings.Contains(strings.ToLower(line), strings.ToLower(phrase)) && strings.ContainsAny(line, "0123456789") {
				if _, err := outfile.Write([]byte(line + "\n")); err != nil {
					fmt.Println("Error Writing to out file")
					os.Exit(1)
				}
				found += 1
			}
			bar.Add(1)
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) <= maxLenint && len(line) >= minLenint && strings.Contains(strings.ToLower(line), strings.ToLower(phrase)) {
				if _, err := outfile.Write([]byte(line + "\n")); err != nil {
					fmt.Println("Error Writing to out file")
					os.Exit(1)
				}
				found += 1
			}
			bar.Add(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading lines")
		os.Exit(1)
	}

	if err := outfile.Close(); err != nil {
		fmt.Println("Error Closing outfile")
		os.Exit(1)
	}

	fmt.Println("Found " + strconv.Itoa(found) + " Passwords")
}

func main() {
	start := time.Now()
	// init parser
	parser := argparse.NewParser("list-parse", "Creates more customized wordlist")

	// flags
	var wordlist *string = parser.String("w", "wordlist", &argparse.Options{Required: true, Help: "Wordlist to parse through; Full path to wordlist"})
	var output *string = parser.String("o", "output", &argparse.Options{Required: true, Help: "File to output to"})
	var minLen *string = parser.String("m", "min-length", &argparse.Options{Required: false, Help: "Minimum length of line", Default: "0"})
	var maxLen *string = parser.String("x", "max-length", &argparse.Options{Required: false, Help: "Maximum length of line", Default: "0"})
	var phrase *string = parser.String("p", "phrase", &argparse.Options{Required: false, Help: "Phrase/word that is required to be in the line", Default: ""})
	var specialChar *bool = parser.Flag("s", "require-special-characters", &argparse.Options{Required: false, Help: "Require special characters"})
	var number *bool = parser.Flag("n", "require-number", &argparse.Options{Required: false, Help: "Require Number"})
	var replcLttrWSpecChar *bool = parser.Flag("r", "replace-letters", &argparse.Options{Required: false, Help: "Replace Letters with special characters in phrase search"})

	// parse through arguments given
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	parse(*wordlist, *output, *minLen, *maxLen, *phrase, *specialChar, *replcLttrWSpecChar, *number)
	elapsed := time.Since(start)
	fmt.Println("\nParsed in", elapsed)
}
