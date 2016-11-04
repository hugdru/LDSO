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

describe('Sub-Group Service', () => {
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

	it('should get sub groups from sub group service', async(() => {
		let mock = [{_id: 5, name: "carlos", weight: 25, main_group: 25},
				{_id: 2, name: "pedro", weight: 30, main_group: 14}];
		let subGroupService: SubGroupService = getTestBed().get(SubGroupService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		subGroupService.getSubGroups().subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('should get one sub group from sub group service using tags', async(() => {
		let mock = {_id: 5, name: "carlos", weight: 25, main_group: 25};
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

	it('should insert a new sub group',
			async(inject([SubGroupService], (subGroupService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Post);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 201})));
		});

		let mock = {_id: 12, name: "rita", weight: 15, main_group: 45};
		subGroupService.setSubGroup(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

});
