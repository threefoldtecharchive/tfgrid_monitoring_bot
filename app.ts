import { MessageBusClientInterface } from "ts-rmb-client-base";
import { HTTPMessageBusClient } from "ts-rmb-http-client";
import { GridClient, NetworkEnv } from "grid3_client";
import { Telegram } from 'telegraf';
import dotenv from 'dotenv';

dotenv.config()

let rmb: MessageBusClientInterface = new HTTPMessageBusClient(0, "", "", "");

const addresses = ["5Gv79EKyAmEAv3h1SyEY9tqq3a9j6NB1rWnkfJsVQBTTRpHc", "5Covf4nJf8gWiqXbyCBfbaFgWoxnnp3F1ZUiDWixx3Lds9bz", "5FfAv7qUb3oa8TfdeLCEicirtdGvXHnV2owtsv5XNa9aDzEu"]

const token: string = process.env.BOT_TOKEN as string;
const telegram: Telegram = new Telegram(token);
const chatId: string = process.env.CHAT_ID as string;
const mnemonics: string = process.env.MNEMONICS as string;
const minutesProvided : string = process.env.MINS as string;
const mins : number = +minutesProvided;
const tftsLimit : string = process.env.TFTS_LIMIT as string;
const tfts : number = +tftsLimit;

async function monitor(client: GridClient) {
    for (let addr of addresses) {
        console.log(`checking {addr}`)
        const balance = await client.tfchain.balanceByAddress({
                address: addr,
        });
            
        if (balance.free < tfts) {
                let msg = `account with address: ${addr} has balance = ${balance.free}`
                telegram.sendMessage(chatId, msg)
        }   
    }
    await client.disconnect();
}

setInterval(() => {
    async() => {
        const gridClient = new GridClient(
            NetworkEnv.dev,
            mnemonics,
            "secret",
            rmb,
        );
        console.log(`client loaded`)
        await gridClient.connect() 
        monitor(gridClient)
    }
    // num of seconds to wait before monitoring addresses
}, mins * 60000)