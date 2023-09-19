INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        answers,
        solutions
    )
VALUES (
        "SAT",
        "easy",
        2,
        60,
        0,
        "A | (B & !B)",
        '["yes", "no"]',
        '[true, false]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        answers,
        solutions
    )
VALUES (
        "SAT",
        "hard",
        2,
        NULL,
        0,
        "(A <-> !B) & (!A -> !B)",
        '["yes", "no"]',
        '[true, false]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        answers,
        solutions
    )
VALUES (
        "EQUIV",
        "hard",
        2,
        30,
        1,
        "A <-> B",
        '["(A -> B) & (!A -> !B)", "(A & B) | (!A & !B)", "(A | B) & (A <-> B)", "(A -> B) & (!B | A)"]',
        '[true, true, false, true]'
    );
INSERT INTO quiz (
        type,
        difficulty,
        num_vars,
        time_limit,
        is_competition_mode,
        question,
        answers,
        solutions
    )
VALUES (
        "TAUT",
        "easy",
        2,
        NULL,
        0,
        "(A | !A) & (B | !B)",
        '["yes", "no"]',
        '[true, false]'
    )