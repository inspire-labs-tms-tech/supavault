create table "public"."environment_variables" (
    "id" uuid not null default gen_random_uuid(),
    "created_at" timestamp with time zone not null default (now() AT TIME ZONE 'utc'::text),
    "environment_id" uuid not null,
    "variable_id" text not null,
    "uid" text not null generated always as ((((environment_id)::text || ':'::text) || variable_id)) stored
);


alter table "public"."environment_variables" enable row level security;

CREATE UNIQUE INDEX environment_variables_pkey ON public.environment_variables USING btree (id);

CREATE UNIQUE INDEX environment_variables_uid_key ON public.environment_variables USING btree (uid);

alter table "public"."environment_variables" add constraint "environment_variables_pkey" PRIMARY KEY using index "environment_variables_pkey";

alter table "public"."environment_variables" add constraint "environment_variables_environment_id_fkey" FOREIGN KEY (environment_id) REFERENCES environments(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."environment_variables" validate constraint "environment_variables_environment_id_fkey";

alter table "public"."environment_variables" add constraint "environment_variables_uid_key" UNIQUE using index "environment_variables_uid_key";

alter table "public"."environment_variables" add constraint "environment_variables_variable_id_fkey" FOREIGN KEY (variable_id) REFERENCES variables(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."environment_variables" validate constraint "environment_variables_variable_id_fkey";

grant delete on table "public"."environment_variables" to "anon";

grant insert on table "public"."environment_variables" to "anon";

grant references on table "public"."environment_variables" to "anon";

grant select on table "public"."environment_variables" to "anon";

grant trigger on table "public"."environment_variables" to "anon";

grant truncate on table "public"."environment_variables" to "anon";

grant update on table "public"."environment_variables" to "anon";

grant delete on table "public"."environment_variables" to "authenticated";

grant insert on table "public"."environment_variables" to "authenticated";

grant references on table "public"."environment_variables" to "authenticated";

grant select on table "public"."environment_variables" to "authenticated";

grant trigger on table "public"."environment_variables" to "authenticated";

grant truncate on table "public"."environment_variables" to "authenticated";

grant update on table "public"."environment_variables" to "authenticated";

grant delete on table "public"."environment_variables" to "service_role";

grant insert on table "public"."environment_variables" to "service_role";

grant references on table "public"."environment_variables" to "service_role";

grant select on table "public"."environment_variables" to "service_role";

grant trigger on table "public"."environment_variables" to "service_role";

grant truncate on table "public"."environment_variables" to "service_role";

grant update on table "public"."environment_variables" to "service_role";


