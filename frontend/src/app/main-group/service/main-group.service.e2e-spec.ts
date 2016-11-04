describe('MainGroup Service e2e',() => {
	it('should click something', () => {
		browser.get('http://localhost:4200');

		pexpect(browser.getTitle()).toEqual('Places4All');
	});
});
