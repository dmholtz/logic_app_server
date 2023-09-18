WITH t as (
    SELECT a0.id as id,
        name,
        description,
        level,
        1 as isAchieved
    FROM achievement a0
        INNER JOIN achieved a1 ON a0.id = a1.achievement_id
    WHERE a1.user_id = ?
)
SELECT name,
    description,
    level,
    isAchieved
FROM t
UNION
SELECT name,
    description,
    level,
    0 as isAchieved
FROM achievement
WHERE id not in (
        SELECT id
        FROM t
    )