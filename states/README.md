# states
Package states implements state list, stack, queue based on a uint64, built for reduce allocations.

# Stack
```go
import "github.com/cosiner/gohper/states"
const (
    STATE1 uint = iota + 1
    STATE2
    STATE3
    STATE4
)

func main() {
stack := states.NewStack(states.UnitSize(STATE4))

stack.
    Push(STATE1).
    Push(STATE2).
    Push(STATE3).
    Push(STATE4).
    IsFull() // false

stack.Pop() // STATE4
stack.Pop() // STATE3
stack.Pop() // STATE2
stack.Pop() // STATE1
stack.IsEmpty() // true
}
```

# Others:
* List
    * `PushBack`
    * `PushFront`
    * `PopBack`
    * `PopFront`
* Queue
    * `Push`
    * `Pop`
