package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/schollz/progressbar/v3"
)

func main() {
	listparse := app.New()
	main_window := listparse.NewWindow("listparse")
	main_window.Resize(fyne.NewSize(600, 600))

	// Widgets
	wordlist_entry := widget.NewEntry()
	fileOpen := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		wordlist_entry.SetText(uc.URI().Path())
	}, main_window)
	wordlist_button := widget.NewButton("open", func() {
		fileOpen.Show()
		fileOpen.Resize(fyne.NewSize(600, 600))
	})

	output_entry := widget.NewEntry()
	outputSave := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		output_entry.SetText(uc.URI().Path())
	}, main_window)
	output_button := widget.NewButton("open", func() {
		outputSave.Show()
		outputSave.Resize(fyne.NewSize(600, 600))
	})

	minLen_entry := widget.NewEntry()
	maxLen_entry := widget.NewEntry()
	phrase_entry := widget.NewEntry()
	specialChar_check := widget.NewCheck("Special Character", func(b bool) {})
	number_check := widget.NewCheck("Number", func(b bool) {})

	progress_widget := widget.NewProgressBarInfinite()

	error_text := widget.NewLabel("")

	found_text := widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Wordlist Path: ", Widget: wordlist_entry},
			{Text: "", Widget: wordlist_button},
			{Text: "Output Path: ", Widget: output_entry},
			{Text: "", Widget: output_button},
			{Text: "Minimum Length: ", Widget: minLen_entry},
			{Text: "Maximum Length: ", Widget: maxLen_entry},
			{Text: "Phrase: ", Widget: phrase_entry},
			{Text: "Special Character: ", Widget: specialChar_check},
			{Text: "Number: ", Widget: number_check},
		},
		OnSubmit: func() {
			progress_widget.Show()

			if minLen_entry.Text == "" {
				minLen_entry.SetText("0")
			}

			if maxLen_entry.Text == "" {
				maxLen_entry.SetText("0")
			}

			count, err := parse(wordlist_entry.Text, output_entry.Text, minLen_entry.Text, maxLen_entry.Text, phrase_entry.Text, specialChar_check.Checked, number_check.Checked, false)
			if err != nil {
				error_text.SetText(err.Error())
				progress_widget.Hide()
			}
			found_text.SetText("Found " + count + " passwords!")
			progress_widget.Hide()
		},
		OnCancel: func() {
			listparse.Quit()
		},
	}

	progress_widget.Hide()
	content := container.New(layout.NewVBoxLayout(), form, found_text, error_text, widget.NewLabel(""), widget.NewLabel(""), progress_widget)
	main_window.SetContent(content)
	main_window.ShowAndRun()
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
