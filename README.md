<p align="center">
    <img width="200" src="https://raw.githubusercontent.com/kyoto-framework/kyoto/master/logo.svg" />
</p>

<h1 align="center">kyoto</h1>

<p align="center">
    Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.
</p>

<p align="center">
    <a href="https://goreportcard.com/report/git.sr.ht/~kyoto-framework/kyoto">
        <img src="https://goreportcard.com/badge/git.sr.ht/~kyoto-framework/kyoto">
    </a>
    <a href="https://codecov.io/gh/kyoto-framework/kyoto">
        <img src="https://codecov.io/gh/kyoto-framework/kyoto/branch/master/graph/badge.svg?token=XVLKT20DP8">
    </a>
    <a href="https://pkg.go.dev/git.sr.ht/~kyoto-framework/kyoto">
        <img src="https://pkg.go.dev/badge/git.sr.ht/~kyoto-framework/kyoto.svg">
    </a>
    <a href="https://opencollective.com/kyoto-framework">
        <img src="https://img.shields.io/opencollective/all/kyoto-framework?label=backers%20%26%20sponsors">
    </a>
    <img src="https://img.shields.io/github/license/kyoto-framework/kyoto">
    <img src="https://visitor-badge.glitch.me/badge?page_id=kyoto-framework&left_color=grey&right_color=green">
</p>

<p align="center">
    <a href="https://pkg.go.dev/git.sr.ht/~kyoto-framework/kyoto">Documentation</a>&nbsp;&bull; <a href="#team">Team</a>&nbsp;&bull; <a href="#who-uses">Who uses?</a>&nbsp;&bull; <a href="#support-us">Support us</a>
</p>

## Motivation

Creating asynchronous and dynamic layout parts is a complex problem for larger projects using `html/template`.
Library tries to simplify this process.

## What kyoto proposes?

- Organize code into configurable and standalone components structure
- Get rid of spaghetti inside of handlers
- Simple asynchronous lifecycle
- Built-in dynamics like Hotwire or Laravel Livewire
- Using a familiar built-in `html/template`
- Full control over project setup (minimal dependencies)
- 0kb JS payload without actions client (~12kb when including a client)
- Minimalistic utility-first package to simplify work with Go
- Internationalizing helper
- Cache control helper package (with a CDN page caching setup guide)

## Reasons to opt out

- API may change drastically between major versions
- You want to develop SPA/PWA
- You're just feeling OK with JS frameworks
- Not situable for a frontend with a lot of client-side logic

## Team

- Yurii Zinets: [email](mailto:yurii.zinets@icloud.com), [telegram](https://t.me/yuriizinets)
- Viktor Korniichuk: [email](mailto:rowdyhcs@gmail.com), [telegram](https://t.me/dinoarmless)

## Who uses?

### Broker One

**Website**: [https://mybrokerone.com](https://mybrokerone.com)

The first version of the site was developed with Vue and suffered from large payload and low performance.
After discussion, it was decided to migrate to Go with a built-in `html/template` due to existing libraries infrastructure inside of the project.  
Despite the good performance result, the code was badly structured and it was very uncomfortable to work in existing paradigm.  
On the basis of these problems, kyoto was born. Now, this library lies in the core of the platform.

### Using the library in your project?

Please tell us about your story! We would love to talk about your usage experience.

## Support us

Any project support is appreciated! Providing a feedback, pull requests, new ideas, whatever. Also, donations and sponsoring will help us to keep high updates frequency. Just send us a quick email or a message on contacts provided above.

If you have an option to donate us with a crypto, we have some addresses.

Bitcoin: `bc1qgxe4u799f8pdyzk65sqpq28xj0yc6g05ckhvkk`  
Ethereum: `0xEB2f24e830223bE081264e0c81fb5FD4DDD2B7B0`

Also, we have a page on open collective for backers support.

Open Collective: [https://opencollective.com/kyoto-framework](https://opencollective.com/kyoto-framework)