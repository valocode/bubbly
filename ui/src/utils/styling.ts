// classNames provides a list of classes to use filtered on a boolean condition
// useful when changing the style of elements based on e.g. selected
export function classNames(...classes): string {
	return classes.filter(Boolean).join(' ');
}
