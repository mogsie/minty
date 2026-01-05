package minty

import (
	"testing"
)

func TestStr(t *testing.T) {
	m := map[string]interface{}{
		"title":  "Hello",
		"views":  42,
		"rating": 4.5,
	}

	tests := []struct {
		key  string
		want string
	}{
		{"title", "Hello"},
		{"views", "42"},
		{"rating", "4.5"},
		{"missing", ""},
	}

	for _, tt := range tests {
		got := Str(m, tt.key)
		if got != tt.want {
			t.Errorf("Str(m, %q) = %q, want %q", tt.key, got, tt.want)
		}
	}

	// Test nil map
	if got := Str(nil, "key"); got != "" {
		t.Errorf("Str(nil, key) = %q, want empty", got)
	}
}

func TestInt(t *testing.T) {
	m := map[string]interface{}{
		"views":   42,
		"rating":  4.7,
		"count":   "100",
		"invalid": "abc",
	}

	tests := []struct {
		key  string
		want int
	}{
		{"views", 42},
		{"rating", 4},
		{"count", 100},
		{"invalid", 0},
		{"missing", 0},
	}

	for _, tt := range tests {
		got := Int(m, tt.key)
		if got != tt.want {
			t.Errorf("Int(m, %q) = %d, want %d", tt.key, got, tt.want)
		}
	}
}

func TestBool(t *testing.T) {
	m := map[string]interface{}{
		"active":   true,
		"disabled": false,
		"notBool":  "yes",
	}

	if !Bool(m, "active") {
		t.Error("Bool(m, active) = false, want true")
	}
	if Bool(m, "disabled") {
		t.Error("Bool(m, disabled) = true, want false")
	}
	if Bool(m, "notBool") {
		t.Error("Bool(m, notBool) = true, want false")
	}
	if Bool(m, "missing") {
		t.Error("Bool(m, missing) = true, want false")
	}
}

func TestTruthy(t *testing.T) {
	tests := []struct {
		val  interface{}
		want bool
	}{
		{nil, false},
		{false, false},
		{true, true},
		{0, false},
		{1, true},
		{-1, true},
		{0.0, false},
		{0.1, true},
		{"", false},
		{"hello", true},
		{[]interface{}{}, false},
		{[]interface{}{1}, true},
		{map[string]interface{}{}, false},
		{map[string]interface{}{"a": 1}, true},
	}

	for _, tt := range tests {
		got := Truthy(tt.val)
		if got != tt.want {
			t.Errorf("Truthy(%v) = %v, want %v", tt.val, got, tt.want)
		}
	}
}

func TestEqGt(t *testing.T) {
	post := map[string]interface{}{
		"status": "published",
		"likes":  42,
	}

	if !Eq(post, "status", "published") {
		t.Error("Eq should match")
	}
	if Eq(post, "status", "draft") {
		t.Error("Eq should not match")
	}

	if !Gt(post, "likes", 0) {
		t.Error("Gt(42 > 0) should be true")
	}
	if !Gt(post, "likes", 41) {
		t.Error("Gt(42 > 41) should be true")
	}
	if Gt(post, "likes", 42) {
		t.Error("Gt(42 > 42) should be false")
	}
	if Gt(post, "likes", 100) {
		t.Error("Gt(42 > 100) should be false")
	}
}

func TestFilterItems(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"id": 1, "status": "published"},
		map[string]interface{}{"id": 2, "status": "draft"},
		map[string]interface{}{"id": 3, "status": "published"},
	}

	published := FilterItems(posts, Where("status", "published"))
	if len(published) != 2 {
		t.Errorf("FilterItems published: got %d, want 2", len(published))
	}

	drafts := FilterItems(posts, Where("status", "draft"))
	if len(drafts) != 1 {
		t.Errorf("FilterItems draft: got %d, want 1", len(drafts))
	}
}

func TestCount(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"status": "published"},
		map[string]interface{}{"status": "draft"},
		map[string]interface{}{"status": "published"},
		map[string]interface{}{"status": "published"},
	}

	if got := Count(posts, Where("status", "published")); got != 3 {
		t.Errorf("Count published = %d, want 3", got)
	}
	if got := Count(posts, Where("status", "draft")); got != 1 {
		t.Errorf("Count draft = %d, want 1", got)
	}
}

