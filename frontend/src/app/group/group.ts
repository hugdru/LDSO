import { SubGroup } from 'sub-group/sub-group';

export class Group {
	id: number;
	name: string;
	weight: number;
	sub_groups: SubGroup[];
}
