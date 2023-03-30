package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/schollz/progressbar/v3"
)

func main() {
	start := time.Now()

	parser := argparse.NewParser("listparse", "Creates more customized wordlist")

	var wordlist *string = parser.String("w", "wordlist", &argparse.Options{Required: true, Help: "Wordlist to parse through; Full path to wordlist"})
	var output *string = parser.String("o", "output", &argparse.Options{Required: true, Help: "File to output to"})
	var minLen *string = parser.String("m", "min-length", &argparse.Options{Required: false, Help: "Minimum length of line", Default: "0"})
	var maxLen *string = parser.String("x", "max-length", &argparse.Options{Required: false, Help: "Maximum length of line", Default: "0"})
	var phrase *string = parser.String("p", "phrase", &argparse.Options{Required: false, Help: "Phrase/word that is required to be in the line", Default: ""})
	var specialChar *bool = parser.Flag("s", "require-special-characters", &argparse.Options{Required: false, Help: "Require special characters"})
	var number *bool = parser.Flag("n", "require-number", &argparse.Options{Required: false, Help: "Require Number"})
	var verbose *bool = parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "Verbose mode"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	count, err := parse(*wordlist, *output, *minLen, *maxLen, *phrase, *specialChar, *number, *verbose)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	fmt.Println("\nParsed in", elapsed)
	fmt.Println("Found: " + count + " passwords")

}

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

func genericCheck(line string, maxLen int, minLen int, phrase string) bool {
	if len(line) <= maxLen && len(line) >= minLen && strings.Contains(strings.ToLower(line), strings.ToLower(phrase)) {
		return true
	} else {
		return false
	}
}

func parse(wordlist string, output string, minLen string, maxLen string, phrase string, specialChar bool, number bool, verbose bool) (string, error) {
	file, err := os.Open(wordlist)
	if err != nil {
		return "0", errors.New("error Opening Wordlist")
	}
	defer file.Close()

	if verbose {
		fmt.Println("Opened: ", wordlist)
	}

	count, err := lineCounter(file)
	if err != nil {
		return "0", errors.New("error Reading Wordlist")
	}

	if verbose {
		fmt.Println("Found ", count, " lines")
	}

	file, err = os.Open(wordlist)
	if err != nil {
		return "0", errors.New("error Opening Wordlist")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if maxLen == "0" {
		maxLen = "65536"
	}

	minLenint, err := strconv.Atoi(minLen)
	if err != nil {
		return "0", errors.New("minlen must be an integer")
	}

	maxLenint, err := strconv.Atoi(maxLen)
	if err != nil {
		return "0", errors.New("maxlen must be an integer")
	}

	outfile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "0", errors.New("unable to Create outfile")
	}

	bar := progressbar.NewOptions(count,
		progressbar.OptionSetDescription("Parsing "+wordlist),
		progressbar.OptionSetWidth(20))

	found := 0

	if specialChar && number {
		for scanner.Scan() {
			line := scanner.Text()
			if genericCheck(line, maxLenint, minLenint, phrase) && !alphaOnly(line) && strings.ContainsAny(line, "0123456789") {
				if verbose {
					fmt.Println("Found: ", line)
				}
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
			if genericCheck(line, maxLenint, minLenint, phrase) && !alphaOnly(line) {
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
			if genericCheck(line, maxLenint, minLenint, phrase) && strings.ContainsAny(line, "0123456789") {
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
			if genericCheck(line, maxLenint, minLenint, phrase) {
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
		return "0", errors.New("error reading lines")
	}

	if err := outfile.Close(); err != nil {
		return "0", errors.New("error closing outfile")
	}

	fmt.Println("Found " + strconv.Itoa(found) + " Passwords")
	return strconv.Itoa(found), nil
}
