package database

import "os"

const (
	RDBMagicNumber      = "REDIS"
	RDBVersion          = "\x00\x00\x00\x06"
	RDBDatabaseSelector = "\xFE\x00"
	RDBEOF              = "\xFF"

	RDBStringType = "\x00"
	RDBListType   = "\x01"

	RDBFilename = "dump.rdb"
)

func RDBFileExists() bool {
	_, err := os.Stat(RDBFilename)
	return !os.IsNotExist(err)
}
