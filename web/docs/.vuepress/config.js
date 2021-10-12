
module.exports = {
    title: 'ssceng',
    description: 'Keep your Go stack united with Server Side Components',

    base: '/docs/',

    themeConfig: {
        navbar: [
            {
                text: 'Home',
                link: 'https://ssceng.codes'
            },
            {
                text: 'GitHub',
                link: 'https://github.com/yuriizinets/ssceng'
            }
        ],
        sidebar: [
            {
                text: 'Overview',
                link: '/overview'
            },
            {
                text: 'Quick Start',
                link: '/quickstart'
            },
            {
                text: 'Concepts',
                link: '/concepts',
                children: [
                    {
                        text: 'Interfaces',
                        link: '/concepts/#interfaces'
                    },
                    {
                        text: 'Lifecycle',
                        link: '/concepts/#lifecycle'
                    }
                ]
            },
            {
                text: 'Core',
                link: '/core',
                children: [
                    {
                        text: 'Pages',
                        link: '/core/#page'
                    },
                    {
                        text: 'Components',
                        link: '/core/#component'
                    },
                    {
                        text: 'Rendering',
                        link: '/core/#render-page'
                    },
                    {
                        text: 'Context',
                        link: '/core/#context'
                    },
                    {
                        text: 'Handler',
                        link: '/core/#handler-factory'
                    },
                    {
                        text: 'Async',
                        link: '/core/#async-components'
                    },
                    {
                        text: 'Flags',
                        link: '/core/#flags'
                    }
                ]
            },
            {
                text: 'Extended',
                link: '/extended',
                children: [
                    {
                        text: 'Meta builder',
                        link: '/extended/#meta-builder'
                    },
                    {
                        text: 'Server Side Actions',
                        link: '/extended/#server-side-actions-ssa'
                    },
                    {
                        text: 'Server Side State',
                        link: '/extended/#server-side-state'
                    },
                    {
                        text: 'Insights',
                        link: '/extended/#insights'
                    }
                ]
            }
        ]
    }
}