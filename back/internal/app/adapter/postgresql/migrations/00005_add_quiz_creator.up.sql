ALTER TABLE "quizzes" ADD "creator_id" bigint not null default 1;
ALTER TABLE "quizzes" ALTER COLUMN "creator_id" DROP DEFAULT;
ALTER TABLE "quizzes" ADD CONSTRAINT "fk_quizzes_user" FOREIGN KEY ("creator_id") REFERENCES "users"("id");
ALTER TABLE "quizzes" rename "label" to "name";