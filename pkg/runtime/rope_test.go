package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyRopeBuilding(t *testing.T) {
	r := NewRope([]R5Node{})

	if r == nil {
		t.Errorf(`Got nil Rope`)
	}

	if r.Len() != 0 {
		t.Error("Empty rope must have zero length")
	}

	if r.Height() != 0 {
		t.Error("Empty rope must have zero height")
	}
}

func TestNonEmptyRopeBuilding(t *testing.T) {
	expected_length := 3
	expected_height := 1
	r := NewRope(
		[]R5Node{&R5NodeChar{Char: 'a'}, &R5NodeNumber{Number: 5}, &R5NodeString{String: "s"}},
	)

	if r == nil {
		t.Errorf(`Got nil Rope`)
	}

	if r.Len() != expected_length {
		t.Errorf("Rope length expected: %d, but got %d", expected_length, r.Len())
	}

	if r.Height() != 0 {
		t.Errorf("Rope heith expected: %d, but got %d", expected_height, r.Height())
	}
}

func TestRopeGetMethod(t *testing.T) {
	expected_first := &R5NodeChar{Char: 'a'}
	expected_second := &R5NodeNumber{Number: 5}
	expected_third := &R5NodeString{String: "s"}
	r := NewRope([]R5Node{expected_first, expected_second, expected_third})

	first := r.Get(0)
	assert.NotNilf(t, first, "Expected %#v, but got nil", *expected_first)

	second := r.Get(1)
	assert.NotNilf(t, second, "Expected %#v, but got nil", *expected_second)

	third := r.Get(2)
	assert.NotNilf(t, third, "Expected %#v, but got nil", *expected_third)

	fourth := r.Get(3)
	assert.Nil(t, fourth)

	assert.Equal(t, expected_first, first)
	assert.Equal(t, expected_second, second)
	assert.Equal(t, expected_third, third)
}

func TestRopeSetMethod(t *testing.T) {
	
}

func TestRopeConcatWithoutRebalance(t *testing.T) {
	
}

func TestRopeBalance(t *testing.T) {
	expectedFirst :=&R5NodeChar{Char: 'a'}
	expectedSecond := &R5NodeNumber{Number: 5}
	expectedThird := &R5NodeString{String: "s"}
	expectedFourth := &R5NodeChar{Char: 'b'}
	expectedBalancedHeight := 2

	r1 := NewRope([]R5Node{expectedFirst})
	r2 := NewRope([]R5Node{expectedSecond})
	r3 := NewRope([]R5Node{expectedThird})
	r4 := NewRope([]R5Node{expectedFourth})

	tmp := r1.Concat(r2)
	assert.NotNil(t, tmp)

	tmp = tmp.Concat(r3)
	assert.NotNil(t, tmp)

	tmp = tmp.Concat(r4)
	assert.NotNil(t, tmp)

	assert.False(t, tmp.IsBalanced())

	balanced := tmp.Balance()
	
	assert.True(t, balanced.IsBalanced())
	assert.Equal(t, tmp.Len(), balanced.Len())
	assert.Equal(t, expectedBalancedHeight, balanced.Height())
	assert.Equal(t, expectedFirst, balanced.Get(0))
	assert.Equal(t, expectedSecond, balanced.Get(1))
	assert.Equal(t, expectedThird, balanced.Get(2))
	assert.Equal(t, expectedFourth, balanced.Get(3))
}

func TestRopeConcatWithRebalance(t *testing.T) {
	
}
