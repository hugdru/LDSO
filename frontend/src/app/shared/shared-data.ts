import {environment} from "../../environments/environment";

export const serverUrl = environment.api;

export const auditTemplatesUrl = serverUrl + '/templates';

export const currentAuditTemplateUrl = serverUrl + '/templates/current';

export const mainGroupsUrl = serverUrl + '/maingroups';

export const subGroupsUrl = serverUrl + '/subgroups';

export const criteriaUrl = serverUrl + '/criteria';

export const accessibilitiesUrl = criteriaUrl + '/#/accessibilities';

export const propertiesUrl = serverUrl + '/properties';

export const auditsUrl = serverUrl + '/audits';

export const auditsFindUrl = auditsUrl + '/find';

export const auditsSubGroupsUrl = auditsUrl + '/subgroups';

export const auditsCriterionUrl = auditsUrl + '/criterion';

export const imageUploadUrl = serverUrl + '/audits/#/criteria/!/remarks';

export const loginUrl = serverUrl + '/entities/login';

export const logoutUrl = serverUrl + '/entities/logout';

export const registerUrl = serverUrl + '/entities/register';
