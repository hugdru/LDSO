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

describe('Group Service', () => {
	let mockBackend: MockBackend;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			providers: [
				GroupService,
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

		TestBed.compileComponents();
		mockBackend = getTestBed().get(MockBackend);
	}));

	it('should receive groups', async(() => {
		let fake = [{ _id: 26, name: "ana", weight: 30}];
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
});
