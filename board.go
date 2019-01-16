package go_in_go


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

// A group of stone
type Group []Point

// A grid that contain the stones (usually 9x9, 11x11 or 19x19)
type Board struct {
	size int
	grid [][]byte
}

// Create a new board with a square grid of some size
func NewBoard(size int) *Board {
	initialGrid := make([][]byte, size)
	for i := range initialGrid {
		initialGrid[i] = make([]byte, size)
	}
	return &Board{size, initialGrid}
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
            if n.value == p.value && !containPoint(group, n) {
                group = append(group, n)
                neighbourGroup := b.UpdateGroup(group)
                for _, v := range neighbourGroup {
                    if !containPoint(group, v) {
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
func (b *Board) GroupLiberty(group Group) (liberty int) {
    for _, s := range group {
        liberty += b.StoneLiberty(s.y, s.x)
    }
    return
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

func (b *Board) String() (repr string) {
    for _, r := range b.grid {
        //repr += i
        repr += "\n"
        for _, v := range r {
            switch v {
            case Empty:
                repr += "."
            case Black:
                repr += "X"
            case White:
                repr += "O"
            }
            repr += " "
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

func containPoint(g Group, p Point) bool {
    for _, e := range g {
        if e == p {
            return true
        }
    }
    return false
}
