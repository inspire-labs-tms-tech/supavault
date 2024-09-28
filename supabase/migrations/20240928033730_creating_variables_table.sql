create table "public"."variables" (
    "id" text not null,
    "created_at" timestamp with time zone not null default (now() AT TIME ZONE 'utc'::text),
    "description" text not null default ''::text
);


alter table "public"."variables" enable row level security;

CREATE UNIQUE INDEX variables_id_key ON public.variables USING btree (id);

CREATE UNIQUE INDEX variables_pkey ON public.variables USING btree (id);

alter table "public"."variables" add constraint "variables_pkey" PRIMARY KEY using index "variables_pkey";

alter table "public"."variables" add constraint "variables_id_check" CHECK ((id ~ '^[A-Z_]{1,}[A-Z0-9_]{0,}$'::text)) not valid;

alter table "public"."variables" validate constraint "variables_id_check";

alter table "public"."variables" add constraint "variables_id_key" UNIQUE using index "variables_id_key";

grant delete on table "public"."variables" to "anon";

grant insert on table "public"."variables" to "anon";

grant references on table "public"."variables" to "anon";

grant select on table "public"."variables" to "anon";

grant trigger on table "public"."variables" to "anon";

grant truncate on table "public"."variables" to "anon";

grant update on table "public"."variables" to "anon";

grant delete on table "public"."variables" to "authenticated";

grant insert on table "public"."variables" to "authenticated";

grant references on table "public"."variables" to "authenticated";

grant select on table "public"."variables" to "authenticated";

grant trigger on table "public"."variables" to "authenticated";

grant truncate on table "public"."variables" to "authenticated";

grant update on table "public"."variables" to "authenticated";

grant delete on table "public"."variables" to "service_role";

grant insert on table "public"."variables" to "service_role";

grant references on table "public"."variables" to "service_role";

grant select on table "public"."variables" to "service_role";

grant trigger on table "public"."variables" to "service_role";

grant truncate on table "public"."variables" to "service_role";

grant update on table "public"."variables" to "service_role";


