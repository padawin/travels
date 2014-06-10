(function(){
	var travelsApp = angular.module('travelsApp', ['ngRoute']),
		travels,
		getTravels;

	travelsApp.config(function($routeProvider, $locationProvider) {
		$locationProvider.hashPrefix('');
		$routeProvider
			.when("/", {
				templateUrl: "partials/travels-list.html",
				controller: "TravelsListCtrl"
			})
			.otherwise({redirectTo: "" });
	});

	getTravels = function($scope, $http, doneCallback)
	{
		if (travels != null) {
			$scope.travels = travels;
			doneCallback && doneCallback($scope);
		}
		else {
			$http.get('data/travels.json').success(function(data) {
				travels = data;
				$scope.travels = data;
				doneCallback && doneCallback($scope);
			});
		}
	};

	travelsApp.controller('TravelsListCtrl', function($scope, $http){
		$scope.orderProp = 'title';
		getTravels($scope, $http, function($scope){
			var t, pics;
			for (t in $scope.travels) {
				$scope.travels[t].id = t;
				$scope.travels[t].previewPics = [].concat.apply([], $scope.travels[t].pics)
					// duplicate initial pics list
					.slice(0)
					// randomize
					.sort(function() { return 0.5 - Math.random() })
					// takes the 4 first
					.splice(0, 4);
			}
		});
	});

	window.travelsApp = travelsApp;
})();
