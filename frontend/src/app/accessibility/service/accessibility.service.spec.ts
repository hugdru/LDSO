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

describe('Accessibility Service', () => {
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

	it('should get accessibilities from accessibilityservice', async(() => {
		let mock = [
			{
				_id: 5,
				name: "carlos",
				weight: 25,
				criterion: 25
			},
			{
				_id: 2,
				name: "pedro",
				weight: 30,
				criterion: 25
			}];
		let accessibilityService: AccessibilityService = getTestBed()
				.get(AccessibilityService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		accessibilityService.getAccessibilities().subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('should get a accessibilities from accessibilities service using tags',
			async(() => {
		let mock =
			{
				_id: 5,
				name: "carlos",
				weight: 25,
				criterion: 25
			};
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

	it('should insert a new accessibility',
			async(inject([AccessibilityService], (accessibilityService) => {
		mockBackend.connections.subscribe((connection: MockConnection) => {
			expect(connection.request.method).toBe(RequestMethod.Post);
			connection.mockRespond(new Response(
					new ResponseOptions({status: 201})));
		});

		let mock =
			{
				_id: 5,
				name: "carlos",
				weight: 25,
				criterion: 25
			};
		accessibilityService.setAccessibility(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

});
