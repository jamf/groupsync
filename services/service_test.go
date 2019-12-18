package services

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDiffWithEmptySrc(t *testing.T) {
	srcGrp := buildMockUsers(0, 0)
	tarGrp := buildMockUsers(0, 3)

	_, err := Diff(srcGrp, tarGrp, "mockservice")
	switch err.(type) {
	case SourceGroupEmptyError:
	default:
		panic("diff should return a SourceGroupEmptyError on empty source group")
	}
}

func TestDiffWithOverlap(t *testing.T) {
	srcGrp := buildMockUsers(0, 3)
	tarGrp := buildMockUsers(1, 4)

	expectedRem := []User{
		tarGrp[2],
	}

	expectedAdd := []User{
		srcGrp[0],
	}

	diff, _ := Diff(srcGrp, tarGrp, "mockservice")

	assertDiff(expectedRem, diff.Rem, expectedAdd, diff.Add)
}

func TestDiffIdenticalGroups(t *testing.T) {
	srcGrp := buildMockUsers(0, 3)
	tarGrp := buildMockUsers(0, 3)

	expectedRem := []User{}

	expectedAdd := []User{}

	diff, _ := Diff(srcGrp, tarGrp, "mockservice")

	assertDiff(expectedRem, diff.Rem, expectedAdd, diff.Add)
}

// Helpers

func buildMockUsers(start, end uint32) []User {
	var result []User

	for i := start; i < end; i++ {
		result = append(result, newUser())

		result[len(result)-1].addIdentity("mockservice", newMockIdentity(i))
	}

	return result
}

func assertDiff(expectedRem, rem, expectedAdd, add []User) {
	expectedRem = sanitize(expectedRem)
	rem = sanitize(rem)
	expectedAdd = sanitize(expectedAdd)
	add = sanitize(add)

	if reflect.DeepEqual(expectedRem, rem) &&
		reflect.DeepEqual(expectedAdd, add) {
		return
	}

	fmt.Printf(
		"Expected rem: %+v\nActual rem: %+v\nExpected add: %+v\nActual add: %+v\n",
		expectedRem,
		rem,
		expectedAdd,
		add,
	)

	panic("expected and actual rem/add slices didn't match")
}

func sanitize(u []User) []User {
	// because the fact slices can be either nil or empty in Golang and they
	// behave differently in only some fringe circumstances is absolutely awful

	if u == nil {
		return make([]User, 0)
	}

	return u
}
