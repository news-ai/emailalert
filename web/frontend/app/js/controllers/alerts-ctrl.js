'use strict';

angular
    .module('app.controllers')
    .controller('AlertsCtrl', function($scope, Restangular, $location) {
        $scope.loadingPromise = Restangular.one('get_all_articles')
            .get()
            .then(function(data) {
                var articles = [];
                for (var company in data) {
                    if (company && data[company]) {
                        for (var article in data[company].HREFs) {
                            if (data[company].HREFs[article] && !data[company].HREFs[article].status) {
                                data[company].HREFs[article].Company = data[company].Keyword;
                                data[company].HREFs[article].CompanyId = data[company].Id;
                                articles.push(data[company].HREFs[article]);
                            } else {
                                console.log(data[company].HREFs[article]);
                            }
                        }
                    }
                }
                $scope.articles = articles;
                $scope.currentArticle = articles[0];
            });
        $scope.removeArticle = function (id, url) {
            Restangular.one('article_status/' + id + '?url=' + window.encodeURIComponent(url)).get().then(function(response){
                $scope.articles.pop();
                $scope.currentArticle = $scope.articles[0];
            });
        }
    });