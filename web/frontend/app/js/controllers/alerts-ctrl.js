'use strict';

angular
    .module('app.controllers')
    .controller('AlertsCtrl', function($scope, Restangular, $location) {
        $scope.loadingPromise = Restangular.one('get_all_articles')
            .get()
            .then(function(data) {
                var articles = [];
                $scope.alerts = data;
                for (var company in data) {
                    if (company && data[company]) {
                        for (var article in data[company].HREFs) {
                            if (data[company].HREFs[article]) {
                                data[company].HREFs[article].Company = data[company].Keyword;
                                articles.push(data[company].HREFs[article]);
                            }
                        }
                    }
                }
                $scope.articles = articles;
            });
    });