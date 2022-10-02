package gocache

import (
	"reflect"
	"testing"
)

func addSth(testCache *cache, key string, value string) {
	testValue := ByteView{
		b: []byte(value),
	}
	testCache.add(key, testValue)
}

func TestAdd(t *testing.T) {
	testCache := cache{}
	for i := 0; i < 5; i++ {
		letter := string(rune(int('a') + i))
		// go addSth(&testCache, letter, "test_"+letter)
		addSth(&testCache, letter, "test_"+letter)
	}
}

func TestGet(t *testing.T) {
	testCache := cache{}
	for i := 0; i < 5; i++ {
		letter := string(rune(int('a') + i))
		// go addSth(&testCache, letter, "test_"+letter)
		addSth(&testCache, letter, "test_"+letter)
	}

	for i := 0; i < 5; i++ {
		key := string(rune(int('a') + i))
		value, ok := testCache.get(key)
		if !ok || string(value.b) != "test_"+key {
			t.Fatalf("get no key " + key)
		}
	}

	//test key not exist conditions
	_, ok := testCache.get("key_not_exist")
	if ok {
		t.Fatalf("get wrong value when key not exist")
	}
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
