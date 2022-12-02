package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	_, err = rdb.Ping(context.Background()).Result()
	return
}

func redisExample() {
	_, err := rdb.Set(context.Background(), "score", 100, 0).Result()
	if err != nil {
		fmt.Printf("set score failed, err:%v\n", err)
		return
	}

	val, err := rdb.Get(context.Background(), "score").Result()
	if err != nil {
		fmt.Printf("get score failed, err:%v\n", err)
		return
	}
	fmt.Println("score", val)

	val2, err := rdb.Get(context.Background(), "name").Result()
	if err != nil {
		fmt.Printf("get name failed, err:%v\n", err)
		return
	}
	fmt.Println("name", val2)
}

func hgetExample() {
	v, err := rdb.HGetAll(context.Background(), "user").Result()
	if err != nil {
		fmt.Printf("hgetall failed. err:%v\n", err)
		return
	}
	fmt.Println(v)

	v2 := rdb.HMGet(context.Background(), "user", "name", "age").Val()
	fmt.Println(v2)

	v3 := rdb.HGet(context.Background(), "user", "age").Val()
	fmt.Println(v3)
}

func zsetExample() {
	zsetKey := "launguage_rank"
	languages := []redis.Z{
		redis.Z{Score: 90.0, Member: "Golang"},
		redis.Z{Score: 98.0, Member: "Java"},
		redis.Z{Score: 95.0, Member: "Python"},
		redis.Z{Score: 97.0, Member: "JavaScrit"},
		redis.Z{Score: 99.0, Member: "C/C++"},
	}
	num, err := rdb.ZAdd(context.Background(), zsetKey, languages...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Printf("zadd %d success.\n", num)

	newScore, err := rdb.ZIncrBy(context.Background(), zsetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 取分数最高的三个
	ret, err := rdb.ZRevRangeWithScores(context.Background(), zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed. err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取95~100分的
	op := redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(context.Background(), zsetKey, &op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}

	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

func pipelineExample() {
	pipe := rdb.Pipeline()
	incr := pipe.Incr(context.Background(), "pipeline_counter")
	pipe.Expire(context.Background(), "pipeline_counter", time.Hour)

	_, err := pipe.Exec(context.Background())
	fmt.Println(incr.Val(), err)
}

func txPipelineExample() {
	pipe := rdb.TxPipeline()
	incr := pipe.Incr(context.Background(), "pipeline_counter")
	pipe.Expire(context.Background(), "pipeline_counter", time.Hour)

	_, err := pipe.Exec(context.Background())
	fmt.Println(incr.Val(), err)
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed, err:%v\n", err)
	}
	fmt.Println("connect redis success...")
	defer rdb.Close()

	// redisExample()
	// hgetExample()
	// zsetExample()
	// pipelineExample()
	txPipelineExample()
}
