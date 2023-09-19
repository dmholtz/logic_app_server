SELECT id,
    type,
    question,
    answers,
    solutions
FROM quiz
WHERE is_competition_mode = 0
    AND type = ?
    AND difficulty = ?
    AND num_vars = ?
    AND id NOT IN (
        SELECT quiz_id
        FROM quiz_participation
        WHERE user_id = ?
    )