import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Card, Button, Flex, Heading, Text } from '@radix-ui/themes';
import docsConfig, { DocPage } from '../docsConfig';

const allPages: DocPage[] = docsConfig.flatMap((g) => g.pages);

const Explore: React.FC<{ onBack: () => void }> = ({ onBack }) => {
	const [selectedId, setSelectedId] = useState(allPages[0].id);
	const [content, setContent] = useState('');
	const [loading, setLoading] = useState(true);

	const selectedPage = allPages.find((p) => p.id === selectedId)!;

	useEffect(() => {
		setLoading(true);
		setContent('');
		fetch(selectedPage.file)
			.then((res) => {
				if (!res.ok) throw new Error(`Failed to load ${selectedPage.file}`);
				return res.text();
			})
			.then((text) => {
				setContent(text);
				setLoading(false);
			})
			.catch(() => {
				setContent('_Could not load page content._');
				setLoading(false);
			});
	}, [selectedPage.file]);

	return (
		<Flex direction="column" width="100%" style={{ padding: '2rem', height: '100%', boxSizing: 'border-box', overflow: 'hidden' }}>
			<Flex style={{ width: '100%', gap: 32, alignItems: 'stretch', flex: 1, overflow: 'hidden' }}>
				{/* Sidebar */}
				<Card style={{ width: 280, minWidth: 220, padding: '2rem 0', display: 'flex', flexDirection: 'column', gap: 0, overflowY: 'auto', flexShrink: 0 }}>
					<Heading size="5" weight="bold" style={{ margin: '0 0 1.5rem 1.5rem' }}>Explore</Heading>
					{docsConfig.map((group) => (
						<div key={group.group} style={{ marginBottom: '1.25rem' }}>
							<Text
								size="1"
								weight="bold"
								style={{
									textTransform: 'uppercase',
									letterSpacing: '0.08em',
									color: '#888',
									padding: '0 0 0.4rem 1.5rem',
									display: 'block',
								}}
							>
								{group.group}
							</Text>
							<ul style={{ listStyle: 'none', margin: 0, padding: 0 }}>
								{group.pages.map((page) => {
									const active = selectedId === page.id;
									return (
										<li key={page.id} style={{ margin: 0, padding: 0 }}>
											<Button
												variant="ghost"
												onClick={() => setSelectedId(page.id)}
												size="2"
												style={{
													width: '100%',
													justifyContent: 'flex-start',
													fontWeight: active ? 700 : 400,
													fontSize: 15,
													borderLeft: active ? '3px solid #ffd700' : '3px solid transparent',
													borderRadius: 0,
													padding: '0.55rem 1.5rem',
													background: active ? 'rgba(255, 215, 0, 0.08)' : undefined,
													color: active ? '#bfa100' : undefined,
													transition: 'background 0.15s, color 0.15s',
												}}
												aria-current={active ? 'page' : undefined}
											>
												{page.title}
											</Button>
										</li>
									);
								})}
							</ul>
						</div>
					))}
					<Flex style={{ padding: '1.5rem 1.5rem 0.5rem' }}>
						<Button
							onClick={onBack}
							variant="outline"
							color="gold"
							size="2"
							style={{ width: '100%', fontWeight: 700 }}
						>
							← Back to Generator
						</Button>
					</Flex>
				</Card>

				{/* Main Content */}
				<Card style={{ flex: 1, minHeight: 0, padding: 0, overflow: 'hidden' }}>
					<div style={{ height: '100%', overflowY: 'auto', padding: '2.5rem 3rem', boxSizing: 'border-box' }}>
						{loading ? (
							<Text size="3" style={{ color: '#888' }}>Loading…</Text>
						) : (
							<div className="docs-content">
								<ReactMarkdown remarkPlugins={[remarkGfm]}>
									{content}
								</ReactMarkdown>
							</div>
						)}
					</div>
				</Card>
			</Flex>
		</Flex>
	);
};

export default Explore;
