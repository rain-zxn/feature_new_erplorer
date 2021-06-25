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

// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta

package model

import (
	"github.com/polynetwork/explorer/basedef"
)

type ExplorerInfoReq struct {
	Start        string    `json:"start"`
	End          string    `json:"end"`
}

type ExplorerInfoResp struct {
	Chains        []*ChainInfoResp `json:"chains"`
	CrossTxNumber int64           `json:"crosstxnumber"`
	Tokens        []*CrossChainTokenResp `json:"tokens"`
}

func getChainStatistic(chainId uint64, statistics []*ChainStatistic) *ChainStatistic {
	for _, statistic := range statistics {
		if statistic.ChainId == chainId {
			return statistic
		}
	}
	return nil
}

func MakeExplorerInfoResp(chains []*Chain, statistics []*ChainStatistic, tokenBasics []*TokenBasic) *ExplorerInfoResp {
	chainInfoResps := make([]*ChainInfoResp, 0)
	for _, chain := range chains {
		chainInfoResp := MakeChainInfoResp(chain)
		for _, statistic := range statistics {
			if statistic.ChainId == *chain.ChainId {
				chainInfoResp.Addresses = statistic.Addresses
				chainInfoResp.In = statistic.In
				chainInfoResp.Out = statistic.Out
			}
		}
		for _, tokenBasic := range tokenBasics {
			for _, token := range tokenBasic.Tokens {
				if token.ChainId == *chain.ChainId {
					chainInfoResp.Tokens = append(chainInfoResp.Tokens, MakeChainTokenResp(token))
				}
			}
		}
		chainInfoResps = append(chainInfoResps, chainInfoResp)
	}
	crossTxNumber := getChainStatistic(basedef.POLY_CROSSCHAIN_ID, statistics).In
	crossChainTokenResp := make([]*CrossChainTokenResp, 0)
	for _, tokenBasic := range tokenBasics {
		crossChainTokenResp = append(crossChainTokenResp, MakeTokenBasicResp(tokenBasic))
	}
	explorerInfoResp := &ExplorerInfoResp{
		Chains: chainInfoResps,
		CrossTxNumber: crossTxNumber,
		Tokens: crossChainTokenResp,
	}
	return explorerInfoResp
}

type ChainInfoResp struct {
	Id        uint32               `json:"chainid"`
	Name      string               `json:"chainname"`
	Height    uint32               `json:"blockheight"`
	In        int64               `json:"in"`
	//InCrossChainTxStatus []*CrossChainTxStatus    `json:"incrosschaintxstatus"`
	Out       int64               `json:"out"`
	//OutCrossChainTxStatus []*CrossChainTxStatus    `json:"outcrosschaintxstatus"`
	Addresses int64               `json:"addresses"`
	//Contracts []*ChainContractResp `json:"contracts"`
	Tokens    []*ChainTokenResp    `json:"tokens"`
}

func MakeChainInfoResp(chain *Chain) *ChainInfoResp {
	chainInfoResp := &ChainInfoResp{
		Id:        0,
		Name:      "",
		Height:    0,
		In:        0,
		Out:       0,
		Addresses: 0,
		Tokens:    nil,
	}
	return chainInfoResp
}

type CrossChainTxStatus struct {
	TT        uint32    `json:"timestamp"`
	TxNumber  uint32    `json:"txnumber"`
}

type ChainContractResp struct {
	Id       uint32 `json:"chainid"`
	Contract string `json:"contract"`
}

