/*
 * Copyright © 2021 Zecrey Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logic

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zecrey-labs/zecrey-eth-rpc/_rpc"
	"github.com/zecrey-labs/zecrey-legend/common/model/account"
	"github.com/zecrey-labs/zecrey-legend/common/model/nft"
	"github.com/zecrey-labs/zecrey-legend/common/model/proofSender"
	"github.com/zeromicro/go-zero/core/logx"
	"sort"
	"time"
)

/*
	MonitorL2BlockEvents: monitor layer-2 block events
*/
func MonitorL2BlockEvents(
	bscCli *_rpc.ProviderClient,
	bscPendingBlocksCount uint64,
	mempoolModel MempoolModel,
	accountModel AccountModel,
	accountHistoryModel AccountHistoryModel,
	nftModel L2NftModel,
	nftHistoryModel L2NftHistoryModel,
	blockModel BlockModel,
	l1TxSenderModel L1TxSenderModel,
) (err error) {
	// get pending transactions from l1TxSender
	pendingSenders, err := l1TxSenderModel.GetL1TxSendersByTxStatus(L1TxSenderPendingStatus)
	if err != nil {
		logx.Errorf("[MonitorL2BlockEvents] unable to get l1 tx senders by tx status: %s", err.Error())
		return err
	}
	// scan each event
	var (
		// pending update blocks
		relatedBlocks                  = make(map[int64]*Block)
		pendingUpdateSenders           []*L1TxSender
		pendingUpdateProofSenderStatus = make(map[int64]int)
	)
	// handle each sender
	for _, pendingSender := range pendingSenders {
		txHash := pendingSender.L1TxHash
		var (
			l1BlockNumber uint64
			receipt       *types.Receipt
		)
		// check if the status of tx is success
		// get latest l1 block height(latest height - pendingBlocksCount)
		latestHeight, err := bscCli.GetHeight()
		if err != nil {
			errInfo := fmt.Sprintf("[MonitorL2BlockEvents]<=>[cli.GetHeight] %s", err.Error())
			logx.Error(errInfo)
			return err
		}
		_, isPending, err := bscCli.GetTransactionByHash(txHash)
		if err != nil {
			logx.Errorf("[MonitorL2BlockEvents] unable to get transaction by hash: %s", err.Error())
			continue
		}
		if isPending {
			logx.Errorf("[MonitorL2BlockEvents] the tx is still pending, just handle next sender: %s", txHash)
			continue
		}
		// get receipt
		receipt, err = bscCli.GetTransactionReceipt(txHash)
		if err != nil {
			logx.Errorf("[MonitorL2BlockEvents] unable to get tx receipt: %s", err.Error())
			continue
		}
		l1BlockNumber = receipt.BlockNumber.Uint64()
		// check if the height is over safe height
		if latestHeight < l1BlockNumber+bscPendingBlocksCount {
			logx.Infof("[MonitorL2BlockEvents] haven't reached to safe block height, should wait: %s", txHash)
			continue
		}
		// get events from the tx
		logs := receipt.Logs
		timeAt := time.Now().UnixMilli()
		var isValidSender bool
		for _, vlog := range logs {
			switch vlog.Topics[0].Hex() {
			case ZecreyLogBlockCommitSigHash.Hex():
				// parse event info
				var event ZecreyLegendBlockCommit
				err = ZecreyContractAbi.UnpackIntoInterface(&event, BlockCommitEventName, vlog.Data)
				if err != nil {
					errInfo := fmt.Sprintf("[MonitorL2BlockEvents]<=>[ZecreyContractAbi.UnpackIntoInterface] %s", err.Error())
					logx.Error(errInfo)
					return err
				}
				// get related blocks
				blockHeight := int64(event.BlockNumber)
				if relatedBlocks[blockHeight] == nil {
					relatedBlocks[blockHeight], err = blockModel.GetBlockByBlockHeightWithoutTx(blockHeight)
					if err != nil {
						logx.Errorf("[MonitorL2BlockEvents] unable to get block by block height: %s", err.Error())
						return err
					}
				}
				// check block height
				if blockHeight == pendingSender.L2BlockHeight {
					isValidSender = true
				}
				relatedBlocks[blockHeight].CommittedTxHash = receipt.TxHash.Hex()
				relatedBlocks[blockHeight].CommittedAt = timeAt
				break
			case ZecreyLogBlockVerificationSigHash.Hex():
				// parse event info
				var event ZecreyLegendBlockVerification
				err = ZecreyContractAbi.UnpackIntoInterface(&event, BlockVerificationEventName, vlog.Data)
				if err != nil {
					errInfo := fmt.Sprintf("[blockMoniter.MonitorL2BlockEvents]<=>[ZecreyContractAbi.UnpackIntoInterface] %s", err.Error())
					logx.Error(errInfo)
					return err
				}
				// get related blocks
				blockHeight := int64(event.BlockNumber)
				if relatedBlocks[blockHeight] == nil {
					relatedBlocks[blockHeight], err = blockModel.GetBlockByBlockHeightWithoutTx(blockHeight)
					if err != nil {
						logx.Errorf("[MonitorL2BlockEvents] unable to get block by block height: %s", err.Error())
						return err
					}
				}
				// check block height
				if blockHeight == pendingSender.L2BlockHeight {
					isValidSender = true
					pendingUpdateProofSenderStatus[blockHeight] = proofSender.ConfirmedOnChain
				}
				// update block status
				relatedBlocks[blockHeight].VerifiedTxHash = receipt.TxHash.Hex()
				relatedBlocks[blockHeight].VerifiedAt = timeAt
				break
			case ZecreyLogBlocksRevertSigHash.Hex():
				// TODO revert
				break
			default:
				break
			}
		}
		if isValidSender {
			// update sender status
			pendingSender.TxStatus = L1TxSenderHandledStatus
			pendingUpdateSenders = append(pendingUpdateSenders, pendingSender)
		}
	}
	// get pending update info
	var (
		pendingUpdateBlocks []*Block
	)
	for _, pendingUpdateBlock := range relatedBlocks {
		pendingUpdateBlocks = append(pendingUpdateBlocks, pendingUpdateBlock)
	}
	// sort for blocks
	if len(pendingUpdateBlocks) != 0 {
		sort.Sort(blockInfosByBlockHeight(pendingUpdateBlocks))
		logx.Info("pending update blocks count: %v and height: %v", len(pendingUpdateBlocks), pendingUpdateBlocks[len(pendingUpdateBlocks)-1].BlockHeight)
	}

	// handle executed blocks
	var (
		pendingUpdateAccountsMap                          = make(map[int64]*account.Account)
		pendingUpdateNftAssetsMap, pendingNewNftAssetsMap = make(map[int64]*nft.L2Nft), make(map[int64]*nft.L2Nft)
		pendingUpdateMempoolTxs                           []*MempoolTx
	)
	for _, pendingUpdateBlock := range pendingUpdateBlocks {
		if pendingUpdateBlock.BlockStatus == BlockVerifiedStatus {
			pendingUpdateAccountHistories, err := accountHistoryModel.GetAccountsByBlockHeight(pendingUpdateBlock.BlockHeight)
			if err != nil {
				logx.Errorf("[MonitorL2BlockEvents] unable to get related account info by height: %s", err.Error())
				return err
			}
			for _, pendingUpdateAccountHistory := range pendingUpdateAccountHistories {
				// get account info by index
				if pendingUpdateAccountsMap[pendingUpdateAccountHistory.AccountIndex] == nil {
					pendingUpdateAccountsMap[pendingUpdateAccountHistory.AccountIndex], err = accountModel.GetAccountByAccountIndex(pendingUpdateAccountHistory.AccountIndex)
					if err != nil {
						logx.Errorf("[MonitorL2BlockEvents] invalid account index: %s", err.Error())
						return err
					}
				}
				pendingUpdateAccountsMap[pendingUpdateAccountHistory.AccountIndex].Nonce = pendingUpdateAccountHistory.Nonce
				pendingUpdateAccountsMap[pendingUpdateAccountHistory.AccountIndex].AssetInfo = pendingUpdateAccountHistory.AssetInfo
				pendingUpdateAccountsMap[pendingUpdateAccountHistory.AccountIndex].AssetRoot = pendingUpdateAccountHistory.AssetRoot
			}
			// get related account nft from account nft history table
			_, pendingUpdateNftAssetsHistory, err := nftHistoryModel.GetNftAssetsByBlockHeight(pendingUpdateBlock.BlockHeight)
			if err != nil {
				errInfo := fmt.Sprintf("[MonitorL2BlockEvents] unable to get related account liquidity assets by height: %s", err.Error())
				logx.Error(errInfo)
				return err
			}
			// get pending nft assets
			for _, pendingUpdateNftAssetHistory := range pendingUpdateNftAssetsHistory {
				var pendingUpdateNftAsset *nft.L2Nft
				if pendingUpdateNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex] == nil && pendingNewNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex] == nil {
					pendingUpdateNftAsset, err = nftModel.GetNftAsset(pendingUpdateNftAssetHistory.NftIndex)
					if err == ErrNotFound {
						pendingNewNftAsset := &nft.L2Nft{
							NftIndex:            pendingUpdateNftAsset.NftIndex,
							CreatorAccountIndex: pendingUpdateNftAsset.CreatorAccountIndex,
							OwnerAccountIndex:   pendingUpdateNftAsset.OwnerAccountIndex,
							NftContentHash:      pendingUpdateNftAsset.NftContentHash,
							NftL1TokenId:        pendingUpdateNftAsset.NftL1TokenId,
							NftL1Address:        pendingUpdateNftAsset.NftL1Address,
							CollectionId:        pendingUpdateNftAsset.CollectionId,
						}
						pendingNewNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex] = pendingNewNftAsset
						continue
					} else if err != nil {
						errInfo := fmt.Sprintf("[MonitorL2BlockEvents] unable to get related account asset: %s", err.Error())
						logx.Error(errInfo)
						return err
					}
					pendingUpdateNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex] = pendingUpdateNftAsset
				} else {
					if pendingUpdateNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex] == nil {
						pendingUpdateNftAsset = pendingNewNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex]
					} else {
						pendingUpdateNftAsset = pendingUpdateNftAssetsMap[pendingUpdateNftAssetHistory.NftIndex]
					}
				}
				pendingUpdateNftAsset = &nft.L2Nft{
					Model:               pendingUpdateNftAsset.Model,
					NftIndex:            pendingUpdateNftAsset.NftIndex,
					CreatorAccountIndex: pendingUpdateNftAssetHistory.CreatorAccountIndex,
					OwnerAccountIndex:   pendingUpdateNftAssetHistory.OwnerAccountIndex,
					NftContentHash:      pendingUpdateNftAssetHistory.NftContentHash,
					NftL1TokenId:        pendingUpdateNftAssetHistory.NftL1TokenId,
					NftL1Address:        pendingUpdateNftAssetHistory.NftL1Address,
					CollectionId:        pendingUpdateNftAssetHistory.CollectionId,
				}
			}
			// delete related mempool txs
			rowsAffected, pendingDeleteMempoolTxs, err := mempoolModel.GetMempoolTxsByBlockHeight(pendingUpdateBlock.BlockHeight)
			if err != nil {
				errInfo := fmt.Sprintf("[MonitorL2BlockEvents] unable to get related mempool txs by height: %s", err.Error())
				logx.Error(errInfo)
				return err
			}
			if rowsAffected == 0 {
				logx.Error("[MonitorL2BlockEvents] invalid txs size or mempool has been deleted")
				continue
			}
			pendingUpdateMempoolTxs = append(pendingUpdateMempoolTxs, pendingDeleteMempoolTxs...)
		}
	}
	var (
		pendingUpdateAccounts                       []*account.Account
		pendingUpdateNftAssets, pendingNewNftAssets []*nft.L2Nft
	)
	for _, pendingUpdateAccount := range pendingUpdateAccountsMap {
		pendingUpdateAccounts = append(pendingUpdateAccounts, pendingUpdateAccount)
	}
	for _, pendingUpdateNftAsset := range pendingUpdateNftAssetsMap {
		pendingUpdateNftAssets = append(pendingUpdateNftAssets, pendingUpdateNftAsset)
	}
	for _, pendingNewNftAsset := range pendingNewNftAssetsMap {
		pendingNewNftAssets = append(pendingNewNftAssets, pendingNewNftAsset)
	}
	// update blocks, blockDetails, updateEvents, sender
	// update assets, locked assets, liquidity
	// delete mempool txs
	err = l1TxSenderModel.UpdateRelatedEventsAndResetRelatedAssetsAndTxs(
		pendingUpdateBlocks,
		pendingUpdateSenders,
		pendingUpdateAccounts,
		pendingUpdateNftAssets, pendingNewNftAssets,
		pendingUpdateMempoolTxs,
		pendingUpdateProofSenderStatus,
	)
	logx.Info("[MonitorL2BlockEvents] update blocks count: %v", len(pendingUpdateBlocks))
	logx.Info("[MonitorL2BlockEvents] update senders count: %v", len(pendingUpdateSenders))
	logx.Info("[MonitorL2BlockEvents] update accounts count: %v", len(pendingUpdateAccounts))
	logx.Info("[MonitorL2BlockEvents] update nft assets count: %v", len(pendingUpdateNftAssets))
	logx.Info("[MonitorL2BlockEvents] new nft assets count: %v", len(pendingNewNftAssets))
	logx.Info("[MonitorL2BlockEvents] update mempool txs count: %v", len(pendingUpdateMempoolTxs))
	if err != nil {
		logx.Errorf("[MonitorL2BlockEvents] unable to update everything: %s", err.Error())
		return err
	}
	return nil
}
