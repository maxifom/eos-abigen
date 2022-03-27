import {Client as EosioClient} from "./generated/eosio";

const {JsonRpc} = require('eosjs');
const f = require('node-fetch');


(async () => {
    const rpc = new JsonRpc('https://eos.greymass.com', {fetch: f});
    const eosioClient = new EosioClient({rpc})
    const {rows: producers} = await eosioClient.producers();
    for (let producer of producers) {
        console.log(producer.owner, producer.last_claim_time)
    }

    const {rows: eosioRes} = await eosioClient.userres();
    for (let res of eosioRes) {
        console.log(res.owner, res.net_weight.raw, res.cpu_weight.raw, res.ram_bytes)
    }

    const {rows: eosioTokenRes} = await eosioClient.userres({scope: "eosio.token"});
    for (let res of eosioTokenRes) {
        console.log(res.owner, res.net_weight.raw, res.cpu_weight.raw, res.ram_bytes)
    }

})()
