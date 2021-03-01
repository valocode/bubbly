module.exports = {
  title: 'Bubbly',
  tagline: 'Release Readiness in a Bubble',
  url: 'https://bubbly.dev',
  baseUrl: '/docs/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  onDuplicateRoutes: 'throw',
  favicon: 'img/logo.svg',
  organizationName: 'verifa', // Usually your GitHub org/user name.
  projectName: 'bubbly', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Bubbly',
      logo: {
        alt: 'Bubbly Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          to: '/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {to: 'blog', position: 'left'},
        {
          href: 'https://github.com/verifa/bubbly',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Introduction',
              to: '/',
            },
            {
              label: 'Quickstart',
              to: 'getting-started/quickstart/',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Stack Overflow',
              href: 'https://stackoverflow.com/questions/tagged/bubbly',
            },
            {
              label: 'Discord',
              href: 'https://discordapp.com/invite/bubbly',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/bubbly',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/verifa/bubbly',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Verifa. Built with docusaurus.`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          routeBasePath: '/',
          editUrl:
            'https://github.com/verifa/bubbly/edit/master/website/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
