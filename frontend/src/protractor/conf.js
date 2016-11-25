// conf.js
exports.config = {
  framework: 'jasmine',
  seleniumAddress: 'http://localhost:4444/wd/hub',
  //rootElement: '../app',
  useAllAngular2AppRoots: true ,
  specs: [
    'spec_title.js',
    'spec_click_edit_crietiar.js',
    'spec_click_adicionar_novo.js'
     ],
  multiCapabilities: [
      { browserName: 'firefox'}
  //, { browserName: 'chrome'}
  ]
}
