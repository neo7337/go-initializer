import React, { useState, useEffect, useCallback } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { Button } from '@radix-ui/themes';
import docsConfig, { DocPage } from '../docsConfig';

const allPages: DocPage[] = docsConfig.flatMap((g) => g.pages);

// ── SVG Icons ─────────────────────────────────────────────────────────────────

const FolderOpenIcon = () => (
	<svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M1.5 4C1.5 3.17 2.17 2.5 3 2.5H6.5L8 4H13C13.83 4 14.5 4.67 14.5 5.5V12C14.5 12.83 13.83 13.5 13 13.5H3C2.17 13.5 1.5 12.83 1.5 12V4Z" stroke="currentColor" strokeWidth="1.2" fill="none" />
	</svg>
);

const FolderClosedIcon = () => (
	<svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M1.5 4C1.5 3.17 2.17 2.5 3 2.5H6.5L8 4H13C13.83 4 14.5 4.67 14.5 5.5V12C14.5 12.83 13.83 13.5 13 13.5H3C2.17 13.5 1.5 12.83 1.5 12V4Z" stroke="currentColor" strokeWidth="1.2" fill="rgba(255,215,0,0.08)" />
	</svg>
);

const FileIcon = () => (
	<svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M4 1.5H9.5L13.5 5.5V14C13.5 14.55 13.05 15 12.5 15H3.5C2.95 15 2.5 14.55 2.5 14V2.5C2.5 1.95 2.95 1.5 3.5 1.5H4Z" stroke="currentColor" strokeWidth="1.2" fill="none" />
		<path d="M9.5 1.5V5.5H13.5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
	</svg>
);

const ChevronRightIcon = () => (
	<svg width="10" height="10" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M6 4L10 8L6 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
	</svg>
);

const ChevronDownIcon = () => (
	<svg width="10" height="10" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M4 6L8 10L12 6" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
	</svg>
);

const CopyIcon = () => (
	<svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<rect x="5" y="5" width="9" height="9" rx="1.5" stroke="currentColor" strokeWidth="1.2" />
		<path d="M3 11V3C3 2.45 3.45 2 4 2H11" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
	</svg>
);

const CheckIcon = () => (
	<svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M3 8L6.5 11.5L13 4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
	</svg>
);

const BreadcrumbSepIcon = () => (
	<svg width="12" height="12" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M6 4L10 8L6 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
	</svg>
);

const HamburgerIcon = () => (
	<svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M2 4H14M2 8H14M2 12H14" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
	</svg>
);

const XCloseIcon = () => (
	<svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden="true">
		<path d="M3 3L13 13M13 3L3 13" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
	</svg>
);

// ── Skeleton Loader ────────────────────────────────────────────────────────────

const SkeletonLoader: React.FC = () => (
	<div className="docs-skeleton" aria-busy="true" aria-label="Loading content">
		<div className="skeleton-title" />
		<div className="skeleton-line" style={{ width: '80%' }} />
		<div className="skeleton-line" style={{ width: '65%' }} />
		<div className="skeleton-line" style={{ width: '90%' }} />
		<div className="skeleton-block" />
		<div className="skeleton-line" style={{ width: '75%' }} />
		<div className="skeleton-line" style={{ width: '55%' }} />
		<div className="skeleton-line" style={{ width: '85%' }} />
	</div>
);

// ── Code Block with copy button ───────────────────────────────────────────────

const CodeBlock: React.FC<{ language: string; value: string }> = ({ language, value }) => {
	const [copied, setCopied] = useState(false);

	const handleCopy = useCallback(() => {
		navigator.clipboard.writeText(value).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		}).catch(() => {});
	}, [value]);

	return (
		<div className="docs-code-block">
			<div className="docs-code-header">
				<span className="docs-code-lang">{language || 'text'}</span>
				<button
					className="docs-code-copy"
					onClick={handleCopy}
					aria-label={copied ? 'Copied' : 'Copy code'}
					title={copied ? 'Copied!' : 'Copy code'}
				>
					{copied ? <CheckIcon /> : <CopyIcon />}
					<span>{copied ? 'Copied!' : 'Copy'}</span>
				</button>
			</div>
			<SyntaxHighlighter
				language={language || 'text'}
				style={oneDark}
				customStyle={{
					margin: 0,
					borderRadius: '0 0 var(--radius-md) var(--radius-md)',
					fontSize: '0.875rem',
					background: 'var(--color-surface-0)',
					border: 'none',
				}}
				codeTagProps={{ style: { fontFamily: "'Geist Mono', 'Fira Mono', 'Consolas', monospace" } }}
			>
				{value}
			</SyntaxHighlighter>
		</div>
	);
};

// ── Breadcrumb ─────────────────────────────────────────────────────────────────

