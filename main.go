
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/Daironode/aingle-crypto/keypair"
	alog "github.com/Daironode/aingle-event/log"
	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/cmd"
	cmdcom " github.com/Daironode/aingle/cmd/common"
	" github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/consensus"
	" github.com/Daironode/aingle/core/genesis"
	" github.com/Daironode/aingle/core/ledger"
	" github.com/Daironode/aingle/events"
	bactor " github.com/Daironode/aingle/http/base/actor"
	" github.com/Daironode/aingle/http/graphql"
	" github.com/Daironode/aingle/http/jsonrpc"
	" github.com/Daironode/aingle/http/localrpc"
	" github.com/Daironode/aingle/http/nodeinfo"
	" github.com/Daironode/aingle/http/restful"
	" github.com/Daironode/aingle/http/websocket"
	" github.com/Daironode/aingle/p2pserver"
	netreqactor " github.com/Daironode/aingle/p2pserver/actor/req"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/txnpool"
	tc " github.com/Daironode/aingle/txnpool/common"
	" github.com/Daironode/aingle/txnpool/proc"
	" github.com/Daironode/aingle/validator/stateful"
	" github.com/Daironode/aingle/validator/stateless"
	"github.com/urfave/cli"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "AIngle CLI"
	app.Action = startAIngle
	app.Version = config.Version
	app.Copyright = "Copyright in 2018 The AIngle Authors"
	app.Commands = []cli.Command{
		cmd.AccountCommand,
		cmd.InfoCommand,
		cmd.AssetCommand,
		cmd.ContractCommand,
		cmd.ImportCommand,
		cmd.ExportCommand,
		cmd.TxCommond,
		cmd.SigTxCommand,
		cmd.MultiSigAddrCommand,
		cmd.MultiSigTxCommand,
		cmd.SendTxCommand,
		cmd.ShowTxCommand,
	}
	app.Flags = []cli.Flag{
		//common setting
		utils.ConfigFlag,
		utils.LogLevelFlag,
		utils.LogDirFlag,
		utils.DisableLogFileFlag,
		utils.DisableEventLogFlag,
		utils.DataDirFlag,
		utils.WasmVerifyMethodFlag,
		//account setting
		utils.WalletFileFlag,
		utils.AccountAddressFlag,
		utils.AccountPassFlag,
		//consensus setting
		utils.EnableConsensusFlag,
		utils.MaxTxInBlockFlag,
		//txpool setting
		utils.GasPriceFlag,
		utils.GasLimitFlag,
		utils.TxpoolPreExecDisableFlag,
		utils.DisableSyncVerifyTxFlag,
		utils.DisableBroadcastNetTxFlag,
		//p2p setting
		utils.ReservedPeersOnlyFlag,
		utils.ReservedPeersFileFlag,
		utils.NetworkIdFlag,
		utils.NodePortFlag,
		utils.HttpInfoPortFlag,
		utils.MaxConnInBoundFlag,
		utils.MaxConnOutBoundFlag,
		utils.MaxConnInBoundForSingleIPFlag,
		//test mode setting
		utils.EnableTestModeFlag,
		utils.TestModeGenBlockTimeFlag,
		//rpc setting
		utils.RPCDisabledFlag,
		utils.RPCPortFlag,
		utils.RPCLocalEnableFlag,
		utils.RPCLocalProtFlag,
		//rest setting
		utils.RestfulEnableFlag,
		utils.RestfulPortFlag,
		utils.RestfulMaxConnsFlag,
		//graphql setting
		utils.GraphQLEnableFlag,
		utils.GraphQLPortFlag,
		utils.GraphQLMaxConnsFlag,
		//ws setting
		utils.WsEnabledFlag,
		utils.WsPortFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}

func startAIngle(ctx *cli.Context) {
	initLog(ctx)

	log.Infof("aingle version %s", config.Version)

	setMaxOpenFiles()

	cfg, err := initConfig(ctx)
	if err != nil {
		log.Errorf("initConfig error: %s", err)
		return
	}
	acc, err := initAccount(ctx)
	if err != nil {
		log.Errorf("initWallet error: %s", err)
		return
	}
	stateHashHeight := config.GetStateHashCheckHeight(cfg.P2PNode.NetworkId)
	ldg, err := initLedger(ctx, stateHashHeight)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	txpool, err := initTxPool(ctx)
	if err != nil {
		log.Errorf("initTxPool error: %s", err)
		return
	}
	p2pSvr, p2p, err := initP2PNode(ctx, txpool, acc)
	if err != nil {
		log.Errorf("initP2PNode error: %s", err)
		return
	}
	_, err = initConsensus(ctx, p2p, txpool, acc)
	if err != nil {
		log.Errorf("initConsensus error: %s", err)
		return
	}
	err = initRpc(ctx)
	if err != nil {
		log.Errorf("initRpc error: %s", err)
		return
	}
	err = initLocalRpc(ctx)
	if err != nil {
		log.Errorf("initLocalRpc error: %s", err)
		return
	}
	initGraphQL(ctx)
	initRestful(ctx)
	initWs(ctx)
	initNodeInfo(ctx, p2pSvr)

	go logCurrBlockHeight()
	waitToExit(ldg)
}

