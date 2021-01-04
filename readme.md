# Altenator - Minimal STG

1. [Start using alternator](#start-using-alternator)
2. [Install as a Systemd service](#install-as-a-systemd-service)

## Start using alternator

1. Add alternator to system path.
2. Init new empty project `altenator --init`.
3. Start watching for file changes and start web server on port 8000 with `altenator --watch --http`.

## Install as a Systemd service

```bash
cp altenator.service /lib/systemd/system/altenator.service

service altenator enable

service altenator start

service altenator status
```
