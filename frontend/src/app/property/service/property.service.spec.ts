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
import {PropertyService} from "property/service/property.service";
import {HandlerService} from "handler.service";

const mockArray = [
    {id: 26, name: "ana", image_path: "bla"},
    {id: 14, name: "joao", image_path: "ble"}
];
const mock = {id: 26, name: "ana", image_path: "bla"};

describe('Property Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                PropertyService,
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

    it('Get all properties', async(() => {
        let propertyService: PropertyService = getTestBed().get(PropertyService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        propertyService.getProperties().subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get some properties', async(() => {
        let propertyService: PropertyService = getTestBed().get(PropertyService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mockArray
            })));
        });

        propertyService.getSomeProperties("name", "ana").subscribe((data) => {
            expect(data).toBe(mockArray);
        });
    }));

    it('Get one property', async(() => {
        let propertyService: PropertyService = getTestBed().get(PropertyService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(new ResponseOptions({
                body: mock
            })));
        });

        propertyService.getProperty(26).subscribe((data) => {
            expect(data).toBe(mock);
        });
    }));

    it('Updade a property', async(() => {
        let propertyService: PropertyService = getTestBed().get(PropertyService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            expect(connection.request.method).toBe(RequestMethod.Put);
            connection.mockRespond(new Response(new ResponseOptions({
                status: 200
            })));
        });

        propertyService.updateProperty(mock).subscribe((result => {
            expect(result).toBeDefined();
            expect(result.status).toBe(200);
        }));
    }));

    it('Add a new property',
            async(inject([PropertyService], (propertyService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                propertyService.setProperty(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Delete a property',
            async(inject([PropertyService], (propertyService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Delete);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 204})));
                        });

                propertyService.removeProperty(15).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(204);
                }));
            })));
});
