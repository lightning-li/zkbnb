package desertexit

import (
	"github.com/bnb-chain/zkbnb/tools/desertexit/config"
	"github.com/bnb-chain/zkbnb/tools/desertexit/desertexit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"math/big"
	"strconv"
)

const CommandActivateDesert = "activateDesert"
const CommandPerform = "perform"
const CommandCancelOutstandingDeposit = "cancelOutstandingDeposit"
const CommandWithdrawNFT = "withdrawNFT"
const CommandWithdrawAsset = "withdrawAsset"
const CommandGetBalance = "getBalance"
const CommandGetPendingBalance = "getPendingBalance"
const CommandGetNftList = "getNftList"

func Perform(configFile string, command string, amount string, nftIndex string, owner string, privateKey string, proof string, token string) error {
	var c config.Config
	conf.MustLoad(configFile, &c)
	logx.MustSetup(c.LogConf)
	logx.DisableStat()

	if privateKey != "" {
		c.ChainConfig.PrivateKey = privateKey
	}
	if owner != "" {
		c.Address = owner
	}
	if token != "" {
		c.Token = token
	}
	if nftIndex != "" {
		var err error
		c.NftIndex, err = strconv.ParseInt(nftIndex, 10, 64)
		if err != nil {
			logx.Severe(err)
			return err
		}
	}
	m, err := desertexit.NewPerformDesert(c)
	if err != nil {
		logx.Severe(err)
		return err
	}

	proc.SetTimeToForceQuit(GracefulShutdownTimeout)
	proc.AddShutdownListener(func() {
		logx.Info("start to shutdown desertexit......")
		_ = logx.Close()
	})

	switch command {
	case CommandActivateDesert:
		err = m.ActivateDesertMode()
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandPerform:
		var performDesertAsset desertexit.PerformDesertAssetData
		conf.MustLoad(proof, &performDesertAsset)
		err = m.PerformDesert(performDesertAsset)
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandCancelOutstandingDeposit:
		err = m.CancelOutstandingDeposit()
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandWithdrawNFT:
		err = m.WithdrawPendingNFTBalance(c.NftIndex)
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandWithdrawAsset:
		bigIntAmount, success := new(big.Int).SetString(amount, 10)
		if !success {
			logx.Severe("failed to transfer big int")
			return nil
		}
		err = m.WithdrawPendingBalance(common.HexToAddress(c.Address), common.HexToAddress(c.Token), bigIntAmount)
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandGetBalance:
		_, err := m.GetBalance(common.HexToAddress(c.Address), common.HexToAddress(c.Token))
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandGetPendingBalance:
		_, err := m.GetPendingBalance(common.HexToAddress(c.Address), common.HexToAddress(c.Token))
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	case CommandGetNftList:
		err := m.GetNftList(c.Address)
		if err != nil {
			logx.Severe(err)
			return err
		}
		break
	default:
		logx.Severef("no %s  command", command)
	}
	return nil
}
