/* jshint undef: true, unused: true, strict: true */
/* global require, window, clearTimeout, document, Babylon*/

requirejs.config({
    baseUrl: "/js",
    paths: {
        babylon: "./lib/babylon.2.3.max"
    },
    shim : {
        'game': ['babylon']
    }
});

require(['game']);
