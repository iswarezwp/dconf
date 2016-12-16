package dconf

import (
	"os"
	"testing"
	"time"
)

const (
	TEST_FILE  = "test.conf"
	TEST_CONF1 = `
[Default]
testKey1=TestValue1
testKey2=TestValue2
`
	TEST_CONF2 = `
[Default]
testKey1=TestValue2
`
	TEST_CONF3 = `
[Default]
testKey1=TestValue1

[TestSec]
testKey2=TestValue2
`

	TEST_CONF4 = `
testKey1=TestValue1

[TestSec]
testKey2=TestValue2
`
)

func createTestConfFile(t *testing.T, content string) {
	fp, err := os.Create(TEST_FILE)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	if _, err := fp.WriteString(content); err != nil {
		t.Fatal(err)
	}
}

func Test(t *testing.T) {
	createTestConfFile(t, TEST_CONF1)
	conf, err := NewDConf(TEST_FILE, true)
	if err != nil {
		t.Fatal(err)
	}

	// ---
	v := conf.Get("testKey1", "defaultValue")
	if v != "TestValue1" {
		t.Fatalf("Expected: TestValue1, got: %s", v)
	}

	v = conf.GetValue("Default", "testKey1", "defaultValue")
	if v != "TestValue1" {
		t.Fatalf("Expected: TestValue1, got: %s", v)
	}

	v = conf.Get("testKey2", "defaultValue")
	if v != "TestValue2" {
		t.Fatalf("Expected: TestValue2, got: %s", v)
	}

	v = conf.GetValue("TestSec", "testKey2", "defaultValue")
	if v != "defaultValue" {
		t.Fatalf("Expected: defaultValue, got: %s", v)
	}

	v = conf.Get("testKey3", "defaultValue")
	if v != "defaultValue" {
		t.Fatalf("Expected: defaultValue, got: %s", v)
	}

	// ---
	createTestConfFile(t, TEST_CONF2)
	time.Sleep(1 * time.Second)
	v = conf.Get("testKey1", "defaultValue")
	if v != "TestValue2" {
		t.Fatalf("Expected: TestValue2, got: %s", v)
	}

	// ---
	createTestConfFile(t, TEST_CONF3)
	time.Sleep(1 * time.Second)
	v = conf.Get("testKey1", "defaultValue")
	if v != "TestValue1" {
		t.Fatalf("Expected: TestValue1, got: %s", v)
	}

	v = conf.Get("testKey2", "defaultValue")
	if v != "defaultValue" {
		t.Fatalf("Expected: defaultValue, got: %s", v)
	}

	v = conf.GetValue("TestSec", "testKey2", "defaultValue")
	if v != "TestValue2" {
		t.Fatalf("Expected: TestValue2, got: %s", v)
	}

	// ---
	createTestConfFile(t, TEST_CONF4)
	time.Sleep(1 * time.Millisecond)
	v = conf.Get("testKey1", "defaultValue")
	if v != "TestValue1" {
		t.Fatalf("Expected: TestValue1, got: %s", v)
	}

	os.Remove(TEST_FILE)
}
