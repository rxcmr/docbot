package commands

const (
	ReadFailed string = "Failed to retrieve documentation.\nCheck your spelling?"
)

// Workaround "constants"
var JavaDocWorkarounds = make(map[string]string)

func init() {
	JavaDocWorkarounds["Interface Deque<E>"] = `A linear collection that supports element insertion and removal at both ends. 
The name deque is short for "double ended queue" and is usually pronounced "deck". 
Most Deque implementations place no fixed limits on the number of elements they may contain, but this interface supports capacity-restricted deques as well as those with no fixed size limit.`
}
