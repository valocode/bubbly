

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
      label: 'Tutorials',
      items: [
        'tutorials/github-metrics',
        'tutorials/snyk-metrics',
        'tutorials/gosec-metrics'
      ],
    },
    {
      type: 'doc',
      id: 'adapters/adapters',
    },
    {
      type: 'doc',
      id: 'policies/policies',
    },
    {
      type: 'doc',
      id: 'schema/schema',
    },
    {
      type: 'doc',
      id: 'api/api',
    },
    {
      type: 'doc',
      id: 'graphql/graphql',
    },
    {
      type: 'category',
      label: 'CLI',
      items: [
        'cli/bubbly',
        'cli/bubbly_adapter',
        'cli/bubbly_demo',
        'cli/bubbly_policy',
        'cli/bubbly_release',
        'cli/bubbly_server',
        'cli/bubbly_version',
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
  ],
};
