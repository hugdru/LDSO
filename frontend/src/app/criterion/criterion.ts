import { Accessibility } from './accessibility/accessibility';

export class Criterion {
	name: string;
	weight: number;
	accessibility: Accessibility[];
	legislation: string;
}
