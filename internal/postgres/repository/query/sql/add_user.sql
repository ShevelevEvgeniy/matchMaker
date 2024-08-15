INSERT INTO users (name, skill, latency, searching_match, search_start_time)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name)
DO UPDATE SET
   searching_match = EXCLUDED.searching_match,
   search_start_time = EXCLUDED.search_start_time;