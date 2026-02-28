/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/spf13/cobra"
)

// BackupSqliteCmd represents the version command
var BackupSqliteCmd = &cobra.Command{
	Use:   "backup-sqlite",
	Short: "Backup the SQLite database",
	Run: func(cmd *cobra.Command, args []string) {
		Init()
		defer Release()

		// 检查数据库类型
		if conf.Conf.Database.Type != "sqlite3" {
			fmt.Printf("Database type '%s' is not supported for backup. Only sqlite3 is supported.\n", conf.Conf.Database.Type)
			os.Exit(1)
		}

		// 检查数据库文件是否存在
		if !utils.Exists(conf.Conf.Database.DBFile) {
			fmt.Printf("Database file '%s' does not exist.\n", conf.Conf.Database.DBFile)
			os.Exit(1)
		}

		// 创建备份文件名（包含时间戳）
		timestamp := time.Now().Format("20060102_150405")
		backupDir := filepath.Dir(conf.Conf.Database.DBFile)
		backupFileName := fmt.Sprintf("data_backup_%s.db", timestamp)
		backupFilePath := filepath.Join(backupDir, backupFileName)

		database := conf.Conf.Database

		sqldb, err := sql.Open("sqlite3", fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
			database.DBFile))
		if err != nil {
			fmt.Printf("Failed to open database: %s\n", err.Error())
			os.Exit(1)
		}
		defer sqldb.Close()

		_, err = sqldb.Exec(fmt.Sprintf("VACUUM INTO '%s'", backupFilePath))
		if err != nil {
			fmt.Printf("Failed to backup database: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Database backup successful: %s\n", backupFilePath)

	},
}

func init() {
	RootCmd.AddCommand(BackupSqliteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
