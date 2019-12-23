package services

import "testing"

func TestWrongTargetName(t *testing.T) {
	var err error

	_, err = TargetFromString("nope")
	switch err.(type) {
	case ServiceNotDefined:
		t.Log("ServiceNotDefined thrown as it should be")
	default:
		panic("TargetFromString() doesn't throw ServiceNotDefined")
	}
}

func TestServiceThatIsNotTarget(t *testing.T) {
	var err error

	_, err = TargetFromString("ldap")
	switch err.(type) {
	case TargetNotDefined:
		t.Log("TargetNotDefined thrown as it should be")
	default:
		panic("TargetFromString() doesn't throw TargetNotDefined")
	}
}
