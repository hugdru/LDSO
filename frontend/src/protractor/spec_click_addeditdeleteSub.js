// spec_click_addeditdeleteSub.js
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
    var MAING = element(by.className('acesso0'));
    var ADDSUBGROUP = element(by.className('acesso0')).element(by.xpath('//list-manage/div/button'));
    var NAMESUB = element(by.className('acesso0')).element(by.xpath('//list-manage/div')).element(by.id('adicionarMainGroup'));
    var WEIGHTSUB = element(by.className('acesso0')).element(by.xpath('//list-manage/div')).element(by.id('adicionarMainGroup'));
    var SCSUB = element(by.className('acesso0')).element(by.xpath('//list-manage/div')).element(by.className('btn btn-success round outline'));

    function add(a,b) {
      EC.click();
      AN.click();
      NAME.sendKeys(a);
      WEIGHT.sendKeys(b);
      SC.click();

    }
  function entreMainGroup(a,b) {
    MAING.click();
    ADDSUBGROUP.click();
    NAMESUB.sendKeys(a);
    WEIGHTSUB.sendKeys(b);
  }

    //for sucess
  beforeEach(function(){
      browser.get('http://lp4a.tk');
    });
  it('should click "Adicionar Acesso0 ao subgrupo "', function() {
    //var len = element(by.className('list')).getSize();
    add('acesso0','0');
    entreMainGroup('massaneta','0');
    //expect(len).toEqual(0);
  });


});

