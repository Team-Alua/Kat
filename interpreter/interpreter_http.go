package interpreter

import (
    "io"
    "net/http"
    "errors"
    "os"
    "github.com/dop251/goja"
)

func (i *Interpreter) LoadHttpBuiltins() {
    vm := i.vm
    vm.Set("download", func(path, url string) goja.Value {
        return i.DownloadFile(path, url)
    });
}

// https://golangcode.com/download-a-file-from-a-url/
// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadToWriter(w io.Writer, url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return errors.New("There was an issue with the provided url.")
    }

    // Write the body to file
    _, err = io.Copy(w, resp.Body)
    return err
}

func DownloadFile(path string, url string) error {
    // Create the file
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    return DownloadToWriter(out, url)
}


func (i *Interpreter) DownloadFile(path, url string) goja.Value {
    fs := i.fs
    fh, err := fs.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0777)
    if err != nil {
        panic(err)
    }
    defer fh.Close()
    DownloadToWriter(fh, url)
    return i.vm.ToValue(nil)
}
