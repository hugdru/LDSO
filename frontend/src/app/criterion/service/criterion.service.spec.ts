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
import {CriterionService} from "criterion/service/criterion.service";
import {HandlerService} from "../../shared/service/handler.service";

const mockArray = [
    {id: 5, name: "carlos", weight: 25, legislation: "aaa", idSubgroup: 25},
    {id: 2, name: "pedro", weight: 30, legislation: "abba", idSubgroup: 33}
];

const mock = {id: 5, name: "ana", weight: 25, legislation: "a", idSubgroup: 25};

describe('Criterion Service w/ Mock Server', () => {
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
            imports: [HttpModule]
        });

        mockBackend = getTestBed().get(MockBackend);
    }));

    it('Get all criteria', async(() => {
        let criterionService: CriterionService = getTestBed()
                .get(CriterionService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        criterionService.getCriteria().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some criteria', async(() => {
        let criterionService: CriterionService = getTestBed()
                .get(CriterionService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        criterionService.getSomeCriteria("idSubgroup", 4)
                .subscribe((data) => {
                    expect(data).toBe(mockArray);
                });
    }));

    it('Get one criterion', async(() => {
        let criterionService: CriterionService = getTestBed()
                .get(CriterionService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        criterionService.getCriterion(5)
                .subscribe((data) => {
                    expect(data).toBe(mock);
                });
    }));

    it('Update a criterion', async(() => {
        let criterionService: CriterionService = getTestBed()
                .get(CriterionService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        criterionService.updateCriterion(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new criterion',
            async(inject([CriterionService], (criterionService) => {
                mockBackend.connections.subscribe((connection: MockConnection) => {
                    expect(connection.request.method).toBe(RequestMethod.Post);
                    connection.mockRespond(new Response(
                            new ResponseOptions({status: 201})));
                });

                criterionService.setCriterion(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a criterion',
            async(inject([CriterionService], (criterionService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                criterionService.removeCriterion(20).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
