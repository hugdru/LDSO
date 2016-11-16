export class AuditCriterion {
	criterion: number;
	rating: number;
}

export class Audit {
	_id: number;
	property: number;
	rating: number;
	criteria: AuditCriterion[];
}
