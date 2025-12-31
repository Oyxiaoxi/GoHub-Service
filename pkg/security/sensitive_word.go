// Package security 敏感词过滤
package security

import (
	"strings"
	"sync"
	"unicode"
)

// SensitiveWordFilter 敏感词过滤器
type SensitiveWordFilter struct {
	mu          sync.RWMutex
	words       []string
	trie        *TrieNode
	replacement string
}

// TrieNode 前缀树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

var (
	once   sync.Once
	filter *SensitiveWordFilter
)

// NewSensitiveWordFilter 创建敏感词过滤器
func NewSensitiveWordFilter(replacement string) *SensitiveWordFilter {
	return &SensitiveWordFilter{
		words:       make([]string, 0),
		trie:        &TrieNode{children: make(map[rune]*TrieNode)},
		replacement: replacement,
	}
}

// GetFilter 获取全局敏感词过滤器实例
func GetFilter() *SensitiveWordFilter {
	once.Do(func() {
		filter = NewSensitiveWordFilter("***")
		// 添加默认敏感词
		filter.LoadDefaultWords()
	})
	return filter
}

// LoadDefaultWords 加载默认敏感词列表
func (f *SensitiveWordFilter) LoadDefaultWords() {
	// 政治敏感词
	politicalWords := []string{
		"法轮功", "法轮大法", "反党", "反共", "反政府",
		"台独", "藏独", "疆独", "港独", "分裂",
	}

	// 色情敏感词
	pornWords := []string{
		"色情", "淫秽", "裸聊", "约炮", "一夜情",
		"卖淫", "嫖娼", "性交易", "黄色网站",
	}

	// 暴力敏感词
	violenceWords := []string{
		"杀人", "自杀", "爆炸", "恐怖袭击", "制造炸弹",
		"持枪", "枪支买卖", "暴力革命",
	}

	// 赌博欺诈
	gamblingWords := []string{
		"网络赌博", "赌场", "六合彩", "时时彩",
		"诈骗", "传销", "非法集资", "高利贷",
		"刷单", "洗钱", "黑客", "盗号",
	}

	// 毒品相关
	drugWords := []string{
		"毒品", "大麻", "海洛因", "冰毒", "可卡因",
		"吸毒", "贩毒", "制毒",
	}

	// 合并所有敏感词
	allWords := append(politicalWords, pornWords...)
	allWords = append(allWords, violenceWords...)
	allWords = append(allWords, gamblingWords...)
	allWords = append(allWords, drugWords...)

	f.AddWords(allWords...)
}

// AddWords 添加敏感词
func (f *SensitiveWordFilter) AddWords(words ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, word := range words {
		if word == "" {
			continue
		}
		f.words = append(f.words, word)
		f.addToTrie(word)
	}
}

// addToTrie 添加词到前缀树
func (f *SensitiveWordFilter) addToTrie(word string) {
	node := f.trie
	runes := []rune(strings.ToLower(word))

	for _, r := range runes {
		if node.children[r] == nil {
			node.children[r] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[r]
	}
	node.isEnd = true
}

// Filter 过滤文本中的敏感词
func (f *SensitiveWordFilter) Filter(text string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if text == "" {
		return text
	}

	runes := []rune(text)
	result := make([]rune, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		maxMatchLength := 0
		node := f.trie

		// 尝试匹配最长的敏感词
		for j := i; j < len(runes); j++ {
			r := unicode.ToLower(runes[j])
			if node.children[r] == nil {
				break
			}
			node = node.children[r]
			if node.isEnd {
				maxMatchLength = j - i + 1
			}
		}

		if maxMatchLength > 0 {
			// 找到敏感词，替换
			result = append(result, []rune(f.replacement)...)
			i += maxMatchLength - 1
		} else {
			result = append(result, runes[i])
		}
	}

	return string(result)
}

// Contains 检查文本是否包含敏感词
func (f *SensitiveWordFilter) Contains(text string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if text == "" {
		return false
	}

	runes := []rune(strings.ToLower(text))

	for i := 0; i < len(runes); i++ {
		node := f.trie

		for j := i; j < len(runes); j++ {
			r := runes[j]
			if node.children[r] == nil {
				break
			}
			node = node.children[r]
			if node.isEnd {
				return true
			}
		}
	}

	return false
}

// FindAll 查找文本中所有的敏感词
func (f *SensitiveWordFilter) FindAll(text string) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	found := make([]string, 0)
	if text == "" {
		return found
	}

	runes := []rune(strings.ToLower(text))
	matched := make(map[string]bool)

	for i := 0; i < len(runes); i++ {
		node := f.trie

		for j := i; j < len(runes); j++ {
			r := runes[j]
			if node.children[r] == nil {
				break
			}
			node = node.children[r]
			if node.isEnd {
				word := string(runes[i : j+1])
				if !matched[word] {
					found = append(found, word)
					matched[word] = true
				}
			}
		}
	}

	return found
}

// RemoveWords 删除敏感词
func (f *SensitiveWordFilter) RemoveWords(words ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 重建前缀树
	f.trie = &TrieNode{children: make(map[rune]*TrieNode)}
	f.words = make([]string, 0)

	// 重新添加除了要删除的词之外的所有词
	removeMap := make(map[string]bool)
	for _, word := range words {
		removeMap[word] = true
	}

	for _, word := range f.words {
		if !removeMap[word] {
			f.words = append(f.words, word)
			f.addToTrie(word)
		}
	}
}

// Clear 清空所有敏感词
func (f *SensitiveWordFilter) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.words = make([]string, 0)
	f.trie = &TrieNode{children: make(map[rune]*TrieNode)}
}

// GetWords 获取所有敏感词
func (f *SensitiveWordFilter) GetWords() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	words := make([]string, len(f.words))
	copy(words, f.words)
	return words
}

// SetReplacement 设置替换字符
func (f *SensitiveWordFilter) SetReplacement(replacement string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.replacement = replacement
}
