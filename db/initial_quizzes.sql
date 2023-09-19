INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        possible_answers
    )
VALUES (
        "SAT",
        "easy",
        2,
        60,
        0,
        "A | (B & !B)",
        '[{"answer": "yes", "solution": true}, {"answer": "no", "solution": false}]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        possible_answers
    )
VALUES (
        "SAT",
        "hard",
        2,
        NULL,
        0,
        "(A <-> !B) & (!A -> !B)",
        '[{"answer": "yes", "solution": true}, {"answer": "no", "solution": false}]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        possible_answers
    )
VALUES (
        "EQUIV",
        "hard",
        2,
        30,
        1,
        "A <-> B",
        '[{"answer": "(A -> B) & (!A -> !B)", "solution": true}, {"answer": "(A & B) | (!A & !B) ", "solution": true}, {"answer": "(A | B) & (A <-> B)", "solution": false}, {"answer": "(A -> B) & (!B | A)", "solution": true}]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        possible_answers
    )
VALUES (
        "TAUT",
        "easy",
        2,
        NULL,
        0,
        "(A | !A) & (B | !B)",
        '[{"answer": "yes", "solution": true}, {"answer": "no", "solution": false}]'
    )