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

        .controller('FilesCtrl', ['$scope',
            '$upload',
            function(
                $scope,
                $upload
            ) {
                $scope.onFileSelect = function($files) {
                    for (var i = 0; i < $files.length; i++) {
                        var file = $files[i];
                        $scope.upload = $upload.upload({
                            url: 'api/files',
                            method: 'POST',
                            file: file,
                        }).progress(function(evt) {
                            console.log('percent: ' + parseInt(100.0 * evt.loaded / evt.total));
                        }).success(function(data, status, headers, config) {
                            // file is uploaded successfully
                            console.log(data);
                        }).error(function(data, status, headers, config) {
                            console.log(data);
                        });
                    }
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
