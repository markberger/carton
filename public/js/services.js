// Services

(function() {
    'use strict';
    angular.module('carton.services', [])
        .service('UserService', ['$http',
			function($http) {

				var sdo = {
					isLogged: false
				}

				$http({method: 'GET', url: '/api/auth/status'})

				.success(function(data, status, headers, config) {
					if (data.status) {
						sdo.isLogged = true;
							return sdo;
					} else {
						return sdo;
					}
				})

				.error(function(data, status, headers, config) {
					return sdo;
				});
        }]);
})();