type ChainTokenResp struct {
	Chain       int32  `json:"chainid"`
	ChainName   string    `json:"chainname"`
	Hash        string `json:"hash"`
	Token       string  `json:"token"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Precision   uint64 `json:"precision"`
	Desc        string `json:"desc"`
}

func MakeChainTokenResp(token *Token) *ChainTokenResp {
	chainTokenResp := &ChainTokenResp{

	}
	return chainTokenResp
}

type CrossChainTokenResp struct {
	Name      string             `json:"name"`
	Tokens    []*ChainTokenResp  `json:"tokens"`
}

func MakeTokenBasicResp(tokenBasic *TokenBasic) *CrossChainTokenResp {
	crossChainTokenResp := &CrossChainTokenResp{
		Name: tokenBasic.Name,
	}
	for _, token := range tokenBasic.Tokens {
		crossChainTokenResp.Tokens = append(crossChainTokenResp.Tokens, MakeChainTokenResp(token))
	}
	return crossChainTokenResp
}

type FChainTxResp struct {
	ChainId    uint32    `json:"chainid"`
	ChainName  string    `json:"chainname"`
	TxHash     string    `json:"txhash"`
	State      byte      `json:"state"`
	TT         uint32    `json:"timestamp"`
	Fee        string    `json:"fee"`
	Height     uint32    `json:"blockheight"`
	User       string    `json:"user"`
	TChainId   uint32    `json:"tchainid"`
	TChainName string    `json:"tchainname"`
	Contract   string    `json:"contract"`
	Key        string    `json:"key"`
	Param      string    `json:"param"`
	Transfer   *FChainTransferResp `json:"transfer"`
}

type FChainTransferResp struct {
	TokenHash    string `json:"tokenhash"`
	TokenName    string `json:"tokenname"`
	TokenType    string `json:"tokentype"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	ToChain      uint32 `json:"tchainid"`
	ToChainName  string `json:"tchainname"`
	ToTokenHash  string `json:"totokenhash"`
	ToTokenName  string `json:"totokenname"`
	ToTokenType  string `json:"totokentype"`
	ToUser       string `json:"tuser"`
}

type MChainTxResp struct {
	ChainId    uint32 `json:"chainid"`
	ChainName  string `json:"chainname"`
	TxHash     string `json:"txhash"`
	State      byte   `json:"state"`
	TT         uint32 `json:"timestamp"`
	Fee        string `json:"fee"`
	Height     uint32 `json:"blockheight"`
	FChainId   uint32 `json:"fchainid"`
	FChainName string `json:"fchainname"`
	FTxHash    string `json:"ftxhash"`
	TChainId   uint32 `json:"tchainid"`
	TChainName string `json:"tchainname"`
	Key        string `json:"key"`
}

type TChainTxResp struct {
	ChainId    uint32    `json:"chainid"`
	ChainName  string    `json:"chainname"`
	TxHash     string    `json:"txhash"`
	State      byte      `json:"state"`
	TT         uint32    `json:"timestamp"`
	Fee        string    `json:"fee"`
	Height     uint32    `json:"blockheight"`
	FChainId   uint32    `json:"fchainid"`
	FChainName string    `json:"fchainname"`
	Contract   string    `json:"contract"`
	RTxHash    string    `json:"mtxhash"`
	Transfer   *TChainTransferResp `json:"transfer"`
}

