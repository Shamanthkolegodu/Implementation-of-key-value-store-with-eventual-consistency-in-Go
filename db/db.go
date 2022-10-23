package main

import (
	"os"
	"sync"
	"time"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
)

type Database struct{
	name string
	dir string
	dbFile *os.File
	keyvalue map[string]string
	mu sync.Mutex
	Addr string
	duration time.Duration
}
const(
	baseDir = "/tmp/database/"
	fileExtension=".database"
)
// open creates data file for newly creating database. If the database file is
// already exists, it returns error without creating anything. name indicates
// database name.
func open(name string, addr string, duration time.Duration) (*Database, error) {
	fullPath := baseDir + name + fileExtension

	// Check database's base directory
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("database directory couldn't created: %s", err.Error())
		}
	}

	// Check database file's directory
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		dbFile, err := os.OpenFile(fullPath, os.O_CREATE, 0777)
		if err != nil {
			return nil, fmt.Errorf("database couldn't created: %s", err.Error())
		}
		m := make(map[string]string)
		k := &Database{
			name:     name + fileExtension,
			dir:      baseDir + name + fileExtension,
			dbFile:   dbFile,
			keyvalue:       m,
			mu:       sync.Mutex{},
			Addr:     addr,
			duration: duration,
		}

		ticker := time.NewTicker(duration)
		go func() {
			for {
				select {
				case t := <-ticker.C:
					err := k.write()
					if err != nil {
						log.Println("Writing file failed at", t.Local())
					} else {
						log.Println("Data saved on the file at", t.Local())
					}
				}
			}
		}()
		return k, nil
	}
	k, err := openAndLoad(name, addr, duration)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				err := k.write()
				if err != nil {
					log.Println("Writing file failed at", t.Local())
				} else {
					log.Println("Data saved on the file at", t.Local())
				}
			}
		}
	}()
	return k, nil
}

func (k *Database) write() error {
	defer k.dbFile.Close()
	d := ""
	for key, val := range k.keyvalue {
		d += fmt.Sprintf("%s=%s\n", key, val)
	}
	return ioutil.WriteFile(k.dir, []byte(d), 0666)
}
// openAndLoad opens the named database file for file operations. Also, loads
// the database file into map to in-memory operations.
func openAndLoad(dbName string, addr string, duration time.Duration) (*Database, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	fullPath := baseDir + dbName + fileExtension
	dbFile, err := os.OpenFile(fullPath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	k := &Database{
		name:     dbName,
		dir:      fullPath,
		dbFile:   dbFile,
		mu:       sync.Mutex{},
		Addr:     addr,
		duration: duration,
	}
	err = k.load()
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Close closes the file.
func (k *Database) Close() error {
	return k.write()
}

// load reads and loads the data from the file into map.
func (k *Database) load() error {
	k.mu.Lock()
	defer k.mu.Unlock()
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
	dataArr := 
	strings.Split(fileData, "\n")
	for i := 0; i < len(dataArr)-1; i++ {
		data := strings.Split(dataArr[i], "=")
		k, v := data[0], data[1]
		m[k] = v
	}
	k.keyvalue = m
	return nil
}