const Breadcrumb: React.FC<{ group: string; page: string }> = ({ group, page }) => (
	<nav aria-label="Breadcrumb" className="docs-breadcrumb">
		<span className="docs-breadcrumb-group">{group}</span>
		<span className="docs-breadcrumb-sep"><BreadcrumbSepIcon /></span>
		<span className="docs-breadcrumb-page">{page}</span>
	</nav>
);

// ── Main Component ─────────────────────────────────────────────────────────────

const Explore: React.FC<{ onBack: () => void }> = ({ onBack }) => {
	const [selectedId, setSelectedId] = useState(allPages[0].id);
	const [content, setContent] = useState('');
	const [loading, setLoading] = useState(true);
	const [collapsedGroups, setCollapsedGroups] = useState<Record<string, boolean>>({});
	const [drawerOpen, setDrawerOpen] = useState(false);

	const closeDrawer = useCallback(() => setDrawerOpen(false), []);

	const toggleGroup = (group: string) =>
		setCollapsedGroups((prev) => ({ ...prev, [group]: !prev[group] }));

	const selectedPage = allPages.find((p) => p.id === selectedId)!;
	const selectedGroup = docsConfig.find((g) => g.pages.some((p) => p.id === selectedId));

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
		<div className="explore-shell">
			{/* Backdrop overlay for mobile drawer */}
			<div
				className={`explore-drawer-overlay${drawerOpen ? ' explore-drawer-overlay--visible' : ''}`}
				onClick={closeDrawer}
				aria-hidden="true"
			/>

			{/* Sidebar */}
			<nav className={`explore-sidebar${drawerOpen ? ' explore-sidebar--open' : ''}`} aria-label="Documentation navigation">
				<div className="explore-sidebar-header">
					<span className="explore-sidebar-title">Explore</span>
					<button
						className="explore-sidebar-close"
						onClick={closeDrawer}
						aria-label="Close navigation"
					>
						<XCloseIcon />
					</button>
				</div>

				<div className="explore-sidebar-tree">
					{docsConfig.map((group) => {
						const isCollapsed = !!collapsedGroups[group.group];
						return (
							<div key={group.group} className="explore-group">
								<button
									className="explore-group-toggle"
									onClick={() => toggleGroup(group.group)}
									aria-expanded={!isCollapsed}
								>
									<span className="explore-group-chevron">
										{isCollapsed ? <ChevronRightIcon /> : <ChevronDownIcon />}
									</span>
									<span className="explore-group-icon">
										{isCollapsed ? <FolderClosedIcon /> : <FolderOpenIcon />}
									</span>
									<span className="explore-group-label">{group.group}</span>
								</button>

								<div
									className="explore-group-pages"
									style={{
										maxHeight: isCollapsed ? 0 : `${group.pages.length * 44}px`,
									}}
								>
									{group.pages.map((page) => {
										const active = selectedId === page.id;
										return (
											<button
												key={page.id}
												className={`explore-page-btn${active ? ' explore-page-btn--active' : ''}`}
												onClick={() => { setSelectedId(page.id); closeDrawer(); }}
												aria-current={active ? 'page' : undefined}
											>
												<span className="explore-page-icon">
													<FileIcon />
												</span>
												<span className="explore-page-title">{page.title}</span>
											</button>
										);
									})}
								</div>
							</div>
						);
					})}
				</div>

				<div className="explore-sidebar-footer">
					<Button
						onClick={onBack}
						variant="outline"
						color="gold"
						size="2"
						style={{ width: '100%', fontWeight: 700 }}
					>
						← Back to Generator
					</Button>
				</div>
			</nav>

			{/* Main content pane */}
			<main className="explore-content" aria-label="Documentation content">
				<div className="explore-content-scroll">
					<div className="explore-topbar">
						<button
							className="explore-mobile-toggle"
							onClick={() => setDrawerOpen(true)}
							aria-label="Open navigation"
							aria-expanded={drawerOpen}
						>
							<HamburgerIcon />
							Menu
						</button>
						{selectedGroup && (
							<Breadcrumb group={selectedGroup.group} page={selectedPage.title} />
						)}
					</div>
					{loading ? (
						<SkeletonLoader />
					) : (
						<div className="docs-content">
							<ReactMarkdown
								remarkPlugins={[remarkGfm]}
								components={{
									pre({ children }) {
										return <>{children}</>;
									},
								code({ className, children, node: _node, ...props }: React.HTMLAttributes<HTMLElement> & { node?: unknown }) {
										const match = /language-(\w+)/.exec(className || '');
										const codeStr = String(children).replace(/\n$/, '');
										const isBlock = !!match || codeStr.includes('\n');
										if (isBlock) {
											return <CodeBlock language={match?.[1] ?? ''} value={codeStr} />;
										}
										return <code className={className} {...props}>{children}</code>;
									},
								}}
							>
								{content}
							</ReactMarkdown>
						</div>
					)}
				</div>
			</main>
		</div>
	);
};

export default Explore;
