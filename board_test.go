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
        if !comparePointSlice(neighbours, tc.neighbours) {
            t.Errorf("Wrong neighbours\ngot:    %v\nexpect: %v",
                     neighbours, tc.neighbours)
        }
    }
}

func TestUpdateGroup(t *testing.T) {
    crossBoard := NewBoard(9)
    surroundBoard := NewBoard(9)
    tt := []struct{
        name string
        group Group
        board *Board
    }{
        {"black cross board",
         Group{{0, 1, 1}, {1, 1, 1}, {1, 0, 1}, {2, 1, 1}, {1, 2, 1}},
         crossBoard},
        {"white cross board",
         Group{{2, 2, 2}, {2, 3, 2}, {1, 3, 2}, {3, 3, 2}, {2, 4, 2}},
         crossBoard},
        {"black surround board",
         Group{{0, 0, 1}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1}},
         surroundBoard},
        {"white surround board",
         Group{{2, 0, 2}, {2, 1, 2}, {2, 2, 2}, {1, 2, 2}, {0, 2, 2}},
         surroundBoard},
    }
    for _, tc := range tt {
        for _, p := range tc.group {
            tc.board.SetPoint(p.y, p.x, p.value)
        }
    }
    for _, tc := range tt {
        group := Group{tc.group[0]}
        group = tc.board.UpdateGroup(group)
        if !comparePointSlice(group, tc.group) {
            t.Errorf("Wrong group on %s\ngot:    %v\nexpect: %v",
                     tc.name, group, tc.group)
        }
    }
}

func TestGroupFrom(t *testing.T) {

}

//func TestGroupLiberty(t *testing.T) {
    //blackSurround := NewBoard(9)
    //whiteSurround := NewBoard(9)
    //tt := []struct{
        //group Group
        //liberty int
        //board *Board
    //}{
        //{Group{{1, 1, 1}, {2, 1, 1}}, 0, blackSurround},
        //{Group{{0, 1, 2}, {1, 0, 2}, {2, 0, 2},
               //{3, 1, 2}, {1, 3, 2}, {2, 3, 2}},
         //7, blackSurround},
        //{Group{{8, 8, 2}, {7, 8, 2}, {6, 8, 2}}, 0, whiteSurround},
        //{Group{{8, 7, 1}, {7, 7, 1}, {6, 7, 1}, {5, 7, 1}, {5, 8, 1}},
         //5, whiteSurround},
    //}
    //for _, tc := range tt {
        //liberty := tc.board.GroupLiberty(tc.group)
        //if liberty != tc.liberty {
            //t.Errorf("Liberty of %v should be %d, is %d",
                     //tc.group, tc.liberty, liberty)
        //}
    //}
//}

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

func comparePointSlice(s1, s2 Group) bool {
    if len(s1) != len(s2) {
        return false
    }
    sort.Slice(s1, func(i, j int) bool {
        return s1[i].y < s1[j].y
    })
    sort.Slice(s2, func(i, j int) bool {
        return s2[i].y < s2[j].y
    })
    for i, p := range s1 {
        if p.y != s2[i].y || p.x != s2[i].x || p.value != s2[i].value {
            return false
        }
    }
    return true
}
