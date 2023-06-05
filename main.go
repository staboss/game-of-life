package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CELL_ALIVE_STATE = 'O'
	CELL_DEAD_STATE  = ' '
)

type cell struct {
	x     int
	y     int
	state rune
}

func (c *cell) getAliveNeighbours(u *universe) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			neighbourX := (c.x + i + u.size) % u.size
			neighbourY := (c.y + j + u.size) % u.size
			if !(i == 0 && j == 0) && u.maps[neighbourX][neighbourY].state == CELL_ALIVE_STATE {
				count++
			}
		}
	}
	return count
}

func (c *cell) getNextCellState(u *universe) cell {
	aliveNeighbors := c.getAliveNeighbours(u)
	nextCell := cell{x: c.x, y: c.y}

	if c.state == CELL_ALIVE_STATE {
		if aliveNeighbors == 2 || aliveNeighbors == 3 {
			nextCell.state = CELL_ALIVE_STATE
		} else {
			nextCell.state = CELL_DEAD_STATE
		}
	} else {
		if aliveNeighbors == 3 {
			nextCell.state = CELL_ALIVE_STATE
		} else {
			nextCell.state = CELL_DEAD_STATE
		}
	}

	return nextCell
}

type universe struct {
	aliveCells        int
	generationsPassed int
	maps              [][]cell
	seed              int64
	size              int
}

func createUniverse(size int, seed int64) *universe {
	var u universe
	u.initialize(size, seed)
	return &u
}

func (u *universe) initialize(size int, seed int64) {
	random := rand.New(rand.NewSource(seed))

	u.seed = seed
	u.size = size
	u.maps = make([][]cell, size)

	for i := 0; i < u.size; i++ {
		u.maps[i] = make([]cell, u.size)
		for j := 0; j < u.size; j++ {
			state := CELL_DEAD_STATE
			if random.Intn(2) == 1 {
				state = CELL_ALIVE_STATE
				u.aliveCells += 1
			}
			u.maps[i][j] = cell{
				x:     i,
				y:     j,
				state: state,
			}
		}
	}

	u.generationsPassed++
}

func (u *universe) display() {
	fmt.Printf("Generation #%d\n", u.generationsPassed)
	fmt.Printf("Alive: %d\n", u.aliveCells)

	for i := 0; i < u.size; i++ {
		for j := 0; j < u.size; j++ {
			fmt.Printf("%c", u.maps[i][j].state)
		}
		fmt.Println()
	}
}

func createNextGeneration(u *universe) *universe {
	nextGen := createUniverse(u.size, u.seed)
	tmpAliveCells := 0

	for i := 0; i < u.size; i++ {
		for j := 0; j < u.size; j++ {
			currCell := u.maps[i][j]
			nextGen.maps[i][j] = currCell.getNextCellState(u)
			if nextGen.maps[i][j].state == CELL_ALIVE_STATE {
				tmpAliveCells += 1
			}
		}
	}

	nextGen.aliveCells = tmpAliveCells
	nextGen.generationsPassed = u.generationsPassed + 1

	return nextGen
}

func main() {
	var (
		generations int
		size        int
		seed        int64
		universe    universe
	)

	r := rand.New(rand.NewSource(99))

	seed = r.Int63()
	generations = 1000

	fmt.Scan(&size)

	universe.initialize(size, seed)
	universe.display()
	for i := 1; i <= generations; i++ {
		nextGen := createNextGeneration(&universe)
		universe = *nextGen
		universe.display()

		time.Sleep(500 * time.Millisecond)
		fmt.Print("\033[H\033[2J")
	}
}
