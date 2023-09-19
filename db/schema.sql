CREATE TABLE "users" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "username" varchar(255) NOT NULL,
    "hashed_password" varchar(255) NOT NULL,
    "salt" varchar(255) NOT NULL
);
CREATE TABLE "sessions" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" integer NOT NULL REFERENCES "users" ("id"),
    "token" varchar(255) NOT NULL
);
CREATE TABLE "quiz" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "type" varchar(4) NOT NULL,
    "difficulty" varchar(6) NOT NULL,
    "num_vars" integer NOT NULL,
    -- time limit in seconds, null if unlimited
    "time_limit" DECIMAL,
    -- true if quiz is in competition mode, false if in practice mode
    "is_competition_mode" boolean NOT NULL,
    -- string representation of the question formula
    "question" text NOT NULL,
    -- json list of possible answers and solutions: [{"answer": "yes", "solution": true}, {"answer": "no", "solution": false}]
    "possible_answers" text NOT NULL
);
-- models the many-to-many relationship between users and quizzes
CREATE TABLE "quiz_participation" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "quiz_id" integer NOT NULL REFERENCES "quiz" ("id"),
    "user_id" integer NOT NULL REFERENCES "users" ("id"),
    "correct" boolean NOT NULL,
    "time" DECIMAL,
    "points" integer
);
CREATE TABLE "achievement" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" varchar(255) NOT NULL,
    "description" text NOT NULL,
    "level" varchar(6) NOT NULL,
    -- SQL query to check if achievement is unlocked
    "sql" text NOT NULL
);
-- models the many-to-many relationship between users and achievements
CREATE TABLE "achieved" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" integer NOT NULL REFERENCES "users" ("id"),
    "achievement_id" integer NOT NULL REFERENCES "achievement" ("id"),
    UNIQUE ("user_id", "achievement_id")
);
-- TODO: Model leaderboard as a view