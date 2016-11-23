import {environment} from "../../environments/environment";

export const serverUrl = environment.api;

export const mainGroupsUrl = serverUrl + '/mainGroups';
export const mainGroupsFindUrl = mainGroupsUrl + '/find';

export const subGroupsUrl = serverUrl + '/subGroups';
export const subGroupsFindUrl = subGroupsUrl + '/find';

export const criteriaUrl = serverUrl + '/criteria';
export const criteriaFindUrl = criteriaUrl + '/find';

export const accessibilitiesUrl = serverUrl + '/accessibilities';
export const accessibilitiesFindUrl = accessibilitiesUrl + '/find';

export const propertiesUrl = serverUrl + '/properties';
export const propertiesFindUrl = propertiesUrl + '/find';

export const auditsUrl = serverUrl + '/audits';
export const auditsFindUrl = auditsUrl + '/find';
