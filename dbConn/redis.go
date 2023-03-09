package dbConn

//var ctx = context.TODO()
//
//func Redis() *redis.Client {
//	return initRedis()
//}
//
//func initRedis() *redis.Client {
//	addr := v.GetString("db.redis.addr")
//	pwd := v.GetString("db.redis.pwd")
//	dbn := v.GetInt("db.redis.db")
//	poolSize := v.GetInt("db.redis.poolSize")
//
//	client := redis.NewClient(&redis.Options{
//		Addr:     addr,     // host:port
//		Password: pwd,      // set password
//		DB:       dbn,      // use default DB
//		PoolSize: poolSize, // 连接池大小
//	})
//
//	//检测心跳
//	_, err := client.Ping(ctx).Result()
//	if err != nil {
//		log.Fatalln("redis connection failed:", err)
//	}
//	return client
//}
