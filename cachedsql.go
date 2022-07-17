package cachedmysql

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type DB struct {
	*sql.DB
	Cache *cache
}

func (db *DB) ExecDel(query string, args ...any) func(keys ...string) error {
	return func(keys ...string) error {
		err := db.Cache.Del(keys...)
		if err != nil {
			return err
		}
		_, err = db.Exec(query, args...)
		return err
	}
}

func (db *DB) ExecSet(query string, args ...any) func(val any, keys ...string) error {
	return func(val any, keys ...string) error {
		if len(keys) < 1 {
			return errors.New("key missing")
		}
		err := db.Cache.Set(keys[0], val)
		if err != nil {
			return err
		}
		_, err = db.Exec(query, args...)
		return err
	}
}

func (db *DB) QueryGet(query string, args ...any) func(val any, keys ...string) error {
	return func(val any, keys ...string) error {
		if len(keys) < 1 {
			return errors.New("key missing")
		}
		return db.Cache.Get(keys[0], val, func(val interface{}) error {
			var data []byte
			rows, err := db.Query(query, args)
			if err != nil {
				return err
			}
			err = rows.Scan(&data)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, val)
		})
	}
}
