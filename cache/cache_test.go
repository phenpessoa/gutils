package cache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	for _, tc := range []struct {
		name     string
		cacheFor time.Duration
		f        func(*Cache[string, any]) error
	}{
		{
			"get set delete",
			0,
			func(c *Cache[string, any]) error {
				c.Set("foo", 1)
				v, ok := c.Get("foo")
				if !ok {
					return fmt.Errorf("key foo not found in Cache")
				}

				if v != 1 {
					return fmt.Errorf("expected foo value to be 1 but got %v instead", 1)
				}

				c.Delete("foo")

				_, ok = c.Get("foo")

				if ok {
					return fmt.Errorf("key foo found in Cache")
				}

				return nil
			},
		},
		{
			"concurrent reads and writes",
			0,
			func(c *Cache[string, any]) error {
				wg := new(sync.WaitGroup)
				wg.Add(200)
				for i := 0; i < 100; i++ {
					go func(i int) {
						defer wg.Done()
						c.Set(strconv.Itoa(i), i)
					}(i)

					go func(i int) {
						defer wg.Done()
						c.Get(strconv.Itoa(i))
					}(i)
				}
				wg.Wait()
				return nil
			},
		},
		{
			"gc",
			25 * time.Millisecond,
			func(c *Cache[string, any]) error {
				c.Set("foo", 1)
				time.Sleep(60 * time.Millisecond)
				_, ok := c.Get("foo")
				if ok {
					return fmt.Errorf("key foo found in Cache")
				}
				return nil
			},
		},
		{
			"get set",
			0,
			func(c *Cache[string, any]) error {
				i := new(int)
				i2 := new(int)

				iGetSet := c.GetSet("foo", i)
				i2GetSet := c.GetSet("foo", i2)

				if iGetSet != i {
					return fmt.Errorf("iGetSet is not equal to i")
				}

				if iGetSet == i2 {
					return fmt.Errorf("iGetSet is equal to i2")
				}

				if i2GetSet != i {
					return fmt.Errorf("i2GetSet is not equal to i")
				}

				if i2GetSet == i2 {
					return fmt.Errorf("i2GetSet is equal to i2")
				}
				return nil
			},
		},
		{
			"wipe and len",
			0,
			func(c *Cache[string, any]) error {
				c.Set("foo", 1)
				c.Set("bar", 2)

				if c.Len() != 2 {
					return fmt.Errorf("cache len is not 2")
				}

				c.Wipe()

				if c.Len() != 0 {
					return fmt.Errorf("cache len is not 0")
				}

				return nil
			},
		},
		{
			"contains",
			0,
			func(c *Cache[string, any]) error {
				ok := c.Contains("foo")
				if ok {
					return fmt.Errorf("key foo found in cache")
				}

				c.Set("foo", 1)

				ok = c.Contains("foo")
				if !ok {
					return fmt.Errorf("key foo not found in cache")
				}

				return nil
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCache[string, any](tc.cacheFor)
			if err := tc.f(c); err != nil {
				t.Errorf("\ntest '%s' failed\nerr: %v", tc.name, err)
			}
		})
	}
}
