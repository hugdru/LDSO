import { Accessibility } from 'accessibility/accessibility';

export class Criterion {
	_id: number;
	name: string;
	weight: number;
	accessibility: Accessibility[];
	legislation: string;
}
