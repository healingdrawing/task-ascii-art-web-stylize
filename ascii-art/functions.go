package ascii_art

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

//go:embed banners/*
var f embed.FS

var BannerFiles = map[string]string{
	"standard":   "banners/standard.txt",
	"thinkertoy": "banners/thinkertoy.txt",
	"shadow":     "banners/shadow.txt"}

var HeadString = "Avoid to use next\n[character](number):\n\n"
var warningString = "[%c](%v)\n"

// check string include not supported characters
func checkText(words []string) (out string) {
	urune := map[rune]int{}
	for _, word := range words {
		for _, r := range word {
			ind := int(r)
			if ind < 32 || ind > 127 { // not printable character
				_, ok := urune[r]
				if !ok {
					urune[r] = 0
					out += fmt.Sprintf(warningString, r, r)
				}
			}
		}
	}
	if out != "" {
		out = HeadString + out
	}
	return
}

// returns ascii-art, warning message, error
func AsciiToString(text string, bannerFile string) (string, string, error) {

	text = strings.ReplaceAll(text, "\r", "") // remove it to properly manage allowed symbols
	// log.Println(text) // uncomment this to print raw input text

	lines, err := readFileIntoSlice(bannerFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("banner '%v' not found", bannerFile)
		}
		return "", "", err
	}

	splittwo := "\n"
	words := strings.Split(text, splittwo)

	if text == "" || text == "\n" {
		return text, "", nil
	}

	out := checkText(words)

	art := ""
	for _, word := range words { // nested loop to print line by line depending on input.
		if word == "" { // the new line "\\n" was at the end of "words" slice, and Split create the "" word
			art += fmt.Sprintln()
		} else { // usual case letter print
			// vertical step to print horizontal sequences of letter ascii art
			for h := 1; h < 9; h++ { // from one to ignore the empty new line from standart.txt
				for _, l := range word {
					ind := (int(l)-32)*9 + h         // potential index (the height from up to bottom) in "lines" for required letter line(because art letter is multilined)
					if ind > 0 && ind < len(lines) { // check the index is inside available ascii art symbols ... f.e. standart.txt
						art += fmt.Sprint(lines[ind]) // print the line from high "h" for the word letter "l"
					}
				}
				art += fmt.Sprintln()
			}
		}
	}
	return art, out, nil
}

func readFileIntoSlice(name string) ([]string, error) {
	file, err := f.ReadFile(name)

	if err != nil {
		return nil, err
	}

	clearFile := strings.ReplaceAll(string(file), "\r\n", "\n")
	lines := strings.Split(clearFile, "\n")
	return lines, nil
}
