import {
	Http,
	BaseRequestOptions,
	Response,
	ResponseOptions,
	RequestMethod
} from '@angular/http';
import { MockBackend, MockConnection } from '@angular/http/testing';

import { inject, TestBed } from '@angular/core/testing';

import { GroupService } from 'group/service/group.service';
import { GroupComponent } from 'group/group.component';

describe('GroupServiceTest', () => {
	let service: GroupService = null;
	let backend: MockBackend = null;

	beforeEach(
		TestBed.configureTestingModule({
		declarations: [
			GroupComponent
		],
		imports: [
		  // HttpModule, etc.
		],
		providers: [
		  // { provide: ServiceA, useClass: TestServiceA }
		]
	  });	
	);
		// inject([GroupService, MockBackend],
		// 	(groupService: GroupService, mockBackend: MockBackend) => {
		// service = groupService;
		// backend = mockBackend;
	// }));

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
