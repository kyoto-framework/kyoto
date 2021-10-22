
module.exports = {
    title: 'kyoto',
    description: 'Keep your Go stack united with Server Side Components',

    base: '/docs/',

    themeConfig: {
        navbar: [
            {
                text: 'Home',
                link: 'https://kyoto.codes'
            },
            {
                text: 'GitHub',
                link: 'https://github.com/yuriizinets/kyoto'
            }
        ],
        sidebar: [
            {
                text: 'Prologue',
                link: '/prologue'
            },
            {
                text: 'Get Started',
                link: '/get-started',
                children: [
                    {
                        text: 'Quick Start',
                        link: '/get-started/#quick-start'
                    },
                    {
                        text: 'Installation',
                        link: '/get-started/#installation'
                    },
                    {
                        text: 'Integration',
                        link: '/get-started/#integration'
                    },
                ],
            },
            {
                text: 'Concepts',
                link: '/concepts',
                children: [
                    {
                        text: 'Structures',
                        link: '/docs/concepts/#structures'
                    },
                    {
                        text: 'Rendering lifecycle',
                        link: '/docs/concepts/#rendering-lifecycle'
                    },
                    {
                        text: 'Lifecycle integration',
                        link: '/docs/concepts/#lifecycle-integration'
                    },
                    {
                        text: 'Methods overloading',
                        link: '/docs/concepts/#methods-overloading'
                    },
                ]
            },
            {
                text: 'Core Features',
                link: '/docs/core-features',
                children: [
                    {
                        text: 'Page rendering',
                        link: '/docs/core-features/#page-rendering'
                    },
                    {
                        text: 'Built-in handler',
                        link: '/docs/core-features/#built-in-handler'
                    },
                    {
                        text: 'Context management',
                        link: '/docs/core-features/#context-management'
                    },
                    {
                        text: 'Component lifecycle',
                        link: '/docs/core-features/#component-lifecycle'
                    }
                ]
            },
            {
                text: 'Extended Features',
                link: '/docs/extended-features',
                children: [
                    {
                        text: 'Server Side Actions',
                        link: '/docs/extended-features/#server-side-actions',
                        children: [
                            {
                                text: 'Installation',
                                link: '/docs/extended-features/#ssa-installation'
                            },
                            {
                                text: 'Usage',
                                link: '/docs/extended-features/#ssa-usage'
                            },
                            {
                                text: 'Lifecycle',
                                link: '/docs/extended-features/#ssa-lifecycle'
                            },
                            {
                                text: 'Notes',
                                link: '/docs/extended-features/#ssa-notes'
                            }
                        ]
                    },
                    {
                        text: 'Server Side State',
                        link: '/docs/extended-features/#server-side-state'
                    },
                    {
                        text: 'Meta builder',
                        link: '/docs/extended-features/#meta-builder'
                    },
                    {
                        text: 'Insights',
                        link: '/docs/extended-features/#insights'
                    }
                ]
            },
            {
                text: 'Additional notes',
                link: '/additional-notes',
                children: [
                    {
                        text: 'Dealing with Go packages',
                        link: '/additional-notes/#dealing-with-go-packages'
                    }
                ]
            },
            {
                text: 'Example with guide',
                link: '/docs/example-with-guide'
            }
        ]
    },

    plugins: [
        ['@vuepress/plugin-search', {
            locales: {
                '/': {
                    placeholder: 'Search'
                }
            }
        }]
    ],
}