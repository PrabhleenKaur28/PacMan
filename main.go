package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func initialise() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin
	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("Unable to activate cbreak mode:", err)
	}
}

func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin
	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("Unable to restore cooked mode:", err)
	}
}

var maze []string

func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maze = nil
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	return nil
}
func printMaze() {
	ClearScreen()
	for _, line := range maze {
		fmt.Println(line)
	}
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	n, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}
	log.Printf("Read %d bytes: %v\n", n, buffer[:n])
	if n == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}
	return "", nil
}

func ClearScreen() {
	fmt.Print("\x1b[2J")
	MoveCursor(0, 0)
}

func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row+1, col+1)
}

const reset = "\x1b[0m"

type Color int

const (
	BLACK Color = iota
	RED
	GREEN
	BROWN
	BLUE
	MAGENTA
	CYAN
	GREY
)

var colors = map[Color]string{
	BLACK:   "\x1b[1;30;40m", // ANSI escape code for black background
	RED:     "\x1b[1;31;41m", //bold red text on red background
	GREEN:   "\x1b[1;32;42m",
	BROWN:   "\x1b[1;33;43m",
	BLUE:    "\x1b[1;34;44m",
	MAGENTA: "\x1b[1;35;45m",
	CYAN:    "\x1b[1;36;46m",
	GREY:    "\x1b[1;37;47m",
}

func WithBlueBackground(text string) string {
	return "\x1b[44m" + text + reset
}

func WithBackground(text string, color Color) string {
	if c, ok := colors[color]; ok {
		return c + text + reset
	}
	return WithBlueBackground(text)
}

func main() {
	// initialize game
	initialise()
	defer cleanup()

	// load resources
	err := loadMaze("maze1.txt")
	if err != nil {
		log.Println("Failed to load maze:", err)
		return
	}

	// game loop
	for {
		// update screen
		printMaze()

		// process input
		input, err := readInput()
		if err != nil {
			log.Println("Error reading input:", err)
			break
		}
		if input == "ESC" {
			log.Println("HELLOOO")
			break
		}

		// process movement

		// process collisions

		// check game over

		// break infinite loop

		// repeat
	}
}
