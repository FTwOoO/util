package word_filter

import (
	"sync"
)

type Trie struct {
	Root           *Node
	Mutex          sync.RWMutex
	CheckWhiteList bool // 是否检查白名单
}

type Node struct {
	SubNodes map[rune]*Node
	End      bool
}

func NewTrie() *Trie {
	t := new(Trie)
	t.Root = NewTrieNode()
	return t
}

func NewTrieNode() *Node {
	n := new(Node)
	n.SubNodes = make(map[rune]*Node)
	n.End = false
	return n
}

func (this *Trie) Import(keywords []string) {
	for _, d := range keywords {
		this.Add(d)
	}
}

// Add 添加一个敏感词(UTF-8的)到Trie树中
func (this *Trie) Add(keyword string) {
	chars := []rune(keyword)

	if len(chars) == 0 {
		return
	}

	this.Mutex.Lock()

	node := this.Root
	for _, unicodeChar := range chars {
		if _, ok := node.SubNodes[unicodeChar]; !ok {
			node.SubNodes[unicodeChar] = NewTrieNode()
		}
		node = node.SubNodes[unicodeChar]
	}
	node.End = true

	this.Mutex.Unlock()
}

// Del 从Trie树中删除一个敏感词
func (this *Trie) Del(keyword string) {
	chars := []rune(keyword)
	if len(chars) == 0 {
		return
	}

	this.Mutex.Lock()
	node := this.Root
	this.cycleDel(node, chars, 0)
	this.Mutex.Unlock()
}

func (this *Trie) cycleDel(node *Node, chars []rune, index int) (shouldDel bool) {
	char := chars[index]
	l := len(chars)

	if n, ok := node.SubNodes[char]; ok {
		if index+1 < l {
			shouldDel = this.cycleDel(n, chars, index+1)
			if shouldDel && len(n.SubNodes) == 0 {
				if n.End { // 说明这是一个敏感词，不能删除
					shouldDel = false
				} else {
					delete(node.SubNodes, char)
				}
			}
		} else if n.End {
			if len(n.SubNodes) == 0 { // 是最后一个节点
				shouldDel = true
				delete(node.SubNodes, char)

			} else { // 不是最后一个节点
				n.End = false
			}
		}
	}

	return
}

// Query 查询敏感词
// 将text中在trie里的敏感字，替换为*号
// 返回结果: 是否有敏感字, 敏感字数组, 替换后的文本
func (this *Trie) Query(text string) (bool, []string, string) {
	return this.QueryAndReplace(text, 42)
}

func (this *Trie) QueryAndReplace(text string, replaceWithChar rune) (bool, []string, string) {
	found := []string{}
	chars := []rune(text)
	l := len(chars)
	if l == 0 {
		return false, found, text
	}

	var (
		i, j, jj int
		ok       bool
	)

	node := this.Root
	for i = 0; i < l; i++ {
		if _, ok = node.SubNodes[chars[i]]; !ok {
			continue
		}

		jj = 0

		node = node.SubNodes[chars[i]]
		for j = i + 1; j < l; j++ {
			if _, ok = node.SubNodes[chars[j]]; !ok {
				if jj > 0 {
					tmpFound := chars[i : jj+1]
					found = append(found, string(tmpFound))

					if replaceWithChar != 0 {
						this.replaceToAsterisk(chars, i, jj, replaceWithChar)
					}
					i = jj
				}
				break
			}

			node = node.SubNodes[chars[j]]
			if node.End {
				jj = j //还有子节点的情况, 记住上次找到的位置, 以匹配最大串 (eg: AV, AV女优)

				if len(node.SubNodes) == 0 || j+1 == l { // 是最后节点或者最后一个字符, break
					tmpFound := chars[i : j+1]
					found = append(found, string(tmpFound))
					if replaceWithChar != 0 {
						this.replaceToAsterisk(chars, i, j, replaceWithChar)
					}
					i = j
					break
				}
			}
		}
		node = this.Root
	}

	exist := false
	if len(found) > 0 {
		exist = true
	}

	return exist, found, string(chars)
}

// 替换为*号
func (this *Trie) replaceToAsterisk(chars []rune, i, j int, replaceRune rune) {
	for k := i; k <= j; k++ {
		chars[k] = replaceRune
	}
	return
}

// ReadAll 返回所有敏感词
func (this *Trie) ReadAll() (words []string) {
	this.Mutex.Lock()
	words = []string{}
	words = this.cycleRead(this.Root, words, "")
	this.Mutex.Unlock()
	return
}

func (this *Trie) cycleRead(node *Node, words []string, parentWord string) []string {
	for char, n := range node.SubNodes {
		if n.End {
			words = append(words, parentWord+string(char))
		}
		if len(n.SubNodes) > 0 {
			words = this.cycleRead(n, words, parentWord+string(char))
		}
	}
	return words
}
