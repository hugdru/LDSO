import {TestBed, getTestBed, async, inject} from "@angular/core/testing";
import {
    BaseRequestOptions,
    Response,
    HttpModule,
    Http,
    XHRBackend,
    RequestMethod,
    ResponseOptions
} from "@angular/http";
import {MockBackend, MockConnection} from "@angular/http/testing";
import {AccessibilityService} from "accessibility/service/accessibility.service";
import {HandlerService} from "handler.service";

const mockArray = [
    {id: 5, name: "carlos", weight: 25},
    {id: 2, name: "pedro", weight: 30}
];
const mock = {id: 5, name: "carlos", weight: 25};

const mock_criterion_id = 4;

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
            imports: [HttpModule]
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

    it('Get some accessibilities', async(() => {
        let accessibilityService: AccessibilityService = getTestBed()
                .get(AccessibilityService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        accessibilityService.getSomeAccessibilities(12)
                .subscribe((data) => {
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

        accessibilityService.getAccessibility(5)
                .subscribe((data) => {
                    expect(data).toBe(mock);
                });
    }));


    it('Update an accessibility', async(() => {
        let accessibilityService: AccessibilityService = getTestBed()
                .get(AccessibilityService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        accessibilityService.updateAccessibility(mock, mock_criterion_id)
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
                accessibilityService.setAccessibility(mock, mock_criterion_id)
                        .subscribe((result => {
                            expect(result).toBeDefined();
                            expect(result.status).toBe(201);
                        }));
            })));

    it('Delete an accessibility',
            async(inject([AccessibilityService], (accessibilityService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                accessibilityService.removeAccessibility(20, 1)
                        .subscribe((result => {
                            expect(result).toBeDefined();
                            expect(result.status).toBe(204);
                }));
            })));
});
