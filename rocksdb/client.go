package rocksdb

/*
#cgo CFLAGS: -I/usr/local/include/rocksdb
#cgo LDFLAGS: -L/usr/local/bin -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd

#include <stdlib.h>
#include "rocksdb/c.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

type RocksDb struct {
	db   *C.rocksdb_t
	path string
}

func OpenDb() *RocksDb {
	db := RocksDb{path: "/tmp/rocksdb_c_simple_example"}
	options := C.rocksdb_options_create()
	C.rocksdb_options_set_create_if_missing(options, 1)
	C.rocksdb_options_set_use_direct_io_for_flush_and_compaction(options, 1)
	C.rocksdb_options_set_use_direct_reads(options, 1)

	pathRef := C.CString(db.path)
	var errRef *C.char
	rocksDb := C.rocksdb_open(options, pathRef, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		panic(err)
	}

	defer func() {
		C.free(unsafe.Pointer(pathRef))
		C.free(unsafe.Pointer(errRef))
	}()
	db.db = rocksDb
	return &db
}

func (db *RocksDb) Put(key, value string) error {
	keyRef := C.CString(key)
	valueRef := C.CString(value)
	var errRef *C.char

	options := C.rocksdb_writeoptions_create()
	C.rocksdb_writeoptions_set_sync(options, 1)

	C.rocksdb_put(db.db, options, keyRef, C.size_t(len(key)), valueRef, C.size_t(len(value)), &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		return err
	}

	defer func() {
		C.free(unsafe.Pointer(keyRef))
		C.free(unsafe.Pointer(valueRef))
		C.free(unsafe.Pointer(options))
	}()
	return nil
}

func (db *RocksDb) Get(key string) (string, error) {
	keyRef := C.CString(key)

	options := C.rocksdb_readoptions_create()
	var returnedLen C.size_t
	var errRef *C.char
	result := C.rocksdb_get(db.db, options, keyRef, C.size_t(len(key)), &returnedLen, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		return "", err
	}

	defer func() {
		C.free(unsafe.Pointer(keyRef))
		C.free(unsafe.Pointer(options))
		//C.free(unsafe.Pointer(returnedLen))
		C.free(unsafe.Pointer(errRef))
	}()

	return C.GoString(result), nil
}

func (db *RocksDb) Delete(key string) error {
	keyRef := C.CString(key)

	options := C.rocksdb_writeoptions_create()
	C.rocksdb_writeoptions_set_sync(options, 1)
	var errRef *C.char
	C.rocksdb_delete(db.db, options, keyRef, C.size_t(len(key)), &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		return err
	}

	defer func() {
		C.free(unsafe.Pointer(keyRef))
		C.free(unsafe.Pointer(errRef))
	}()
	return nil
}

func (db *RocksDb) DeleteDb() {
	C.rocksdb_close(db.db)
	options := C.rocksdb_options_create()
	pathRef := C.CString(db.path)
	var errRef *C.char
	C.rocksdb_destroy_db(options, pathRef, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		panic(err)
	}
}
