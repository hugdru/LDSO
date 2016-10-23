import { FrontendPage } from './app.po';

describe('frontend App', function() {
  let page: FrontendPage;

  beforeEach(() => {
    page = new FrontendPage();
  });

  it('should display message saying request handled by nginx', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('request handled by nginx');
  });
});
