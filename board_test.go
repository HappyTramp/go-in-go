package go_in_go

import (
    "testing"
    "reflect"
    "sort"
)


func TestNewBoard(t *testing.T) {
    tt := []struct{
        name string
        board *Board
        expected_size int
    }{
        {"9x9 board", NewBoard(9), 9},
        {"11x11 board", NewBoard(11), 11},
        {"19x19 board", NewBoard(19), 19},
    }
    for _, tc := range tt {
        if tc.board.size != tc.expected_size {
            t.Errorf("%s size != %d, is %d",
                     tc.name, tc.expected_size, tc.board.size)
        }
        for i := range tc.board.grid {
            for j := range tc.board.grid[i] {
                if tc.board.grid[i][j] != 0 {
                    t.Errorf("%s grid at [%d, %d] is %d, should be 0",
                             tc.name, i, j, tc.board.grid[i][j])
                }
            }
        }
    }
}


func TestGetPoint(t *testing.T) {
    board := NewBoard(9)
    board.grid[0][0] = White
    board.grid[1][1] = Black
    tt := []struct{
        y, x int
        value byte
        err error
    }{
        {0, 0, White, nil},
        {1, 1, Black, nil},
        {3, 3, Empty, nil},
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
        {0, 0, White, nil},
        {1, 1, Black, nil},
        {3, 3, White, nil},
        {3, 3, Empty, nil},
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
                t.Errorf("Point at [%d, %d] should be set to %d, is %d",
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
        neighbours Group
    }{
        {0, 0, Group{{0, 1, 0}, {1, 0, 0}}},
        {2, 2, Group{{2, 1, 0}, {2, 3, 2}, {1, 2, 0}, {3, 2, 1}}},
        {3, 3, Group{{3, 2, 1}, {3, 4, 0}, {2, 3, 2}, {4, 3, 0}}},
        {3, 4, Group{{3, 3, 0}, {3, 5, 0}, {2, 4, 0}, {4, 4, 2}}},
    }
    for _, tc := range tt {
        neighbours := board.Neighbours(tc.y, tc.x)
        if !compareGroup(neighbours, tc.neighbours) {
            t.Errorf("Wrong neighbours\ngot:    %v\nexpect: %v",
                     neighbours, tc.neighbours)
        }
    }
}


func TestUpdateGroup(t *testing.T) {
    board := NewBoard(9)
    tt_groups := []Group{
        Group{{0, 1, 1}, {1, 1, 1}, {1, 0, 1}, {2, 1, 1}, {1, 2, 1}},
        Group{{2, 2, 2}, {2, 3, 2}, {1, 3, 2}, {3, 3, 2}, {2, 4, 2}},
        Group{{8, 8, 1}, {8, 7, 1}, {7, 8, 1}, {7, 7, 1}},
        Group{{8, 6, 2}, {7, 6, 2}, {6, 6, 2}, {6, 7, 2}, {6, 8, 2}},
    }
    for _, tc_group := range tt_groups {
        for _, p := range tc_group {
            board.SetPoint(p.y, p.x, p.value)
        }
    }
    for _, tc_group := range tt_groups {
        group := Group{tc_group[0]}
        group = board.UpdateGroup(group)
        if !compareGroup(group, tc_group) {
            t.Errorf("Wrong group\ngot:    %v\nexpect: %v",
                     group, tc_group)
        }
    }
}


func TestGroupFrom(t *testing.T) {
    board := NewBoard(9)
    tt := []struct{
        y, x int
        group Group
    }{
        {0, 0, Group{{0, 0, 1}, {0, 1, 1}, {1, 1, 1}}},
        {5, 4, Group{{5, 4, 1}, {6, 4, 1}, {7, 4, 1}, {8, 4, 1}}},
        {3, 3, Group{{3, 3, 2}, {3, 2, 2}, {3, 1, 2}}},
        {8, 8, Group{{8, 8, 2}, {7, 8, 2}, {8, 7, 2}}},
    }
    for _, tc := range tt {
        for _, p := range tc.group {
            board.SetPoint(p.y, p.x, p.value)
        }
    }
    for _, tc := range tt {
        group := board.GroupFrom(tc.y, tc.x)
        if !compareGroup(group, tc.group) {
            t.Errorf("Wrong group\ngot:    %v\nexpect: %v",
                     group, tc.group)
        }
    }
}


func TestGroupLiberty(t *testing.T) {
    blackSurround := NewBoard(9)
    whiteSurround := NewBoard(9)
    tt := []struct{
        group Group
        liberty int
        board *Board
    }{
        {Group{{1, 1, 1}, {2, 1, 1}}, 0, blackSurround},
        {Group{{0, 0, 2}, {0, 1, 2}, {1, 0, 2}, {2, 0, 2}, {0, 2, 2},
               {3, 1, 2}, {1, 2, 2}, {2, 2, 2}, {3, 0, 2}, {3, 2, 2}},
         7, blackSurround},
        {Group{{8, 8, 2}, {7, 8, 2}, {6, 8, 2}}, 0, whiteSurround},
        {Group{{8, 7, 1}, {7, 7, 1}, {6, 7, 1}, {5, 7, 1}, {5, 8, 1}},
         6, whiteSurround},
    }
    for _, tc := range tt {
        for _, p := range tc.group {
            tc.board.SetPoint(p.y, p.x, p.value)
        }
    }
    for _, tc := range tt {
        liberty := tc.board.GroupLiberty(tc.group)
        if liberty != tc.liberty {
            t.Errorf("Liberty of %v should be %d, is %d",
                     tc.group, tc.liberty, liberty)
        }
    }
}


func TestStoneLiberty(t *testing.T) {
    board := NewBoard(9)
    board.SetPoint(0, 0, White)
    board.SetPoint(1, 1, White)
    board.SetPoint(2, 0, White)
    tt := []struct{
        y, x int
        liberty int
    }{
        {0, 0, 2},
        {1, 0, 0},
        {3, 0, 2},
        {5, 5, 4},
        {1, 2, 3},
        {4, 8, 3},
    }

    for _, tc := range tt {
        liberty := board.StoneLiberty(tc.y, tc.x)
        if liberty != tc.liberty {
            t.Errorf("Expected %d liberty at [%d, %d], got: %d",
                     tc.liberty, tc.y, tc.x, liberty)
        }
    }
}


func TestDeleteGroup(t *testing.T) {
    board := NewBoard(9)
    blackGroup := Group{{0, 0, 1}, {0, 1, 1}, {1, 1, 1}}
    whiteGroup := Group{{1, 0, 2}, {2, 0, 2}, {2, 1, 2}}
    for i := 0; i < len(blackGroup); i++ {
        s := blackGroup[i]
        board.SetPoint(s.y, s.x, s.value)
        s = whiteGroup[i]
        board.SetPoint(s.y, s.x, s.value)
    }
    board.DeleteGroup(blackGroup)
    for _, p := range blackGroup {
        if v, _ := board.GetPoint(p.y, p.x); v != Empty {
            t.Errorf("Member [%d, %d] of group %v should be deleted",
                     p.y, p.x, blackGroup)
        }
    }
    board.DeleteGroup(whiteGroup)
    for _, p := range whiteGroup {
        if v, _ := board.GetPoint(p.y, p.x); v != Empty {
            t.Errorf("Member [%d, %d] of group %v should be deleted",
                     p.y, p.x, whiteGroup)
        }
    }
}


func TestUpdateScore(t *testing.T) {
    board := NewBoard(9)
    board.UpdateScore(1, 100)
    board.UpdateScore(2, 150)
    if board.player1Score != 100 {
        t.Errorf("Player 1 score should be 100, is %d",
                 board.player1Score)
    }
    if board.player2Score != 150 {
        t.Errorf("Player 2 score should be 150, is %d",
                 board.player2Score)
    }
}


func TestKoDetection(t *testing.T) {
    board := NewBoard(9)
    blackGroup := Group{{0, 1, 1}, {1, 0, 1}, {2, 1, 1}, {1, 2, 0}}
    whiteGroup := Group{{0, 2, 2}, {2, 2, 2}, {1, 3, 2}}
    koMove := Point{1, 1, White}
    for _, v := range blackGroup {
        board.SetPoint(v.y, v.x, v.value)
    }
    for _, v := range whiteGroup {
        board.SetPoint(v.y, v.x, v.value)
    }
    previousBoard := board.clone()
    previousBoard.SetPoint(1, 2, Empty)
    previousBoard.SetPoint(1, 1, White)
    if board.KoDetection(previousBoard, koMove) != true {
        t.Errorf("ko move %v hasn't been detected", koMove)
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

func TestString(t *testing.T) {
    board9 := NewBoard(9)
    board9.SetPoint(0, 0, 1)
    board9.SetPoint(5, 4, 1)
    board9.SetPoint(8, 8, 2)
    board11 := NewBoard(11)
    board11.SetPoint(4, 0, 1)
    board11.SetPoint(5, 0, 1)
    board11.SetPoint(6, 0, 2)
    board19 := NewBoard(19)
    board19.SetPoint(18, 17, 2)
    board19.SetPoint(18, 18, 1)
    board19.SetPoint(17, 18, 2)
    tt := []struct{
        board *Board
        stringRepr string
    }{
        {board9, `
   A B C D E F G H J
1  X . . . . . . . .
2  . . . . . . . . .
3  . . . . . . . . .
4  . . . . . . . . .
5  . . . . . . . . .
6  . . . . X . . . .
7  . . . . . . . . .
8  . . . . . . . . .
9  . . . . . . . . O`},
        {board11, `
   A B C D E F G H J K L
1  . . . . . . . . . . .
2  . . . . . . . . . . .
3  . . . . . . . . . . .
4  . . . . . . . . . . .
5  X . . . . . . . . . .
6  X . . . . . . . . . .
7  O . . . . . . . . . .
8  . . . . . . . . . . .
9  . . . . . . . . . . .
10 . . . . . . . . . . .
11 . . . . . . . . . . .`},
        {board19, `
   A B C D E F G H J K L M N O P Q R S T
1  . . . . . . . . . . . . . . . . . . .
2  . . . . . . . . . . . . . . . . . . .
3  . . . . . . . . . . . . . . . . . . .
4  . . . . . . . . . . . . . . . . . . .
5  . . . . . . . . . . . . . . . . . . .
6  . . . . . . . . . . . . . . . . . . .
7  . . . . . . . . . . . . . . . . . . .
8  . . . . . . . . . . . . . . . . . . .
9  . . . . . . . . . . . . . . . . . . .
10 . . . . . . . . . . . . . . . . . . .
11 . . . . . . . . . . . . . . . . . . .
12 . . . . . . . . . . . . . . . . . . .
13 . . . . . . . . . . . . . . . . . . .
14 . . . . . . . . . . . . . . . . . . .
15 . . . . . . . . . . . . . . . . . . .
16 . . . . . . . . . . . . . . . . . . .
17 . . . . . . . . . . . . . . . . . . .
18 . . . . . . . . . . . . . . . . . . O
19 . . . . . . . . . . . . . . . . . O X`},
    }
    for _, tc := range tt {
        stringRepr := tc.board.String()
        if stringRepr != tc.stringRepr {
            t.Errorf("Wrong string representation\ngot: %v\nexpected: %v",
                    stringRepr, tc.stringRepr)
        }
    }
}

func compareGroup(s1, s2 Group) bool {
    if len(s1) != len(s2) {
        return false
    }
    sort.Slice(s1, func(i, j int) bool {
        return s1[i].y < s1[j].y
    })
    sort.Slice(s2, func(i, j int) bool {
        return s2[i].y < s2[j].y
    })
    for i := range s1 {
        if s1[i] != s2[i] {
            return false
        }
    }
    return true
}