func TestSum(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"views": 100},
		map[string]interface{}{"views": 200},
		map[string]interface{}{"views": 50},
	}

	if got := Sum(posts, "views"); got != 350 {
		t.Errorf("Sum views = %d, want 350", got)
	}
}

func TestSortBy(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"title": "C", "views": 100},
		map[string]interface{}{"title": "A", "views": 300},
		map[string]interface{}{"title": "B", "views": 200},
	}

	// Sort by title ascending
	byTitle := SortBy(posts, "title", Asc)
	if Str(byTitle[0].(map[string]interface{}), "title") != "A" {
		t.Error("SortBy title Asc: first should be A")
	}
	if Str(byTitle[2].(map[string]interface{}), "title") != "C" {
		t.Error("SortBy title Asc: last should be C")
	}

	// Sort by views descending
	byViews := SortBy(posts, "views", Desc)
	if Int(byViews[0].(map[string]interface{}), "views") != 300 {
		t.Error("SortBy views Desc: first should have 300 views")
	}
	if Int(byViews[2].(map[string]interface{}), "views") != 100 {
		t.Error("SortBy views Desc: last should have 100 views")
	}

	// Original should be unchanged
	if Str(posts[0].(map[string]interface{}), "title") != "C" {
		t.Error("Original slice was modified")
	}
}

func TestPredicateCombinators(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"status": "published", "views": 100},
		map[string]interface{}{"status": "published", "views": 500},
		map[string]interface{}{"status": "draft", "views": 200},
	}

	// AND: published AND views > 200
	popularPublished := FilterItems(posts, And(
		Where("status", "published"),
		WhereGt("views", 200),
	))
	if len(popularPublished) != 1 {
		t.Errorf("And filter: got %d, want 1", len(popularPublished))
	}

	// OR: published OR views > 150
	combined := FilterItems(posts, Or(
		Where("status", "published"),
		WhereGt("views", 150),
	))
	if len(combined) != 3 {
		t.Errorf("Or filter: got %d, want 3", len(combined))
	}
}

func TestGroupItems(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"id": 1, "status": "published"},
		map[string]interface{}{"id": 2, "status": "draft"},
		map[string]interface{}{"id": 3, "status": "published"},
	}

	groups := GroupItems(posts, "status")

	if len(groups["published"]) != 2 {
		t.Errorf("GroupItems published: got %d, want 2", len(groups["published"]))
	}
	if len(groups["draft"]) != 1 {
		t.Errorf("GroupItems draft: got %d, want 1", len(groups["draft"]))
	}
}

func TestPluck(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"title": "Hello"},
		map[string]interface{}{"title": "World"},
	}

	titles := Pluck(posts, "title")
	if len(titles) != 2 {
		t.Errorf("Pluck: got %d items, want 2", len(titles))
	}
	if titles[0] != "Hello" || titles[1] != "World" {
		t.Errorf("Pluck: got %v, want [Hello World]", titles)
	}
}

func TestFind(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"id": 1, "status": "draft"},
		map[string]interface{}{"id": 2, "status": "published"},
		map[string]interface{}{"id": 3, "status": "published"},
	}

	found := Find(posts, Where("status", "published"))
	if found == nil {
		t.Error("Find: should have found item")
	}
	if Int(found.(map[string]interface{}), "id") != 2 {
		t.Error("Find: should return first match (id=2)")
	}

	notFound := Find(posts, Where("status", "archived"))
	if notFound != nil {
		t.Error("Find: should return nil for no match")
	}
}

func TestAnyAll(t *testing.T) {
	posts := []interface{}{
		map[string]interface{}{"status": "published"},
		map[string]interface{}{"status": "published"},
	}

	if !Any(posts, Where("status", "published")) {
		t.Error("Any: should be true")
	}
	if Any(posts, Where("status", "draft")) {
		t.Error("Any: should be false for draft")
	}

	if !All(posts, Where("status", "published")) {
		t.Error("All: should be true when all published")
	}

	posts = append(posts, map[string]interface{}{"status": "draft"})
	if All(posts, Where("status", "published")) {
		t.Error("All: should be false when one is draft")
	}
}
