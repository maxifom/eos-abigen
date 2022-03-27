import {Api} from "eosjs";
import {JsSignatureProvider} from "eosjs/dist/eosjs-jssig";
import {ActionBuilder} from "./generated/greymassnoop";

const {JsonRpc} = require('eosjs');
const f = require('node-fetch');


(async () => {
    const rpc = new JsonRpc('https://eos.greymass.com', {fetch: f});
    const api = new Api({
        rpc,
        // ATTENTION!: JSSignatureProvider is insecure and should not be used in production. Used for demonstration purposes
        signatureProvider: new JsSignatureProvider(["SOME_PRIVATE_KEY"]),
        textDecoder: new TextDecoder(),
        textEncoder: new TextEncoder(),
    })

    const res = await api.transact({
        actions: new ActionBuilder().noop({}).build()
    }, {
        blocksBehind: 3,
        expireSeconds: 30
    })
    console.log(JSON.stringify(res))
})()
