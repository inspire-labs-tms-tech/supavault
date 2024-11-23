create policy "select based on client"
on "public"."environment_variables"
as permissive
for select
to authenticated
using ((environment_id = ( SELECT client.environment_id
   FROM clients client
 LIMIT 1)));


create policy "select based on client"
on "public"."variables"
as permissive
for select
to authenticated
using ((project_id = ( SELECT env.project_id
   FROM environments env
 LIMIT 1)));



