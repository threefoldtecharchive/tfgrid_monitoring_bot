# tfgrid monitoring bot

This is a bot to monitor the balance in accounts and send warnings if it is under some threshold.

## How to start

- Create a new [telegram bot](README.md#create-a-bot-if-you-dont-have) if you don't have.
- Create a new env file `.env`, for example:

```env
TESTNET_MNEMONIC=<your mainnet mnemonic>
MAINNET_MNEMONIC=<your testnet mnemonic>
BOT_TOKEN=<your token>
CHAT_ID=<your chat ID>
MINS=<number of minutes between each message>
```

- Create a new json file `wallets.json` and add the list of addresses you want to monitor, for example:

```json
{ 
    "testnet": [{ 
        "name": "<your wallet name>", 
        "address": "<your tfchain address>", 
        "threshold": 700 
    }],

    "mainnet": [{ 
        "name": "<your wallet name>", 
        "address": "<your tfchain address>", 
        "threshold": 700 
    }]
}
```

- Get the binary

> Download the latest from the [releases page](https://github.com/threefoldtech/tfgrid_monitoring_bot/releases)

- Run the bot

After downloading the binary

```bash
sudo cp tfgrid_monitoring_bot /usr/local/bin
tfgrid_monitoring_bot -e .env -w wallets.json
```

Where

- `.env` is the environment file
- `wallets.json` is the json file of wallets to be monitored  

## Create a bot if you don't have

- Open telegram app
- Create a new bot
  
```ordered
1. Find telegram bot named "@botfarther"
2. Type /newbot
```

- Get the bot token
  
```ordered
1. In the same bot named "@botfarther"
2. Type /token
3. Choose your bot
```

- Get your chat ID

```ordered
1. Search for @RawDataBot and select Telegram Bot Raw from the drop-down list.
2. In the json returned, you will find it in section message -> chat -> id
```

## Test

```bash
make test
```

## Release

- Check `goreleaser check`
- Create a tag `git tag -a v1.0.2 -m "release v1.0.2"`
- Push the tag `git push origin v1.0.2`
- export GITHUB_TOKEN=<GITHUB_TOKEN>
- goreleaser release
