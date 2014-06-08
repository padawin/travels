(function(){
var travelsApp = angular.module('travelsApp', []);

	travelsApp.controller('TravelsListCtrl', function($scope, $http){
		$http.get('data/travels.json').success(function(data) {
			$scope.travels = data;
			$scope.orderProp = 'title';
		});
	});

	window.travelsApp = travelsApp;
})();
