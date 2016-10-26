import { FrontendPage } from './app.po';

describe('frontend App', function() {
  let page: FrontendPage;

  beforeEach(() => {
    page = new FrontendPage();
  });

  it('should display message saying p4a works!', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('p4a works!');
  });
});
