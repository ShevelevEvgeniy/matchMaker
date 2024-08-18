SELECT id, name, skill, latency, search_start_time
FROM users
WHERE search_match = true
ORDER BY search_start_time ASC
LIMIT $1;