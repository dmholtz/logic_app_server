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
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Rare Competitor",
        "Solve a quiz in competition mode",
        "bronze",
        "SELECT COUNT(qp.id) as cnt FROM quiz_participation qp INNER JOIN quiz q ON qp.quiz_id = q.id WHERE user_id = ? AND correct = 1 AND is_competition_mode = 1 GROUP BY user_id HAVING cnt >= 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Scorer",
        "Collect 20 points in competition mode",
        "bronze",
        "SELECT SUM(qp.points) as score FROM quiz_participation qp WHERE user_id = ? GROUP BY user_id HAVING score >= 20"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Gladiator",
        "Collect 100 points in competition mode",
        "bronze",
        "SELECT SUM(qp.points) as score FROM quiz_participation qp WHERE user_id = ? GROUP BY user_id HAVING score >= 100"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "All the same",
        "Solve a hard EQUIV quiz",
        "gold",
        "SELECT COUNT(qp.id) as cnt FROM quiz_participation qp INNER JOIN quiz q ON qp.quiz_id = q.id WHERE user_id = ? AND correct = 1 AND difficulty = 'hard' AND type = 'EQUIV' GROUP BY user_id HAVING cnt >= 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Fortuneteller",
        "Solve 7 TAUT quizzes of any difficulty",
        "silver",
        "SELECT COUNT(qp.id) as cnt FROM quiz_participation qp INNER JOIN quiz q ON qp.quiz_id = q.id WHERE user_id = ? AND correct = 1 AND type = 'TAUT' GROUP BY user_id HAVING cnt >= 7"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "Second Chance",
        "Correct a mistake without looking at the solution",
        "bronze",
        "SELECT COUNT(qp1.id) as cnt FROM quiz_participation qp1 INNER JOIN quiz_participation qp2 ON qp1.quiz_id = qp2.quiz_id WHERE qp1.user_id = ? AND qp2.user_id = qp1.user_id AND qp1.correct = 1 AND qp2.correct = 0 GROUP BY qp1.user_id HAVING cnt >= 1"
    );
INSERT INTO achievement (name, description, level, sql)
VALUES (
        "h-Index = 4",
        "Solve at least 4 quizzes with at least 4 variables",
        "gold",
        "SELECT COUNT(qp.id) as cnt FROM quiz_participation qp INNER JOIN quiz q ON qp.quiz_id = q.id WHERE user_id = ? AND correct = 1 AND num_vars >= 4 GROUP BY user_id HAVING cnt >= 4"
    );