package database

import (
	"fmt"
	"log"
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
	file, err := os.Create(RDBFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	db.mu.RLock()
	defer db.mu.RUnlock()

	// Magic number
	_, err = file.Write([]byte(RDBMagicNumber))
	if err != nil {
		return err
	}
	// RDB Version number
	_, err = file.Write([]byte(RDBVersion))
	if err != nil {
		return err
	}
	// Database 0 selector
	_, err = file.Write([]byte(RDBDatabaseSelector))
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
			err = rdbWriteStringValue(v.value, file)
			if err != nil {
				return err
			}
		case *dblist:
			err = rdbWriteListValue(v, file)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported value type %T", v)
		}
	}
	// End of the RDB file
	_, err = file.Write([]byte(RDBEOF))
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

func rdbWriteStringValue(s string, f *os.File) error {
	_, err := f.Write([]byte(RDBStringType)) // String type
	if err != nil {
		return err
	}
	return rdbWriteString(s, f)
}

func rdbWriteString(s string, f *os.File) error {
	// Length Prefixed String for the key
	// TODO handle keys with length > 63
	_, err := f.Write([]byte{byte(len(s))})
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

func rdbWriteListValue(l *dblist, f *os.File) error {
	// Encoded as a list
	_, err := f.Write([]byte(RDBListType)) // List type
	if err != nil {
		return err
	}
	listSize := 0
	for n := l.head; n != nil; n = n.next {
		listSize++
	}
	log.Printf("List size: %d", listSize)
	// TODO handle lists with length > 63
	_, err = f.Write([]byte{byte(listSize)})
	if err != nil {
		return err
	}
	for n := l.head; n != nil; n = n.next {
		// Length Prefixed String for the list element
		err = rdbWriteString(n.value, f)
		if err != nil {
			return err
		}
	}
	return nil
}
