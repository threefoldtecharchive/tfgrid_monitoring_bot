# This is a bot to monitor tfchain balances by providing a list of addresses

# To run the project
- clone it in your machine
- Run 
```
cd tfgrid_monitoring_bot
```
- create a .env file containing [BOT_TOKEN, CHAT_ID, MNEMONICS, MINS, TFTS_LIMIT] where mins the time between each call
- provide the list of addresses you want to monitor in app.ts file
- Run
```
yarn install
```
- Run 
```
yarn ts-node app.ts
```
