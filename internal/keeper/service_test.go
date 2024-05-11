package keeper

import (
	"testing"
	"time"
)

func TestStoring(t *testing.T) {
	k := NewService(10 * time.Second)
	k.Set("key1", []byte("data"), 0)
	value := k.Get("key1")
	if len(value) == 0 {
		t.Error("value of key1 empty but shoudnt")
	}
}

func TestKeyOverwrite(t *testing.T) {
	k := NewService(10 * time.Second)
	k.Set("key1", []byte("data1"), 0)

	wantData := []byte("data2")
	k.Set("key1", wantData, 0)
	gotData := k.Get("key1")
	if string(gotData) != string(wantData) {
		t.Errorf("want data %s, but got %s", string(wantData), string(gotData))
	}
}

func TestDelete(t *testing.T) {
	k := NewService(10 * time.Second)
	k.Set("key1", []byte("data1"), 0)

	k.Delete("key1")
	value := k.Get("key1")
	if len(value) != 0 {
		t.Errorf("value of key1 not empty but shoud")
	}
}

func TestExpirationOneEntry(t *testing.T) {
	t.Parallel()

	k := NewService(0)

	k.Set("key1", []byte("data"), 50*time.Millisecond)
	k.Run()

	time.Sleep(51 * time.Millisecond)

	value := k.Get("key1")
	if len(value) != 0 {
		t.Error("value of key1 not empty but shoud")
	}
}

func TestExpirationMultipleEntries(t *testing.T) {
	t.Parallel()

	k := NewService(0)

	k.Set("key1", []byte("data"), 50*time.Millisecond)
	k.Set("key2", []byte("data"), 1*time.Second)
	k.Set("key3", []byte("data"), 10*time.Second)
	k.Set("key4", []byte("data"), 1*time.Nanosecond)

	k.Run()

	time.Sleep(500 * time.Millisecond)
	value := k.Get("key1")
	if len(value) != 0 {
		t.Error("value of key1 not empty but shoud")
	}

	value = k.Get("key2")
	if len(value) == 0 {
		t.Error("value of key2 empty but shoudnt")
	}

	value = k.Get("key3")
	if len(value) == 0 {
		t.Error("value of key3 empty but shoudnt")
	}

	value = k.Get("key4")
	if len(value) != 0 {
		t.Error("value of key4 not empty but shoud")
	}
}
