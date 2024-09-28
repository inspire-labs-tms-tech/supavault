alter table "public"."secrets" add column "default" text;

alter table "public"."secrets" add column "project_id" uuid not null;

alter table "public"."secrets" enable row level security;

alter table "public"."variables" add column "default" text;

alter table "public"."variables" add column "project_id" uuid not null;

alter table "public"."secrets" add constraint "secrets_project_id_fkey" FOREIGN KEY (project_id) REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."secrets" validate constraint "secrets_project_id_fkey";

alter table "public"."variables" add constraint "variables_project_id_fkey" FOREIGN KEY (project_id) REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."variables" validate constraint "variables_project_id_fkey";


