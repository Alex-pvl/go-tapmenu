box.watch('box.status', function()
    if box.info.ro then
        return
    end

    box.schema.user.create('admin', {
        password = 'admin',
        if_not_exists = true
    })
    box.schema.user.grant('admin', 'super')

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
    box.space.tables:insert { 'hash1', 'menu.ru', 'rest_name', 1, datetime }
    box.space.tables:insert { 'hash2', 'menu.ru', 'rest_name', 2, datetime }
    box.space.tables:insert { 'hash3', 'menu.ru', 'rest_name', 3, datetime }
    box.space.tables:insert { 'hash4', 'menu.ru', 'rest_name', 4, datetime }
    box.space.tables:insert { 'hash5', 'menu.ru', 'rest_name', 5, datetime }

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
        { name = 'username',        type = 'string' },
        { name = 'hashed_password', type = 'string' },
        { name = 'session_token',   type = 'string' },
        { name = 'csrf_token',      type = 'string' },
    })
    box.space.waiters:create_index('pk', { parts = { 'username' }, if_not_exists = true })
end)
