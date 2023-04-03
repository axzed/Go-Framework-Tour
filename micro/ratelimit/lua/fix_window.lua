local val = redis.call('get', KEYS[1])
local limit = tonumber(ARGV[2])
if val == false then
    if limit < 1 then
        -- 执行限流
        return "true"
    else
        -- key 不存在，设置初始值 1，并且设置过期时间
        redis.call('set', KEYS[1], 1, 'PX', ARGV[1])
        -- 不执行限流
        return "false"
    end
elseif tonumber(val) < limit then
    -- 自增 1
    redis.call('incr', KEYS[1])
    -- 不需要限流
    return "false"
else
    -- 限流
    return "true"
end