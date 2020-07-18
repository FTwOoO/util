package word_filter

func NewDefaultTrie() *Trie {
	blackTrie := NewTrie()
	blackTrie.CheckWhiteList = true
	blackTrie.Import(BlacksNew)
	return blackTrie
}
