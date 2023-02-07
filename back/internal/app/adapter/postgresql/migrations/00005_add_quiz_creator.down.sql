ALTER TABLE "quizzes" DROP CONSTRAINT "fk_quizzes_user";
ALTER TABLE "quizzes" DROP COLUMN "creator_id";
ALTER TABLE "quizzes" rename "name" to "label";