package main

import "io/ioutil"

import "github.com/go-martini/martini"
import "github.com/coreos/go-etcd/etcd"

func main() {
    etcd_client := etcd.NewClient([]string{"http://localhost:4001"})

    m := martini.Classic()

    m.Get("/", func() string {
        indexFile, err := ioutil.ReadFile("templates/index.html")
        if (err != nil) {
            panic(err)
        }

        return string(indexFile)
    })

    m.Get("/api/list", func() string {
        key, err := etcd_client.Get("mykey", true, false)
        if (err != nil) {
            panic(err)
        }
        return key.Node.Value
    })

    m.Run()
}
