ALTER TABLE "questions" ADD "question_type" text;

CREATE TABLE "response_choices" (
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    "question_id" bigint,
    "type" text,
    "value" text,
    CONSTRAINT "fk_response_choices_question" FOREIGN KEY ("question_id") REFERENCES "questions"("id")
);

CREATE INDEX IF NOT EXISTS "idx_response_choices_deleted_at" ON "response_choices" ("deleted_at");

ALTER TABLE "responses" ADD "response_choice_id" bigint;

ALTER TABLE "responses" ADD CONSTRAINT "fk_responses_response_choice"
    FOREIGN KEY ("response_choice_id") REFERENCES "response_choices"("id");
