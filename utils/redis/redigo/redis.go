package redigo

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

func (r *redis) PING() (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("PING"))
}

func (r *redis) TTL(key string) (time.Duration, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	dur, err := redigo.Int64(conn.Do("TTL", key))
	return time.Duration(dur), err
}

func (r *redis) EXISTS(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redigo.Bool(conn.Do("EXISTS", key))
}

func (r *redis) EXPIRE(key string, duration time.Duration) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redigo.Int64(conn.Do("EXPIRE", key, strconv.FormatInt(int64(duration/time.Second), 10)))
	return err
}

func (r *redis) INCR(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	result, err := redigo.Int64(conn.Do("INCR", key))
	return result, err
}

func (r *redis) INCRBY(key string, val int) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redigo.Int64(conn.Do("INCRBY", key, val))
	return err
}

func (r *redis) SET(key string, value interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (r *redis) SETEX(key string, dur time.Duration, value interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, dur.Seconds(), value)
	return err
}

func (r *redis) GET(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("GET", key))
}

func (r *redis) MGET(key ...string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	payload := make([]interface{}, len(key))
	for n, c := range key {
		payload[n] = c
	}

	resp, err := redigo.Strings(conn.Do("MGET", payload...))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *redis) DEL(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redigo.Int64(conn.Do("DEL", key))
	return err
}

func (r *redis) HSET(key string, field string, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	return err
}

func (r *redis) HSETEX(key string, duration time.Duration, field string, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", key, field, value)
	if err != nil {
		return err
	} else {
		err = r.EXPIRE(key, duration)
	}

	return err
}

func (r *redis) HMSET(key string, data map[string]interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()
	var args []interface{}

	args = append(args, key)
	for field, val := range data {
		args = append(args, field)
		args = append(args, val)
	}

	_, err := conn.Do("HMSET", args...)
	return err
}

func (r *redis) HGET(key string, field string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("HGET", key, field))
}

//
func (r *redis) HGETP(keys []string, field string) (map[string]interface{}, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	// open pipeline, generate chain of commands
	for _, key := range keys {
		cmd := append([]string{key}, field)
		payload := make([]interface{}, len(cmd))
		for n, c := range cmd {
			payload[n] = c
		}
		// Send writes the command to the connection's output buffer
		if err := conn.Send("HGET", payload...); err != nil {
			return nil, err
		}

	}

	// Flush flushes the connection's output buffer to the server.
	if err := conn.Flush(); err != nil {
		return nil, err
	}

	//Receive reads a single reply from the server.
	replies := make([]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		reply, err := conn.Receive()
		if err != nil {
			return nil, err
		}
		replies[i] = reply
	}
	resp, err := redigo.Values(replies, nil)
	if err != nil {
		return nil, err
	}

	// build final result, assign each array string to each corresponding field
	results := make(map[string]interface{}, len(keys))
	for i, v := range resp {
		s, err := redigo.String(v, nil)
		if err != nil {
			continue
		}
		results[keys[i]] = s
	}
	return results, nil
}

func (r *redis) HMGET(key string, fields ...string) (map[string]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	l := len(fields)
	cmd := append([]string{key}, fields...)

	payload := make([]interface{}, len(cmd))
	for n, c := range cmd {
		payload[n] = c
	}

	result := make(map[string]string)
	resp, err := redigo.Strings(conn.Do("HMGET", payload...))
	if err != nil {
		return nil, err
	}

	var ctrFail int
	for i := 0; i < l; i++ {
		result[fields[i]] = resp[i]
		if resp[i] == "" {
			ctrFail++
		}
	}

	if ctrFail == l {
		return result, errors.New("Redis no result")
	}

	return result, nil
}

// HMGETP : is HMGET with pipeline redis
func (r *redis) HMGETP(keys []string, fields ...string) (map[string]interface{}, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	// open pipeline, generate chain of commands
	for _, key := range keys {
		cmd := append([]string{key}, fields...)
		payload := make([]interface{}, len(cmd))
		for n, c := range cmd {
			payload[n] = c
		}
		// Send writes the command to the connection's output buffer
		if err := conn.Send("HMGET", payload...); err != nil {
			return nil, err
		}

	}

	// Flush flushes the connection's output buffer to the server.
	if err := conn.Flush(); err != nil {
		return nil, err
	}

	//Receive reads a single reply from the server.
	replies := make([]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		reply, err := conn.Receive()
		if err != nil {
			return nil, err
		}
		replies[i] = reply
	}
	resp, err := redigo.Values(replies, nil)
	if err != nil {
		return nil, err
	}

	// build final result, assign each array string to each corresponding field
	results := make(map[string]interface{}, len(keys))
	for i, v := range resp {
		s, err := redigo.Strings(v, nil)
		if err != nil {
			continue
		}
		resultField := make(map[string]string, len(s))
		if len(s) == len(fields) {
			for i, sk := range s {
				resultField[fields[i]] = sk
			}
		}
		results[keys[i]] = resultField
	}

	return results, nil
}

func (r *redis) LPUSH(key string, value interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := redigo.Ints(conn.Do("LPUSH", key, value))
	if err != nil {
		return nil
	}
	return err
}

func (r *redis) LLEN(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	size, err := redigo.Int64(conn.Do("LLEN", key))
	return size, err
}

func (r *redis) LPUSHEX(key string, duration time.Duration, value interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := redigo.Int64(conn.Do("LPUSH", key, value))
	if err != nil {
		return err
	} else {
		err = r.EXPIRE(key, duration)
	}
	return err
}

