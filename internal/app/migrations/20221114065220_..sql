-- modify "user_yoga" table
ALTER TABLE "user_yoga" ALTER COLUMN "user_type" TYPE smallint USING(user_type::smallint);
