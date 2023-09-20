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
        1,
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
        "medium",
        3,
        30,
        1,
        "(A | B | !C) & !A & C & (!B | !C)",
        '["yes", "no"]',
        '[false, true]'
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
        20,
        1,
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
        "SAT",
        "medium",
        2,
        30,
        1,
        "(A <-> !B) & (A -> B)",
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
        "TAUT",
        "easy",
        2,
        60,
        1,
        "(A | !A) & B",
        '["yes", "no"]',
        '[false, true]'
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
        "medium",
        2,
        30,
        1,
        "(A | !A) & (B | !B)",
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
        "TAUT",
        "hard",
        3,
        20,
        1,
        "(A -> B) | (B -> C) | (C -> A)",
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
        "easy",
        3,
        60,
        1,
        "(A & B) | C",
        '["(A | C) & (B | C)", "(A & C) | (B & C)", "A & (B | C)", "(C | B) & (A | C)"]',
        '[true, false, false, true]'
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
        "medium",
        3,
        40,
        1,
        "(A & B) | (C & A)",
        '["(A & C) | (B & C)", "(A & B) | (C & B)", "A & (B | C)", "A & !(!B & !C)"]',
        '[false, false, true, true]'
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