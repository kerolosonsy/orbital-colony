# Security Policy

## Supported code

Security fixes target the latest commit on `main`. The Go game is for local use. Keep it on loopback and run player code you trust. Avoid exposing the game server through a reverse proxy, tunnel, or public port.

The First Signal browser demo runs from static files. It uses a small simulation of the direction buttons and sends no game or profile data to a backend.

## Report a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/kerolosonsy/orbital-colony/security/advisories/new) to contact the maintainer. Include the affected version, component, and impact. Keep credentials, player data, and full logs out of public issues. Redact personal data in private reports and coordinate disclosure with the maintainer.
