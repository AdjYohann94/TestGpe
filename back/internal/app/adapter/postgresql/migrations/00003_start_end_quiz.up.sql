alter table quizzes
    add started_at timestamp with time zone;

alter table quizzes
    add closed_at timestamp with time zone;
