package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	HEIGHT = 28 // Height in board (1~20 line is showed)
	WIDTH  = 18 // Width in board (1~10 line is showed)
	BLACK  = 9632
	WHITE  = 9633
	UP     = 0
	DOWN   = 1
	RIGHT  = 2
	LEFT   = 3
)

const (
	I = 0 + iota
	L
	J
	T
	SQUARE // square
)

var dy = [4]int{0, 1, 0, 0}
var dx = [4]int{0, 0, 1, -1}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// Vertex is position [y][x]
type Vertex struct {
	y int
	x int
}

// Block has shape I,L,T
type Block struct {
	shape int
	pos   [4]Vertex
}

func (b *Block) init() bool {

	b.shape = rand.Intn(5) // TODO : random
	switch b.shape {
	case I:
		b.pos[0].y, b.pos[0].x = 4, 8
		b.pos[1].y, b.pos[1].x = 5, 8
		b.pos[2].y, b.pos[2].x = 6, 8
		b.pos[3].y, b.pos[3].x = 7, 8
	case L:
		b.pos[0].y, b.pos[0].x = 6, 8
		b.pos[1].y, b.pos[1].x = 6, 9
		b.pos[2].y, b.pos[2].x = 5, 8
		b.pos[3].y, b.pos[3].x = 4, 8
	case J:
		b.pos[0].y, b.pos[0].x = 6, 9
		b.pos[1].y, b.pos[1].x = 6, 8
		b.pos[2].y, b.pos[2].x = 5, 9
		b.pos[3].y, b.pos[3].x = 4, 9
	case T:
		b.pos[0].y, b.pos[0].x = 5, 8
		b.pos[1].y, b.pos[1].x = 5, 7
		b.pos[2].y, b.pos[2].x = 5, 9
		b.pos[3].y, b.pos[3].x = 4, 8
	case SQUARE:
		b.pos[0].y, b.pos[0].x = 4, 8
		b.pos[1].y, b.pos[1].x = 4, 9
		b.pos[2].y, b.pos[2].x = 5, 8
		b.pos[3].y, b.pos[3].x = 5, 9
	}
	for i := 0; i < 4; i++ {
		if board[b.pos[i].y][b.pos[i].x] == BLACK {
			return false
		}
	}
	show()
	return true
}

func (b *Block) autoDown() bool {
	// can go down -> true
	// can't go down -> false
	if !b.canGo(DOWN) {
		return false
	}
	for i := 0; i < 4; i++ {
		b.pos[i].y++
	}
	show()
	return true
}
func (b *Block) down(c chan bool) {
	for {
		<-c
		mux.Lock()
		if !b.canGo(DOWN) {
			mux.Unlock()
			continue
		}
		for i := 0; i < 4; i++ {
			b.pos[i].y++
		}
		show()
		mux.Unlock()
	}
}

func (b *Block) rotate(c chan bool) {
	for {
		<-c
		mux.Lock()
		if b.shape == SQUARE {
			mux.Unlock()
			continue
		}
		if !b.canRotate() {
			mux.Unlock()
			continue
		}
		for i := 0; i < 4; i++ {
			tempX := curr.pos[i].x - curr.pos[0].x
			tempY := curr.pos[i].y - curr.pos[0].y
			curr.pos[i].x = -tempY + curr.pos[0].x
			curr.pos[i].y = tempX + curr.pos[0].y

		}

		mux.Unlock()
	}
}
func (b *Block) left(c chan bool) {
	for {
		<-c
		mux.Lock()
		if !b.canGo(LEFT) {
			mux.Unlock()
			continue
		}
		for i := 0; i < 4; i++ {
			b.pos[i].x--
		}
		show()
		mux.Unlock()
	}
}
func (b *Block) right(c chan bool) {
	for {
		<-c
		mux.Lock()
		if !b.canGo(RIGHT) {
			mux.Unlock()
			continue
		}
		for i := 0; i < 4; i++ {
			b.pos[i].x++
		}
		show()
		mux.Unlock()
	}
}

