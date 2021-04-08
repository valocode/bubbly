module.exports = {
  title: 'Bubbly',
  tagline: 'Release Readiness in a Bubble',
  url: 'https://bubbly.dev',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  onDuplicateRoutes: 'throw',
  favicon: 'img/logo.svg',
  organizationName: 'valocode', // Usually your GitHub org/user name.
  projectName: 'bubbly', // Usually your repo name.
  themeConfig: {
    prism: {
      additionalLanguages: ['hcl'],
    },
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
        {
          href: 'https://github.com/valocode/bubbly',
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
              label: 'Getting Started',
              to: 'getting-started/getting-started/',
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
              href: 'https://github.com/valocode/bubbly',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Valocode. Built with docusaurus.`,
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
            'https://github.com/valocode/bubbly/edit/main/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
