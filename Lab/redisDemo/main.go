package redisDemo

//******连接
//******基本使用
//******其他数据结构
//******pipeline
//******事务

//******连接*********************************************
//普通连接模式：连接到由 Redis Sentinel 管理的 Redis 服务器。
//rdb := redis.NewClient(&redis.Options{
//	Addr:     "localhost:6379",
//	Password: "", // 密码
//	DB:       0,  // 数据库
//	PoolSize: 20, // 连接池大小
//})
////TLS连接模式
//rdb := redis.NewClient(&redis.Options{
//	TLSConfig: &tls.Config{
//	MinVersion: tls.VersionTLS12,
//	},
//})
////redis sentinel模式：连接到由 Redis Sentinel 管理的 Redis 服务器。
//rdb := redis.NewFailoverClient(&redis.FailoverOptions{
//	MasterName:    "master-name",
//	SentinelAddrs: []string{":9126", ":9127", ":9128"},
//})
////redis cluster模式：go-redis 支持按延迟或随机路由命令。
//rdb := redis.NewClusterClient(&redis.ClusterOptions{
//	Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
//	// 若要根据延迟或随机路由命令，请启用以下命令之一
//	// RouteByLatency: true,
//	// RouteRandomly: true,
//})

//******基本使用*******************************************
//执行命令
//func doCommand() {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	// 执行命令获取结果
//	val, err := rdb.Get(ctx, "key").Result()
//	fmt.Println(val, err)
//
//	// 先获取到命令对象
//	cmder := rdb.Get(ctx, "key")
//	fmt.Println(cmder.Val()) // 获取值
//	fmt.Println(cmder.Err()) // 获取错误
//
//	// 直接执行命令获取错误
//	err = rdb.Set(ctx, "key", 10, time.Hour).Err()
//
//	// 直接执行命令获取值
//	value := rdb.Get(ctx, "key").Val()
//	fmt.Println(value)
//}
//
////执行任意命令
//// doDemo rdb.Do 方法使用示例
//func doDemo() {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	// 直接执行命令获取错误
//	err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
//	fmt.Println(err)
//
//	// 执行命令获取结果
//	val, err := rdb.Do(ctx, "get", "key").Result()
//	fmt.Println(val, err)
//}
//
////redis.nil:redis.Nil 错误来表示 Key 不存在的错误
//func getValueFromRedis(key, defaultValue string) (string, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	val, err := rdb.Get(ctx, key).Result()
//	if err != nil {
//		// 如果返回的错误是key不存在
//		if errors.Is(err, redis.Nil) {
//			return defaultValue, nil
//		}
//		// 出其他错了
//		return "", err
//	}
//	return val, nil
//}

//**********其他数据结构*************************************
//zset示例
//func zsetDemo() {
//	// key
//	zsetKey := "language_rank"
//	// value
//	languages := []*redis.Z{
//		{Score: 90.0, Member: "Golang"},
//		{Score: 98.0, Member: "Java"},
//		{Score: 95.0, Member: "Python"},
//		{Score: 97.0, Member: "JavaScript"},
//		{Score: 99.0, Member: "C/C++"},
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	// ZADD
//	err := rdb.ZAdd(ctx, zsetKey, languages...).Err()
//	if err != nil {
//		fmt.Printf("zadd failed, err:%v\n", err)
//		return
//	}
//	fmt.Println("zadd success")
//
//	// 把Golang的分数加10
//	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
//	if err != nil {
//		fmt.Printf("zincrby failed, err:%v\n", err)
//		return
//	}
//	fmt.Printf("Golang's score is %f now.\n", newScore)
//
//	// 取分数最高的3个
//	ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
//	for _, z := range ret {
//		fmt.Println(z.Member, z.Score)
//	}
//
//	// 取95~100分的
//	op := &redis.ZRangeBy{
//		Min: "95",
//		Max: "100",
//	}
//	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
//	if err != nil {
//		fmt.Printf("zrangebyscore failed, err:%v\n", err)
//		return
//	}
//	for _, z := range ret {
//		fmt.Println(z.Member, z.Score)
//	}
//}
//
////遍历所有key:使用KEYS prefix:* 命令按前缀获取所有 key。
//vals, err := rdb.Keys(ctx, "prefix*").Result()
//// scanKeysDemo1 按前缀查找所有key示例
//func scanKeysDemo1() {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	var cursor uint64
//	for {
//		var keys []string
//		var err error
//		// 按前缀扫描key
//		keys, cursor, err = rdb.Scan(ctx, cursor, "prefix:*", 0).Result()
//		if err != nil {
//			panic(err)
//		}
//
//		for _, key := range keys {
//			fmt.Println("key", key)
//		}
//
//		if cursor == 0 { // no more keys
//			break
//		}
//	}
//}
//// scanKeysDemo2 按前缀扫描key示例
//func scanKeysDemo2() {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//	// 按前缀扫描key
//	iter := rdb.Scan(ctx, 0, "prefix:*", 0).Iterator()
//	for iter.Next(ctx) {
//		fmt.Println("keys", iter.Val())
//	}
//	if err := iter.Err(); err != nil {
//		panic(err)
//	}
//}

