alter table "public"."environment_variables" add column "value" text not null default ''::text;

alter table "public"."variables" alter column "default" set default ''::text;

alter table "public"."variables" alter column "default" set not null;


