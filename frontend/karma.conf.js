// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
    config.set({
        basePath: '',
        frameworks: ['jasmine', 'angular-cli'],
        plugins: [
            require('karma-jasmine'),
            require('karma-chrome-launcher'),
            require('karma-firefox-launcher'),
            require('karma-remap-istanbul'),
            require('angular-cli/plugins/karma'),
            require('karma-jasmine-html-reporter'),
            require('karma-story-reporter'),
            require('karma-phantomjs-launcher')
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
        reporters: ['progress', 'karma-remap-istanbul', 'story', 'kjhtml'],
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
