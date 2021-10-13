# Prologue

This library implements an HTML render engine concept that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (`net/http`/`gin`/etc), that provides io.Writer in response.

::: warning
This project in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/yuriizinets/ssceng/issues/new)
:::

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
