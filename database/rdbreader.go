package database

import (
	"fmt"
	"io"
	"os"
)

type RDBReader struct {
	db *DB
}

func NewRDBReader(db *DB) *RDBReader {
	return &RDBReader{db: db}
}

func (r *RDBReader) Read() error {
	// https://rdb.fnordig.de/file_format.html#redis-rdb-file-format
	rdb, err := os.Open(RDBFilename)
	if err != nil {
		return err
	}
	// Magic number
	magic := make([]byte, len(RDBMagicNumber))
	_, err = io.ReadFull(rdb, magic)
	if err != nil {
		return err
	}
	if string(magic) != RDBMagicNumber {
		return fmt.Errorf("invalid RDB file format")
	}

	// RDB Version number
	version := make([]byte, len(RDBVersion))
	_, err = io.ReadFull(rdb, version)
	if err != nil {
		return err
	}
	if string(version) != RDBVersion {
		return fmt.Errorf("invalid RDB version number")
	}

	// Database 0 selector
	selector := make([]byte, len(RDBDatabaseSelector))
	_, err = io.ReadFull(rdb, selector)
	if err != nil {
		return err
	}
	if string(selector) != RDBDatabaseSelector {
		return fmt.Errorf("invalid database selector")
	}

	// Key-value pairs
	for {
		key, err := rdbReadString(rdb)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		valueType := make([]byte, 1)
		_, err = io.ReadFull(rdb, valueType)
		if err != nil {
			return err
		}

		switch string(valueType[0]) {
		case RDBStringType:
			value, err := rdbReadString(rdb)
			if err != nil {
				return err
			}
			r.db.Set(key, value, nil)
		case RDBListType:
			list, err := rdbReadList(rdb)
			if err != nil {
				return err
			}
			for _, value := range list {
				r.db.ListLPush(key, value)
			}
		default:
			return fmt.Errorf("unsupported value type %d", valueType[0])
		}
	}

	return nil
}

func rdbReadList(rdb *os.File) ([]string, error) {
	length := make([]byte, 1)
	_, err := io.ReadFull(rdb, length)
	if err != nil {
		return nil, err
	}

	// TODO support lengths > 63
	data := make([]string, int(length[0]))
	for i := 0; i < int(length[0]); i++ {
		value, err := rdbReadString(rdb)
		if err != nil {
			return nil, err
		}
		data[i] = value
	}
	return data, nil
}

func rdbReadString(rdb *os.File) (string, error) {
	length := make([]byte, 1)
	_, err := io.ReadFull(rdb, length)
	if err != nil {
		return "", err
	}

	// TODO support lengths > 63
	data := make([]byte, int(length[0]))
	_, err = io.ReadFull(rdb, data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
