// Controllers

(function() {
    'use strict';
    angular.module('carton.controllers', [])
        .controller('loginCtrl', ['$scope',
            '$state',
            '$http',
            'UserService',
            function(
                $scope,
                $http,
                User
            ) {
                $scope.login = function () {
                    $http({method: 'GET', url: '/auth/status'})
                        .success(function(data, status, headers, config) {
                            if (data.status) {
                                User.isLogged = true;
                            } else {
                                User.isLogged = false;
                            }
                        })
                        .error(function(data, status, headers, config) {
                            User.isLogged = false;
                        });
                }
            }
        ])

        .controller('NavController', ['$scope',
            'UserService',
            function(
                $scope,
                userSrv)
            {
                $scope.isLogged = function() {
                    return userSrv.isLogged;
                }
            }
        ]);
})();
