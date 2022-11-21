package src

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"sync"
)

type Db struct{
	name string
	dir string
	dbFile *os.File
	keyvalue map[string]string
	mu sync.Mutex
}

const (
	baseDir       = "/tmp/Db/"
	fileExtension = ".db"
)

// open creates data file for newly creating database. If the database file is
// already exists, it returns error without creating anything. name indicates
// database name.
func(k *Db) Set(key, value string) {
	k.mu.Lock()
	k.keyvalue[key] = value
	k.mu.Unlock()
}
func (k *Db) Get(key string) string {
	k.mu.Lock()
	val := k.keyvalue[key]
	// fmt.Println(k.keyvalue)
	k.mu.Unlock()
	return val
}
func Open(name string) (*Db, error) {
	fullPath := baseDir + name + fileExtension

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("database directory couldn't created: %s", err.Error())
		}
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		dbFile, err := os.OpenFile(fullPath, os.O_CREATE, 0777)
			if err != nil {
				return nil, fmt.Errorf("database couldn't created: %s", err.Error())
			}
			m := make(map[string]string)
			k := &Db{
				name:     name + fileExtension,
				dir:      baseDir + name + fileExtension,
				dbFile:   dbFile,
				keyvalue:       m,
			}
			return k, nil
		}
	k, err := OpenAndLoad(name)
	if err != nil {
		return nil, err
	}
	return k,nil
}
// OpenAndLoad opens the named database file for file operations. Also, loads
// the database file into map to in-memory operations.
func OpenAndLoad(dbName string) (*Db, error) {
	fullPath := baseDir + dbName + fileExtension
	dbFile, err := os.OpenFile(fullPath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	k := &Db{
		name:     dbName,
		dir:      fullPath,
		dbFile:   dbFile,
	}
	err = k.Load()
	if err != nil {
		return nil, err
	}
	return k, nil
}

// Close closes the file.
func (k *Db) Close() error {
	return k.Write()
}

func (k *Db) Load() error {
	m := make(map[string]string)
	buf, err := os.ReadFile(k.dir)

	if err != nil {
		return err
	}

	fileData := string(buf[:])
	if fileData == "" {
		k.keyvalue = m
		return nil
	}
	dataArr := strings.Split(fileData, "\n")
	for i := 0; i < len(dataArr)-1; i++ {
		data := strings.Split(dataArr[i], "=")
		k, v := data[0], data[1]
		m[k] = v
	}
	k.keyvalue = m
	return nil
}
// write saves data into file. It writes the data in map to the file.
func (k *Db) Write() error {
	defer k.dbFile.Close()
	d := ""
	for key, val := range k.keyvalue {
		d += fmt.Sprintf("%s=%s\n", key, val)
	}
	return ioutil.WriteFile(k.dir, []byte(d), 0666)
}
func CreateServer(serverName string) {
	k,_:=Open(serverName)
	// fmt.Println(k)
	k.Close()
}