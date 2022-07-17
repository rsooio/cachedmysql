package cachedmysql

import (
	"sync"
)

type (
	call struct {
		wg  sync.WaitGroup
		val any
		err error
	}

	singlefight struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

var sf = singlefight{
	calls: make(map[string]*call),
}

func (sf *singlefight) do(key string, fn func() (any, error)) (any, error) {
	c, done := sf.newCall(key)
	if done {
		return c.val, c.err
	}

	sf.Call(c, key, fn)
	return c.val, c.err
}

func (sf *singlefight) newCall(key string) (c *call, done bool) {
	sf.lock.Lock()
	if c, ok := sf.calls[key]; ok {
		sf.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	sf.calls[key] = c
	sf.lock.Unlock()

	return c, false
}

func (sf *singlefight) Call(c *call, key string, fn func() (any, error)) {
	defer func() {
		sf.lock.Lock()
		delete(sf.calls, key)
		sf.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
