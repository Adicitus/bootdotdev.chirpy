package trie_test

import (
	"testing"

	. "github.com/Adicitus/bootdotdev.chirpy/trie"
)

func TestTrieNodeAdd(t *testing.T) {
	var ok bool
	var n, top *TrieNode

	top = NewTrie()

	top.Add("a")
	n, ok = top.Next['a']
	if !ok {
		t.Log("Failed to add top 'a' node")
		t.Fail()
	} else {
		if !n.End {
			t.Log("Top 'a' node not marked as end of word")
			t.Fail()
		}
	}

	top.Add("b")
	n, ok = top.Next['b']
	if !ok {
		t.Log("Failed to add top 'b' node")
		t.Fail()
	} else {
		if !n.End {
			t.Log("Top 'b' node notmarked as end of word")
			t.FailNow()
		}
	}

	top.Add("abc")
	n, ok = top.Next['a']
	if !ok {
		t.Log("Failed to get top 'b' node")
		t.FailNow()
	}

	n, ok = n.Next['b']
	if !ok {
		t.Log("Failed to get second 'b' node")
		t.Fail()
	} else {
		if n.End {
			t.Log("Second 'b' node marked as end of word")
			t.Fail()
		}
	}

	n, ok = n.Next['c']
	if !ok {
		t.Log("Failed to get 'c' node")
		t.Fail()
	} else {
		if !n.End {
			t.Log("'c' node not marked as end of word")
			t.Fail()
		}
	}
}

func TestTrieNodeContains(t *testing.T) {
	var top *TrieNode

	top = NewTrie()

	top.Add("ab")
	top.Add("abcd")

	if !top.Contains("ab") {
		t.Log("'ab' not found in 'ab'")
		t.Fail()
	}

	if top.Contains("abc") {
		t.Log("'ab' not found in 'abc'")
		t.Fail()
	}

	if !top.Contains("abcd") {
		t.Log("'abcd' not found in 'abcd'")
		t.Fail()
	}
}

func TestTrieNodePresentAt(t *testing.T) {
	var p bool
	var l int
	var top *TrieNode

	top = NewTrie()

	top.Add("a")

	p, l = top.PresentAt("abc", 0)

	if !p {
		t.Log("'a' not found at the start of 'abc'")
		t.Fail()
	}

	if l != 1 {
		t.Logf("Reported incorrect length (%d) when finding 'a' in 'abc'", l)
		t.Fail()
	}

	p, l = top.PresentAt("abc", 1)

	if p {
		t.Log("'a' found in 'bc'")
		t.Fail()
	}

	if l != 0 {
		t.Logf("reported length %d when no word detected", l)
		t.Fail()
	}

	top.Add("ab")

	p, l = top.PresentAt("abcd", 0)

	if !p {
		t.Log("neither 'a' nor 'ab' not found at start of 'abcd'")
		t.Fail()
	}

	if l != 1 {
		t.Logf("reported prefix length %d when 'a' should be found in 'abcd', expected 1", l)
		t.Fail()
	}

}

func TestTrieNodeLPresentAt(t *testing.T) {
	var p bool
	var l int
	var top *TrieNode

	top = NewTrie()

	top.Add("a")
	top.Add("b")
	top.Add("ab")

	p, l = top.LPresentAt("aba", 0)
	if !p {
		t.Log("Failed to find either 'a', 'b' or 'ab' in string 'aba'")
		t.Fail()
	}

	if l != 2 {
		t.Logf("Reported incorrect length when looking for 'a', 'b' or 'ab' in 'aba'. Found %d,  expected %d", l, 2)
		t.Fail()
	}

	p, l = top.LPresentAt("aba", 1)

	if !p {
		t.Log("Found neither 'a', 'b' nor 'ab' in substring 'ba'")
		t.Fail()
	}

	if l != 1 {
		t.Logf("Reported incorrect length when looking for 'a' or 'ab' in substring 'ba'. Found %d, expected %d.", l, 1)
		t.Fail()
	}
}

func TestTrieNodeReplaceEmpty(t *testing.T) {
	words := NewTrie()

	words.Add("a")

	v := struct {
		s1 string
		r  string
		s2 string
	}{s1: "", r: "b", s2: ""}

	if s, err := words.Replace(v.s1, v.r); err != nil {
		t.Fail()
		t.Logf("Error: %s", err)
	} else {
		if s != v.s2 {
			t.Fail()
			t.Logf("Expected '%s', found '%s'", v.s2, s)
		}
	}
}

func TestTrieNodeReplaceSimple(t *testing.T) {
	words := NewTrie()

	words.Add("a")

	v := struct {
		s1 string
		r  string
		s2 string
	}{s1: "a", r: "b", s2: "b"}

	if s, err := words.Replace(v.s1, v.r); err != nil {
		t.Fail()
		t.Logf("Error: %s", err)
	} else {
		if s != v.s2 {
			t.Fail()
			t.Logf("Expected '%s', found '%s'", v.s2, s)
		}
	}
}

func TestTrieNodeReplacePreferLongest(t *testing.T) {
	words := NewTrie()

	words.Add("a")
	words.Add("ab")

	v := struct {
		s1 string
		r  string
		s2 string
	}{s1: "aba", r: "c", s2: "cc"}

	if s, err := words.Replace(v.s1, v.r); err != nil {
		t.Fail()
		t.Logf("Error: %s", err)
	} else {
		if s != v.s2 {
			t.Fail()
			t.Logf("Expected '%s', found '%s'", v.s2, s)
		}
	}
}

func TestTrieNodeReplaceCaseInsensitive(t *testing.T) {
	words := NewTrie()
	words.CaseInsensitive = true

	words.Add("Fornax")
	words.Add("kerfuffle")

	v := struct {
		s1 string
		r  string
		s2 string
	}{s1: "I really need a kerfuffle to go to bed sooner, Fornax !", r: "****", s2: "I really need a **** to go to bed sooner, **** !"}

	if s, err := words.Replace(v.s1, v.r); err != nil {
		t.Fail()
		t.Logf("Error: %s", err)
	} else {
		if s != v.s2 {
			t.Fail()
			t.Logf("Expected '%s', found '%s'", v.s2, s)
		}
	}
}

func TestTrieNodeCompletions(t *testing.T) {
	words := NewTrie()

	words.Add("fox")
	words.Add("foxbat")
	words.Add("foxhound")

	c, err := words.Complete("fox")

	if err != nil {
		t.Fail()
		t.Logf("Error: %v\n", err)
		return
	}

	if len(c) != 3 {
		t.Fail()
		t.Logf("Expected %d completions, found %d", 3, len(c))
	}

	c, err = words.Complete("Fox")

	if err != nil {
		t.Fail()
		t.Logf("Error: %v\n", err)
		return
	}

	if len(c) != 0 {
		t.Fail()
		t.Logf("Expected %d completions, found %d", 0, len(c))
	}
}

func TestTrieNodeCompletionsCaseInsensitive(t *testing.T) {
	words := NewTrie()
	words.CaseInsensitive = true

	words.Add("fox")
	words.Add("foxbat")
	words.Add("foxhound")

	c, err := words.Complete("Fox")

	if err != nil {
		t.Fail()
		t.Logf("Error: %v\n", err)
		return
	}

	if len(c) != 3 {
		t.Fail()
		t.Logf("Expected %d completions, found %d", 3, len(c))
	}
}
