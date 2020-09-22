package wordfilter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter_Replace(t *testing.T) {
	f := NewFilter("")
	var patterns = []string{"ah", "ahk", "oars", "soar"}
	for _, v := range patterns {
		f.trie.Add(v)
	}
	f.trie.BuildFailurePointer()
	assert := assert.New(t)
	var tests = []struct {
		input    string
		expected string
	}{
		{"aht", "**t"},
		{"ahkr", "***r"},
		{"a", "a"},
		{"soars", "*****"},
		{"soarsoars", "*********"},
	}
	for _, test := range tests {
		assert.Equal(f.trie.Replace(test.input, '*'), test.expected)
	}
}

func TestFilter_Replace2(t *testing.T) {
	f := NewFilter("[0-9]") // 屏蔽数字
	var patterns = []string{"ah", "ahk", "oars", "soar"}
	for _, v := range patterns {
		f.trie.Add(v)
	}
	f.trie.BuildFailurePointer()
	assert := assert.New(t)
	var tests = []struct {
		input    string
		expected string
	}{
		{"ah45t", "**t"},
		{"ah2kr", "***r"},
		{"5a6", "a"},
		{"s6oa234rs", "*****"},
		{"so234ar234so234ars", "*********"},
	}
	for _, test := range tests {
		assert.Equal(f.trie.Replace(f.RemoveNoise(test.input), '*'), test.expected)
	}
}

func TestFilter_RemoveNoise(t *testing.T) {
	f := NewFilter("[0-9]") // 屏蔽数字
	assert := assert.New(t)
	var tests = []struct {
		input    string
		expected string
	}{
		{"aa23", "aa"},
		{"a23g", "ag"},
		{"222", ""},
		{"2t26d", "td"},
	}
	for _, test := range tests {
		assert.Equal(f.RemoveNoise(test.input), test.expected)
	}
}

func TestFilter_LoadAndReplace(t *testing.T) {
	f := NewFilter("[`~!@#$%^&*()+=|{}':;',\\[\\].<>/?~！@#￥%……&*（）——+|{}【】‘；：”“’。，、？]")
	err := f.LoadLocalWordFile("profanityworddict/list.txt")
	if err != nil {
		panic(err)
	}
	f.trie.BuildFailurePointer()
	assert := assert.New(t)
	var tests = []struct {
		input    string
		expected string
	}{
		{"hellboy", "****boy"},
		{"hello bitch world!", "****o ***** world"},
		{"hello bi@#tch world!", "****o ***** world"},
	}
	for _, test := range tests {
		assert.Equal(f.trie.Replace(f.RemoveNoise(test.input), '*'), test.expected)
	}
}
