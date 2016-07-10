'use strict';

angular
    .module('app.controllers')
    .controller('AlertsCtrl', function($scope, Restangular, $location) {
        $scope.loadingPromise = Restangular.one('get_all_articles')
            .get()
            .then(function(data) {
                var articlesLeft = [];
                var articlesApproved = [];
                for (var company in data) {
                    if (company && data[company]) {
                        for (var article in data[company].HREFs) {
                            console.log(data[company].HREFs[article])
                            if (data[company].HREFs[article] && data[company].HREFs[article].approved) {
                                articlesApproved.push(data[company].HREFs[article]);
                            }
                            if (data[company].HREFs[article] && !data[company].HREFs[article].status && !data[company].HREFs[article].approved) {
                                data[company].HREFs[article].Company = data[company].Keyword;
                                data[company].HREFs[article].CompanyId = data[company].Id;
                                articles.push(data[company].HREFs[article]);
                            } else {
                                console.log(data[company].HREFs[article]);
                            }
                        }
                    }
                }
                $scope.articles = articlesLeft;
                $scope.currentArticle = articlesLeft[0];
                console.log(articlesApproved);
                $scope.articlesApproved = articlesApproved;
            });
        $scope.removeArticle = function (id, url) {
            Restangular.one('article_status/' + id + '?url=' + window.encodeURIComponent(url)).get().then(function(response){
                $scope.articles.pop();
                $scope.currentArticle = $scope.articles[0];
                window.location.reload()
            });
        }

        $scope.approveArticle = function (id, url, sentiment) {
            Restangular.one('article_approve/' + id + '/' + sentiment + '?url=' + window.encodeURIComponent(url)).get().then(function(response){
                $scope.articles.pop();
                $scope.currentArticle = $scope.articles[0];
                window.location.reload()
            });
        }
    });