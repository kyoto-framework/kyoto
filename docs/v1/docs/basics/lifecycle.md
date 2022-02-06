
# Lifecycle

We tried to get from component approach as much as we can.
That's why we also borrowed lifecycle idea from popular JS frameworks.
Why do we need it on server side?
Let's figure it out.

We will take known example from components documentation.

=== "component.uuid.go"

	```go
	package main

	import (
	    "net/http"
	    "encoding/json"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/lifecycle"
	)

	func ComponentUUID(core *kyoto.Core) {
	    lifecycle.Init(core, func() {
	        core.State.Set("UUID", "")
	    })
	    lifecycle.Async(core, func() error {
	        resp, _ := http.Get("http://httpbin.org/uuid")
	        data := map[string]string{}
	        json.NewDecoder(resp.Body).Decode(&data)
	        c.State.Set("UUID", data["uuid"])
	        return nil
	    })
	}
	```

First, let's figure out what stages lifecycle have.  
Usually, everything is executing in this order: **init** -> **async** -> **after async**.  
I'm avoiding internal steps in documentation to not confuse people,
but you can explore lifecycle in details with "Concepts" documentation category.
Each step is executed asynchronously.
That means library spawns goroutine for each instance of core receiver.

To integrate our component into lifecycle, we are using adapters from `lifecycle` module:
`Init`, `Async`, `AfterAsync`. Only `Async` and `AfterAsync` are allowed to return an error.
It was designed in this way to ensure that `Init` will be used as intended (only for state and components initialization).
