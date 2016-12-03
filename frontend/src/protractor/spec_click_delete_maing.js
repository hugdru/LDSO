// spec_click_edit_crietiar.js
describe('Test if i click "Adicionar Novo" in  centarbar ', function() {
    //go to sidebar link criteria
    var EC = element(by.linkText('Editar Criterios'));
    //click add "adicionar novo"
    var AN = element(by.buttonText('Adicionar Novo'));
    var SC =  element(by.className('btn btn-success round outline'));
    var SC =  element(by.className('btn btn-success round outline'));
    //
    var name =   element(by.className('ng-pristine ng-valid ng-touched') );
    var weight = element(by.className('ng-pristine ng-valid ng-touched') );

    function add(a,b) {
      EC.click();
      AN.click();
     //  name.sendKeys(a);
     //  weight.sendKeys(b);
      SC.click();

    }
   // name.sendKeys('Acesso2');
   // weight.sendKeys('0');

    //for sucess
  beforeEach(function(){
      browser.get('http://localhost:4200/main-group');
    });
  it('should click "Adicionar Novo"', function() {
    add('acesso0','0');

  });
});
