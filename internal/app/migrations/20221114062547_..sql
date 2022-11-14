-- modify "users" table
ALTER TABLE "users" DROP COLUMN "social_name_change", DROP COLUMN "new", ADD COLUMN "social_name" smallint NULL;