type TChainTransferResp struct {
	TokenHash    string `json:"tokenhash"`
	TokenName    string `json:"tokenname"`
	TokenType    string `json:"tokentype"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
}

type CrossTransferResp struct {
	CrossTxType  uint32 `json:"crosstxtype"`
	CrossTxName  string `json:"crosstxname"`
	FromChainId  uint32 `json:"fromchainid"`
	FromChain    string `json:"fromchainname"`
	FromAddress  string `json:"fromaddress"`
	ToChainId    uint32 `json:"tochainid"`
	ToChain      string `json:"tochainname"`
	ToAddress    string `json:"toaddress"`
	TokenHash    string `json:"tokenhash"`
	TokenName    string `json:"tokenname"`
	TokenType    string `json:"tokentype"`
	Amount       string `json:"amount"`
}

// swagger:parameters CrossTxReq
type CrossTxReq struct {
	// in: query
	TxHash    string       `json:"txhash"`
}

type CrossTxResp struct {
	Transfer       *CrossTransferResp `json:"crosstransfer"`
	Fchaintx       *FChainTxResp      `json:"fchaintx"`
	Fchaintx_valid bool               `json:"fchaintx_valid"`
	Mchaintx       *MChainTxResp      `json:"mchaintx"`
	Mchaintx_valid bool               `json:"mchaintx_valid"`
	Tchaintx       *TChainTxResp      `json:"tchaintx"`
	Tchaintx_valid bool               `json:"tchaintx_valid"`
}

// getcrosstx response
// swagger:response CrossTxResponse
type CrossTxResponse struct {
	// response body
	// in: body
	Body struct {
		Code          int                    `json:"code"`
		Action        string                 `json:"action"`
		Desc          string                 `json:"desc"`
		Version       string                 `json:"version"`
		Result        CrossTxResp            `json:"result"`
	}
}

type CrossTxListReq struct {
	Start        string    `json:"start"`
	End          string    `json:"end"`
}

// swagger:parameters CrossTxListRequest
type CrossTxListRequest struct {
	// in: body
	Body CrossTxListReq
}

type CrossTxOutlineResp struct {
	TxHash     string        `json:"txhash"`
	State      byte          `json:"state"`
	TT         uint32        `json:"timestamp"`
	Fee        uint64        `json:"fee"`
	Height     uint32        `json:"blockheight"`
	FChainId   uint32        `json:"fchainid"`
	FChainName string        `json:"fchainname"`
	TChainId   uint32        `json:"tchainid"`
	TChainName string        `json:"tchainname"`
}

type CrossTxListResp struct {
	CrossTxList       []*CrossTxOutlineResp     `json:"crosstxs"`
}

type TokenTxListReq struct {
	PageSize int
	PageNo   int
	ChainId uint64
	Token       string     `json:"token"`
}

type TokenTxResp struct {
	TxHash       string `json:"txhash"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Amount       string  `json:"amount"`
	TT           uint32   `json:"timestamp"`
	Height       uint32  `json:"blockheight"`
	Direct       uint32  `json:"direct"`
}

type TokenTxListResp struct {
	TokenTxList       []*TokenTxResp     `json:"tokentxs"`
	Total             int64             `json:"total"`
}

func MakeTokenTxList(transactoins []*TransactionOnToken, tokenStatistic *TokenStatistic) *TokenTxListResp {
	tokenTxListResp := &TokenTxListResp {

	}
	tokenTxListResp.Total = tokenStatistic.InCounter + tokenStatistic.OutCounter
	tokenTxListResp.TokenTxList = make([]*TokenTxResp, 0)
	for _, transactoin := range transactoins {
		tokenTxListResp.TokenTxList = append(tokenTxListResp.TokenTxList, &TokenTxResp{
			TxHash: transactoin.Hash,
		})
	}
	return tokenTxListResp
}



type AddressTxListReq struct {
	PageSize int
	PageNo   int
	Address       string     `json:"address"`
	ChainId         string     `json:"chain"`
}

type AddressTxResp struct {
	TxHash       string `json:"txhash"`
	From         string  `json:"from"`
	To           string  `json:"to"`
	Amount       string  `json:"amount"`
	TT           uint32   `json:"timestamp"`
	Height       uint32  `json:"blockheight"`
	TokenHash    string  `json:"tokenhash"`
	TokenName    string  `json:"tokenname"`
	TokenType    string  `json:"tokentype"`
	Direct       uint32  `json:"direct"`
}

type AddressTxListResp struct {
	AddressTxList       []*AddressTxResp     `json:"addresstxs"`
	Total               int64               `json:"total"`
}

func MakeAddressTxList(transactoins []*TransactionOnAddress, counter int64) *AddressTxListResp {
	addressTxListResp := &AddressTxListResp {

	}
	addressTxListResp.Total = counter
	addressTxListResp.AddressTxList = make([]*AddressTxResp, 0)
	for _, transactoin := range transactoins {
		addressTxListResp.AddressTxList = append(addressTxListResp.AddressTxList, &AddressTxResp{
			TxHash: transactoin.Hash,
		})
	}
	return addressTxListResp
}


