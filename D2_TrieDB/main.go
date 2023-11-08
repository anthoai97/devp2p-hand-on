package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
)

func updateString(trie *trie.Trie, k, v string) {
	trie.MustUpdate([]byte(k), []byte(v))
}

func deleteString(trie *trie.Trie, k string) {
	trie.MustDelete([]byte(k))
}

// TrieDB usage
func main() {
	diskdb := rawdb.NewMemoryDatabase()
	trie := trie.NewEmpty(trie.NewDatabase(diskdb, nil))

	// trie.MustUpdate([]byte("120000"), []byte("qwerqwerqwerqwerqwerqwerqwerqwer"))
	root := trie.Hash()
	fmt.Println(root)
	fmt.Println("Update trie DB")

	updateString(trie, "doe", "reindeer")
	trie.Hash()
	// updateString(trie, "dog", "puppy")
	// updateString(trie, "dogglesworth", "cat")

	// root, _, _ = trie.Commit(false)
	// fmt.Println(root)

}