func (b *Block) space(c chan bool) {
	for {
		<-c
		mux.Lock()
		for {
			if !b.canGo(DOWN) {
				break
			}
			for i := 0; i < 4; i++ {
				b.pos[i].y++
			}
		}

		lowest := 0
		for i := 0; i < 4; i++ {
			board[curr.pos[i].y][curr.pos[i].x] = BLACK
			if lowest < curr.pos[i].y {
				lowest = curr.pos[i].y
			}
		}
		lineClear(lowest)
		if !curr.init() {
			mux.Unlock()
			showEnd()
			for {

			}

		}
		mux.Unlock()
	}
}
func (b *Block) canGo(direction int) bool {

	for i := 0; i < 4; i++ {
		nextY := b.pos[i].y + dy[direction]
		nextX := b.pos[i].x + dx[direction]
		if board[nextY][nextX] == BLACK {
			return false
		}
	}
	return true
}
func (b *Block) canRotate() bool {
	// 이거하려면 사이즈를 재조정해야할듯.. 주위로 4칸씩 블랙으로 감자
	for i := 0; i < 4; i++ {
		tempX := curr.pos[i].x - curr.pos[0].x
		tempY := curr.pos[i].y - curr.pos[0].y
		nextX := -tempY + curr.pos[0].x
		nextY := tempX + curr.pos[0].y
		if board[nextY][nextX] == BLACK {
			return false
		}
	}
	return true
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////

var board [HEIGHT][WIDTH]rune
var curr Block
var mux sync.Mutex

func main() {

	start()
	curr.init()
	go getInput()

	for {
		time.Sleep(500 * time.Millisecond)
		// 2초마다 일어나는 동작 수행
		mux.Lock()
		if !curr.autoDown() {
			lowest := 0
			for i := 0; i < 4; i++ {
				board[curr.pos[i].y][curr.pos[i].x] = BLACK
				if lowest < curr.pos[i].y {
					lowest = curr.pos[i].y
				}
			}
			lineClear(lowest)
			if !curr.init() {
				mux.Unlock()
				showEnd()
				for {

				}

			}
		}
		mux.Unlock()
	}

}
func lineClear(lowest int) {
	deleted := 0
	for y := lowest; y >= 4; y-- {
		complete := true
		for x := 4; x < showW; x++ {
			if board[y][x] == WHITE {
				complete = false
				break
			}
		}
		if complete {
			deleted++
		} else {
			for x := 4; x < showW; x++ {
				board[y+deleted][x] = board[y][x]
			}
		}
	}
}
func getInput() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	input := make([]byte, 4)

	cUp := make(chan bool)
	cDown := make(chan bool)
	cRight := make(chan bool)
	cLeft := make(chan bool)
	cSPACE := make(chan bool)
	go curr.rotate(cUp)
	go curr.down(cDown)
	go curr.right(cRight)
	go curr.left(cLeft)
	go curr.space(cSPACE)
	for {
		os.Stdin.Read(input)
		if input[0] == ' ' {
			cSPACE <- true
			continue
		}
		switch input[2] - 65 {
		case UP:
			cUp <- true
		case DOWN:
			cDown <- true
		case RIGHT:
			cRight <- true
		case LEFT:
			cLeft <- true
		}
	}
}

func start() {
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			board[y][x] = WHITE
		}
	}
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < 4; x++ {
			board[y][x] = BLACK
		}
		for x := WIDTH - 4; x < WIDTH; x++ {
			board[y][x] = BLACK
		}
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < WIDTH; x++ {
			board[y][x] = BLACK
		}
	}
	for y := HEIGHT - 4; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			board[y][x] = BLACK
		}
	}
}

const showH = HEIGHT - 4
const showW = WIDTH - 4

func show() {
	fmt.Println("*************************************")
	for i := 0; i < 4; i++ {
		board[curr.pos[i].y][curr.pos[i].x] = BLACK
	}
	for y := 4; y < showH; y++ {
		for x := 4; x < showW; x++ {
			fmt.Printf("%c ", board[y][x])
		}
		fmt.Printf("\n")
	}
	for i := 0; i < 4; i++ {
		board[curr.pos[i].y][curr.pos[i].x] = WHITE
	}
}

func showEnd() {

	for y := 4; y < showH; y++ {
		for x := 4; x < showW; x++ {
			board[y][x] = BLACK
		}
	}
	for x := 5; x < showW-1; x++ {
		board[7][x] = WHITE
	}
	for x := 4; x < showW; x++ {
		board[13][x] = WHITE
	}
	for x := 6; x < showW-1; x++ {
		board[16][x] = WHITE
	}
	for x := 6; x < showW-1; x++ {
		board[18][x] = WHITE
	}
	for x := 6; x < showW-1; x++ {
		board[20][x] = WHITE
	}
	for y := 8; y <= 10; y++ {
		board[y][8] = WHITE
	}
	for y := 8; y <= 10; y++ {
		board[y][12] = WHITE
	}
	board[17][6] = WHITE
	board[19][6] = WHITE
	fmt.Println("*** get out a here!! (ctrl + c) ***")
	for y := 4; y < showH; y++ {
		for x := 4; x < showW; x++ {
			fmt.Printf("%c ", board[y][x])
		}
		fmt.Printf("\n")
	}
}
