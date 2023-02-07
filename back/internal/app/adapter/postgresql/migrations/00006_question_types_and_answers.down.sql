ALTER TABLE "questions" DROP COLUMN "question_type";

DROP TABLE "response_choices";

ALTER TABLE "responses" DROP CONSTRAINT "fk_responses_response_choice";
ALTER TABLE "responses" DROP COLUMN "response_choice_id";
