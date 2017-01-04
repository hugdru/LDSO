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
import {LoginService} from "login/service/login.service";
import {HandlerService} from "../../shared/service/handler.service";

const mock = {id: 5, username: "carlos", password: "bla"};

describe('Session Service w/ Mock Service', () => {
    let mockBackend: MockBackend;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                LoginService,
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

    it('Add a new login',
            async(inject([LoginService], (loginService) => {
                mockBackend.connections
                        .subscribe((connection: MockConnection) => {
                            expect(connection.request.method).toBe(RequestMethod.Post);
                            connection.mockRespond(new Response(
                                    new ResponseOptions({status: 201})));
                        });

                loginService.setLogin(mock).subscribe((result => {
                    expect(result).toBeDefined();
                    expect(result.status).toBe(201);
                }));
            })));

    it('Get logout', async(() => {
        let loginService: LoginService = getTestBed().get(LoginService);

        mockBackend.connections.subscribe((connection: MockConnection) => {
            connection.mockRespond(new Response(
                    new ResponseOptions({status: 200})));
        });

        loginService.getLogout().subscribe((response => {
            expect(response).toBeDefined();
            expect(response.status).toBe(200);
        }));
    }));

});
