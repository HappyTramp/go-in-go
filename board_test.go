package go_in_go

import (
    "fmt"
    "testing"
)

func TestNewBoard(t *testing.T) {
    tt := []struct{
        name string
        value *Board
        expected_size uint
    }{
        {"9x9 board", NewBoard(9), 9},
        {"11x11 board", NewBoard(11), 11},
        {"19x19 board", NewBoard(19), 19},
    }
    for _, tc := range tt {
        if tc.value.size != tc.expected_size {
            t.Errorf("%s size != %d, is %d",
                     tc.name, tc.expected_size, tc.value.size,)
        }
        for i := range tc.value.grid {
            for j := range tc.value.grid[i] {
                if tc.value.grid[i][j] != 0 {
                    t.Errorf("%s grid at [%d, %d] is %d, should be 0",
                             tc.name, i, j, tc.value.grid[i][j])
                }
            }
        }
    }
}

func TestGetPoint(t *testing.T) {
    board := NewBoard(9)
    board.grid[0][0] = 2
    board.grid[1][1] = 1
    tt := []struct{
        y, x uint
        value byte
    }{
        {0, 0, 2},
        {1, 1, 1},
        {3, 3, 0},
    }
    for _, tc := range tt {
        p_value, err := board.GetPoint(tc.y, tc.x)
        if err != nil {
            t.Fatal(err)
        }
        if p_value != tc.value {
            t.Errorf("Point at [%d, %d] != %d, is %d",
                     tc.y, tc.x, tc.value, p_value)
        }
    }
}

func TestSetPoint(t *testing.T) {
    board := NewBoard(9)
    tt := []struct{
        y, x uint
        value byte
    }{
        {0, 0, 2},
        {1, 1, 1},
        {3, 3, 2},
        {3, 3, 0},
    }
    for _, tc := range tt {
        err := board.SetPoint(tc.y, tc.x, tc.value)
        if err != nil {
            t.Fatal(err)
        }
        p_value := board.grid[tc.y][tc.x]
        if p_value != tc.value {
            t.Errorf("Point at [%d, %d] should have been set to %d, is %d",
                     tc.y, tc.x, tc.value, p_value)
        }
    }
}

func TestCorrectIndex(t *testing.T) {
    tt := []struct{
        y, x uint
        board *Board
        description string
    }{
        {0, 9, NewBoard(9), "[0, 9]"},
        {9, 0, NewBoard(9), "[9, 0]"},
        {0, 19, NewBoard(19), "[0, 19]"},
        {11, 0, NewBoard(11), "[11, 0]"},
    }
    for _, tc := range tt {
        err := tc.board.correctIndex(tc.y, tc.x)
        if err == nil {
            t.Errorf("Index %s should return an error", tc.description)
        }
        expected_format := fmt.Sprintf(
            "BoardIndexError: %s is not a valid index",
            tc.description)
        if err.Error() != expected_format {
            t.Errorf("Error %s is not formated correctly", tc.description)
        }
    }
}
