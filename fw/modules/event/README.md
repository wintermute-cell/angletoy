# Module: event

The event module was dapted from:
[AlexanderGrom/go-event](https://github.com/AlexanderGrom/go-event).

It enables you to register and trigger globally available events.

## Performance

The event module heavily uses reflection to provide *some* level of type safety
(although only at runtime). But thanks to Go reflection being pretty fast this
should only become a problem if you trigger >100000 events per frame.

## Usage

The event systems is built around a dispatcher. The module provides global
default dispatcher, but you can also create your own.

```
                                  ┌──────► Listen()
                                  │                
                   ┌────────────┐ ├──────► Listen()
                   │            │ │                
Trigger() ───────► │ Dispatcher ├─┼──────► Listen()
                   │            │ │                
                   └────────────┘ ├──────► Listen()
                                  │                
                                  └──────► Listen()
```

In practice/code this looks like this:

```go
err := event.Listen("test", func(a int) error {
    fmt.Println("event 'test' was triggered!")
    return nil // we could return an error here. This would stop any other listeners from being called.
})
if err != nil { // always check your errors.
    // this can help you catch a faulty call, since events don't have compile time type safety.
    fmt.Println("Error listening:", err)
}

// ...

for rl.WindowShouldClose() {
    // ...

    if rl.IsKeyPressed(rl.KeyE) {
        err := event.Trigger("test")
        if err != nil { // don't forget to check the error! 
            // this can help you catch a faulty call, since events don't have compile time type safety.
            fmt.Println("Error triggering event:", err)
        }
    }

    // ...
}
```

### Creating Your Own Dispatcher
Should you need multiple separate dispatchers or want to manage the dispatchers
lifetime yourself, you may create your own like this:

```go
myDispatcher := event.NewDispatcher()

// we call the functions on our dispatcher instance instead of on the imported
// module.
err := myDispatcher.Listen("test", func(a int) error {
    fmt.Println("event 'test' was triggered!")
    return nil
})

```
