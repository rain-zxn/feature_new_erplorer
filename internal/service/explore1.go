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

package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/polynetwork/explorer/internal/common"
	"github.com/polynetwork/explorer/internal/log"
	"github.com/polynetwork/explorer/internal/model"
	myerror "github.com/polynetwork/explorer/internal/server/restful/error"
	"strconv"
	"strings"
)


type ExplorerController struct {
	beego.Controller
}

// GetExplorerInfo shows explorer information, such as current blockheight (the number of blockchain and so on) on the home page.
func (c *ExplorerController) GetExplorerInfo() () {
	// get parameter
	var explorerReq model.ExplorerInfoReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &explorerReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//get all chains
	chains := make([]*model.Chain, 0)
	res := db.Find(&chains)
	if res.RowsAffected == 0 {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("chain does not exist"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	// get all chains statistic
	chainStatistics := make([]*model.ChainStatistic, 0)

	// get all tokens
	tokenBasics := make([]*model.TokenBasic, 0)
	res = db.Find(&tokenBasics)
	if res.RowsAffected == 0 {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("chain does not exist"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}

	c.Data["json"] = model.MakeExplorerInfoResp(chains, chainStatistics, tokenBasics)
	c.ServeJSON()
}

func (c *ExplorerController) GetTokenTxList() {
	// get parameter
	var tokenTxListReq model.TokenTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &tokenTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//
	transactionOnTokens := make([]*model.TransactionOnToken, 0)
	db.Raw("select a.hash, a.height, a.time, a.chain_id, b.from, b.to, b.amount, 1 as direct from src_transactions a inner join src_transfers b on a.hash = b.tx_hash where b.asset = ? and b.chain_id = ?" +
		"union select c.hash, c.height, c.time, c.chain_id, d.from, d.to, d.amount, 2 as direct from dst_transactions c inner join dst_transfers d on c.hash = d.tx_hash where d.asset = ? and d.chain_id = ?" +
		"order by height desc limit ?,?",
		tokenTxListReq.Token, tokenTxListReq.ChainId, tokenTxListReq.Token, tokenTxListReq.ChainId, tokenTxListReq.PageSize * tokenTxListReq.PageNo, tokenTxListReq.PageSize).Find(&transactionOnTokens)
	//
	tokenStatistic := new(model.TokenStatistic)
	db.Where("chain_id = ? and hash = ?", tokenTxListReq.ChainId, tokenTxListReq.Token).Find(tokenStatistic)
	//
	c.Data["json"] = model.MakeTokenTxList(transactionOnTokens, tokenStatistic)
	c.ServeJSON()
}

func (c *ExplorerController) GetAddressTxList() {
	// get parameter
	var addressTxListReq model.AddressTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &addressTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}

	//
	transactionOnAddresses := make([]*model.TransactionOnAddress, 0)
	db.Raw("select a.hash, a.height, a.time, a.chain_id, b.from, b.to, b.amount, c.hash as token_hash, c.type as token_type, c.name as token_name 1 as direct from src_transactions a inner join src_transfers b on a.hash = b.tx_hash inner join tokens c on b.asset = c.hash and b.chain_id = c.chain_id where b.from = ? and b.chain_id = ?" +
		"union select d.hash, d.height, d.time, d.chain_id, e.from, e.to, e.amount, f.hash as token_hash, f.type as token_type, f.name as token_name, 2 as direct from dst_transactions d inner join dst_transfers e on d.hash = e.tx_hash inner join tokens f on e.asset = f.hash and e.chain_id = f.chain_id where e.to = ? and e.chain_id = ?" +
		"order by height desc limit ?,?",
		addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.PageSize * addressTxListReq.PageNo, addressTxListReq.PageSize).Find(&transactionOnAddresses)
	//

	counter := new(model.Counter)
	db.Raw("select sum(cnt) as counter from (select count(*) as cnt from src_transactions a inner join src_transfers b on a.hash = b.tx_hash inner join tokens c on b.asset = c.hash and b.chain_id = c.chain_id where b.from = ? and b.chain_id = ?" +
		"union count(*) as cnt from dst_transactions d inner join dst_transfers e on d.hash = e.tx_hash inner join tokens f on e.asset = f.hash and e.chain_id = f.chain_id where e.to = ? and e.chain_id = ?) as u",
		addressTxListReq.Address, addressTxListReq.ChainId, addressTxListReq.Address, addressTxListReq.ChainId).Find(counter)
	//
	//
	c.Data["json"] = model.MakeAddressTxList(transactionOnAddresses, counter.Counter)
	c.ServeJSON()
}

// TODO GetCrossTxList gets Cross transaction list from start to end (to be optimized)
func (c *ExplorerController) GetCrossTxList()  {
	// get parameter
	var crossTxListReq model.CrossTxListReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &crossTxListReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}


	tokenStatistics := make([]*model.TokenStatistic, 0)
	db.Raw("select count(*) in_counter,  sum(amount) as in_amount, asset as hash, chain_id as chain_id from dst_transfers group by chain_id, asset").
		Preload("Token").Preload("Token.TokenBasic").
		Find(&tokenStatistics)
	for _, tokenStatistic := range tokenStatistics {
		tokenStatistic.InAmountUsdt = tokenStatistic.InAmount/tokenStatistic.Token.Precision * tokenStatistic.Token.TokenBasic.Price
	}
	db.Save(&tokenStatistics)

	db.Table("dst_transfers").Select("")
	db.Model(&model.Token{}).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash").
		Where("src_transactions.standard = ?", 0).
		Joins("left join src_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Limit(crossTxListReq.PageSize).Offset(crossTxListReq.PageSize * crossTxListReq.PageNo).
		Find(&srcPolyDstRelations)


	srcPolyDstRelations := make([]*model.SrcPolyDstRelation, 0)
	db.Model(&model.PolyTransaction{}).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash").
		Where("src_transactions.standard = ?", 0).
		Joins("left join src_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Limit(crossTxListReq.PageSize).Offset(crossTxListReq.PageSize * crossTxListReq.PageNo).
		Find(&srcPolyDstRelations)

	var transactionNum int64
	db.Model(&model.PolyTransaction{}).Where("src_transactions.standard = ?", 0).
		Joins("left join src_transactions on src_transactions.hash = poly_transactions.src_hash").Count(&transactionNum)

	c.Data["json"] = model.MakeCrossTxListResp(srcPolyDstRelations)
	c.ServeJSON()
}

// GetCrossTx gets cross tx by Tx
func (c *ExplorerController) GetCrossTx() {
	var crossTxReq model.CrossTxReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &crossTxReq); err != nil {
		c.Data["json"] = model.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}
	srcPolyDstRelations := make([]*model.SrcPolyDstRelation, 0)
	db.Model(&model.SrcTransaction{}).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash, src_transactions.chain_id as chain_id, src_transactions.asset as token_hash").
		Where("src_transactions.standard = ? and (src_transactions.hash = ? or poly_transactions.hash = ? or dst_transactions.hash = ?)", 0, crossTxReq.TxHash, crossTxReq.TxHash, crossTxReq.TxHash).
		Joins("left join src_transfers on src_transactions.hash = src_transfers.tx_hash").
		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Preload("Token").
		Preload("Token.TokenBasic").
		Find(&srcPolyDstRelations)
	c.Data["json"] = model.MakeCrossTxResp(srcPolyDstRelations)
	c.ServeJSON()
}

func (c *ExplorerController) GetAssetStatistic() {

}

func (c *ExplorerController) GetTransferStatistic() {

}

func (exp *Service) outputChainInfos(chainInfos []*model.ChainInfo) []*model.ChainInfoResp {
	chainInfoResps := make([]*model.ChainInfoResp, 0)
	for _, chainInfo := range chainInfos {
		chainInfoResp := &model.ChainInfoResp{
			Id:     chainInfo.Id,
			Name:   chainInfo.Name,
			Height: chainInfo.Height,
			In:     chainInfo.In,
			Out:    chainInfo.Out,
		}
		chainInfoResps = append(chainInfoResps, chainInfoResp)
	}
	return chainInfoResps
}

func (exp *Service) outputChainContracts(chainContracts []*model.ChainContract) []*model.ChainContractResp {
	chainContractResps := make([]*model.ChainContractResp, 0)
	for _, chainContract := range chainContracts {
		chainContractResp := &model.ChainContractResp{
			Id:       chainContract.Id,
			Contract: chainContract.Contract,
		}
		chainContractResps = append(chainContractResps, chainContractResp)
	}
	return chainContractResps
}

func (exp *Service) outputChainTokens(chainTokens []*model.ChainToken) []*model.ChainTokenResp {
	chainTokenResps := make([]*model.ChainTokenResp, 0)
	for _, chainToken := range chainTokens {
		chainTokenResp := &model.ChainTokenResp{
			Chain:     chainToken.Id,
			ChainName: exp.ChainId2Name(uint32(chainToken.Id)),
			Hash:      chainToken.Hash,
			Token:     chainToken.Token,
			Name:      chainToken.Name,
			Type:      chainToken.Type,
			Precision: chainToken.Precision,
			Desc:      chainToken.Desc,
		}
		chainTokenResps = append(chainTokenResps, chainTokenResp)
	}
	return chainTokenResps
}

func (exp *Service) outputCrossChainTxStatus(status []*model.CrossChainTxStatus, start uint32, end uint32) []*model.CrossChainTxStatus {
	status_new := make([]*model.CrossChainTxStatus, 0)
	current_day := exp.DayOfTime(start)
	end_day := exp.DayOfTime(end)
	for current_day <= end_day {
		status_new = append(status_new, &model.CrossChainTxStatus{
			TT:       current_day,
			TxNumber: 0,
		})
		current_day = exp.DayOfTimeAddOne(current_day)
	}
	i := 0
	j := 0
	for i < len(status) && j < len(status_new) {
		if status[i].TT == status_new[j].TT {
			status_new[j].TxNumber = status[i].TxNumber
			i++
			j++
		} else {
			j++
		}
	}
	return status_new
}

func (exp *Service) outputCrossChainTxStatus1(status []*model.CrossChainTxStatus, start uint32, end uint32, total uint32) []*model.CrossChainTxStatus {
	if len(status) == 0 {
		return nil
	}
	current_txnumber := total
	current_tt := uint32(0)
	status_new := make([]*model.CrossChainTxStatus, 0)
	for _, s := range status {
		status_new = append(status_new, &model.CrossChainTxStatus{
			TT:       s.TT,
			TxNumber: current_txnumber,
		})
		current_txnumber = current_txnumber - s.TxNumber
		current_tt = s.TT
	}
	status_new = append(status_new, &model.CrossChainTxStatus{
		TT:       exp.DayOfTimeSubOne(current_tt),
		TxNumber: current_txnumber,
	})

	status_new1 := make([]*model.CrossChainTxStatus, 0)
	current_txnumber = status_new[0].TxNumber
	current_tt = status_new[0].TT
	for _, s := range status_new {
		for s.TT < current_tt {
			status_new1 = append(status_new1, &model.CrossChainTxStatus{
				TT:       current_tt,
				TxNumber: current_txnumber,
			})
			current_tt = exp.DayOfTimeSubOne(current_tt)
		}

		current_txnumber = s.TxNumber
		status_new1 = append(status_new1, &model.CrossChainTxStatus{
			TT:       current_tt,
			TxNumber: current_txnumber,
		})
		current_tt = exp.DayOfTimeSubOne(current_tt)
	}

	status_new = make([]*model.CrossChainTxStatus, 0)
	ss := status_new1[len(status_new1)-1]
	current_txnumber = ss.TxNumber
	current_tt = exp.DayOfTime(start)
	for current_tt < ss.TT {
		status_new = append(status_new, &model.CrossChainTxStatus{
			TT:       current_tt,
			TxNumber: current_txnumber,
		})
		current_tt = exp.DayOfTimeAddOne(current_tt)
	}
	for i := 0; i < len(status_new1); i++ {
		bb := status_new1[len(status_new1)-1-i]
		current_tt = bb.TT
		current_txnumber = bb.TxNumber
		if current_tt > exp.DayOfTimeUp(end) {
			break
		}
		status_new = append(status_new, bb)
	}
	for current_tt < exp.DayOfTimeUp(end) {
		status_new = append(status_new, &model.CrossChainTxStatus{
			TT:       current_tt,
			TxNumber: current_txnumber,
		})
		current_tt = exp.DayOfTimeAddOne(current_tt)
	}
	return status_new
}

func (exp *Service) outputCrossTransfer(chainid uint32, user string, transfer *model.SrcTransfer) *model.CrossTransferResp {
	if transfer == nil {
		return nil
	}
	crossTransfer := new(model.CrossTransferResp)
	crossTransfer.CrossTxType = 1
	crossTransfer.CrossTxName = exp.TxType2Name(crossTransfer.CrossTxType)
	crossTransfer.FromChainId = chainid
	crossTransfer.FromChain = exp.ChainId2Name(crossTransfer.FromChainId)
	crossTransfer.FromAddress = user
	crossTransfer.ToChainId = transfer.ToChain
	crossTransfer.ToChain = exp.ChainId2Name(crossTransfer.ToChainId)
	crossTransfer.ToAddress = transfer.ToUser
	token := exp.GetToken(transfer.Asset)
	if token != nil {
		crossTransfer.TokenHash = token.Hash
		crossTransfer.TokenName = token.Name
		crossTransfer.TokenType = token.Type
		crossTransfer.Amount = exp.FormatAmount(token.Precision, transfer.Amount)
	}
	return crossTransfer
}

func (exp *Service) outputFChainTx(fChainTx *model.SrcTransaction) *model.FChainTxResp {
	fChainTxResp := &model.FChainTxResp{
		ChainId:    fChainTx.Chain,
		ChainName:  exp.ChainId2Name(fChainTx.Chain),
		TxHash:     fChainTx.TxHash,
		State:      fChainTx.State,
		TT:         fChainTx.TT,
		Fee:        exp.FormatFee(fChainTx.Chain, fChainTx.Fee),
		Height:     fChainTx.Height,
		User:       fChainTx.User,
		TChainId:   fChainTx.TChain,
		TChainName: exp.ChainId2Name(fChainTx.TChain),
		Contract:   fChainTx.Contract,
		Key:        fChainTx.Key,
		Param:      fChainTx.Param,
	}
	if fChainTx.Transfer != nil {
		fChainTxResp.Transfer = &model.FChainTransferResp{
			From:        fChainTx.Transfer.From,
			To:          fChainTx.Transfer.To,
			Amount:      strconv.FormatUint(fChainTx.Transfer.Amount, 10),
			ToChain:     fChainTx.Transfer.ToChain,
			ToChainName: exp.ChainId2Name(fChainTx.Transfer.ToChain),
			ToUser:      fChainTx.Transfer.ToUser,
		}
		token := exp.GetToken(fChainTx.Transfer.Asset)
		fChainTxResp.Transfer.TokenHash = fChainTx.Transfer.Asset
		if token != nil {
			fChainTxResp.Transfer.TokenHash = token.Hash
			fChainTxResp.Transfer.TokenName = token.Name
			fChainTxResp.Transfer.TokenType = token.Type
			fChainTxResp.Transfer.Amount = exp.FormatAmount(token.Precision, fChainTx.Transfer.Amount)
		}
		totoken := exp.GetToken(fChainTx.Transfer.ToAsset)
		fChainTxResp.Transfer.ToTokenHash = fChainTx.Transfer.ToAsset
		if totoken != nil {
			fChainTxResp.Transfer.ToTokenHash = totoken.Hash
			fChainTxResp.Transfer.ToTokenName = totoken.Name
			fChainTxResp.Transfer.ToTokenType = totoken.Type
		}
	}
	if fChainTx.Chain == common.CHAIN_ETH {
		fChainTxResp.TxHash = "0x" + fChainTx.Key
	} else if fChainTx.Chain == common.CHAIN_COSMOS {
		fChainTxResp.TxHash = strings.ToUpper(fChainTxResp.TxHash)
	}
	return fChainTxResp
}

func (exp *Service) outputMChainTx(mChainTx *model.PolyTransaction) *model.MChainTxResp {
	mChainTxResp := &model.MChainTxResp{
		ChainId:    mChainTx.Chain,
		ChainName:  exp.ChainId2Name(mChainTx.Chain),
		TxHash:     mChainTx.TxHash,
		State:      mChainTx.State,
		TT:         mChainTx.TT,
		Fee:        exp.FormatFee(mChainTx.Chain, mChainTx.Fee),
		Height:     mChainTx.Height,
		FChainId:   mChainTx.FChain,
		FChainName: exp.ChainId2Name(mChainTx.FChain),
		FTxHash:    mChainTx.FTxHash,
		TChainId:   mChainTx.TChain,
		TChainName: exp.ChainId2Name(mChainTx.TChain),
		Key:        mChainTx.Key,
	}
	return mChainTxResp
}

func (exp *Service) outputTChainTx(tChainTx *model.DstTransaction) *model.TChainTxResp {
	tChainTxResp := &model.TChainTxResp{
		ChainId:    tChainTx.Chain,
		ChainName:  exp.ChainId2Name(tChainTx.Chain),
		TxHash:     tChainTx.TxHash,
		State:      tChainTx.State,
		TT:         tChainTx.TT,
		Fee:        exp.FormatFee(tChainTx.Chain, tChainTx.Fee),
		Height:     tChainTx.Height,
		FChainId:   tChainTx.FChain,
		FChainName: exp.ChainId2Name(tChainTx.FChain),
		Contract:   tChainTx.Contract,
		RTxHash:    tChainTx.RTxHash,
	}
	if tChainTx.Transfer != nil {
		tChainTxResp.Transfer = &model.TChainTransferResp{
			From:   tChainTx.Transfer.From,
			To:     tChainTx.Transfer.To,
			Amount: strconv.FormatUint(tChainTx.Transfer.Amount, 10),
		}
		token := exp.GetToken(tChainTx.Transfer.Asset)
		tChainTxResp.Transfer.TokenHash = tChainTx.Transfer.Asset
		if token != nil {
			tChainTxResp.Transfer.TokenHash = token.Hash
			tChainTxResp.Transfer.TokenName = token.Name
			tChainTxResp.Transfer.TokenType = token.Type
			tChainTxResp.Transfer.Amount = exp.FormatAmount(token.Precision, tChainTx.Transfer.Amount)
		}
	}
	if tChainTx.Chain == common.CHAIN_ETH {
		tChainTxResp.TxHash = "0x" + tChainTxResp.TxHash
	} else if tChainTx.Chain == common.CHAIN_COSMOS {
		tChainTxResp.TxHash = strings.ToUpper(tChainTxResp.TxHash)
	}
	return tChainTxResp
}

func (exp *Service) outputCrossTxList(crossTxs []*model.PolyTransaction) *model.CrossTxListResp {
	var crossTxListResp model.CrossTxListResp
	crossTxListResp.CrossTxList = make([]*model.CrossTxOutlineResp, 0)
	for _, mChainTx := range crossTxs {
		crossTxListResp.CrossTxList = append(crossTxListResp.CrossTxList, &model.CrossTxOutlineResp{
			TxHash:     mChainTx.TxHash,
			State:      mChainTx.State,
			TT:         mChainTx.TT,
			Fee:        mChainTx.Fee,
			Height:     mChainTx.Height,
			FChainId:   mChainTx.FChain,
			FChainName: exp.ChainId2Name(mChainTx.FChain),
			TChainId:   mChainTx.TChain,
			TChainName: exp.ChainId2Name(mChainTx.TChain),
		})
	}
	return &crossTxListResp
}

func (exp *Service) outputTokenTxList(tokenHash string, tokenTxs []*model.TokenTx, tokenTxTotal uint32) *model.TokenTxListResp {
	var tokenTxListResp model.TokenTxListResp
	tokenTxListResp.Total = tokenTxTotal
	tokenTxListResp.TokenTxList = make([]*model.TokenTxResp, 0)
	token := exp.GetToken(tokenHash)
	for _, tokenTx := range tokenTxs {
		amount := strconv.FormatUint(tokenTx.Amount, 10)
		if token != nil {
			amount = exp.FormatAmount(token.Precision, tokenTx.Amount)
		}
		tokenTxListResp.TokenTxList = append(tokenTxListResp.TokenTxList, &model.TokenTxResp{
			TxHash: tokenTx.TxHash,
			From:   tokenTx.From,
			To:     tokenTx.To,
			Amount: amount,
			Height: tokenTx.Height,
			TT:     tokenTx.TT,
			Direct: tokenTx.Direct,
		})
	}
	return &tokenTxListResp
}

func (exp *Service) outputAddressTxList(addressTxs []*model.AddressTx, addressTxTotal uint32) *model.AddressTxListResp {
	var addressTxListResp model.AddressTxListResp
	addressTxListResp.Total = addressTxTotal
	addressTxListResp.AddressTxList = make([]*model.AddressTxResp, 0)
	for _, addressTx := range addressTxs {
		txresp := &model.AddressTxResp{
			TxHash:    addressTx.TxHash,
			From:      addressTx.From,
			To:        addressTx.To,
			Amount:    strconv.FormatUint(addressTx.Amount, 10),
			Height:    addressTx.Height,
			TT:        addressTx.TT,
			Direct:    addressTx.Direct,
			TokenHash: addressTx.Asset,
		}
		token := exp.GetToken(addressTx.Asset)
		if token != nil {
			txresp.Amount = exp.FormatAmount(token.Precision, addressTx.Amount)
			txresp.TokenName = token.Name
			txresp.TokenType = token.Type
		}
		addressTxListResp.AddressTxList = append(addressTxListResp.AddressTxList, txresp)
	}
	return &addressTxListResp
}
