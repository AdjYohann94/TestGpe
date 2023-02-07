create table work_categories
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    label      text not null
        unique
);

alter table work_categories
    owner to postgres;

create index idx_work_categories_deleted_at
    on work_categories (deleted_at);

create table users
(
    id               bigserial
        primary key,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone,
    first_name       text                        not null,
    last_name        text                        not null,
    email            text                        not null,
    password         text                        not null,
    role             text default 'member'::text not null,
    status           text default 'active'::text,
    phone_number     text,
    zip_code         text,
    address          text,
    city             text,
    work_category_id bigint
        constraint fk_work_categories_users
            references work_categories
);

alter table users
    owner to postgres;

create unique index idx_users_email
    on users (email);

create index idx_users_deleted_at
    on users (deleted_at);

create table quizzes
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    label       text                       not null,
    description text                       not null,
    status      text default 'draft'::text not null
);

alter table quizzes
    owner to postgres;

create table questions
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    description text,
    quiz_id     bigint
        constraint fk_quizzes_questions
            references quizzes
);

alter table questions
    owner to postgres;

create table responses
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    user_id     bigint
        constraint fk_responses_user
            references users,
    value       text,
    question_id bigint
        constraint fk_questions_responses
            references questions
);

alter table responses
    owner to postgres;

create index idx_responses_deleted_at
    on responses (deleted_at);

create index idx_questions_deleted_at
    on questions (deleted_at);

create index idx_quizzes_deleted_at
    on quizzes (deleted_at);

create table quiz_work_categories
(
    quiz_id          bigint not null
        constraint fk_quiz_work_categories_quiz
            references quizzes,
    work_category_id bigint not null
        constraint fk_quiz_work_categories_work_category
            references work_categories,
    primary key (quiz_id, work_category_id)
);

alter table quiz_work_categories
    owner to postgres;
