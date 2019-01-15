package go_in_go

import (
    "testing"
    "reflect"
    "sort"
)

func TestNewBoard(t *testing.T) {
    tt := []struct{
        name string
        value *Board
        expected_size int
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
        y, x int
        value byte
        err error
    }{
        {0, 0, 2, nil},
        {1, 1, 1, nil},
        {3, 3, 0, nil},
        {10, 0, 200, &BoardIndexError{10, 0}},
        {0, 9, 200, &BoardIndexError{0, 9}},
    }
    for _, tc := range tt {
        p_value, err := board.GetPoint(tc.y, tc.x)
        if reflect.TypeOf(err) != reflect.TypeOf(tc.err) {
            t.Fatalf("Didnt return correct error got: %T, expected: %T",
                     err, tc.err)
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
        y, x int
        value byte
        err error
    }{
        {0, 0, 2, nil},
        {1, 1, 1, nil},
        {3, 3, 2, nil},
        {3, 3, 0, nil},
        {9, 0, 0, &BoardIndexError{9, 0}},
        {0, 11, 0, &BoardIndexError{0, 11}},
    }
    for _, tc := range tt {
        err := board.SetPoint(tc.y, tc.x, tc.value)
        if reflect.TypeOf(err) != reflect.TypeOf(tc.err) {
            t.Fatalf("Didnt return correct error got: %T, expected: %T",
                     err, tc.err)
        }
        if tc.err == nil {
            p_value := board.grid[tc.y][tc.x]
            if p_value != tc.value {
                t.Errorf("Point at [%d, %d] should have been set to %d, is %d",
                         tc.y, tc.x, tc.value, p_value)
            }
        }
    }
}

func TestNeighbours(t *testing.T) {
    board := NewBoard(9)
    board.SetPoint(3, 2, 1)
    board.SetPoint(2, 3, 2)
    board.SetPoint(4, 4, 2)
    tt := []struct{
        y, x int
        neighbours []Point
    }{
        {0, 0, []Point{{0, 1, 0}, {1, 0, 0}}},
        {2, 2, []Point{{2, 1, 0}, {2, 3, 2}, {1, 2, 0}, {3, 2, 1}}},
        {3, 3, []Point{{3, 2, 1}, {3, 4, 0}, {2, 3, 2}, {4, 3, 0}}},
        {3, 4, []Point{{3, 3, 0}, {3, 5, 0}, {2, 4, 0}, {4, 4, 2}}},
    }
    for _, tc := range tt {
        sort.Slice(tc.neighbours, func(i, j int) bool {
            return tc.neighbours[i].y < tc.neighbours[j].y
        })
        neighbours := board.Neighbours(tc.y, tc.x)
        sort.Slice(neighbours, func(i, j int) bool {
            return neighbours[i].y < neighbours[j].y
        })
        for i, n := range tc.neighbours {
            if n != neighbours[i] {
                t.Errorf("Wrong neighbours\ngot:      %v\nexpected: %v",
                         neighbours, tc.neighbours)
            }
        }
    }
}

func TestCorrectIndex(t *testing.T) {
    tt := []struct{
        y, x int
        board *Board
        value bool
    }{
        {0, 4, NewBoard(9), true},
        {5, 0, NewBoard(9), true},
        {0, 9, NewBoard(9), false},
        {9, 0, NewBoard(9), false},
        {0, -1, NewBoard(9), false},
        {-1, 0, NewBoard(9), false},
        {0, 19, NewBoard(19), false},
        {11, 0, NewBoard(11), false},
    }
    for _, tc := range tt {
        if v := tc.board.correctIndex(tc.y, tc.x); v != tc.value {
            t.Errorf("index [%d, %d] of board size %d is %t, got: %t",
                    tc.y, tc.x, tc.board.size, tc.value, v)
        }

    }
}

