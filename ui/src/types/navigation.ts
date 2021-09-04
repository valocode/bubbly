// a navigationItem is used in the Sidebar components to represent a navigation
// option
export interface navigationItem {
	name: string;
	href: string;
	icon?: string;
	children?: navigationItem[];
}

// a breadcrumbItem is used in bubbly's Breadcrumb component to represent the
// text (name) and link (href) to be used at each hierarchical level within the
// Breadcrumb
export interface breadcrumbItem {
	name: string;
	href: string;
}
