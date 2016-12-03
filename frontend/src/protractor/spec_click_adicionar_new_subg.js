// spec_click_edit_crietiar.js
describe('Test if i click "Adicionar Novo" in  centarbar ', function() {
    //go to sidebar link criteria
    var EC = element(by.linkText('Editar Criterios'));
    //click add "adicionar novo"
    var AN = element(by.buttonText('Adicionar Novo'));
    var SC =  element(by.className('btn btn-success round outline'));
    //var NAME = element(by.model('selectedObject.name'));
   // var WEIGHT = element(by.id('adicionarweight'));

    function add(a,b) {
      EC.click();
      AN.click();
      //NAME.sendKeys(a);
     // WEIGHT.sendKeys(b);
      SC.click();

    }


    //for sucess
  beforeEach(function(){
      browser.get('http://lp4a.tk/main-group');
    });
  it('should click "Adicionar Novo"', function() {
    var len = element(by.className('list')).getSize();
    add('acesso0','0');

    expect(len).toEqual(0);
  });
});
