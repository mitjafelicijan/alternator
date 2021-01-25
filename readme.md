# Altenator - Minimal STG

- [Download binary](#download-binary)
- [Start using alternator](#start-using-alternator)
- [Install as a Systemd service](#install-as-a-systemd-service)

## Download binary

```sh
# amd64
wget https://github.com/mitjafelicijan/alternator/raw/master/release/linux-amd64/alternator

# arm
wget https://github.com/mitjafelicijan/alternator/raw/master/release/linux-arm/alternator

# add executable bit
chmod +x alternator

# check if bin is working
./alternator --version
```

## Start using alternator

1. Add alternator to system path.
2. Init new empty project `altenator --init`.
3. Start watching for file changes and start web server on port 8000 with `altenator --watch --http`.

## Install as a Systemd service

**Do not use yet!**

```bash
cp altenator.service /lib/systemd/system/altenator.service

service altenator enable

service altenator start

service altenator status
```
