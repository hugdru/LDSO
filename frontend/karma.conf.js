// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
  config.set({
      // base path that will be used to resolve all patterns (eg. files, exclude)
      basePath: '',
      // frameworks to use
      // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
      frameworks: ['jasmine', 'angular-cli'],
      plugins: [
        require('karma-jasmine'),
        require('karma-chrome-launcher'),
        require('karma-firefox-launcher'),
        require('karma-phantomjs-launcher'),
        require('karma-remap-istanbul'),
        require('angular-cli/plugins/karma')
      ],
      files: [
        { pattern: './src/test.ts', watched: false }
      ],
      preprocessors: {
        './src/test.ts': ['angular-cli']
      },
      remapIstanbulReporter: {
        reports: {
          html: 'coverage',
          lcovonly: './coverage/coverage.lcov'
        }
      },
      angularCli: {
        config: './angular-cli.json',
        environment: 'dev'
      },
      reporters: ['progress', 'karma-remap-istanbul'],
      port: 9876,
      colors: true,
      logLevel: config.LOG_INFO,
      // enable / disable watching file and executing tests whenever any file changes
      autoWatch: false,
      // There is a problem with chrome in travis, Chrome have not captured in 60000 ms, killing.
      browsers: ["Chrome", "Firefox"],
      // Continuous Integration mode
      // if true, Karma captures browsers, runs the tests and exits
      singleRun: false
  });

  if (process.env.TRAVIS) {
    config.browsers = ['PhantomJS'];
  }

};
