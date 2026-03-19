export interface DocPage {
	id: string;
	title: string;
	/** Path relative to /public, fetched at runtime */
	file: string;
}

export interface DocGroup {
	group: string;
	pages: DocPage[];
}

/**
 * Central configuration for the Explore documentation pages.
 *
 * To add a new page:
 *   1. Drop a .md file into frontend/public/docs/
 *   2. Add an entry here pointing to that file
 *
 * To add a new group, append a new DocGroup object.
 */
const docsConfig: DocGroup[] = [
	{
		group: 'Getting Started',
		pages: [
			{
				id: 'introduction',
				title: 'Introduction',
				file: '/docs/introduction.md',
			},
			{
				id: 'quick-start',
				title: 'Quick Start',
				file: '/docs/quick-start.md',
			},
			{
				id: 'installation',
				title: 'Installation',
				file: '/docs/installation.md',
			},
		],
	},
	{
		group: 'Application Guide',
		pages: [
			{
				id: 'app-architecture',
				title: 'Architecture Overview',
				file: '/docs/app-architecture.md',
			},
			{
				id: 'project-types',
				title: 'Project Types',
				file: '/docs/project-types.md',
			},
			{
				id: 'frameworks',
				title: 'Choosing a Framework',
				file: '/docs/frameworks.md',
			},
			{
				id: 'microservices',
				title: 'Microservices Best Practices',
				file: '/docs/microservices.md',
			},
		],
	},
	{
		group: 'How-To',
		pages: [
			{
				id: 'how-to-web-ui',
				title: 'Using the Web UI',
				file: '/docs/how-to-web-ui.md',
			},
			{
				id: 'how-to-cli',
				title: 'Using the goini CLI',
				file: '/docs/how-to-cli.md',
			},
			{
				id: 'how-to-addons',
				title: 'Working with Add-ons',
				file: '/docs/how-to-addons.md',
			},
			{
				id: 'docker-setup',
				title: 'Docker Setup',
				file: '/docs/docker-setup.md',
			},
			{
				id: 'how-to-self-host',
				title: 'Self-Hosting',
				file: '/docs/how-to-self-host.md',
			},
		],
	},
	{
		group: 'Reference',
		pages: [
			{
				id: 'api-reference',
				title: 'REST API Reference',
				file: '/docs/api-reference.md',
			},
			{
				id: 'configuration',
				title: 'Configuration',
				file: '/docs/configuration.md',
			},
			{
				id: 'troubleshooting',
				title: 'Troubleshooting',
				file: '/docs/troubleshooting.md',
			},
		],
	},
	{
		group: 'Community',
		pages: [
			{
				id: 'contributing',
				title: 'Contributing',
				file: '/docs/contributing.md',
			},
			{
				id: 'community-resources',
				title: 'Community Resources',
				file: '/docs/community-resources.md',
			},
		],
	},
];

export default docsConfig;
