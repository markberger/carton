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
            '$http',
            '$upload',
            'ngDialog',
            function(
                $scope,
                $http,
                $upload,
                ngDialog
            ) {

                var filesCtrl = this;
                $scope.selected = null;

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
                            $scope.apiGetFiles();
                        }).error(function(data, status, headers, config) {
                            console.log(data);
                        });
                    }
                }

                $scope.setSelected = function(file) {
                    if ($scope.selected === file) {
                        $scope.selected = null;
                    } else {
                        $scope.selected = file;
                    }
                }

                $scope.deleteSelected = function() {
                    var q = ngDialog.openConfirm({
                        template: 'partials/confirmDelete.html',
                        className: 'ngdialog-theme-default',
                        scope: $scope
                    });

                    q.then(function() {
                        $scope._deleteSelected();
                    })
                }

                $scope._deleteSelected = function() {
                    var hash = $scope.selected.hash;
                    $http.delete('/api/files/'+hash)

                    .success(function(data, status, headers, config) {
                        $scope.selected = null;
                        $scope.apiGetFiles();
                    })

                    .error(function(data, status, headers, config) {
                        console.log('failed to delete file');
                    })
                }

                $scope.apiGetFiles = function() {
                    $http.get('/api/files')

                    .success(function(data, status, headers, config) {
                        filesCtrl.files = data;
                    })

                    .error(function(data, status, headers, config) {
                        filesCtrl.files = {};
                    })
                }

                $scope.getFiles = function() {
                    return filesCtrl.files;
                };

                $scope.download = function(file) {
                    var path = '/api/files/' + file.hash;
                    window.open(path, '_blank', '');
                }

                $scope.apiGetFiles();
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
