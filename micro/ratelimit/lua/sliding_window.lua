-- KEY 只有一个
-- 参数（按照顺序）：阈值，窗口大小（毫秒），当前时间戳（毫秒）。
-- 你们也可以考虑使用秒或者纳秒作为单位，差异不大
local threshold = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local min = now - window
local key = KEYS[1]

-- 先挪动窗口
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

-- 看看还有没有容量
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')

if cnt >= threshold then
    -- 限流
    return "true"
else
    -- 值和优先级我们都设置成当前时间戳
    redis.call('ZADD', key, now, now)
    -- 这里设不设置过期时间影响不大，设置了过期时间可以防止长期没有人访问的 key 正常被删除
    redis.call('PEXPIRE', key, window)
    -- 不限流
    return "false"
end

