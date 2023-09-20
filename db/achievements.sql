INSERT INTO achievement (name, description, level, sql)
VALUES (
        "First Steps",
        "Solve your 5 first quizzes",
        "bronze",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 GROUP BY user_id HAVING cnt >= 5"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Quiz Master",
        "Complete 10 quizzes",
        "silver",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 GROUP BY user_id HAVING cnt >= 10"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Quiz Champion",
        "Complete 20 quizzes",
        "gold",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 GROUP BY user_id HAVING cnt >= 20"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Errare Humanum Est",
        "Get 5 questions wrong",
        "bronze",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 0 GROUP BY user_id HAVING cnt >= 5"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Can't wait forever",
        "Solve a quiz in less than 60 seconds",
        "bronze",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 AND time <= 60 GROUP BY user_id HAVING cnt >= 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Better Be Fast",
        "Solve a quiz in less than 20 seconds",
        "silver",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 AND time <= 20 GROUP BY user_id HAVING cnt >= 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Speed Demon",
        "Solve a quiz in less than 10 seconds",
        "gold",
        "SELECT COUNT(id) as cnt FROM quiz_participation WHERE user_id = ? AND correct = 1 AND time <= 10 GROUP BY user_id HAVING cnt >= 1"
    );