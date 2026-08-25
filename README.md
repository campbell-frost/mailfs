# MailFS

A file store that persists files over email messages. This is theoretically infinite free cloud storage, at the risk of TOS violation.

## Requirements

- Go
- Node + [pnpm](https://pnpm.io)
- [Task](https://taskfile.dev)
- A Gmail account with 2FA enabled and an [app password](https://myaccount.google.com/apppasswords).

## Setup

```bash
cp .env.example .env
go mod download
pnpm --dir web i
```

Fill in `.env` with your email credentials. Make sure you create an app password, dont use your real password.

## Run

```bash
task dev
```

## Debugging

To debug the API, you have two options:

1. Start the API via the `mailfs API (launch)` Zed debugger config and separately start the frontend with `task web`
2. Run `task dev:debug`, then attach with the `mailfs API (attach :2112)` config.

   SIGINT (Ctrl-C) wont stop `dlv` cleanly, it forwards the signal to the app instead.
   Stop the debugger from Zed, or send a SIGTERM.

TODO: add vscode launch commands
