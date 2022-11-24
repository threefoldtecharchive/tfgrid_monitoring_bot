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
    "testnet": [{ "addr": "<your address>" , "threshold": 700 }],
    "mainnet": [{ "addr": "<your address>" , "threshold": 700 }]
}
```

- Get the binary
Download the latest from the [releases page](https://github.com/threefoldtech/tfgrid_monitoring_bot/releases/tag/v1.0.0) 

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
