package main

import (
	"fmt"
)

/*
  create a node struct
  each letter of word is a node
  at the end of word link to longest suffix in trie except itself

  ex:
   - acc
   - atc
   - cat
   - gcg

                   **********
		/     \        \
	        A      C        G
	       / \      \        \
	      C   T      A        C
	     /     \      \         \
	    C       C      C         G
*/
type Node struct {
	letter     string
	child      [26]*Node
	suffix     *Node
	output     *Node
	pattern_id int
}

// buildTrie build a trie from word received
func (t *Node) buildTrie(word []string) {
	/*
	   for each letter of dictionnary word, create a node if it does not exist
	   if so, add a child from the existing node
	*/
	for i := 0; i < len(word); i++ {
		wordLength := len(word[i])
		currentNode := t
		for j := 0; j < wordLength; j++ {
			charIndex := word[i][j] - 'a'
			if currentNode.child[charIndex] == nil {
				currentNode.child[charIndex] = &Node{letter: string(word[i][j]), suffix: nil, output: nil, pattern_id: -1}
			}
			currentNode = currentNode.child[charIndex]
		}
		/*
			update pattern_id at the end of word in trie with index of word
		*/
		currentNode.pattern_id = i
	}
}
func buildSuffix(root *Node) {
	root.suffix = root
	var queue []*Node
	/*
		link to root each direct child of root
	*/
	for i := 0; i < len(root.child); i++ {
		if root.child[i] != nil {
			root.child[i].suffix = root
			queue = append(queue, root.child[i])
		}
	}
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		for i := 0; i < len(currentNode.child); i++ {
			if currentNode.child[i] != nil {
				letter := currentNode.child[i].letter
				temp := currentNode.suffix
				index := []rune(letter)[0] - 'a'
				/*
					if char is not found in current node, try to start over from suffix
					until it's found out
				*/
				for temp.child[index] == nil && temp != root {
					temp = temp.suffix
				}
				/*
				  if root is reached, add suffix to root
				*/
				if temp.child[index] != nil {
					currentNode.child[i].suffix = temp.child[index]
				} else {
					currentNode.child[i].suffix = root
				}
				queue = append(queue, currentNode.child[i])
			}

		}
		/*
		 if word is found (pattern_id > 0), add output to current suffix
		*/
		if currentNode.suffix.pattern_id >= 0 {
			currentNode.output = currentNode.suffix
		} else {
			currentNode.output = currentNode.suffix.output
		}
	}
}
func (t *Node) search(word string) map[int][]int {
	parent := t
	indices := make(map[int][]int)
	for i := 0; i < len(word); i++ {
		charStr := word[i] - 'a'
		if parent.child[charStr] != nil {
			parent = parent.child[charStr]
			if parent.pattern_id >= 0 {
				indices[parent.pattern_id] = append(indices[parent.pattern_id], i)
			}
			temp := parent.output
			for temp != t {
				indices[temp.pattern_id] = append(indices[temp.pattern_id], i)
				temp = temp.output
			}
		} else {
			for parent != t && parent.child[charStr] == nil {
				parent = parent.suffix
			}
			if parent.child[charStr] != nil {
				i--
			}
		}
	}

	return indices
}
func (t *Node) display() {
	for _, v := range t.child {
		if v != nil {
			fmt.Println("------------------------------")
			fmt.Printf("node address %p \n", v)
			fmt.Println("letter", v.letter, "suffix", v.suffix.letter)
			fmt.Printf("suffix address %p \n", v.suffix)
			fmt.Println("pattern_id", v.pattern_id)
			fmt.Printf("output address %p \n", v.output)
			fmt.Println("------------------------------")
			v.display()
		}
	}
}
func main() {
	root := &Node{
		letter: "root",
	}

	data := []string{"a", "ab", "bab", "bc", "bca", "c", "caa"}
	fmt.Println("dictionnary", data)
	root.buildTrie(data)
	buildSuffix(root)
	lookingFor := "abccab"
	result := root.search(lookingFor)
	fmt.Println("result", result, "len", len(lookingFor))
	for k, _ := range data {
		if result[k] != nil {
			fmt.Println("data", data[k], "is found on position", result[k])
		}
	}
	root.display()
}
