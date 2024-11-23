alter table "public"."clients" add column "environment_id" uuid;

alter table "public"."clients" add constraint "clients_environment_id_fkey" FOREIGN KEY (environment_id) REFERENCES environments(id) ON UPDATE CASCADE not valid;

alter table "public"."clients" validate constraint "clients_environment_id_fkey";

create policy "select own client"
on "public"."clients"
as permissive
for select
to authenticated
using ((( SELECT auth.uid() AS uid) = id));


create policy "select based on client"
on "public"."environments"
as permissive
for select
to authenticated
using ((id = ( SELECT client.environment_id
   FROM clients client
 LIMIT 1)));



