// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
  config.set({
      basePath: '',
      frameworks: ['jasmine', '@angular/cli'],
      plugins: [
          require('karma-jasmine'),
          require('karma-chrome-launcher'),
          require('karma-coverage-istanbul-reporter'),
          require('@angular/cli/plugins/karma')
          require('karma-firefox-launcher'),
          require('karma-phantomjs-launcher'),
          require('karma-jasmine-html-reporter'),
          require('karma-story-reporter')
      ],
      files: [
          { pattern: './src/test.ts', watched: false }
      ],
      preprocessors: {
          './src/test.ts': ['@angular-cli']
      },
      mime: {
        'text/x-typescript': ['ts','tsx']
      },
      coverageIstanbulReporter: {
        reports: [ 'html', 'lcovonly' ],
        fixWebpackSourcePaths: true
      },
      angularCli: {
          config: './angular-cli.json',
          environment: 'dev'
      },
      reporters: config.angularCli && config.angularCli.codeCoverage
            ? ['progress', 'coverage-istanbul']
            : ['progress'],
      port: 9876,
      colors: true,
      logLevel: config.LOG_INFO,
      autoWatch: true,
      browsers: ['Firefox'],
      singleRun: false
  });

  if (process.env.TRAVIS) {
      config.browsers = ['PhantomJS'];
  }
};
