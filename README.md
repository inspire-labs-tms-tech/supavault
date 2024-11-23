# supavault

A Supabase Key-Store

## Dependencies

A clean Supabase instance to install Supavault into. This can be either a Supabase-hosted or Self-hosted instance, so
long as it is accessible on the network (or over the internet is used publicly).

## Getting Started

- Install the Supavault CLI 

- Run the `update` command for the first time (this will apply all migrations)

Example:

```shell
supavault \
  --db-host="<region-info>.pooler.supabase.com" \
  --db-port=6543 \
  --db-user="postgres.<project-id>" \
  --db-pass="super-secret-password-here" \
  admin update
```

Or, if running Supabase locally, just use:

```shell
supavault admin update
```


