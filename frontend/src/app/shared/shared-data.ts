import {environment} from "../../environments/environment";

export const serverUrl = environment.api;

export const mainGroupsUrl = serverUrl + '/maingroups';

export const subGroupsUrl = serverUrl + '/subgroups';

export const criteriaUrl = serverUrl + '/criteria';

export const accessibilitiesUrl = criteriaUrl + '/#/accessibilities';

export const propertiesUrl = serverUrl + '/properties';

export const auditsUrl = serverUrl + '/audits';
