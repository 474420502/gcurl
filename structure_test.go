package gcurl

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPQueue(t *testing.T) {
	PQExec := newPQueueExecute()
	PQExec.Push(&parseFunction{Priority: 5})
	PQExec.Push(&parseFunction{Priority: 10})
	PQExec.Push(&parseFunction{Priority: 4})
	PQExec.Push(&parseFunction{Priority: 4})
	PQExec.Push(&parseFunction{Priority: 20})
	PQExec.Push(&parseFunction{Priority: 10})
	PQExec.Push(&parseFunction{Priority: 15})

	content := ""
	for PQExec.Len() > 0 {
		content += strconv.Itoa(PQExec.Pop().Priority)
		content += " "
	}
	if content != "4 4 5 10 10 15 20 " {
		t.Error(content)
	}
}

type Word string

func (w Word) GetWord() string {
	return string(w)
}

func TestTrie(t *testing.T) {
	// OptionTrie 设置的前缀树
	var trie *hTrie = newTrie()

	trie.Insert(Word("123"))
	trie.Insert(Word("12"))

	if fmt.Sprintf("%v", trie.AllWords()) != "[12 123]" {
		t.Error(trie.AllWords())
	}

	trie.Remove("12")

	if fmt.Sprintf("%v", trie.AllWords()) != "[123]" {
		t.Error(trie.AllWords())
	}

	trie = newTrie()
	for i := 0; i < 100; i++ {
		trie.Insert(Word(strconv.Itoa(i)))
	}

	m := trie.Match(Word("12"))
	if m == nil && m != "12" {
		t.Error("match error", m)
	}

	m = trie.Match(Word("100"))
	if m != nil {
		t.Error("match error", m)
	}

	for i := 0; i < 50; i++ {
		trie.Remove(strconv.Itoa(i))
	}

	if !trie.StartsWith("4") {
		t.Error("start error")
	}

}
