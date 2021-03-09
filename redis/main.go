// docker exec -it a6e9c0148dd9 redis-cli

/*
SET
Команда используется для установки ключа и его значения,
с дополнительными необязательными параметрами для указания срока действия записи значения ключа.
Параметр EX указывает время жизни объекта в секундах, PX в милисекундах

127.0.0.1:6379> SET foo "hello world"
OK
127.0.0.1:6379> SET foo1 "hello world" ex 5
OK

GET
Команда используется для получения значения, связанного с ключом.
Если запись значения ключа превысила срок действия, будет возвращено nil:

127.0.0.1:6379> GET foo
"hello world"
# если истечет время жизни записи
127.0.0.1:6379> GET foo
(nil)
По умолчанию все значение в Redis сохраняются как строки.

EXISTS
Эта команда проверяет, существует ли что то с данным ключом.
Она возвращает 1 если объект существует или 0 если нет. Boolean типа в Redis нет.

*/

package main

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val, err = rdb.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
}
