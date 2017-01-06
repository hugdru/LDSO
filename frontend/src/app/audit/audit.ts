export class AuditCriterion {
    criterion: number;
    rating: number;
}

export class Audit {
    id: number;
    property: number;
    rating: number;
    criteria: AuditCriterion[];
}

export class AuditSubgrups {
    idProperty: number;
    subgroups: number[];
}