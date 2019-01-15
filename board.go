package go_in_go


const (
	EMPTY byte = 0
	BLACK byte = 1
	WHITE byte = 2
)

type Point struct {
    x, y int
    value byte
}

type Board struct {
	size int
	grid [][]byte
}

func NewBoard(size int) *Board {
	initial_grid := make([][]byte, size)
	for i := range initial_grid {
		initial_grid[i] = make([]byte, size)
	}
	return &Board{size, initial_grid}
}

func (b *Board) GetPoint(y, x int) (byte, error) {
    if !b.correctIndex(y, x) {
        return 200, &BoardIndexError{y, x}
    }
    return b.grid[y][x], nil
}

func (b *Board) SetPoint(y, x int, value byte) error {
    if !b.correctIndex(y, x) {
        return &BoardIndexError{y, x}
    }
    b.grid[y][x] = value
    return nil
}

func (b *Board) Neighbours(y, x int) []Point {
    neighbours := make([]Point, 0)
    for _, mod := range []int{-1, 1} {
        mod_y := y + mod
        mod_x := x + mod
        mod_y_value, mod_y_err := b.GetPoint(mod_y, x)
        mod_x_value, mod_x_err := b.GetPoint(y, mod_x)
        if mod_y_err == nil {
            neighbours = append(neighbours, Point{mod_y, x, mod_y_value})
        }
        if mod_x_err == nil {
            neighbours = append(neighbours, Point{y, mod_x, mod_x_value})
        }
    }
    return neighbours
}

func (b *Board) correctIndex(y, x int) bool {
    if y < 0 || y > b.size - 1 || x < 0 || x > b.size - 1 {
        return false
    }
    return true
}
