
# Server Side State

!!! danger
    Not implemented yet. Check [issue](https://github.com/kyoto-framework/kyoto/issues/28) state

This feature is useful in case of large state payloads.
Instead of saving state inline as html tag, store state on server side and inject state hash as html tag.
Using this, you will decrease amount of data sent with SSA request and total HTML document size.
