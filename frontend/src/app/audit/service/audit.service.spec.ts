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
import {AuditService} from "audit/service/audit.service";
import {HandlerService} from "../../shared/service/handler.service";

const mockArray = [
    {
        id: 26, property: 10, rating: 85,
        criteria: [{criterion: 1, rating: 10}, {criterion: 2, rating: 25}]
    },
    {
        id: 26, property: 12, rating: 76,
        criteria: [{criterion: 10, rating: 10}, {criterion: 2, rating: 25}]
    }
];
const mock = {
    id: 26, property: 10, rating: 85,
    criteria: [{criterion: 1, rating: 10}, {criterion: 2, rating: 25}]
};


describe('Audit Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                AuditService,
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

    it('Get all audits', async(() => {
        let auditService: AuditService = getTestBed().get(AuditService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        auditService.getAudits().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some audits', async(() => {
        let auditService: AuditService = getTestBed().get(AuditService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        auditService.getSomeAudits("weight", "ana").subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get one audit', async(() => {
        let auditService: AuditService = getTestBed().get(AuditService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        auditService.getAudit(26).subscribe((data) => {
            expect(data).toBe(mock);
        });
    }));

    it('Updade a audit', async(() => {
        let auditService: AuditService = getTestBed().get(AuditService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        auditService.updateAudit(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new audit',
            async(inject([AuditService], (auditService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                auditService.setAudit(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a audit',
            async(inject([AuditService], (auditService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                auditService.removeAudit(15).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
