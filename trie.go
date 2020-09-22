package wordfilter

// AcNode
type AcNode struct {
	character    rune
	isEndingChar bool             // 结尾字符为true
	length       int              // 当 isEndingChar = true时，记录从root到node长度
	fail         *AcNode          // 失败指针
	children     map[rune]*AcNode // 子节点map
}

// Trie
type Trie struct {
	root *AcNode
}

// NewTrie 创建一个Trie
func NewTrie() *Trie {
	return &Trie{
		root: NewNode(0)}
}

// NewNode // 创建一个Node
func NewNode(ch rune) *AcNode {
	return &AcNode{
		character: ch,
		children:  make(map[rune]*AcNode, 0),
	}
}

// Add 添加Profanity word
func (T *Trie) Add(word string) {
	var p = T.root
	runes := []rune(word)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if next, ok := p.children[r]; ok {
			p = next
		} else {
			newNode := NewNode(r)
			p.children[r] = newNode
			p = newNode
		}
	}
	p.isEndingChar = true
	p.length = len(runes)
}

// BuildFailurePointer 构建失败指针
// 逐层依次来求解每个节点的失败指针
func (T *Trie) BuildFailurePointer() {
	//logger.Infof("BuildFailurePointer....")
	var queue = make([]*AcNode, 0)
	var root = T.root

	queue = append(queue, root)
	for len(queue) > 0 {
		// get head
		var p = queue[0]
		queue = queue[1:]

		for _, pc := range p.children {
			if pc == nil {
				continue
			}

			if p == root {
				pc.fail = root
			} else {
				var q = p.fail
				for q != nil {
					// 如果找到相等
					if qc, ok := q.children[pc.character]; ok {
						pc.fail = qc
						break
					}
					// 没有找到，继续上面查找
					q = q.fail
				}
				// 最后没有找到，指向root
				if q == nil {
					pc.fail = root
				}
			}
			queue = append(queue, pc)
		}
	}
}

// Partial 返回Profanity word的在text中的位置
// 左闭，右开
type Partial struct {
	start int
	end   int
}

// mergePartial 合并Partial，
//func mergePartial(src []Partial) []Partial {
//	if len(src) == 0 {
//		return nil
//	}
//
//	res := make([]Partial, len(src))
//	res[0] = src[0]
//	ind := 0
//
//	for i := 1; i < len(src); i++ {
//		if src[i].start < res[ind].end {
//			res[ind].end = src[i].end
//		} else {
//			res = append(res, src[i])
//			ind++
//		}
//	}
//	return res
//}

// GetDirties 返回text包含的Profanity word
func (T *Trie) GetDirties(text string) []string {
	partials := T.match(text)
	var res = make([]string, len(partials))

	for i := 0; i < len(partials); i++ {
		res[i] = text[partials[i].start:partials[i].end]
	}
	return res
}

// Replace 使用rep替换text中包含的Profanity word
func (T *Trie) Replace(text string, rep rune) string {
	partials := T.match(text)
	if len(partials) == 0 {
		return text
	}

	var (
		runes = []rune(text)
		ind   = 0
	)

	for i := 0; i < len(runes); i++ {
		for ind < len(partials)-1 && partials[ind].end <= i {
			ind++
		}

		if i >= partials[ind].start && i < partials[ind].end {
			runes[i] = rep
		}
	}
	return string(runes)
}

// match 匹配所有包含的Profanity word
func (T *Trie) match(text string) []Partial {
	var (
		res   = make([]Partial, 0)
		runes = []rune(text)
		n     = len(runes)
		root  = T.root
		p     = root
	)

	for i := 0; i < n; i++ {
		pc, found := p.children[runes[i]]
		for !found && p != root {
			p = p.fail // 没有找到，通过失败指针往下继续匹配
			pc, found = p.children[runes[i]]
		}
		p = pc
		// 没有匹配的，从root开始重新匹配
		if p == nil {
			p = root
		}
		var tmp = p
		// 返回可以匹配的Profanity word
		for tmp != root {
			if tmp.isEndingChar {
				var startPos = i - tmp.length + 1
				res = append(res, Partial{startPos, i + 1})
			}
			tmp = tmp.fail
		}
	}
	return res
}
