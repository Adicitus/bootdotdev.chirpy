package trie

import (
	"io"
	"strings"
	"unicode"
)

type TrieNode struct {
	End             bool
	Next            map[rune]*TrieNode
	CaseInsensitive bool
}

func NewTrie() *TrieNode {
	n := new(TrieNode)
	n.Next = make(map[rune]*TrieNode)
	return n
}

func (n *TrieNode) Add(s string) error {
	return n.add(*strings.NewReader(s))
}

func (n *TrieNode) Contains(s string) (present bool) {
	present, _ = n.present(strings.NewReader(s), false, 0)
	return
}

func (n *TrieNode) Prefix(s string) (bool, int) {
	return n.PresentAt(s, 0)
}

func (n *TrieNode) PresentAt(s string, offset int) (bool, int) {
	r := strings.NewReader(s)
	r.Seek(int64(offset), io.SeekStart)
	return n.present(r, true, 0)
}

func (n *TrieNode) LPresentAt(s string, offset int) (bool, int) {
	r := strings.NewReader(s)
	r.Seek(int64(offset), io.SeekStart)
	return n.lPresent(r, 0, 0)
}

func (n *TrieNode) Replace(s string, replacer string) (string, error) {
	old_s := strings.NewReader(s)
	new_s := new(strings.Builder)

	lim := int(old_s.Size())

	for old_s.Len() > 0 {
		i := lim - old_s.Len()
		p, l := n.LPresentAt(s, i)

		if p {
			old_s.Seek(int64(l), io.SeekCurrent)
			new_s.WriteString(replacer)
			continue
		}

		r, _, err := old_s.ReadRune()

		if err != nil {
			return "", err
		}

		new_s.WriteRune(r)
	}

	return new_s.String(), nil
}

func (n *TrieNode) Complete(sample string) ([]string, error) {
	c, err := n.complete(strings.NewReader(sample), new(strings.Builder), []string{})

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (n *TrieNode) add(tail strings.Reader) error {
	if tail.Len() == 0 {
		n.End = true
		return nil
	}

	r, _, err := tail.ReadRune()

	if err != nil {
		return nil
	}

	if n.CaseInsensitive {
		r = unicode.ToLower(r)
	}

	next, ok := n.Next[r]

	if !ok {
		next = NewTrie()
		n.Next[r] = next
	}

	return next.add(tail)
}

/*
Checks if any of the words in Trie occur in the strings.Reader at its current offset.

This method is greedy and will look for the shortest possible match.

So given a the search string "abcd" Trie containing the strings "ab" and "abc", this
method will match "ab" and return (true, 2)
*/
func (n *TrieNode) present(tail *strings.Reader, partial bool, agg int) (found bool, length int) {
	if n.End && (tail.Len() == 0 || partial) {
		return true, agg
	}

	r, _, err := tail.ReadRune()

	if n.CaseInsensitive {
		r = unicode.ToLower(r)
	}

	if err != nil {
		return false, 0
	}

	next, ok := n.Next[r]

	if !ok {
		return false, 0
	}

	return next.present(tail, partial, agg+1)
}

/*
Checks if any of the words in Trie occur in the strings.Reader at its current offset.

This method will try to find the longest possible word matching word.

So given a the search string "abcd" Trie containing the strings "ab" and "abc", this
method will match "abc" and return (true, 3)
*/
func (n *TrieNode) lPresent(tail *strings.Reader, agg int, longest int) (found bool, length int) {

	if n.End {
		// We've found a new longest word:
		longest = agg
	}

	r, _, err := tail.ReadRune()

	if err != nil {
		if err == io.EOF {
			// End of input, there won't be any more matches so return the best match:
			return longest > 0, longest
		}

		return false, 0
	}

	if n.CaseInsensitive {
		r = unicode.ToLower(r)
	}

	next, ok := n.Next[r]

	if !ok {
		// We've reached the end of all possible candidate words, return the best match:
		return longest > 0, longest
	}

	return next.lPresent(tail, agg+1, longest)
}

func (n *TrieNode) complete(sample *strings.Reader, agg *strings.Builder, completions []string) ([]string, error) {
	if n.End {
		completions = append(completions, agg.String())
	}

	if sample.Len() > 0 {

		r, _, err := sample.ReadRune()

		if err != nil {
			return nil, err
		}

		agg.WriteRune(r)

		if n.CaseInsensitive {
			r = unicode.ToLower(r)
		}

		next, ok := n.Next[r]

		if !ok {
			return completions, nil
		}

		return next.complete(sample, agg, completions)
	}

	for r, next := range n.Next {

		agg.WriteRune(r)

		if n.CaseInsensitive {
			r = unicode.ToLower(r)
		}

		new_agg := new(strings.Builder)
		new_agg.WriteString(agg.String())

		var err error

		completions, err = next.complete(sample, new_agg, completions)

		if err != nil {
			return nil, err
		}
	}

	return completions, nil
}
