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
import {MainGroupService} from "main-group/service/main-group.service";
import {HandlerService} from "handler.service";

const mockArray = [
    {_id: 26, name: "ana", weight: 30},
    {_id: 14, name: "joao", weight: 25}
];
const mock = {_id: 26, name: "ana", weight: 30};

describe('MainGroup Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                MainGroupService,
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
        let groupService: MainGroupService = getTestBed().get(MainGroupService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getMainGroups().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some main groups', async(() => {
        let groupService: MainGroupService = getTestBed().get(MainGroupService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        groupService.getSomeMainGroups("weight", "string", "ana").subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get one main group', async(() => {
        let groupService: MainGroupService = getTestBed().get(MainGroupService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        groupService.getMainGroup("weight", "int", 1).subscribe((data) => {
            expect(data).toBe(mock);
        });
    }));

    it('Updade a main group', async(() => {
        let groupService: MainGroupService = getTestBed().get(MainGroupService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        groupService.updateMainGroup(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new main group',
            async(inject([MainGroupService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                groupService.setMainGroup(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a main group',
            async(inject([MainGroupService], (groupService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                groupService.removeMainGroup(15).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
