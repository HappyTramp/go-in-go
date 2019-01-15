package go_in_go


const (
	EMPTY byte = 0
	BLACK byte = 1
	WHITE byte = 2
)

type Board struct {
	size uint
	grid [][]byte
}

func NewBoard(size uint) *Board {
	initial_grid := make([][]byte, size)
	for i := range initial_grid {
		initial_grid[i] = make([]byte, size)
	}
	return &Board{size, initial_grid}
}

func (b *Board) GetPoint(y, x uint) (byte, error) {
    if err := b.correctIndex(y, x); err != nil {
        return 200, err
    }
    return b.grid[y][x], nil
}

func (b *Board) SetPoint(y, x uint, value byte) error {
    if err := b.correctIndex(y, x); err != nil {
        return err
    }
    b.grid[y][x] = value
    return nil
}

func (b *Board) correctIndex(y, x uint) error {
    if y < 0 || y > b.size - 1 || x < 0 || x > b.size - 1 {
        return &BoardIndexError{y, x}
    }
    return nil
}
