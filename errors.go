package go_in_go

import (
    "fmt"
)

type BoardIndexError struct {
    y, x int
}

func (e *BoardIndexError) Error() string {
    return fmt.Sprintf("BoardIndexError: [%d, %d] is not a valid index",
                       e.y, e.x)
}
