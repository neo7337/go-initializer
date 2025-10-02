import React, { useState } from 'react';
import { Card, Button, Flex, Heading, Text } from '@radix-ui/themes';

const blogPosts = [
	{
		id: 1,
		title: 'Getting Started with Go Initializer',
		content: `Go Initializer helps you scaffold Go projects with best practices and modern frameworks. Select your project type, Go version, and framework to generate a ready-to-use starter kit.`,
	},
	{
		id: 2,
		title: 'Microservices Best Practices',
		content: `Learn how to structure your Go microservices for scalability and maintainability. Explore recommended folder structures, dependency management, and testing strategies.`,
	},
	{
		id: 3,
		title: 'Choosing a Go Web Framework',
		content: `Compare popular Go web frameworks like Gin, Echo, and Fiber. Understand their strengths, community support, and when to use each for your project.`,
	},
	{
		id: 4,
		title: 'Community Resources',
		content: `Find curated links to Go tutorials, open-source projects, and community forums to accelerate your learning and development.`,
	},
];

const Explore: React.FC<{ onBack: () => void }> = ({ onBack }) => {
	const [selectedId, setSelectedId] = useState(blogPosts[0].id);
	const selectedPost = blogPosts.find((post) => post.id === selectedId);

	return (
		<Flex direction="column" width="100%" style={{ padding: '2rem 0' }}>
			<Flex style={{ maxWidth: 1200, width: '100%', margin: '0 auto', gap: 32 }}>
				{/* Sidebar: Titles */}
				<Card style={{ width: 300, minWidth: 220, padding: '2rem 0', height: 'fit-content', display: 'flex', flexDirection: 'column', gap: 0 }}>
					<Heading size="5" weight="bold" style={{ margin: '0 0 1.5rem 2rem' }}>Explore</Heading>
					<ul style={{ listStyle: 'none', margin: 0, padding: 0, display: 'flex', flexDirection: 'column', gap: 8 }}>
						{blogPosts.map((post, idx) => (
							<li key={post.id} style={{ margin: 0, padding: 0 }}>
								<Button
									variant={selectedId === post.id ? 'solid' : 'ghost'}
									color={selectedId === post.id ? 'gold' : 'gray'}
									onClick={() => setSelectedId(post.id)}
									size="3"
									radius="large"
									style={{
										display: 'flex',
										alignItems: 'center',
										gap: 16,
										justifyContent: 'flex-start',
										fontWeight: selectedId === post.id ? 700 : 500,
										fontSize: 18,
										borderLeft: selectedId === post.id ? '4px solid #ffd700' : '4px solid transparent',
										borderRadius: 12,
										padding: '0.9rem 2rem 0.9rem 2.2rem',
										background: selectedId === post.id ? '#fffbe6' : undefined,
										color: selectedId === post.id ? '#bfa100' : undefined,
										textAlign: 'left',
										boxShadow: selectedId === post.id ? '0 2px 8px 0 #ffd70022' : undefined,
										transition: 'background 0.18s, color 0.18s, box-shadow 0.18s',
										outline: selectedId === post.id ? '2px solid #ffd700' : undefined,
									}}
									aria-current={selectedId === post.id ? 'true' : undefined}
									tabIndex={0}
								>
									<span style={{ fontSize: 20, opacity: selectedId === post.id ? 1 : 0.5, transition: 'opacity 0.2s', marginRight: 4 }}>
										{selectedId === post.id ? '★' : '•'}
									</span>
									<span style={{ whiteSpace: 'pre-line' }}>{post.title}</span>
								</Button>
							</li>
						))}
					</ul>
				</Card>
				{/* Main Content: Selected Post */}
				<Card style={{ flex: 1, minHeight: 400, padding: '2.5rem 2.5rem 2.5rem 4.5rem', display: 'flex', flexDirection: 'column', justifyContent: 'flex-start', position: 'relative' }}>
					<Heading size="7" weight="bold" style={{ marginBottom: 8, marginTop: 0 }}>{selectedPost?.title}</Heading>
					<Text size="5" style={{ lineHeight: 1.7 }}>{selectedPost?.content}</Text>
					<Flex style={{ marginTop: 'auto', textAlign: 'center', justifyContent: 'center' }}>
						<Button
							onClick={onBack}
							variant="outline"
							color="gold"
							size="3"
							style={{ fontWeight: 700, fontSize: 17 }}
						>
							← Back to Generator
						</Button>
					</Flex>
				</Card>
			</Flex>
		</Flex>
	);
};

export default Explore;
