UPDATE "users" SET work_category_id = 1 WHERE work_category_id is null;
ALTER TABLE "users" ALTER COLUMN work_category_id set not null;
