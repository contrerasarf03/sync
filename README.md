sync
====

`sync` is an extension of the sync package of go with specific implementations. (eg. Waiting for OS interrupt and Parallel work)

To use:

```
go get -u github.com/contrerasarf03/sync.git
```

```
import "github.com/contrerasarf03/sync.git"

func test() {

    go doSomeThingToInterrupt()

    cleanup := func() error {
        // do some cleanup like closing kafka readers
        return nil
    }

    // will wait until interrupt then run the cleanup code
    sync.WaitForInterrupt(cleanup)
}
```