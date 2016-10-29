import { Property } from './property/property';
import { Note } from './note/note';

export class Audit {
	property: Property;
	notes: Note[];
	value: number;
}
