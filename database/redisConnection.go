package database

import (
	"auction/secret"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RCredentials struct {
	Addr     string
	Password string
	DB       int
}

type RData struct {
	redisCredentials RCredentials
	Client           *redis.Client
}

var (
	Rcn *RData
)

func (rdb *RData) NewRedis() error {
	rdb = &RData{}
	secretRedis := secret.Env["Redis"].(map[string]interface{})
	rdb.redisCredentials = RCredentials{
		Addr:     secretRedis["Addr"].(string) + ":" + secretRedis["Port"].(string),
		Password: secretRedis["Password"].(string),
		DB:       int(secretRedis["DB"].(float64)),
	}

	client := redis.NewClient(&redis.Options{
		Addr:     rdb.redisCredentials.Addr,
		Password: rdb.redisCredentials.Password,
		DB:       rdb.redisCredentials.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return err
	}
	rdb.Client = client
	Rcn = rdb
	return nil
}

func GetName[T any](ctx context.Context, nconst string, setModel *T) error {
	cmd := Rcn.Client.Get(ctx, nconst)
	cmdb, err := cmd.Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			return errors.New("nil")
		}
		return err
	}
	err = json.Unmarshal(cmdb, &setModel)
	if err != nil {
		return err
	}
	return nil
}

func SetName[T any](ctx context.Context, value *T, id string, time time.Duration) error {

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Rcn.Client.Set(ctx, id, val, time).Err()
}
