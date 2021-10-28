
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
                        link: '/concepts/#structures'
                    },
                    {
                        text: 'Rendering lifecycle',
                        link: '/concepts/#rendering-lifecycle'
                    },
                    {
                        text: 'Lifecycle integration',
                        link: '/concepts/#lifecycle-integration'
                    },
                    {
                        text: 'Methods overloading',
                        link: '/concepts/#methods-overloading'
                    },
                ]
            },
            {
                text: 'Core Features',
                link: '/core-features',
                children: [
                    {
                        text: 'Page rendering',
                        link: '/core-features/#page-rendering'
                    },
                    {
                        text: 'Built-in handler',
                        link: '/core-features/#built-in-handler'
                    },
                    {
                        text: 'Context management',
                        link: '/core-features/#context-management'
                    },
                    {
                        text: 'Component lifecycle',
                        link: '/core-features/#component-lifecycle'
                    }
                ]
            },
            {
                text: 'Extended Features',
                link: '/extended-features',
                children: [
                    {
                        text: 'Server Side Actions',
                        link: '/extended-features/#server-side-actions',
                        children: [
                            {
                                text: 'Installation',
                                link: '/extended-features/#ssa-installation'
                            },
                            {
                                text: 'Usage',
                                link: '/extended-features/#ssa-usage'
                            },
                            {
                                text: 'Lifecycle',
                                link: '/extended-features/#ssa-lifecycle'
                            },
                            {
                                text: 'Notes',
                                link: '/extended-features/#ssa-notes'
                            },
                            {
                                text: 'Limitations',
                                link: '/extended-features/#ssa-limitations'
                            }
                        ]
                    },
                    {
                        text: 'Server Side State',
                        link: '/extended-features/#server-side-state'
                    },
                    {
                        text: 'Meta builder',
                        link: '/extended-features/#meta-builder'
                    },
                    {
                        text: 'Insights',
                        link: '/extended-features/#insights'
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
                    },
                    {
                        text: 'Downsides of the library',
                        link: '/additional-notes/#downsides-of-the-library'
                    }
                ]
            },
            {
                text: 'Example with guide',
                link: '/example-with-guide'
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