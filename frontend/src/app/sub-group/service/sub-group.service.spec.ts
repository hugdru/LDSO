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
import { SubGroupService } from 'sub-group/service/sub-group.service';
import { HandlerService } from 'handler.service';

const mockArray = [
	{_id: 5, name: "carlos", weight: 25, main_group: 25},
	{_id: 2, name: "pedro", weight: 30, main_group: 14}
];
const mock = {_id: 5, name: "carlos", weight: 25, main_group: 25};

describe('SubGroup Service w/ Mock Service', () => {
	let mockBackend: MockBackend;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			providers: [
				SubGroupService,
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

	it('Get all subgroups', async(() => {
		let subGroupService: SubGroupService = getTestBed().get(SubGroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mockArray
			})));
		});

		subGroupService.getSubGroups().subscribe((data) => {
			expect(data).toBe(mockArray);
		});
	}));

	it('Get some subgroups', async(() => {
		let subGroupService: SubGroupService = getTestBed().get(SubGroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mockArray
			})));
		});

		subGroupService.getSomeSubGroups("main_group", "int", 5)
				.subscribe((data) => {
			expect(data).toBe(mockArray);
		});
	}));

	it('Get one subgroup', async(() => {
		let subGroupService: SubGroupService = getTestBed().get(SubGroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		subGroupService.getSubGroup("name", "string", "clarlos")
				.subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('Update a subgroup',
	   		async(() => {
		let subGroupService: SubGroupService = getTestBed()
				.get(SubGroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Put);
			connection.mockRespond(new Response(new ResponseOptions({
				status: 200
			})));
		});

		subGroupService.updateSubGroup(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(200);
		}));
	}));

	it('Add a new subgroup',
			async(inject([SubGroupService], (subGroupService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Post);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 201})));
		});

		subGroupService.setSubGroup(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

	it('Delete a subgroup',
			async(inject([SubGroupService], (subGroupService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Delete);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 204})));
		});

		subGroupService.removeSubGroup(20).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(204);
		}));
	})));
});