func initLog(ctx *cli.Context) {
	//init log module
	logLevel := ctx.GlobalInt(utils.GetFlagName(utils.LogLevelFlag))
	//if true, the log will not be output to the file
	disableLogFile := ctx.GlobalBool(utils.GetFlagName(utils.DisableLogFileFlag))
	if disableLogFile {
		log.InitLog(logLevel, log.Stdout)
	} else {
		logFileDir := ctx.GlobalString(utils.GetFlagName(utils.LogDirFlag))
		logFileDir = filepath.Join(logFileDir, "") + string(os.PathSeparator)
		alog.InitLog(logFileDir)
		log.InitLog(logLevel, logFileDir, log.Stdout)
	}
}

func initConfig(ctx *cli.Context) (*config.AIngleConfig, error) {
	//init aingle config from cli
	cfg, err := cmd.SetAIngleConfig(ctx)
	if err != nil {
		return nil, err
	}
	log.Infof("Config init success")
	return cfg, nil
}

func initAccount(ctx *cli.Context) (*account.Account, error) {
	if !config.DefConfig.Consensus.EnableConsensus {
		return nil, nil
	}
	walletFile := ctx.GlobalString(utils.GetFlagName(utils.WalletFileFlag))
	if walletFile == "" {
		return nil, fmt.Errorf("please config wallet file using --wallet flag")
	}
	if !common.FileExisted(walletFile) {
		return nil, fmt.Errorf("cannot find wallet file: %s. Please create a wallet first", walletFile)
	}

	acc, err := cmdcom.GetAccount(ctx)
	if err != nil {
		return nil, fmt.Errorf("get account error: %s", err)
	}
	pubKey := hex.EncodeToString(keypair.SerializePublicKey(acc.PublicKey))
	log.Infof("Using account: %s, pubkey: %s", acc.Address.ToBase58(), pubKey)

	if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		config.DefConfig.Genesis.SOLO.Bookkeepers = []string{pubKey}
	}

	log.Infof("Account init success")
	return acc, nil
}

func initLedger(ctx *cli.Context, stateHashHeight uint32) (*ledger.Ledger, error) {
	events.Init() //Init event hub

	var err error
	dbDir := utils.GetStoreDirPath(config.DefConfig.Common.DataDir, config.DefConfig.P2PNode.NetworkName)
	ledger.DefLedger, err = ledger.NewLedger(dbDir, stateHashHeight)
	if err != nil {
		return nil, fmt.Errorf("NewLedger error: %s", err)
	}
	bookKeepers, err := config.DefConfig.GetBookkeepers()
	if err != nil {
		return nil, fmt.Errorf("GetBookkeepers error: %s", err)
	}
	genesisConfig := config.DefConfig.Genesis
	genesisBlock, err := genesis.BuildGenesisBlock(bookKeepers, genesisConfig)
	if err != nil {
		return nil, fmt.Errorf("genesisBlock error %s", err)
	}
	err = ledger.DefLedger.Init(bookKeepers, genesisBlock)
	if err != nil {
		return nil, fmt.Errorf("init ledger error: %s", err)
	}

	log.Infof("Ledger init success")
	return ledger.DefLedger, nil
}

func initTxPool(ctx *cli.Context) (*proc.TXPoolServer, error) {
	disablePreExec := ctx.GlobalBool(utils.GetFlagName(utils.TxpoolPreExecDisableFlag))
	bactor.DisableSyncVerifyTx = ctx.GlobalBool(utils.GetFlagName(utils.DisableSyncVerifyTxFlag))
	disableBroadcastNetTx := ctx.GlobalBool(utils.GetFlagName(utils.DisableBroadcastNetTxFlag))
	txPoolServer, err := txnpool.StartTxnPoolServer(disablePreExec, disableBroadcastNetTx)
	if err != nil {
		return nil, fmt.Errorf("init txpool error: %s", err)
	}
	stlValidator, _ := stateless.NewValidator("stateless_validator")
	stlValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator2, _ := stateless.NewValidator("stateless_validator2")
	stlValidator2.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stfValidator, _ := stateful.NewValidator("stateful_validator")
	stfValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))

	bactor.SetTxnPoolPid(txPoolServer.GetPID(tc.TxPoolActor))
	bactor.SetTxPid(txPoolServer.GetPID(tc.TxActor))

	log.Infof("TxPool init success")
	return txPoolServer, nil
}

