package models

//
//var (
//	conn *redis.Client
//	Ctx  = context.Background()
//)
//
//func Redis() *redis.Client {
//	return conn
//}
//
//func InitRedis() *redis.Client {
//	addr := v.GetString("db.redis.addr")
//	pwd := v.GetString("db.redis.pwd")
//	dbn := v.GetInt("db.redis.db")
//	poolSize := v.GetInt("db.redis.poolSize")
//
//	conn = redis.NewClient(&redis.Options{
//		Addr:     addr,     // host:port
//		Password: pwd,      // set password
//		DB:       dbn,      // use default DB
//		PoolSize: poolSize, // 连接池大小
//	})
//
//	//检测心跳
//	_, err := conn.Ping(Ctx).Result()
//	if err != nil {
//		log.Fatalln("redis connection failed:", err)
//	}
//	return conn
//}
