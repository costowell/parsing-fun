package common

import "testing"

func TestOrderedSetInsert(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Insert(10)
	if set.Data[0] != 10 {
		t.Error("Failed to insert single element")
	}
	set.Insert(12)
	if set.Data[1] != 12 {
		t.Error("Failed to insert second element")
	}
	set.Insert(10)
	if len(set.Data) != 2 {
		t.Error("Duplicate data found")
	}
}

func TestOrderedSetContains(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Insert(10)
	set.Insert(12)
	if !set.Contains(10) {
		t.Error("Set does not contain element just inserted")
	}
	if set.Contains(13) {
		t.Error("Set contains element never inserted")
	}
}

func TestOrderedSetRemove(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Insert(10)
	set.Insert(12)
	set.Insert(14)
	set.Remove(12)
	if set.Contains(12) {
		t.Error("Failed to remove from set")
	}
	if !set.Contains(14) || !set.Contains(10) {
		t.Error("Extra elements removed")
	}
}
