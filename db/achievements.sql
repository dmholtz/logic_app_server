INSERT INTO achievement (name, description, level, sql)
VALUES (
        "First Steps",
        "Solve your 5 first quizzes",
        "bronze",
        "SELECT COUNT(*) as cnt FROM quiz_participation WHERE user_id = ? AND cnt >= 5 AND correct = 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Quiz Master",
        "Complete 10 quizzes",
        "silver",
        "SELECT COUNT(*) as cnt FROM quiz_participation WHERE user_id = ? AND cnt >= 10 AND correct = 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Quiz Champion",
        "Complete 20 quizzes",
        "gold",
        "SELECT COUNT(*) as cnt FROM quiz_participation WHERE user_id = ? AND cnt >= 20 AND correct = 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Errare Humanum Est",
        "Get 5 questions wrong",
        "bronze",
        "SELECT COUNT(*) as cnt FROM quiz_participation WHERE user_id = ? AND cnt >= 5 AND correct = 0"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Better Be Fast",
        "Solve a quiz in less than 20 seconds",
        "silver",
        "SELECT COUNT(*) as cnt FROM quiz_participation qp, quiz q WHERE qp.quiz_id = q.id AND qp.user_id = ? AND q.time_limit IS NOT NULL AND q.time_limit <= 20 AND qp.correct = 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Speed Demon",
        "Solve a quiz in less than 10 seconds",
        "gold",
        "SELECT COUNT(qp.id) as cnt FROM quiz_participation qp, quiz q WHERE qp.quiz_id = q.id AND qp.user_id = ? AND q.time_limit IS NOT NULL AND q.time_limit <= 10 AND qp.correct = 1"
    );