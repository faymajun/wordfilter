package wordfilter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrie_GetDirties(t *testing.T) {
	trie := NewTrie()
	var patterns = []string{"ah", "ahk", "oars", "soar"}
	for _, v := range patterns {
		trie.Add(v)
	}
	trie.BuildFailurePointer()
	assert := assert.New(t)
	var tests = []struct {
		input    string
		expected []string
	}{
		{"aht", []string{"ah"}},
		{"ahkr", []string{"ah", "ahk"}},
		{"a", []string{}},
		{"soars", []string{"soar", "oars"}},
		{"soarsoars", []string{"soar", "oars", "soar", "oars"}},
	}
	for _, test := range tests {
		assert.Equal(trie.GetDirties(test.input), test.expected)
	}
}
