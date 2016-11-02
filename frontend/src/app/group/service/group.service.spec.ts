import {
	Http,
	BaseRequestOptions,
	Response,
	ResponseOptions,
	RequestMethod
} from '@angular/http';
import { MockBackend, MockConnection } from '@angular/http/testing';
import { inject } from '@angular/core/testing';

import { GroupService } from 'group/service/group.service';

describe('GroupServiceTest', () => {
	let service: GroupService = null;
	let backend: MockBackend = null;

	beforeEach(inject([GroupService, MockBackend],
			(groupService: GroupService, mockBackend: MockBackend) => {
		service = groupService;
		backend = mockBackend;
	}));

	it('#getGroups', (done) => {
		let fake: Object[] = [{_id: 1, name: "Casa", weight: 30, sub_groups: null}];
		backend.connections.subscribe((connection: MockConnection) => {
			let options = new ResponseOptions({
				body: JSON.stringify(fake)
			});
		connection.mockRespond(new Response(options));
		});

		service.getGroups().subscribe((response) => {
			expect(response).toEqual(fake);
			done();
		});
	});
});
