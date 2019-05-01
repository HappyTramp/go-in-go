package go_in_go

import (
    "fmt"
)

/* TODO:
 * Ko rule detection and prevention
 * Group suicide prevention
 * Score counting
 * game flow
 */

// The three state a point can be in
const (
	Empty byte = 0
	Black byte = 1
	White byte = 2
)

// A point is defined by his position and value
type Point struct {
    y, x int
    value byte
}

// Group of stone
type Group []Point

// A board that contain the stones (usually 9x9, 11x11 or 19x19)
type Board struct {
	size int
	grid [][]byte
    player1Score, player2Score int
}


// Create a new board with a square grid of some size
func NewBoard(size int) *Board {
	initialGrid := make([][]byte, size)
	for i := range initialGrid {
		initialGrid[i] = make([]byte, size)
	}
	return &Board{size, initialGrid, 0, 0}
}


// Get a point at an index (BoardIndexError if index off boundary)
func (b *Board) GetPoint(y, x int) (byte, error) {
    if !b.correctIndex(y, x) {
        return 200, &BoardIndexError{y, x}
    }
    return b.grid[y][x], nil
}


// Set a point at an index (BoardIndexError if index off boundary)
func (b *Board) SetPoint(y, x int, value byte) error {
    if !b.correctIndex(y, x) {
        return &BoardIndexError{y, x}
    }
    b.grid[y][x] = value
    return nil
}


// Slice of points neighbours to a point
func (b *Board) Neighbours(y, x int) Group {
    neighbours := make(Group, 0)
    for _, mod := range []int{-1, 1} {
        mY, mX := y + mod, x + mod
        mYValue, mYErr := b.GetPoint(mY, x)
        mXValue, mXErr := b.GetPoint(y, mX)
        if mYErr == nil {
            neighbours = append(neighbours, Point{mY, x, mYValue})
        }
        if mXErr == nil {
            neighbours = append(neighbours, Point{y, mX, mXValue})
        }
    }
    return neighbours
}


// Update a group of stone
func (b *Board) UpdateGroup(group Group) Group {
    for _, p := range group {
        for _, n := range b.Neighbours(p.y, p.x) {
            if n.value == p.value && !group.contain(n) {
                group = append(group, n)
                neighbourGroup := b.UpdateGroup(group)
                for _, v := range neighbourGroup {
                    if !group.contain(v) {
                        group = append(group, v)
                    }
                }
            }
        }
    }
    return group
}


// Return the group from the position
func (b *Board) GroupFrom(y, x int) Group {
    initialStone, err := b.GetPoint(y, x)
    if err != nil {
        return Group{}
    }
    return b.UpdateGroup(Group{{y, x, initialStone}})
}


// Liberty of a group of stone
func (b *Board) GroupLiberty(group Group) int {
    libertyPoints := Group{}
    for _, s := range group {
        for _, n := range b.Neighbours(s.y, s.x) {
            if n.value == Empty && !libertyPoints.contain(s) {
                libertyPoints = append(libertyPoints, n)

            }
        }
    }
    return len(libertyPoints)
}


// Liberty of one stone
func (b *Board) StoneLiberty(y, x int) (liberty int) {
    for _, n := range b.Neighbours(y, x) {
        if n.value == Empty {
            liberty += 1
        }
    }
    return
}


// Delete an entire group of stone
func (b *Board) DeleteGroup(group Group) error {
    for _, s := range group {
        err := b.SetPoint(s.y, s.x, Empty)
        if err != nil {
            return err
        }
    }
    return nil
}


func (b *Board) UpdateScore(player, points int) {
    switch player {
    case 1:
        b.player1Score += points
    case 2:
        b.player2Score += points
    }
}


// Detect if a move result in a ko
func (b *Board) KoDetection(previousBoard Board, move Point) bool {
    nextBoard := b.clone()
    nextBoard.SetPoint(move.y, move.x, move.value)
    for i, row := range previousBoard.grid {
        for j, p := range row {
            if v, _ := nextBoard.GetPoint(i, j); v != p {
                return false
            }
        }
    }
    return true
}


const notationLetters string = "ABCDEFGHJKLMNOPQRST"
// String representation of Board with algebreaic notation
// for terminal purposes
func (b *Board) String() (repr string) {
    repr += "\n  "
    for i := 0; i < b.size; i++ {
        repr += fmt.Sprintf(" %c", notationLetters[i])
    }
    for i, r := range b.grid {
        repr += fmt.Sprintf("\n%-2d", i + 1)
        for _, v := range r {
            repr += " "
            switch v {
            case Empty:
                repr += "."
            case Black:
                repr += "X"
            case White:
                repr += "O"
            }
        }
    }
    return
}


// Verify that the index is in boundary
func (b *Board) correctIndex(y, x int) bool {
    if y < 0 || y > b.size - 1 || x < 0 || x > b.size - 1 {
        return false
    }
    return true
}


// Check if a group contain a point
func (g Group) contain(p Point) bool {
    for _, e := range g {
        if e == p {
            return true
        }
    }
    return false
}


// Clone/Copy a board
func (b *Board) clone() Board {
    gridClone := make([][]byte, b.size)
    for i := range b.grid {
        gridClone[i] = make([]byte, b.size)
        copy(gridClone[i], b.grid[i])
    }
    return Board{b.size, gridClone, b.player1Score, b.player2Score}
}
