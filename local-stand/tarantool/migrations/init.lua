box.watch('box.status', function()
    if box.info.ro then
        return
    end

    box.schema.user.create('api_user', {
        password = 'api_user',
        if_not_exists = true
    })
    box.schema.user.grant('api_user', 'super')

    box.schema.space.create('tables', { if_not_exists = true })
    box.space.tables:format({
        { name = 'hash',            type = 'string' },
        { name = 'menu_url',        type = 'string' },
        { name = 'restaurant_name', type = 'string' },
        { name = 'table_number',    type = 'number' },
        { name = 'last_call',       type = 'datetime' },
    })
    box.space.tables:create_index('pk', { parts = { 'hash' }, if_not_exists = true })

    local datetime = require("datetime").parse("2025-01-01T00:00:00.000000000-00:00")
    for i = 1, 20 do
        if i > 10 then
            box.space.tables:insert { string.format("hash%d", i), 'menu.sabaibar.ru', 'SabaiBar', i, datetime }
        else
            box.space.tables:insert { string.format("hash%d", i), 'menu.ru', 'rest_name', i, datetime }
        end
    end

    box.schema.space.create('orders', { if_not_exists = true })
    box.space.orders:format({
        { name = 'id',              type = 'string' },
        { name = 'restaurant_name', type = 'string' },
        { name = 'table_number',    type = 'number' },
        { name = 'created_at',      type = 'datetime' },
        { name = 'updated_at',      type = 'datetime' },
        { name = 'accepted',        type = 'boolean' },
    })
    box.space.orders:create_index('pk', { parts = { 'id' }, if_not_exists = true })

    box.schema.space.create('waiters', { if_not_exists = true })
    box.space.waiters:format({
        { name = 'waiter_id',       type = 'string' },
        { name = 'username',        type = 'string' },
        { name = 'hashed_password', type = 'string' },
        { name = 'restaurant_name', type = 'string' },
    })
    box.space.waiters:create_index('pk', { parts = { 'waiter_id' }, if_not_exists = true })
    box.space.waiters:create_index('waiter_username_idx', { parts = { 'username' }, if_not_exists = true })
    box.space.waiters:insert { 'dcd4ca2b-60c6-43a4-aa43-3dbfdbf73982', 'admin',
        '$2a$10$V16aIiBMif9aYH0FUOOjG.RsgSU9jBhGNOpfdr5ChnqHpSLF2e3YG', 'rest_name' }
end)
