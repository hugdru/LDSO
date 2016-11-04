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
import {
	AccessibilityService
} from 'accessibility/service/accessibility.service';
import { HandlerService } from 'handler.service';

const mockArray = [
	{_id: 5, name: "carlos", weight: 25, criterion: 25},
	{_id: 2, name: "pedro", weight: 30, criterion: 25}
];
const mock = {_id: 5, name: "carlos", weight: 25, criterion: 25};

describe('Accessibility Service w/ Mock Server', () => {
	let mockBackend: MockBackend;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			providers: [
				AccessibilityService,
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

	it('Get all accessibilities', async(() => {
		let accessibilityService: AccessibilityService = getTestBed()
				.get(AccessibilityService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mockArray
			})));
		});

		accessibilityService.getAccessibilities().subscribe((data) => {
			expect(data).toBe(mockArray);
		});
	}));

	it('Get one accessibility', async(() => {
		let accessibilityService: AccessibilityService = getTestBed()
				.get(AccessibilityService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		accessibilityService.getAccessibility("name", "string", "clarlos")
				.subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));


	it('Update an accessibility', async(() => {
		let accessibilityService: AccessibilityService = getTestBed()
				.get(AccessibilityService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				status: 200
			})));
		});

		accessibilityService.updateAccessibility(5, "name", "string", "henrique")
				.subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(200);
		}));
	}));

	it('Add a new accessibiliy',
			async(inject([AccessibilityService], (accessibilityService) => {
		mockBackend.connections.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Post);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 201})));
		});
		accessibilityService.setAccessibility(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

	it('Delete an accessibility',
			async(inject([AccessibilityService], (accessibilityService) => {
		mockBackend.connections
				.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Get);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 204})));
		});

		accessibilityService.removeAccessibility(20).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(204);
		}));
	})));
});
