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
import { HandlerService } from 'handler.service';

const mockArray = [
	{_id: 26, name: "ana", weight: 30},
	{_id: 14, name: "joao", weight: 25}
];
const mock = { _id: 26, name: "ana", weight: 30};

describe('Group Service w/ Mock Service', () => {
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

	it('Get all groups', async(() => {
		let groupService: GroupService = getTestBed().get(GroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mockArray
			})));
		});

		groupService.getGroups().subscribe((data) => {
			expect(data).toBe(mockArray);
		});
	}));

	it('Get one group', async(() => {
		let groupService: GroupService = getTestBed().get(GroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		groupService.getGroup("name", "int", 1).subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('Updade a group', async(() => {
		let groupService: GroupService = getTestBed().get(GroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				status: 200
			})));
		});

		groupService.updateGroup(5, "name", "string", "henrique")
				.subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(200);
		}));
	}));

	it('Add a new group',
			async(inject([GroupService], (groupService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Post);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 201})));
		});

		groupService.setGroup(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

	it('Delete a group',
			async(inject([GroupService], (groupService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Get);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 204})));
		});

		groupService.removeGroup(15).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(204);
		}));
	})));
});
