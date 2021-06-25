/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package dao

import (
	"database/sql"
	"github.com/polynetwork/explorer/internal/common"
	"github.com/polynetwork/explorer/internal/model"
)

const (
	//_insertTChainTx                 = "insert into tchain_tx(chain_id, txhash, state, tt, fee, height, fchain, contract, rtxhash) values (?,?,?,?,?,?,?,?,?)"
	//_insertTChainTransfer          = "insert into tchain_transfer(txhash, asset, xfrom, xto, amount) values(?,?,?,?,?)"
	//_insertFChainTx                 = "insert into fchain_tx(chain_id, txhash, state, tt, fee, height, xuser, tchain, contract, xkey, xparam) values (?,?,?,?,?,?,?,?,?,?,?)"
	//_insertFChainTransfer          = "insert into fchain_transfer(txhash, asset, xfrom, xto, amount, tochainid, toasset, touser) values (?,?,?,?,?,?,?,?)"
	//_insertMChainTx             = "insert into mchain_tx(chain_id, txhash, state, tt, fee, height, fchain, ftxhash, tchain, xkey) values (?,?,?,?,?,?,?,?,?,?)"
	//无使用_selectMChainTxCount        = "select count(*) from mchain_tx"
	//重写_selectMChainTxByLimit      = "select A.chain_id, A.txhash, case when B.txhash is null OR C.txhash is null THEN 0 ELSE 1 END as state, A.tt, A.fee, A.height, A.fchain, A.tchain from mchain_tx A left join tchain_tx B on A.txhash = B.rtxhash left join fchain_tx C on A.ftxhash = C.txhash order by A.height desc limit ?,?;"
	//重写_selectMChainTxByHash       = "select chain_id, txhash, state, tt, fee, height, fchain, ftxhash, tchain, xkey from mchain_tx where txhash = ?"
	//重写_selectMChainTxByFHash      = "select chain_id, txhash, state, tt, fee, height, fchain, ftxhash, tchain, xkey from mchain_tx where ftxhash = ?"
	//重写_selectFChainTxByHash       = "select chain_id, txhash, state, tt, fee, height, xuser, tchain, contract, xkey, xparam from fchain_tx where case when chain_id = ? then xkey = ? else txhash = ? end"
	//重写_selectFChainTxByTime       = "select unix_timestamp(FROM_UNIXTIME(tt,'%Y%m%d')) days, count(*) from fchain_tx where chain_id = ? and tt > ? and tt < ? group by chain_id,days order by days desc"
	//重写_selectFChainTransferByHash = "select txhash, asset, xfrom, xto, amount, tochainid, toasset, touser from fchain_transfer where txhash = ?" // TODO
	//重写_selectTChainTxByHash       = "select chain_id, txhash, state, tt, fee, height, fchain, contract,rtxhash from tchain_tx where txhash = ?"
	//重写_selectTChainTxByMHash      = "select chain_id, txhash, state, tt, fee, height, fchain, contract,rtxhash from tchain_tx where rtxhash = ?"
	//重写_selectTChainTxByTime       = "select unix_timestamp(FROM_UNIXTIME(tt,'%Y%m%d')) days, count(*) from tchain_tx where chain_id = ? and tt > ? and tt < ? group by chain_id,days order by days desc"
	//重写_selectTChainTransferByHash = "select txhash, asset, xfrom, xto, amount from tchain_transfer where txhash = ?"
	//下方new _selectChainAddresses       = "select count(distinct a) from (select distinct xfrom as a from fchain_transfer aa left join fchain_tx bb on aa.txhash = bb.txhash where bb.chain_id = ? union all select distinct xto as a from tchain_transfer aa left join tchain_tx bb on aa.txhash = bb.txhash where bb.chain_id = ?) c"
	//下方不变 _selectChainInfoById        = "select xname, id, url, xtype, height, txin, txout from chain_info where id = ?"
	//下方不变 _selectAllChainInfos        = "select xname, id, url, xtype, height, txin, txout from chain_info order by id"
	//下方不变 _selectContractById         = "select id, contract from chain_contract where id = ?"
	//下方不变 _selectTokenById            = "select id, xtoken, hash, xname, xtype, xprecision, xdesc from chain_token where id = ?"
	//无使用_selectTokenCount           = "select count(distinct xtoken) from chain_token"
	//无使用_updateChainInfoById        = "update chain_info set xname = ?, url = ?, height = ?, txin = ?, txout = ? where id = ?"
	//无使用_selectAllianceTx           = "select chain_id, txhash, state, tt, fee, height, fchain, ftxhash, tchain,xkey from mchain_tx where (tchain = ? or fchain = ?) and height > ? order by height"
	//无使用_selectBitcoinUnconfirmTx   = "select txhash from tchain_tx where chain_id = ? and tt = ?"
	//_updateBitcoinConfirmTx         = "update tchain_tx set tt = ?, height = ?, state = 1, fee = ? where txhash = ?"
	//_updateBitcoinConfirmTransfer  = "update tchain_transfer set xto = ?, amount = ? where txhash = ?"
	//下方new _selectTokenTxList    = "select case when b.chain_id = ? then b.xkey else b.txhash end as txhash,a.xfrom,a.xto,a.amount,b.height,b.tt,1 as direct from fchain_transfer a left join fchain_tx b on a.txhash = b.txhash where a.asset = ? union all select d.txhash,c.xfrom,c.xto,c.amount,d.height,d.tt,2 as direct from tchain_transfer c left join tchain_tx d on c.txhash = d.txhash where c.asset = ? order by height desc limit ?,?;"
	//下方new _selectTokenTxTotal   = "select sum(cnt) from (select count(*) as cnt from fchain_transfer a left join fchain_tx b on a.txhash = b.txhash where a.asset = ? union all select count(*) as cnt from tchain_transfer c left join tchain_tx d on c.txhash = d.txhash where c.asset = ?) t"
	//下方new _selectAddressTxList  = "select case when b.chain_id = ? then b.xkey else b.txhash end as txhash,a.xfrom,a.xto,a.asset,a.amount,b.height,b.tt,1 as direct from fchain_transfer a left join fchain_tx b on a.txhash = b.txhash where a.xfrom = ? and b.chain_id = ? union all select d.txhash,c.xfrom,c.xto,c.asset,c.amount,d.height,d.tt,2 as direct from tchain_transfer c left join tchain_tx d on c.txhash = d.txhash where c.xto = ? and d.chain_id = ? order by height desc limit ?,?;"
	//下方new _selectAddressTxTotal = "select sum(cnt) from (select count(*) as cnt from fchain_transfer a left join fchain_tx b on a.txhash = b.txhash where a.xfrom = ? and b.chain_id = ? union all select count(*) as cnt from tchain_transfer c left join tchain_tx d on c.txhash = d.txhash where c.xto = ? and d.chain_id = ?) t"
	//无使用_insertPolyValidator  = "insert into poly_validators(height, validators) values(?,?)"

	//new
	//_selectMChainTxByLimit
	_selectPolyTransactionByLimit = "select A.chain_id, A.hash, case when B.hash is null OR C.hash is null THEN 0 ELSE 1 END as state, A.time, A.fee, A.height, A.src_chain_id, A.dst_chain_id from poly_transactions A left join dst_transactions B on A.hash = B.hash left join src_transactions C on A.src_hash = C.hash order by A.height desc limit ?,?;"
	//_selectMChainTxByHash
	_selectPolyTransactionByHash = "select chain_id, hash, state, time, fee, height, src_chain_id, src_hash, dst_chain_id, `key` from poly_transactions where hash = ?"
	//_selectMChainTxByFHash
	_selectPolyTransactionBySrcHash = "select chain_id, hash, state, time, fee, height, src_chain_id, src_hash, dst_chain_id, `key` from poly_transactions where src_hash = ?"
	//_selectFChainTxByHash
	_selectSrcTransactionByHash = "select chain_id, hash, state, time, fee, height, user, dst_chain_id, contract, key, param from src_transactions where case when chain_id = ? then `key` = ? else hash = ? end"
	//_selectFChainTxByTime
	_selectSrcTransactionByTime = "select unix_timestamp(FROM_UNIXTIME(time,'%Y%m%d')) days, count(*) from src_transactions where chain_id = ? and time > ? and time < ? group by chain_id,days order by days desc"
	//_selectFChainTransferByHash
	_selectSrcTransferByHash = "select tx_hash, asset, `from`, `to`, amount, dst_chain_id, dst_asset, dst_user from src_transfers where tx_hash = ?"
	//_selectTChainTxByHash
	_selectDstTransactionByHash = "select chain_id, hash, state, time, fee, height, src_chain_id, contract,poly_hash from dst_transactions where hash = ?"
	//_selectTChainTxByMHash
	_selectDstTransactionByPolyHash = "select chain_id, hash, state, hash, fee, height, src_chain_id, contract,poly_hash from dst_transactions where poly_hash = ?"
	//_selectTChainTxByTime
	_selectDstTransactionByTime = "select unix_timestamp(FROM_UNIXTIME(time,'%Y%m%d')) days, count(*) from dst_transactions where chain_id = ? and time > ? and time < ? group by chain_id,days order by days desc"
	//_selectTChainTransferByHash
	_selectDstTransferByHash = "select tx_hash, asset, `from`, `to`, amount from dst_transfers where tx_hash = ?"
	_selectChainAddresses    = "select count(distinct a) from (select distinct `from` as a from src_transfers aa left join src_transactions bb on aa.tx_hash = bb.hash where bb.chain_id = ? union all select distinct `to` as a from dst_transfers aa left join dst_transactions bb on aa.tx_hash = bb.hash where bb.chain_id = ?) c"
	//_selectChainInfoById     = "select xname, id, url, xtype, height, txin, txout from chain_info where id = ?"
	_selectChainStatisticByChainId = "select chain_id, txin, txout from chain_statistic where chain_id = ?"

	_selectAllChainInfos  = "select xname, id, url, xtype, height, txin, txout from chain_info order by id"
	_selectContractById   = "select id, contract from chain_contract where id = ?"
	_selectTokenById      = "select id, xtoken, hash, xname, xtype, xprecision, xdesc from chain_token where id = ?"
	_selectTokenTxList    = "select case when b.chain_id = ? then b.`key` else b.hash end as txhash,a.`from`,a.`to`,a.amount,b.height,b.time,1 as direct from src_transfers a left join src_transactions b on a.tx_hash = b.hash where a.asset = ? union all select d.hash,c.`from`,c.`to`,c.amount,d.height,d.time,2 as direct from dst_transfers c left join dst_transactions d on c.tx_hash = d.hash where c.asset = ? order by height desc limit ?,?;"
	_selectTokenTxTotal   = "select sum(cnt) from (select count(*) as cnt from src_transfers a left join src_transactions b on a.tx_hash = b.hash where a.asset = ? union all select count(*) as cnt from dst_transfers c left join dst_transactions d on c.tx_hash = d.hash where c.asset = ?) t"
	_selectAddressTxList  = "select case when b.chain_id = ? then b.`key` else b.hash end as txhash,a.`from`,a.`to`,a.asset,a.amount,b.height,b.time,1 as direct from src_transfers a left join src_transactions b on a.tx_hash = b.hash where a.`from` = ? and b.chain_id = ? union all select d.hash,c.`from`,c.`to`,c.asset,c.amount,d.height,d.time,2 as direct from dst_transfers c left join dst_transactions d on c.tx_hash = d.hash where c.`to` = ? and d.chain_id = ? order by height desc limit ?,?;"
	_selectAddressTxTotal = "select sum(cnt) from (select count(*) as cnt from src_transfers a left join src_transactions b on a.tx_hash = b.hash where a.`from` = ? and b.chain_id = ? union all select count(*) as cnt from dst_transfers c left join dst_transactions d on c.tx_hash = d.hash where c.`to` = ? and d.chain_id = ?) t"
)

