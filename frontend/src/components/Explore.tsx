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
	const [collapsedGroups, setCollapsedGroups] = useState<Record<string, boolean>>({});

	const toggleGroup = (group: string) =>
		setCollapsedGroups((prev) => ({ ...prev, [group]: !prev[group] }));

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
						<div key={group.group} style={{ marginBottom: '0.25rem' }}>
							{/* Folder row */}
							<button
								onClick={() => toggleGroup(group.group)}
								style={{
									all: 'unset',
									cursor: 'pointer',
									display: 'flex',
									alignItems: 'center',
									gap: 6,
									width: '100%',
									boxSizing: 'border-box',
									padding: '0.4rem 1rem',
									borderRadius: 4,
									transition: 'background 0.15s',
								}}
								onMouseEnter={e => (e.currentTarget.style.background = 'rgba(255,255,255,0.05)')}
								onMouseLeave={e => (e.currentTarget.style.background = 'transparent')}
							>
								<span style={{ fontSize: 11, color: '#888', width: 12, flexShrink: 0 }}>
									{collapsedGroups[group.group] ? '▶' : '▼'}
								</span>
								<span style={{ fontSize: 13, color: '#ffd700' }}>📁</span>
								<Text
									size="1"
									weight="bold"
									style={{
										textTransform: 'uppercase',
										letterSpacing: '0.08em',
										color: '#aaa',
									}}
								>
									{group.group}
								</Text>
							</button>

							{/* File rows */}
							{!collapsedGroups[group.group] && (
								<ul style={{ listStyle: 'none', margin: 0, padding: 0 }}>
									{group.pages.map((page, idx, arr) => {
										const active = selectedId === page.id;
										const isLast = idx === arr.length - 1;
										return (
											<li key={page.id} style={{ margin: 0, padding: 0, display: 'flex', alignItems: 'stretch' }}>
												{/* Tree lines */}
												<div style={{ width: 36, flexShrink: 0, position: 'relative', marginLeft: 16 }}>
													{/* Vertical line */}
													<div style={{
														position: 'absolute',
														left: 10,
														top: 0,
														bottom: isLast ? '50%' : 0,
														width: 1,
														background: '#444',
													}} />
													{/* Horizontal line */}
													<div style={{
														position: 'absolute',
														left: 10,
														top: '50%',
														width: 14,
														height: 1,
														background: '#444',
													}} />
												</div>
												<button
													onClick={() => setSelectedId(page.id)}
													style={{
														all: 'unset',
														cursor: 'pointer',
														flex: 1,
														display: 'flex',
														alignItems: 'center',
														gap: 6,
														padding: '0.4rem 0.75rem 0.4rem 0',
														fontSize: 13,
														fontWeight: active ? 700 : 400,
														color: active ? '#ffd700' : 'inherit',
														borderLeft: active ? '2px solid #ffd700' : '2px solid transparent',
														background: active ? 'rgba(255, 215, 0, 0.08)' : 'transparent',
														borderRadius: '0 4px 4px 0',
														transition: 'background 0.15s, color 0.15s',
													}}
													onMouseEnter={e => { if (!active) e.currentTarget.style.background = 'rgba(255,255,255,0.05)'; }}
													onMouseLeave={e => { if (!active) e.currentTarget.style.background = 'transparent'; }}
													aria-current={active ? 'page' : undefined}
												>
													<span style={{ fontSize: 13 }}>📄</span>
													{page.title}
												</button>
											</li>
										);
									})}
								</ul>
							)}
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
