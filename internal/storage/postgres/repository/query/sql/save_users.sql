INSERT INTO users (name, skill, latency, search_match, search_start_time)
VALUES %s
ON CONFLICT (name)
DO UPDATE SET
    skill = EXCLUDED.skill,
    latency = EXCLUDED.latency,
    search_match = EXCLUDED.search_match,
    search_start_time = EXCLUDED.search_start_time;