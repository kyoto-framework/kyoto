
# SSC Engine Webpage

Project landing page and documentation.  
Landing page built with SSC library itself. Documentation uses VuePress to simplify setup.

## Working with documentation

First, you'll need VuePress installed.  

```bash
$ npm i -g vuepress
```

Running documentation in dev mode:

```bash
$ vuepress dev docs
```

Building documentation in static mode:

```bash
$ vuepress build docs
```

After building in static mode, documentation can be used in landing page with `/docs/` route.
