// persist package handles data saved for Nut, on the file system (in .nut/ folder)
package persist

import (
    // log "github.com/Sirupsen/logrus"
    // "io/ioutil"
    "os"
    "errors"
    // "path/filepath"
)

const Name = ".nut"

type Store interface {
    // returns the root path of the store
    getPath() string
}

type StoreBase struct {
    path string
}
    func (self *StoreBase) getPath() string {
        return self.path
    }

// Initialize the store on the hard drive if it does not
// exist yet, and return it. Throws if it couldn't create files.
func InitStore(path string) (Store, error) {
    _, err := os.Stat(path)
    if err != nil {
        err = os.MkdirAll(path, 0644) // TODO: discuss this permission level
        if err != nil {
            return nil, errors.New("Folder " + path +
                " does not exit, and could be created.")
        }
    }
    store := &StoreBase {
        path: path,
    }
    return store, nil
}

// Remove the store from hard drive
func ClearStore(store Store) {
    os.RemoveAll(store.getPath())
}

// // return a saved value
// // define for string, int, []string, []int
// func read(key, defaultValue) {}
// // save a value
// // define for several types
// func save(key) {}


