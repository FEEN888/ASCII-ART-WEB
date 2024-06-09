package asciiart

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func AsciiTable(input, filename string) (string, error) {
	lines := strings.Split(input, "\n") // Split the input by \n to handle multi-line input
	var result strings.Builder

	for i, line := range lines {
		lineResult, err := processLine(line, filename)
		if err != nil {
			return "", err
		}
		result.WriteString(lineResult)
		if i < len(lines)-1 {
			result.WriteString("\n") // Add a newline between lines but not after the last one
		}
	}

	return result.String(), nil
}

func processLine(line, filename string) (string, error) {
	str := []rune(line)
	lnum := []int{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	// Calculate the line numbers for each character
	for i := 0; i < len(str); i++ {
		if int(str[i]) < 32 || int(str[i]) > 126 {
			return "", fmt.Errorf("character out of range: %v", str[i])
		}
		fline := ((int(str[i]) - 32) * 9) + 2
		lnum = append(lnum, fline)
	}

	var parts [][]int
	var currentPart []int
	for _, num := range lnum {
		if num == 0 && len(currentPart) > 0 && currentPart[len(currentPart)-1] == 0 {
			parts = append(parts, currentPart)
			currentPart = []int{}
		} else {
			currentPart = append(currentPart, num)
		}
	}
	if len(currentPart) > 0 {
		parts = append(parts, currentPart)
	}

	var result strings.Builder
	for _, part := range parts {
		if EqualToZero(part) {
			result.WriteString("\n")
		} else {
			tableOutput, err := Table(part, data)
			if err != nil {
				return "", err
			}
			result.WriteString(tableOutput)
		}
	}
	if checkLastElement(parts) {
		result.WriteString("\n")
	}
	return result.String(), nil
}

func Table(lnum []int, data []byte) (string, error) {
	var result strings.Builder
	text := string(data)
	lines := strings.Split(text, "\n")
	if len(lines) < 9 {
		return "", fmt.Errorf("banner file does not contain enough lines")
	}

	for k := 0; k < 8; k++ {
		for _, lineNum := range lnum {
			if lineNum != 0 && lineNum-1 < len(lines) {
				result.WriteString(strings.Replace(lines[lineNum-1], "\r", "", -1))
			} else {
				return "", fmt.Errorf("line number out of range: %d (total lines: %d)", lineNum, len(lines))
			}
		}
		if k < 7 {
			result.WriteString("\n") // Add a newline after each line except the last one
		}
		for j := 0; j < len(lnum); j++ {
			if lnum[j] != 0 {
				lnum[j]++
			}
		}
	}
	return result.String(), nil
}

func EqualToZero(arr []int) bool {
	if len(arr) != 1 {
		return false
	}
	return arr[0] == 0
}

func checkLastElement(arrays [][]int) bool {
	if len(arrays) == 0 {
		return false
	}
	lastArray := arrays[len(arrays)-1]
	if len(lastArray) == 1 {
		return false
	}
	if len(lastArray) == 0 {
		return false
	}

	lastElement := lastArray[len(lastArray)-1]
	return lastElement == 0
}
