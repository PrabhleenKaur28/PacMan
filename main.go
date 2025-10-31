package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
)

var maze []string

func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close() 

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}
	
	return nil
}
func printMaze() {
	for _, line := range maze { 
		fmt.Println(line)
	}
}

func main() {
    // initialize game

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

        // process movement

        // process collisions

        // check game over

        // Temp: break infinite loop
        break

        // repeat
    }
}