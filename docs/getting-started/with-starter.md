# With Starter

The quickest way to get started with `kyoto` is a pre-prepared starter project.

## Installation

First, you need to clone the starter as a new project.  
`<app name>` - desired project name

```bash
git clone --recursive https://github.com/kyoto-framework/starter <app name>
```

Then, we need to install node dependencies and build the CSS.  
All statics and assets used for building are placed in the `static` directory.

```bash
(cd static; npm i; npm run build)
```

For the final step, you'll need to set the git origin URL to your own repository's URL.  
`<repo>` - your's project repository's URL

```bash
git remote set-url origin <repo>
```

## What's Included?

- `kyoto` - core library
- `uikit` - [UI Kit](https://github.com/kyoto-framework/uikit), built on top of `kyoto`
- `tailwindcss` - [Tailwind CSS](https://imgur.com/RN4YbvR.png) library
