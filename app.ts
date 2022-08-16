import { MessageBusClientInterface } from "ts-rmb-client-base";
import { HTTPMessageBusClient } from "ts-rmb-http-client";
import { GridClient, NetworkEnv } from "grid3_client";
import { Telegram } from 'telegraf';
import dotenv from 'dotenv';

dotenv.config()


let rmb: MessageBusClientInterface = new HTTPMessageBusClient(0, "", "", "");

const token: string = process.env.BOT_TOKEN as string;

const telegram: Telegram = new Telegram(token);

const chatId: string = process.env.CHAT_ID as string;

const mnemonics: string = process.env.MNEMONICS as string;

const minutesProvided : string = process.env.MINS as string;
const mins : number = +minutesProvided;

const tftsLimit : string = process.env.TFTS_LIMIT as string;
const tfts : number = +tftsLimit;

const addresses = ["5Gv79EKyAmEAv3h1SyEY9tqq3a9j6NB1rWnkfJsVQBTTRpHc", "5Gv79EKyAmEAv3h1SyEY9tqq3a9j6NB1rWnkfJsVQBTTRpHc"]



async function monitor(address:string) {
    const gridClient = new GridClient(
            NetworkEnv.dev,
            mnemonics,
            "secret",
            rmb,
    );
    await gridClient.connect()
    

    
        const balance = await gridClient.tfchain.balanceByAddress({
                address: address,
            });
            
            if (balance.free < tfts) {
                let msg = `account with address: ${address} has balance = ${balance.free}`
                telegram.sendMessage(chatId, msg)
            }
        await gridClient.disconnect()
    
    
}

setInterval(() => {
    for(let i = 0; i < addresses.length; i++){
        monitor(addresses[i])
    } 
    // num of mins to wait before monitoring addresses
}, mins * 60000)