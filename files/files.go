package files

import (
	"bufio"
	"fmt"
	"os"
)

func FileScanner(filename string) (*bufio.Scanner, error) {
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}

	// defer file.Close()

	return bufio.NewScanner(file), nil
}

func FileReader(filename string) (*bufio.Reader, error) {
	fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}

	// defer file.Close()

	return bufio.NewReader(file), nil
}
