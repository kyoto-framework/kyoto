
# Prologue

This library implements an HTML render engine concept that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (`net/http`/`gin`/etc), that provides io.Writer in response.

!!! warning
    This project in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/kyoto-framework/kyoto/issues/new)

## Motivation

Main motivation is to reduce usage of popular SPA/PWA frameworks where it's not needed because it adds a lot of complexity and overhead. There is no reason to bring significant runtime, VirtualDOM, and Webpack into the project with minimal dynamic frontend behavior. This project proves the possibility of keeping most of the logic on the server's side.

## What problems does it solve?

While developing the website's frontend with plain Go templates, I discovered some of the downsides of this approach:

- With plain `html/template` you're starting to repeat yourself. It's harder to define reusable parts.
- You must repeat DTO calls for each page, where you're using reusable parts.
- With Go's routines approach it's hard to make async-like DTO calls in the handlers.
- For dynamic things, you still need to use JS and client-side DOM modification.

Complexity is much higher when all of them get combined.

This engine tries to bring components and async experience to the traditional server-side rendering.

## Zen

For contributors:  

- Don't replace Go features that exist already
- Don't do work that's already done
- Don't force developers to use a specific solution (Gin/Chi/GORM/sqlx/etc). Let them choose
- Rely on the server to do the rendering, minimum JS specifics or client-side only behavior

## Alternatives?

### Notes

- I'm not going to compare languages pros & cons, just purpose and ideology
- Also, I'm not going to take into account existing codebases

### Laravel Livewire

Laravel is awesome on my opinion and livewire makes it even more awesome! Nice choise for people who want to have "battries included" framework. Kyoto is not a framework, it's just a small library and tries to solve another kind of problem - components and asynchronous operations organization. Features like context, Server Side Actions, Server Side State, Insights, are just extensions to Core library purpose. Also, Kyoto not delivered with batteries "included", it gives more control to developer.

Differences (Kyoto vs Livewire):

- Minimalistic over "batteries included"
- Another approach regarding client-server communication
- Another purpose

### Elixir Phoenix

To be honest, I'm far away from Elixir and Erlang ecosystem generally. If you have some time to tell me more about Phoenix, I'll be very grateful!  

Differences (Kyoto vs Phoenix):

- Need more details

### JavaScript Frameworks

The most delicious piece of cake. Please, check "Motivation" part. I'd like to notice, that Kyoto not tries to replace popular PWA/SPA approach, but to reduce it usage where it's not needed. If any of JS Frameworks works for you, so, why not?

Differences (Kyoto vs JS Frameworks):

- Give more control to developer
- Reduce amount of dependencies
- Make SSR and debugging easier

![meme](https://imgur.com/RN4YbvR.png)

### Any other?

Just create an issue or contact with email if you'll find something interseting!
