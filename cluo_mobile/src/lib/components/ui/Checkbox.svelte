<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	interface Props {
		id?: string;
		name?: string;
		checked?: boolean;
		disabled?: boolean;
		readonly?: boolean;
		required?: boolean;
		class?: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let {
		id,
		name,
		checked = false,
		disabled = false,
		readonly = false,
		required = false,
		class: className = '',
		size = 'md',
		...restProps
	}: Props = $props();

	const dispatch = createEventDispatcher();

	const baseClasses = 'rounded border-gray-300 bg-background text-primary focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50';

	const sizeClasses = {
		sm: 'h-4 w-4',
		md: 'h-4 w-4',
		lg: 'h-5 w-5'
	};

	const checkboxClass = `${baseClasses} ${sizeClasses[size]} ${className}`;

	function handleChange(event: Event) {
		const target = event.target as HTMLInputElement;
		dispatch('change', { checked: target.checked });
	}
</script>

<input
	type="checkbox"
	{id}
	{name}
	{checked}
	{disabled}
	{readonly}
	{required}
	class={checkboxClass}
	onchange={handleChange}
	{...restProps}
/>