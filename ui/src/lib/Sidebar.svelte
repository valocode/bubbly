<script lang="ts">
	import { base } from '$app/paths';
	import Icon from './Icon.svelte';

	interface navigationItem {
		name: string;
		href: string;
		icon?: string;
		children?: navigationItem[];
	}

	const items: navigationItem[] = [
		{
			name: 'Releases',
			href: `${base}/releases`,
			icon: 'hashtag'
		}
		// {
		// 	name: 'Repos',
		// 	href: `${base}/repositories`,
		// 	icon: 'sparkles'
		// }
	];
</script>

<!-- This example requires Tailwind CSS v2.0+ -->
<div class="h-screen flex overflow-hidden bg-gray-100">
	<!-- Off-canvas menu for mobile, show/hide based on off-canvas menu state. -->
	<div class="fixed inset-0 flex z-40 md:hidden" role="dialog" aria-modal="true">
		<!--
        Off-canvas menu overlay, show/hide based on off-canvas menu state.
  
        Entering: "transition-opacity ease-linear duration-300"
          From: "opacity-0"
          To: "opacity-100"
        Leaving: "transition-opacity ease-linear duration-300"
          From: "opacity-100"
          To: "opacity-0"
      -->
		<div class="fixed inset-0 bg-gray-600 bg-opacity-75" aria-hidden="true" />

		<!--
        Off-canvas menu, show/hide based on off-canvas menu state.
  
        Entering: "transition ease-in-out duration-300 transform"
          From: "-translate-x-full"
          To: "translate-x-0"
        Leaving: "transition ease-in-out duration-300 transform"
          From: "translate-x-0"
          To: "-translate-x-full"
      -->
		<div class="relative flex-1 flex flex-col max-w-xs w-full bg-white">
			<!--
          Close button, show/hide based on off-canvas menu state.
  
          Entering: "ease-in-out duration-300"
            From: "opacity-0"
            To: "opacity-100"
          Leaving: "ease-in-out duration-300"
            From: "opacity-100"
            To: "opacity-0"
        -->
			<div class="absolute top-0 right-0 -mr-12 pt-2">
				<button
					type="button"
					class="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
				>
					<span class="sr-only">Close sidebar</span>
					<!-- Heroicon name: outline/x -->
					<svg
						class="h-6 w-6 text-white"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			<div class="flex-1 h-0 pt-5 pb-4 overflow-y-auto">
				<div class="flex-shrink-0 flex items-center px-4">
					<img class="h-8 w-auto" src="{base}/logo.svg" alt="Bubbly" />
				</div>
				<nav class="mt-5 px-2 space-y-1">
					{#each items as item}
						<!-- Current: "bg-gray-100 text-gray-900", Default: "text-gray-600 hover:bg-gray-50 hover:text-gray-900" -->
						<a
							href={item.href}
							class="bg-gray-100 text-gray-900 group flex items-center px-2 py-2 text-base font-medium rounded-md"
						>
							<!--
                Heroicon name: outline/home
  
                Current: "text-gray-500", Default: "text-gray-400 group-hover:text-gray-500"
              -->
							<svg
								class="text-gray-500 mr-4 flex-shrink-0 h-6 w-6"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								aria-hidden="true"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
								/>
							</svg>
							{item.name}
						</a>
					{/each}
				</nav>
			</div>
			<div class="flex-shrink-0 flex border-t border-gray-200 p-4">
				<!-- <a href="#" class="flex-shrink-0 group block"> -->
				<div class="flex items-center">
					<div>
						<img class="inline-block h-10 w-10 rounded-full" src="{base}/empty-user.svg" alt="" />
					</div>
					<div class="ml-3">
						<p class="text-base font-medium text-gray-700 group-hover:text-gray-900">Annonymous</p>
						<!-- <p class="text-sm font-medium text-gray-500 group-hover:text-gray-700">
								View profile
							</p> -->
					</div>
				</div>
				<!-- </a> -->
			</div>
		</div>

		<div class="flex-shrink-0 w-14">
			<!-- Force sidebar to shrink to fit close icon -->
		</div>
	</div>

	<!-- Static sidebar for desktop -->
	<div class="hidden md:flex md:flex-shrink-0">
		<div class="flex flex-col w-64">
			<!-- Sidebar component, swap this element with another sidebar if you like -->
			<div class="flex-1 flex flex-col min-h-0 border-r border-gray-200 bg-bubbly-dark">
				<div class="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
					<a href="{base}/">
						<div class="flex items-center flex-shrink-0 px-4">
							<img class="h-16 w-auto" src="{base}/logo.svg" alt="Bubbly" />
						</div>
					</a>
					<nav class="mt-5 flex-1 px-2 bg-bubbly-dark space-y-1">
						{#each items as item}
							<!-- Current: "bg-gray-100 text-gray-900", Default: "text-gray-600 hover:bg-gray-50 hover:text-gray-900" -->
							<a
								href={item.href}
								class="bg-bubbly-dark text-white group flex items-center px-2 py-2 text-sm font-medium rounded-md"
							>
								<!--
                  Heroicon name: outline/home
  
                  Current: "text-gray-500", Default: "text-gray-400 group-hover:text-gray-500"
                -->
								<Icon name={item.icon} class="text-indigo-300 mr-3 flex-shrink-0 h-6 w-6" />
								{item.name}
							</a>
						{/each}
					</nav>
				</div>
				<div class="flex-shrink-0 flex border-t border-bubbly-darker p-4">
					<!-- <a href="" class="flex-shrink-0 w-full group block"> -->
					<div class="flex items-center">
						<div>
							<Icon name="user" class="h-9 w-9 text-gray-300" />
						</div>
						<div class="ml-3">
							<p class="text-sm font-medium text-white group-hover:text-gray-900">Annonymous</p>
							<!-- <p class="text-xs font-medium text-gray-500 group-hover:text-gray-700">
									View profile
								</p> -->
						</div>
					</div>
					<!-- </a> -->
				</div>
			</div>
		</div>
	</div>
	<div class="flex flex-col w-0 flex-1 overflow-hidden">
		<div class="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3">
			<button
				type="button"
				class="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
			>
				<span class="sr-only">Open sidebar</span>
				<!-- Heroicon name: outline/menu -->
				<svg
					class="h-6 w-6"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 6h16M4 12h16M4 18h16"
					/>
				</svg>
			</button>
		</div>
		<main class="flex-1 relative z-0 overflow-y-auto focus:outline-none">
			<slot />
		</main>
	</div>
</div>
