'use strict';

angular.module('app.services', ['restangular']);
angular.module('app.controllers', ['app.services']);
angular.module('app', ['ngSanitize', 'ui.router', 'ui.bootstrap',
    'restangular', 'app.services', 'app.controllers', 'cgBusy',
]).config(function(RestangularProvider) {
    RestangularProvider.setBaseUrl('http://104.196.156.136:8000/v1');
    RestangularProvider.setDefaultHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    });  
});