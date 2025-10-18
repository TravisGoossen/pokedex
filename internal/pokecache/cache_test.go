package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "www.testurl.net",
			val: []byte("this is the testurl.net response!"),
		},
		{
			key: "www.youtube.com",
			val: []byte("blablah blah videos blah blah blah"),
		},
		{
			key: "www.cats.pic",
			val: []byte("cute cats!!"),
		},
	}

	for _, c := range cases {
		cache := NewCache(5 * time.Second)
		cache.Add(c.key, c.val)
		value, found := cache.Get(c.key)
		if !found {
			t.Error("cache entry not found")
			return
		}
		if string(value) != string(c.val) {
			t.Error("value isn't correct")
			return
		}
		if len(cache.Entries) != 1 {
			t.Error("len of cache entries is incorrect")
			return
		}
	}
}

func TestReapLoop(t *testing.T) {
	cases := []struct {
		interval time.Duration
		waitTime time.Duration
		key      string
		val      []byte
	}{
		{
			500 * time.Millisecond,
			600 * time.Millisecond,
			"https://www.cachetester.org",
			[]byte("this is cached text"),
		},
		{
			200 * time.Millisecond,
			350 * time.Millisecond,
			"https://www.LadyGaga.com",
			[]byte("Lots of fun songs here"),
		},
		{
			10 * time.Millisecond,
			11 * time.Millisecond,
			"https://www.thisOneIsFast!.com",
			[]byte("Hopefully go can get the 1 millisecond difference!"),
		},
	}

	for _, c := range cases {
		cache := NewCache(c.interval)
		cache.Add(c.key, c.val)
		time.Sleep(c.waitTime)
		_, found := cache.Get(c.key)
		if found {
			t.Error("cache found when it should have been reaped")
		}
	}
}
