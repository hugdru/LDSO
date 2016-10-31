import { GroupService } from './group.service';

describe('GroupService without the TestBed', () => {
	let service: GroupService;


	it('#getGroups should return faked value from a fake object', () => {
		const fake =  {
			getGroups: () => [{name: "jose", weight: 30, sub_groups: null}]
		};
		service = new GroupService(fake as GroupService);
		expect(service.getValue())
				.toBe([{name: "jose", weight: 30, sub_groups: null}]);
	});

	  // it('#getValue should return stubbed value from a FancyService spy', () => {
		// const fancy = new FancyService();
		// const stubValue = 'stub value';
		// const spy = spyOn(fancy, 'getValue').and.returnValue(stubValue);
		// service = new DependentService(fancy);
		// expect(service.getValue()).toBe(stubValue, 'service returned stub value');
		// expect(spy.calls.count()).toBe(1, 'stubbed method was called once');
		// expect(spy.calls.mostRecent().returnValue).toBe(stubValue);
	  // });
});
