SELECT type,
    coalesce(COUNT(DISTINCT q.id), 0),
    coalesce(COUNT(qp.id), 0),
    coalesce(AVG(qp.time), 0)
FROM quiz q
    LEFT JOIN quiz_participation qp ON q.id = qp.quiz_id
GROUP BY type