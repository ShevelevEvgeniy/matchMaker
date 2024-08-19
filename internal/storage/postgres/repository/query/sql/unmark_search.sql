UPDATE users SET search_match = false, search_start_time = NULL
WHERE id = ANY($1);