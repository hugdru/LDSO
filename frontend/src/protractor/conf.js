// conf.js
exports.config = {
  framework: 'jasmine',
  seleniumAddress: 'http://localhost:4444/wd/hub',
  //rootElement: '../app',
  useAllAngular2AppRoots: true ,
  specs: [
    'spec_title.js',
    'spec_click_addeditdelete.js',
    'spec_click_addeditdeleteSub.js'
     ],
  multiCapabilities: [
      { browserName: 'firefox'}
  //, { browserName: 'chrome'}
  ]
}
