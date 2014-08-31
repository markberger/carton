// Directives

(function () {
	'use strict';
	angular.module('carton.directives', []).

	directive('selectOnClick', function () {
	    return {
	        restrict: 'A',
	        link: function (scope, element, attrs) {
	            element.on('click', function () {
	                this.select();
	            });
	        }
	    };
	});
})();