//func (d *Dao) SelectChainInfoById(id uint32) (res *model.ChainInfo, err error) {
//	res = new(model.ChainInfo)
//	row := d.db.QueryRow(_selectChainInfoById, id)
//	if err = row.Scan(&res.Name, &res.Id, &res.Url, &res.XType, &res.Height, &res.In, &res.Out); err != nil {
//		if err == sql.ErrNoRows {
//			err = nil
//		}
//		res = nil
//	}
//	return
//}

func (d *Dao) SelectChainStatisticByChainId(id uint32) (res *model.ChainStatistic, err error) {
	res = new(model.ChainStatistic)
	row := d.db.QueryRow(_selectChainStatisticByChainId, id)
	if err = row.Scan(&res.Chain, &res.In, &res.Out); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		res = nil
	}
	return
}

func (d *Dao) SelectAllChainInfos() (c []*model.ChainInfo, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectAllChainInfos); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ChainInfo)
		if err = rows.Scan(&r.Name, &r.Id, &r.Url, &r.XType, &r.Height, &r.In, &r.Out); err != nil {
			c = nil
			return
		}
		c = append(c, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectContractById(id uint32) (c []*model.ChainContract, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectContractById, id); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ChainContract)
		if err = rows.Scan(&r.Id, &r.Contract); err != nil {
			c = nil
			return
		}
		c = append(c, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectTokenById(id uint32) (c []*model.ChainToken, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectTokenById, id); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ChainToken)
		if err = rows.Scan(&r.Id, &r.Token, &r.Hash, &r.Name, &r.Type, &r.Precision, &r.Desc); err != nil {
			c = nil
			return
		}
		c = append(c, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectMChainTxByHash(hash string) (res *model.PolyTransaction, err error) {
	res = new(model.PolyTransaction)
	row := d.db.QueryRow(_selectPolyTransactionByHash, hash)
	if err = row.Scan(&res.Chain, &res.TxHash, &res.State, &res.TT, &res.Fee, &res.Height, &res.FChain, &res.FTxHash, &res.TChain, &res.Key); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		res = nil
	}
	return
}

