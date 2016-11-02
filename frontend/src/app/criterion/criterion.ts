import { Accessibility } from 'accessibility/accessibility';

export class Criterion {
	id: number;
	name: string;
	weight: number;
	accessibility: Accessibility[];
	legislation: string;
}
