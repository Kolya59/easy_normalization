package redisdriver

import (
	"time"

	"github.com/go-redis/cache"
	"github.com/psu/easy_normalization/pkg/car"
)

// Send info to Redis database
func SetCar(newCar *car.Car, codec *cache.Codec, index string) error {
	err = codec.Set(&cache.Item{
		Ctx:        nil,
		Key:        index,
		Object:     newCar,
		Func:       nil,
		Expiration: time.Minute,
	})

	if err != nil {
		return err
	}
	return nil
}

// Get info from Redis data base
func GetCar(codec *cache.Codec, index string) (newCar *car.Car, err error) {
	newCar = &car.Car{}
	err = codec.Get(index, newCar)
	if err != nil {
		return nil, err
	}
	return newCar, nil
}