func (d *Dao) SelectMChainTxByFHash(hash string) (res *model.PolyTransaction, err error) {
	res = new(model.PolyTransaction)
	row := d.db.QueryRow(_selectPolyTransactionBySrcHash, hash)
	if err = row.Scan(&res.Chain, &res.TxHash, &res.State, &res.TT, &res.Fee, &res.Height, &res.FChain, &res.FTxHash, &res.TChain, &res.Key); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		res = nil
	}
	return
}

func (d *Dao) SelectMChainTxByLimit(start int, limit int) (res []*model.PolyTransaction, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectPolyTransactionByLimit, start, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PolyTransaction)
		if err = rows.Scan(&r.Chain, &r.TxHash, &r.State, &r.TT, &r.Fee, &r.Height, &r.FChain, &r.TChain); err != nil {
			rows = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectFChainTxByHash(hash string, chain uint32) (res *model.SrcTransaction, err error) {
	res = new(model.SrcTransaction)
	res.Transfer = new(model.SrcTransfer)
	{
		row := d.db.QueryRow(_selectSrcTransactionByHash, chain, hash, hash)
		if err = row.Scan(&res.Chain, &res.TxHash, &res.State, &res.TT, &res.Fee, &res.Height, &res.User, &res.TChain, &res.Contract, &res.Key, &res.Param); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			res = nil
			return
		}
	}
	{
		row := d.db.QueryRow(_selectSrcTransferByHash, hash)
		if err = row.Scan(&res.Transfer.TxHash, &res.Transfer.Asset, &res.Transfer.From, &res.Transfer.To, &res.Transfer.Amount, &res.Transfer.ToChain, &res.Transfer.ToAsset, &res.Transfer.ToUser); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
		}
	}
	return
}