func (r *redis) LPUSHEXP(value map[string][]interface{}, duration time.Duration) error {
	conn := r.Pool.Get()
	defer conn.Close()

	for key, val := range value {
		val = append([]interface{}{key}, val...)
		errPush := conn.Send("LPUSH", val...)
		if errPush != nil {
			return errPush
		}

		errExpire := conn.Send("EXPIRE", key, duration.Seconds())
		if errExpire != nil {
			return errExpire
		}
	}

	if err := conn.Flush(); err != nil {
		return err
	}

	return nil
}

func (r *redis) LREM(key string, count interface{}, value interface{}) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	removedIdx, err := redigo.Int64(conn.Do("LREM", key, count, value))
	if err != nil {
		return 0, nil
	}
	return removedIdx, err
}

func (r *redis) LRANGE(key string, start interface{}, stop interface{}) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	res, err := redigo.Strings(conn.Do("LRANGE", key, start, stop))
	if err != nil {
		return nil, nil
	}
	return res, nil
}

func (r *redis) LRANGEP(keys []string, fields ...string) (map[string][]int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	for _, key := range keys {
		cmd := append([]string{key}, fields...)
		payload := make([]interface{}, len(cmd))
		for n, c := range cmd {
			payload[n] = c
		}
		// Send writes the command to the connection's output buffer
		if err := conn.Send("LRANGE", payload...); err != nil {
			return nil, err
		}
	}

	// Flush flushes the connection's output buffer to the server.
	if err := conn.Flush(); err != nil {
		return nil, err
	}

	//Receive reads a single reply from the server.
	replies := make([]interface{}, len(keys))
	for i := 0; i < len(keys); i++ {
		reply, err := conn.Receive()
		if err != nil {
			return nil, err
		}
		replies[i] = reply
	}

	resp, err := redigo.Values(replies, nil)
	if err != nil {
		return nil, err
	}

	// build final result, assign each array string to each corresponding field
	results := make(map[string][]int64, len(keys))
	for i, arrOfID := range resp {
		var int64arr []int64
		arrOfID, _ := redigo.Ints(arrOfID, nil)
		for j := range arrOfID {
			if int64(arrOfID[j]) > 0 {
				int64arr = append(int64arr, int64(arrOfID[j]))
			}
		}
		results[keys[i]] = int64arr
	}
	return results, nil
}

func (r *redis) HGETALL(key string) (map[string]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.StringMap(conn.Do("HGETALL", key))
}

func (r *redis) HDEL(key string, fields ...string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	var args []interface{}
	args = append(args, key)
	for field, val := range fields {
		args = append(args, field)
		args = append(args, val)
	}

	_, err := conn.Do("HDEL", args...)
	return err
}

func (r *redis) HEXISTS(key string, field string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Bool(conn.Do("HEXISTS", key, field))
}

//ZADD
func (r *redis) ZADD(key string, fields ...Z) error {
	conn := r.Pool.Get()
	defer conn.Close()

	var err error
	for _, z := range fields {
		err = conn.Send("ZADD", key, fmt.Sprintf("%v", z.Score), z.Member)
	}
	defer conn.Flush()

	return err
}

func (r *redis) ZSCORE(key string, value interface{}) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redigo.Int(conn.Do("ZSCORE", key, value))
}

func (r *redis) ZREM(key string, fields ...string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	var sr []string
	sr = append(sr, key)
	sr = append(sr, fields...)

	args := make([]interface{}, len(sr))
	for i, s := range sr {
		args[i] = s
	}

	err := conn.Send("ZREM", args...)
	defer conn.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) ZRANGE(key string, start int, end int) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	resp, err := redigo.Strings(conn.Do("ZRANGE", key, start, end))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *redis) ZRANGEBYSCORE(key string, opt ZRangeByScore) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	args := []interface{}{key, opt.Min, opt.Max}
	if opt.Offset != 0 || opt.Count != 0 {
		args = append(
			args,
			"LIMIT",
			opt.Offset,
			opt.Count,
		)
	}

	resp, err := redigo.Strings(conn.Do("ZRANGEBYSCORE", args...))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *redis) LTRIM(key string, start interface{}, end interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redigo.Strings(conn.Do("LTRIM", key, start, end))
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) LPOP(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	val, err := redigo.Strings(conn.Do("LPOP", key))
	if err != nil {
		return "", err
	}
	if len(val) == 0 {
		return "", nil
	}
	return val[0], nil
}

func (r *redis) SADD(key string, val ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	var args []interface{}
	args = append(args, key)
	args = append(args, val...)

	_, err := redigo.Int64(conn.Do("SADD", args...))
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) SREM(key string, val ...interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()

	var args []interface{}
	args = append(args, key)
	args = append(args, val...)

	_, err := redigo.Int64(conn.Do("SREM", args...))
	if err != nil {
		return err
	}

	return nil
}

func (r *redis) SMEMBERS(key string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	resp, err := redigo.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *redis) SCARD(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	resp, err := redigo.Int64(conn.Do("SCARD", key))
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *redis) SISMEMBERS(key string, field interface{}) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	resp, err := redigo.Bool(conn.Do("SISMEMBER", key, field))
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (r *redis) RPUSH(key string, values ...string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	args := []interface{}{key}
	for _, v := range values {
		args = append(args, v)
	}

	return redigo.Int64(conn.Do("RPUSH", args...))
}
