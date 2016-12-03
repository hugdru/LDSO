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
import {CtemplateService} from "ctemplate/service/ctemplate.service";
import {HandlerService} from "handler.service";

const mockArray = [
    {id: 26, name: "ana"},
    {id: 14, name: "joao"}
];
const mock = {id: 26, name: "ana"};

describe('Ctemplate Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                CtemplateService,
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
        let groupService: CtemplateService = getTestBed().get(CtemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getCtemplates().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some main groups', async(() => {
        let groupService: CtemplateService = getTestBed().get(CtemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getSomeCtemplates("name", "ana").subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get one main group', async(() => {
        let groupService: CtemplateService = getTestBed().get(CtemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        groupService.getCtemplate(26).subscribe((data) => {
            expect(data).toBe(mock);
        });
    }));

    it('Updade a main group', async(() => {
        let groupService: CtemplateService = getTestBed().get(CtemplateService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        groupService.updateCtemplate(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new main group',
            async(inject([CtemplateService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                groupService.setCtemplate(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a main group',
            async(inject([CtemplateService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                groupService.removeCtemplate(15).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
