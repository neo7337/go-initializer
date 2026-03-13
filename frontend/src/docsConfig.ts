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
		],
	},
	{
		group: 'Guides',
		pages: [
			{
				id: 'microservices',
				title: 'Microservices Best Practices',
				file: '/docs/microservices.md',
			},
			{
				id: 'frameworks',
				title: 'Choosing a Framework',
				file: '/docs/frameworks.md',
			},
		],
	},
	{
		group: 'Community',
		pages: [
			{
				id: 'community-resources',
				title: 'Community Resources',
				file: '/docs/community-resources.md',
			},
		],
	},
];

export default docsConfig;
