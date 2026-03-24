package math

import "testing"

func TestAdd(t *testing.T) {
    if got := Add(1, 1); got != 2 {
        t.Errorf("1 + 1 = %d; want 2", got)
    }
}