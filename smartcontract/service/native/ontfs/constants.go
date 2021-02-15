
package ontfs

const (
	DefaultPerBlockSize = 256 //kb.
)

const (
	DefaultMinTimeForFileStorage = 60 * 60 * 24 //1day
	DefaultContractInvokeGasFee  = 10000000     //0.01ong
	DefaultChallengeReward       = 100000000    //0.1ong
	DefaultFilePerServerPdpTimes = 2
	DefaultPassportExpire        = 9           //block count. passport expire for GetFileHashList
	DefaultNodeMinVolume         = 1024 * 1024 //kb. min total volume with single fsNode
	DefaultChallengeInterval     = 1 * 60 * 60 //1hour
	DefaultNodePerKbPledge       = 1024 * 100  //fsNode's pledge for participant
	DefaultFilePerBlockFeeRate   = 60          //file mode cost of per block save from fsNode for one minute
	DefaultSpacePerBlockFeeRate  = 60          //space mode cost of per block save from fsNode for one hour
	DefaultGasPerBlockForRead    = 256         //cost of per block read from fsNode
)

//challenge state
const (
	Judged = iota
	NoReplyAndValid
	NoReplyAndExpire
	RepliedAndSuccess
	RepliedButVerifyError
)
