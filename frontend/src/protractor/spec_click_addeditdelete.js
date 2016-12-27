// spec_click_edit_crietiar.js
describe('Test if i click "Adicionar/editar/delete Novo" in  centarbar ', function() {
    //go to sidebar link criteria
    var EC = element(by.linkText('Editar Criterios'));
    //click add "adicionar novo"
    var AN = element(by.buttonText('Adicionar Novo'));
    var SC =  element(by.className('btn btn-success round outline'));
    var NAME = element(by.id('adicionarMainGroup'));
    var WEIGHT = element(by.id('adicionarweight'));
    var EDBT = element(by.className('btn btn-neutral round outline'));
    var DELBT = element(by.className('btn btn-delete round outline'));



    function add(a,b) {
      EC.click();
      AN.click();
      NAME.sendKeys(a);
      WEIGHT.sendKeys(b);
      SC.click();

    }
  function editName(a,b) {
    EDBT.click();
    NAME.clear();
    NAME.sendKeys(a);
    WEIGHT.sendKeys(b);
    SC.click();

  }
  function deleteName() {
    DELBT.click();

  }


    //for sucess
  beforeEach(function(){
      browser.get('http://lp4a.tk');
    });
  it('should click "Adicionar Acesso0 "', function() {
    add('acesso0','0');


  },1000000);
  it('should click "Adicionar Acesso0 Acess02"', function() {
    add('acesso0','0');
    add('acesso2','10');

    //expect(len).toEqual(0);
  });
  it('should click "Edita Acesso0"', function() {
    //var len = element(by.className('list')).getSize();
    add('acesso0','0');
    editName('acesso1','1');

    //expect(len).toEqual(0);
  });
  it('should click "Delelte Acesso0"', function() {
    //var len = element(by.className('list')).getSize();
    add('acesso0','0');
    deleteName();

    //expect(len).toEqual(0);
  });
});

