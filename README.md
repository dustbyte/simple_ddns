# Simple DDNS

A [DNSimple](https://dnsimple.com/) based Dynamic DNS implementation

## Usage

```
Usage: simple_ddns [-h] [-t TOKEN] [--ttl TTL] ARGS...

DynDNS for mere mortals

argument details:
	-h, --help=false        Show this help
	-t, --token=""          DNSimple API token
	--ttl=60                TTL of the record in seconds
```

The domain/zone NS records must be setup with DNSimple. It updates the A record to the current IP address if it already exists, otherwise it creates it.

The token can be provided throught the environment variable `DNSIMPLE_TOKEN`.
