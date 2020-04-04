package commands

import "github.com/anaskhan96/soup"

const (
	ReadFailed string = "Failed to retrieve documentation.\nCheck your spelling?"
)

// Workaround functions, since some of them won't cooperate because of their inconsistency.
var JavaDocWorkarounds = make(map[string]func(root soup.Root, inheritanceTree *string) string)

func init() {
	JavaDocWorkarounds["Interface Deque<E>"] = func(root soup.Root, inheritanceTree *string) string {
		*inheritanceTree = "java.lang.Object\njava.lang.Iterable\njava.util.Collection\njava.util.Queue\njava.util.Deque"
		return root.Text() +
			"deque is short for \"double ended queue\"\n and is usually pronounced \"deck\".  Most Deque " +
			"implementations place no fixed limits on the number of elements\n they may contain," +
			" but this interface supports capacity-restricted\n deques as well as those with no fixed size limit."
	}
}