//*********************pipeline********************************************
//Redis Pipeline 允许通过使用单个 client-server-client 往返执行多个命令来提高性能。区别于一个接一个地执行100个命令，
//你可以将这些命令放入 pipeline 中，然后使用1次读写操作像执行单个命令一样执行它们。这样做的好处是节省了执行命令的网络往返时间（RTT）
//一次性执行多个命令场景下，用pipeline优化！！！！！
//pipe := rdb.Pipeline()
//
//incr := pipe.Incr(ctx, "pipeline_counter") //命令一
//pipe.Expire(ctx, "pipeline_counter", time.Hour)  //命令二
//
//cmds, err := pipe.Exec(ctx) //一次性执行
//if err != nil {
//panic(err)
//}
//
//// 在执行pipe.Exec之后才能获取到结果
//fmt.Println(incr.Val())

//***********事务************************************
//Redis 是单线程执行命令的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，例如在它们之间交替执行。
//Multi/exec能够确保在multi/exec两个语句之间的命令之间没有其他客户端正在执行命令。
//其实和pipeline原理一模一样：两种使用方式
// TxPipeline demo
//pipe := rdb.TxPipeline()
//incr := pipe.Incr(ctx, "tx_pipeline_counter")
//pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
//_, err := pipe.Exec(ctx)
//fmt.Println(incr.Val(), err)
//
//// TxPipelined demo
//var incr2 *redis.IntCmd
//_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
//	incr2 = pipe.Incr(ctx, "tx_pipeline_counter")
//	pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
//	return nil
//})
//fmt.Println(incr2.Val(), err)
//上述命令实际上就是以下
//MULTI
//INCR pipeline_counter
//EXPIRE pipeline_counts 3600
//EXEC

//watch:搭配 WATCH命令来执行事务操作。从使用WATCH命令监视某个 key 开始，直到执行EXEC命令的这段时间里，
//如果有其他用户抢先对被监视的 key 进行了替换、更新、删除等操作，那么当用户尝试执行EXEC的时候，
//事务将失败并返回一个错误，用户可以根据这个错误选择重试事务或者放弃事务。
// watchDemo 在key值不变的情况下将其值+1

//func watchDemo(ctx context.Context, key string) error {
//	return rdb.Watch(ctx, func(tx *redis.Tx) error {
//		n, err := tx.Get(ctx, key).Int()
//		if err != nil && err != redis.Nil {
//			return err
//		}
//		// 假设操作耗时5秒
//		// 5秒内我们通过其他的客户端修改key，当前事务就会失败
//		time.Sleep(5 * time.Second)
//		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
//			pipe.Set(ctx, key, n+1, time.Hour)
//			return nil
//		})
//		return err
//	}, key)
//}
//将上面的函数执行并打印其返回值，如果我们在程序运行后的5秒内修改了被 watch 的 key 的值，
//那么该事务操作失败，返回redis: transaction failed错误。
