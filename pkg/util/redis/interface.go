package redis

type Interface interface {
	Get(key string) (reply interface{}, err error)
	Set(key string, value interface{}, expireSeconds int) (err error)
}