func (d *Dao) SelectFChainTxByTime(chainId uint32, start uint32, end uint32) (res []*model.CrossChainTxStatus, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectSrcTransactionByTime, chainId, start, end); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := new(model.CrossChainTxStatus)
		if err = rows.Scan(&r.TT, &r.TxNumber); err != nil {
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectTChainTxByHash(hash string) (res *model.DstTransaction, err error) {
	res = new(model.DstTransaction)
	res.Transfer = new(model.DstTransfer)
	{
		row := d.db.QueryRow(_selectDstTransactionByHash, hash)
		if err = row.Scan(&res.Chain, &res.TxHash, &res.State, &res.TT, &res.Fee, &res.Height, &res.FChain, &res.Contract, &res.RTxHash); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			res = nil
			return
		}
		{
			row := d.db.QueryRow(_selectDstTransferByHash, hash)
			if err = row.Scan(&res.Transfer.TxHash, &res.Transfer.Asset, &res.Transfer.From, &res.Transfer.To, &res.Transfer.Amount); err != nil {
				if err == sql.ErrNoRows {
					err = nil
				}
			}
		}
	}
	return
}

func (d *Dao) SelectTChainTxByMHash(hash string) (res *model.DstTransaction, err error) {
	res = new(model.DstTransaction)
	res.Transfer = new(model.DstTransfer)
	{
		row := d.db.QueryRow(_selectDstTransactionByPolyHash, hash)
		if err = row.Scan(&res.Chain, &res.TxHash, &res.State, &res.TT, &res.Fee, &res.Height, &res.FChain, &res.Contract, &res.RTxHash); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			res = nil
			return
		}
	}
	{
		row := d.db.QueryRow(_selectDstTransferByHash, res.TxHash)
		if err = row.Scan(&res.Transfer.TxHash, &res.Transfer.Asset, &res.Transfer.From, &res.Transfer.To, &res.Transfer.Amount); err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
		}
	}
	return
}

