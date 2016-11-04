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
import { CriterionService } from 'criterion/service/criterion.service';
import { HandlerService } from 'handler.service';

describe('Criterion Service', () => {
	let mockBackend: MockBackend;

	beforeEach(async(() => {
		TestBed.configureTestingModule({
			providers: [
				CriterionService,
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

	it('should get criterias from Criterionservice', async(() => {
		let mock = [
			{
				_id: 5,
				name: "carlos",
				weight: 25,
				legislation: "aaa",
				sub_group: 25
			},
			{
				_id: 2,
				name: "pedro",
				weight: 30,
				legislation: "abba",
				sub_group: 33
			}];
		let criterionService: CriterionService = getTestBed()
				.get(CriterionService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		criterionService.getCriteria().subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('should get a criterion from criterion service using tags', async(() => {
		let mock =
			{
				_id: 5,
				name: "carlos",
				weight: 25,
				legislation: "aaa",
				sub_group: 25
			};
		let criterionService: CriterionService = getTestBed()
				.get(CriterionService);

		mockBackend.connections.subscribe((connection: MockConnection) => {
			connection.mockRespond(new Response(new ResponseOptions({
				body: mock
			})));
		});

		criterionService.getCriterion("name", "string", "clarlos")
				.subscribe((data) => {
			expect(data).toBe(mock);
		});
	}));

	it('should insert a new criterion',
			async(inject([CriterionService], (criterionService) => {
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
				legislation: "aaa",
				sub_group: 25
			};
		criterionService.setCriterion(mock).subscribe((result => {
			expect(result).toBeDefined();
			expect(result.status).toBe(201);
		}));
	})));

});
