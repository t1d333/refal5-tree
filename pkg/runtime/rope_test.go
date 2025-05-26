package runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyRopeBuilding(t *testing.T) {
	r := NewRope([]R5Node{})

	if r.Len() != 0 {
		t.Error("Empty rope must have zero length")
	}
}

// func TestRopeSplit(t *testing.T) {
// 	r1 := NewRope([]R5Node{&R5NodeChar{Char: '1'}})
// 	r2 := NewRope([]R5Node{&R5NodeNumber{Number: 2}})
// 	r3 := NewRope([]R5Node{&R5NodeString{Value: "3"}})
// 	r4 := NewRope([]R5Node{&R5NodeChar{Char: '4'}})
// 	r5 := NewRope([]R5Node{&R5NodeChar{Char: '5'}})
// 	r6 := NewRope([]R5Node{&R5NodeChar{Char: '6'}})
// 	r7 := NewRope([]R5Node{&R5NodeChar{Char: '7'}})
// 	r8 := NewRope([]R5Node{&R5NodeChar{Char: '8'}})
// 	r9 := NewRope([]R5Node{&R5NodeChar{Char: '9'}})
// 	r10 := NewRope([]R5Node{&R5NodeChar{Char: 'z'}})
// 	r11 := NewRope([]R5Node{&R5NodeChar{Char: 'z'}})
//
// 	tmp := r1.ConcatAVL(r2)
// 	tmp = tmp.ConcatAVL(r3)
// 	tmp = tmp.ConcatAVL(r4)
// 	tmp = tmp.ConcatAVL(r4)
// 	tmp = tmp.ConcatAVL(r5)
// 	tmp = tmp.ConcatAVL(r6)
// 	tmp = tmp.ConcatAVL(r7)
// 	tmp = tmp.ConcatAVL(r8)
// 	tmp = tmp.ConcatAVL(r9)
// 	tmp = tmp.ConcatAVL(r10)
// 	tmp = tmp.ConcatAVL(r11)
//
// 	l, a := tmp.Split(5)
//
// 	fmt.Println(a.Len(), a.String(), "-----", l.String())
// 	// assert.Fail(t, "fail")
// }

func TestNonEmptyRopeBuilding(t *testing.T) {
	expected_length := 3
	expected_height := 1
	r := NewRope(
		[]R5Node{&R5NodeChar{Char: 'a'}, &R5NodeNumber{Number: 5}, &R5NodeString{Value: "s"}},
	)

	if r.Len() != expected_length {
		t.Errorf("Rope length expected: %d, but got %d", expected_length, r.Len())
	}

	if r.Height() != 1 {
		t.Errorf("Rope height expected: %d, but got %d", expected_height, r.Height())
	}
}

func TestRopeGetMethod(t *testing.T) {
	expected_first := &R5NodeChar{Char: 'a'}
	expected_second := &R5NodeNumber{Number: 5}
	expected_third := &R5NodeString{Value: "s"}
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

func TestRopeBalanceAVLFactor(t *testing.T) {

	r1 := NewRope([]R5Node{&R5NodeString{Value: "a"}})
	r2 := NewRope([]R5Node{&R5NodeString{Value: "b"}})
	r3 := NewRope([]R5Node{&R5NodeString{Value: "c"}})
	r4 := NewRope([]R5Node{&R5NodeString{Value: "d"}})
	r5 := NewRope([]R5Node{&R5NodeString{Value: "d"}})
	r6 := NewRope([]R5Node{&R5NodeString{Value: "e"}})

	tmp := r1.ConcatAVL(r2)
	assert.NotNil(t, tmp)
	//
	tmp = tmp.ConcatAVL(r3)
	assert.NotNil(t, tmp)
	//
	tmp = tmp.ConcatAVL(r4)
	assert.NotNil(t, tmp)
	
	b := r5.ConcatAVL(tmp)
	b = r6.ConcatAVL(b)
	VisualizeRope(b, 0)
	fmt.Println(12321312, b.String())

	k, p := b.Split(5)
	VisualizeRope(k, 0)
	fmt.Println("------------")
	VisualizeRope(p, 0)
	assert.False(t, tmp.IsAVLBalanced())

	// assert.False(t, tmp.IsAVLBalanced())
	//
	// balanced := tmp.balanceAVL()
	//
	// assert.True(t, balanced.IsAVLBalanced())
	// assert.Equal(t, tmp.Len(), balanced.Len())
	// assert.Equal(t, expectedBalancedHeight, balanced.Height())
	//
	// for i := 0; i < tmp.Len(); i++ {
	// 	assert.Equal(t, tmp.Get(i), balanced.Get(i))
	// }
}
//
// func TestRopeConcatWithRebalance2(t *testing.T) {
// 	// expectedBalancedHeight := 2
//
// 	r1 := NewRope([]R5Node{&R5NodeString{String: "a"}})
// 	r2 := NewRope([]R5Node{&R5NodeString{String: "b"}})
// 	r3 := NewRope([]R5Node{&R5NodeString{String: "c"}})
// 	r4 := NewRope([]R5Node{&R5NodeString{String: "d"}})
// 	r5 := NewRope([]R5Node{&R5NodeString{String: "e"}})
// 	// r6 := NewRope([]R5Node{&R5NodeString{String: "f"}})
//
// 	tmp := r1.Concat(r2)
// 	fmt.Println(tmp.Height())
// 	assert.NotNil(t, tmp)
//
// 	tmp = r3.Concat(tmp)
// 	fmt.Println(tmp.Height())
// 	assert.NotNil(t, tmp)
//
// 	tmp = tmp.Concat(r4)
// 	tmp = tmp.balanceAVL()
// 	assert.NotNil(t, tmp)
// 	fmt.Println("213", tmp.Height())
//
// 	tmp = tmp.Concat(r5)
// 	assert.NotNil(t, tmp)
//
// 	assert.False(t, tmp.IsAVLBalanced())
//
// 	fmt.Println(tmp.Height())
// 	balanced := tmp.balanceAVL()
// 	fmt.Println(balanced.Height())
//
// 	assert.True(t, balanced.IsAVLBalanced())
// 	assert.Equal(t, 1, 2)
// 	// assert.Equal(t, tmp.Len(), balanced.Len())
// 	// assert.Equal(t, expectedBalancedHeight, balanced.Height())
//
// 	// for i := 0; i < tmp.Len(); i++ {
// 		// assert.Equal(t, tmp.Get(i), balanced.Get(i))
// 	// }
// }
