# rest-api-service

# REST API

GET /api/items - list of items +
GET /api/items/{item_id} - item +
DELETE /api/items/{item_id} +
POST /api/items +
GET /api/items/{item_id}/trades
GET /api/trades
POST /api/trades
DELETE /api/trades/{trade_id}
GET /api/trades/{trade_id}
PUT /api/trades/{trade_id}
GET /api/users/{user_id}/trades

GET /api/users -- 200, 404, 500
POST /api/users/{user_id} -- 204, 4xx, Header Location: url
DELETE /api/users/{user_id} -- 204, 404, 400
GET /api/users/{user_id} -- 200, 404, 500
PUT /api/users/{user_id} -- 204/200
PATCH /api/users/{user_id} -- 204/200, 404, 400, 500 