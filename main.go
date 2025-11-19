package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"math/rand"
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

	for row, line := range maze {
		for col, ch := range line {
			switch ch {
			case 'P':
				player = sprite{row, col}	
			case 'G':
				ghosts = append(ghosts, &sprite{row, col} )				
			}
		}
	}

	return nil
}
func printMaze() {
	ClearScreen()
	for _, line := range maze {
		for _, ch := range line {
			switch ch {
			case '#':
				fmt.Printf("%c", ch)
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	MoveCursor(player.row, player.col)
	fmt.Print("P")

	for _, ghost := range ghosts {
		MoveCursor(ghost.row, ghost.col)
		fmt.Print("G")
	}

	MoveCursor(len(maze) + 1, 0) // moving cursor outside maze
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	n, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}
	if n == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}else if n >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}	
	}
	return "", nil
}

func makeMove(oldRow, oldCol int, direction string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch direction {
	case "UP":
		newRow--
		if newRow < 0{
			newRow = len(maze)-1
		}
	case "DOWN":
		newRow++
		if newRow == len(maze){
			newRow = 0
		}
	case "LEFT":
		newCol--
		if newCol < 0{
			newCol = len(maze[0])-1
		}
	case "RIGHT":
		newCol++
		if newCol == len(maze[0]){
			newCol = 0
		}		
	}

	if maze[newRow][newCol] == '#' { // wall (collision)
		newRow, newCol = oldRow, oldCol
	}
	return
}

func movePlayer(direction string) {
	player.row, player.col = makeMove(player.row, player.col, direction)
}

func getRandomDirection() string {
	direction := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "LEFT",
		3: "RIGHT",
	}
	return move[direction] 
}

func moveGhosts() {
	for _, ghost := range ghosts {
		direction := getRandomDirection()
		ghost.row, ghost.col = makeMove(ghost.row, ghost.col, direction)
	}
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

type sprite struct{
	row, col int
}
var player sprite
var ghosts []*sprite

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
		movePlayer(input)
		moveGhosts()

		if input == "ESC" {
			break
		}

		// process movement

		// process collisions

		// check game over

		// break infinite loop

		// repeat
	}
}
