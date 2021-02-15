
package sigsvr

import (
	"fmt"

	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/cmd"
	" github.com/Daironode/aingle/cmd/sigsvr/store"
	" github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	"github.com/urfave/cli"
)

var ImportWalletCommand = cli.Command{
	Name:      "import",
	Usage:     "Import accounts from a wallet file",
	ArgsUsage: "",
	Action:    importWallet,
	Flags: []cli.Flag{
		utils.CliWalletDirFlag,
		utils.WalletFileFlag,
	},
	Description: "",
}

func importWallet(ctx *cli.Context) error {
	walletDirPath := ctx.String(utils.GetFlagName(utils.CliWalletDirFlag))
	walletFilePath := ctx.String(utils.GetFlagName(utils.WalletFileFlag))
	if walletDirPath == "" || walletFilePath == "" {
		cmd.PrintErrorMsg("Missing %s or %s flag.", utils.CliWalletDirFlag.Name, utils.WalletFileFlag.Name)
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	if !common.FileExisted(walletFilePath) {
		return fmt.Errorf("wallet file:%s does not exist", walletFilePath)
	}
	walletStore, err := store.NewWalletStore(walletDirPath)
	if err != nil {
		return fmt.Errorf("NewWalletStore dir path:%s error:%s", walletDirPath, err)
	}
	wallet, err := account.Open(walletFilePath)
	if err != nil {
		return fmt.Errorf("open wallet:%s error:%s", walletFilePath, err)
	}
	walletData := wallet.GetWalletData()
	if *walletStore.WalletScrypt != *walletData.Scrypt {
		return fmt.Errorf("import account failed, wallet scrypt:%+v != %+v", walletData.Scrypt, walletStore.WalletScrypt)
	}
	addNum := 0
	updateNum := 0
	for i := 0; i < len(walletData.Accounts); i++ {
		ok, err := walletStore.AddAccountData(walletData.Accounts[i])
		if err != nil {
			return fmt.Errorf("import account address:%s error:%s", walletData.Accounts[i].Address, err)
		}
		if ok {
			addNum++
		} else {
			updateNum++
		}
	}
	cmd.PrintInfoMsg("Import account success.")
	cmd.PrintInfoMsg("Total account number:%d", len(walletData.Accounts))
	cmd.PrintInfoMsg("Add account number:%d", addNum)
	cmd.PrintInfoMsg("Update account number:%d", updateNum)
	return nil
}
