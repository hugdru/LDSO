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
import {AuditTemplateService} from "./audit-template.service";
import {HandlerService} from "../../shared/service/handler.service";

const mockArray = [
    {id: 26, name: "ana"},
    {id: 14, name: "joao"}
];
const mock = {id: 26, name: "ana"};

describe('AuditTemplate Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                AuditTemplateService,
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

    it('Get all main groups', async(() => {
        let groupService: AuditTemplateService = getTestBed().get(AuditTemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getAuditTemplates().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some main groups', async(() => {
        let groupService: AuditTemplateService = getTestBed().get(AuditTemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getSomeAuditTemplates("name", "ana").subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get one main group', async(() => {
        let groupService: AuditTemplateService = getTestBed().get(AuditTemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        groupService.getAuditTemplate(26).subscribe((data) => {
            expect(data).toBe(mock);
        });
    }));

    it('Updade a main group', async(() => {
        let groupService: AuditTemplateService = getTestBed().get(AuditTemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        groupService.updateAuditTemplate(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new main group',
            async(inject([AuditTemplateService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                groupService.setAuditTemplate(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a main group',
            async(inject([AuditTemplateService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                groupService.removeAuditTemplate(15).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
