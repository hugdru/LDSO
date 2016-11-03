import {
	TestBed,
	getTestBed,
	async,
	inject
} from '@angular/core/testing';

import {
	Headers, BaseRequestOptions,
	Response, HttpModule, Http, XHRBackend, RequestMethod
} from '@angular/http';

import { ResponseOptions } from '@angular/http';
import { MockBackend, MockConnection } from '@angular/http/testing';
import { GroupService } from 'group/service/group.service';
import { GroupComponent } from 'group/group.component';
import { HandlerService } from 'handler.service';

describe('Group Service', () => {
	let mockBackend: MockBackend;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			providers: [
				GroupService,
				HandlerService,
				MockBackend,
				BaseRequestOptions,
				{
					provide: Http,
					deps: [MockBackend, BaseRequestOptions],
					useFactory: (backend: XHRBackend,
								 defaultOptions: BaseRequestOptions) => {
						return new Http(backend, defaultOptions);
					},
				}
			],
			imports: [ HttpModule ]
		});

		mockBackend = getTestBed().get(MockBackend);
	}));

	it('should receive groups', async(() => {
		let fake = [{ _id: 26, name: "ana", weight: 30},
				{_id: 14, name: "joao", weight: 25}];
		let groupService: GroupService = getTestBed().get(GroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: fake
			})));
		});

		groupService.getGroups().subscribe((data) => {
			expect(data).toBe(fake);
		});
	}));

	// it('should insert a new group',
	// 		async(inject([GroupService], (groupService) => {
	// 	mockBackend.connections
	// 			.subscribe((connection: MockConnection) => {
	// 		expect(connection.request.method).toBe(RequestMethod.Post);
	// 		connection.mockRespond(new Response(
	// 				new ResponseOptions({status: 201})));
	// 	});

	// 	let fake = {_id: 24, name: "ana", weight: 20};
	// 	groupService.setGroup(fake).subscribe((result => {
	// 		expect(result).toBeDefined();
	// 		expect(result.status).toBe(201);
	// 	}));
	// })));
});
