# Steampipe Plugin for People Data Labs

The People Data Labs Steampipe plugin lets you query enriched person intelligence directly from Steampipe. It uses the [People Data Labs Person Enrichment API](https://docs.peopledatalabs.com/docs/person-enrichment-api) to resolve a profile from identifiers such as email, phone number, profile URL, or name.

## Getting Started

1. **Install Steampipe** – Follow the [official installation guide](https://steampipe.io/docs/install).
2. **Clone this repository** – Place the plugin source code in your Steampipe plugins directory.
3. **Build the plugin** – Run `go build` inside this repository to produce the plugin binary (for example `steampipe-plugin-pdl`).
4. **Install the plugin locally** – Copy the compiled binary to your Steampipe plugins directory. One easy approach is:

   ```bash
   mkdir -p ~/.steampipe/plugins/local/pdl
   go build -o ~/.steampipe/plugins/local/pdl/steampipe-plugin-pdl
   ```

   If you prefer to keep the binary in the repository directory, you can instead export `STEAMPIPE_PLUGIN_INSTALL_DIR` and point it at the build output before running Steampipe.
5. **Configure the connection** – Add a connection block to your `~/.steampipe/config/pdl.spc` file:

   ```hcl
   connection "pdl" {
     plugin  = "pdl"
     api_key = "${env.PDL_API_KEY}"
   }
   ```

   Alternatively, set the `PDL_API_KEY` environment variable.

## Tables

| Table | Description |
|-------|-------------|
| `pdl_person` | Enrich a single person profile using one or more identifiers. |

### `pdl_person`

Retrieve a single enriched person profile. You must provide at least one of the supported key columns: `email`, `phone`, `profile`, or `name`.

```sql
select
  data->>'full_name' as full_name,
  data->>'job_title' as job_title,
  likelihood,
  status
from
  pdl_person
where
  email = 'someone@example.com';
```

## Development & Local Testing

This plugin uses Go 1.21 and the Steampipe Plugin SDK v5. To iterate locally:

1. Install the Go toolchain (Go 1.21 or later) and ensure you can download modules from the Go proxy (or vendor dependencies locally).
2. From this repository, run `go build` to compile the plugin. You can also build directly into your Steampipe plugin directory as shown above to skip the manual copy step.
3. Launch Steampipe with `steampipe query`. The CLI will load the plugin from `~/.steampipe/plugins/local/pdl` (or from the directory referenced by `STEAMPIPE_PLUGIN_INSTALL_DIR`).
4. Run queries such as `select * from pdl_person where email = 'someone@example.com';` to verify that enrichment works end-to-end.
5. Execute unit tests (if any are added) using `go test ./...`.

To remove the local build, delete the `~/.steampipe/plugins/local/pdl` directory or rebuild with new changes.

## Contributing

Issues and pull requests are welcome! Let us know what other API endpoints or tables would be useful.
