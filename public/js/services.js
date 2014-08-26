// Services

(function() {
    'use strict';
    angular.module('carton.services', [])
        .service('UserService', [function() {
            var sdo = {
                isLogged: false,
                username: ''
            };

            return sdo;
        }]);
})();
