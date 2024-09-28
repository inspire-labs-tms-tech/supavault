create table "public"."clients" (
    "id" uuid not null
);


alter table "public"."clients" enable row level security;

CREATE UNIQUE INDEX clients_id_key ON public.clients USING btree (id);

CREATE UNIQUE INDEX clients_pkey ON public.clients USING btree (id);

CREATE UNIQUE INDEX environment_variables_id_key ON public.environment_variables USING btree (id);

CREATE UNIQUE INDEX environments_id_key ON public.environments USING btree (id);

CREATE UNIQUE INDEX projects_id_key ON public.projects USING btree (id);

alter table "public"."clients" add constraint "clients_pkey" PRIMARY KEY using index "clients_pkey";

alter table "public"."clients" add constraint "clients_id_fkey" FOREIGN KEY (id) REFERENCES auth.users(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."clients" validate constraint "clients_id_fkey";

alter table "public"."clients" add constraint "clients_id_key" UNIQUE using index "clients_id_key";

alter table "public"."environment_variables" add constraint "environment_variables_id_key" UNIQUE using index "environment_variables_id_key";

alter table "public"."environments" add constraint "environments_id_key" UNIQUE using index "environments_id_key";

alter table "public"."projects" add constraint "projects_id_key" UNIQUE using index "projects_id_key";

grant delete on table "public"."clients" to "anon";

grant insert on table "public"."clients" to "anon";

grant references on table "public"."clients" to "anon";

grant select on table "public"."clients" to "anon";

grant trigger on table "public"."clients" to "anon";

grant truncate on table "public"."clients" to "anon";

grant update on table "public"."clients" to "anon";

grant delete on table "public"."clients" to "authenticated";

grant insert on table "public"."clients" to "authenticated";

grant references on table "public"."clients" to "authenticated";

grant select on table "public"."clients" to "authenticated";

grant trigger on table "public"."clients" to "authenticated";

grant truncate on table "public"."clients" to "authenticated";

grant update on table "public"."clients" to "authenticated";

grant delete on table "public"."clients" to "service_role";

grant insert on table "public"."clients" to "service_role";

grant references on table "public"."clients" to "service_role";

grant select on table "public"."clients" to "service_role";

grant trigger on table "public"."clients" to "service_role";

grant truncate on table "public"."clients" to "service_role";

grant update on table "public"."clients" to "service_role";


