package config

// DefaultValues is the default configuration
const DefaultValues = `
[Aggregator]
Host = "0.0.0.0"
Port = 50081
RetryTime = "5s"
VerifyProofInterval = "10s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"
ProofStatePollingInterval = "5s"
SenderAddress = ""
CleanupLockedProofsInterval = "2m"
GeneratingProofCleanupThreshold = "10m"
ForkId = 9
GasOffset = 0
WitnessURL = "localhost:8123"
UseL1BatchData = true
UseFullWitness = false
SettlementBackend = "l1"
AggLayerTxTimeout = "5m"
AggLayerURL = ""
SequencerPrivateKey = {}
	[Aggregator.DB]
		Name = "aggregator_db"
		User = "aggregator_user"
		Password = "aggregator_password"
		Host = "zkevm-db"
		Port = "5432"
		EnableLog = false	
		MaxConns = 200
	[Aggregator.Log]
		Environment = "development" # "production" or "development"
		Level = "info"
		Outputs = ["stderr"]
	[Aggregator.StreamClient]
		Server = "localhost:6900"
	[Aggregator.EthTxManager]
		FrequencyToMonitorTxs = "1s"
		WaitTxToBeMined = "2m"
		GetReceiptMaxTime = "250ms"
		GetReceiptWaitInterval = "1s"
		PrivateKeys = [
			{Path = "/pk/aggregator.keystore", Password = "testonly"},
		]
		ForcedGas = 0
		GasPriceMarginFactor = 1
		MaxGasPriceLimit = 0
		PersistenceFilename = ""
		ReadPendingL1Txs = false
		SafeStatusL1NumberOfBlocks = 0
		FinalizedStatusL1NumberOfBlocks = 0
			[Aggregator.EthTxManager.Etherman]
				URL = ""
				L1ChainID = 11155111
				HTTPHeaders = []
	[Aggregator.Synchronizer]
		[Aggregator.Synchronizer.DB]
			Name = "sync_db"
			User = "sync_user"
			Password = "sync_password"
			Host = "zkevm-db"
			Port = "5432"
			EnableLog = false
			MaxConns = 10
		[Aggregator.Synchronizer.Synchronizer]
			SyncInterval = "10s"
			SyncChunkSize = 1000
			GenesisBlockNumber = 5511080
			SyncUpToBlock = "finalized"
			BlockFinality = "finalized"
			OverrideStorageCheck = false
		[Aggregator.Synchronizer.Etherman]
			[Aggregator.Synchronizer.Etherman.Validium]
				Enabled = false
[EventLog]
	[EventLog.DB]
`
