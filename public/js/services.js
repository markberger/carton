// Services

(function() {
    'use strict';
    angular.module('carton.services', [])
        .service('UserService', ['$http',
			function($http) {

				var userSrv = this;
				userSrv.isLogged = false;

				$http({method: 'GET', url: '/api/auth/status'})

				.success(function(data, status, headers, config) {
					userSrv.isLogged = true;
				});
        }]);
})();
