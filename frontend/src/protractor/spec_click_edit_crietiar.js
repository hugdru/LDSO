// spec_click_edit_crietiar.js
describe('Test if i click "Editar Criterios" in sidebar change the centarbar ', function() {
  it('should click "Editar Criterios"', function() {
    browser.get('http://localhost:4200/main-group');

    //shearch the link with the class active and click
    element(by.linkText('Editar Criterios')).click();
    expect(browser.getCurrentUrl()).toEqual("http://localhost:4200/main-group");
  });
});
