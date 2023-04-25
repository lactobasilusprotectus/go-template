package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"log"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	run, err := miniredis.Run()

	if err != nil {
		log.Fatalf("miniredis.Run returns err: %+v\n", err)
	}

	defer run.Close()

	NewRedisClient(config.RedisConfig{
		Host:     run.Addr(),
		Password: "",
		DB:       0,
	})

}
