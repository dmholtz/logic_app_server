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
-- model quiz type enum as a table
CREATE TABLE "quiz_type" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" varchar(4) NOT NULL
);
-- model quiz difficulty enum as a table
CREATE TABLE "quiz_difficulty" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" varchar(10) NOT NULL
);
CREATE TABLE "quiz" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "type" integer NOT NULL REFERENCES "quiz_type" ("id"),
    "difficulty" integer NOT NULL REFERENCES "quiz_difficulty" ("id"),
    "num_vars" integer NOT NULL,
    -- time limit in seconds, null if unlimited
    "time_limit" integer,
    -- true if quiz is in competition mode, false if in practice mode
    "is_competition_mode" boolean NOT NULL,
    -- json representation of the question formula
    "question" text NOT NULL,
    -- json list of possible answers and solutions
    "answer" text NOT NULL
);
CREATE TABLE "quiz_participation" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "quiz_id" integer NOT NULL REFERENCES "quiz" ("id"),
    "user_id" integer NOT NULL REFERENCES "users" ("id"),
    "timestamp" datetime NOT NULL,
    "score" integer NOT NULL
);
-- TODO: Model leaderboard as a view