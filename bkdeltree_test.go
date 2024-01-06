package bkdeltree

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// GenerateRandomString generates a random string of length between 6 and 8 characters.
func GenerateRandomString(seed int) string {
	var r *rand.Rand
	if seed != 0 {
		r = rand.New(rand.NewSource(int64(seed)))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Define the character set from which to generate the random string
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Generate a random string with a length between 6 and 8 characters
	length := r.Intn(3) + 6
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}

func TestDistance(t *testing.T) {
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test1", args: args{x: "abc", y: "abc"}, want: 0},
		{name: "test2", args: args{x: "abc", y: "abcd"}, want: 1},
		{name: "test3", args: args{x: "abc", y: "ab"}, want: 1},
		{name: "test4", args: args{x: "abc", y: "abcef"}, want: 2},
		{name: "test5", args: args{x: "abc", y: "abd"}, want: 1},
		{name: "test6", args: args{x: "abc", y: "abdef"}, want: 3},
		{name: "test7", args: args{x: "abc", y: "xyzijk"}, want: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distance(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBKDelTree(t *testing.T) {
	bktree := NewBKDelTree()
	if bktree == nil {
		t.Fatalf("NewBKDelTree() = %v, want %v", bktree, "not nil")
	}
	if bktree.root != nil {
		t.Errorf("NewBKDelTree(root) = %v, want %v", bktree.root, "nil")
	}
	if bktree.numNodes != 0 {
		t.Errorf("NewBKDelTree(numNodes) = %v, want %v", bktree.numNodes, 0)
	}
}

func TestBKDelTreeInsert(t *testing.T) {
	t.Run("StandardBKTree", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abc", "abcd", "ab", "abcef", "abd", "abdef", "xyzijk"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}
		if bktree.numNodes != len(words) {
			t1.Errorf("BKDelTree.Insert(numNodes) = %v, want %v", bktree.numNodes, len(words))
		}
		if bktree.root.word != "abc" {
			t1.Errorf("BKDelTree.Insert(root.word) = %v, want %v", bktree.root.word, "abc")
		}
	})
	t.Run("SameString", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abc", "abc", "abc", "abc", "abc", "abc", "abc"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}
		if bktree.numNodes != 1 {
			t1.Errorf("BKDelTree.Insert(numNodes) = %v, want %v", bktree.numNodes, 1)
		}
		if bktree.root.word != "abc" {
			t1.Errorf("BKDelTree.Insert(root.word) = %v, want %v", bktree.root.word, "abc")
		}
	})

	t.Run("SomeDuplicateString", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abc", "abcd", "abce", "xyz", "abc", "abce", "abced"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}
		if bktree.numNodes != 5 {
			t1.Errorf("BKDelTree.Insert(numNodes) = %v, want %v", bktree.numNodes, 5)
		}
	})
}

func TestBKDelTreeSearch(t *testing.T) {
	nilbktree := NewBKDelTree()
	if nilbktree.Search("abc", 0) != nil {
		t.Errorf("BKDelTree.Search() = %v, want %v", nilbktree.Search("abc", 0), nil)
	}

	bktree := NewBKDelTree()
	words := []string{"abc", "abcd", "ab", "abcef", "abd", "abdef", "xyzijk"}
	for _, word := range words {
		_ = bktree.Insert(word)
	}

	type args struct {
		word    string
		maxDist int
	}
	type gots map[string]bool

	tests := []struct {
		name string
		args args
		want gots
	}{
		{
			name: "test1",
			args: args{word: "abc", maxDist: 0},
			want: gots{"abc": true},
		},
		{
			name: "test2",
			args: args{word: "abc", maxDist: 1},
			want: gots{"abc": true, "abcd": true, "ab": true, "abd": true},
		},
		{
			name: "test3",
			args: args{word: "abc", maxDist: 2},
			want: gots{"abc": true, "abcd": true, "ab": true, "abd": true, "abcef": true, "abdef": true},
		},
		{
			name: "test4",
			args: args{word: "xyzijk", maxDist: 1},
			want: gots{"xyzijk": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			if got := bktree.Search(tt.args.word, tt.args.maxDist); got != nil {
				for _, g := range got {
					if _, ok := tt.want[g.word]; !ok {
						t1.Errorf("Search() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}

func TestBKDelTreeGetParent(t *testing.T) {
	t.Run("NilBKTree", func(t1 *testing.T) {
		nilbktree := NewBKDelTree()
		if parent, _ := nilbktree.GetParent("abc"); parent != nil {
			t1.Errorf("BKDelTree.Search() = %v, want %v", parent, nil)
		}
	})

	t.Run("NonExistentWord", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abc", "abcd", "ab", "abcef", "abd", "abdef", "xyzijk"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}
		if parent, _ := bktree.GetParent("abcde"); parent != nil {
			t1.Errorf("BKDelTree.Search() = %v, want %v", parent, nil)
		}
	})

	t.Run("StandardBKTree", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abcde", "abcdd", "abcdf", "acbed", "acebd", "cabde", "dadbc", "cdacb", "cadbc", "bacde"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}

		wantChildParents := map[string]string{
			"abcdd": "abcde",
			"cabde": "abcde",
			"acbed": "abcde",
			"acebd": "abcde",
			"cdacb": "abcde",
			"abcdf": "abcdd",
			"bacde": "cabde",
			"dadbc": "acebd",
			"cadbc": "dadbc",
		}
		for _, word := range words {
			if parent, _ := bktree.GetParent(word); parent != nil {
				if word == "abcde" { // special case for root node
					if parent != nil {
						t1.Errorf("GetParent(%v) = %v, want %v", word, parent.word, nil)
					}
				} else {
					if parent.word != wantChildParents[word] {
						t1.Errorf("GetParent(%v) = %v, want %v", word, parent.word, wantChildParents[word])
					}
				}
			}
		}
	})
}

func TestBKNodeCollectFamily(t *testing.T) {
	t.Run("NilBKNode", func(t1 *testing.T) {
		nilnode := &BKNode{}
		if len(nilnode.collectFamily()) != 1 {
			t1.Errorf("BKNode.collectFamily() = %v, want %v", len(nilnode.collectFamily()), 1)
		}
	})

	bktree := NewBKDelTree()
	words := []string{"abcde", "abcdd", "abcdf", "acbed", "acebd", "cabde", "dadbc", "cdacb", "cadbc", "bacde"}
	for _, word := range words {
		_ = bktree.Insert(word)
	}

	tests := []struct {
		name string
		word string
		want int
	}{
		{name: "test1", word: "abcde", want: 10},
		{name: "test2", word: "abcdd", want: 2},
		{name: "test3", word: "abcdf", want: 1},
		{name: "test4", word: "acbed", want: 1},
		{name: "test5", word: "acebd", want: 3},
		{name: "test6", word: "cabde", want: 2},
		{name: "test7", word: "dadbc", want: 2},
		{name: "test8", word: "cdacb", want: 1},
		{name: "test9", word: "cadbc", want: 1},
		{name: "test10", word: "bacde", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			if node := bktree.Search(tt.word, 0); len(node) != 0 {
				if len(node[0].collectFamily()) != tt.want {
					t1.Errorf("BKNode.collectFamily(%s) = %v, want %v", node[0].word, len(node[0].collectFamily()), tt.want)
				}
			}
		})
	}
}

func TestBKDelTreeDelete(t *testing.T) {
	t.Run("NilBKTree", func(t1 *testing.T) {
		nilbktree := NewBKDelTree()
		if nilbktree.Delete("abc") != nil {
			t1.Errorf("BKDelTree.Delete() = %v, want %v", nilbktree.Delete("abc"), nil)
		}
	})

	bktree := NewBKDelTree()
	words := []string{"abcde", "abcdd", "abcdf", "acbed", "acebd", "cabde", "dadbc", "cdacb", "cadbc", "bacde"}
	for _, word := range words {
		_ = bktree.Insert(word)
	}
	t.Run("deleteLeaf", func(t1 *testing.T) {
		if err := bktree.Delete("abcdf"); err != nil {
			t1.Errorf("BKDelTree.Delete() = %v, want %v", err, nil)
		}
		if len(bktree.Search("abcdf", 0)) != 0 {
			t1.Errorf("BKDelTree.Delete(abcdf) = %v, want %v", len(bktree.Search("abcdf", 0)), 0)
		}
	})

	t.Run("deleteMid", func(t1 *testing.T) {
		if err := bktree.Delete("acebd"); err != nil {
			t1.Errorf("BKDelTree.Delete(acebd) = %v, want %v", err, nil)
		}
		if len(bktree.Search("acebd", 0)) != 0 {
			t1.Errorf("BKDelTree.Delete(acebd) = %v, want %v", len(bktree.Search("acebd", 0)), 0)
		}
	})

	t.Run("deleteRoot", func(t1 *testing.T) {
		if err := bktree.Delete("abcde"); err != nil {
			t1.Errorf("BKDelTree.Delete(abcde) = %v, want %v", err, nil)
		}
		if len(bktree.Search("abcde", 0)) != 0 {
			t1.Errorf("BKDelTree.Delete(abcde) = %v, want %v", len(bktree.Search("abcde", 0)), 0)
		}
	})

	t.Run("deleteNonExistent", func(t1 *testing.T) {
		if err := bktree.Delete("xyz"); err == nil {
			t1.Errorf("BKDelTree.Delete(abcde) = %v, want %v", err, "non nil")
		}
	})
}

func TestBKDelTree_SPPrintf(t *testing.T) {
	t.Run("NilBKTree", func(t1 *testing.T) {
		nilbktree := NewBKDelTree()
		if nilbktree.SPPrintf("") != "" {
			t1.Errorf("BKDelTree.SPPrintf() = %v, want %v", nilbktree.SPPrintf(" "), "")
		}
	})
	t.Run("StandardBKTree", func(t1 *testing.T) {
		bktree := NewBKDelTree()
		words := []string{"abcde", "abcdd", "abcdf", "acbed", "acebd", "cabde", "dadbc", "cdacb", "cadbc", "bacde"}
		for _, word := range words {
			_ = bktree.Insert(word)
		}
		got := bktree.SPPrintf(". ")
		// count num of lines in got
		numlines := strings.Count(got, "\n")
		if numlines != 10 {
			t1.Errorf("BKDelTree.SPPrintf() = %v, want %v", numlines, 10)
		}
	})
}

func benchmarkBKDelTreeInsertN(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		bktree := NewBKDelTree()
		for i := 0; i < n; i++ {
			word := GenerateRandomString(i)
			_ = bktree.Insert(word)
		}
	}
}

func BenchmarkBKDelTreeInsert(b *testing.B) {
	numNodes := []int{1, 10, 100, 1000, 10_000, 20_000, 30_000, 40_000, 50_000} //, 75_000, 100_000, 1000_000}
	for _, n := range numNodes {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			benchmarkBKDelTreeInsertN(n, b)
		})
	}
}

func benchmarkBKDelTreeSearchN(bktree *BKDelTree, word string, distance int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		bktree.Search(word, distance)
	}
}

func BenchmarkBKDelTreeSearch(b *testing.B) {
	bktree := NewBKDelTree()
	for i := 0; i < 1000; i++ {
		word := GenerateRandomString(i)
		_ = bktree.Insert(word)
	}

	wordsDistance := map[string]int{ // word: distance
		"M1YK1G":   0,
		"I4FH9EV":  1,
		"ZUHX11":   2,
		"C6OJ6JI8": 0,
		"Z15BT6VJ": 1,
		"17569R1":  2,
		"DEADBEEF": 0,
	}

	for word, distance := range wordsDistance {
		b.Run(fmt.Sprintf("%s:%d", word, distance), func(b1 *testing.B) {
			for i := 0; i < b1.N; i++ {
				benchmarkBKDelTreeSearchN(bktree, word, distance, b)
			}
		})
	}
}

func BenchmarkBKDelTree_Delete(b *testing.B) {
	bktree := NewBKDelTree()
	for i := 0; i < 1000; i++ {
		word := GenerateRandomString(i)
		_ = bktree.Insert(word)
	}

	words := []string{
		"M1YK1G",
		"I4FH9EV",
		"ZUHX11",
		"C6OJ6JI8",
		"Z15BT6VJ",
		"17569R1",
		"DEADBEEF",
		"7NBOZQR", // root
	}

	for _, word := range words {
		b.Run(word, func(b1 *testing.B) {
			for i := 0; i < b1.N; i++ {
				_ = bktree.Delete(word)
			}
		})
	}
}
