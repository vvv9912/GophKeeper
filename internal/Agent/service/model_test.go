package service

import "testing"

func TestNewPath(t *testing.T) {

	NewPath("testPathStorage")
	if PathStorage != "testPathStorage" {
		t.Errorf("PathStorage is not set correctly")
	}

	NewPath("testPathStorage", "testPathTmp")
	if PathStorage != "testPathStorage" {
		t.Errorf("PathStorage is not set correctly")
	}
	if PathTmp != "testPathTmp" {
		t.Errorf("testPathTmp is not set correctly")
	}

	NewPath("testPathStorage", "testPathTmp", "testPathUserData")
	if PathStorage != "testPathStorage" {
		t.Errorf("PathStorage is not set correctly")
	}
	if PathTmp != "testPathTmp" {
		t.Errorf("testPathTmp is not set correctly")
	}
	if PathUserData != "testPathUserData" {
		t.Errorf("PathUserData is not set correctly")
	}
}
