CREATE SEQUENCE message_tasks_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."message_tasks" (
  "id" bigint DEFAULT nextval('message_tasks_id_seq') NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "message_content" text,
  "to" text,
  "status" text,
  CONSTRAINT "message_tasks_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
CREATE INDEX "idx_message_tasks_deleted_at" ON "public"."message_tasks" USING btree ("deleted_at");
CREATE SEQUENCE task_states_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."task_states" (
  "id" bigint DEFAULT nextval('task_states_id_seq') NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "last_successful_query_time" timestamptz,
  "status" text,
  "active" boolean,
  CONSTRAINT "task_states_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
CREATE INDEX "idx_task_states_deleted_at" ON "public"."task_states" USING btree ("deleted_at");
CREATE TYPE task_status AS ENUM (
  'WAITING',
  'PROCESSING',
  'COMPLETED',
  'FAILED'
);