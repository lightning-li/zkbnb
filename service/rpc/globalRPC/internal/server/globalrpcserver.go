// Code generated by goctl. DO NOT EDIT!
// Source: globalRPC.proto

package server

import (
	"context"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/internal/logic/sendRawTypeTx"

	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/globalRPCProto"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/internal/logic"
	"github.com/zecrey-labs/zecrey-legend/service/rpc/globalRPC/internal/svc"
)

type GlobalRPCServer struct {
	svcCtx *svc.ServiceContext
	globalRPCProto.UnimplementedGlobalRPCServer
}

func NewGlobalRPCServer(svcCtx *svc.ServiceContext) *GlobalRPCServer {
	return &GlobalRPCServer{
		svcCtx: svcCtx,
	}
}

//  Asset
func (s *GlobalRPCServer) GetLatestAssetsListByAccountIndex(ctx context.Context, in *globalRPCProto.ReqGetLatestAssetsListByAccountIndex) (*globalRPCProto.RespGetLatestAssetsListByAccountIndex, error) {
	l := logic.NewGetLatestAssetsListByAccountIndexLogic(ctx, s.svcCtx)
	return l.GetLatestAssetsListByAccountIndex(in)
}

//  Liquidity
func (s *GlobalRPCServer) GetLatestPairInfo(ctx context.Context, in *globalRPCProto.ReqGetLatestPairInfo) (*globalRPCProto.RespGetLatestPairInfo, error) {
	l := logic.NewGetLatestPairInfoLogic(ctx, s.svcCtx)
	return l.GetLatestPairInfo(in)
}

func (s *GlobalRPCServer) GetSwapAmount(ctx context.Context, in *globalRPCProto.ReqGetSwapAmount) (*globalRPCProto.RespGetSwapAmount, error) {
	l := logic.NewGetSwapAmountLogic(ctx, s.svcCtx)
	return l.GetSwapAmount(in)
}

func (s *GlobalRPCServer) GetLpValue(ctx context.Context, in *globalRPCProto.ReqGetLpValue) (*globalRPCProto.RespGetLpValue, error) {
	l := logic.NewGetLpValueLogic(ctx, s.svcCtx)
	return l.GetLpValue(in)
}

//  Transaction
func (s *GlobalRPCServer) SendTx(ctx context.Context, in *globalRPCProto.ReqSendTx) (*globalRPCProto.RespSendTx, error) {
	l := logic.NewSendTxLogic(ctx, s.svcCtx)
	return l.SendTx(in)
}

func (s *GlobalRPCServer) SendCreateCollectionTx(ctx context.Context, in *globalRPCProto.ReqSendCreateCollectionTx) (*globalRPCProto.RespSendCreateCollectionTx, error) {
	l := logic.NewSendCreateCollectionTxLogic(ctx, s.svcCtx)
	return l.SendCreateCollectionTx(in)
}

func (s *GlobalRPCServer) SendMintNftTx(ctx context.Context, in *globalRPCProto.ReqSendMintNftTx) (*globalRPCProto.RespSendMintNftTx, error) {
	l := logic.NewSendMintNftTxLogic(ctx, s.svcCtx)
	return l.SendMintNftTx(in)
}

func (s *GlobalRPCServer) GetNextNonce(ctx context.Context, in *globalRPCProto.ReqGetNextNonce) (*globalRPCProto.RespGetNextNonce, error) {
	l := logic.NewGetNextNonceLogic(ctx, s.svcCtx)
	return l.GetNextNonce(in)
}

//  NFT
func (s *GlobalRPCServer) GetMaxOfferId(ctx context.Context, in *globalRPCProto.ReqGetMaxOfferId) (*globalRPCProto.RespGetMaxOfferId, error) {
	l := logic.NewGetMaxOfferIdLogic(ctx, s.svcCtx)
	return l.GetMaxOfferId(in)
}

func (s *GlobalRPCServer) SendAddLiquidityTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendAddLiquidityTxLogic(ctx, s.svcCtx)
	return l.SendAddLiquidityTx(in)
}

func (s *GlobalRPCServer) SendAtomicMatchTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendAtomicMatchTxLogic(ctx, s.svcCtx)
	return l.SendAtomicMatchTx(in)
}

func (s *GlobalRPCServer) SendCancelOfferTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendCancelOfferTxLogic(ctx, s.svcCtx)
	return l.SendCancelOfferTx(in)
}

func (s *GlobalRPCServer) SendRemoveLiquidityTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendRemoveLiquidityTxLogic(ctx, s.svcCtx)
	return l.SendRemoveLiquidityTx(in)
}

func (s *GlobalRPCServer) SendSwapTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendSwapTxLogic(ctx, s.svcCtx)
	return l.SendSwapTx(in)
}

func (s *GlobalRPCServer) SendTransferNftTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendTransferNftTxLogic(ctx, s.svcCtx)
	return l.SendTransferNftTx(in)
}

func (s *GlobalRPCServer) SendTransferTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendTransferTxLogic(ctx, s.svcCtx)
	return l.SendTransferTx(in)
}

func (s *GlobalRPCServer) SendWithdrawNftTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendWithdrawNftTxLogic(ctx, s.svcCtx)
	return l.SendWithdrawNftTx(in)
}

func (s *GlobalRPCServer) SendWithdrawTx(ctx context.Context, in *globalRPCProto.ReqSendTxByRawInfo) (*globalRPCProto.RespSendTx, error) {
	l := sendRawTypeTx.NewSendWithdrawTxLogic(ctx, s.svcCtx)
	return l.SendWithdrawTx(in)
}
