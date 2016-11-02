import { SubGroup } from 'sub-group/sub-group';

export class Group {
	_id: number;
	name: string;
	weight: number;
	sub_groups: SubGroup[];
}
