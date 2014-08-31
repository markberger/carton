(function() {
    'use strict';

    angular.module('carton', [
        'ui.router',
        'angularFileUpload',
        'ngDialog',
        'carton.controllers',
        'carton.services'
    ]).
    config([
        '$stateProvider',
        '$urlRouterProvider',
        function(
            $stateProvider,
            $urlRouterProvider
        ) {
            $urlRouterProvider.otherwise('/');
            $stateProvider.
            state('files', {
                url: '/',
                templateUrl: 'partials/files.html',
                controller: 'FilesCtrl',
                data: {
                    access: {
                        isFree: false
                    }
                }
            }).

            state('login', {
                url: '/login',
                templateUrl: 'partials/login.html',
                controller: 'LoginCtrl',
                data: {
                    access: {
                        isFree: true
                    }
                }
            }).

            state('register', {
                url: '/register',
                templateUrl: 'partials/register.html',
                controller: 'RegisterCtrl',
                data: {
                    access: {
                        isFree: true
                    }
                }
            });

        }
    ])

    .run(['$rootScope', '$state', 'UserService',
        function($root, $state, userSrv) {
            $root.$on(
                '$stateChangeSuccess',
                function(event, toState, toParams, fromState, fromParams) {
                    if (!userSrv.isLogged &&
                        !toState.data.access.isFree) {
                        event.preventDefault();
                        $state.go('login')
                    }
                }
            )
        }
    ]);
})();