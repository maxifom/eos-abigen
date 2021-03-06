// Generated by eos-abigen version master
import {JsonRpc} from "eosjs";
import * as types from "./types";

export interface ClientOpts {
    rpc: JsonRpc;
}

export class Client {
    private readonly rpc: JsonRpc;

    public constructor(opts: ClientOpts) {
        this.rpc = opts.rpc;
    }
    
    public async abihash(params?: types.GetTableRowsParams): Promise<types.AbiHashRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.AbiHashRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "abihash",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapAbiHash(r))
        };
    }

    public async bidrefunds(params?: types.GetTableRowsParams): Promise<types.BidRefundRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.BidRefundRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "bidrefunds",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapBidRefund(r))
        };
    }

    public async cpuloan(params?: types.GetTableRowsParams): Promise<types.RexLoanRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexLoanRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "cpuloan",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexLoan(r))
        };
    }

    public async delband(params?: types.GetTableRowsParams): Promise<types.DelegatedBandwidthRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.DelegatedBandwidthRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "delband",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapDelegatedBandwidth(r))
        };
    }

    public async global(params?: types.GetTableRowsParams): Promise<types.EosioGlobalStateRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.EosioGlobalStateRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "global",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapEosioGlobalState(r))
        };
    }

    public async global2(params?: types.GetTableRowsParams): Promise<types.EosioGlobalState2Rows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.EosioGlobalState2RowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "global2",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapEosioGlobalState2(r))
        };
    }

    public async global3(params?: types.GetTableRowsParams): Promise<types.EosioGlobalState3Rows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.EosioGlobalState3RowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "global3",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapEosioGlobalState3(r))
        };
    }

    public async global4(params?: types.GetTableRowsParams): Promise<types.EosioGlobalState4Rows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.EosioGlobalState4RowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "global4",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapEosioGlobalState4(r))
        };
    }

    public async namebids(params?: types.GetTableRowsParams): Promise<types.NameBidRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.NameBidRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "namebids",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapNameBid(r))
        };
    }

    public async netloan(params?: types.GetTableRowsParams): Promise<types.RexLoanRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexLoanRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "netloan",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexLoan(r))
        };
    }

    public async powuporder(params?: types.GetTableRowsParams): Promise<types.PowerupOrderRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.PowerupOrderRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "powup.order",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapPowerupOrder(r))
        };
    }

    public async powupstate(params?: types.GetTableRowsParams): Promise<types.PowerupStateRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.PowerupStateRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "powup.state",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapPowerupState(r))
        };
    }

    public async producers(params?: types.GetTableRowsParams): Promise<types.ProducerInfoRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.ProducerInfoRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "producers",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapProducerInfo(r))
        };
    }

    public async producers2(params?: types.GetTableRowsParams): Promise<types.ProducerInfo2Rows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.ProducerInfo2RowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "producers2",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapProducerInfo2(r))
        };
    }

    public async rammarket(params?: types.GetTableRowsParams): Promise<types.ExchangeStateRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.ExchangeStateRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rammarket",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapExchangeState(r))
        };
    }

    public async refunds(params?: types.GetTableRowsParams): Promise<types.RefundRequestRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RefundRequestRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "refunds",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRefundRequest(r))
        };
    }

    public async retbuckets(params?: types.GetTableRowsParams): Promise<types.RexReturnBucketsRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexReturnBucketsRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "retbuckets",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexReturnBuckets(r))
        };
    }

    public async rexbal(params?: types.GetTableRowsParams): Promise<types.RexBalanceRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexBalanceRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rexbal",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexBalance(r))
        };
    }

    public async rexfund(params?: types.GetTableRowsParams): Promise<types.RexFundRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexFundRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rexfund",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexFund(r))
        };
    }

    public async rexpool(params?: types.GetTableRowsParams): Promise<types.RexPoolRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexPoolRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rexpool",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexPool(r))
        };
    }

    public async rexqueue(params?: types.GetTableRowsParams): Promise<types.RexOrderRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexOrderRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rexqueue",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexOrder(r))
        };
    }

    public async rexretpool(params?: types.GetTableRowsParams): Promise<types.RexReturnPoolRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.RexReturnPoolRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "rexretpool",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapRexReturnPool(r))
        };
    }

    public async userres(params?: types.GetTableRowsParams): Promise<types.UserResourcesRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.UserResourcesRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "userres",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapUserResources(r))
        };
    }

    public async voters(params?: types.GetTableRowsParams): Promise<types.VoterInfoRows> {
        let lower_bound = params?.lower_bound || params?.bounds || undefined;
        let upper_bound = params?.upper_bound || params?.bounds || undefined;
		let result: types.VoterInfoRowsInterm = await this.rpc.get_table_rows({
            json: true,
            code: params?.code || types.CONTRACT_NAME,
            scope: params?.scope || types.CONTRACT_NAME,
            table: "voters",
            lower_bound: lower_bound,
            upper_bound: upper_bound,
            index_position: params?.index_position,
            key_type: params?.key_type,
            limit: params?.limit,
            reverse: params?.reverse,
            show_payer: params?.show_payer,
        });

		return {
			more: result.more,
			next_key: result.next_key,
			rows: result.rows.map(r => types.mapVoterInfo(r))
        };
    }
}