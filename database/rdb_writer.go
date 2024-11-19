package database

import (
	"fmt"
	"os"
)

type RDBWriter struct {
	db *DB
}

func NewRDBWriter(db *DB) *RDBWriter {
	return &RDBWriter{db: db}
}

func (r *RDBWriter) Write() error {
	// https://rdb.fnordig.de/file_format.html#redis-rdb-file-format
	file, err := os.Create("dump.rdb")
	if err != nil {
		return err
	}
	defer file.Close()

	db.mu.RLock()
	defer db.mu.RUnlock()

	// Magic number
	_, err = file.Write([]byte("REDIS"))
	if err != nil {
		return err
	}
	// RDB Version number
	_, err = file.Write([]byte{0, 0, 0, 6})
	if err != nil {
		return err
	}
	// Database 0 selector
	_, err = file.Write([]byte{0xFE, 00})
	if err != nil {
		return err
	}
	// Key-value pairs
	for key, value := range db.data {
		// Encoding the Key
		err = rdbWriteString(key, file)
		if err != nil {
			return err
		}

		// Encoding the Value
		switch v := value.(type) {
		case dbstring:
			err = rdbWriteString(v.value, file)
			if err != nil {
				return err
			}
		case *dblist:
			err = rdbWriteList(v, file)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported value type %T", v)
		}
	}
	// End of the RDB file
	_, err = file.Write([]byte{0xFF})
	if err != nil {
		return err
	}
	// CRC64 checksum
	// TODO implement CRC64 checksum
	// Disable for now
	_, err = file.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}
	return nil
}

func rdbWriteString(s string, f *os.File) error {
	// Encoded as a length-prefixed string
	// https://rdb.fnordig.de/file_format.html#encoding-strings
	_, err := f.Write([]byte{0}) // String type
	if err != nil {
		return err
	}
	// Length Prefixed String for the key
	// TODO handle keys with length > 63
	_, err = f.Write([]byte{byte(len(s))})
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

func rdbWriteList(l *dblist, f *os.File) error {
	// Encoded as a list
	// https://rdb.fnordig.de/file_format.html#list-encoding
	_, err := f.Write([]byte{1}) // List type
	if err != nil {
		return err
	}
	listSize := 0
	for n := l.head; n != nil; n = n.next {
		listSize++
	}
	// TODO handle lists with length > 63
	_, err = f.Write([]byte{byte(listSize)})
	if err != nil {
		return err
	}
	for n := l.head; n != nil; n = n.next {
		// Length Prefixed String for the list element
		// TODO handle list values with length > 63
		_, err = f.Write([]byte{byte(len(n.value))})
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(n.value))
		if err != nil {
			return err
		}
	}
	return nil
}
