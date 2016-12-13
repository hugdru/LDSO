// spec_title.js
describe('Check if the title is Places4All', function() {
  it('should have a title', function() {
    browser.get('http://lp4a.tk:8080');
    expect(browser.getTitle()).toEqual('Places4All');
  });
});
