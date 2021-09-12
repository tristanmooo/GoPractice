package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	errors "github.com/pkg/errors"
)

type RedisConn struct {
	host string
	port int
	pswd string
	conn *redis.Conn
}

func (redisConn *RedisConn) init(host string, port int, pswd string) {
	redisConn.host = host
	redisConn.port = port
	redisConn.pswd = pswd
	redisConn.conn = new(redis.Conn)
	fmt.Printf("Redis conn init succ, host: %s, port: %d\n", redisConn.host, redisConn.port)
}

func (redisConn *RedisConn) connect() (c redis.Conn, err error) {
	connString := fmt.Sprintf("%s:%d", redisConn.host, redisConn.port)
	c, err = redis.Dial("tcp", connString)
	if err != nil {
		// fmt.Println(err)
		return *new(redis.Conn), errors.Wrap(err, "Fail to establish redis connection.")
	}
	redis.DialPassword(redisConn.pswd)
	fmt.Printf("Connection to redis established successfully. \n")
	return c, nil
}

func (redisConn *RedisConn) set(key interface{}, value interface{}) (res interface{}, err error) {
	res, err = (*redisConn.conn).Do("set", key, value)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to set value.")
	}
	fmt.Printf("Operation: set, key: %s, value: %v\n", key, value)
	return res, err
}

func (redisConn *RedisConn) get(key interface{}) (res interface{}, err error) {
	res, err = redis.String((*redisConn.conn).Do("get", key))
	if err != nil {
		return nil, errors.Wrap(err, "Fail to get result.")
	}
	fmt.Printf("Operation: get, key: %s, result: %v\n", key, res)
	return res, err
}

func logging(err error) {
	fmt.Printf("original error: %T, %v\n", errors.Cause(err), errors.Cause(err))
	fmt.Printf("statck trace: \n%+v\n", err)
}

func main() {

	redisConn := new(RedisConn)
	redisConn.init("192.168.66.172", 6379, "mypass")

	var err error

	*redisConn.conn, err = redisConn.connect()
	if err != nil {
		logging(err)
		return
	}

	_, err = redisConn.set("abc", 200)
	if err != nil {
		logging(err)
		return
	}

	_, err = redisConn.get("abc")
	if err != nil {
		logging(err)
		return
	}

}
