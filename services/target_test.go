package services

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDiffWithEmptySrc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("diff should have panicked when given an empty source group")
		}
	}()

	srcGrp := buildMockUsers(0, 0)
	tarGrp := buildMockUsers(0, 3)

	Diff(srcGrp, tarGrp, "mockservice")
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

	rem, add := Diff(srcGrp, tarGrp, "mockservice")

	if len(rem) != len(expectedRem) {
		panic(fmt.Sprintf("The users-to-remove slice has the wrong length: %v", len(rem)))
	}

	if len(add) != len(expectedAdd) {
		panic(fmt.Sprintf("The users-to-add slice has the wrong length: %v", len(add)))
	}

	if !reflect.DeepEqual(expectedRem, rem) ||
		!reflect.DeepEqual(expectedAdd, add) {
		panic("!")
	}
}

func TestDiffIdenticalGroups(t *testing.T) {
	srcGrp := buildMockUsers(0, 3)
	tarGrp := buildMockUsers(0, 3)

	expectedRem := []User{}

	expectedAdd := []User{}

	rem, add := Diff(srcGrp, tarGrp, "mockservice")

	if len(rem) != len(expectedRem) {
		panic(fmt.Sprintf("The users-to-remove slice has the wrong length: %v", len(rem)))
	}

	if len(add) != len(expectedAdd) {
		panic(fmt.Sprintf("The users-to-add slice has the wrong length: %v", len(add)))
	}
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
