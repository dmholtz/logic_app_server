with xp_table as (
    SELECT username,
        coalesce(COUNT(a.id), 0) as xp
    FROM users
        LEFT JOIN achieved a ON users.id = a.user_id
    GROUP BY username
),
score_table as (
    SELECT username,
        coalesce(SUM(qp.points), 0) as score
    FROM users
        LEFT JOIN quiz_participation qp ON users.id = qp.user_id
    GROUP BY username
)
SELECT xt.username,
    xt.xp as xp,
    st.score as score
FROM xp_table xt
    INNER JOIN score_table st ON st.username = xt.username
ORDER BY st.score DESC,
    xt.xp DESC,
    xt.username ASC;