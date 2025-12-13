package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
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

type sprite struct {
	row, col           int
	startRow, startCol int
}

var maze []string
var player sprite
var ghosts []*sprite

var score int
var numDots int
var lives = 3

var currentDirection string = ""

type Config struct {
	Player   string `json:"player"`
	Ghost    string `json:"ghost"`
	Wall     string `json:"wall"`
	Dot      string `json:"dot"`
	Pill     string `json:"pill"`
	Death    string `json:"death"`
	Space    string `json:"space"`
	UseEmoji bool   `json:"use_emoji"`
}

var cfg Config

const dataDir = "/usr/share/pacman"

func loadConfig(file string) error {
	path := filepath.Join(dataDir, file)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func loadMaze(file string) error {
	path := filepath.Join(dataDir, file)

	f, err := os.Open(path)
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
				player = sprite{row, col, row, col}
			case 'G':
				ghosts = append(ghosts, &sprite{row, col, row, col})
			case '.':
				numDots++
			}
		}
	}

	return nil
}

func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row+1, col+1)
}

func moveCursor(row, col int) {
	if cfg.UseEmoji {
		MoveCursor(row, col*2)
	} else {
		MoveCursor(row, col)
	}
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

func ClearScreen() {
	fmt.Print("\x1b[2J")
	MoveCursor(0, 0)
}

func printMaze() {
	ClearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Print(WithBackground(cfg.Wall, BLUE))
			case '.':
				fmt.Print(cfg.Dot)
			case 'X':
				fmt.Print(cfg.Pill)
			default:
				fmt.Print(cfg.Space)
			}
		}
		fmt.Println()
	}
}

func eraseSprite(row, col int) {
	moveCursor(row, col)
	chr := rune(maze[row][col])
	switch chr {
	case '#':
		fmt.Print(WithBackground(cfg.Wall, CYAN))
	case '.':
		fmt.Print(cfg.Dot)
	case 'X':
		fmt.Print(cfg.Pill)
	default:
		fmt.Print(cfg.Space)
	}
}

func drawSprite(row, col int, sprite string) {
	moveCursor(row, col)
	fmt.Print(sprite)
}

func updateStatus() {
	moveCursor(len(maze)+1, 0)
	fmt.Print("\x1b[K") // Clear the line
	fmt.Print("Score:", score, "\tLives:", lives)
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	n, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}
	if n == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if n >= 3 { //arrow key escape sequence is 3 bytes
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
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow++
		if newRow == len(maze) {
			newRow = 0
		}
	case "LEFT":
		newCol--
		if newCol < 0 {
			newCol = len(maze[oldRow]) - 1
		}
	case "RIGHT":
		newCol++
		if newCol == len(maze[oldRow]) {
			newCol = 0
		}
	}

	if maze[newRow][newCol] == '#' { // wall (collision)
		newRow, newCol = oldRow, oldCol
	}
	return
}

func movePlayer(dir string) {
	oldRow, oldCol := player.row, player.col
	player.row, player.col = makeMove(player.row, player.col, dir)

	if oldRow != player.row || oldCol != player.col {
		eraseSprite(oldRow, oldCol)
		drawSprite(player.row, player.col, cfg.Player)
	}

	removeDot := func(row, col int) {
		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
	}

	switch maze[player.row][player.col] {
	case '.':
		numDots--
		score++
		removeDot(player.row, player.col)
		updateStatus()
	case 'X':
		score += 10
		removeDot(player.row, player.col)
		updateStatus()
	}
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
	type position struct {
		row, col int
	}
	oldPositions := make([]position, len(ghosts))

	for i, ghost := range ghosts {
		oldPositions[i] = position{ghost.row, ghost.col}
		ghost.row, ghost.col = makeMove(ghost.row, ghost.col, getRandomDirection())
	}

	for i, ghost := range ghosts {
		if oldPositions[i].row != ghost.row || oldPositions[i].col != ghost.col {
			eraseSprite(oldPositions[i].row, oldPositions[i].col)
		}
	}

	for _, ghost := range ghosts {
		drawSprite(ghost.row, ghost.col, cfg.Ghost)
	}

	moveCursor(len(maze)+2, 0)
}

func main() {
	// initialize game
	initialise()
	defer cleanup()

	// load resources
	var prompt = promptui.Select{
		Label: "Choose Level",
		Items: []string{"Easy", "Medium", "Hard"},
	}
	_, result, _ := prompt.Run()

	mazeFile := map[string]string{
		"Easy":   "maze1.txt",
		"Medium": "maze2.txt",
		"Hard":   "maze3.txt",
	}[result]

	loadMaze(mazeFile)

	err := loadConfig("emojis.json")
	if err != nil {
		log.Println("failed to load configuration:", err)
		return
	}

	printMaze()
	drawSprite(player.row, player.col, cfg.Player)
	for _, g := range ghosts {
		drawSprite(g.row, g.col, cfg.Ghost)
	}
	updateStatus()

	//process input (async)
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			inp, err := readInput()
			if err != nil {
				log.Println("Error reading input:", err)
				ch <- "ESC"
			}
			ch <- inp
		}
	}(input)

	// game loop
	for {
		// process movement
		select {
		case c := <-input:
			if c == "ESC" {
				lives = 0
			}

			if c == "UP" || c == "DOWN" || c == "LEFT" || c == "RIGHT" {
				currentDirection = c
			}
		default:
		}

		if currentDirection != "" {
			movePlayer(currentDirection)
		}

		moveGhosts()

		// process collisions
		for _, ghost := range ghosts {
			if ghost.row == player.row && ghost.col == player.col {
				lives--
				eraseSprite(player.row, player.col)
				player.row, player.col = player.startRow, player.startCol
				drawSprite(player.row, player.col, cfg.Player)
				updateStatus()
			}
		}

		// Check game over
		if numDots == 0 || lives == 0 {
			if lives == 0 {
				moveCursor(player.row, player.col)
				fmt.Print(cfg.Death)
				moveCursor(len(maze)+2, 0)
				fmt.Println("GAME OVER!!")
			} else {
				moveCursor(len(maze)+2, 0)
				fmt.Println("Congratulations! You collected all the coins!!")
			}
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}
