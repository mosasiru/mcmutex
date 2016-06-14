## About

This is a simple mutex library using memcached for the Go programming language
(http://golang.org/).

Memcached client is golibmc (https://github.com/douban/libmc/blob/master/src/golibmc.go)

GoDoc is Here (https://godoc.org/github.com/mosasiru/mcmutex)

## Installing

### Using *go get*

    $ go get github.com/mosasiru/mcmutex

After this command *mutex* is ready to use. Its source will be in:

    $GOPATH/src/github.com/mosasiru/mcmutex

## Example

    import (
        "github.com/douban/libmc/golibmc"
        "github.com/mosasiru/mcmutex"
    )

    func main() {
        mc := golibmc.SimpleNew([]string{"127.0.0.1:11211"})
        mutex := mcmutex.NewMCMutex(mc)
        defer mutex.Unlock("key")
        if err := mutex.Lock("key"); err != nil {
            println("lock failed")
            return
        }
        ...
    }

## Configure

### Retry

retry count before acquisition lock (default: 0)

### Interval

retry interval (default: 10ms)

###  Expiration

lock will be expired after Expiration time (default: 30s)


## Appendix

This library features perl's KeyedMutex::Memcached (http://search.cpan.org/~zigorou/KeyedMutex-Memcached-0.01/lib/KeyedMutex/Memcached.pm)
