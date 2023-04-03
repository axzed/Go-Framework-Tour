local val = redis.call("get", KEYS[1])
if not val then
    --    key 不存在
    return redis.call('set', KEYS[1], ARGV[1], 'PX', ARGV[2])
elseif val == ARGV[1] then
-- 你上一次加锁成功
    redis.call('expire', KEYS[1], ARGV[2])
    return  "OK"
else
    -- 被人持有锁
    return ""
end
