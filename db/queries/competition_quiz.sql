-- Get one quiz in competition mode that the user has not yet participated in
SELECT id,
    type,
    time_limit,
    question,
    possible_answers
FROM quiz
WHERE is_competition_mode = 1
    AND id NOT IN (
        SELECT quiz_id
        FROM quiz_participation
        WHERE user_id = ?
    )
ORDER BY id ASC
LIMIT 1;