func initP2PNode(ctx *cli.Context, txpoolSvr *proc.TXPoolServer, acct *account.Account) (*p2pserver.P2PServer, p2p.P2P, error) {
	if config.DefConfig.Genesis.ConsensusType == config.CONSENSUS_TYPE_SOLO {
		return nil, nil, nil
	}
	p2p, err := p2pserver.NewServer(acct)
	if err != nil {
		return nil, nil, err
	}

	err = p2p.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2p service start error %s", err)
	}
	netreqactor.SetTxnPoolPid(txpoolSvr.GetPID(tc.TxActor))
	txpoolSvr.Net = p2p.GetNetwork()
	bactor.SetNetServer(p2p.GetNetwork())
	p2p.WaitForPeersStart()
	log.Infof("P2P init success")
	return p2p, p2p.GetNetwork(), nil
}

func initConsensus(ctx *cli.Context, net p2p.P2P, txpoolSvr *proc.TXPoolServer, acc *account.Account) (consensus.ConsensusService, error) {
	if !config.DefConfig.Consensus.EnableConsensus {
		return nil, nil
	}
	pool := txpoolSvr.GetPID(tc.TxPoolActor)

	consensusType := strings.ToLower(config.DefConfig.Genesis.ConsensusType)
	consensusService, err := consensus.NewConsensusService(consensusType, acc, pool, nil, net)
	if err != nil {
		return nil, fmt.Errorf("NewConsensusService %s error: %s", consensusType, err)
	}
	consensusService.Start()

	netreqactor.SetConsensusPid(consensusService.GetPID())
	bactor.SetConsensusPid(consensusService.GetPID())

	log.Infof("Consensus init success")
	return consensusService, nil
}

func initRpc(ctx *cli.Context) error {
	if !config.DefConfig.Rpc.EnableHttpJsonRpc {
		return nil
	}
	var err error
	exitCh := make(chan interface{}, 0)
	go func() {
		err = jsonrpc.StartRPCServer()
		close(exitCh)
	}()

	flag := false
	select {
	case <-exitCh:
		if !flag {
			return err
		}
	case <-time.After(time.Millisecond * 5):
		flag = true
	}
	log.Infof("Rpc init success")
	return nil
}

func initLocalRpc(ctx *cli.Context) error {
	if !ctx.GlobalBool(utils.GetFlagName(utils.RPCLocalEnableFlag)) {
		return nil
	}
	var err error
	exitCh := make(chan interface{}, 0)
	go func() {
		err = localrpc.StartLocalServer()
		close(exitCh)
	}()

	flag := false
	select {
	case <-exitCh:
		if !flag {
			return err
		}
	case <-time.After(time.Millisecond * 5):
		flag = true
	}

	log.Infof("Local rpc init success")
	return nil
}

func initGraphQL(ctx *cli.Context) {
	if !config.DefConfig.GraphQL.EnableGraphQL {
		return
	}
	go graphql.StartServer(config.DefConfig.GraphQL)

	log.Infof("GraphQL init success")
}

func initRestful(ctx *cli.Context) {
	if !config.DefConfig.Restful.EnableHttpRestful {
		return
	}
	go restful.StartServer()

	log.Infof("Restful init success")
}

func initWs(ctx *cli.Context) {
	if !config.DefConfig.Ws.EnableHttpWs {
		return
	}
	websocket.StartServer()

	log.Infof("Ws init success")
}

func initNodeInfo(ctx *cli.Context, p2pSvr *p2pserver.P2PServer) {
	// testmode has no p2pserver(see function initP2PNode for detail), simply ignore httpInfoPort in testmode
	if ctx.Bool(utils.GetFlagName(utils.EnableTestModeFlag)) || config.DefConfig.P2PNode.HttpInfoPort == 0 {
		return
	}
	go nodeinfo.StartServer(p2pSvr.GetNetwork())

	log.Infof("Nodeinfo init success")
}

func logCurrBlockHeight() {
	ticker := time.NewTicker(config.DEFAULT_GEN_BLOCK_TIME * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Infof("CurrentBlockHeight = %d", ledger.DefLedger.GetCurrentBlockHeight())
			log.CheckRotateLogFile()
		}
	}
}

func setMaxOpenFiles() {
	max, err := fdlimit.Maximum()
	if err != nil {
		log.Errorf("failed to get maximum open files: %v", err)
		return
	}
	_, err = fdlimit.Raise(uint64(max))
	if err != nil {
		log.Errorf("failed to set maximum open files: %v", err)
		return
	}
}

func waitToExit(db *ledger.Ledger) {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("AIngle received exit signal: %v.", sig.String())
			log.Infof("closing ledger...")
			db.Close()
			close(exit)
			break
		}
	}()
	<-exit
}
