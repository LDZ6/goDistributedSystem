# Redis 数据类型及命令详解

## 数据类型及命令
### 通用命令
- `help @generic`：查看通用命令
- `help keys`：查看 `keys` 命令的详细用法
- `KEYS *`：查看所有 key
- `DEL key`：删除 key（多个 key 用空格隔开）
- `TYPE key`：查看 key 的类型
- `EXISTS key`：判断 key 是否存在
- `EXPIRE key seconds`：设置 key 的过期时间
- `TTL key`：查看 key 的剩余过期时间

### String
- `int`：整型
- `string`：字符串型
- `float`：浮点型

#### String 命令
- `SET key value`：设置 key 的值
- `GET key`：获取 key 的值
- `MSET key1 value1 key2 value2`：设置多个 key 的值
- `MGET key1 key2`：获取多个 key 的值
- `APPEND key value`：追加值到 key 的值末尾
- `STRLEN key`：获取 key 值的长度
- `INCR key`：自增 key 的值
- `DECR key`：自减 key 的值
- `INCRBY key increment`：自增 key 的值
- `INCRBYFLOAT key increment`：自增 key 的值
- `SETNX key value`：设置 key 的值，但只有在 key 不存在时才设置成功
- `SETEX key seconds value`：设置 key 的值并设置过期时间
- `DECRBY key decrement`：自减 key 的值
- `GETRANGE key start end`：获取 key 值的子串
- `SETRANGE key offset value`：设置 key 值的子串
- `GETSET key value`：设置 key 的值并返回旧值

### Hash
#### Hash 命令
- `HSET key field value`：设置 hash 表的 field 的值
- `HGET key field`：获取 hash 表的 field 的值
- `HMSET key field1 value1 field2 value2`：设置多个 hash 表的 field 的值
- `HMGET key field1 field2`：获取多个 hash 表的 field 的值
- `HGETALL key`：获取 hash 表的所有 field 和 value
- `HKEYS key`：获取 hash 表的所有 field
- `HVALS key`：获取 hash 表的所有 value
- `HLEN key`：获取 hash 表的 field 数量
- `HINCRBY key field increment`：自增 hash 表的 field 的值
- `HINCRBYFLOAT key field increment`：自增 hash 表的 field 的值
- `HSETNX key field value`：设置 hash 表的 field 的值，但只有在 field 不存在时才设置成功
- `HDEL key field1 field2`：删除 hash 表的 field
- `HEXISTS key field`：判断 hash 表的 field 是否存在
- `HSCAN key cursor [MATCH pattern] [COUNT count]`：迭代 hash 表的 field 和 value

### List
#### List 命令
- `LPUSH key value1 value2`：向列表左侧添加元素
- `RPUSH key value1 value2`：向列表右侧添加元素
- `LPOP key`：弹出列表左侧的元素
- `RPOP key`：弹出列表右侧的元素
- `LINDEX key index`：获取列表指定索引的元素
- `LLEN key`：获取列表的长度
- `LRANGE key start end`：获取列表指定范围的元素
- `LREM key count value`：删除列表指定元素
- `LSET key index value`：设置列表指定索引的元素的值
- `LTRIM key start end`：截取列表指定范围的元素
- `LINSERT key BEFORE|AFTER pivot value`：在列表指定元素前或后插入元素
- `BLPOP key1 key2 timeout`：弹出列表左侧第一个非空列表的元素，如果列表为空则阻塞

### Set
#### Set 命令
- `SADD key member1 member2`：向集合添加元素
- `SPOP key`：随机弹出集合的一个元素
- `SREM key member1 member2`：删除集合的元素
- `SCARD key`：获取集合的元素数量
- `SISMEMBER key member`：判断元素是否在集合中
- `SMEMBERS key`：获取集合的所有元素
- `SRANDMEMBER key [count]`：随机获取集合的元素
- `SINTER key1 key2`：交集
- `SINTERSTORE destination key1 key2`：交集并存储到 destination 集合
- `SUNION key1 key2`：并集
- `SUNIONSTORE destination key1 key2`：并集并存储到 destination 集合
- `SDIFF key1 key2`：差集
- `SDIFFSTORE destination key1 key2`：差集并存储到 destination 集合
- `SMOVE source destination member`：移动元素
- `SPOP key`：随机弹出集合的一个元素

### SortedSet
#### SortedSet 命令
- `ZADD key score1 member1 score2 member2`：向有序集合添加元素
- `ZPOPMIN key [count]`：弹出有序集合中分数最小的元素
- `ZPOPMAX key [count]`：弹出有序集合中分数最大的元素
- `ZRANGE key start end [WITHSCORES]`：获取有序集合指定范围的元素
- `ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]`：获取有序集合指定分数范围的元素
- `ZRANK key member`：获取元素在有序集合中的排名
- `ZREVRANK key member`：获取元素在有序集合中的倒排名
- `ZREM key member1 member2`：删除有序集合的元素
- `ZREMRANGEBYRANK key start end`：删除有序集合指定排名范围的元素
- `ZREMRANGEBYSCORE key min max`：删除有序集合指定分数范围的元素
- `ZCARD key`：获取有序集合的元素数量
- `ZCOUNT key min max`：获取有序集合指定分数范围的元素数量