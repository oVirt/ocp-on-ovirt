1. Get the go-ovirt
```
$ go get github.com/ovirt/go-ovirt
```

2. Go to example dir
```
    $ cd network
```

2. Edit list_networks.go:
```
    $ vi list_networks.go
        1. Change domain: foobar.mydomain.home
        2. Username
        3. Password
```
3. Execute
```
$ go run ./list_networks.go
  Cluster: Default id: 40285c82-00e0-11ea-acb1-ecf4bbf47b7c networks: [net1 net2 ovirtmgmt]
```
