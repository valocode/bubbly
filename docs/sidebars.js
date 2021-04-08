

module.exports = {
  docs: [
    {
      type: 'category',
      label: 'Bubbly',
      items: [
        'introduction/introduction',
        'introduction/use-cases',
        'introduction/core-concepts',
      ],
    },
    {
      type: 'doc',
      id: 'getting-started/getting-started',
    },
    {
      type: 'category',
      label: 'Resources',
      items: [
        'resources/overview',
        'resources/kinds',
      ],
    },
    {
      type: 'doc',
      id: 'schema/schema'
    },
    {
      type: 'doc',
      id: 'api/api'
    },
    {
      type: 'category',
      label: 'CLI',
      items: [
        'cli/bubbly',
        'cli/bubbly-agent',
        'cli/bubbly-apply',
        'cli/bubbly-get',
        'cli/bubbly-schema',
        'cli/schema/bubbly-schema-apply',
      ],
    },
    {
      type: 'category',
      label: 'Tutorials',
      items: [
        'tutorials/github-metrics',
        'tutorials/snyk-metrics',
        'tutorials/gosec-metrics'
      ],
    },
    {
      type: 'doc',
      id: 'current-status/status',
    },
    {
      type: 'doc',
      id: 'contributing/contributing',
    },
    {
      type: 'doc',
      id: 'future/future',
    },
  ],
};