func (d *Dao) SelectTChainTxByTime(chainId uint32, start uint32, end uint32) (res []*model.CrossChainTxStatus, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectDstTransactionByTime, chainId, start, end); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		r := new(model.CrossChainTxStatus)
		if err = rows.Scan(&r.TT, &r.TxNumber); err != nil {
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectChainAddresses(chainId uint32) (uint32, error) {
	row := d.db.QueryRow(_selectChainAddresses, chainId, chainId)
	var counter uint32
	if err := row.Scan(&counter); err != nil {
		return 0, err
	}
	return counter, nil
}

func (d *Dao) SelectTokenTxList(token string, start uint32, end uint32) (res []*model.TokenTx, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectTokenTxList, common.CHAIN_ETH, token, token, start, end); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.TokenTx)
		if err = rows.Scan(&r.TxHash, &r.From, &r.To, &r.Amount, &r.Height, &r.TT, &r.Direct); err != nil {
			rows = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectTokenTxTotal(token string) (count *uint32, err error) {
	count = new(uint32)
	row := d.db.QueryRow(_selectTokenTxTotal, token, token)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		*count = uint32(0)
	}
	return
}

func (d *Dao) SelectAddressTxList(chainId uint32, addr string, start uint32, end uint32) (res []*model.AddressTx, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(_selectAddressTxList, common.CHAIN_ETH, addr, chainId, addr, chainId, start, end); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AddressTx)
		if err = rows.Scan(&r.TxHash, &r.From, &r.To, &r.Asset, &r.Amount, &r.Height, &r.TT, &r.Direct); err != nil {
			rows = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

func (d *Dao) SelectAddressTxTotal(chainId uint32, addr string) (count *uint32, err error) {
	count = new(uint32)
	row := d.db.QueryRow(_selectAddressTxTotal, addr, chainId, addr, chainId)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		*count = uint32(0)
	}
	return
}
