// Controllers

(function() {
    'use strict';
    angular.module('carton.controllers', [])
        .controller('LoginCtrl', ['$scope',
            '$http',
            '$state',
            'UserService',
            function(
                $scope,
                $http,
                $state,
                userSrv
            ) {
                $scope.login = function (user) {
                    $http.post('/api/auth/login', user)

                    .success(function(data, status, headers, config) {
                        userSrv.isLogged = true;
                        $state.go('files');
                    })

                    .error(function(data, status, headers, config) {
                        userSrv.isLogged = false;
                    });
                }
            }
        ])

        .controller('RegisterCtrl', ['$scope',
            '$http',
            '$state',
            'UserService',
            function(
                $scope,
                $http,
                $state,
                userSrv
            ) {
                $scope.register = function(user) {
                    $http.post('/api/auth/register', user)

                    .success(function(data, status, headers, config) {
                        userSrv.isLogged = true;
                        $state.go('files');
                    })
                    .error(function(data, status, headers, config) {
                        console.log(data);
                        $state.go('register');
                    });
                }
            }
        ])

        .controller('NavController', ['$scope',
            '$http',
            '$state',
            'UserService',
            function(
                $scope,
                $http,
                $state,
                userSrv)
            {
                $scope.isLogged = function() {
                    return userSrv.isLogged;
                }

                $scope.logout = function() {
                    $http.post('/api/auth/logout')

                    .success(function(data, status, headers, config) {
                        userSrv.isLogged = false;
                        $state.go('login');
                    })
                }
            }
        ]);
})();
